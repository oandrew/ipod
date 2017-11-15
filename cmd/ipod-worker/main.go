package main

import (
	"io"
	"log"
	"os"

	"git.andrewo.pw/andrew/ipod"
	"git.andrewo.pw/andrew/ipod/proto"
)

type dev struct {
	uimode ipod.UIMode
}

func (d *dev) UIMode() ipod.UIMode {
	return d.uimode
}

func (d *dev) SetUIMode(mode ipod.UIMode) {
	d.uimode = mode
}

func (d *dev) Name() string {
	return "Andrew"
}

func (d *dev) SoftwareVersion() (major uint8, minor uint8, rev uint8) {
	return 1, 1, 1
}

func (d *dev) SerialNum() string {
	return "abcd1234"
}

func (d *dev) LingoProtocolVersion(lingo uint8) (major uint8, minor uint8) {
	return 1, 0
}

func (d *dev) LingoOptions(ling uint8) uint64 {
	return 0
}

func (d *dev) PrefSettingID(classID uint8) uint8 {
	return 0
}

func (d *dev) SetPrefSettingID(classID uint8, settingID uint8, restoreOnExit bool) {
}

func (d *dev) StartIDPS() {
}

func (d *dev) SetEventNotificationMask(mask uint64) {

}

func (d *dev) EventNotificationMask() uint64 {
	return 0
}

func (d *dev) SupportedEventNotificationMask() uint64 {
	return 0
}

func (d *dev) CancelCommand(lingo uint8, cmd uint16, transaction uint16) {

}

type transport struct {
	t proto.TransportEncoder
}

func (t *transport) Send(packet *ipod.Packet) {
	log.Printf("SEND: %#v", packet)
	if err := t.t.Encode(packet); err != nil {
		log.Printf("SEND error: %v", err)
	}
}

func (t *transport) MaxPayload() uint16 {
	return 0xfff9
}

func main() {
	p := "/dev/iap0"
	if len(os.Args) > 1 {
		p = os.Args[1]
	}
	f, err := os.OpenFile(p, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}

	hidDec := proto.NewHidDecoderDefault(f)
	hidEnc := proto.NewHidEncoderDefault(f)
	t := &transport{hidEnc}
	dev := &dev{}

	var pkt ipod.Packet

	for {
		err := hidDec.Decode(&pkt)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("err: %v", err)
			continue
		}
		log.Printf("Packet: %#v", pkt)

		ipod.HandleGeneral(pkt, t, dev)
	}
}
