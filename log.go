package ipod

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/sirupsen/logrus"
)

type LoggingPacketReadWriter struct {
	RW PacketReadWriter
	L  *logrus.Logger
}

func packetLogEntry(e *logrus.Entry, pkt Packet) *logrus.Entry {
	return e.WithFields(logrus.Fields{
		"id":   pkt.ID,
		"trx":  pkt.Transaction,
		"type": fmt.Sprintf("%T", pkt.Payload),
	})
}

func (lprw *LoggingPacketReadWriter) ReadPacket() (pkt Packet, err error) {
	pkt, err = lprw.RW.ReadPacket()
	entry := packetLogEntry(logrus.NewEntry(lprw.L), pkt)
	if err != nil {
		entry.WithError(err).Errorf("recv error")
		return
	}
	entry.Debugf("packet recv\n%s", spew.Sdump(pkt.Payload))
	return
}

func (lprw *LoggingPacketReadWriter) WritePacket(pkt Packet) (err error) {
	err = lprw.RW.WritePacket(pkt)
	entry := packetLogEntry(logrus.NewEntry(lprw.L), pkt)
	if err != nil {
		entry.WithError(err).Errorf("send error")
		return
	}
	entry.Debugf("packet send\n%s", spew.Sdump(pkt.Payload))
	return
}
