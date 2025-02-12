package yeti

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
)

func GetPlaylistMetadata(playlistPage *youtube_urls.PlaylistInitialData, playlistId string, expand bool, rdx redux.Writeable) error {

	gppma := nod.Begin(" metadata for %s", playlistId)
	defer gppma.Done()

	if err := rdx.MustHave(
		data.PlaylistTitleProperty,
		data.PlaylistChannelProperty,
		data.PlaylistVideosProperty,
		data.VideoTitleProperty,
		data.VideoDurationProperty,
		data.VideoOwnerChannelNameProperty); err != nil {
		return err
	}

	var err error
	if playlistPage == nil {
		playlistPage, err = youtube_urls.GetPlaylistPage(http.DefaultClient, playlistId)
		if err != nil {
			return err
		}
	}

	phr := playlistPage.PlaylistHeaderRenderer()
	if phr.Title.SimpleText != "" {
		if err = rdx.AddValues(data.PlaylistTitleProperty, playlistId, phr.Title.SimpleText); err != nil {
			return err
		}
	}

	if playlistPage.PlaylistOwner() != "" {
		if err = rdx.AddValues(data.PlaylistChannelProperty, playlistId, playlistPage.PlaylistOwner()); err != nil {
			return err
		}
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
			videoLengths[videoId] = []string{video.LengthSeconds}
		}

		if expand && playlistPage.HasContinuation() {
			if err = playlistPage.Continue(http.DefaultClient); err != nil {
				return err
			}
		} else {
			playlistPage = nil
		}
	}

	if err = rdx.ReplaceValues(data.PlaylistVideosProperty, playlistId, playlistVideos...); err != nil {
		return err
	}

	if err = rdx.BatchAddValues(data.VideoTitleProperty, videoTitles); err != nil {
		return err
	}

	if err = rdx.BatchAddValues(data.VideoDurationProperty, videoLengths); err != nil {
		return err
	}

	if err = rdx.BatchAddValues(data.VideoOwnerChannelNameProperty, videoChannels); err != nil {
		return err
	}

	return nil
}
