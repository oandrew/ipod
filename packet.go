// Package ipod implements the iPod Accessory protocol (iap)
package ipod

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Packet is a decoded iap packet
type Packet struct {
	ID          LingoCmdID
	Transaction *Transaction
	Payload     interface{}
}

func (p Packet) WithTransaction(t uint16) *Packet {
	tr := Transaction(t)
	p.Transaction = &tr
	return &p
}

func BuildPacket(payload interface{}) (*Packet, error) {
	id, ok := LookupID(payload)
	if !ok {
		return nil, errors.New("payload not known")
	}
	return &Packet{
		ID:      id,
		Payload: payload,
	}, nil

}

func Respond(req *Packet, pw PacketWriter, payload interface{}) {
	p, err := BuildPacket(payload)
	if err != nil {
		return
	}
	p.Transaction = req.Transaction
	pw.WritePacket(p)
}

// UnknownPayload is a payload  that represents an unknown command
type UnknownPayload []byte

type Transaction uint16

func (tr Transaction) GoString() string {
	return fmt.Sprintf("%#04x", tr)
}

func (tr Transaction) String() string {
	return fmt.Sprintf("%#04x", uint16(tr))
}

// PayloadUnmarshaler is the interface implemented by a payload
// that can unmarshal itself
type PayloadUnmarshaler interface {
	UnmarshalPayload(r io.Reader) error
}

// PayloadMarshaler is the interface implemented by a payload
// that can marshal itself
type PayloadMarshaler interface {
	MarshalPayload(w io.Writer) error
}

type PacketReader interface {
	ReadPacket() (*Packet, error)
}

type PacketWriter interface {
	WritePacket(*Packet) error
}

type PacketReadWriter interface {
	PacketReader
	PacketWriter
}

// RawPacket is an iap packet with encoded payload
type RawPacket struct {
	ID   LingoCmdID
	Data PacketPayload
}

func (p *RawPacket) Length() int {
	return p.ID.len() + len(p.Data)
}

type PacketPayload []byte

// func (pp PacketPayload) String() string {
// 	return fmt.Sprintf("(%d)[% 02x]", len([]byte(pp)), []byte(pp))
// }

type PayloadStringer struct {
	P interface{}
}

func (ps PayloadStringer) String() string {
	return fmt.Sprintf("%+v", ps.P)
}

const (
	PacketStartByte byte = 0x55
)
const (
	rawSmallPacketMinLen = 1 + 1 + 2 // start + len + ids
	rawLargePacketMinLen = 1 + 3 + 2 // start + len + ids
	largePacketMinLen    = 256
)

func binWrite(w io.Writer, v interface{}) error {
	if err := binary.Write(w, binary.BigEndian, v); err != nil {
		panic(err)
	}
	return nil
}

func binRead(r io.Reader, v interface{}) error {
	if err := binary.Read(r, binary.BigEndian, v); err != nil {
		panic(err)
	}
	return nil
}

func catchPanicErr(prefix string, dst *error) {
	if r := recover(); r != nil {
		*dst = fmt.Errorf("%s%v", prefix, r)
	}
}

type RawPacketReader struct {
	r *bufio.Reader
}

func NewRawPacketReader(r io.Reader) *RawPacketReader {
	return &RawPacketReader{
		r: bufio.NewReader(r),
	}
}

func (pd *RawPacketReader) ReadPayload() (PacketPayload, error) {
	//r := bytes.NewReader(d)
	var header [4]byte
	if _, err := pd.r.Read(header[:2]); err != nil {
		return nil, err
	}

	if header[0] == 0x00 && header[1] == 0x00 {
		return nil, io.EOF
	}

	if header[0] != PacketStartByte {
		return nil, errors.New("raw packet unmarshal: start byte not found")
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
		return nil, errors.New("raw packet unmarshal: short read")
	}

	payload := payloadWithCrc[:payLen]
	crc := payloadWithCrc[payLen]
	crc8 := NewCRC8()
	crc8.Write(header[1:])
	crc8.Write(payload)
	calcCrc := crc8.Sum8()
	if crc != calcCrc {
		return nil, fmt.Errorf("small packet: crc mismatch: recv %02x != calc %02x", crc, calcCrc)
	}

	// var id LingoCmdID
	// unmarshalLingoCmdID(payloadBuf, &id)
	return payload, nil

}

type RawPacketWriter struct {
	w *bufio.Writer
}

func NewRawPacketWriter(w io.Writer) *RawPacketWriter {
	return &RawPacketWriter{
		w: bufio.NewWriter(w),
	}
}

func (pe *RawPacketWriter) WritePayload(payload PacketPayload) error {

	binary.Write(pe.w, binary.BigEndian, PacketStartByte)

	crc := NewCRC8()
	mw := io.MultiWriter(pe.w, crc)
	if len(payload) > largePacketMinLen {
		binary.Write(mw, binary.BigEndian, byte(0x00))
		binary.Write(mw, binary.BigEndian, uint16(len(payload)))
	} else {
		binary.Write(mw, binary.BigEndian, byte(len(payload)))
	}
	mw.Write(payload)

	pe.w.Write([]byte{crc.Sum8()})

	return pe.w.Flush()
}

func MarshalPacket(p *Packet) (payload PacketPayload, err error) {
	defer catchPanicErr(" Packet marshal: ", &err)

	payloadBuf := bytes.Buffer{}

	marshalLingoCmdID(&payloadBuf, p.ID)

	if p.Transaction != nil {
		binWrite(&payloadBuf, *p.Transaction)
	}
	if d, ok := p.Payload.(PayloadMarshaler); ok {
		if err := d.MarshalPayload(&payloadBuf); err != nil {
			return nil, err
		}

	} else {
		binWrite(&payloadBuf, p.Payload)
	}

	return PacketPayload(payloadBuf.Bytes()), nil

}

func UnmarshalPacket(payload PacketPayload) (p *Packet, err error) {
	defer catchPanicErr(" Packet unmarshal: ", &err)

	p = &Packet{}

	r := bytes.NewBuffer(payload)
	unmarshalLingoCmdID(r, &p.ID)

	lookup, ok := Lookup(p.ID, r.Len())
	if !ok {
		p.Payload = UnknownPayload(r.Bytes())
		return nil, fmt.Errorf("unknown command id/size: %#v", p)
	}

	if lookup.Transaction {
		var tr Transaction
		binRead(r, &tr)
		p.Transaction = &tr
	}

	if d, ok := lookup.Payload.(PayloadUnmarshaler); ok {
		if err := d.UnmarshalPayload(r); err != nil {
			return nil, err
		}

	} else {
		binRead(r, lookup.Payload)
	}
	p.Payload = reflect.ValueOf(lookup.Payload).Elem().Interface()
	return

}
