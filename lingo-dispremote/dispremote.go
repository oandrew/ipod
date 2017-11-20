package dispremote

import (
	"git.andrewo.pw/andrew/ipod"
)

func init() {
	ipod.RegisterLingos(LingoDisplayRemoteID, Lingos)
}

const LingoDisplayRemoteID = 0x03

var Lingos struct {
	ACK                        `id:"0x00"`
	GetCurrentEQProfileIndex   `id:"0x01"`
	RetCurrentEQProfileIndex   `id:"0x02"`
	SetCurrentEQProfileIndex   `id:"0x03"`
	GetNumEQProfiles           `id:"0x04"`
	RetNumEQProfiles           `id:"0x05"`
	GetIndexedEQProfileName    `id:"0x06"`
	RetIndexedEQProfileName    `id:"0x07"`
	SetRemoteEventNotification `id:"0x08"`
	RemoteEventNotification    `id:"0x09"`
	GetRemoteEventStatus       `id:"0x0A"`
	RetRemoteEventStatus       `id:"0x0B"`
	GetiPodStateInfo           `id:"0x0C"`
	RetiPodStateInfo           `id:"0x0D"`
	SetiPodStateInfo           `id:"0x0E"`
	GetPlayStatus              `id:"0x0F"`
	RetPlayStatus              `id:"0x10"`
	SetCurrentPlayingTrack     `id:"0x11"`
	GetIndexedPlayingTrackInfo `id:"0x12"`
	RetIndexedPlayingTrackInfo `id:"0x13"`
	GetNumPlayingTracks        `id:"0x14"`
	RetNumPlayingTracks        `id:"0x15"`
	GetArtworkFormats          `id:"0x16"`
	RetArtworkFormats          `id:"0x17"`
	GetTrackArtworkData        `id:"0x18"`
	RetTrackArtworkData        `id:"0x19"`
	GetPowerBatteryState       `id:"0x1A"`
	RetPowerBatteryState       `id:"0x1B"`
	GetSoundCheckState         `id:"0x1C"`
	RetSoundCheckState         `id:"0x1D"`
	SetSoundCheckState         `id:"0x1E"`
	GetTrackArtworkTimes       `id:"0x1F"`
	RetTrackArtworkTimes       `id:"0x20"`
}

type ACKStatus uint8

const (
	ACKStatusSuccess ACKStatus = 0x00
	ACKStatusPending ACKStatus = 0x06
)

type ACK struct {
	Status ACKStatus
	CmdID  uint8
}
type GetCurrentEQProfileIndex struct {
}
type RetCurrentEQProfileIndex struct {
	CurrentEQIndex uint32
}
type SetCurrentEQProfileIndex struct {
	CurrentEQIndex uint32
	RestoreOnExit  bool
}
type GetNumEQProfiles struct {
}
type RetNumEQProfiles struct {
	NumEQProfiles uint32
}
type GetIndexedEQProfileName struct {
	EQProfileIndex uint32
}
type RetIndexedEQProfileName struct {
	EQProfileName []byte
}
type SetRemoteEventNotification struct {
	EventMask uint32
}
type RemoteEventNotification struct {
	EventNum  byte
	EventData []byte
}
type GetRemoteEventStatus struct {
}
type RetRemoteEventStatus struct {
	EventStatus uint32
}
type GetiPodStateInfo struct {
	InfoType byte
}
type RetiPodStateInfo struct {
	InfoType byte
	InfoData []byte
}
type SetiPodStateInfo struct {
	InfoType byte
	InfoData byte // todo
}
type GetPlayStatus struct {
}
type RetPlayStatus struct {
	PlayState   byte
	TrackIndex  uint32
	TrackLength uint32
	TrackPos    uint32
}
type SetCurrentPlayingTrack struct {
	TrackIndex uint32
}
type GetIndexedPlayingTrackInfo struct {
	InfoType     byte
	TrackIndex   uint32
	ChapterIndex uint16
}
type RetIndexedPlayingTrackInfo struct {
	InfoType byte
	InfoData []byte
}
type GetNumPlayingTracks struct {
}
type RetNumPlayingTracks struct {
	NumPlayTracks uint32
}
type GetArtworkFormats struct {
}

type ArtworkFormat struct {
	FormatID    uint16
	PixelFormat uint8
	ImageWidth  uint16
	ImageHeight uint16
}
type RetArtworkFormats struct {
	Formats []ArtworkFormat
}
type GetTrackArtworkData struct {
	TrackIndex uint32
	FormatID   uint16
	TimeOffset uint32
}
type RetTrackArtworkData struct {
	//todo
}
type GetPowerBatteryState struct {
}
type RetPowerBatteryState struct {
	PowerState   byte
	BatteryLevel uint8
}
type GetSoundCheckState struct {
}
type RetSoundCheckState struct {
	Enabled bool
}
type SetSoundCheckState struct {
	Enabled       bool
	RestoreOnExit bool
}
type GetTrackArtworkTimes struct {
	TrackIndex   uint32
	FormatID     uint16
	ArtworkIndex uint16
	ArtworkCount uint16
}
type RetTrackArtworkTimes struct {
	TimeOffset []uint32
}
