package main

import (
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

		packetReader := ipod.NewPacketReader(inFrame)
		for {
			inPacket, err := packetReader.ReadPacket()
			inPktLogE := PacketLogEntry(logrus.NewEntry(log), inPacket)

			if err != nil {
				if err == io.EOF {
					break
				}
				inPktLogE.WithError(err).Errorf("<< PACKET READ ERROR")
				continue
			}
			inPktLogE.Infof("<< PACKET READ")
			if log.Level == logrus.DebugLevel {
				log.Debug(spew.Sdump(inPacket))
			}

			packetBuf := ipod.PacketBuffer{}
			//todo: check return error
			handlePacket(&packetBuf, inPacket)
			if len(packetBuf.Packets) == 0 {
				continue
			}
			frameBuilder := ipod.NewFrameBuilder()
			for i := range packetBuf.Packets {
				outPacket := packetBuf.Packets[i]
				outPktLogE := PacketLogEntry(logrus.NewEntry(log), outPacket)
				if err := frameBuilder.WritePacket(outPacket); err != nil {
					outPktLogE.WithError(err).Errorf(">> PACKET WRITE ERROR")
				}
				outPktLogE.Infof(">> PACKET WRITE")
				if log.Level == logrus.DebugLevel {
					log.Debug(spew.Sdump(outPacket.Payload))
				}
			}
			outFrame := frameBuilder.Frame()
			if len(outFrame) > 0 {
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

func handlePacket(packetRW ipod.PacketWriter, packet ipod.Packet) {
	switch packet.ID.LingoID() {
	case general.LingoGeneralID:
		general.HandleGeneral(packet, packetRW, devGeneral)
		if _, ok := packet.Payload.(general.RetDevAuthenticationSignature); ok {
			audio.Start(packetRW)
		}
	case simpleremote.LingoSimpleRemotelID:
		//todo
		log.Warn("Lingo SimpleRemote is not supported yet")
	case dispremote.LingoDisplayRemoteID:
		dispremote.HandleDispRemote(packet, packetRW, nil)
	case extremote.LingoExtRemotelID:
		extremote.HandleExtRemote(packet, packetRW, nil)
	case audio.LingoAudioID:
		audio.HandleAudio(packet, packetRW, nil)
	}
}
