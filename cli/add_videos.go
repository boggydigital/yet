package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"golang.org/x/exp/maps"
	"net/url"
	"strings"
)

func AddVideosHandler(u *url.URL) error {
	q := u.Query()

	downloadQueue := strings.Split(q.Get("download-queue"), ",")
	watchlist := strings.Split(q.Get("watchlist"), ",")
	ended := strings.Split(q.Get("ended"), ",")
	raw := q.Has("raw")

	return AddVideos(map[string][]string{
		data.VideosDownloadQueueProperty: downloadQueue,
		data.VideosWatchlistProperty:     watchlist,
		data.VideoEndedProperty:          ended,
	}, raw)
}

func AddVideos(propertyValues map[string][]string, raw bool) error {

	ava := nod.NewProgress("adding videos...")
	defer ava.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return ava.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir,
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoEndedProperty)
	if err != nil {
		return ava.EndWithError(err)
	}

	ava.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := addPropertyValues(rxa, raw, property, values...); err != nil {
			return ava.EndWithError(err)
		}
		ava.Increment()
	}

	if !raw {
		// get metadata for the videos when adding them
		uniqueVideos := make(map[string]interface{})

		for _, values := range propertyValues {
			for _, v := range values {
				uniqueVideos[v] = nil
			}
		}

		if len(uniqueVideos) > 0 {
			if err := GetVideoMetadata(false, maps.Keys(uniqueVideos)...); err != nil {
				return ava.EndWithError(err)
			}
		}
	}

	ava.EndWithResult("done")

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
		if id == "" {
			continue
		}
		tv[id] = []string{data.TrueValue}
	}
	return tv
}
