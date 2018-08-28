package extremote

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/oandrew/ipod"
)

func init() {
	ipod.RegisterLingos(LingoExtRemotelID, Lingos)
}

const LingoExtRemotelID = 0x04

var Lingos struct {
	ACK                                        `id:"0x0001"`
	GetCurrentPlayingTrackChapterInfo          `id:"0x0002"`
	ReturnCurrentPlayingTrackChapterInfo       `id:"0x0003"`
	SetCurrentPlayingTrackChapter              `id:"0x0004"`
	GetCurrentPlayingTrackChapterPlayStatus    `id:"0x0005"`
	ReturnCurrentPlayingTrackChapterPlayStatus `id:"0x0006"`
	GetCurrentPlayingTrackChapterName          `id:"0x0007"`
	ReturnCurrentPlayingTrackChapterName       `id:"0x0008"`
	GetAudiobookSpeed                          `id:"0x0009"`
	ReturnAudiobookSpeed                       `id:"0x000A"`
	SetAudiobookSpeed                          `id:"0x000B"`
	GetIndexedPlayingTrackInfo                 `id:"0x000C"`
	ReturnIndexedPlayingTrackInfo              `id:"0x000D"`
	GetArtworkFormats                          `id:"0x000E"`
	RetArtworkFormats                          `id:"0x000F"`
	GetTrackArtworkData                        `id:"0x0010"`
	RetTrackArtworkData                        `id:"0x0011"`
	ResetDBSelection                           `id:"0x0016"`
	SelectDBRecord                             `id:"0x0017"`
	GetNumberCategorizedDBRecords              `id:"0x0018"`
	ReturnNumberCategorizedDBRecords           `id:"0x0019"`
	RetrieveCategorizedDatabaseRecords         `id:"0x001A"`
	ReturnCategorizedDatabaseRecord            `id:"0x001B"`
	GetPlayStatus                              `id:"0x001C"`
	ReturnPlayStatus                           `id:"0x001D"`
	GetCurrentPlayingTrackIndex                `id:"0x001E"`
	ReturnCurrentPlayingTrackIndex             `id:"0x001F"`
	GetIndexedPlayingTrackTitle                `id:"0x0020"`
	ReturnIndexedPlayingTrackTitle             `id:"0x0021"`
	GetIndexedPlayingTrackArtistName           `id:"0x0022"`
	ReturnIndexedPlayingTrackArtistName        `id:"0x0023"`
	GetIndexedPlayingTrackAlbumName            `id:"0x0024"`
	ReturnIndexedPlayingTrackAlbumName         `id:"0x0025"`
	SetPlayStatusChangeNotification            `id:"0x0026"`
	SetPlayStatusChangeNotificationShort       `id:"0x0026"`
	PlayStatusChangeNotification               `id:"0x0027"`
	PlayCurrentSelection                       `id:"0x0028"`
	PlayControl                                `id:"0x0029"`
	GetTrackArtworkTimes                       `id:"0x002A"`
	RetTrackArtworkTimes                       `id:"0x002B"`
	GetShuffle                                 `id:"0x002C"`
	ReturnShuffle                              `id:"0x002D"`
	SetShuffle                                 `id:"0x002E"`
	GetRepeat                                  `id:"0x002F"`
	ReturnRepeat                               `id:"0x0030"`
	SetRepeat                                  `id:"0x0031"`
	SetDisplayImage                            `id:"0x0032"`
	GetMonoDisplayImageLimits                  `id:"0x0033"`
	ReturnMonoDisplayImageLimits               `id:"0x0034"`
	GetNumPlayingTracks                        `id:"0x0035"`
	ReturnNumPlayingTracks                     `id:"0x0036"`
	SetCurrentPlayingTrack                     `id:"0x0037"`
	SelectSortDBRecord                         `id:"0x0038"`
	GetColorDisplayImageLimits                 `id:"0x0039"`
	ReturnColorDisplayImageLimits              `id:"0x003A"`
	ResetDBSelectionHierarchy                  `id:"0x003B"`
	GetDBiTunesInfo                            `id:"0x003C"`
	RetDBiTunesInfo                            `id:"0x003D"`
	GetUIDTrackInfo                            `id:"0x003E"`
	RetUIDTrackInfo                            `id:"0x003F"`
	GetDBTrackInfo                             `id:"0x0040"`
	RetDBTrackInfo                             `id:"0x0041"`
	GetPBTrackInfo                             `id:"0x0042"`
	RetPBTrackInfo                             `id:"0x0043"`
}

type ACKStatus uint8

const (
	ACKStatusSuccess ACKStatus = 0x00
	ACKStatusFailed  ACKStatus = 0x02
)

type ACK struct {
	Status ACKStatus
	CmdID  uint16
}
type GetCurrentPlayingTrackChapterInfo struct {
}
type ReturnCurrentPlayingTrackChapterInfo struct {
	CurrentChapterIndex int32
	ChapterCount        int32
}
type SetCurrentPlayingTrackChapter struct {
	ChapterIndex int32
}
type GetCurrentPlayingTrackChapterPlayStatus struct {
	CurrentChapterIndex int32
}
type ReturnCurrentPlayingTrackChapterPlayStatus struct {
	ChapterLength   uint32
	ChapterPosition uint32
}
type GetCurrentPlayingTrackChapterName struct {
	ChapterIndex int32
}
type ReturnCurrentPlayingTrackChapterName struct {
	ChapterName []byte
}
type GetAudiobookSpeed struct {
}

type ReturnAudiobookSpeed struct {
	Speed byte //add enums
}
type SetAudiobookSpeed struct {
	Speed byte //add enums
}

type TrackInfoType byte

const (
	TrackInfoCaps         TrackInfoType = 0x00
	TrackInfoPodcastName  TrackInfoType = 0x01
	TrackInfoReleaseDate  TrackInfoType = 0x02
	TrackInfoDescription  TrackInfoType = 0x03
	TrackInfoLyrics       TrackInfoType = 0x04
	TrackInfoGenre        TrackInfoType = 0x05
	TrackInfoComposer     TrackInfoType = 0x06
	TrackInfoArtworkCount TrackInfoType = 0x07
)

type TrackCaps struct {
	Caps         uint32
	TrackLength  uint32
	ChapterCount uint16
}

type TrackReleaseDate struct {
	Pad uint64
}

type TrackLongText struct {
	Flags       byte
	PacketIndex uint16
	Text        byte // string
}

type GetIndexedPlayingTrackInfo struct {
	InfoType     TrackInfoType
	TrackIndex   int32
	ChapterIndex int16
}
type ReturnIndexedPlayingTrackInfo struct {
	InfoType TrackInfoType
	Info     interface{}
}

func (s ReturnIndexedPlayingTrackInfo) MarshalBinary() ([]byte, error) {
	w := bytes.Buffer{}
	if err := binary.Write(&w, binary.BigEndian, s.InfoType); err != nil {
		return nil, err
	}
	if err := binary.Write(&w, binary.BigEndian, s.Info); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (s *ReturnIndexedPlayingTrackInfo) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, &s.InfoType); err != nil {
		return err
	}

	switch s.InfoType {
	case TrackInfoCaps:
		s.Info = &TrackCaps{}
	case TrackInfoDescription, TrackInfoLyrics:
		s.Info = &TrackLongText{}
	default:
		s.Info = &struct{}{}
	}
	if err := binary.Read(r, binary.BigEndian, s.Info); err != nil {
		return err
	}
	return nil
}

type GetArtworkFormats struct {
}

type ArtworkFormat struct {
	FormatID    uint16
	PixelFormat byte
	ImageWidth  uint16
	ImageHeight uint16
}
type RetArtworkFormats struct {
	Formats []ArtworkFormat
}

func (s RetArtworkFormats) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	for i := range s.Formats {
		binary.Write(&buf, binary.BigEndian, s.Formats[i])
	}
	return buf.Bytes(), nil
}

func (s *RetArtworkFormats) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	for {
		var f ArtworkFormat
		err := binary.Read(r, binary.BigEndian, &f)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		s.Formats = append(s.Formats, f)
	}
	return nil
}

type GetTrackArtworkData struct {
	TrackIndex int32
	FormatID   uint16
	Offset     uint32
}
type RetTrackArtworkData struct {
	PacketIndex uint16
	PixelFormat byte
	ImageWidth  uint16
	ImageHeight uint16

	TopLeftX     uint16
	TopLeftY     uint16
	BottomRightX uint16
	BottomRightY uint16
	RowSize      uint32
	Data         []byte
}

//ack
type ResetDBSelection struct {
}

type DBCategoryType byte

const (
	DbCategoryPlaylist       DBCategoryType = 0x01
	DbCategoryArtist         DBCategoryType = 0x02
	DbCategoryAlbum          DBCategoryType = 0x03
	DbCategoryGenre          DBCategoryType = 0x04
	DbCategoryTrack          DBCategoryType = 0x05
	DbCategoryComposer       DBCategoryType = 0x06
	DbCategoryAudiobook      DBCategoryType = 0x07
	DbCategoryPodcast        DBCategoryType = 0x08
	DbCategoryNestedPlaylist DBCategoryType = 0x09
)

type SelectDBRecord struct {
	CategoryType DBCategoryType
	RecordIndex  int32
}
type GetNumberCategorizedDBRecords struct {
	CategoryType DBCategoryType
}
type ReturnNumberCategorizedDBRecords struct {
	RecordCount int32
}
type RetrieveCategorizedDatabaseRecords struct {
	CategoryType DBCategoryType
	Offset       uint32
	Count        int32
}
type ReturnCategorizedDatabaseRecord struct {
	RecordCategoryIndex uint32
	String              [16]byte //fix length
}
type GetPlayStatus struct {
}

type PlayerState byte

const (
	PlayerStateStopped PlayerState = 0x00
	PlayerStatePlaying PlayerState = 0x01
	PlayerStatePaused  PlayerState = 0x02
	PlayerStateError   PlayerState = 0xff
)

type ReturnPlayStatus struct {
	TrackLength   uint32
	TrackPosition uint32
	State         PlayerState
}

type GetCurrentPlayingTrackIndex struct {
}
type ReturnCurrentPlayingTrackIndex struct {
	TrackIndex int32
}
type GetIndexedPlayingTrackTitle struct {
	TrackIndex int32
}
type ReturnIndexedPlayingTrackTitle struct {
	Title []byte
}

func (s ReturnIndexedPlayingTrackTitle) MarshalBinary() ([]byte, error) {
	return s.Title, nil
}

func (s *ReturnIndexedPlayingTrackTitle) UnmarshalBinary(data []byte) error {
	s.Title = make([]byte, len(data))
	copy(s.Title, data)
	return nil
}

type GetIndexedPlayingTrackArtistName struct {
	TrackIndex int32
}
type ReturnIndexedPlayingTrackArtistName struct {
	ArtistName []byte
}

func (s ReturnIndexedPlayingTrackArtistName) MarshalBinary() ([]byte, error) {
	return s.ArtistName, nil
}

func (s *ReturnIndexedPlayingTrackArtistName) UnmarshalBinary(data []byte) error {
	s.ArtistName = make([]byte, len(data))
	copy(s.ArtistName, data)
	return nil
}

type GetIndexedPlayingTrackAlbumName struct {
	TrackIndex int32
}
type ReturnIndexedPlayingTrackAlbumName struct {
	AlbumName []byte
}

func (s ReturnIndexedPlayingTrackAlbumName) MarshalBinary() ([]byte, error) {
	return s.AlbumName, nil
}

func (s *ReturnIndexedPlayingTrackAlbumName) UnmarshalBinary(data []byte) error {
	s.AlbumName = make([]byte, len(data))
	copy(s.AlbumName, data)
	return nil
}

type SetPlayStatusChangeNotification struct {
	EventMask uint32
}

// SetPlayStatusChangeNotificationShort is another possible version of SetPlayStatusChangeNotification,
// that uses a single bit instead of a bitmask to enable/disable all play-status-change notifications
type SetPlayStatusChangeNotificationShort struct {
	Enabled bool
}
type PlayStatusChangeNotification struct {
	Status byte // finish
}
type PlayCurrentSelection struct {
	SelectedTrackIndex int32
}

type PlayControlCmd byte

const (
	PlayControlToggle      PlayControlCmd = 0x01
	PlayControlStop        PlayControlCmd = 0x02
	PlayControlNextTrack   PlayControlCmd = 0x03
	PlayControlPrevTrack   PlayControlCmd = 0x04
	PlayControlStartFF     PlayControlCmd = 0x05
	PlayControlStartRew    PlayControlCmd = 0x06
	PlayControlEndFFRew    PlayControlCmd = 0x07
	PlayControlNext        PlayControlCmd = 0x08
	PlayControlPrev        PlayControlCmd = 0x09
	PlayControlPlay        PlayControlCmd = 0x0a
	PlayControlPause       PlayControlCmd = 0x0b
	PlayControlNextChapter PlayControlCmd = 0x0c
	PlayControlPrevChapter PlayControlCmd = 0x0d
)

type PlayControl struct {
	Cmd PlayControlCmd
}
type GetTrackArtworkTimes struct {
	TrackIndex   int32
	FormatID     uint16
	ArtworkIndex uint16
	ArtworkCount int16
}
type RetTrackArtworkTimes struct {
	// empty for now
}

type ShuffleMode byte

const (
	ShuffleOff    ShuffleMode = 0x00
	ShuffleTracks ShuffleMode = 0x01
	ShuffleAlbums ShuffleMode = 0x02
)

type GetShuffle struct {
}
type ReturnShuffle struct {
	Mode ShuffleMode
}
type SetShuffle struct {
	Mode ShuffleMode
	//restore on exit
}

type RepeatMode byte

const (
	RepeatOff RepeatMode = 0x00
	RepeatOne RepeatMode = 0x01
	RepeatAll RepeatMode = 0x02
)

type GetRepeat struct {
}
type ReturnRepeat struct {
	Mode RepeatMode
}
type SetRepeat struct {
	Mode RepeatMode
	//restore on exit
}
type SetDisplayImage struct {
	//todo
}
type GetMonoDisplayImageLimits struct {
}
type ReturnMonoDisplayImageLimits struct {
	MaxWidth    uint16
	MaxHeight   uint16
	PixelFormat byte
}
type GetNumPlayingTracks struct {
}
type ReturnNumPlayingTracks struct {
	NumTracks uint32
}
type SetCurrentPlayingTrack struct {
	TrackIndex int32
}
type SelectSortDBRecord struct {
	CategoryType DBCategoryType
	RecordIndex  int32
	SortType     byte // add enum
}
type GetColorDisplayImageLimits struct {
}
type ReturnColorDisplayImageLimits struct {
	MaxWidth    uint16
	MaxHeight   uint16
	PixelFormat byte
}

type ResetDBSelectionHierarchy struct {
	Selection byte
}
type GetDBiTunesInfo struct {
}
type RetDBiTunesInfo struct {
}
type GetUIDTrackInfo struct {
}
type RetUIDTrackInfo struct {
}

type GetDBTrackInfo struct {
}
type RetDBTrackInfo struct {
}
type GetPBTrackInfo struct {
}
type RetPBTrackInfo struct {
}
