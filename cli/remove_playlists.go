package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func RemovePlaylistsHandler(u *url.URL) error {
	q := u.Query()

	watchlist := strings.Split(q.Get("watchlist"), ",")

	return RemovePlaylists(map[string][]string{
		data.VideosWatchlistProperty: watchlist,
	})
}

func RemovePlaylists(propertyValues map[string][]string) error {
	rpa := nod.NewProgress("removing playlists...")
	defer rpa.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return rpa.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir,
		data.PlaylistWatchlistProperty)
	if err != nil {
		return rpa.EndWithError(err)
	}

	rpa.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := removePropertyValues(rxa, yeti.ParsePlaylistIds, property, values...); err != nil {
			return rpa.EndWithError(err)
		}
		rpa.Increment()
	}

	rpa.EndWithResult("done")

	return nil
}
