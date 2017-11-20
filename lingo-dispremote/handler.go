package dispremote

import (
	"git.andrewo.pw/andrew/ipod"
)

type DeviceDispRemote interface {
}

func ackSuccess(req ipod.Packet) ACK {
	return ACK{Status: ACKStatusSuccess, CmdID: uint8(req.ID.CmdID())}
}

func HandleDispRemote(req ipod.Packet, tr ipod.PacketWriter, dev DeviceDispRemote) error {
	switch msg := req.Payload.(type) {

	case GetCurrentEQProfileIndex:
		ipod.Respond(req, tr, RetCurrentEQProfileIndex{
			CurrentEQIndex: 0,
		})

	case SetCurrentEQProfileIndex:
		ipod.Respond(req, tr, ackSuccess(req))

	case GetNumEQProfiles:
		ipod.Respond(req, tr, RetNumEQProfiles{
			NumEQProfiles: 1,
		})
	case GetIndexedEQProfileName:
		ipod.Respond(req, tr, RetIndexedEQProfileName{
			EQProfileName: ipod.StringToBytes("Default"),
		})
	case SetRemoteEventNotification:
		ipod.Respond(req, tr, ackSuccess(req))

	case GetRemoteEventStatus:
		ipod.Respond(req, tr, RetRemoteEventStatus{
			EventStatus: 0,
		})

	case GetiPodStateInfo:
	// ipod.Respond(req, tr, RetiPodStateInfo{
	// })
	// todo

	case SetiPodStateInfo:
		ipod.Respond(req, tr, ackSuccess(req))

	case GetPlayStatus:
		ipod.Respond(req, tr, RetPlayStatus{
			PlayState: 0, //stopped
		})

	case SetCurrentPlayingTrack:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetIndexedPlayingTrackInfo:
		// RetIndexedPlayingTrackInfo:
		ipod.Respond(req, tr, RetIndexedPlayingTrackInfo{
			InfoType: msg.InfoType,
			InfoData: []byte{0x00}, //no data
		})
	case GetNumPlayingTracks:
		ipod.Respond(req, tr, RetNumPlayingTracks{
			NumPlayTracks: 0,
		})
	case GetArtworkFormats:
		ipod.Respond(req, tr, RetArtworkFormats{})
	case GetTrackArtworkData:
	// RetTrackArtworkData:
	//todo
	case GetPowerBatteryState:
		ipod.Respond(req, tr, RetPowerBatteryState{
			BatteryLevel: 255, // 100%
			PowerState:   0x01,
		})
	case GetSoundCheckState:
		ipod.Respond(req, tr, RetSoundCheckState{
			Enabled: false,
		})
	case SetSoundCheckState:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetTrackArtworkTimes:
		ipod.Respond(req, tr, RetTrackArtworkTimes{})

	default:
		_ = msg
	}
	return nil
}
