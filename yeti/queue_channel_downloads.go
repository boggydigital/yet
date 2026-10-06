package yeti

import (
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

// QueueChannelDownloads goes through channel videos according to the download policy,
// skips ended and previously queued videos and queues the rest
func QueueChannelDownloads(rdx redux.Writeable, channelId string) error {

	queue := make(map[string][]string)

	for _, videoId := range ChannelNotEndedVideos(channelId, data.RecentDownloadsLimit, rdx) {
		if rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			continue
		}
		queue[videoId] = []string{FmtNow()}
	}

	return rdx.BatchAddValues(data.VideoDownloadQueuedProperty, queue)
}
