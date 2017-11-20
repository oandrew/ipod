package simpleremote

import (
	"git.andrewo.pw/andrew/ipod"
)

func init() {
	ipod.RegisterLingos(LingoSimpleRemotelID, Lingos)
}

const LingoSimpleRemotelID = 0x02

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

type ContextButtonStatus struct {
	ButtonStates uint8 // add optional
}

type ACK struct{}
type VideoButtonStatus struct {
	ButtonStates uint8 // add optional
}
type AudioButtonStatus struct {
	ButtonStates uint8 // add optional
}
type iPodOutButtonStatus struct {
	ButtonSource byte
	ButtonMask   uint8 // add optional
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
