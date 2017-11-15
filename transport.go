package ipod

import (
	"bytes"
	"io"
	"log"
)

type Transport interface {
	Send(p *Packet)
	MaxPayload() uint16
}

type TransportEncoder interface {
	WriteFrame(data []byte) error
}

type TransportDecoder interface {
	ReadFrame() ([]byte, error)
}

type transportPacketReader struct {
	tr TransportDecoder
	r  *bytes.Reader
}

func (tpr *transportPacketReader) ReadPacket() (Packet, error) {

	for {
		if tpr.r == nil || tpr.r.Len() == 0 {
			frame, err := tpr.tr.ReadFrame()
			if err != nil {
				return Packet{}, err
			}
			log.Printf("frame: [% 02x]", frame)
			tpr.r = bytes.NewReader(frame)
		}
		log.Printf("leftover: %d", tpr.r.Len())
		var pkt Packet
		err := UnmarshalPacket(tpr.r, &pkt)
		if err == io.EOF {
			continue
		}
		return pkt, err
	}

}

func NewTransportPacketReader(tr TransportDecoder) PacketReader {
	return &transportPacketReader{
		tr: tr,
	}

}

type transportPacketWriter struct {
	tr TransportEncoder
}

func (tpw *transportPacketWriter) WritePacket(pkt Packet) error {
	buf := bytes.Buffer{}
	if err := MarshalPacket(&buf, &pkt); err != nil {
		return err
	}
	return tpw.tr.WriteFrame(buf.Bytes())

}

func NewTransportPacketWriter(tr TransportEncoder) PacketWriter {
	return &transportPacketWriter{
		tr: tr,
	}
}
