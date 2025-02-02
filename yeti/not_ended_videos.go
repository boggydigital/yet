package yeti

import (
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

func ChannelNotEndedVideos(channelId string, rdx redux.Readable) []string {
	return notEndedVideos(
		channelId,
		data.ChannelVideosProperty,
		data.ChannelDownloadPolicyProperty,
		rdx)
}

func PlaylistNotEndedVideos(playlistId string, rdx redux.Readable) []string {
	return notEndedVideos(
		playlistId,
		data.PlaylistVideosProperty,
		data.PlaylistDownloadPolicyProperty,
		rdx)
}

func notEndedVideos(id string, videosProperty, downloadPolicyProperty string, rdx redux.Readable) []string {

	videos, ok := rdx.GetAllValues(videosProperty, id)
	if !ok {
		return nil
	}

	policy := data.DefaultDownloadPolicy
	if dp, ok := rdx.GetLastVal(downloadPolicyProperty, id); ok {
		policy = data.ParseDownloadPolicy(dp)
	}

	limitVideos := data.RecentDownloadsLimit
	if policy == data.All || limitVideos > len(videos) {
		limitVideos = len(videos)
	}

	videoIds := make([]string, 0, limitVideos)

	for ii := 0; ii < limitVideos; ii++ {

		videoId := videos[ii]

		if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
			continue
		}
		videoIds = append(videoIds, videoId)
	}

	return videoIds
}
