// Package ipod implements the iPod Accessory protocol (iap)
package ipod

import (
	"bufio"
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
	r io.Reader
}

func NewPacketReader(r io.Reader) *PacketReader {
	return &PacketReader{
		//r: bufio.NewReader(r),
		r: r,
	}
}

func (pd *PacketReader) ReadPacket() ([]byte, error) {
	var header [4]byte
	if _, err := pd.r.Read(header[:2]); err != nil {
		return nil, err
	}

	if header[0] == 0x00 && header[1] == 0x00 {
		return nil, io.EOF
	}

	if header[0] == 0xff && header[1] == PacketStartByte {
		header[0] = header[1]
		if _, err := pd.r.Read(header[1:2]); err != nil {
			return nil, err
		}
	}

	if header[0] != PacketStartByte {
		return nil, errors.New("packet decode: start byte not found")
	}

	var payLen int
	if header[1] == 0x00 {
		if _, err := pd.r.Read(header[2:4]); err != nil {
			return nil, err
		}
		payLen = int(binary.BigEndian.Uint16(header[2:4]))
	} else {
		payLen = int(header[1])
	}

	payloadWithCrc := make([]byte, payLen+1)
	if n, err := pd.r.Read(payloadWithCrc); n != len(payloadWithCrc) || err != nil {
		return nil, errors.New("packet decode: short read")
	}

	payload := payloadWithCrc[:payLen]
	crc := payloadWithCrc[payLen]
	crc8 := NewCRC8()
	crc8.Write(header[1:])
	crc8.Write(payload)
	calcCrc := crc8.Sum8()
	if crc != calcCrc {
		return nil, fmt.Errorf("packet decode: crc mismatch: recv %02x != calc %02x", crc, calcCrc)
	}
	return payload, nil

}

type PacketWriter struct {
	w *bufio.Writer
}

func NewPacketWriter(w io.Writer) *PacketWriter {
	return &PacketWriter{
		w: bufio.NewWriterSize(w, minPacketBufSize),
	}
}

func (pw *PacketWriter) WritePacket(pkt []byte) error {
	if len(pkt) == 0 {
		return fmt.Errorf("packet encode: empty packet")
	}
	binary.Write(pw.w, binary.BigEndian, PacketStartByte)

	crc := NewCRC8()
	mw := io.MultiWriter(pw.w, crc)
	if len(pkt) > largePacketMinLen {
		binary.Write(mw, binary.BigEndian, byte(0x00))
		binary.Write(mw, binary.BigEndian, uint16(len(pkt)))
	} else {
		binary.Write(mw, binary.BigEndian, byte(len(pkt)))
	}
	mw.Write(pkt)

	pw.w.Write([]byte{crc.Sum8()})

	return pw.w.Flush()
}
