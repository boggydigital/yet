package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RefreshPlaylistsMetadataHandler(_ *url.URL) error {
	return RefreshPlaylistsMetadata(nil)
}

func RefreshPlaylistsMetadata(rdx redux.Writeable) error {

	upma := nod.NewProgress("updating all playlists metadata...")
	defer upma.Done()

	var err error
	rdx, err = validateWritableRedux(rdx, data.PlaylistProperties()...)
	if err != nil {
		return err
	}

	// update auto-refresh playlists metadata
	upma.TotalInt(rdx.Len(data.PlaylistAutoRefreshProperty))

	refreshOptions := &PlaylistOptions{
		Force: true,
	}

	for playlistId := range rdx.Keys(data.PlaylistAutoRefreshProperty) {

		if err := GetPlaylistsMetadata(rdx, refreshOptions, playlistId); err != nil {
			return err
		}

		upma.Increment()
	}

	return nil
}
