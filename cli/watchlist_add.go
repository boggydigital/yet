package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"strings"
)

func WatchlistAddHandler(u *url.URL) error {

	ids := strings.Split(u.Query().Get("id"), ",")
	return WatchlistAdd(ids)
}

func WatchlistAdd(ids []string) error {

	wlaa := nod.Begin("adding to watchlist...")
	defer wlaa.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return wlaa.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.VideosWatchlistProperty)
	if err != nil {
		return wlaa.EndWithError(err)
	}

	for _, id := range ids {
		if err := rxa.AddValues(data.VideosWatchlistProperty, id, data.TrueValue); err != nil {
			return wlaa.EndWithError(err)
		}
	}

	wlaa.EndWithResult("done")

	return nil
}
