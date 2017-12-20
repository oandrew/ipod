// Package ipod implements the iPod Accessory protocol (iap)
package ipod

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func (pp PacketPayload) String() string {
	return fmt.Sprintf("(%d)[% 02x]", len([]byte(pp)), []byte(pp))
}

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

func MarshalRawPacket(p *RawPacket) (data []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, p.Length()+3))

	binWrite(buf, PacketStartByte)

	crc := NewCRC8()
	mw := io.MultiWriter(buf, crc)
	if p.Length() > largePacketMinLen {
		binWrite(mw, byte(0x00))
		binWrite(mw, uint16(p.Length()))
	} else {
		binWrite(mw, byte(p.Length()))
	}
	marshalLingoCmdID(mw, p.ID)
	binWrite(mw, p.Data)

	binWrite(buf, crc.Sum8())
	data = buf.Bytes()

	return
}

func UnmarshalRawPacket(data []byte) (*RawPacket, error) {
	if len(data) < 4 {
		return nil, errors.New("raw packet unmarshal: too small")
	}
	if data[0] != PacketStartByte {
		return nil, errors.New("raw packet unmarshal: start byte not found")
	}

	var payOff int
	var payLen int

	if data[1] == 0x00 {
		payLen = int(binary.BigEndian.Uint16(data[2:4]))
		payOff = 4
	} else {
		payLen = int(data[1])
		payOff = 2
	}
	if len(data) < payOff+payLen+1 {
		return nil, errors.New("raw packet unmarshal: packet short")
	}

	payload := data[payOff : payOff+payLen]
	crc := data[payOff+payLen]
	calcCrc := Checksum(data[1 : payOff+payLen])
	if crc != calcCrc {
		return nil, fmt.Errorf("small packet: crc mismatch: recv %02x != calc %02x", crc, calcCrc)
	}

	payloadBuf := bytes.NewBuffer(payload)
	var id LingoCmdID
	unmarshalLingoCmdID(payloadBuf, &id)
	return &RawPacket{
		ID:   id,
		Data: payloadBuf.Bytes(),
	}, nil

}

func MarshalPacket(w io.Writer, p *Packet) (err error) {
	defer catchPanicErr(" Packet marshal: ", &err)

	payloadBuf := bytes.Buffer{}
	if p.Transaction != nil {
		binWrite(&payloadBuf, *p.Transaction)
	}
	if d, ok := p.Payload.(PayloadMarshaler); ok {
		if err := d.MarshalPayload(&payloadBuf); err != nil {
			return err
		}

	} else {
		binWrite(&payloadBuf, p.Payload)
	}

	rawPkt := &RawPacket{
		ID:   p.ID,
		Data: payloadBuf.Bytes(),
	}

	rawData, e := MarshalRawPacket(rawPkt)
	if e != nil {
		err = e
		return
	}
	_, err = w.Write(rawData)
	return

}

func UnmarshalPacket(r io.Reader, p *Packet) (err error) {
	defer catchPanicErr(" Packet unmarshal: ", &err)

	rawData, e := ioutil.ReadAll(r)
	if e != nil {
		err = e
		return
	}
	rawPkt, e := UnmarshalRawPacket(rawData)
	if e != nil {
		err = e
		return
	}

	p.ID = rawPkt.ID

	lookup, ok := Lookup(rawPkt.ID, len(rawPkt.Data))
	if !ok {
		p.Payload = UnknownPayload(rawPkt.Data)
		return fmt.Errorf("unknown command id/size: %#v", p)
	}

	dr := bytes.NewReader(rawPkt.Data)
	if lookup.Transaction {
		var tr Transaction
		binRead(dr, &tr)
		p.Transaction = &tr
	}

	if d, ok := lookup.Payload.(PayloadUnmarshaler); ok {
		if err := d.UnmarshalPayload(dr); err != nil {
			return err
		}

	} else {
		binRead(dr, lookup.Payload)
	}
	p.Payload = reflect.ValueOf(lookup.Payload).Elem().Interface()
	return nil

}
