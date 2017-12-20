package extremote

import (
	"github.com/oandrew/ipod"
)

type DeviceExtRemote interface {
	PlaybackStatus() (trackLength, trackPos uint32, state PlayerState)
}

func ackSuccess(req *ipod.Command) ACK {
	return ACK{Status: ACKStatusSuccess, CmdID: req.ID.CmdID()}
}

// func ackPending(req ipod.Packet, maxWait uint32) ACKPending {
// 	return ACKPending{Status: ACKStatusPending, CmdID: uint8(req.ID.CmdID()), MaxWait: maxWait}
// }

func HandleExtRemote(req *ipod.Command, tr ipod.CommandWriter, dev DeviceExtRemote) error {
	//log.Printf("Req: %#v", req)
	switch msg := req.Payload.(type) {

	case GetCurrentPlayingTrackChapterInfo:
		ipod.Respond(req, tr, ReturnCurrentPlayingTrackChapterInfo{
			CurrentChapterIndex: -1,
			ChapterCount:        0,
		})
	case SetCurrentPlayingTrackChapter:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetCurrentPlayingTrackChapterPlayStatus:
		ipod.Respond(req, tr, ReturnCurrentPlayingTrackChapterPlayStatus{
			ChapterPosition: 0,
			ChapterLength:   0,
		})
	case GetCurrentPlayingTrackChapterName:
		ipod.Respond(req, tr, ReturnCurrentPlayingTrackChapterName{
			ChapterName: ipod.StringToBytes("chapter"),
		})
	case GetAudiobookSpeed:
		ipod.Respond(req, tr, ReturnAudiobookSpeed{
			Speed: 0,
		})
	case SetAudiobookSpeed:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetIndexedPlayingTrackInfo:
		var info interface{}
		switch msg.InfoType {
		case TrackInfoCaps:
			info = &TrackCaps{
				Caps:         0x0,
				TrackLength:  300 * 1000,
				ChapterCount: 0,
			}
		case TrackInfoDescription, TrackInfoLyrics:
			info = &TrackLongText{
				Flags:       0x0,
				PacketIndex: 0,
				Text:        0x00,
			}
		case TrackInfoArtworkCount:
			info = struct{}{}
		default:
			info = []byte{0x00}

		}
		ipod.Respond(req, tr, ReturnIndexedPlayingTrackInfo{
			InfoType: msg.InfoType,
			Info:     info,
		})
	case GetArtworkFormats:
		ipod.Respond(req, tr, RetArtworkFormats{})
	case GetTrackArtworkData:
		ipod.Respond(req, tr, ACK{
			Status: ACKStatusFailed,
			CmdID:  req.ID.CmdID(),
		})
	case ResetDBSelection:
		ipod.Respond(req, tr, ackSuccess(req))
	case SelectDBRecord:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetNumberCategorizedDBRecords:
		ipod.Respond(req, tr, ReturnNumberCategorizedDBRecords{
			RecordCount: 0,
		})
	case RetrieveCategorizedDatabaseRecords:
		ipod.Respond(req, tr, ReturnCategorizedDatabaseRecord{})
	case GetPlayStatus:
		ipod.Respond(req, tr, ReturnPlayStatus{
			TrackLength:   300 * 1000,
			TrackPosition: 20 * 1000,
			State:         PlayerStatePaused,
		})
	case GetCurrentPlayingTrackIndex:
		ipod.Respond(req, tr, ReturnCurrentPlayingTrackIndex{
			TrackIndex: 0,
		})
	case GetIndexedPlayingTrackTitle:
		ipod.Respond(req, tr, ReturnIndexedPlayingTrackTitle{
			Title: ipod.StringToBytes("title"),
		})
	case GetIndexedPlayingTrackArtistName:
		ipod.Respond(req, tr, ReturnIndexedPlayingTrackArtistName{
			ArtistName: ipod.StringToBytes("artist"),
		})
	case GetIndexedPlayingTrackAlbumName:
		ipod.Respond(req, tr, ReturnIndexedPlayingTrackAlbumName{
			AlbumName: ipod.StringToBytes("album"),
		})
	case SetPlayStatusChangeNotification:
		ipod.Respond(req, tr, ackSuccess(req))
	//case PlayStatusChangeNotification:
	case PlayCurrentSelection:
		ipod.Respond(req, tr, ackSuccess(req))
	case PlayControl:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetTrackArtworkTimes:
		ipod.Respond(req, tr, RetTrackArtworkTimes{})
	case GetShuffle:
		ipod.Respond(req, tr, ReturnShuffle{Mode: ShuffleOff})
	case SetShuffle:
		ipod.Respond(req, tr, ackSuccess(req))

	case GetRepeat:
		ipod.Respond(req, tr, ReturnRepeat{Mode: RepeatOff})
	case SetRepeat:
		ipod.Respond(req, tr, ackSuccess(req))

	case SetDisplayImage:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetMonoDisplayImageLimits:
		ipod.Respond(req, tr, ReturnMonoDisplayImageLimits{
			MaxWidth:    640,
			MaxHeight:   960,
			PixelFormat: 0x01,
		})
	case GetNumPlayingTracks:
		ipod.Respond(req, tr, ReturnNumPlayingTracks{
			NumTracks: 0,
		})
	case SetCurrentPlayingTrack:
	case SelectSortDBRecord:
	case GetColorDisplayImageLimits:
		ipod.Respond(req, tr, ReturnColorDisplayImageLimits{
			MaxWidth:    640,
			MaxHeight:   960,
			PixelFormat: 0x01,
		})
	case ResetDBSelectionHierarchy:
		//noop

	case GetDBiTunesInfo:
	// RetDBiTunesInfo:
	case GetUIDTrackInfo:
	// RetUIDTrackInfo:
	case GetDBTrackInfo:
	// RetDBTrackInfo:
	case GetPBTrackInfo:
	// RetPBTrackInfo:

	default:
		_ = msg
	}
	return nil
}
