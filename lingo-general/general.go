package general

import (
	"bufio"
	"io"
	"io/ioutil"

	"git.andrewo.pw/andrew/ipod"
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
	ACKStatusSuccess ACKStatus = 0x00
	ACKStatusPending ACKStatus = 0x06
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

func (s ReturniPodName) MarshalPayload(w io.Writer) error {
	w.Write(s.Name)
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

func (s ReturniPodSerialNum) MarshalPayload(w io.Writer) error {
	w.Write(s.Serial)
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

func (s *RetDevAuthenticationInfo) UnmarshalPayload(r io.Reader) error {
	br := bufio.NewReader(r)
	var err error
	s.Major, err = br.ReadByte()
	if err != nil {
		return err
	}

	s.Minor, err = br.ReadByte()
	if err != nil {
		return err
	}

	if s.Major >= 0x02 {
		s.CertCurrentSection, err = br.ReadByte()
		if err != nil {
			return err
		}

		s.CertMaxSection, err = br.ReadByte()
		if err != nil {
			return err
		}

		s.CertData, err = ioutil.ReadAll(br)
		if err != nil {
			return err
		}
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

func (s *RetDevAuthenticationSignature) UnmarshalPayload(r io.Reader) error {
	var err error
	s.Signature, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}
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

type FIDTokenValue struct {
	Len        byte
	FIDType    byte
	FIDSubtype byte
	Data       []byte
}

type SetFIDTokenValues struct {
	NumFIDTokenValues byte
	FIDTokenValues    []FIDTokenValue
}

func (s *SetFIDTokenValues) UnmarshalPayload(r io.Reader) error {
	br := bufio.NewReader(r)
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
		v.Data = make([]byte, v.Len-2)
		_, err = br.Read(v.Data)
		if err != nil {
			return err
		}
	}
	return nil

}

type RetFIDTokenValueACKs struct {
	NumFIDTokenValueACKs byte
	FIDTokenValueACKs    []byte
}

func (s RetFIDTokenValueACKs) MarshalPayload(w io.Writer) error {
	w.Write([]byte{s.NumFIDTokenValueACKs})
	w.Write(s.FIDTokenValueACKs)
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
