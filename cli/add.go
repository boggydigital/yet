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

func AddHandler(u *url.URL) error {
	q := u.Query()

	downloadQueue := strings.Split(q.Get("download-queue"), ",")
	watchlist := strings.Split(q.Get("watchlist"), ",")
	ended := strings.Split(q.Get("ended"), ",")
	raw := q.Has("raw")

	return Add(map[string][]string{
		data.VideosDownloadQueueProperty: downloadQueue,
		data.VideosWatchlistProperty:     watchlist,
		data.VideoEndedProperty:          ended,
	}, raw)
}

func Add(propertyValues map[string][]string, raw bool) error {

	wlaa := nod.NewProgress("adding...")
	defer wlaa.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return wlaa.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir,
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoEndedProperty)
	if err != nil {
		return wlaa.EndWithError(err)
	}

	wlaa.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := addPropertyValues(rxa, raw, property, values...); err != nil {
			return wlaa.EndWithError(err)
		}
		wlaa.Increment()
	}

	wlaa.EndWithResult("done")

	return nil
}

func addPropertyValues(rxa kvas.ReduxAssets, raw bool, property string, values ...string) error {
	apva := nod.Begin(" %s", property)
	defer apva.End()

	if !raw {
		var err error
		if values, err = yeti.ParseVideoIds(values...); err != nil {
			return apva.EndWithError(err)
		}
	}

	if err := rxa.BatchAddValues(property, trueValues(values...)); err != nil {
		return apva.EndWithError(err)
	}

	result := "done "
	if len(values) > 0 {
		result += strings.Join(values, ",")
	}
	apva.EndWithResult(result)
	return nil
}

func trueValues(ids ...string) map[string][]string {
	tv := make(map[string][]string)
	for _, id := range ids {
		tv[id] = []string{data.TrueValue}
	}
	return tv
}
