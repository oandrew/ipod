package extremote

import (
	"git.andrewo.pw/andrew/ipod"
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

type ACK struct {
}
type GetCurrentPlayingTrackChapterInfo struct {
}
type ReturnCurrentPlayingTrackChapterInfo struct {
}
type SetCurrentPlayingTrackChapter struct {
}
type GetCurrentPlayingTrackChapterPlayStatus struct {
}
type ReturnCurrentPlayingTrackChapterPlayStatus struct {
}
type GetCurrentPlayingTrackChapterName struct {
}
type ReturnCurrentPlayingTrackChapterName struct {
}
type GetAudiobookSpeed struct {
}
type ReturnAudiobookSpeed struct {
}
type SetAudiobookSpeed struct {
}
type GetIndexedPlayingTrackInfo struct {
}
type ReturnIndexedPlayingTrackInfo struct {
}
type GetArtworkFormats struct {
}
type RetArtworkFormats struct {
}
type GetTrackArtworkData struct {
}
type RetTrackArtworkData struct {
}
type ResetDBSelection struct {
}
type SelectDBRecord struct {
}
type GetNumberCategorizedDBRecords struct {
}
type ReturnNumberCategorizedDBRecords struct {
}
type RetrieveCategorizedDatabaseRecords struct {
}
type ReturnCategorizedDatabaseRecord struct {
}
type GetPlayStatus struct {
}
type ReturnPlayStatus struct {
}
type GetCurrentPlayingTrackIndex struct {
}
type ReturnCurrentPlayingTrackIndex struct {
}
type GetIndexedPlayingTrackTitle struct {
}
type ReturnIndexedPlayingTrackTitle struct {
}
type GetIndexedPlayingTrackArtistName struct {
}
type ReturnIndexedPlayingTrackArtistName struct {
}
type GetIndexedPlayingTrackAlbumName struct {
}
type ReturnIndexedPlayingTrackAlbumName struct {
}
type SetPlayStatusChangeNotification struct {
}
type PlayStatusChangeNotification struct {
}
type PlayCurrentSelection struct {
}
type PlayControl struct {
}
type GetTrackArtworkTimes struct {
}
type RetTrackArtworkTimes struct {
}
type GetShuffle struct {
}
type ReturnShuffle struct {
}
type SetShuffle struct {
}
type GetRepeat struct {
}
type ReturnRepeat struct {
}
type SetRepeat struct {
}
type SetDisplayImage struct {
}
type GetMonoDisplayImageLimits struct {
}
type ReturnMonoDisplayImageLimits struct {
}
type GetNumPlayingTracks struct {
}
type ReturnNumPlayingTracks struct {
}
type SetCurrentPlayingTrack struct {
}
type SelectSortDBRecord struct {
}
type GetColorDisplayImageLimits struct {
}
type ReturnColorDisplayImageLimits struct {
}
type ResetDBSelectionHierarchy struct {
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
