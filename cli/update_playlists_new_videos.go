package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
)

func UpdatePlaylistsNewVideosHandler(u *url.URL) error {
	return UpdatePlaylistsNewVideos()
}

func UpdatePlaylistsNewVideos() error {

	upmnva := nod.NewProgress("updating playlists new videos (new since last ended)...")
	defer upmnva.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return upmnva.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return upmnva.EndWithError(err)
	}

	playlistIds := rdx.Keys(data.PlaylistWatchlistProperty)

	upmnva.TotalInt(len(playlistIds))

	pnv := make(map[string][]string, len(playlistIds))

	for _, playlistId := range playlistIds {
		pnv[playlistId] = playlistNewVideos(rdx, playlistId)
		upmnva.Increment()
	}

	if err := rdx.BatchReplaceValues(data.PlaylistNewVideosProperty, pnv); err != nil {
		return upmnva.EndWithError(err)
	}

	upmnva.EndWithResult("done")

	return nil
}

func playlistNewVideos(rdx kvas.ReadableRedux, playlistId string) []string {

	newVideos := make([]string, 0)

	if playlistVideos, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok {
		for _, videoId := range playlistVideos {
			if rdx.HasKey(data.VideoEndedProperty, videoId) {
				break
			}
			newVideos = append(newVideos, videoId)
		}
	}

	return newVideos
}
