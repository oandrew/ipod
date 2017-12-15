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
	fr FrameReader
	r  *bytes.Reader
}

func (pr *packetReader) ReadPacket() (Packet, error) {

	for {
		if pr.r == nil || pr.r.Len() == 0 {
			frame, err := pr.fr.ReadFrame()
			if err != nil {
				return Packet{}, err
			}
			//log.Printf("frame: [% 02x]", frame)
			pr.r = bytes.NewReader(frame)
		}
		//log.Printf("leftover: %d", tpr.r.Len())
		var pkt Packet
		err := UnmarshalPacket(pr.r, &pkt)
		if err == io.EOF {
			continue
		}
		return pkt, err
	}

}

func NewPacketReader(fr FrameReader) PacketReader {
	return &packetReader{
		fr: fr,
	}

}

type packetWriter struct {
	fw FrameWriter
}

func (pw *packetWriter) WritePacket(pkt Packet) error {
	buf := bytes.Buffer{}
	if err := MarshalPacket(&buf, &pkt); err != nil {
		return err
	}
	return pw.fw.WriteFrame(buf.Bytes())

}

func NewPacketWriter(fw FrameWriter) PacketWriter {
	return &packetWriter{
		fw: fw,
	}
}

type BufferedPacketWriter struct {
	buf bytes.Buffer
	fw  FrameWriter
}

func (bpw *BufferedPacketWriter) WritePacket(pkt Packet) error {
	return MarshalPacket(&bpw.buf, &pkt)
}

func (bpw *BufferedPacketWriter) Flush() error {
	frame := bpw.buf.Bytes()
	bpw.buf.Reset()
	return bpw.fw.WriteFrame(frame)
}

func NewBufferedPacketWriter(fw FrameWriter) *BufferedPacketWriter {
	return &BufferedPacketWriter{
		fw: fw,
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
