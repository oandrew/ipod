package simpleremote

import (
	"encoding/binary"
	"errors"
	"math/bits"
	"strings"

	"github.com/oandrew/ipod"
)

func init() {
	ipod.RegisterLingos(ipod.LingoSimpleRemoteID, Lingos)
}

var Lingos struct {
	ContextButtonStatus       `id:"0x00"`
	ACK                       `id:"0x01"`
	VideoButtonStatus         `id:"0x03"`
	AudioButtonStatus         `id:"0x04"`
	iPodOutButtonStatus       `id:"0x0B"`
	RotationInputStatus       `id:"0x0C"`
	RadioButtonStatus         `id:"0x0D"`
	CameraButtonStatus        `id:"0x0E"`
	RegisterDescriptor        `id:"0x0F"`
	SendHIDReportToiPod       `id:"0x10"`
	SendHIDReportToAcc        `id:"0x11"`
	UnregisterDescriptor      `id:"0x12"`
	AccessibilityEvent        `id:"0x13"`
	GetAccessibilityParameter `id:"0x14"`
	RetAccessibilityParameter `id:"0x15"`
	SetAccessibilityParameter `id:"0x16"`
	GetCurrentItemProperty    `id:"0x17"`
	RetCurrentItemProperty    `id:"0x18"`
	SetContext                `id:"0x19"`
	AccParameterChanged       `id:"0x1A"`
	DevACK                    `id:"0x81"`
}

type ButtonStates struct {
	ButtonStates uint32
	// force binary.Size() == -1
	_ []byte
}

func (s *ButtonStates) MarshalBinary() ([]byte, error) {
	var mask [4]byte
	binary.LittleEndian.PutUint32(mask[:], s.ButtonStates)
	byteLen := (bits.Len32(s.ButtonStates) + 7) / 8
	if byteLen == 0 {
		byteLen = 1
	}
	return mask[:byteLen], nil
}

func (s *ButtonStates) UnmarshalBinary(data []byte) error {
	var mask [4]byte
	switch copy(mask[:], data) {
	case 1, 2, 3, 4:
		s.ButtonStates = binary.LittleEndian.Uint32(mask[:])
		return nil
	default:
		return errors.New("invalid data")
	}
}

//go:generate stringer -type=ContextButtonBit
type ContextButtonBit uint32

const (
	ContextButtonVolumeUp ContextButtonBit = 1 << iota
	ContextButtonVolumeDown
	ContextButtonNextTrack
	ContextButtonPreviousTrack
	ContextButtonNextAlbum
	ContextButtonPreviousAlbum
	ContextButtonStop
	ContextButtonPlayResume
	ContextButtonPause
	ContextButtonMuteToggle
	ContextButtonNextChapter
	ContextButtonPreviousChapter
	ContextButtonNextPlaylist
	ContextButtonPreviousPlaylist
	ContextButtonShuffleSettingAdvance
	ContextButtonRepeatSettingAdvance
	ContextButtonPowerOn
	ContextButtonPowerOff
	ContextButtonBacklightfor30Seconds
	ContextButtonBeginFastForward
	ContextButtonBeginRewind
	ContextButtonMenu
	ContextButtonSelect
	ContextButtonUpArrow
	ContextButtonDownArrow
	ContextButtonBacklightOff
)

type ContextButtonMask uint32

func (m ContextButtonMask) String() string {
	labels := make([]string, 0, 32)
	for i := 0; i < 32; i++ {
		bit := uint32(1 << i)
		if uint32(m)&bit != 0 {
			labels = append(labels, ContextButtonBit(bit).String())
		}
	}
	return strings.Join(labels, " | ")
}

type ContextButtonStatus struct {
	State ContextButtonMask
}

func (s *ContextButtonStatus) MarshalBinary() ([]byte, error) {
	tmp := ButtonStates{ButtonStates: uint32(s.State)}
	return tmp.MarshalBinary()
}

func (s *ContextButtonStatus) UnmarshalBinary(data []byte) error {
	var tmp ButtonStates
	if err := tmp.UnmarshalBinary(data); err != nil {
		return err
	}
	s.State = ContextButtonMask(tmp.ButtonStates)
	return nil
}

type ACK struct{}
type VideoButtonStatus struct {
	ButtonStates
}
type AudioButtonStatus struct {
	ButtonStates
}
type iPodOutButtonStatus struct {
	ButtonSource byte
	// add ButtonStates
}

type RotationInputStatus struct{}
type RadioButtonStatus struct{}
type CameraButtonStatus struct{}
type RegisterDescriptor struct{}
type SendHIDReportToiPod struct{}
type SendHIDReportToAcc struct{}
type UnregisterDescriptor struct{}
type AccessibilityEvent struct{}
type GetAccessibilityParameter struct{}
type RetAccessibilityParameter struct{}
type SetAccessibilityParameter struct{}
type GetCurrentItemProperty struct{}
type RetCurrentItemProperty struct{}
type SetContext struct{}
type AccParameterChanged struct{}
type DevACK struct{}
