package ipod

import (
	"bytes"
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
	if tpr.r == nil || tpr.r.Len() == 0 {
		frame, err := tpr.tr.ReadFrame()
		if err != nil {
			return Packet{}, err
		}
		tpr.r = bytes.NewReader(frame)
	}

	var pkt Packet
	if err := UnmarshalPacket(tpr.r, &pkt); err != nil {
		return Packet{}, err
	}
	return pkt, nil

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
