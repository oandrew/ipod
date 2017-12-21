package trace_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/oandrew/ipod/trace"
)

func TestWriteRead(t *testing.T) {

	tests := []struct {
		name    string
		m       trace.Msg
		wantErr bool
	}{
		{"simple-in", trace.Msg{Dir: trace.DirIn, Data: []byte{0x00}}, false},
		{"simple-out", trace.Msg{Dir: trace.DirOut, Data: []byte{0x00}}, false},
		{"bad-dir", trace.Msg{Dir: trace.Dir(0xaa), Data: []byte{0x00}}, true},
		{"no-data", trace.Msg{Dir: trace.DirOut, Data: []byte{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			w := trace.NewWriter(&buf)
			t.Logf("msg: %#v", tt.m)
			if err := w.WriteMsg(&tt.m); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			t.Logf("marshaled: %s", buf.String())

			r := trace.NewReader(&buf)
			var mm trace.Msg
			if err := r.ReadMsg(&mm); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(tt.m, mm) {
				t.Errorf("msg1 != msg2: m1=%#v, m2=%#v", tt.m, mm)
			}

		})
	}
}

func TestReadWrite(t *testing.T) {

	tests := []struct {
		name    string
		t       []byte
		wantErr bool
	}{
		{"simple-in", []byte("< 01 02 03\n"), false},
		{"simple-out", []byte("> 01 02 03\n"), false},
		{"bad-dir", []byte("? 01 02 03\n"), true},
		{"no-data", []byte(">\n"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf1 := bytes.NewReader(tt.t)
			r := trace.NewReader(buf1)
			var m trace.Msg
			if err := r.ReadMsg(&m); (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			buf2 := bytes.Buffer{}
			w := trace.NewWriter(&buf2)

			if err := w.WriteMsg(&m); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(tt.t, buf2.Bytes()) {
				t.Errorf("msg1 != msg2: m1: %s, m2: %s", tt.t, buf2.Bytes())
			}

		})
	}
}

func TestTracer(t *testing.T) {
	tbuf := bytes.Buffer{}
	buf := bytes.Buffer{}
	tr := trace.NewTracer(&tbuf, &buf)

	t.Run("write", func(t *testing.T) {
		buf.Reset()
		tbuf.Reset()

		io.WriteString(tr, "ab")
		if buf.String() != "ab" {
			t.Errorf("dest: %s != %s", buf.String(), "ab")
		}

		if tbuf.String() != "> 61 62\n" {
			t.Errorf("trace: %s != %s", tbuf.String(), "> 61 62")
		}
	})

	t.Run("read", func(t *testing.T) {
		buf.Reset()
		tbuf.Reset()

		buf.WriteString("ab")

		data, _ := ioutil.ReadAll(tr)

		if string(data) != "ab" {
			t.Errorf("dest: %s != %s", string(data), "ab")
		}

		if tbuf.String() != "< 61 62\n" {
			t.Errorf("trace: %s != %s", tbuf.String(), "< 61 62")
		}
	})
}

var testReports = `
< 01
> 02
< 03
< 05
< 07
> 04
< 09
< 0A
> 06
< 0B
< 0C
> 08
`

func TestTraceRead(t *testing.T) {
	r := trace.NewReader(strings.NewReader(testReports))
	for {
		var msg trace.Msg
		err := r.ReadMsg(&msg)
		if err == io.EOF {
			break
		}
		t.Logf("msg: %#v", msg)
	}
}

func TestTraceQueue(t *testing.T) {
	r := trace.NewReader(strings.NewReader(testReports))
	var q1 trace.Queue
	var q2 trace.Queue
	var q3 trace.Queue
	var q4 trace.Queue
	for {
		var msg trace.Msg
		err := r.ReadMsg(&msg)
		if err == io.EOF {
			break
		}

		q1.Enqueue(&msg)
		q2.Enqueue(&msg)
		q3.Enqueue(&msg)
		q4.Enqueue(&msg)
	}

	for {
		m := q1.Dequeue()
		if m == nil {
			break
		}
		t.Logf("msg: %#v", m)
	}
	t.Logf("in first")
	for {
		m := q2.DequeueDir(trace.DirIn)
		if m == nil {
			break
		}
		t.Logf("msg: %#v", m)
	}
	t.Logf("out first")
	for {
		m := q3.DequeueDir(trace.DirOut)
		if m == nil {
			break
		}
		t.Logf("msg: %#v", m)
	}

	t.Logf("in,out,in,...")
	dir := trace.DirIn
	for {
		m := q4.DequeueDir(dir)
		if m == nil && q4.Head() == nil {
			break
		}
		if dir == trace.DirIn {
			dir = trace.DirOut
		} else {
			dir = trace.DirIn
		}
		t.Logf("msg: %#v", m)
	}
}
