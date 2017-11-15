package ipod

import (
	"io"
	"log"
)

type loggingPacketReadWriter struct {
	r      PacketReader
	w      PacketWriter
	logger *log.Logger
}

func (lprw *loggingPacketReadWriter) ReadPacket() (Packet, error) {
	pkt, err := lprw.r.ReadPacket()
	lprw.logger.Printf("RECV %#+v", pkt)
	if err != nil {
		lprw.logger.Printf("RECV ERR: %v", err)
	}
	return pkt, err
}

func (lprw *loggingPacketReadWriter) WritePacket(pkt Packet) error {
	lprw.logger.Printf("SEND %#+v", pkt)
	err := lprw.w.WritePacket(pkt)
	if err != nil {
		lprw.logger.Printf("SEND ERR: %v", err)
	}
	return err
}

func NewLoggingPacketReadWriter(r PacketReader, w PacketWriter, logw io.Writer) PacketReadWriter {
	return &loggingPacketReadWriter{
		r:      r,
		w:      w,
		logger: log.New(logw, "log>", log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}
