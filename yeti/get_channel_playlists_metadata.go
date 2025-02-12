package yeti

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
)

func GetChannelPlaylistsMetadata(channelPlaylistsPage *youtube_urls.ChannelPlaylistsInitialData, channelId string, rdx redux.Writeable) error {

	gcpma := nod.Begin(" channel playlists metadata for %s", channelId)
	defer gcpma.Done()

	if err := rdx.MustHave(
		data.ChannelTitleProperty,
		data.ChannelDescriptionProperty,
		data.ChannelPlaylistsProperty,
		data.PlaylistTitleProperty,
		data.PlaylistChannelProperty); err != nil {
		return err
	}

	var err error
	if channelPlaylistsPage == nil {
		channelPlaylistsPage, err = youtube_urls.GetChannelPlaylistsPage(http.DefaultClient, channelId)
		if err != nil {
			return err
		}
	}

	channelTitle := ""
	if channelTitle = channelPlaylistsPage.Metadata.ChannelMetadataRenderer.Title; channelTitle != "" {
		if err := rdx.ReplaceValues(data.ChannelTitleProperty, channelId, channelTitle); err != nil {
			return err
		}
	}

	if description := channelPlaylistsPage.Metadata.ChannelMetadataRenderer.Description; description != "" {
		if err = rdx.ReplaceValues(data.ChannelDescriptionProperty, channelId, description); err != nil {
			return err
		}
	}

	channelPlaylists := make([]string, 0)
	playlistsTitles := make(map[string][]string)
	playlistsChannels := make(map[string][]string)

	for _, playlist := range channelPlaylistsPage.Playlists() {
		playlistId := playlist.PlaylistId
		channelPlaylists = append(channelPlaylists, playlistId)
		playlistsTitles[playlistId] = []string{playlist.Title.String()}
		playlistsChannels[playlistId] = []string{channelTitle}
	}

	if err = rdx.ReplaceValues(data.ChannelPlaylistsProperty, channelId, channelPlaylists...); err != nil {
		return err
	}

	if err = rdx.BatchReplaceValues(data.PlaylistTitleProperty, playlistsTitles); err != nil {
		return err
	}

	if err = rdx.BatchReplaceValues(data.PlaylistChannelProperty, playlistsChannels); err != nil {
		return err
	}

	return nil
}
