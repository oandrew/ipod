package ipod

import (
	"bytes"
	"io"
)

type FrameReader interface {
	// ReadFrame reads a frame that contains
	// one or more iap packets
	ReadFrame() ([]byte, error)
}

type FrameWriter interface {
	// WriteFrame writes a frame that contains
	// one or more iap packets
	WriteFrame(data []byte) error
}

// FrameReadWriter is the interface
// implemented by iap transports i.e. usbhid
type FrameReadWriter interface {
	FrameReader
	FrameWriter
}

// DummyFrameReadWriter is a no-op implementation of FrameReadWriter
type DummyFrameReadWriter struct{}

func (d *DummyFrameReadWriter) ReadFrame() ([]byte, error) {
	return []byte{}, nil
}

func (d *DummyFrameReadWriter) WriteFrame([]byte) error {
	return nil
}

type packetReader struct {
	r *bytes.Reader
}

func (pr *packetReader) ReadPacket() (*Packet, error) {
	if pr.r.Len() == 0 {
		return nil, io.EOF
	}

	var pkt Packet
	err := UnmarshalPacket(pr.r, &pkt)
	return &pkt, err
}

func NewPacketReader(frame []byte) PacketReader {
	return &packetReader{
		r: bytes.NewReader(frame),
	}

}

type PacketBuffer struct {
	Packets []*Packet
}

func (pb *PacketBuffer) WritePacket(pkt *Packet) error {
	pb.Packets = append(pb.Packets, pkt)
	return nil
}

type FrameBuilder struct {
	buf *bytes.Buffer
}

func (fb *FrameBuilder) WritePacket(pkt *Packet) error {
	return MarshalPacket(fb.buf, pkt)
}

func (fb *FrameBuilder) Frame() []byte {
	return fb.buf.Bytes()
}

func (fb *FrameBuilder) Reset() {
	fb.buf.Reset()
}

func NewFrameBuilder() *FrameBuilder {
	return &FrameBuilder{
		buf: bytes.NewBuffer(make([]byte, 0, 1024)),
	}
}
