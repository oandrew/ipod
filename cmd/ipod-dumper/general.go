package main

import (
	"git.andrewo.pw/andrew/ipod/lingo-general"
)

type DevGeneral struct {
	uimode general.UIMode
}

var _ general.DeviceGeneral = &DevGeneral{}

func (d *DevGeneral) UIMode() general.UIMode {
	return d.uimode
}

func (d *DevGeneral) SetUIMode(mode general.UIMode) {
	d.uimode = mode
}

func (d *DevGeneral) Name() string {
	return "Andrew"
}

func (d *DevGeneral) SoftwareVersion() (major uint8, minor uint8, rev uint8) {
	return 1, 1, 1
}

func (d *DevGeneral) SerialNum() string {
	return "abcd1234"
}

func (d *DevGeneral) LingoProtocolVersion(lingo uint8) (major uint8, minor uint8) {
	switch lingo {
	case 0x0a:
		return 1, 3
	default:
		return 1, 0
	}
}

func (d *DevGeneral) LingoOptions(lingo uint8) uint64 {
	switch lingo {
	case 0x00:
		return 0x000000063DEF73FF

	default:
		return 0
	}
}

func (d *DevGeneral) PrefSettingID(classID uint8) uint8 {
	return 0
}

func (d *DevGeneral) SetPrefSettingID(classID uint8, settingID uint8, restoreOnExit bool) {
}

func (d *DevGeneral) StartIDPS() {
}

func (d *DevGeneral) SetEventNotificationMask(mask uint64) {

}

func (d *DevGeneral) EventNotificationMask() uint64 {
	return 0
}

func (d *DevGeneral) SupportedEventNotificationMask() uint64 {
	return 0
}

func (d *DevGeneral) CancelCommand(lingo uint8, cmd uint16, transaction uint16) {

}

func (d *DevGeneral) MaxPayload() uint16 {
	return 65535
}
