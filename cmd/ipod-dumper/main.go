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

	reportReader := hid.NewRawReportReader(f)
	hidDecoder := hid.NewDecoderDefault(reportReader)
	packetReader := ipod.NewTransportPacketReader(hidDecoder)

	reportWriter := hid.NewRawReportWriter(f)
	hidEncoder := hid.NewEncoderDefault(reportWriter)
	packetWriter := ipod.NewTransportPacketWriter(hidEncoder)

	devGeneral := &DevGeneral{}

	packetRW := ipod.NewLoggingPacketReadWriter(packetReader, packetWriter, os.Stderr)

	errCnt := 0
	for {
		//var fdr syscall.FdSet
		//syscall.Select(int(f.Fd()+1), &fdr, nil, nil, &syscall.Timeval{Sec: 1})
		packet, err := packetRW.ReadPacket()
		if err != nil {
			log.Printf("Error: %v", err)
			//time.Sleep(10 * time.Millisecond)
			errCnt++
		} else {
			errCnt = 0
		}

		if errCnt == 50 {
			break
		}
		//log.Printf("err: %v packet go-syntax: %#v", err, packet)
		//log.Printf("packet stringer: %v", packet)
		// buf := bytes.Buffer{}
		// ipod.NewRawPacketWriter(&buf).WritePacket(packet)
		// log.Printf("encoded again: [% 02x]", buf.Bytes())

		general.HandleGeneral(packet, packetRW, devGeneral)

	}
}
