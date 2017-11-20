package ipod

import (
	"bytes"
	"io"
)

type FrameReader interface {
	ReadFrame() ([]byte, error)
}

type FrameWriter interface {
	WriteFrame(data []byte) error
}

type FrameReadWriter interface {
	FrameReader
	FrameWriter
}

// type Transport struct {
// 	TransportReader
// 	TransportWriter
// }

type DummyFrameReadWriter struct{}

func (d *DummyFrameReadWriter) ReadFrame() ([]byte, error) {
	return []byte{}, nil
}

func (d *DummyFrameReadWriter) WriteFrame([]byte) error {
	return nil
}

type transportPacketReader struct {
	tr FrameReader
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

func NewPacketReader(tr FrameReader) PacketReader {
	return &transportPacketReader{
		tr: tr,
	}

}

type transportPacketWriter struct {
	tr FrameWriter
}

func (tpw *transportPacketWriter) WritePacket(pkt Packet) error {
	buf := bytes.Buffer{}
	if err := MarshalPacket(&buf, &pkt); err != nil {
		return err
	}
	return tpw.tr.WriteFrame(buf.Bytes())

}

func NewPacketWriter(tr FrameWriter) PacketWriter {
	return &transportPacketWriter{
		tr: tr,
	}
}

type packetRW struct {
	PacketReader
	PacketWriter
}

func NewPacketTransport(tr FrameReadWriter) PacketReadWriter {
	return &packetRW{
		PacketReader: NewPacketReader(tr),
		PacketWriter: NewPacketWriter(tr),
	}

}
