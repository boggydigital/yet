package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"strings"
)

func GetPlaylistMetadataHandler(u *url.URL) error {
	q := u.Query()
	ids := strings.Split(q.Get("id"), ",")
	allVideos := q.Has("all-videos")
	force := q.Has("force")
	return GetPlaylistMetadata(allVideos, force, ids...)
}

func GetPlaylistMetadata(allVideos, force bool, ids ...string) error {
	gpma := nod.NewProgress("getting playlist metadata...")
	defer gpma.End()

	gpma.TotalInt(len(ids))

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return gpma.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.AllProperties()...)
	if err != nil {
		return gpma.EndWithError(err)
	}

	for _, playlistId := range ids {

		if rxa.HasKey(data.PlaylistTitleProperty, playlistId) && !force {
			continue
		}

		if err := getPlaylistPageMetadata(nil, playlistId, allVideos, rxa); err != nil {
			gpma.Error(err)
		}

		gpma.Increment()
	}

	gpma.EndWithResult("done")

	return nil
}

func getPlaylistPageMetadata(playlistPage *yt_urls.ContextualPlaylist, playlistId string, allVideos bool, rxa kvas.ReduxAssets) error {

	gppma := nod.Begin(" metadata for %s", playlistId)
	defer gppma.End()

	var err error
	if playlistPage == nil {
		playlistPage, err = yt_urls.GetPlaylistPage(http.DefaultClient, playlistId)
		if err != nil {
			return gppma.EndWithError(err)
		}
	}

	if err := rxa.ReplaceValues(data.PlaylistTitleProperty, playlistId, playlistPage.Playlist.Header.Title.SimpleText); err != nil {
		return gppma.EndWithError(err)
	}

	playlistVideos := make([]string, 0)
	videoTitles := make(map[string][]string)
	videoChannels := make(map[string][]string)

	for playlistPage != nil &&
		len(playlistPage.Videos()) > 0 {

		for _, video := range playlistPage.Videos() {
			videoId := video.VideoId
			playlistVideos = append(playlistVideos, videoId)
			videoTitles[videoId] = []string{video.Title}
			videoChannels[videoId] = []string{video.Channel}
		}

		if allVideos && playlistPage.HasContinuation() {
			playlistPage, err = playlistPage.Continue(http.DefaultClient)
			if err != nil {
				return gppma.EndWithError(err)
			}
		} else {
			playlistPage = nil
		}
	}

	if err := rxa.ReplaceValues(data.PlaylistVideosProperty, playlistId, playlistVideos...); err != nil {
		return gppma.EndWithError(err)
	}

	if err := rxa.BatchAddValues(data.VideoTitleProperty, videoTitles); err != nil {
		return gppma.EndWithError(err)
	}

	if err := rxa.BatchAddValues(data.VideoOwnerChannelNameProperty, videoChannels); err != nil {
		return gppma.EndWithError(err)
	}

	gppma.EndWithResult("done")

	return nil
}
