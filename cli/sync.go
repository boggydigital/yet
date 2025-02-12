package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func SyncHandler(u *url.URL) error {
	q := u.Query()

	options := &VideoOptions{
		Force: q.Has("Force"),
	}
	return Sync(nil, options)
}

func Sync(rdx redux.Writeable, opt *VideoOptions) error {

	sa := nod.Begin("syncing yet data...")
	defer sa.Done()

	if opt == nil {
		opt = DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return err
	}

	if err := UpdateYtDlp(false); err != nil {
		return err
	}

	if err := RefreshChannelsMetadata(rdx); err != nil {
		return err
	}
	if err := QueueChannelsDownloads(rdx); err != nil {
		return err
	}

	if err := RefreshPlaylistsMetadata(rdx); err != nil {
		return err
	}
	if err := QueuePlaylistsDownloads(rdx); err != nil {
		return err
	}

	if err := ProcessQueue(rdx, opt); err != nil {
		return err
	}

	if err := DehydratePosters(false); err != nil {
		return err
	}

	if err := ScrubEndedProperties(rdx); err != nil {
		return err
	}
	if err := ScrubDepositionProperties(rdx); err != nil {
		return err
	}
	if err := CleanupEndedVideos(false, rdx); err != nil {
		return err
	}

	if err := Backup(); err != nil {
		return err
	}

	return nil
}
