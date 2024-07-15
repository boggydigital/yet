package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RefreshPlaylistsMetadataHandler(_ *url.URL) error {
	return RefreshPlaylistsMetadata(nil)
}

func RefreshPlaylistsMetadata(rdx kevlar.WriteableRedux) error {

	upma := nod.NewProgress("updating all playlists metadata...")
	defer upma.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.PlaylistProperties()...)
	if err != nil {
		return upma.EndWithError(err)
	}

	// update auto-refresh playlists metadata
	playlistIds := rdx.Keys(data.PlaylistAutoRefreshProperty)
	upma.TotalInt(len(playlistIds))

	refreshOptions := &PlaylistOptions{
		Force: true,
	}

	for _, playlistId := range playlistIds {

		if err := GetPlaylistsMetadata(rdx, refreshOptions, playlistId); err != nil {
			return upma.EndWithError(err)
		}

		upma.Increment()
	}

	upma.EndWithResult("done")

	return nil
}
