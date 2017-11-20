package audio

import (
	"encoding/binary"
	"io"

	"github.com/oandrew/ipod"
)

func init() {
	ipod.RegisterLingos(LingoAudioID, Lingos)
}

const LingoAudioID = 0x0a

var Lingos struct {
	AccAck                  `id:"0x00"`
	iPodAck                 `id:"0x01"`
	GetAccSampleRateCaps    `id:"0x02"`
	RetAccSampleRateCaps    `id:"0x03"`
	TrackNewAudioAttributes `id:"0x04"`
	SetVideoDelay           `id:"0x05"`
}

type ACKStatus uint8

const (
	ACKStatusSuccess ACKStatus = 0x00
)

type AccAck struct {
	Status ACKStatus
	CmdID  uint8
}
type iPodAck struct {
	Status ACKStatus
	CmdID  uint8
}
type GetAccSampleRateCaps struct {
}
type RetAccSampleRateCaps struct {
	SampleRates []uint32
}

func (s *RetAccSampleRateCaps) UnmarshalPayload(r io.Reader) error {
	for {
		var rate uint32
		if err := binary.Read(r, binary.BigEndian, &rate); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		s.SampleRates = append(s.SampleRates, rate)
	}
	return nil
}

type TrackNewAudioAttributes struct {
	SampleRate       uint32
	SoundCheckValue  uint32
	VolumeAdjustment uint32
}
type SetVideoDelay struct {
	Delay uint32
}
