package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathology"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"golang.org/x/exp/maps"
	"net/url"
	"strings"
	"time"
)

func AddVideosHandler(u *url.URL) error {
	q := u.Query()

	downloadQueue := strings.Split(q.Get("download-queue"), ",")
	watchlist := strings.Split(q.Get("watchlist"), ",")
	ended := strings.Split(q.Get("ended"), ",")

	return AddVideos(map[string][]string{
		data.VideosDownloadQueueProperty: downloadQueue,
		data.VideosWatchlistProperty:     watchlist,
		data.VideoEndedProperty:          ended,
	})
}

func AddVideos(propertyValues map[string][]string) error {

	ava := nod.NewProgress("adding videos...")
	defer ava.End()

	metadataDir, err := pathology.GetAbsDir(paths.Metadata)
	if err != nil {
		return ava.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir,
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoEndedProperty)
	if err != nil {
		return ava.EndWithError(err)
	}

	ava.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := addPropertyValues(rdx, yeti.ParseVideoIds, property, values...); err != nil {
			return ava.EndWithError(err)
		}
		ava.Increment()
	}

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

	ava.EndWithResult("done")

	return nil
}

func addPropertyValues(rdx kvas.WriteableRedux, parseDelegate func(...string) ([]string, error), property string, values ...string) error {
	apva := nod.Begin(" %s", property)
	defer apva.End()

	var err error
	if values, err = parseDelegate(values...); err != nil {
		return apva.EndWithError(err)
	}

	valuesDelegate := trueValues

	switch property {
	case data.VideoEndedProperty:
		valuesDelegate = timestampValues
	default:
		// do nothing, trueValues is already the default
	}

	if err := rdx.BatchAddValues(property, valuesDelegate(values...)); err != nil {
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

func timestampValues(ids ...string) map[string][]string {
	tv := make(map[string][]string)
	for _, id := range ids {
		if id == "" {
			continue
		}
		tv[id] = []string{time.Now().Format(time.RFC3339)}
	}
	return tv
}
