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

func RemoveHandler(u *url.URL) error {
	q := u.Query()

	downloadQueue := strings.Split(q.Get("download-queue"), ",")
	watchlist := strings.Split(q.Get("watchlist"), ",")
	progress := strings.Split(q.Get("progress"), ",")
	ended := strings.Split(q.Get("ended"), ",")
	raw := q.Has("raw")

	return Remove(map[string][]string{
		data.VideosDownloadQueueProperty: downloadQueue,
		data.VideosWatchlistProperty:     watchlist,
		data.VideoProgressProperty:       progress,
		data.VideoEndedProperty:          ended,
	}, raw)
}

func Remove(propertyValues map[string][]string, raw bool) error {

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
		if err := removePropertyValues(rxa, raw, property, values...); err != nil {
			return wlra.EndWithError(err)
		}
		wlra.Increment()
	}

	wlra.EndWithResult("done")

	return nil
}

func removePropertyValues(rxa kvas.ReduxAssets, raw bool, property string, values ...string) error {
	rpva := nod.Begin(" %s", property)
	defer rpva.End()

	if !raw {
		var err error
		if values, err = yeti.ParseVideoIds(values...); err != nil {
			return rpva.EndWithError(err)
		}
	}

	if err := rxa.BatchCutValues(property, trueValues(values...)); err != nil {
		return rpva.EndWithError(err)
	}

	result := "done "
	if len(values) > 0 {
		result += strings.Join(values, ",")
	}
	rpva.EndWithResult(result)
	return nil
}
