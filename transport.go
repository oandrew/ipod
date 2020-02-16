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
