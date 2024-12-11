package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func RemovePlaylistHandler(u *url.URL) error {
	q := u.Query()

	playlistId := q.Get("playlist-id")
	options := &PlaylistOptions{
		AutoRefresh:  q.Has("auto-refresh"),
		AutoDownload: q.Has("auto-download"),
		Expand:       q.Has("expand"),
		Force:        q.Has("force"),
	}

	return RemovePlaylist(nil, playlistId, options)
}

func RemovePlaylist(rdx kevlar.WriteableRedux, playlistId string, opt *PlaylistOptions) error {

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
	if opt.DownloadPolicy != data.DefaultDownloadPolicy {
		propertyKeys[data.PlaylistDownloadPolicyProperty] = playlistId
	}
	if opt.Expand {
		propertyKeys[data.PlaylistExpandProperty] = playlistId
	}

	for property, key := range propertyKeys {
		if err := rdx.CutKeys(property, key); err != nil {
			return rpa.EndWithError(err)
		}
	}

	rpa.EndWithResult("done")

	return nil
}
