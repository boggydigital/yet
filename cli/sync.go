package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func SyncHandler(u *url.URL) error {
	q := u.Query()

	options := &VideoOptions{
		PreferSingleFormat: q.Has("prefer-single-format"),
		Force:              q.Has("Force"),
	}
	return Sync(nil, options)
}

func Sync(rdx kevlar.WriteableRedux, opt *VideoOptions) error {

	sa := nod.Begin("syncing playlists...")
	defer sa.End()

	if opt == nil {
		opt = DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return sa.EndWithError(err)
	}

	if err := RefreshPlaylistsMetadata(rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := QueuePlaylistsDownloads(rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := DownloadQueue(rdx, opt); err != nil {
		return sa.EndWithError(err)
	}

	if err := CleanupEnded(false, rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := Backup(); err != nil {
		return sa.EndWithError(err)
	}

	return nil
}
