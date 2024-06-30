package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func RemovePlaylistHandler(u *url.URL) error {
	q := u.Query()

	playlistId := q.Get("playlist-id")
	options := &PlaylistOptions{
		AutoRefresh:        q.Has("auto-refresh"),
		AutoDownload:       q.Has("auto-download"),
		PreferSingleFormat: q.Has("prefer-single-format"),
		Expand:             q.Has("expand"),
		Force:              q.Has("force"),
	}

	return RemovePlaylist(nil, playlistId, options)
}

func RemovePlaylist(rdx kvas.WriteableRedux, playlistId string, opt *PlaylistOptions) error {

	rpa := nod.Begin("removing playlist %s...", playlistId)
	defer rpa.End()

	if opt == nil {
		opt = DefaultPlaylistOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.PlaylistProperties()...)
	if err != nil {
		return rpa.EndWithError(err)
	}

	playlistId, err = yeti.ParsePlaylistId(playlistId)
	if err != nil {
		return rpa.EndWithError(err)
	}

	propertyKeys := make(map[string]string)

	if opt.AutoRefresh {
		propertyKeys[data.PlaylistAutoRefreshProperty] = playlistId
	}
	if opt.AutoDownload {
		propertyKeys[data.PlaylistAutoDownloadProperty] = playlistId
	}
	if opt.DownloadPolicy != data.Unset {
		propertyKeys[data.PlaylistDownloadPolicyProperty] = playlistId
	}
	if opt.Expand {
		propertyKeys[data.PlaylistExpandProperty] = playlistId
	}
	if opt.PreferSingleFormat {
		propertyKeys[data.PlaylistPreferSingleFormatProperty] = playlistId
	}

	for property, key := range propertyKeys {
		if err := rdx.CutKeys(property, key); err != nil {
			return rpa.EndWithError(err)
		}
	}

	rpa.EndWithResult("done")

	return nil
}
