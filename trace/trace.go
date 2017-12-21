package trace

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	//"github.com/oandrew/ipod/trace"
)

const (
	DirIn Dir = iota
	DirOut
)

type Dir byte

func (d Dir) MarshalText() ([]byte, error) {
	switch d {
	case DirIn:
		return []byte{'<'}, nil
	case DirOut:
		return []byte{'>'}, nil
	}

	return nil, fmt.Errorf("bad dir: %v", d)
}

func (d *Dir) UnmarshalText(text []byte) error {
	if len(text) != 1 {
		return fmt.Errorf("trace dir unmarshal: bad value %v", text)
	}
	switch text[0] {
	case '<':
		*d = DirIn
		return nil
	case '>':
		*d = DirOut
		return nil
	}

	return fmt.Errorf("trace dir unmarshal: unknown symbol '%c'", text[0])
}

type Msg struct {
	Dir  Dir
	TS   uint
	Data []byte
}

func (m Msg) MarshalText() ([]byte, error) {
	dt, err := m.Dir.MarshalText()
	if err != nil {
		return nil, err
	}
	if len(m.Data) == 0 {
		return nil, fmt.Errorf("trace marshal: no data")
	}

	t := fmt.Sprintf("%c % 02X", dt[0], m.Data)
	return []byte(t), nil
}

func (m *Msg) UnmarshalText(text []byte) error {
	if len(text) < 4 {
		return fmt.Errorf("trace unmarshal: short msg")
	}
	if err := m.Dir.UnmarshalText(text[0:1]); err != nil {
		return err
	}

	h := bytes.Join(bytes.Split(text[2:], []byte{' '}), []byte{})
	var data []byte
	_, err := fmt.Sscanf(string(h), "%x", &data)
	if err != nil {
		return fmt.Errorf("trace unmarshal: bad data")
	}
	m.Data = data[:]
	return nil

}

type Reader struct {
	s   *bufio.Scanner
	err error
	ts  uint
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		s: bufio.NewScanner(r),
	}
}

func (r *Reader) ReadMsg(m *Msg) error {
	if r.err != nil {
		return r.err
	}
	for r.s.Scan() {
		text := r.s.Bytes()
		if len(text) == 0 {
			continue
		}
		err := m.UnmarshalText(text)
		if err == nil {
			m.TS = r.ts
			r.ts++
		}
		return err
	}
	r.err = r.s.Err()
	if r.err == nil {
		r.err = io.EOF
	}
	return r.err
}

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) WriteMsg(m *Msg) error {
	t, err := m.MarshalText()
	if err != nil {
		return err
	}
	t = append(t, '\n')
	n, err := w.w.Write(t)
	_ = n
	return err
}

type tracer struct {
	tw *Writer
	rw io.ReadWriter
}

func (t *tracer) Write(p []byte) (n int, err error) {
	n, err = t.rw.Write(p)
	if err == nil {
		t.tw.WriteMsg(&Msg{Dir: DirOut, Data: p[:n]})
	}
	return
}

func (t *tracer) Read(p []byte) (n int, err error) {
	n, err = t.rw.Read(p)
	if err == nil {
		t.tw.WriteMsg(&Msg{Dir: DirIn, Data: p[:n]})
	}
	return
}

func NewTracer(tw io.Writer, rw io.ReadWriter) io.ReadWriter {
	return &tracer{
		tw: NewWriter(tw),
		rw: rw,
	}
}

type queueItem struct {
	msg        *Msg
	allE, dirE *list.Element
}

type Queue struct {
	all     list.List
	in, out list.List
}

func (q *Queue) dirList(dir Dir) *list.List {
	switch dir {
	case DirIn:
		return &q.in
	case DirOut:
		return &q.out
	}
	panic("bad dir")
}

func (q *Queue) Enqueue(msg *Msg) {
	qi := &queueItem{msg: msg}
	qi.allE = q.all.PushBack(qi)
	qi.dirE = q.dirList(msg.Dir).PushBack(qi)
}

func (q *Queue) Head() *Msg {
	qie := q.all.Front()
	if qie == nil {
		return nil
	}
	qi, ok := qie.Value.(*queueItem)
	if !ok {
		return nil
	}
	return qi.msg
}

func (q *Queue) Dequeue() *Msg {
	qie := q.all.Front()
	if qie == nil {
		return nil
	}
	qi, ok := qie.Value.(*queueItem)
	if !ok {
		return nil
	}
	q.all.Remove(qi.allE)
	q.dirList(qi.msg.Dir).Remove(qi.dirE)

	return qi.msg
}

func (q *Queue) DequeueDir(dir Dir) *Msg {
	qie := q.dirList(dir).Front()
	if qie == nil {
		return nil
	}
	qi, ok := qie.Value.(*queueItem)
	if !ok {
		return nil
	}
	q.all.Remove(qi.allE)
	q.dirList(qi.msg.Dir).Remove(qi.dirE)

	return qi.msg
}

type queueDirReader struct {
	q   *Queue
	dir Dir
}

func NewQueueDirReader(q *Queue, dir Dir) io.Reader {
	return &queueDirReader{
		q:   q,
		dir: dir,
	}
}

func (qdr *queueDirReader) Read(p []byte) (n int, err error) {
	msg := qdr.q.DequeueDir(qdr.dir)
	if msg == nil {
		return 0, io.EOF
	}
	return copy(p, msg.Data), nil
}
