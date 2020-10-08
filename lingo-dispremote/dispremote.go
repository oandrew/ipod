package dispremote

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/oandrew/ipod"
)

func init() {
	ipod.RegisterLingos(ipod.LingoDisplayRemoteID, Lingos)
}

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

//go:generate stringer -type=InfoType
type InfoType uint8

const (
	InfoTypeTrackPositionMs InfoType = iota
	InfoTypeTrackIndex
	InfoTypeChapterIndex
	InfoTypePlayStatus
	InfoTypeVolume
	InfoTypePower
	InfoTypeEqualizer
	InfoTypeShuffle
	InfoTypeRepeat
	InfoTypeDateTime
	_ //InfoTypeAlarm
	InfoTypeBacklight
	InfoTypeHoldSwitch
	InfoTypeSoundCheck
	InfoTypeAudiobookSpeed
	InfoTypeTrackPositionSec
	InfoTypeVolume2
)

type InfoTrackPositionMs struct {
	TrackPositionMs uint32
}
type InfoTrackIndex struct {
	TrackIndex uint32
}
type InfoChapterIndex struct {
	TrackIndex   uint32
	ChapterCount uint16
	ChapterIndex uint16
}

//go:generate stringer -type=PlayStatusType
type PlayStatusType uint8

const (
	PlayStatusStopped PlayStatusType = iota
	PlayStatusPlaying
	PlayStatusPaused
	PlayStatusFF
	PlayStatusREW
	PlayStatusEndFFREW
)

type InfoPlayStatus struct {
	PlayStatus PlayStatusType
}
type InfoVolume struct {
	MuteState     uint8
	UIVolumeLevel uint8
}
type InfoPower struct {
	PowerState   uint8
	BatteryLevel uint8
}
type InfoEqualizer struct {
	EqIndex uint32
}
type InfoShuffle struct {
	ShuffleState uint8
}
type InfoRepeat struct {
	RepeatState uint8
}
type InfoDateTime struct {
	Year   uint16
	Month  uint8
	Day    uint8
	Hour   uint8
	Minute uint8
}
type InfoBacklight struct {
	BacklightLevel uint8
}
type InfoHoldSwitch struct {
	HoldSwitchState uint8
}
type InfoSoundCheck struct {
	SoundCheckState uint8
}
type InfoAudiobookSpeed struct {
	PlaybackSpeed uint8
}
type InfoTrackPositionSec struct {
	TrackPositionSec uint16
}
type InfoVolume2 struct {
	MuteState           uint8
	UIVolumeLevel       uint8
	AbsoluteVolumeLevel uint8
}

type GetiPodStateInfo struct {
	InfoType InfoType
}
type RetiPodStateInfo struct {
	InfoType InfoType
	InfoData interface{}
}

func (t *RetiPodStateInfo) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, t.InfoType)
	binary.Write(&buf, binary.BigEndian, t.InfoData)
	return buf.Bytes(), nil
}

func (t *RetiPodStateInfo) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &t.InfoType)
	switch t.InfoType {
	case InfoTypeTrackPositionMs:
		t.InfoData = &InfoTrackPositionMs{}
	case InfoTypeTrackIndex:
		t.InfoData = &InfoTrackIndex{}
	case InfoTypeChapterIndex:
		t.InfoData = &InfoChapterIndex{}
	case InfoTypePlayStatus:
		t.InfoData = &InfoPlayStatus{}
	case InfoTypeVolume:
		t.InfoData = &InfoVolume{}
	case InfoTypePower:
		t.InfoData = &InfoPower{}
	case InfoTypeEqualizer:
		t.InfoData = &InfoEqualizer{}
	case InfoTypeShuffle:
		t.InfoData = &InfoShuffle{}
	case InfoTypeRepeat:
		t.InfoData = &InfoRepeat{}
	case InfoTypeDateTime:
		t.InfoData = &InfoDateTime{}
	case InfoTypeBacklight:
		t.InfoData = &InfoBacklight{}
	case InfoTypeHoldSwitch:
		t.InfoData = &InfoHoldSwitch{}
	case InfoTypeSoundCheck:
		t.InfoData = &InfoSoundCheck{}
	case InfoTypeAudiobookSpeed:
		t.InfoData = &InfoAudiobookSpeed{}
	case InfoTypeTrackPositionSec:
		t.InfoData = &InfoTrackPositionSec{}
	case InfoTypeVolume2:
		t.InfoData = &InfoVolume2{}
	default:
		return errors.New("unknown info type")
	}
	return binary.Read(r, binary.BigEndian, t.InfoData)
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

//go:generate stringer -type=TrackInfoType
type TrackInfoType uint8

const (
	TrackInfoTypeCaps TrackInfoType = iota
	TrackInfoTypeChapterTimeName
	TrackInfoTypeArtist
	TrackInfoTypeAlbum
	TrackInfoTypeGenre
	TrackInfoTypeTrack
	TrackInfoTypeComposer
	TrackInfoTypeLyrics
	TrackInfoTypeArtworkCount
)

type GetIndexedPlayingTrackInfo struct {
	InfoType     TrackInfoType
	TrackIndex   uint32
	ChapterIndex uint16
}

type TrackInfoCaps struct {
	Caps         uint32
	TrackTotalMs uint32
	ChapterCount uint16
}
type TrackInfoChapterTimeName struct {
	ChapterTime uint32
	ChapterName []byte
}
type TrackInfoArtist struct {
	Name []byte
}
type TrackInfoAlbum struct {
	Name []byte
}
type TrackInfoGenre struct {
	Name []byte
}
type TrackInfoTrack struct {
	Title []byte
}
type TrackInfoComposer struct {
	Name []byte
}
type TrackInfoLyrics struct {
	Flags       uint8
	PacketIndex uint16
	Lyrics      []byte
}
type TrackInfoArtworkCount struct {
	None byte // empty = 0x08
}

type RetIndexedPlayingTrackInfo struct {
	InfoType TrackInfoType
	InfoData interface{}
}

func (t *RetIndexedPlayingTrackInfo) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, t.InfoType)
	binary.Write(&buf, binary.BigEndian, t.InfoData)
	return buf.Bytes(), nil
}

func (t *RetIndexedPlayingTrackInfo) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &t.InfoType)
	switch t.InfoType {
	case TrackInfoTypeCaps:
		t.InfoData = &TrackInfoCaps{}
	case TrackInfoTypeChapterTimeName:
		t.InfoData = &TrackInfoChapterTimeName{
			ChapterTime: 0,
			ChapterName: make([]byte, 0),
		}
	case TrackInfoTypeArtist:
		t.InfoData = &TrackInfoArtist{
			Name: make([]byte, 0),
		}
	case TrackInfoTypeAlbum:
		t.InfoData = &TrackInfoAlbum{
			Name: make([]byte, 0),
		}
	case TrackInfoTypeGenre:
		t.InfoData = &TrackInfoGenre{
			Name: make([]byte, 0),
		}
	case TrackInfoTypeTrack:
		t.InfoData = &TrackInfoTrack{
			Title: make([]byte, 0),
		}
	case TrackInfoTypeComposer:
		t.InfoData = &TrackInfoComposer{
			Name: make([]byte, 0),
		}
	case TrackInfoTypeLyrics:
		t.InfoData = &TrackInfoLyrics{
			Flags:       0x00,
			PacketIndex: 0,
			Lyrics:      make([]byte, 0),
		}
	case TrackInfoTypeArtworkCount:
		t.InfoData = &TrackInfoArtworkCount{}
	default:
		return errors.New("unknown info type")
	}
	return binary.Read(r, binary.BigEndian, t.InfoData)
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
