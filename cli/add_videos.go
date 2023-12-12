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

	if !raw {
		// get metadata for the videos, playlists upon adding them
		uniqueVideos := make(map[string]interface{})
		uniquePlaylists := make(map[string]interface{})

		var unique map[string]interface{}
		for property, values := range propertyValues {
			switch property {
			case data.PlaylistWatchlistProperty:
				unique = uniquePlaylists
			default:
				unique = uniqueVideos
			}
			for _, v := range values {
				unique[v] = nil
			}
		}

		if len(uniqueVideos) > 0 {
			if err := GetVideoMetadata(false, maps.Keys(uniqueVideos)...); err != nil {
				return wlaa.EndWithError(err)
			}
		}
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
		if id == "" {
			continue
		}
		tv[id] = []string{data.TrueValue}
	}
	return tv
}
