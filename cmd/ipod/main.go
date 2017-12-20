package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/oandrew/ipod"
	"github.com/oandrew/ipod/hid"
	"github.com/oandrew/ipod/lingo-audio"
	"github.com/oandrew/ipod/lingo-dispremote"
	"github.com/oandrew/ipod/lingo-extremote"
	"github.com/oandrew/ipod/lingo-general"
	"github.com/oandrew/ipod/lingo-simpleremote"
	"github.com/oandrew/ipod/trace"
)

var devicePath = flag.String("d", "", "iap device i.e. /dev/iap0")
var readTracePath = flag.String("r", "", "Respond to requests  from a trace file instead of device")
var writeTracePath = flag.String("w", "", "Save traces to a file")

var verbose = flag.Bool("v", false, "Enable verbose logging")

var log = logrus.StandardLogger()

func open() io.ReadWriter {
	if *readTracePath != "" {
		f, err := os.Open(*readTracePath)
		e := log.WithField("path", *readTracePath)
		if err != nil {
			e.WithError(err).Fatalf("Couldn't open the trace file")
		}
		e.Warningf("Using trace file")
		t := trace.NewReader(f)
		return ReadWriter{NewTraceInputReader(t), ioutil.Discard}

	} else if *devicePath != "" {
		dev, err := os.OpenFile(*devicePath, os.O_RDWR, os.ModePerm)
		e := log.WithField("path", *devicePath)
		if err != nil {
			e.WithError(err).Fatalf("Couldn't not open the device")
		}
		stat, _ := dev.Stat()
		if stat.Mode()&os.ModeCharDevice != os.ModeCharDevice {
			e.Fatalf("Not a device")
		}
		e.Infof("Device opened")
		return dev
	}
	return nil
}

func main() {
	flag.Parse()
	if *readTracePath == "" && *devicePath == "" || *readTracePath != "" && *devicePath != "" {
		fmt.Fprintf(os.Stderr, "Specify either a device or a trace file\n\n")
		flag.Usage()
		os.Exit(2)

	}

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
	}
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	}

	log.Debugf("Registered lingos:\n%s", ipod.DumpLingos())

	f := open()
	if *writeTracePath != "" {
		traceFile, err := os.OpenFile(*writeTracePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.WithError(err).Fatal("Couldn't open the save trace file")
		}
		f = trace.NewTracer(traceFile, f)
	}
	hidReportTransport := hid.NewCharDevReportTransport(f)
	frameTransport := hid.NewTransport(hidReportTransport, hid.DefaultReportDefs)

	processFrames(frameTransport)

}

func processFrames(frameTransport ipod.FrameReadWriter) {
	for {
		inFrame, err := frameTransport.ReadFrame()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.WithError(err).Errorf("<< FRAME READ ERROR")
			continue
		}

		FrameLogEntry(logrus.NewEntry(log), inFrame).Infof("<< FRAME READ")
		if log.Level == logrus.DebugLevel {
			log.Debug(spew.Sdump(inFrame))
		}

		packetReader := ipod.NewPacketReader(bytes.NewReader(inFrame))
		for {
			inPacket, err := packetReader.ReadPacket()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.WithError(err).Errorf("<< PACKET READ ERROR")
				continue
			}
			log.Infof("<< PACKET READ")
			if log.Level == logrus.DebugLevel {
				log.Debug(spew.Sdump(inPacket))
			}

			var inCmd ipod.Command
			if err := inCmd.UnmarshalBinary(inPacket); err != nil {
				CommandLogEntry(logrus.NewEntry(log), &inCmd).WithError(err).Errorf("<< CMD DECODE ERROR")

			}
			CommandLogEntry(logrus.NewEntry(log), &inCmd).Infof("<< CMD READ")
			if log.Level == logrus.DebugLevel {
				log.Debug(spew.Sdump(inCmd))
			}

			cmdBuf := ipod.CmdBuffer{}
			//todo: check return error
			handlePacket(&cmdBuf, &inCmd)
			if len(cmdBuf.Commands) == 0 {
				continue
			}
			frameBuf := bytes.Buffer{}
			packetWriter := ipod.NewPacketWriter(&frameBuf)
			for i := range cmdBuf.Commands {
				outCmd := cmdBuf.Commands[i]
				CommandLogEntry(logrus.NewEntry(log), outCmd).Infof(">> CMD WRITE")
				if log.Level == logrus.DebugLevel {
					log.Debug(spew.Sdump(outCmd))
				}

				outPacket, err := outCmd.MarshalBinary()
				if err != nil {
					log.WithError(err).Errorf(">> CMD ENCODE ERROR")
					continue
				}

				log.Infof(">> PACKET WRITE")
				if log.Level == logrus.DebugLevel {
					log.Debug(spew.Sdump(outPacket))
				}
				packetWriter.WritePacket(outPacket)
			}

			if frameBuf.Len() > 0 {
				outFrame := frameBuf.Bytes()
				outFrameLogE := FrameLogEntry(logrus.NewEntry(log), outFrame)
				if err := frameTransport.WriteFrame(outFrame); err != nil {
					outFrameLogE.WithError(err).Errorf(">> FRAME WRITE ERROR")
				}
				outFrameLogE.Infof(">> FRAME WRITE")
				if log.Level == logrus.DebugLevel {
					log.Debug(spew.Sdump(outFrame))
				}

			}

		}
	}
	log.Warnf("EOF")
}

var devGeneral = &DevGeneral{}

func handlePacket(cmdWriter ipod.CommandWriter, cmd *ipod.Command) {
	switch cmd.ID.LingoID() {
	case general.LingoGeneralID:
		general.HandleGeneral(cmd, cmdWriter, devGeneral)
		if _, ok := cmd.Payload.(general.RetDevAuthenticationSignature); ok {
			audio.Start(cmdWriter)
		}
	case simpleremote.LingoSimpleRemotelID:
		//todo
		log.Warn("Lingo SimpleRemote is not supported yet")
	case dispremote.LingoDisplayRemoteID:
		dispremote.HandleDispRemote(cmd, cmdWriter, nil)
	case extremote.LingoExtRemotelID:
		extremote.HandleExtRemote(cmd, cmdWriter, nil)
	case audio.LingoAudioID:
		audio.HandleAudio(cmd, cmdWriter, nil)
	}
}
