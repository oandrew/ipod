package hid_test

import (
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/oandrew/ipod/hid"
)

type testReportWriter struct {
	reports []hid.Report
}

func (rw *testReportWriter) WriteReport(report hid.Report) error {
	rw.reports = append(rw.reports, report)
	return nil
}

type testReportReader struct {
	offset  int
	reports []hid.Report
}

func (rr *testReportReader) ReadReport() (hid.Report, error) {
	if rr.offset >= len(rr.reports) {
		return hid.Report{}, io.EOF
	}
	defer func() { rr.offset++ }()
	return rr.reports[rr.offset], nil
}

var testReportDefs1 = hid.ReportDefs{
	hid.ReportDef{ID: 0x01, Len: 2, Dir: hid.ReportDirAccIn},
}

var testReportDefs2 = hid.ReportDefs{
	hid.ReportDef{ID: 0x01, Len: 2, Dir: hid.ReportDirAccIn},
	hid.ReportDef{ID: 0x02, Len: 3, Dir: hid.ReportDirAccIn},
}

var testReportDefs3 = hid.ReportDefs{
	hid.ReportDef{ID: 0x01, Len: 4, Dir: hid.ReportDirAccIn},
}

func TestEncoder(t *testing.T) {
	tests := []struct {
		name       string
		reportDefs hid.ReportDefs
		data       []byte
		want       []hid.Report
		wantErr    bool
	}{
		{"report-1-pkt-1", testReportDefs1, []byte{0x01}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlDone, Data: []byte{0x01}},
		}, false},
		{"report-1-pkt-2", testReportDefs1, []byte{0x01, 0x02}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlMoreToFollow, Data: []byte{0x01}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue, Data: []byte{0x02}},
		}, false},
		{"report-1-pkt-3", testReportDefs1, []byte{0x01, 0x02, 0x03}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlMoreToFollow, Data: []byte{0x01}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue | hid.LinkControlMoreToFollow, Data: []byte{0x02}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue, Data: []byte{0x03}},
		}, false},

		{"report-2-pkt-1", testReportDefs2, []byte{0x01}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlDone, Data: []byte{0x01}},
		}, false},
		{"report-2-pkt-2", testReportDefs2, []byte{0x01, 0x02}, []hid.Report{
			{ID: 0x02, LinkControl: hid.LinkControlDone, Data: []byte{0x01, 0x02}},
		}, false},
		{"report-2-pkt-3", testReportDefs2, []byte{0x01, 0x02, 0x03}, []hid.Report{
			{ID: 0x02, LinkControl: hid.LinkControlMoreToFollow, Data: []byte{0x01, 0x02}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue, Data: []byte{0x03}},
		}, false},

		{"report-3-pkt-1", testReportDefs3, []byte{0x01, 0x02}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlDone, Data: []byte{0x01, 0x02, 0x00}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &testReportWriter{}
			enc := hid.NewEncoder(rw, tt.reportDefs)
			err := enc.WriteFrame(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestHidEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(rw.reports, tt.want) {
				t.Errorf("TestHidEncoder() = [%+v], want [%+v]", rw.reports, tt.want)
			}
		})
	}
}

func TestHidDecoder(t *testing.T) {
	tests := []struct {
		name       string
		reportDefs hid.ReportDefs
		want       []byte
		reports    []hid.Report

		wantErr bool
	}{
		{"report-1-pkt-1", testReportDefs1, []byte{0x01}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlDone, Data: []byte{0x01}},
		}, false},
		{"report-1-pkt-2", testReportDefs1, []byte{0x01, 0x02}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlMoreToFollow, Data: []byte{0x01}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue, Data: []byte{0x02}},
		}, false},
		{"report-1-pkt-3", testReportDefs1, []byte{0x01, 0x02, 0x03}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlMoreToFollow, Data: []byte{0x01}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue | hid.LinkControlMoreToFollow, Data: []byte{0x02}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue, Data: []byte{0x03}},
		}, false},

		{"report-2-pkt-1", testReportDefs2, []byte{0x01}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlDone, Data: []byte{0x01}},
		}, false},
		{"report-2-pkt-2", testReportDefs2, []byte{0x01, 0x02}, []hid.Report{
			{ID: 0x02, LinkControl: hid.LinkControlDone, Data: []byte{0x01, 0x02}},
		}, false},
		{"report-2-pkt-3", testReportDefs2, []byte{0x01, 0x02, 0x03}, []hid.Report{
			{ID: 0x02, LinkControl: hid.LinkControlMoreToFollow, Data: []byte{0x01, 0x02}},
			{ID: 0x01, LinkControl: hid.LinkControlContinue, Data: []byte{0x03}},
		}, false},

		{"report-3-pkt-1", testReportDefs3, []byte{0x01, 0x02, 0x00}, []hid.Report{
			{ID: 0x01, LinkControl: hid.LinkControlDone, Data: []byte{0x01, 0x02, 0x00}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := &testReportReader{
				reports: tt.reports,
			}
			dec := hid.NewDecoder(rr, tt.reportDefs)
			payload, err := dec.ReadFrame()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestHidDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(payload, tt.want) {
				t.Errorf("TestHidDecoder() = %v, want %v", payload, tt.want)
			}
		})
	}
}

func BenchmarkDecoder(b *testing.B) {
	report := []byte{
		0x12, 0x00, 0x55, 0x28, 0x0a, 0x03, 0x03, 0xe7,
		0x00, 0x00, 0x1f, 0x40, 0x00, 0x00, 0x2b, 0x11,
		0x00, 0x00, 0x2e, 0xe0, 0x00, 0x00, 0x3e, 0x80,
		0x00, 0x00, 0x56, 0x22, 0x00, 0x00, 0x5d, 0xc0,
		0x00, 0x00, 0x7d, 0x00, 0x00, 0x00, 0xac, 0x44,
		0x00, 0x00, 0xbb, 0x80, 0x3d, 0x00, 0x00, 0x00,
		0x00,
	}
	r := hid.SingleReport(report)
	d := hid.NewDecoderDefault(r)
	for i := 0; i < b.N; i++ {
		d.ReadFrame()
	}
}

func BenchmarkEncoder(b *testing.B) {
	frame := []byte{
		0x55, 0x28, 0x0a, 0x03, 0x03, 0xe7, 0x00, 0x00,
		0x1f, 0x40, 0x00, 0x00, 0x2b, 0x11, 0x00, 0x00,
		0x2e, 0xe0, 0x00, 0x00, 0x3e, 0x80, 0x00, 0x00,
		0x56, 0x22, 0x00, 0x00, 0x5d, 0xc0, 0x00, 0x00,
		0x7d, 0x00, 0x00, 0x00, 0xac, 0x44, 0x00, 0x00,
		0xbb, 0x80, 0x3d, 0x00, 0x00, 0x00, 0x00,
	}
	e := hid.NewEncoderDefault(hid.NewReportWriter(ioutil.Discard))
	for i := 0; i < b.N; i++ {
		e.WriteFrame(frame)
	}
}
