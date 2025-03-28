package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func QueuePlaylistsDownloadsHandler(_ *url.URL) error {
	return QueuePlaylistsDownloads(nil)
}

func QueuePlaylistsDownloads(rdx redux.Writeable) error {

	qpda := nod.NewProgress("queueing playlists downloads...")
	defer qpda.Done()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return err
	}

	qpda.TotalInt(rdx.Len(data.PlaylistAutoDownloadProperty))

	for playlistId := range rdx.Keys(data.PlaylistAutoDownloadProperty) {

		if err := queuePlaylistDownloads(rdx, playlistId); err != nil {
			return err
		}

		qpda.Increment()
	}

	return nil
}

// queuePlaylistDownloads goes through playlist videos according to the download policy,
// skips ended and previously queued videos and queues the rest
func queuePlaylistDownloads(rdx redux.Writeable, playlistId string) error {

	queue := make(map[string][]string)

	for _, videoId := range yeti.PlaylistNotEndedVideos(playlistId, rdx) {
		if rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			continue
		}
		queue[videoId] = []string{yeti.FmtNow()}
	}

	return rdx.BatchAddValues(data.VideoDownloadQueuedProperty, queue)
}
