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
		entry.WithError(err).Errorf("PACKET READ ERROR")
		return
	}
	entry.Infof("<< PACKET READ")
	if lprw.L.Level == logrus.DebugLevel {
		lprw.L.Debug(spew.Sdump(pkt.Payload))
	}

	return
}

func (lprw *LoggingPacketReadWriter) WritePacket(pkt Packet) (err error) {
	err = lprw.RW.WritePacket(pkt)
	entry := packetLogEntry(logrus.NewEntry(lprw.L), pkt)
	if err != nil {
		entry.WithError(err).Errorf("PACKET WRITE ERROR")
		return
	}
	entry.Infof(">> PACKET WRITE")
	if lprw.L.Level == logrus.DebugLevel {
		lprw.L.Debug(spew.Sdump(pkt.Payload))
	}
	return
}

type LoggingFrameReadWriter struct {
	RW FrameReadWriter
	L  *logrus.Logger
}

func frameLogEntry(e *logrus.Entry, data []byte) *logrus.Entry {
	return e.WithFields(logrus.Fields{
		"len": len(data),
	})
}

func (lfrw *LoggingFrameReadWriter) ReadFrame() (data []byte, err error) {
	data, err = lfrw.RW.ReadFrame()
	entry := frameLogEntry(logrus.NewEntry(lfrw.L), data)
	if err != nil {
		entry.WithError(err).Errorf("FRAME READ ERROR")
		return
	}
	entry.Infof("<< FRAME READ")
	if lfrw.L.Level == logrus.DebugLevel {
		lfrw.L.Debug(spew.Sdump(data))
	}

	return
}

func (lfrw *LoggingFrameReadWriter) WriteFrame(data []byte) (err error) {
	err = lfrw.RW.WriteFrame(data)
	entry := frameLogEntry(logrus.NewEntry(lfrw.L), data)
	if err != nil {
		entry.WithError(err).Errorf("FRAME WRITE ERROR")
		return
	}
	entry.Infof(">> FRAME WRITE")
	if lfrw.L.Level == logrus.DebugLevel {
		lfrw.L.Debug(spew.Sdump(data))
	}
	return
}
