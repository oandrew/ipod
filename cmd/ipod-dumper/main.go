package main

import (
	"flag"
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

var devicePath = flag.String("device", "/dev/iap0", "iap char device path")

var log = logrus.StandardLogger()

func main() {
	flag.Parse()

	log.SetLevel(logrus.DebugLevel)
	log.Formatter = &logrus.TextFormatter{}

	log.Debugf("Registered lingos:\n%s", ipod.DumpLingos())

	dev, err := os.OpenFile(*devicePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.WithError(err).Fatalf("Coult not open device %s", *devicePath)
	}

	log.Infof("Device %s opened", *devicePath)

	reportTransport := hid.NewCharDevReportTransport(dev)
	rw := ipod.NewPacketReadWriter(&ipod.Transport{
		TransportReader: hid.NewDecoder(reportTransport, hid.DefaultReportDefs),
		TransportWriter: hid.NewEncoder(reportTransport, hid.DefaultReportDefs),
	})
	packetRW := &ipod.LoggingPacketReadWriter{
		RW: rw,
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
