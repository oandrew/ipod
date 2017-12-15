package main

import (
	"io"

	"github.com/oandrew/ipod/trace"
)

type ReadWriter struct {
	io.Reader
	io.Writer
}

type traceInputReader struct {
	t *trace.Reader
}

func NewTraceInputReader(t *trace.Reader) io.Reader {
	return &traceInputReader{
		t: t,
	}
}

func (tir *traceInputReader) Read(p []byte) (n int, err error) {
	var m trace.Msg
	for {
		if err := tir.t.ReadMsg(&m); err != nil {
			return 0, err
		}
		if m.Dir == trace.DirIn {
			return copy(p, m.Data), nil
		}
	}
}
