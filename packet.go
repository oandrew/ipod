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

type Packet struct {
	ID          LingoCmdID
	Transaction *Transaction
	Payload     interface{}
}

func (p Packet) WithTransaction(t uint16) Packet {
	tr := Transaction(t)
	p.Transaction = &tr
	return p
}

func BuildPacket(payload interface{}) (Packet, error) {
	id, ok := LookupID(payload)
	if !ok {
		return Packet{}, errors.New("payload not known")
	}
	return Packet{
		ID:      id,
		Payload: payload,
	}, nil

}

func Respond(req Packet, pw PacketWriter, payload interface{}) {
	p, err := BuildPacket(payload)
	if err != nil {
		return
	}
	p.Transaction = req.Transaction
	pw.WritePacket(p)
}

type LingoCmdID uint32

func (id LingoCmdID) LingoID() uint16 {
	return uint16(id >> 16 & 0xffff)
}

func (id LingoCmdID) CmdID() uint16 {
	return uint16(id & 0xffff)
}

func (id LingoCmdID) GoString() string {
	return fmt.Sprintf("(%#02x|%#0*x)", id.LingoID(), cmdIDLen(id.LingoID())*2, id.CmdID())
}

func (id LingoCmdID) len() int {
	return 1 + cmdIDLen(id.LingoID())
}

func cmdIDLen(lingoID uint16) int {
	switch lingoID {
	case 0x04:
		return 2
	default:
		return 1
	}
}

func marshalLingoCmdID(w io.Writer, id LingoCmdID) {
	binWrite(w, byte(id.LingoID()))
	switch cmdIDLen(id.LingoID()) {
	case 2:
		binWrite(w, uint16(id.CmdID()))
	default:
		binWrite(w, byte(id.CmdID()))
	}
}

func unmarshalLingoCmdID(r io.Reader, id *LingoCmdID) {
	var lingoID byte
	binRead(r, &lingoID)
	switch cmdIDLen(uint16(lingoID)) {
	case 2:
		var cmdID uint16
		binRead(r, &cmdID)
		*id = NewLingoCmdID(uint16(lingoID), uint16(cmdID))
	default:
		var cmdID uint8
		binRead(r, &cmdID)
		*id = NewLingoCmdID(uint16(lingoID), uint16(cmdID))
	}
}

func NewLingoCmdID(lingo, cmd uint16) LingoCmdID {
	return LingoCmdID(uint32(lingo)<<16 | uint32(cmd))
}

type Transaction uint16

func (tr Transaction) GoString() string {
	return fmt.Sprintf("%#04x", tr)
}

type PayloadUnmarshaler interface {
	UnmarshalPayload(r io.Reader) error
}

type PayloadMarshaler interface {
	MarshalPayload(w io.Writer) error
}

type PacketReader interface {
	ReadPacket() (Packet, error)
}

type PacketWriter interface {
	WritePacket(Packet) error
}

type PacketReadWriter interface {
	PacketReader
	PacketWriter
}

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

func MarshalSmallPacket(w io.Writer, p *RawPacket) (err error) {
	if p.Length() >= largePacketMinLen {
		return errors.New("small Packet: payload is too large")
	}
	defer catchPanicErr("small Packet marshal: ", &err)

	binWrite(w, PacketStartByte)

	crc := NewCRC8()
	mw := io.MultiWriter(w, crc)
	binWrite(mw, byte(p.Length()))
	marshalLingoCmdID(mw, p.ID)
	binWrite(mw, p.Data)

	binWrite(w, crc.Sum8())

	return
}

func UnmarshalSmallPacket(r io.Reader, p *RawPacket) (err error) {
	defer catchPanicErr("small Packet unmarshal: ", &err)

	var startByte byte
	binRead(r, &startByte)
	if startByte != PacketStartByte {
		return errors.New("small Packet: start byte not found")
	}
	crc := NewCRC8()
	tr := io.TeeReader(r, crc)
	var payloadLen byte
	binRead(tr, &payloadLen)
	if payloadLen < 2 {
		return errors.New("small Packet: wrong length")
	}
	payloadData := make([]byte, int(payloadLen))
	binRead(tr, payloadData)

	var crcWant byte
	binRead(r, &crcWant)
	if crc.Sum8() != crcWant {
		return errors.New("small Packet: crc mismatch")
	}

	payloadBuf := bytes.NewBuffer(payloadData)
	var id LingoCmdID
	unmarshalLingoCmdID(payloadBuf, &id)
	*p = RawPacket{
		ID:   id,
		Data: payloadBuf.Bytes(),
	}
	return

}

func MarshalLargePacket(w io.Writer, p *RawPacket) (err error) {
	defer catchPanicErr("large Packet marshal: ", &err)

	if p.Length() < largePacketMinLen {
		return errors.New("large Packet: payload too small")
	}
	if p.Length() > 65535 {
		return errors.New("large Packet: payload too large")
	}

	binWrite(w, PacketStartByte)
	binWrite(w, byte(0x00)) //len marker

	crc := NewCRC8()
	mw := io.MultiWriter(w, crc)
	binWrite(mw, uint16(p.Length()))
	marshalLingoCmdID(mw, p.ID)
	binWrite(mw, p.Data)

	binWrite(w, crc.Sum8())

	return nil
}

func UnmarshalLargePacket(r io.Reader, p *RawPacket) (err error) {
	defer catchPanicErr("large Packet unmarshal: ", &err)

	var startByte byte
	binRead(r, &startByte)

	if startByte != PacketStartByte {
		return errors.New("large Packet: start byte not found")
	}

	var marker byte
	if err := binRead(r, &marker); marker != 0x00 || err != nil {
		return errors.New("large Packet: payload len marker not found")
	}

	crc := NewCRC8()
	tr := io.TeeReader(r, crc)

	var payloadLen uint16
	binRead(tr, &payloadLen)

	payloadData := make([]byte, int(payloadLen))
	binRead(tr, payloadData)

	var crcWant byte
	binRead(r, &crcWant)

	if crc.Sum8() != crcWant {
		return errors.New("large Packet: crc mismatch")
	}

	payloadBuf := bytes.NewBuffer(payloadData)
	var id LingoCmdID
	unmarshalLingoCmdID(payloadBuf, &id)

	*p = RawPacket{
		ID:   id,
		Data: payloadBuf.Bytes(),
	}
	return nil

}

func MarshalPacket(w io.Writer, p *Packet) (err error) {
	defer catchPanicErr(" Packet marshal: ", &err)

	payloadBuf := bytes.Buffer{}
	if p.Transaction != nil {
		binWrite(&payloadBuf, *p.Transaction)
	}
	if d, ok := p.Payload.(PayloadMarshaler); ok {
		//log.Printf("Custom PayloadMarshaler")
		if err := d.MarshalPayload(&payloadBuf); err != nil {
			return err
		}

	} else {
		binWrite(&payloadBuf, p.Payload)
	}

	RawPacket := &RawPacket{
		ID:   p.ID,
		Data: payloadBuf.Bytes(),
	}

	if RawPacket.Length() < largePacketMinLen {
		return MarshalSmallPacket(w, RawPacket)
	} else {
		return MarshalLargePacket(w, RawPacket)
	}

}

func UnmarshalPacket(r io.Reader, pp *Packet) (err error) {
	defer catchPanicErr(" Packet unmarshal: ", &err)

	br := bufio.NewReader(r)
	for {
		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b == PacketStartByte {
			br.UnreadByte()
			break
		}
	}
	header, err := br.Peek(2)
	if err != nil {
		return err
	}

	var p RawPacket
	if header[1] == 0x00 {
		if err := UnmarshalLargePacket(br, &p); err != nil {
			return err
		}
	} else {
		if err := UnmarshalSmallPacket(br, &p); err != nil {
			return err
		}
	}

	lookup, ok := Lookup(p.ID, len(p.Data))
	if !ok {
		return fmt.Errorf("id/size not found: %#v", p)
	}

	//log.Printf("lookup: %#v", lookup)

	pp.ID = p.ID
	dr := bytes.NewReader(p.Data)
	if lookup.Transaction {
		var tr Transaction
		binRead(dr, &tr)
		pp.Transaction = &tr
	}

	if d, ok := lookup.Payload.(PayloadUnmarshaler); ok {
		if err := d.UnmarshalPayload(dr); err != nil {
			//log.Print(err)
			return err
		}

	} else {
		binRead(dr, lookup.Payload)
	}
	pp.Payload = reflect.ValueOf(lookup.Payload).Elem().Interface()
	//log.Printf("%#v", pp)
	return nil

}

type RawPacketWriter struct {
	w io.Writer
}

func (rpw *RawPacketWriter) WritePacket(pkt Packet) error {
	return MarshalPacket(rpw.w, &pkt)
}

func NewRawPacketWriter(w io.Writer) PacketWriter {
	return &RawPacketWriter{
		w: w,
	}
}
