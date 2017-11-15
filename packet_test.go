package ipod_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"git.andrewo.pw/andrew/ipod"
	"git.andrewo.pw/andrew/ipod/lingo-general"
)

type ShortWriter struct {
	max     int
	written int
}

func (sw *ShortWriter) Write(p []byte) (int, error) {
	if sw.written+len(p) > sw.max {
		return 0, io.ErrShortWrite
	}
	sw.written += len(p)
	return len(p), nil
}

func NewShortWriter(max int) io.Writer {
	return &ShortWriter{
		max: max,
	}
}
func TestBinWritePanic(t *testing.T) {
	buf := NewShortWriter(1)
	err := ipod.MarshalSmallPacket(buf, &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{0x03, 0x04}})
	if err == nil {
		t.Error("error == nil")
	}
	t.Logf("binWrite err = %v", err)
}

func TestMarshalSmallPacket(t *testing.T) {
	tests := []struct {
		name    string
		pkt     *ipod.RawPacket
		want    []byte
		wantErr bool
	}{
		{"nil-data", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), nil}, []byte{0x55, 0x02, 0x01, 0x02, 256 - 0x05}, false},
		{"empty-data", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{}}, []byte{0x55, 0x02, 0x01, 0x02, 256 - 0x05}, false},
		{"with-data", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{0xfd}}, []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0xfd}, false},
		{"wrong-size", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), make([]byte, 254)}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bytes.Buffer{}
			err := ipod.MarshalSmallPacket(&got, tt.pkt)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalSmallPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Bytes(), tt.want) {
				t.Errorf("MarshalSmallPacket() = [% 02x], want [% 02x]", got.Bytes(), tt.want)
			}
		})
	}
}

func TestUnmarshalSmallPacket(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		want    *ipod.RawPacket
		wantErr bool
	}{
		{"no-data", []byte{0x55, 0x02, 0x01, 0x02, 256 - 0x05}, &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{}}, false},
		{"with-data", []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0xfd}, &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{0xfd}}, false},
		{"bad-crc", []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0x22}, nil, true},
		{"wrong-start-byte", []byte{0xff, 0x03, 0x01, 0x02, 0xfd, 0xfd}, nil, true},
		{"too-small-payload-length", []byte{0xff, 0x01, 0x01, 0x02, 0xfd, 0xfd}, nil, true},
		{"too-large-payload-length", []byte{0xff, 0x33, 0x01, 0x02, 0xfd, 0xfd}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got = ipod.RawPacket{}
			err := ipod.UnmarshalSmallPacket(bytes.NewReader(tt.raw), &got)
			t.Logf("%+v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalSmallPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("UnmarshalSmallPacket() = %v, want %v", &got, tt.want)
			}
		})
	}
}

func TestMarshalLargePacket(t *testing.T) {
	largeData := bytes.Repeat([]byte{0xee}, 255)
	tests := []struct {
		name    string
		pkt     *ipod.RawPacket
		want    []byte
		wantErr bool
	}{
		{"nil-data", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), nil}, nil, true},
		{"empty-data", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{}}, nil, true},
		{"with-data", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), largeData}, append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0xe9)...), false},
		//{"wrong-size", &Packet{0x01, 0x02, make([]byte, 254)}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			err := ipod.MarshalLargePacket(&buf, tt.pkt)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalLargePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("MarshalLargePacket() = [% 02x], want [% 02x]", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestUnmarshalLargePacket(t *testing.T) {
	largeData := bytes.Repeat([]byte{0xee}, 255)

	tests := []struct {
		name    string
		raw     []byte
		want    *ipod.RawPacket
		wantErr bool
	}{

		{"no-len-marker", []byte{0x55, 0xdd, 0x01, 0x01, 0x1, 0x02}, nil, true},
		{"bad-crc", append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0x22)...), nil, true},
		{"with-data", append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0xe9)...), &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), largeData}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ipod.RawPacket
			err := ipod.UnmarshalLargePacket(bytes.NewReader(tt.raw), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalLargePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("UnmarshalLargePacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalPacket(t *testing.T) {
	tests := []struct {
		name        string
		pkt         *ipod.Packet
		wantLenByte byte
		wantErr     bool
	}{
		{"small", &ipod.Packet{ipod.NewLingoCmdID(0x00, 0x02), nil, general.ACK{0xfd, 0x00}}, 0x04, false},
		//{"large", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), make([]byte, 254)}, 0x00, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bytes.Buffer{}
			err := ipod.MarshalPacket(&got, tt.pkt)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMarshalPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			lenByte := got.Bytes()[1]
			if lenByte != tt.wantLenByte {
				t.Errorf("TestMarshalPacket() = %02x, want %02x", lenByte, tt.wantLenByte)
			}
		})
	}
}

func TestUnmarshalPacket(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		want    *ipod.Packet
		wantErr bool
	}{

		{"small", []byte{0x55, 0x04, 0x00, 0x02, 0xfd, 0x00, 0xfd}, &ipod.Packet{ipod.NewLingoCmdID(0x00, 0x02), nil, general.ACK{0xfd, 0x00}}, false},
		{"large", []byte{0x55, 0x00, 0x00, 0x04, 0x00, 0x02, 0xfd, 0x00, 0xfd}, &ipod.Packet{ipod.NewLingoCmdID(0x00, 0x02), nil, general.ACK{0xfd, 0x00}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ipod.Packet{}
			err := ipod.UnmarshalPacket(bytes.NewReader(tt.raw), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("UnmarshalPacket() = %v, want %v", &got, tt.want)
			}
		})
	}
}

// func TestPacketPrint(t *testing.T) {
// 	p := Packet{LingoID: Lingo0General, CmdID: CmdRequestIdentify, Data: []byte{0xff}}

// 	t.Logf("val: string (%%s): %s", p)
// 	t.Logf("ptr: string (%%s): %s", &p)
// 	t.Logf("val: value (%%v): %v", p)
// 	t.Logf("ptr: value (%%v): %v", &p)
// 	t.Logf("val: value+fields (%%+v): %+v", p)
// 	t.Logf("ptr: value+fields (%%+v): %+v", &p)
// 	t.Logf("val: go-syntax (%%#v): %#v", p)
// 	t.Logf("ptr: go-syntax (%%#v): %#v", &p)

// 	t.Logf("enum %v", LingoCmdID(0x00))

// }
