package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"strings"
)

func RemoveHandler(u *url.URL) error {
	q := u.Query()

	downloadQueue := strings.Split(q.Get("download-queue"), ",")
	watchlist := strings.Split(q.Get("watchlist"), ",")
	progress := strings.Split(q.Get("progress"), ",")
	ended := strings.Split(q.Get("ended"), ",")

	return Add(map[string][]string{
		data.VideosDownloadQueueProperty: downloadQueue,
		data.VideosWatchlistProperty:     watchlist,
		data.VideoProgressProperty:       progress,
		data.VideoEndedProperty:          ended,
	})
}

func Remove(propertyValues map[string][]string) error {

	wlra := nod.NewProgress("removing...")
	defer wlra.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return wlra.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir,
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoProgressProperty,
		data.VideoEndedProperty)
	if err != nil {
		return wlra.EndWithError(err)
	}

	wlra.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := removePropertyValues(rxa, property, values...); err != nil {
			return wlra.EndWithError(err)
		}
		wlra.Increment()
	}

	wlra.EndWithResult("done")

	return nil
}

func removePropertyValues(rxa kvas.ReduxAssets, property string, values ...string) error {
	rpva := nod.Begin(" %s", property)
	defer rpva.End()

	return rxa.BatchCutValues(property, trueValues(values...))
}
