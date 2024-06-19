package cli

import (
	"errors"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"golang.org/x/exp/maps"
	"net/url"
)

type playlistOptions struct {
	refresh      bool
	download     bool
	singleFormat bool
	expand       bool
}

func defaultPlaylistOptions() *playlistOptions {
	return &playlistOptions{
		refresh:      false,
		download:     false,
		singleFormat: true,
		expand:       false,
	}
}

func AddPlaylistsHandler(u *url.URL) error {
	q := u.Query()

	playlistId := q.Get("playlist-id")
	refresh := q.Has("refresh")
	download := q.Has("download")
	singleFormat := q.Has("single-format")
	expand := q.Has("expand")

	options := &playlistOptions{
		refresh:      refresh,
		download:     download,
		singleFormat: singleFormat,
		expand:       expand,
	}

	return AddPlaylists(nil, playlistId, options)
}

func AddPlaylists(rdx kvas.WriteableRedux, playlistId string, options *playlistOptions) error {

	apa := nod.Begin("adding playlist %s...", playlistId)
	defer apa.End()

	if rdx == nil {

		metadataDir, err := pathways.GetAbsDir(paths.Metadata)
		if err != nil {
			return apa.EndWithError(err)
		}

		rdx, err = kvas.NewReduxWriter(metadataDir, data.PlaylistProperties()...)
		if err != nil {
			return apa.EndWithError(err)
		}
	}

	// validate input playlist-id
	if ppids, err := yeti.ParsePlaylistIds(playlistId); err == nil {
		if len(ppids) > 0 {
			playlistId = ppids[0]
		} else {
			err = errors.New("invalid playlist id")
			return apa.EndWithError(err)
		}
	} else {
		return apa.EndWithError(err)
	}

	//for property, values := range propertyValues {
	//	if err := addPropertyValues(rdx, yeti.ParsePlaylistIds, property, values...); err != nil {
	//		return apa.EndWithError(err)
	//	}
	//	apa.Increment()
	//}

	// get metadata for the playlists when adding them
	uniquePlaylists := make(map[string]interface{})

	for _, values := range propertyValues {
		for _, v := range values {
			uniquePlaylists[v] = nil
		}
	}

	if len(uniquePlaylists) > 0 {
		if err := GetPlaylistMetadata(rdx, allVideos, false, maps.Keys(uniquePlaylists)...); err != nil {
			return apa.EndWithError(err)
		}
	}

	apa.EndWithResult("done")

	return nil
}
