package ipod

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
)

// LingoCmdID represents Lingo ID and Command ID
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

func (id LingoCmdID) String() string {
	return fmt.Sprintf("%#02x,%#0*x", id.LingoID(), cmdIDLen(id.LingoID())*2, id.CmdID())
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

func marshalLingoCmdID(w io.Writer, id LingoCmdID) error {
	err := binary.Write(w, binary.BigEndian, byte(id.LingoID()))
	if err != nil {
		return err
	}
	switch cmdIDLen(id.LingoID()) {
	case 2:
		return binary.Write(w, binary.BigEndian, uint16(id.CmdID()))
	default:
		return binary.Write(w, binary.BigEndian, byte(id.CmdID()))
	}
}

func unmarshalLingoCmdID(r io.Reader, id *LingoCmdID) error {
	var lingoID byte
	err := binary.Read(r, binary.BigEndian, &lingoID)
	if err != nil {
		return err
	}

	switch cmdIDLen(uint16(lingoID)) {
	case 2:
		var cmdID uint16
		err := binary.Read(r, binary.BigEndian, &cmdID)
		if err != nil {
			return err
		}
		*id = NewLingoCmdID(uint16(lingoID), uint16(cmdID))
		return nil
	default:
		var cmdID uint8
		err := binary.Read(r, binary.BigEndian, &cmdID)
		if err != nil {
			return err
		}
		*id = NewLingoCmdID(uint16(lingoID), uint16(cmdID))
		return nil
	}
}

func NewLingoCmdID(lingo, cmd uint16) LingoCmdID {
	return LingoCmdID(uint32(lingo)<<16 | uint32(cmd))
}

func parseIdTag(tag *reflect.StructTag) (uint16, error) {
	id, err := strconv.ParseUint(tag.Get("id"), 0, 16)
	return uint16(id), err
}

var mIDToType = make(map[LingoCmdID][]reflect.Type)
var mTypeToID = make(map[reflect.Type]LingoCmdID)

func storeMapping(cmd LingoCmdID, t reflect.Type) {
	mIDToType[cmd] = append(mIDToType[cmd], t)
	mTypeToID[t] = cmd
}

// RegisterLingos registers a set of lingo commands
func RegisterLingos(lingoID uint8, m interface{}) error {
	lingos := reflect.TypeOf(m)

	for i := 0; i < lingos.NumField(); i++ {
		cmd := lingos.Field(i)
		cmdID, err := parseIdTag(&cmd.Tag)
		if err != nil {
			return fmt.Errorf("register lingos: parse id tag err: %v", err)
		}

		storeMapping(NewLingoCmdID(uint16(lingoID), cmdID), cmd.Type)
	}
	return nil

}

// DumpLingos returns a list of all registered lingos and commands
func DumpLingos() string {
	type cmd struct {
		id   LingoCmdID
		name string
	}
	var cmds []cmd
	for id, types := range mIDToType {
		cmds = append(cmds, cmd{id, types[0].String()})
	}
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].id < cmds[j].id
	})
	buf := bytes.Buffer{}
	for _, cmd := range cmds {
		fmt.Fprintf(&buf, "%s\t%s\n", cmd.id.GoString(), cmd.name)
	}
	return buf.String()

}

// LookupID finds a registered LingoCmdID by the type of v
// i.e. reverse to Lookup
func LookupID(v interface{}) (id LingoCmdID, ok bool) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("payload is not pointer: %v", v))
	}
	id, ok = mTypeToID[t.Elem()]
	return
}

// LookupResult contains the result of a Lookup.
// Payload is a pointer to a new zero value of the found type
// Transaction specifies if the Transaction should be present in the packet.
type LookupResult struct {
	Payload     interface{}
	Transaction bool
}

// Lookup finds a the payload by LingoCmdID using payloadSize as a hint
func Lookup(id LingoCmdID, payloadSize int) (LookupResult, bool) {
	payloads, ok := mIDToType[id]
	if !ok {
		return LookupResult{}, false
	}
	for _, p := range payloads {
		switch {
		case p.Size() == uintptr(payloadSize):
			return LookupResult{
				Payload:     reflect.New(p).Interface(),
				Transaction: false,
			}, true
		case p.Size() == uintptr(payloadSize-2):
			return LookupResult{
				Payload:     reflect.New(p).Interface(),
				Transaction: true,
			}, true
		}
	}
	if len(payloads) == 1 {
		return LookupResult{
			Payload:     reflect.New(payloads[0]).Interface(),
			Transaction: true,
		}, true
	}

	return LookupResult{}, false
}
