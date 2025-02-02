package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func AddPlaylistHandler(u *url.URL) error {
	q := u.Query()

	playlistId := q.Get("playlist-id")
	options := &PlaylistOptions{
		AutoRefresh:    q.Has("auto-refresh"),
		AutoDownload:   q.Has("auto-download"),
		DownloadPolicy: data.ParseDownloadPolicy(q.Get("download-policy")),
		Expand:         q.Has("expand"),
		Force:          q.Has("force"),
	}

	return AddPlaylist(nil, playlistId, options)
}

func AddPlaylist(rdx redux.Writeable, playlistId string, opt *PlaylistOptions) error {

	apa := nod.Begin("adding playlist %s...", playlistId)
	defer apa.End()

	if opt == nil {
		opt = DefaultPlaylistOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return apa.EndWithError(err)
	}

	playlistId, err = yeti.ParsePlaylistId(playlistId)
	if err != nil {
		return apa.EndWithError(err)
	}

	propertyValues := make(map[string]map[string][]string)

	if opt.AutoRefresh {
		propertyValues[data.PlaylistAutoRefreshProperty] = map[string][]string{
			playlistId: {data.TrueValue},
		}
	}
	if opt.AutoDownload {
		propertyValues[data.PlaylistAutoDownloadProperty] = map[string][]string{
			playlistId: {data.TrueValue},
		}
	}
	if opt.DownloadPolicy != data.DefaultDownloadPolicy {
		propertyValues[data.PlaylistDownloadPolicyProperty] = map[string][]string{
			playlistId: {string(opt.DownloadPolicy)},
		}
	}
	if opt.Expand {
		propertyValues[data.PlaylistExpandProperty] = map[string][]string{
			playlistId: {data.TrueValue},
		}
	}

	for property, idValues := range propertyValues {
		if err := rdx.BatchAddValues(property, idValues); err != nil {
			return apa.EndWithError(err)
		}
	}

	if err := GetPlaylistsMetadata(rdx, opt, playlistId); err != nil {
		return apa.EndWithError(err)
	}

	apa.EndWithResult("done")

	return nil
}
