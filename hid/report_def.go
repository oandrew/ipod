package hid

import (
	"errors"
	"fmt"
)

type ReportDir uint8

const (
	ReportDirAccIn  ReportDir = 0
	ReportDirAccOut ReportDir = 1
)

type ReportDef struct {
	ID  int
	Len int
	Dir ReportDir
}

func (def ReportDef) MaxPayload() int {
	return def.Len - 1
}

// Sorted list of reports
type ReportDefs []ReportDef

var DefaultReportDefs = ReportDefs{
	ReportDef{ID: 0x01, Len: 5, Dir: ReportDirAccIn},
	ReportDef{ID: 0x02, Len: 9, Dir: ReportDirAccIn},
	ReportDef{ID: 0x03, Len: 13, Dir: ReportDirAccIn},
	ReportDef{ID: 0x04, Len: 17, Dir: ReportDirAccIn},
	ReportDef{ID: 0x05, Len: 25, Dir: ReportDirAccIn},
	ReportDef{ID: 0x06, Len: 49, Dir: ReportDirAccIn},
	ReportDef{ID: 0x07, Len: 95, Dir: ReportDirAccIn},
	ReportDef{ID: 0x08, Len: 193, Dir: ReportDirAccIn},
	ReportDef{ID: 0x09, Len: 257, Dir: ReportDirAccIn},
	ReportDef{ID: 0x0A, Len: 385, Dir: ReportDirAccIn},
	ReportDef{ID: 0x0B, Len: 513, Dir: ReportDirAccIn},
	ReportDef{ID: 0x0C, Len: 767, Dir: ReportDirAccIn},

	ReportDef{ID: 0x0D, Len: 0x05, Dir: ReportDirAccOut},
	ReportDef{ID: 0x0E, Len: 0x09, Dir: ReportDirAccOut},
	ReportDef{ID: 0x0F, Len: 0x0D, Dir: ReportDirAccOut},
	ReportDef{ID: 0x10, Len: 0x11, Dir: ReportDirAccOut},
	ReportDef{ID: 0x11, Len: 0x19, Dir: ReportDirAccOut},
	ReportDef{ID: 0x12, Len: 0x31, Dir: ReportDirAccOut},
	ReportDef{ID: 0x13, Len: 0x5F, Dir: ReportDirAccOut},
	ReportDef{ID: 0x14, Len: 0xC1, Dir: ReportDirAccOut},
	ReportDef{ID: 0x15, Len: 0xFF, Dir: ReportDirAccOut},
}

func (defs ReportDefs) Pick(payloadSize int, dir ReportDir) (ReportDef, error) {
	var def *ReportDef

	for i := range defs {
		if defs[i].Dir == dir {
			def = &defs[i]
			if defs[i].MaxPayload() >= payloadSize {
				break
			}
		}
	}

	if def == nil {
		return ReportDef{}, errors.New("no matching report found")
	} else {
		return *def, nil
	}
}

func (defs ReportDefs) Find(id int) (ReportDef, error) {
	for i := range defs {
		if defs[i].ID == id {
			return defs[i], nil
		}
	}
	return ReportDef{}, fmt.Errorf("report id no found: %#v", id)
}
