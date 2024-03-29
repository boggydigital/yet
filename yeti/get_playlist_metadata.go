package yeti

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yt_urls"
	"net/http"
)

func GetPlaylistPageMetadata(playlistPage *yt_urls.PlaylistInitialData, playlistId string, allVideos bool, rdx kvas.WriteableRedux) error {

	gppma := nod.Begin(" metadata for %s", playlistId)
	defer gppma.End()

	var err error
	if playlistPage == nil {
		playlistPage, err = yt_urls.GetPlaylistPage(http.DefaultClient, playlistId)
		if err != nil {
			return gppma.EndWithError(err)
		}
	}

	phr := playlistPage.PlaylistHeaderRenderer()
	if err := rdx.ReplaceValues(data.PlaylistTitleProperty, playlistId, phr.Title.SimpleText); err != nil {
		return gppma.EndWithError(err)
	}

	if err := rdx.ReplaceValues(data.PlaylistChannelProperty, playlistId, playlistPage.PlaylistOwner()); err != nil {
		return gppma.EndWithError(err)
	}

	playlistVideos := make([]string, 0)
	videoTitles := make(map[string][]string)
	videoChannels := make(map[string][]string)
	videoLengths := make(map[string][]string)

	for playlistPage != nil &&
		len(playlistPage.Videos()) > 0 {

		for _, video := range playlistPage.Videos() {
			videoId := video.VideoId
			playlistVideos = append(playlistVideos, videoId)
			videoTitles[videoId] = []string{video.Title}
			videoChannels[videoId] = []string{video.Channel}
			videoLengths[videoId] = []string{video.Length}
		}

		if allVideos && playlistPage.HasContinuation() {
			if err = playlistPage.Continue(http.DefaultClient); err != nil {
				return gppma.EndWithError(err)
			}
		} else {
			playlistPage = nil
		}
	}

	if err := rdx.ReplaceValues(data.PlaylistVideosProperty, playlistId, playlistVideos...); err != nil {
		return gppma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.VideoTitleProperty, videoTitles); err != nil {
		return gppma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.VideoDurationProperty, videoLengths); err != nil {
		return gppma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.VideoOwnerChannelNameProperty, videoChannels); err != nil {
		return gppma.EndWithError(err)
	}

	gppma.EndWithResult("done")

	return nil
}
