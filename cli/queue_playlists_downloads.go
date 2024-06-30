package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func QueuePlaylistsDownloadsHandler(u *url.URL) error {
	return QueuePlaylistsDownloads(nil)
}

func QueuePlaylistsDownloads(rdx kvas.WriteableRedux) error {

	qpda := nod.NewProgress("queueing playlists downloads...")
	defer qpda.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return qpda.EndWithError(err)
	}

	playlistIds := rdx.Keys(data.PlaylistAutoDownloadProperty)
	qpda.TotalInt(len(playlistIds))

	for _, playlistId := range playlistIds {

		if err := queuePlaylistDownloads(rdx, playlistId); err != nil {
			return qpda.EndWithError(err)
		}

		qpda.Increment()
	}

	//notNewIndicatorProperties := []string{
	//	data.VideosDownloadQueueProperty,
	//	data.VideosWatchlistProperty,
	//	data.VideoEndedProperty,
	//	data.VideoProgressProperty}
	//
	//for _, pdq := range rdx.Keys(data.PlaylistDownloadQueueProperty) {
	//	if newVideos, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, pdq); ok && len(newVideos) > 0 {
	//		for _, videoId := range newVideos {
	//
	//			skipVideo := false
	//			// don't add videos already in download queue, watch, Ended, in progress
	//			for _, nnip := range notNewIndicatorProperties {
	//				if rdx.HasKey(nnip, videoId) {
	//					skipVideo = true
	//					break
	//				}
	//			}
	//
	//			if skipVideo {
	//				break
	//			}
	//
	//			if err := rdx.AddValues(data.VideosDownloadQueueProperty, videoId, data.TrueValue); err != nil {
	//				return qpda.EndWithError(err)
	//			}
	//
	//			// set video to download as single format if playlist has that flag set
	//			if rdx.HasKey(data.PlaylistSingleFormatDownloadProperty, pdq) {
	//				if err := rdx.AddValues(data.VideoSingleFormatDownloadProperty, videoId, data.TrueValue); err != nil {
	//					return qpda.EndWithError(err)
	//				}
	//			}
	//		}
	//	}
	//}

	qpda.EndWithResult("done")

	return nil
}

func queuePlaylistDownloads(rdx kvas.WriteableRedux, playlistId string) error {
	return nil
}
