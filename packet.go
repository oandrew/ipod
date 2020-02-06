// Package ipod implements the iPod Accessory protocol (iap)
package ipod

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	PacketStartByte byte = 0x55
)
const (
	rawSmallPacketMinLen = 1 + 1 + 2 // start + len + ids
	rawLargePacketMinLen = 1 + 3 + 2 // start + len + ids
	largePacketMinLen    = 256
	minPacketBufSize     = 1024
)

type PacketReader struct {
	//r *bufio.Reader
	//r io.Reader
	frame []byte
}

func NewPacketReader(frame []byte) *PacketReader {
	return &PacketReader{
		//r: bufio.NewReaderSize(r, 512),
		//r: r,
		frame: frame,
	}
}

func parseHeader(data []byte) (payloadOffset, payloadLen int) {
	if data[0] == 0x00 {
		payloadOffset = 3
		payloadLen = int(binary.BigEndian.Uint16(data[1:3]))
	} else {
		payloadOffset = 1
		payloadLen = int(data[0])
	}
	return
}

//func parse(data []byte) (pkt, payload []byte, err error) {

// if data[0] == 0x00 {
// 	payloadLen := binary.BigEndian.Uint16(data[1:3])
// 	pkt = data[:3+payloadLen+1]
// 	payload = data[3 : 3+payloadLen]
// 	if len(payload) != int(payloadLen) {
// 		err = io.ErrUnexpectedEOF
// 	}
// } else {
// 	payloadLen := data[0]
// 	pkt = data[:1+payloadLen+1]
// 	payload = data[1 : 1+payloadLen]
// 	if len(payload) != int(payloadLen) {
// 		err = io.ErrUnexpectedEOF
// 	}
// }
//}

func (pd *PacketReader) ReadPacket() ([]byte, error) {
	next := bytes.IndexByte(pd.frame, PacketStartByte)
	if next == -1 {
		return nil, io.EOF
	}
	pd.frame = pd.frame[next+1:]
	payOff, payLen := parseHeader(pd.frame)
	pktLen := payOff + payLen + 1
	pkt := pd.frame[:pktLen]
	if len(pd.frame) < pktLen {
		return nil, io.ErrUnexpectedEOF
	}
	pd.frame = pd.frame[pktLen:]

	if Checksum(pkt) != 0x00 {
		return nil, errors.New("invalid checksum")
	}

	return pkt[payOff : payOff+payLen], nil
}

type PacketWriter struct {
	//w io.Writer
	frame []byte
}

func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		//w: w,
	}
}

func (pw *PacketWriter) WritePacket(payload []byte) error {
	if len(payload) == 0 {
		return fmt.Errorf("packet encode: empty packet")
	}

	pw.frame = append(pw.frame, PacketStartByte)
	pktStart := len(pw.frame)

	if len(payload) > largePacketMinLen {
		var pktLen [3]byte
		binary.BigEndian.PutUint16(pktLen[1:], uint16(len(payload)))
		pw.frame = append(pw.frame, pktLen[:]...)
	} else {
		pw.frame = append(pw.frame, byte(len(payload)))
	}

	pw.frame = append(pw.frame, payload...)
	pw.frame = append(pw.frame, Checksum(pw.frame[pktStart:]))
	return nil
}

func (pw *PacketWriter) Bytes() []byte {
	return pw.frame
}
