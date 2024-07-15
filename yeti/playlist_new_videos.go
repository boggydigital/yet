package yeti

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/yet/data"
)

func PlaylistNotEndedVideos(playlistId string, rdx kevlar.ReadableRedux) []string {

	playlistVideos, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId)
	if !ok {
		return nil
	}

	policy := data.DefaultDownloadPolicy
	if dp, ok := rdx.GetLastVal(data.PlaylistDownloadPolicyProperty, playlistId); ok {
		policy = data.ParseDownloadPolicy(dp)
	}

	limitVideos := data.RecentDownloadsLimit
	if policy == data.All || limitVideos > len(playlistVideos) {
		limitVideos = len(playlistVideos)
	}

	videoIds := make([]string, 0, limitVideos)

	for ii := 0; ii < limitVideos; ii++ {

		videoId := playlistVideos[ii]

		if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
			continue
		}
		videoIds = append(videoIds, videoId)
	}

	return videoIds
}
