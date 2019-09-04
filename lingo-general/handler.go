package general

import (
	"bytes"

	"github.com/oandrew/ipod"
)

type DeviceGeneral interface {
	UIMode() UIMode
	SetUIMode(UIMode)
	Name() string
	SoftwareVersion() (major, minor, rev uint8)
	SerialNum() string

	LingoProtocolVersion(lingo uint8) (major, minor uint8)
	LingoOptions(ling uint8) uint64

	PrefSettingID(classID uint8) uint8
	SetPrefSettingID(classID, settingID uint8, restoreOnExit bool)

	StartIDPS()
	EndIDPS(status AccEndIDPSStatus)
	SetToken(token FIDTokenValue) error
	AccAuthCert(cert []byte)

	SetEventNotificationMask(mask uint64)
	EventNotificationMask() uint64
	SupportedEventNotificationMask() uint64

	CancelCommand(lingo uint8, cmd uint16, transaction uint16)

	MaxPayload() uint16
}

func ackSuccess(req *ipod.Command) *ACK {
	return &ACK{Status: ACKStatusSuccess, CmdID: uint8(req.ID.CmdID())}
}

func ackPending(req *ipod.Command, maxWait uint32) *ACKPending {
	return &ACKPending{Status: ACKStatusPending, CmdID: uint8(req.ID.CmdID()), MaxWait: maxWait}
}

func ack(req *ipod.Command, status ACKStatus) *ACK {
	return &ACK{Status: status, CmdID: uint8(req.ID.CmdID())}
}

func ackFIDTokens(tokens *SetFIDTokenValues) *RetFIDTokenValueACKs {
	resp := &RetFIDTokenValueACKs{NumFIDTokenValueACKs: tokens.NumFIDTokenValues}
	buf := bytes.Buffer{}
	for _, token := range tokens.FIDTokenValues {

		//after subtype
		ackBuf := bytes.Buffer{}
		ackBuf.Write([]byte{token.FIDType, token.FIDSubtype})

		switch t := token.Token.(type) {
		case *FIDIdentifyToken:
			ackBuf.Write([]byte{0x00})
		case *FIDAccCapsToken:
			ackBuf.Write([]byte{0x00})
		case *FIDAccInfoToken:
			ackBuf.Write([]byte{0x00, t.AccInfoType})
		case *FIDiPodPreferenceToken:
			ackBuf.Write([]byte{0x00, t.PrefClass})
		case *FIDEAProtocolToken:
			ackBuf.Write([]byte{0x00, t.ProtocolIndex})
		case *FIDBundleSeedIDPrefToken:
			ackBuf.Write([]byte{0x00})
		case *FIDScreenInfoToken:
			ackBuf.Write([]byte{0x00})
		case *FIDEAProtocolMetadataToken:
			ackBuf.Write([]byte{0x00})

		case *FIDMicrophoneCapsToken:
			ackBuf.Write([]byte{0x00})
		}
		buf.WriteByte(byte(ackBuf.Len()))
		buf.ReadFrom(&ackBuf)

	}
	resp.FIDTokenValueACKs = buf.Bytes()
	return resp
}

var accCertBuf bytes.Buffer

func HandleGeneral(req *ipod.Command, tr ipod.CommandWriter, dev DeviceGeneral) error {
	switch msg := req.Payload.(type) {
	case *RequestRemoteUIMode:
		ipod.Respond(req, tr, &ReturnRemoteUIMode{
			Mode: ipod.BoolToByte(dev.UIMode() == UIModeExtended),
		})
	case *EnterRemoteUIMode:
		if dev.UIMode() == UIModeExtended {
			ipod.Respond(req, tr, ackSuccess(req))
		} else {
			ipod.Respond(req, tr, ackPending(req, 300))
			dev.SetUIMode(UIModeExtended)
			ipod.Respond(req, tr, ackSuccess(req))
		}
	case *ExitRemoteUIMode:
		if dev.UIMode() != UIModeExtended {
			ipod.Respond(req, tr, ackSuccess(req))
		} else {
			ipod.Respond(req, tr, ackPending(req, 300))
			dev.SetUIMode(UIModeStandart)
			ipod.Respond(req, tr, ackSuccess(req))
		}
	case *RequestiPodName:
		ipod.Respond(req, tr, &ReturniPodName{Name: ipod.StringToBytes(dev.Name())})
	case *RequestiPodSoftwareVersion:
		var resp ReturniPodSoftwareVersion
		resp.Major, resp.Minor, resp.Rev = dev.SoftwareVersion()
		ipod.Respond(req, tr, &resp)
	case *RequestiPodSerialNum:
		ipod.Respond(req, tr, &ReturniPodSerialNum{Serial: ipod.StringToBytes(dev.SerialNum())})
	case *RequestLingoProtocolVersion:
		var resp ReturnLingoProtocolVersion
		resp.Lingo = msg.Lingo
		resp.Major, resp.Minor = dev.LingoProtocolVersion(msg.Lingo)
		ipod.Respond(req, tr, &resp)
	case *RequestTransportMaxPayloadSize:
		ipod.Respond(req, tr, &ReturnTransportMaxPayloadSize{MaxPayload: dev.MaxPayload()})
	case *IdentifyDeviceLingoes:
		ipod.Respond(req, tr, ackSuccess(req))
		ipod.Respond(req, tr, &GetDevAuthenticationInfo{})

	//GetDevAuthenticationInfo
	case *RetDevAuthenticationInfo:
		if msg.Major >= 2 {
			if msg.CertCurrentSection == 0 {
				accCertBuf.Reset()
			}
			accCertBuf.Write(msg.CertData)
			if msg.CertCurrentSection < msg.CertMaxSection {
				ipod.Respond(req, tr, ackSuccess(req))
			} else {
				ipod.Respond(req, tr, &AckDevAuthenticationInfo{Status: DevAuthInfoStatusSupported})
				dev.AccAuthCert(accCertBuf.Bytes())
				ipod.Respond(req, tr, &GetDevAuthenticationSignatureV2{Counter: 0})
			}
		} else {
			ipod.Respond(req, tr, &AckDevAuthenticationInfo{Status: DevAuthInfoStatusSupported})
		}

	// GetDevAuthenticationSignatureV1
	// case *RetDevAuthenticationSignatureV1:
	// 	ipod.Respond(req, tr, ackDevAuthenticationStatus{Status: DevAuthStatusPassed})
	// // GetDevAuthenticationSignatureV2
	// case *RetDevAuthenticationSignatureV2:
	// 	ipod.Respond(req, tr, ackDevAuthenticationStatus{Status: DevAuthStatusPassed})

	case *RetDevAuthenticationSignature:
		ipod.Respond(req, tr, &AckDevAuthenticationStatus{Status: DevAuthStatusPassed})

	case *GetiPodAuthenticationInfo:
		ipod.Respond(req, tr, &RetiPodAuthenticationInfo{
			Major: 1, Minor: 1,
			CertCurrentSection: 0, CertMaxSection: 0, CertData: []byte{},
		})

	case *AckiPodAuthenticationInfo:
		// pass

	case *GetiPodAuthenticationSignature:
		ipod.Respond(req, tr, &RetiPodAuthenticationSignature{Signature: msg.Challenge})

	case *AckiPodAuthenticationStatus:
		// pass

	// revisit
	case *GetiPodOptions:
		ipod.Respond(req, tr, &RetiPodOptions{Options: 0x00})

	// GetAccessoryInfo
	// check back might be useful
	case *RetAccessoryInfo:
		// pass

	case *GetiPodPreferences:
		ipod.Respond(req, tr, &RetiPodPreferences{
			PrefClassID:        msg.PrefClassID,
			PrefClassSettingID: dev.PrefSettingID(msg.PrefClassID),
		})

	case *SetiPodPreferences:
		dev.SetPrefSettingID(msg.PrefClassID, msg.PrefClassSettingID, ipod.ByteToBool(msg.RestoreOnExit))
		ipod.Respond(req, tr, ackSuccess(req))

	case *GetUIMode:
		ipod.Respond(req, tr, &RetUIMode{UIMode: dev.UIMode()})
	case *SetUIMode:
		ipod.Respond(req, tr, ackSuccess(req))

	case *StartIDPS:
		ipod.TrxReset()
		dev.StartIDPS()
		ipod.Respond(req, tr, ackSuccess(req))
	case *SetFIDTokenValues:
		for _, token := range msg.FIDTokenValues {
			dev.SetToken(token)
		}
		ipod.Respond(req, tr, ackFIDTokens(msg))
	case *EndIDPS:
		dev.EndIDPS(msg.AccEndIDPSStatus)
		switch msg.AccEndIDPSStatus {
		case AccEndIDPSStatusContinue:
			ipod.Respond(req, tr, &IDPSStatus{Status: IDPSStatusOK})
			ipod.Send(tr, &GetDevAuthenticationInfo{})

			// get dev auth info
		case AccEndIDPSStatusReset:
			ipod.Respond(req, tr, &IDPSStatus{Status: IDPSStatusTimeLimitNotExceeded})
		case AccEndIDPSStatusAbandon:
			ipod.Respond(req, tr, &IDPSStatus{Status: IDPSStatusWillNotAccept})
		case AccEndIDPSStatusNewLink:
			//pass
		}

	// SetAccStatusNotification, RetAccStatusNotification
	case *AccessoryStatusNotification:

	// iPodNotification later
	case *SetEventNotification:
		dev.SetEventNotificationMask(msg.EventMask)
		ipod.Respond(req, tr, ackSuccess(req))

	case *GetiPodOptionsForLingo:
		ipod.Respond(req, tr, &RetiPodOptionsForLingo{
			LingoID: msg.LingoID,
			Options: dev.LingoOptions(msg.LingoID),
		})

	case *GetEventNotification:
		ipod.Respond(req, tr, &RetEventNotification{
			EventMask: dev.EventNotificationMask(),
		})

	case *GetSupportedEventNotification:
		ipod.Respond(req, tr, &RetSupportedEventNotification{
			EventMask: dev.SupportedEventNotificationMask(),
		})

	case *CancelCommand:
		dev.CancelCommand(msg.LingoID, msg.CmdID, msg.TransactionID)
		ipod.Respond(req, tr, ackSuccess(req))

	case *SetAvailableCurrent:
		// notify acc

	case *RequestApplicationLaunch:
		ipod.Respond(req, tr, ackSuccess(req))

	case *GetNowPlayingFocusApp:
		ipod.Respond(req, tr, &RetNowPlayingFocusApp{AppID: ipod.StringToBytes("")})

	case ipod.UnknownPayload:
		ipod.Respond(req, tr, ack(req, ACKStatusUnkownID))
	default:
		_ = msg
	}
	return nil
}
