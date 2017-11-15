package ipod

import (
	"hash"
)

type CRC8 interface {
	hash.Hash
	Sum8() uint8
}

type crc8 struct {
	crc byte
}

func (c *crc8) Write(p []byte) (n int, err error) {
	for _, v := range p {
		c.crc += v
	}
	return len(p), nil
}

func (c *crc8) Sum8() byte { return -c.crc }

func (c *crc8) Sum(in []byte) []byte {
	return append(in, c.Sum8())
}
func (c *crc8) Reset()         { c.crc = 0x00 }
func (c *crc8) Size() int      { return 1 }
func (c *crc8) BlockSize() int { return 1 }

func NewCRC8() CRC8 {
	return &crc8{}
}

func Checksum(p []byte) uint8 {
	crc := NewCRC8()
	crc.Write(p)
	return crc.Sum8()

}
