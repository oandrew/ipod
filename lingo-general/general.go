package general

import (
	"bufio"
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"

	"github.com/oandrew/ipod"
)

func init() {
	ipod.RegisterLingos(LingoGeneralID, Lingos)
}

const LingoGeneralID = 0x00

var Lingos struct {
	RequestIdentify                `id:"0x00"`
	ACK                            `id:"0x02"`
	ACKPending                     `id:"0x02"`
	ACKDataDropped                 `id:"0x02"`
	RequestRemoteUIMode            `id:"0x03"`
	ReturnRemoteUIMode             `id:"0x04"`
	EnterRemoteUIMode              `id:"0x05"`
	ExitRemoteUIMode               `id:"0x06"`
	RequestiPodName                `id:"0x07"`
	ReturniPodName                 `id:"0x08"`
	RequestiPodSoftwareVersion     `id:"0x09"`
	ReturniPodSoftwareVersion      `id:"0x0A"`
	RequestiPodSerialNum           `id:"0x0B"`
	ReturniPodSerialNum            `id:"0x0C"`
	RequestLingoProtocolVersion    `id:"0x0F"`
	ReturnLingoProtocolVersion     `id:"0x10"`
	RequestTransportMaxPayloadSize `id:"0x11"`
	ReturnTransportMaxPayloadSize  `id:"0x12"`
	IdentifyDeviceLingoes          `id:"0x13"`
	GetDevAuthenticationInfo       `id:"0x14"`
	//RetDevAuthenticationInfoV1      `id:"0x15"`
	//RetDevAuthenticationInfoV2      `id:"0x15"`
	RetDevAuthenticationInfo        `id:"0x15"`
	AckDevAuthenticationInfo        `id:"0x16"`
	GetDevAuthenticationSignatureV1 `id:"0x17"`
	GetDevAuthenticationSignatureV2 `id:"0x17"`
	//RetDevAuthenticationSignatureV1 `id:"0x18"`
	//RetDevAuthenticationSignatureV2 `id:"0x18"`
	RetDevAuthenticationSignature  `id:"0x18"`
	AckDevAuthenticationStatus     `id:"0x19"`
	GetiPodAuthenticationInfo      `id:"0x1A"`
	RetiPodAuthenticationInfo      `id:"0x1B"`
	AckiPodAuthenticationInfo      `id:"0x1C"`
	GetiPodAuthenticationSignature `id:"0x1D"`
	RetiPodAuthenticationSignature `id:"0x1E"`
	AckiPodAuthenticationStatus    `id:"0x1F"`
	NotifyiPodStateChange          `id:"0x23"`
	GetiPodOptions                 `id:"0x24"`
	RetiPodOptions                 `id:"0x25"`
	GetAccessoryInfo               `id:"0x27"`
	RetAccessoryInfo               `id:"0x28"`
	GetiPodPreferences             `id:"0x29"`
	RetiPodPreferences             `id:"0x2A"`
	SetiPodPreferences             `id:"0x2B"`
	GetUIMode                      `id:"0x35"`
	RetUIMode                      `id:"0x36"`
	SetUIMode                      `id:"0x37"`
	StartIDPS                      `id:"0x38"`
	SetFIDTokenValues              `id:"0x39"`
	RetFIDTokenValueACKs           `id:"0x3A"`
	EndIDPS                        `id:"0x3B"`
	IDPSStatus                     `id:"0x3C"`
	OpenDataSessionForProtocol     `id:"0x3F"`
	CloseDataSession               `id:"0x40"`
	DevACK                         `id:"0x41"`
	DevDataTransfer                `id:"0x42"`
	IPodDataTransfer               `id:"0x43"`
	SetAccStatusNotification       `id:"0x46"`
	RetAccStatusNotification       `id:"0x47"`
	AccessoryStatusNotification    `id:"0x48"`
	SetEventNotification           `id:"0x49"`
	IPodNotification               `id:"0x4A"`
	GetiPodOptionsForLingo         `id:"0x4B"`
	RetiPodOptionsForLingo         `id:"0x4C"`
	GetEventNotification           `id:"0x4D"`
	RetEventNotification           `id:"0x4E"`
	GetSupportedEventNotification  `id:"0x4F"`
	CancelCommand                  `id:"0x50"`
	RetSupportedEventNotification  `id:"0x51"`
	SetAvailableCurrent            `id:"0x54"`
	RequestApplicationLaunch       `id:"0x64"`
	GetNowPlayingFocusApp          `id:"0x65"`
	RetNowPlayingFocusApp          `id:"0x66"`
}

type RequestIdentify struct{}

type ACKStatus uint8

const (
	ACKStatusSuccess  ACKStatus = 0x00
	ACKStatusUnkownID ACKStatus = 0x05
	ACKStatusPending  ACKStatus = 0x06
)

type ACK struct {
	Status ACKStatus
	CmdID  uint8
}

type ACKPending struct {
	Status  ACKStatus
	CmdID   uint8
	MaxWait uint32
}

type ACKDataDropped struct {
	Status          ACKStatus
	CmdID           uint8
	SessionID       uint16
	NumBytesDropped uint32
}

type RequestRemoteUIMode struct{}

type ReturnRemoteUIMode struct {
	Mode byte
}

type EnterRemoteUIMode struct{}

type ExitRemoteUIMode struct{}

type RequestiPodName struct{}

type ReturniPodName struct {
	Name []byte
}

func (s ReturniPodName) MarshalBinary() ([]byte, error) {
	return s.Name, nil
}

func (s *ReturniPodName) UnmarshalBinary(data []byte) error {
	s.Name = make([]byte, len(data))
	copy(s.Name, data)
	return nil
}

type RequestiPodSoftwareVersion struct{}

type ReturniPodSoftwareVersion struct {
	Major byte
	Minor byte
	Rev   byte
}
type RequestiPodSerialNum struct {
}
type ReturniPodSerialNum struct {
	Serial []byte
}

func (s ReturniPodSerialNum) MarshalBinary() ([]byte, error) {
	return s.Serial, nil
}

func (s *ReturniPodSerialNum) UnmarshalBinary(data []byte) error {
	s.Serial = make([]byte, len(data))
	copy(s.Serial, data)
	return nil
}

type RequestLingoProtocolVersion struct {
	Lingo byte
}

type ReturnLingoProtocolVersion struct {
	Lingo byte
	Major byte
	Minor byte
}

type RequestTransportMaxPayloadSize struct{}

type ReturnTransportMaxPayloadSize struct {
	MaxPayload uint16
}

type IdentifyDeviceLingoes struct {
	Lingos   uint32
	Options  uint32
	DeviceID uint32
}

type GetDevAuthenticationInfo struct{}

// type RetDevAuthenticationInfo struct {
// 	Major byte
// 	Minor byte
// }

type RetDevAuthenticationInfo struct {
	Major              byte
	Minor              byte
	CertCurrentSection byte
	CertMaxSection     byte
	CertData           []byte
}

func (s *RetDevAuthenticationInfo) UnmarshalBinary(r []byte) error {
	if len(r) < 2 {
		return errors.New("short packet")
	}
	s.Major, s.Minor = r[0], r[1]

	if s.Major >= 0x02 {
		// if len(r) < 5 {
		// 	return errors.New("short packet")
		// }
		s.CertCurrentSection, s.CertMaxSection = r[2], r[3]

		s.CertData = make([]byte, len(r[4:]))
		copy(s.CertData, r[4:])
	}

	return nil

}

type DevAuthInfoStatus uint8

const (
	DevAuthInfoStatusSupported DevAuthInfoStatus = 0x00
)

type AckDevAuthenticationInfo struct {
	Status DevAuthInfoStatus
}
type GetDevAuthenticationSignatureV1 struct {
	Challenge [16]byte
	Counter   byte
}

type GetDevAuthenticationSignatureV2 struct {
	Challenge [20]byte
	Counter   byte
}

type RetDevAuthenticationSignature struct {
	Signature []byte
}

func (s *RetDevAuthenticationSignature) UnmarshalBinary(r []byte) error {
	s.Signature = make([]byte, len(r))
	copy(s.Signature, r)
	return nil
}

type DevAuthStatus uint8

const (
	DevAuthStatusPassed DevAuthStatus = 0x00
	DevAuthStatusFailed DevAuthStatus = 0x01
)

type AckDevAuthenticationStatus struct {
	Status DevAuthStatus
}
type GetiPodAuthenticationInfo struct{}

type RetiPodAuthenticationInfo struct {
	Major              byte
	Minor              byte
	CertCurrentSection byte
	CertMaxSection     byte
	CertData           []byte
}
type AckiPodAuthenticationInfo struct {
	Status byte
}

type GetiPodAuthenticationSignature struct {
	Challenge [20]byte
	Counter   byte
}
type RetiPodAuthenticationSignature struct {
	Signature [20]byte
}

type AckiPodAuthenticationStatus struct {
	Status byte
}
type NotifyiPodStateChange struct {
	StateChange byte
}
type GetiPodOptions struct{}

type RetiPodOptions struct {
	Options uint64
}

type GetAccessoryInfo struct {
	InfoType byte
}

type GetAccessoryInfo2 struct {
	InfoType byte
	ModelID  uint32
	Major    byte
	Minor    byte
	Rev      byte
}

type GetAccessoryInfo3 struct {
	InfoType byte
	LingoID  byte
}

type RetAccessoryInfo struct {
	InfoType byte
	Data     []byte
}

// type RetAccessoryInfo0 struct {
// 	InfoType byte
// 	Caps uint32
// }

// type RetAccessoryInfo1678 struct {
// 	InfoType byte
// 	Data []byte
// }

// type RetAccessoryInfo2 struct {
// 	InfoType byte
// 	ModelID uint32
// 	MinMajor byte
// 	MinMinor byte
// 	MinRev   byte
// }

// type RetAccessoryInfo3 struct {
// 	InfoType byte
// 	ModelID uint32
// 	MinMajor byte
// 	MinMinor byte
// 	MinRev   byte
// }

type GetiPodPreferences struct {
	PrefClassID byte
}
type RetiPodPreferences struct {
	PrefClassID        byte
	PrefClassSettingID byte
}

type SetiPodPreferences struct {
	PrefClassID        byte
	PrefClassSettingID byte
	RestoreOnExit      byte
}

type UIMode uint8

const (
	UIModeStandart UIMode = 0x00
	UIModeExtended UIMode = 0x01
	UIModeiPodOut  UIMode = 0x02
)

type GetUIMode struct{}
type RetUIMode struct {
	UIMode UIMode
}

type SetUIMode struct {
	UIMode UIMode
}

type StartIDPS struct{}

type FIDIdentifyToken struct {
	NumLingoes    byte
	AccLingoes    []byte
	DeviceOptions uint32
	DeviceID      uint32
}

func (t *FIDIdentifyToken) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &t.NumLingoes)
	t.AccLingoes = make([]byte, t.NumLingoes)
	binary.Read(r, binary.BigEndian, &t.AccLingoes)
	binary.Read(r, binary.BigEndian, &t.DeviceOptions)
	binary.Read(r, binary.BigEndian, &t.DeviceID)
	return nil
}

//go:generate stringer -type=AccCapBit
type AccCapBit uint32

const (
	AccCapAnalogLineOut AccCapBit = 1 << iota
	AccCapAnalogLineIn
	AccCapAnalogVideoOut
	_
	AccCapUSBAudio
	_
	_
	_
	_
	AccCapAppComm
	_
	AccCapCheckVolume
)

var AccCaps = []AccCapBit{
	AccCapAnalogLineOut, AccCapAnalogLineIn,
	AccCapAnalogVideoOut, AccCapUSBAudio,
	AccCapAppComm, AccCapCheckVolume,
}

type FIDAccCapsToken struct {
	AccCapsBitmask uint64
}

func (t *FIDAccCapsToken) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &t.AccCapsBitmask)
	return nil
}

//go:generate stringer -type=AccInfoType
type AccInfoType uint8

const (
	AccInfoName       AccInfoType = 0x01
	AccInfoFirmware   AccInfoType = 0x04
	AccInfoHardware   AccInfoType = 0x05
	AccInfoMfr        AccInfoType = 0x06
	AccInfoModel      AccInfoType = 0x07
	AccInfoSerial     AccInfoType = 0x08
	AccInfoMaxPayload AccInfoType = 0x09
)

type FIDAccInfoToken struct {
	AccInfoType byte
	Value       interface{}
}

func (t *FIDAccInfoToken) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &t.AccInfoType)
	switch t.AccInfoType {
	//name
	case 0x01, 0x06, 0x07, 0x08:
		t.Value, _ = bufio.NewReader(r).ReadBytes(0x00)
	case 0x04, 0x05:
		v := make([]byte, 3)
		r.Read(v)
		t.Value = v
	case 0x09:
		v := make([]byte, 2)
		r.Read(v)
		t.Value = v
	case 0x0b, 0x0c:
		v := make([]byte, 4)
		r.Read(v)
		t.Value = v
	default:
		return errors.New("unknown AccInfoToken type")
	}
	return nil
}

type FIDiPodPreferenceToken struct {
	PrefClass        byte
	PrefClassSetting byte
	RestoreOnExit    byte
}

func (t *FIDiPodPreferenceToken) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	return binary.Read(r, binary.BigEndian, t)
}

type FIDEAProtocolToken struct {
	ProtocolIndex  byte
	ProtocolString []byte
}

func (t *FIDEAProtocolToken) UnmarshalBinary(data []byte) error {
	t.ProtocolIndex = data[0]
	t.ProtocolString = data[1:]
	return nil
}

type FIDBundleSeedIDPrefToken struct {
	BundleSeedIDString [11]byte
}

func (t *FIDBundleSeedIDPrefToken) UnmarshalBinary(data []byte) error {
	copy(t.BundleSeedIDString[:], data)
	return nil
}

type FIDScreenInfoToken struct {
	ScreenWidthInches  uint16
	ScreenHeightInches uint16
	ScreenWidthPixels  uint16
	ScreenHeightPixels uint16

	IpodScreenWidthPixels  uint16
	IpodScreenHeightPixels uint16

	ScreenFeaturesMask byte
	ScreenGammaValue   byte
}

func (t *FIDScreenInfoToken) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	return binary.Read(r, binary.BigEndian, t)
}

type FIDEAProtocolMetadataToken struct {
	ProtocolIndex byte
	MetadataType  byte
}

func (t *FIDEAProtocolMetadataToken) UnmarshalBinary(data []byte) error {
	t.ProtocolIndex = data[0]
	t.MetadataType = data[1]
	return nil
}

type FIDMicrophoneCapsToken struct {
	MicCapsBitmask uint32
}

func (t *FIDMicrophoneCapsToken) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	return binary.Read(r, binary.BigEndian, t)
}

type FIDTokenValue struct {
	Len        byte
	FIDType    byte
	FIDSubtype byte
	Token      interface{}
}

type SetFIDTokenValues struct {
	NumFIDTokenValues byte
	FIDTokenValues    []FIDTokenValue
}

func (s *SetFIDTokenValues) UnmarshalBinary(data []byte) error {
	br := bytes.NewReader(data)
	var err error
	s.NumFIDTokenValues, err = br.ReadByte()
	if err != nil {
		return err
	}
	s.FIDTokenValues = make([]FIDTokenValue, s.NumFIDTokenValues)
	for i := range s.FIDTokenValues {

		v := &s.FIDTokenValues[i]

		v.Len, err = br.ReadByte()
		if err != nil {
			return err
		}
		v.FIDType, err = br.ReadByte()
		if err != nil {
			return err
		}
		v.FIDSubtype, err = br.ReadByte()
		if err != nil {
			return err
		}
		data := make([]byte, v.Len-2)
		_, err = br.Read(data)
		if err != nil {
			return err
		}

		switch v.FIDType {
		case 0x00:
			switch v.FIDSubtype {
			case 0x00:
				//identify
				v.Token = &FIDIdentifyToken{}
			case 0x01:
				//acc caps
				v.Token = &FIDAccCapsToken{}
			case 0x02:
				//accinfo
				v.Token = &FIDAccInfoToken{}
			case 0x03:
				//ipod pref
				v.Token = &FIDiPodPreferenceToken{}
			case 0x04:
				//sdk proto
				v.Token = &FIDEAProtocolToken{}
			case 0x05:
				// bundleseed
				v.Token = &FIDBundleSeedIDPrefToken{}
			case 0x07:
				// screen info
				v.Token = &FIDScreenInfoToken{}
			case 0x08:
				// eaprotometadata
				v.Token = &FIDEAProtocolMetadataToken{}

			}
		case 0x01:
			//mic
			v.Token = &FIDMicrophoneCapsToken{}
		}

		if bu, ok := v.Token.(encoding.BinaryUnmarshaler); ok {
			if err := bu.UnmarshalBinary(data); err != nil {
				return err
			}
		} else {
			v.Token = data
		}

	}
	return nil

}

type RetFIDTokenValueACKs struct {
	NumFIDTokenValueACKs byte
	FIDTokenValueACKs    []byte
}

func (s RetFIDTokenValueACKs) MarshalBinary() ([]byte, error) {
	return append([]byte{s.NumFIDTokenValueACKs}, s.FIDTokenValueACKs...), nil

}

func (s *RetFIDTokenValueACKs) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return errors.New("RetFIDTokenValueACKs: short payload")
	}
	s.NumFIDTokenValueACKs = data[0]
	if len(data) > 1 {
		s.FIDTokenValueACKs = make([]byte, len(data[1:]))
		copy(s.FIDTokenValueACKs, data[1:])
	}
	return nil
}

type AccEndIDPSStatus uint8

const (
	AccEndIDPSStatusContinue AccEndIDPSStatus = 0x00
	AccEndIDPSStatusReset    AccEndIDPSStatus = 0x01
	AccEndIDPSStatusAbandon  AccEndIDPSStatus = 0x02
	AccEndIDPSStatusNewLink  AccEndIDPSStatus = 0x03
)

type EndIDPS struct {
	AccEndIDPSStatus AccEndIDPSStatus
}

type IDPSStatusEnum uint8

const (
	IDPSStatusOK                   IDPSStatusEnum = 0x00
	IDPSStatusTimeLimitNotExceeded IDPSStatusEnum = 0x04
	IDPSStatusWillNotAccept        IDPSStatusEnum = 0x06
)

type IDPSStatus struct {
	Status IDPSStatusEnum
}

type OpenDataSessionForProtocol struct {
	SessionID     uint16
	ProtocolIndex byte
}
type CloseDataSession struct {
	SessionID uint16
}
type DevACK struct {
	AckStatus byte
	CmdID     byte
}

type DevDataTransfer struct {
	SessionID uint16
	Data      []byte
}
type IPodDataTransfer struct {
	SessionID uint16
	Data      []byte
}
type SetAccStatusNotification struct {
	StatusMask uint32
}
type RetAccStatusNotification struct {
	StatusMask uint32
}
type AccessoryStatusNotification struct {
	StatusType   byte
	StatusParams []byte
}

type SetEventNotification struct {
	EventMask uint64
}
type IPodNotification struct {
	NotificationType byte
	Data             []byte
}

type GetiPodOptionsForLingo struct {
	LingoID byte
}
type RetiPodOptionsForLingo struct {
	LingoID byte
	Options uint64
}
type GetEventNotification struct{}

type RetEventNotification struct {
	EventMask uint64
}
type GetSupportedEventNotification struct{}

type CancelCommand struct {
	LingoID       byte
	CmdID         uint16
	TransactionID uint16
}
type RetSupportedEventNotification struct {
	EventMask uint64
}
type SetAvailableCurrent struct {
	CurrentLimit uint16
}
type RequestApplicationLaunch struct {
	_     [3]byte
	AppID []byte
}
type GetNowPlayingFocusApp struct{}

type RetNowPlayingFocusApp struct {
	AppID []byte
}
