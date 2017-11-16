package ipod

import (
	"log"
	"reflect"
	"strconv"
)

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

func RegisterLingos(lingoID uint8, m interface{}) {
	lingos := reflect.TypeOf(m)

	for i := 0; i < lingos.NumField(); i++ {
		cmd := lingos.Field(i)
		//cmdName := cmd.Name
		cmdID, err := parseIdTag(&cmd.Tag)
		if err != nil {
			log.Printf("bad cmd id: %v", err)
			continue
		}
		//log.Printf("  %s %#02x %#02x", cmdName, lingoID, cmdID)

		storeMapping(NewLingoCmdID(uint16(lingoID), cmdID), cmd.Type)
	}

}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)

}

func LookupID(v interface{}) (id LingoCmdID, ok bool) {
	id, ok = mTypeToID[reflect.TypeOf(v)]
	return
}

type LookupResult struct {
	Payload     interface{}
	Transaction bool
}

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
