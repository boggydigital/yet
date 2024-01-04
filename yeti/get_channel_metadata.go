package yeti

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yt_urls"
	"net/http"
)

const (
	videosPerTab = 12
)

func GetChannelPageMetadata(channelPage *yt_urls.ChannelInitialData, channelId string, rdx kvas.WriteableRedux) error {

	gchpma := nod.NewProgress(" metadata for %s", channelId)
	defer gchpma.End()

	var err error
	if channelPage == nil {
		channelPage, err = yt_urls.GetChannelPage(http.DefaultClient, channelId)
		if err != nil {
			return gchpma.EndWithError(err)
		}
	}

	chmd := channelPage.ChannelMetadataRenderer()

	if err := rdx.ReplaceValues(data.ChannelTitleProperty, channelId, chmd.Title); err != nil {
		return gchpma.EndWithError(err)
	}

	if err := rdx.ReplaceValues(data.ChannelDescriptionProperty, channelId, chmd.Description); err != nil {
		return gchpma.EndWithError(err)
	}

	tabs := channelPage.Tabs()

	gchpma.TotalInt(len(tabs))

	channelPlaylists := make([]string, 0, len(tabs))
	playlistsTitles := make(map[string][]string, len(tabs))
	playlistsVideos := make(map[string][]string, videosPerTab)
	videosTitles := make(map[string][]string, len(tabs)*videosPerTab)
	videosPublishedTimes := make(map[string][]string, len(tabs)*videosPerTab)
	videosViewCounts := make(map[string][]string, len(tabs)*videosPerTab)
	videosOwnerChannel := make(map[string][]string, len(tabs)*videosPerTab)

	for _, tab := range tabs {
		for _, section := range tab.Sections() {

			playlistId := section.PlaylistId()
			channelPlaylists = append(channelPlaylists, playlistId)
			playlistsTitles[playlistId] = []string{section.Title.String()}

			for _, gvr := range section.GridVideoRenderers() {
				playlistsVideos[playlistId] = append(playlistsVideos[playlistId], gvr.VideoId)
				videosTitles[gvr.VideoId] = []string{gvr.Title.SimpleText}
				videosPublishedTimes[gvr.VideoId] = []string{gvr.PublishedTimeText.SimpleText}
				videosViewCounts[gvr.VideoId] = []string{gvr.ViewCountText.SimpleText}
				videosOwnerChannel[gvr.VideoId] = []string{chmd.Title}
			}

		}

		gchpma.Increment()
	}

	if err := rdx.AddValues(data.ChannelPlaylistsProperty, channelId, channelPlaylists...); err != nil {
		return gchpma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.PlaylistTitleProperty, playlistsTitles); err != nil {
		return gchpma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.PlaylistVideosProperty, playlistsVideos); err != nil {
		return gchpma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.VideoTitleProperty, videosTitles); err != nil {
		return gchpma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.VideoPublishTimeTextProperty, videosPublishedTimes); err != nil {
		return gchpma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.VideoViewCountProperty, videosViewCounts); err != nil {
		return gchpma.EndWithError(err)
	}

	if err := rdx.BatchAddValues(data.VideoOwnerChannelNameProperty, videosOwnerChannel); err != nil {
		return gchpma.EndWithError(err)
	}

	gchpma.EndWithResult("done")

	return nil
}
