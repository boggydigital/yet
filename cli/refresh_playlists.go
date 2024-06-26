package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RefreshPlaylistsMetadataHandler(u *url.URL) error {
	return RefreshPlaylistsMetadata(nil)
}

func RefreshPlaylistsMetadata(rdx kvas.WriteableRedux) error {

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

		if err := GetPlaylistMetadata(rdx, refreshOptions, playlistId); err != nil {
			return upma.EndWithError(err)
		}

		upma.Increment()
	}

	upma.EndWithResult("done")

	return nil
}
