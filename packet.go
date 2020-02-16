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
	frame []byte
}

func NewPacketReader(frame []byte) *PacketReader {
	return &PacketReader{
		frame: frame,
	}
}

func parseHeader(data []byte) (payOff, payLen int, err error) {
	if len(data) < 3 {
		err = io.ErrUnexpectedEOF
		return
	}
	if data[0] == 0x00 {
		payOff = 3
		payLen = int(binary.BigEndian.Uint16(data[1:3]))
	} else {
		payOff = 1
		payLen = int(data[0])
	}
	return
}

func parsePacket(data []byte) (int, []byte, error) {
	payOff, payLen, err := parseHeader(data)
	if err != nil {
		return 0, nil, err
	}
	pktLen := payOff + payLen + 1
	pkt := data[:pktLen]
	if len(pkt) < pktLen {
		return pktLen, nil, io.ErrUnexpectedEOF
	}
	if Checksum(pkt) != 0x00 {
		return pktLen, nil, errors.New("invalid checksum")
	}

	return pktLen, pkt[payOff : payOff+payLen], nil
}

func (pd *PacketReader) ReadPacket() ([]byte, error) {
	next := bytes.IndexByte(pd.frame, PacketStartByte)
	if next == -1 {
		return nil, io.EOF
	}
	pd.frame = pd.frame[next+1:]
	pktLen, payload, err := parsePacket(pd.frame)
	if err != nil {
		return nil, err
	}
	pd.frame = pd.frame[pktLen:]
	return payload, nil
}

type PacketWriter struct {
	frame []byte
}

func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		frame: make([]byte, 0, 512),
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
