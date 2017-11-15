package main

import (
	"log"
	"os"

	"git.andrewo.pw/andrew/ipod"
	"git.andrewo.pw/andrew/ipod/hid"
	"git.andrewo.pw/andrew/ipod/lingo-general"
)

func main() {
	p := "/dev/iap0"
	if len(os.Args) > 1 {
		p = os.Args[1]
	}
	f, err := os.OpenFile(p, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	reportReader, reportWriter := hid.NewRawReportReader(f), hid.NewRawReportWriter(f)
	hidDecoder, hidEncoder := hid.NewDecoderDefault(reportReader), hid.NewEncoderDefault(reportWriter)
	packetReader, packetWriter := ipod.NewTransportPacketReader(hidDecoder), ipod.NewTransportPacketWriter(hidEncoder)

	devGeneral := &DevGeneral{}

	packetRW := ipod.NewLoggingPacketReadWriter(packetReader, packetWriter, os.Stderr)

	for {
		packet, err := packetRW.ReadPacket()
		if err != nil {
			log.Printf("Error: %v", err)
		}
		//log.Printf("err: %v packet go-syntax: %#v", err, packet)
		//log.Printf("packet stringer: %v", packet)
		// buf := bytes.Buffer{}
		// ipod.NewRawPacketWriter(&buf).WritePacket(packet)
		// log.Printf("encoded again: [% 02x]", buf.Bytes())

		general.HandleGeneral(packet, packetRW, devGeneral)

	}
}
