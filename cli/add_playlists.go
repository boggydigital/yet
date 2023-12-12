package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"golang.org/x/exp/maps"
	"net/url"
	"strings"
)

func AddPlaylistsHandler(u *url.URL) error {
	q := u.Query()

	watchlist := strings.Split(q.Get("watchlist"), ",")
	allVideos := q.Has("all-videos")

	return AddPlaylists(allVideos, map[string][]string{
		data.PlaylistWatchlistProperty: watchlist,
	})
}

func AddPlaylists(allVideos bool, propertyValues map[string][]string) error {

	apa := nod.NewProgress("adding playlists...")
	defer apa.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return apa.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir,
		data.PlaylistWatchlistProperty)
	if err != nil {
		return apa.EndWithError(err)
	}

	apa.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := addPropertyValues(rxa, true, property, values...); err != nil {
			return apa.EndWithError(err)
		}
		apa.Increment()
	}

	// get metadata for the playlists when adding them
	uniquePlaylists := make(map[string]interface{})

	for _, values := range propertyValues {
		for _, v := range values {
			uniquePlaylists[v] = nil
		}
	}

	if len(uniquePlaylists) > 0 {
		if err := GetPlaylistMetadata(allVideos, false, maps.Keys(uniquePlaylists)...); err != nil {
			return apa.EndWithError(err)
		}
	}

	apa.EndWithResult("done")

	return nil
}
