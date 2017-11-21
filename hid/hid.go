// Package hid implements iap over hid transport
package hid

import (
	"bytes"
)

type Report struct {
	ID          byte
	LinkControl LinkControl
	Data        []byte
}

type LinkControl byte

const (
	LinkControlDone         LinkControl = 0x00
	LinkControlContinue     LinkControl = 0x01
	LinkControlMoreToFollow LinkControl = 0x02
)

type Encoder struct {
	reportDefs ReportDefs
	w          ReportWriter
}

func (e *Encoder) WriteFrame(data []byte) error {
	offset := 0
	bytesLeft := len(data)
	for bytesLeft > 0 {
		reportDef, err := e.reportDefs.Pick(bytesLeft, ReportDirAccIn)
		if err != nil {
			return err
		}
		payloadLen := bytesLeft
		linkControl := LinkControlDone
		if bytesLeft > reportDef.MaxPayload() {
			payloadLen = reportDef.MaxPayload()
			if offset == 0 {
				linkControl = LinkControlMoreToFollow
			} else {
				linkControl = LinkControlContinue | LinkControlMoreToFollow
			}
		} else if offset > 0 {
			linkControl = LinkControlContinue
		}
		reportData := make([]byte, reportDef.MaxPayload())
		copy(reportData, data[offset:offset+payloadLen])
		report := Report{
			ID:          byte(reportDef.ID),
			LinkControl: linkControl,
			Data:        reportData,
		}
		if err := e.w.WriteReport(report); err != nil {
			return err
		}
		bytesLeft -= payloadLen
		offset += payloadLen
	}
	return nil

}

func NewEncoder(w ReportWriter, defs ReportDefs) *Encoder {
	return &Encoder{
		reportDefs: defs,
		w:          w,
	}
}

func NewEncoderDefault(w ReportWriter) *Encoder {
	return NewEncoder(w, DefaultReportDefs)
}

type Decoder struct {
	reportDefs ReportDefs
	r          ReportReader
}

func (e *Decoder) ReadFrame() ([]byte, error) {
	buf := bytes.Buffer{}
	done := false
	for !done {
		report, err := e.r.ReadReport()
		if err != nil {
			return nil, err
		}
		reportDef, err := e.reportDefs.Find(int(report.ID))
		if err != nil {
			return nil, err
		}

		reportData := make([]byte, reportDef.MaxPayload())
		n := copy(reportData, report.Data)
		switch report.LinkControl {
		case LinkControlDone:
			buf.Reset()
			buf.Write(reportData[:n])
			done = true
		case LinkControlMoreToFollow:
			buf.Reset()
			buf.Write(reportData[:n])
		case LinkControlContinue | LinkControlMoreToFollow:
			buf.Write(reportData[:n])
		case LinkControlContinue:
			buf.Write(reportData[:n])
			done = true
		}
	}
	return buf.Bytes(), nil
}

func NewDecoder(r ReportReader, defs ReportDefs) *Decoder {
	return &Decoder{
		r:          r,
		reportDefs: defs,
	}
}

func NewDecoderDefault(r ReportReader) *Decoder {
	return NewDecoder(r, DefaultReportDefs)
}

type Transport struct {
	*Decoder
	*Encoder
}

func NewTransport(rw ReportReadWriter, defs ReportDefs) *Transport {
	return &Transport{
		Decoder: NewDecoder(rw, defs),
		Encoder: NewEncoder(rw, defs),
	}
}
