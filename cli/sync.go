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
		// TODO: remove this if better options are available
		// (temporary?) workaround - force single format for all videos
		// to mitigate new visitorData, poToken requirements at the cost of video quality
		PreferSingleFormat: true, //q.Has("prefer-single-format"),
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

	if err := RefreshChannelsMetadata(rdx); err != nil {
		return sa.EndWithError(err)
	}
	if err := QueueChannelsDownloads(rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := RefreshPlaylistsMetadata(rdx); err != nil {
		return sa.EndWithError(err)
	}
	if err := QueuePlaylistsDownloads(rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := ProcessQueue(rdx, opt); err != nil {
		return sa.EndWithError(err)
	}

	if err := ScrubEndedProperties(rdx); err != nil {
		return sa.EndWithError(err)
	}
	if err := ScrubDepositionProperties(rdx); err != nil {
		return sa.EndWithError(err)
	}
	if err := CleanupEndedVideos(false, rdx); err != nil {
		return sa.EndWithError(err)
	}

	if err := Backup(); err != nil {
		return sa.EndWithError(err)
	}

	return nil
}
