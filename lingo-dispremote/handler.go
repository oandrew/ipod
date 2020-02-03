package dispremote

import (
	"errors"
	"time"

	"github.com/oandrew/ipod"
)

type DeviceDispRemote interface {
}

func ackSuccess(req *ipod.Command) *ACK {
	return &ACK{Status: ACKStatusSuccess, CmdID: uint8(req.ID.CmdID())}
}

func HandleDispRemote(req *ipod.Command, tr ipod.CommandWriter, dev DeviceDispRemote) error {
	switch msg := req.Payload.(type) {

	case *GetCurrentEQProfileIndex:
		ipod.Respond(req, tr, &RetCurrentEQProfileIndex{
			CurrentEQIndex: 0,
		})

	case *SetCurrentEQProfileIndex:
		ipod.Respond(req, tr, ackSuccess(req))

	case *GetNumEQProfiles:
		ipod.Respond(req, tr, &RetNumEQProfiles{
			NumEQProfiles: 1,
		})
	case *GetIndexedEQProfileName:
		ipod.Respond(req, tr, &RetIndexedEQProfileName{
			EQProfileName: ipod.StringToBytes("Default"),
		})
	case *SetRemoteEventNotification:
		ipod.Respond(req, tr, ackSuccess(req))

	case *GetRemoteEventStatus:
		ipod.Respond(req, tr, &RetRemoteEventStatus{
			EventStatus: 0,
		})

	case *GetiPodStateInfo:
		t := &RetiPodStateInfo{
			InfoType: msg.InfoType,
		}

		switch msg.InfoType {
		case InfoTypeTrackPositionMs:
			t.InfoData = &InfoTrackPositionMs{TrackPositionMs: 0}
		case InfoTypeTrackIndex:
			t.InfoData = &InfoTrackIndex{TrackIndex: 1}
		case InfoTypeChapterIndex:
			t.InfoData = &InfoChapterIndex{
				TrackIndex:   1,
				ChapterCount: 0,
				ChapterIndex: 0,
			}
		case InfoTypePlayStatus:
			t.InfoData = &InfoPlayStatus{
				PlayStatus: PlayStatusPlaying,
			}
		case InfoTypeVolume:
			t.InfoData = &InfoVolume{MuteState: 0x00, UIVolumeLevel: 255}
		case InfoTypePower:
			t.InfoData = &InfoPower{PowerState: 0x05, BatteryLevel: 255}
		case InfoTypeEqualizer:
			t.InfoData = &InfoEqualizer{0x00}
		case InfoTypeShuffle:
			t.InfoData = &InfoShuffle{0x00}
		case InfoTypeRepeat:
			t.InfoData = &InfoRepeat{0x00}
		case InfoTypeDateTime:
			d := time.Now()
			t.InfoData = &InfoDateTime{
				Year:   uint16(d.Year()),
				Month:  uint8(d.Month()),
				Day:    uint8(d.Day()),
				Hour:   uint8(d.Hour()),
				Minute: uint8(d.Minute()),
			}
		case InfoTypeBacklight:
			t.InfoData = &InfoBacklight{BacklightLevel: 255}
		case InfoTypeHoldSwitch:
			t.InfoData = &InfoHoldSwitch{HoldSwitchState: 0x00}
		case InfoTypeSoundCheck:
			t.InfoData = &InfoSoundCheck{SoundCheckState: 0x00}
		case InfoTypeAudiobookSpeed:
			t.InfoData = &InfoAudiobookSpeed{0x00}
		case InfoTypeTrackPositionSec:
			t.InfoData = &InfoTrackPositionSec{0}
		case InfoTypeVolume2:
			t.InfoData = &InfoVolume2{
				MuteState:           0x00,
				UIVolumeLevel:       255,
				AbsoluteVolumeLevel: 255,
			}
		default:
			return errors.New("unknown info type")
		}

		ipod.Respond(req, tr, t)

	case *SetiPodStateInfo:
		ipod.Respond(req, tr, ackSuccess(req))

	case *GetPlayStatus:
		ipod.Respond(req, tr, &RetPlayStatus{
			PlayState: 0, //stopped
		})

	case *SetCurrentPlayingTrack:
		ipod.Respond(req, tr, ackSuccess(req))
	case *GetIndexedPlayingTrackInfo:
		// RetIndexedPlayingTrackInfo:
		ipod.Respond(req, tr, &RetIndexedPlayingTrackInfo{
			InfoType: msg.InfoType,
			InfoData: 0x00, //todo
		})
	case *GetNumPlayingTracks:
		ipod.Respond(req, tr, &RetNumPlayingTracks{
			NumPlayTracks: 0,
		})
	case *GetArtworkFormats:
		ipod.Respond(req, tr, &RetArtworkFormats{})
	case *GetTrackArtworkData:
	// RetTrackArtworkData:
	//todo
	case *GetPowerBatteryState:
		ipod.Respond(req, tr, &RetPowerBatteryState{
			BatteryLevel: 255, // 100%
			PowerState:   0x01,
		})
	case *GetSoundCheckState:
		ipod.Respond(req, tr, &RetSoundCheckState{
			Enabled: false,
		})
	case *SetSoundCheckState:
		ipod.Respond(req, tr, ackSuccess(req))
	case *GetTrackArtworkTimes:
		ipod.Respond(req, tr, &RetTrackArtworkTimes{})

	default:
		_ = msg
	}
	return nil
}
