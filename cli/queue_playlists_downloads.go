package cli

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
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

	qpda.EndWithResult("done")

	return nil
}

// queuePlaylistDownloads goes through playlist videos according to the download policy,
// skips ended and previously queued videos and queues the rest
func queuePlaylistDownloads(rdx kvas.WriteableRedux, playlistId string) error {

	playlistVideos, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId)
	if !ok {
		return fmt.Errorf("cannot queue downloads for an empty playlist %s", playlistId)
	}

	policy := data.Unset
	if dp, ok := rdx.GetLastVal(data.PlaylistDownloadPolicyProperty, playlistId); ok {
		policy = data.ParsePlaylistDownloadPolicy(dp)
	}

	limitVideos := data.RecentDownloadsLimit
	if policy == data.All || limitVideos > len(playlistVideos) {
		limitVideos = len(playlistVideos)
	}

	queue := make(map[string][]string)

	for ii := 0; ii < limitVideos; ii++ {

		videoId := playlistVideos[ii]

		if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
			continue
		}
		if rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			continue
		}
		queue[videoId] = []string{yeti.FmtNow()}
	}

	return rdx.BatchAddValues(data.VideoDownloadQueuedProperty, queue)
}
