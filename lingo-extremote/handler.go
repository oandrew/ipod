package extremote

import (
	"github.com/oandrew/ipod"
)

type DeviceExtRemote interface {
	PlaybackStatus() (trackLength, trackPos uint32, state PlayerState)
}

func ackSuccess(req ipod.Packet) ACK {
	return ACK{Status: ACKStatusSuccess, CmdID: req.ID.CmdID()}
}

// func ackPending(req ipod.Packet, maxWait uint32) ACKPending {
// 	return ACKPending{Status: ACKStatusPending, CmdID: uint8(req.ID.CmdID()), MaxWait: maxWait}
// }

func HandleExtRemote(req ipod.Packet, tr ipod.PacketWriter, dev DeviceExtRemote) error {
	//log.Printf("Req: %#v", req)
	switch msg := req.Payload.(type) {

	case GetCurrentPlayingTrackChapterInfo:
	// ReturnCurrentPlayingTrackChapterInfo:
	case SetCurrentPlayingTrackChapter:
	case GetCurrentPlayingTrackChapterPlayStatus:
	// ReturnCurrentPlayingTrackChapterPlayStatus:
	case GetCurrentPlayingTrackChapterName:
	// ReturnCurrentPlayingTrackChapterName:
	case GetAudiobookSpeed:
	// ReturnAudiobookSpeed:
	case SetAudiobookSpeed:
	case GetIndexedPlayingTrackInfo:
	// ReturnIndexedPlayingTrackInfo:
	case GetArtworkFormats:
	// RetArtworkFormats:
	case GetTrackArtworkData:
	// RetTrackArtworkData:
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
		// var resp ReturnPlayStatus
		// resp.TrackLength, resp.TrackPosition, resp.State = dev.PlaybackStatus()
		// ipod.Respond(req, tr, resp)
		ipod.Respond(req, tr, ReturnPlayStatus{
			TrackLength:   300 * 1000,
			TrackPosition: 20 * 1000,
			State:         PlayerStatePaused,
		})
	case GetCurrentPlayingTrackIndex:
	// ReturnCurrentPlayingTrackIndex:
	case GetIndexedPlayingTrackTitle:
	// ReturnIndexedPlayingTrackTitle:
	case GetIndexedPlayingTrackArtistName:
	// ReturnIndexedPlayingTrackArtistName:
	case GetIndexedPlayingTrackAlbumName:
	// ReturnIndexedPlayingTrackAlbumName:
	case SetPlayStatusChangeNotification:
		ipod.Respond(req, tr, ackSuccess(req))
	//case PlayStatusChangeNotification:
	case PlayCurrentSelection:
	case PlayControl:
		ipod.Respond(req, tr, ackSuccess(req))
	case GetTrackArtworkTimes:
	// RetTrackArtworkTimes:

	case GetShuffle:
		ipod.Respond(req, tr, ReturnShuffle{Mode: ShuffleOff})
	case SetShuffle:
		ipod.Respond(req, tr, ackSuccess(req))

	case GetRepeat:
		ipod.Respond(req, tr, ReturnRepeat{Mode: RepeatOff})
	case SetRepeat:
		ipod.Respond(req, tr, ackSuccess(req))

	case SetDisplayImage:
	case GetMonoDisplayImageLimits:
	// ReturnMonoDisplayImageLimits:
	case GetNumPlayingTracks:
	// ReturnNumPlayingTracks:
	case SetCurrentPlayingTrack:
	case SelectSortDBRecord:
	case GetColorDisplayImageLimits:
	// ReturnColorDisplayImageLimits:
	case ResetDBSelectionHierarchy:
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
