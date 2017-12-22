package ipod

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"

	"log"
)

// UnknownPayload is a payload  that represents an unknown command
type UnknownPayload []byte

type CommandReader interface {
	ReadCommand() (*Command, error)
}

type CommandWriter interface {
	WriteCommand(*Command) error
}

// Command represents iap packet payload
type Command struct {
	ID LingoCmdID
	//Optional
	Transaction *Transaction
	Payload     interface{}
}

type Transaction uint16

func NewTransaction(t uint16) *Transaction {
	tr := Transaction(t)
	return &tr
}

func (tr Transaction) GoString() string {
	return fmt.Sprintf("%#04x", tr)
}

func (tr Transaction) String() string {
	return fmt.Sprintf("%#04x", uint16(tr))
}

func (tr *Transaction) Copy() *Transaction {
	if tr != nil {
		ctr := Transaction(*tr)
		return &ctr
	}
	return nil
}
func (tr *Transaction) Delta(d int) *Transaction {
	if tr != nil {
		return NewTransaction(uint16(int(*tr) + d))
	}
	return nil
}

func (cmd *Command) MarshalBinary() ([]byte, error) {
	pktBuf := bytes.NewBuffer(make([]byte, 0, 1024))

	if err := marshalLingoCmdID(pktBuf, cmd.ID); err != nil {
		return nil, fmt.Errorf("ipod.Command marshal: %v", err)
	}

	if cmd.Transaction != nil {
		binary.Write(pktBuf, binary.BigEndian, *cmd.Transaction)
	}
	if cmd.Payload == nil {
		return nil, fmt.Errorf("ipod.Command marshal: nil payload")
	}

	if d, ok := cmd.Payload.(encoding.BinaryMarshaler); ok {
		payload, err := d.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("ipod.Command marshal: BinaryMarshaler: %v", err)
		}
		pktBuf.Write(payload)
	} else {
		err := binary.Write(pktBuf, binary.BigEndian, cmd.Payload)
		if err != nil {
			return nil, fmt.Errorf("ipod.Command marshal: binary.Write: %v", err)
		}
	}

	return pktBuf.Bytes(), nil

}

func (cmd *Command) UnmarshalBinary(pkt []byte) error {
	pktBuf := bytes.NewBuffer(pkt)
	if err := unmarshalLingoCmdID(pktBuf, &cmd.ID); err != nil {
		return fmt.Errorf("ipod.Command unmarshal: %v", err)
	}

	lookup, ok := Lookup(cmd.ID, pktBuf.Len())
	if !ok {
		cmd.Payload = UnknownPayload(pktBuf.Bytes())
		return fmt.Errorf("ipod.Command unmarshal: unknown cmd %v", cmd.ID)
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
			return fmt.Errorf("ipod.Command unmarshal: BinaryUnmarshaler: %v", err)
		}

	} else {
		err := binary.Read(pktBuf, binary.BigEndian, lookup.Payload)
		if err != nil {
			return fmt.Errorf("ipod.Command unmarshal: binary.Read: %v", err)
		}
	}
	//cmd.Payload = reflect.Indirect(reflect.ValueOf(lookup.Payload)).Interface()
	cmd.Payload = lookup.Payload
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
		log.Printf("BuildCommand err: %v", err)
		return
	}
	cmd.Transaction = req.Transaction.Copy()
	pw.WriteCommand(cmd)
}

func Send(pw CommandWriter, payload interface{}, tr *Transaction) {
	cmd, err := BuildCommand(payload)
	if err != nil {
		return
	}
	cmd.Transaction = tr
	pw.WriteCommand(cmd)
}

type CmdBuffer struct {
	Commands []*Command
}

func (cbuf *CmdBuffer) WriteCommand(cmd *Command) error {
	cbuf.Commands = append(cbuf.Commands, cmd)
	return nil
}
