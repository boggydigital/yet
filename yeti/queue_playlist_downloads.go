package yeti

import (
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

// QueuePlaylistDownloads goes through playlist videos according to the download policy,
// skips ended and previously queued videos and queues the rest
func QueuePlaylistDownloads(rdx redux.Writeable, playlistId string) error {

	queue := make(map[string][]string)

	for _, videoId := range PlaylistNotEndedVideos(playlistId, data.RecentDownloadsLimit, rdx) {
		if rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			continue
		}
		queue[videoId] = []string{FmtNow()}
	}

	return rdx.BatchAddValues(data.VideoDownloadQueuedProperty, queue)
}
