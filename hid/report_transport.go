package hid

import (
	"bytes"
	"io"
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
		Data:        rr.buf[2:n],
	}, nil
}

func NewReportReader(r io.Reader) ReportReader {
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

func NewReportWriter(w io.Writer) ReportWriter {
	return &rawReportWriter{
		w: w,
	}
}
