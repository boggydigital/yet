package cli

import (
	"net/url"

	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func QueuePlaylistsDownloadsHandler(_ *url.URL) error {
	return QueuePlaylistsDownloads(nil)
}

func QueuePlaylistsDownloads(rdx redux.Writeable) error {

	qpda := nod.NewProgress("queueing playlists downloads...")
	defer qpda.Done()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return err
	}

	qpda.TotalInt(rdx.Len(data.PlaylistAutoDownloadProperty))

	for playlistId := range rdx.Keys(data.PlaylistAutoDownloadProperty) {

		if err = yeti.QueuePlaylistDownloads(rdx, playlistId); err != nil {
			return err
		}

		qpda.Increment()
	}

	return nil
}
