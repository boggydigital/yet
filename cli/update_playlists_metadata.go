package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
)

func UpdatePlaylistsMetadataHandler(u *url.URL) error {
	return UpdatePlaylistsMetadata()
}

func UpdatePlaylistsMetadata() error {

	upma := nod.NewProgress("updating all playlists metadata...")
	defer upma.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return upma.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.PlaylistWatchlistProperty)
	if err != nil {
		return upma.EndWithError(err)
	}

	if err := GetPlaylistMetadata(false, true, rdx.Keys(data.PlaylistWatchlistProperty)...); err != nil {
		return upma.EndWithError(err)
	}

	return nil
}
