package cli

import (
	"github.com/boggydigital/nod"
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

	if err := UpdatePlaylistsMetadata(); err != nil {
		return sa.EndWithError(err)
	}

	if err := UpdatePlaylistsNewVideos(); err != nil {
		return sa.EndWithError(err)
	}

	if err := QueuePlaylistsNewVideos(); err != nil {
		return sa.EndWithError(err)
	}

	if err := Download(nil, true, force, singleFormat); err != nil {
		return sa.EndWithError(err)
	}

	if err := Backup(); err != nil {
		return sa.EndWithError(err)
	}

	return nil
}
