package ipod

import (
	"bytes"
	"io"
)

type TransportReader interface {
	ReadFrame() ([]byte, error)
}

type TransportWriter interface {
	WriteFrame(data []byte) error
}

type Transport struct {
	TransportReader
	TransportWriter
}

type DummyTransport struct{}

func (d *DummyTransport) ReadFrame() ([]byte, error) {
	return []byte{}, nil
}

func (d *DummyTransport) WriteFrame([]byte) error {
	return nil
}

type transportPacketReader struct {
	tr TransportReader
	r  *bytes.Reader
}

func (tpr *transportPacketReader) ReadPacket() (Packet, error) {

	for {
		if tpr.r == nil || tpr.r.Len() == 0 {
			frame, err := tpr.tr.ReadFrame()
			if err != nil {
				return Packet{}, err
			}
			//log.Printf("frame: [% 02x]", frame)
			tpr.r = bytes.NewReader(frame)
		}
		//log.Printf("leftover: %d", tpr.r.Len())
		var pkt Packet
		err := UnmarshalPacket(tpr.r, &pkt)
		if err == io.EOF {
			continue
		}
		return pkt, err
	}

}

func NewPacketReader(tr TransportReader) PacketReader {
	return &transportPacketReader{
		tr: tr,
	}

}

type transportPacketWriter struct {
	tr TransportWriter
}

func (tpw *transportPacketWriter) WritePacket(pkt Packet) error {
	buf := bytes.Buffer{}
	if err := MarshalPacket(&buf, &pkt); err != nil {
		return err
	}
	return tpw.tr.WriteFrame(buf.Bytes())

}

func NewPacketWriter(tr TransportWriter) PacketWriter {
	return &transportPacketWriter{
		tr: tr,
	}
}

type packetRW struct {
	PacketReader
	PacketWriter
}

func NewPacketReadWriter(tr *Transport) PacketReadWriter {
	return packetRW{
		NewPacketReader(tr),
		NewPacketWriter(tr),
	}

}
