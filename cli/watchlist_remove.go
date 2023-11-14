package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"strings"
)

func WatchlistRemoveHandler(u *url.URL) error {

	ids := strings.Split(u.Query().Get("id"), ",")
	return WatchlistRemove(ids)
}

func WatchlistRemove(ids []string) error {

	wlra := nod.Begin("removing from watchlist...")
	defer wlra.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return wlra.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.VideosWatchlistProperty)
	if err != nil {
		return wlra.EndWithError(err)
	}

	for _, id := range ids {
		if err := rxa.CutVal(data.VideosWatchlistProperty, id, data.TrueValue); err != nil {
			return wlra.EndWithError(err)
		}
	}

	wlra.EndWithResult("done")

	return nil
}
