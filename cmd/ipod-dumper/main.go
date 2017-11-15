package main

import (
	"bytes"
	"log"
	"os"

	"git.andrewo.pw/andrew/ipod"
	"git.andrewo.pw/andrew/ipod/hid"
	_ "git.andrewo.pw/andrew/ipod/lingo-general"
)

func main() {
	p := "/dev/iap0"
	if len(os.Args) > 1 {
		p = os.Args[1]
	}
	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}

	reportReader := hid.NewRawReportReader(f)
	hidDecoder := hid.NewDecoderDefault(reportReader)
	packetReader := ipod.NewTransportPacketReader(hidDecoder)
	for {
		// frame, err := hidDecoder.ReadFrame()
		// log.Printf("err: %v frame: [% 02x]", err, frame)
		// err2 := ipod.UnmarshalPacket(bytes.NewReader(frame), &pkt)
		packet, err := packetReader.ReadPacket()
		if err != nil {
			break
		}
		log.Printf("err: %v packet: %#v", err, packet)
		buf := bytes.Buffer{}
		ipod.NewRawPacketWriter(&buf).WritePacket(packet)
		log.Printf("encoded again: [% 02x]", buf.Bytes())

	}
}
