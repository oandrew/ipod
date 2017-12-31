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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	buf        bytes.Buffer
}

func (e *Decoder) ReadFrame() ([]byte, error) {
	buf := &e.buf
	buf.Reset()
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

		n := min(len(report.Data), reportDef.MaxPayload())
		reportData := report.Data[:n]
		switch report.LinkControl {
		case LinkControlDone:
			buf.Reset()
			buf.Write(reportData)
			done = true
		case LinkControlMoreToFollow:
			buf.Reset()
			buf.Write(reportData)
		case LinkControlContinue | LinkControlMoreToFollow:
			buf.Write(reportData)
		case LinkControlContinue:
			buf.Write(reportData)
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

func NewTransport(r ReportReader, w ReportWriter, defs ReportDefs) *Transport {
	return &Transport{
		Decoder: NewDecoder(r, defs),
		Encoder: NewEncoder(w, defs),
	}
}
