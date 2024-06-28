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
	q := u.Query()

	options := &DownloadVideoOptions{
		PreferSingleFormat: q.Has("prefer-single-format"),
		Force:              q.Has("force"),
	}
	return Sync(options)
}

func Sync(options *DownloadVideoOptions) error {

	sa := nod.Begin("syncing playlists subscriptions...")
	defer sa.End()

	if options == nil {
		options = DefaultDownloadVideoOptions()
	}

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

	if err := DownloadQueue(rdx, options); err != nil {
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
