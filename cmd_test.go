package ipod_test

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/oandrew/ipod"
)

type CustomPayload struct {
	V uint32
}

var TestLingoID uint8 = 0xaa

var TestLingos struct {
	CustomPayload `id:"0x01"`
}

func (p *CustomPayload) UnmarshalBinary(data []byte) error {
	return binary.Read(bytes.NewReader(data), binary.BigEndian, &p.V)
}

func (p CustomPayload) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, p.V)
	return buf.Bytes(), err
}

func TestCommand_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		cmd     ipod.Command
		want    []byte
		wantErr bool
	}{
		{"no-tr-no-payload", ipod.Command{ipod.NewLingoCmdID(0x01, 0x01), nil, nil}, nil, true},
		{"no-tr-with-simple-payload", ipod.Command{ipod.NewLingoCmdID(0x01, 0x02), nil, uint32(0x03)}, []byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x03}, false},
		{"no-tr-with-custom-payload", ipod.Command{ipod.NewLingoCmdID(0x01, 0x02), nil, CustomPayload{0x03}}, []byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x03}, false},
		{"with-tr-with-simple-payload", ipod.Command{ipod.NewLingoCmdID(0x01, 0x02), ipod.NewTransaction(0x01), uint32(0x03)}, []byte{0x01, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("Command.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_UnmarshalBinary(t *testing.T) {
	ipod.RegisterLingos(TestLingoID, TestLingos)

	tests := []struct {
		name    string
		data    []byte
		want    ipod.Command
		wantErr bool
	}{
		{"no-tr-unknown-payload", []byte{0xee, 0x01}, ipod.Command{ipod.NewLingoCmdID(0xee, 0x01), nil, ipod.UnknownPayload{}}, true},
		{"with-tr-unknown-payload", []byte{0xee, 0x01, 0x00, 0x03}, ipod.Command{ipod.NewLingoCmdID(0xee, 0x01), nil, ipod.UnknownPayload{0x00, 0x03}}, true},
		{"no-tr-known-payload", []byte{0xaa, 0x01, 0x00, 0x00, 0x00, 0x03}, ipod.Command{ipod.NewLingoCmdID(0xaa, 0x01), nil, CustomPayload{V: 0x03}}, false},
		{"with-tr-known-payload", []byte{0xaa, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03}, ipod.Command{ipod.NewLingoCmdID(0xaa, 0x01), ipod.NewTransaction(0x02), CustomPayload{V: 0x03}}, false},

		{"no-tr-known-payload-short", []byte{0xaa, 0x01, 0x00, 0x00, 0x03}, ipod.Command{ipod.NewLingoCmdID(0xaa, 0x01), ipod.NewTransaction(0x00), nil}, true},
		{"with-tr-known-payload-short", []byte{0xaa, 0x01, 0x00, 0x02, 0x00, 0x00, 0x03}, ipod.Command{ipod.NewLingoCmdID(0xaa, 0x01), ipod.NewTransaction(0x02), nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ipod.Command
			if err := got.UnmarshalBinary(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Command.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.UnmarshalBinary() = %v, want %v", got, tt.want)
			}

		})
	}
}
