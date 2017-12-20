package ipod

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// UnknownPayload is a payload  that represents an unknown command
type UnknownPayload []byte

// // PayloadUnmarshaler is the interface implemented by a payload
// // that can unmarshal itself
// type PayloadUnmarshaler interface {
// 	UnmarshalPayload(r io.Reader) error
// }

// // PayloadMarshaler is the interface implemented by a payload
// // that can marshal itself
// type PayloadMarshaler interface {
// 	MarshalPayload(w io.Writer) error
// }

type CommandReader interface {
	ReadCommand() (*Command, error)
}

type CommandWriter interface {
	WriteCommand(*Command) error
}

// type PacketReadWriter interface {
// 	PacketReader
// 	PacketWriter
// }

// Packet is a decoded iap packet
type Command struct {
	ID          LingoCmdID
	Transaction *Transaction
	Payload     interface{}
}

type Transaction uint16

func (tr Transaction) GoString() string {
	return fmt.Sprintf("%#04x", tr)
}

func (tr Transaction) String() string {
	return fmt.Sprintf("%#04x", uint16(tr))
}

func NewTransaction(t uint16) *Transaction {
	tr := Transaction(t)
	return &tr
}

func (cmd *Command) MarshalBinary() ([]byte, error) {
	pktBuf := bytes.NewBuffer(make([]byte, 0, 1024))

	if err := marshalLingoCmdID(pktBuf, cmd.ID); err != nil {
		return nil, err
	}

	if cmd.Transaction != nil {
		binary.Write(pktBuf, binary.BigEndian, *cmd.Transaction)
	}

	if d, ok := cmd.Payload.(encoding.BinaryMarshaler); ok {
		payload, err := d.MarshalBinary()
		if err != nil {
			return nil, err
		}
		pktBuf.Write(payload)
	} else {
		err := binary.Write(pktBuf, binary.BigEndian, cmd.Payload)
		if err != nil {
			return nil, err
		}
	}

	return pktBuf.Bytes(), nil

}

func (cmd *Command) UnmarshalBinary(pkt []byte) error {
	pktBuf := bytes.NewBuffer(pkt)
	if err := unmarshalLingoCmdID(pktBuf, &cmd.ID); err != nil {
		return err
	}

	lookup, ok := Lookup(cmd.ID, pktBuf.Len())
	if !ok {
		cmd.Payload = UnknownPayload(pktBuf.Bytes())
		return fmt.Errorf("unknown command id/size: %#v", cmd)
	}

	if lookup.Transaction {
		var tr Transaction
		err := binary.Read(pktBuf, binary.BigEndian, &tr)
		if err != nil {
			return err
		}
		cmd.Transaction = &tr
	}

	if d, ok := lookup.Payload.(encoding.BinaryUnmarshaler); ok {
		err := d.UnmarshalBinary(pktBuf.Bytes())
		if err != nil {
			return fmt.Errorf("payload unmarshaler: %v", err)
		}

	} else {
		err := binary.Read(pktBuf, binary.BigEndian, lookup.Payload)
		if err != nil {
			return fmt.Errorf("payload simple read: %v", err)
		}
	}
	cmd.Payload = reflect.Indirect(reflect.ValueOf(lookup.Payload)).Interface()
	return nil

}

func BuildCommand(payload interface{}) (*Command, error) {
	id, ok := LookupID(payload)
	if !ok {
		return nil, errors.New("payload not known")
	}
	return &Command{
		ID:      id,
		Payload: payload,
	}, nil

}

func Respond(req *Command, pw CommandWriter, payload interface{}) {
	cmd, err := BuildCommand(payload)
	if err != nil {
		return
	}
	cmd.Transaction = req.Transaction
	pw.WriteCommand(cmd)
}

type CmdBuffer struct {
	Commands []*Command
}

func (cbuf *CmdBuffer) WriteCommand(cmd *Command) error {
	cbuf.Commands = append(cbuf.Commands, cmd)
	return nil
}
