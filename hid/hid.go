package hid

import (
	"bytes"
	"io"

	"git.andrewo.pw/andrew/ipod"
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

type ReportReader interface {
	ReadReport() (Report, error)
}

type ReportWriter interface {
	WriteReport(Report) error
}

type SingleReport []byte

func (s SingleReport) ReadReport() (Report, error) {
	return Report{
		ID:          s[0],
		LinkControl: LinkControl(s[1]),
		Data:        s[2:],
	}, nil
}

const maxRawSize = 1024

type rawReportReader struct {
	r   io.Reader
	buf []byte
}

func (rr *rawReportReader) ReadReport() (Report, error) {
	n, err := rr.r.Read(rr.buf)
	if err != nil {
		return Report{}, err
	}
	if n < 3 {
		return Report{}, io.ErrNoProgress
	}
	return Report{
		ID:          rr.buf[0],
		LinkControl: LinkControl(rr.buf[1]),
		Data:        rr.buf[2:],
	}, nil
}

func NewRawReportReader(r io.Reader) ReportReader {
	return &rawReportReader{
		r:   r,
		buf: make([]byte, maxRawSize),
	}
}

type rawReportWriter struct {
	w   io.Writer
	buf bytes.Buffer
}

func (rw *rawReportWriter) WriteReport(report Report) error {
	rw.buf.Reset()
	rw.buf.WriteByte(report.ID)
	rw.buf.WriteByte(byte(report.LinkControl))
	rw.buf.Write(report.Data)
	_, err := rw.buf.WriteTo(rw.w)
	return err
}

func NewRawReportWriter(w io.Writer) ReportWriter {
	return &rawReportWriter{
		w: w,
	}
}

type hidEncoder struct {
	reportDefs ReportDefs
	w          ReportWriter
}

func (e *hidEncoder) WriteFrame(data []byte) error {
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

func NewEncoder(w ReportWriter, defs ReportDefs) ipod.TransportEncoder {
	return &hidEncoder{
		reportDefs: defs,
		w:          w,
	}
}

func NewEncoderDefault(w ReportWriter) ipod.TransportEncoder {
	return NewEncoder(w, DefaultReportDefs)
}

type hidDecoder struct {
	reportDefs ReportDefs
	r          ReportReader
}

func (e *hidDecoder) ReadFrame() ([]byte, error) {
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
		copy(reportData, report.Data)
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

func NewDecoder(r ReportReader, defs ReportDefs) ipod.TransportDecoder {
	return &hidDecoder{
		r:          r,
		reportDefs: defs,
	}
}

func NewDecoderDefault(r ReportReader) ipod.TransportDecoder {
	return NewDecoder(r, DefaultReportDefs)
}
