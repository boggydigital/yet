package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func QueuePlaylistsDownloadsHandler(u *url.URL) error {
	return QueuePlaylistsDownloads(nil)
}

func QueuePlaylistsDownloads(rdx kevlar.WriteableRedux) error {

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

	qpda.EndWithResult("done")

	return nil
}

// queuePlaylistDownloads goes through playlist videos according to the download policy,
// skips ended and previously queued videos and queues the rest
func queuePlaylistDownloads(rdx kevlar.WriteableRedux, playlistId string) error {

	queue := make(map[string][]string)

	for _, videoId := range yeti.PlaylistNotEndedVideos(playlistId, rdx) {
		if rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			continue
		}
		queue[videoId] = []string{yeti.FmtNow()}
	}

	return rdx.BatchAddValues(data.VideoDownloadQueuedProperty, queue)
}
