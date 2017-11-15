package ipod_test

import (
	"testing"

	"git.andrewo.pw/andrew/ipod"
)

func TestChecksum(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		payload []byte
		wantCrc byte
	}{
		{"empty", []byte{}, 0x00},
		{"zero crc", []byte{0x00, 0x00}, 0x00},
		{"simple overflow", []byte{0xfe, 0x01}, 0x01},
		{"overflow", []byte{0xff, 0xff, 0x02}, 0x00},
		{"random", []byte{0x04, 0x00, 0x11, 0x00, 0x01}, 0xea},
		{"long", []byte{0x0d, 0x00, 0x4c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x06, 0x3d, 0xef, 0x73, 0xff}, 0xff},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipod.Checksum(tt.payload); got != tt.wantCrc {
				t.Errorf("checksum() = %v, want %v", got, tt.wantCrc)
			}
		})
	}
}
