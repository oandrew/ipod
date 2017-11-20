package main

import (
	"flag"
	"fmt"
	"io"

	"os"

	"github.com/sirupsen/logrus"

	"git.andrewo.pw/andrew/ipod/lingo-extremote"
	"git.andrewo.pw/andrew/ipod/lingo-simpleremote"

	"git.andrewo.pw/andrew/ipod"
	"git.andrewo.pw/andrew/ipod/hid"
	"git.andrewo.pw/andrew/ipod/lingo-audio"
	_ "git.andrewo.pw/andrew/ipod/lingo-extremote"
	"git.andrewo.pw/andrew/ipod/lingo-general"
	_ "git.andrewo.pw/andrew/ipod/lingo-simpleremote"
)

var devicePath = flag.String("d", "", "iap device")
var readTracePath = flag.String("r", "", "Read traces from a file instead of device")
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
		return NewLoadTraceReadWriter(f)

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
	log.Formatter = &logrus.TextFormatter{}

	log.Debugf("Registered lingos:\n%s", ipod.DumpLingos())

	f := open()
	if *writeTracePath != "" {
		traceFile, err := os.OpenFile(*writeTracePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.WithError(err).Fatal("Couldn't open the save trace file")
		}
		f = &tracingReadWriter{
			r:     f,
			w:     f,
			trace: traceFile,
		}
	}
	hidReportTransport := hid.NewCharDevReportTransport(f)
	frameTransport := hid.NewTransport(hidReportTransport, hid.DefaultReportDefs)
	logFrameTransport := &ipod.LoggingFrameReadWriter{
		RW: frameTransport,
		L:  log,
	}

	packetTransport := ipod.NewPacketTransport(logFrameTransport)
	packetRW := &ipod.LoggingPacketReadWriter{
		RW: packetTransport,
		L:  log,
	}

	devGeneral := &DevGeneral{}

	for {
		packet, err := packetRW.ReadPacket()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Error(err)
			}
		}

		switch packet.ID.LingoID() {
		case general.LingoGeneralID:
			general.HandleGeneral(packet, packetRW, devGeneral)
			if _, ok := packet.Payload.(general.RetDevAuthenticationSignature); ok {
				audio.Start(packetRW)
			}
		case simpleremote.LingoSimpleRemotelID:
			//todo
		case extremote.LingoExtRemotelID:
			extremote.HandleExtRemote(packet, packetRW, nil)
		case audio.LingoAudioID:
			audio.HandleAudio(packet, packetRW, nil)
		}

	}
}
