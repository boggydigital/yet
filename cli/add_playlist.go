package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

type addPlaylistOptions struct {
	autoRefresh        bool
	autoDownload       bool
	downloadPolicy     data.PlaylistDownloadPolicy
	preferSingleFormat bool
	expand             bool
	force              bool
}

func defaultAddPlaylistOptions() *addPlaylistOptions {
	return &addPlaylistOptions{
		autoRefresh:        false,
		autoDownload:       false,
		downloadPolicy:     data.Recent,
		preferSingleFormat: true,
		expand:             false,
		force:              false,
	}
}

func AddPlaylistHandler(u *url.URL) error {
	q := u.Query()

	playlistId := q.Get("playlist-id")
	options := &addPlaylistOptions{
		autoRefresh:        q.Has("auto-refresh"),
		autoDownload:       q.Has("auto-download"),
		downloadPolicy:     data.ParsePlaylistDownloadPolicy(q.Get("download-policy")),
		preferSingleFormat: q.Has("prefer-single-format"),
		expand:             q.Has("expand"),
		force:              q.Has("force"),
	}

	return AddPlaylist(nil, playlistId, options)
}

func AddPlaylist(rdx kvas.WriteableRedux, playlistId string, options *addPlaylistOptions) error {

	apa := nod.Begin("adding playlist %s...", playlistId)
	defer apa.End()

	if options == nil {
		options = defaultAddPlaylistOptions()
	}

	if rdx == nil {
		metadataDir, err := pathways.GetAbsDir(paths.Metadata)
		if err != nil {
			return apa.EndWithError(err)
		}

		rdx, err = kvas.NewReduxWriter(metadataDir, data.PlaylistProperties()...)
		if err != nil {
			return apa.EndWithError(err)
		}
	} else if err := rdx.MustHave(data.PlaylistProperties()...); err != nil {
		return apa.EndWithError(err)
	}

	var err error
	playlistId, err = yeti.ParsePlaylistId(playlistId)
	if err != nil {
		return apa.EndWithError(err)
	}

	propertyValues := make(map[string]map[string][]string)

	if options.autoRefresh {
		propertyValues[data.PlaylistAutoRefreshProperty] = map[string][]string{
			playlistId: {data.TrueValue},
		}
	}
	if options.autoDownload {
		propertyValues[data.PlaylistAutoDownloadProperty] = map[string][]string{
			playlistId: {data.TrueValue},
		}
	}
	if options.downloadPolicy != data.Unset {
		propertyValues[data.PlaylistDownloadPolicyProperty] = map[string][]string{
			playlistId: {string(options.downloadPolicy)},
		}
	}
	if options.preferSingleFormat {
		propertyValues[data.PlaylistPreferSingleFormatProperty] = map[string][]string{
			playlistId: {data.TrueValue},
		}
	}

	for property, idValues := range propertyValues {
		if err := rdx.BatchAddValues(property, idValues); err != nil {
			return apa.EndWithError(err)
		}
	}

	if err := GetPlaylistMetadata(rdx, options.expand, options.force, playlistId); err != nil {
		return apa.EndWithError(err)
	}

	apa.EndWithResult("done")

	return nil
}
