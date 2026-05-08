package yeti

import (
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

func ChannelNotEndedVideos(channelId string, limitVideos int, rdx redux.Readable) []string {
	return notEndedVideos(
		channelId,
		data.ChannelVideosProperty,
		data.ChannelDownloadPolicyProperty,
		limitVideos,
		rdx)
}

func PlaylistNotEndedVideos(playlistId string, limitVideos int, rdx redux.Readable) []string {
	return notEndedVideos(
		playlistId,
		data.PlaylistVideosProperty,
		data.PlaylistDownloadPolicyProperty,
		limitVideos,
		rdx)
}

func notEndedVideos(id string, videosProperty, downloadPolicyProperty string, limitVideos int, rdx redux.Readable) []string {

	videos, ok := rdx.GetAllValues(videosProperty, id)
	if !ok {
		return nil
	}

	policy := data.DefaultDownloadPolicy
	if dp, sure := rdx.GetLastVal(downloadPolicyProperty, id); sure {
		policy = data.ParseDownloadPolicy(dp)
	}

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
