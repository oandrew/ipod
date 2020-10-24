package ipod

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"sync/atomic"

	"log"

	"github.com/sirupsen/logrus"
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

type CommandSerde struct {
	TrxEnabled bool
}

func (s *CommandSerde) handleCmdID(cmdID LingoCmdID) {
	prev := s.TrxEnabled
	switch cmdID {
	// RequestIdentify
	case NewLingoCmdID(LingoGeneralID, 0x00):
		s.TrxEnabled = false
	// IdentifyDeviceLingoes
	case NewLingoCmdID(LingoGeneralID, 0x13):
		s.TrxEnabled = false
	// StartIDPS
	case NewLingoCmdID(LingoGeneralID, 0x38):
		s.TrxEnabled = true
	}
	if prev != s.TrxEnabled {
		if s.TrxEnabled {
			logrus.Warn("Enabling transaction support")
		} else {
			logrus.Warn("Disabling transaction support")
		}
	}
}

func (s *CommandSerde) MarshalCmd(cmd *Command) ([]byte, error) {
	pktBuf := &bytes.Buffer{}

	if err := marshalLingoCmdID(pktBuf, cmd.ID); err != nil {
		return nil, fmt.Errorf("ipod.Command marshal: %v", err)
	}

	s.handleCmdID(cmd.ID)

	if s.TrxEnabled && cmd.Transaction != nil {
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

func (s *CommandSerde) UnmarshalCmd(pkt []byte) (*Command, error) {
	var cmd Command
	pktBuf := bytes.NewBuffer(pkt)
	if err := unmarshalLingoCmdID(pktBuf, &cmd.ID); err != nil {
		return &cmd, fmt.Errorf("ipod.Command unmarshal: %v", err)
	}

	s.handleCmdID(cmd.ID)

	lookup, ok := Lookup(cmd.ID, pktBuf.Len(), s.TrxEnabled)
	if !ok {
		cmd.Payload = UnknownPayload(pktBuf.Bytes())
		return &cmd, fmt.Errorf("ipod.Command unmarshal: unknown cmd %v", cmd.ID)
	}

	if lookup.Transaction {
		var tr Transaction
		err := binary.Read(pktBuf, binary.BigEndian, &tr)
		if err != nil {
			return &cmd, err
		}
		cmd.Transaction = &tr
	}

	if d, ok := lookup.Payload.(encoding.BinaryUnmarshaler); ok {
		err := d.UnmarshalBinary(pktBuf.Bytes())
		if err != nil {
			return &cmd, fmt.Errorf("ipod.Command unmarshal: BinaryUnmarshaler: %v", err)
		}

	} else {
		err := binary.Read(pktBuf, binary.BigEndian, lookup.Payload)
		if err != nil {
			return &cmd, fmt.Errorf("ipod.Command unmarshal: binary.Read: %v", err)
		}
	}
	//cmd.Payload = reflect.Indirect(reflect.ValueOf(lookup.Payload)).Interface()
	cmd.Payload = lookup.Payload
	return &cmd, nil

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

var trxCounter uint32

func TrxReset() {
	atomic.StoreUint32(&trxCounter, 0)
}
func TrxNext() *Transaction {
	trx := atomic.AddUint32(&trxCounter, 1)
	return NewTransaction(uint16(trx))
}

func Send(pw CommandWriter, payload interface{}) {
	cmd, err := BuildCommand(payload)
	if err != nil {
		return
	}
	cmd.Transaction = TrxNext()
	pw.WriteCommand(cmd)
}

type CmdBuffer struct {
	Commands []*Command
}

func (cbuf *CmdBuffer) WriteCommand(cmd *Command) error {
	cbuf.Commands = append(cbuf.Commands, cmd)
	return nil
}
