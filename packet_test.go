package ipod_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/oandrew/ipod"
	_ "github.com/oandrew/ipod/lingo-general"
)

func TestPacketWriter_WritePacket(t *testing.T) {
	largeData := bytes.Repeat([]byte{0xee}, 255)

	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{"no-data", []byte{}, nil, true},
		{"with-data", []byte{0x01, 0x02, 0xfd}, []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0xfd}, false},
		{"large-with-data", append([]byte{0x1, 0x02}, largeData...), append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0xe9)...), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := ipod.NewPacketWriter()
			err := w.WritePacket(tt.data)
			got := w.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("PacketWriter.WritePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PacketWriter.WritePacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacketReader_ReadPacket(t *testing.T) {
	largeData := bytes.Repeat([]byte{0xee}, 255)

	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{"no-data", []byte{0x55, 0x02, 0x01, 0x02, 256 - 0x05}, []byte{0x01, 0x02}, false},
		{"with-data", []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0xfd}, []byte{0x01, 0x02, 0xfd}, false},
		{"bad-crc", []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0x22}, nil, true},
		{"wrong-start-byte", []byte{0xff, 0x03, 0x01, 0x02, 0xfd, 0xfd}, nil, true},

		{"large-with-data", append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0xe9)...), append([]byte{0x1, 0x02}, largeData...), false},
		{"large-with-data-short", append([]byte{0x55, 0x00, 0x01, 0x02, 0x1, 0x02}, append(largeData, 0xe9)...), nil, true},
		{"large-bad-crc", append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0x22)...), nil, true},
		{"short-packet", []byte{0x55}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ipod.NewPacketReader(tt.data)
			got, err := r.ReadPacket()
			if (err != nil) != tt.wantErr {
				t.Errorf("PacketReader.ReadPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PacketReader.ReadPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacketRoundtrip(t *testing.T) {
	pw := [][]byte{
		[]byte("packet1"),
		[]byte("_packet2"),
	}

	w := ipod.NewPacketWriter()
	for i := range pw {
		w.WritePacket(pw[i])
	}

	r := ipod.NewPacketReader(w.Bytes())

	for i := range pw {
		pr, err := r.ReadPacket()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(pr, pw[i]) {
			t.Fail()
		}
	}
}

func BenchmarkPacketReader(b *testing.B) {
	frame := []byte{
		0x55, 0x28, 0x0a, 0x03, 0x03, 0xe7, 0x00, 0x00,
		0x1f, 0x40, 0x00, 0x00, 0x2b, 0x11, 0x00, 0x00,
		0x2e, 0xe0, 0x00, 0x00, 0x3e, 0x80, 0x00, 0x00,
		0x56, 0x22, 0x00, 0x00, 0x5d, 0xc0, 0x00, 0x00,
		0x7d, 0x00, 0x00, 0x00, 0xac, 0x44, 0x00, 0x00,
		0xbb, 0x80, 0x3d, 0x00, 0x00, 0x00, 0x00,
	}

	for i := 0; i < b.N; i++ {
		r := ipod.NewPacketReader(frame)
		r.ReadPacket()
	}
}

func BenchmarkPacketWriter(b *testing.B) {
	packet := []byte{
		0x0a, 0x03, 0x03, 0xe7, 0x00, 0x00, 0x1f, 0x40,
		0x00, 0x00, 0x2b, 0x11, 0x00, 0x00, 0x2e, 0xe0,
		0x00, 0x00, 0x3e, 0x80, 0x00, 0x00, 0x56, 0x22,
		0x00, 0x00, 0x5d, 0xc0, 0x00, 0x00, 0x7d, 0x00,
		0x00, 0x00, 0xac, 0x44, 0x00, 0x00, 0xbb, 0x80,
	}
	for i := 0; i < b.N; i++ {
		pw := ipod.NewPacketWriter()
		pw.WritePacket(packet)
	}
}
