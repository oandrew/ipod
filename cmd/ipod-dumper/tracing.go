package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type nopWriter struct {
	io.Reader
}

func (w *nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type tracingReadWriter struct {
	r     io.Reader
	w     io.Writer
	trace *os.File
}

func (t *tracingReadWriter) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if err == nil {
		fmt.Fprintf(t.trace, "< % 02X\n", p[:n])
	}

	return

}

func (t *tracingReadWriter) Write(p []byte) (n int, err error) {
	n, err = t.w.Write(p)
	fmt.Fprintf(t.trace, "> % 02X\n", p[:n])
	return
}

type untracingReadWriter struct {
	s *bufio.Scanner
}

func NewLoadTraceReadWriter(r io.Reader) io.ReadWriter {
	return &untracingReadWriter{
		s: bufio.NewScanner(r),
	}
}

func (t *untracingReadWriter) Read(p []byte) (n int, err error) {

	for t.s.Scan() {
		line := t.s.Text()
		if !strings.HasPrefix(line, "<") {
			continue
		}
		h := strings.Join(strings.Split(line[2:], " "), "")
		log.Debug(h)
		data, err := hex.DecodeString(h)
		if err != nil {
			return 0, err
		}
		return copy(p, data), nil
	}

	return 0, io.EOF

}

func (t *untracingReadWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
