package ipod

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
// type DummyFrameReadWriter struct{}

// func (d *DummyFrameReadWriter) ReadFrame() ([]byte, error) {
// 	return []byte{}, nil
// }

// func (d *DummyFrameReadWriter) WriteFrame([]byte) error {
// 	return nil
// }

// type packetReader struct {
// 	r *RawPacketReader
// }

// func (pr *packetReader) ReadPacket() (*Packet, error) {
// 	payload, err := pr.r.ReadPayload()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return UnmarshalPacket(payload)
// }

// func NewPacketReader(frame []byte) PacketReader {
// 	return &packetReader{
// 		r: NewRawPacketReader(bytes.NewReader(frame)),
// 	}

// }

// type FrameBuilder struct {
// 	buf *bytes.Buffer
// 	w   *RawPacketWriter
// }

// func (fb *FrameBuilder) WritePacket(pkt *Packet) error {
// 	payload, err := MarshalPacket(pkt)
// 	if err != nil {
// 		return err
// 	}
// 	return fb.w.WritePayload(payload)
// }

// func (fb *FrameBuilder) Frame() []byte {
// 	return fb.buf.Bytes()
// }

// func (fb *FrameBuilder) Reset() {
// 	fb.buf.Reset()
// }

// func NewFrameBuilder() *FrameBuilder {
// 	buf := bytes.NewBuffer(make([]byte, 0, 1024))
// 	return &FrameBuilder{
// 		buf: buf,
// 		w:   NewRawPacketWriter(buf),
// 	}
// }
