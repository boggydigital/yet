package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
)

func SyncHandler(u *url.URL) error {
	force := u.Query().Has("force")
	singleFormat := u.Query().Has("single-format")
	return Sync(force, singleFormat)
}

func Sync(force, singleFormat bool) error {

	sa := nod.Begin("syncing playlists subscriptions...")
	defer sa.End()

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		return sa.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return sa.EndWithError(err)
	}

	if err := UpdatePlaylistsMetadata(rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := QueuePlaylistsNewVideos(rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := Download(rdx, true, force, singleFormat); err != nil {
		return sa.EndWithError(err)
	}

	if err := CleanupEnded(rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := Backup(); err != nil {
		return sa.EndWithError(err)
	}

	return nil
}
