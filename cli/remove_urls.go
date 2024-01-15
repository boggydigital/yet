package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"strings"
)

func RemoveUrlsHandler(u *url.URL) error {
	q := u.Query()

	downloadQueue := strings.Split(q.Get("download-queue"), ",")
	watchlist := strings.Split(q.Get("watchlist"), ",")
	progress := strings.Split(q.Get("progress"), ",")
	ended := strings.Split(q.Get("ended"), ",")

	return RemoveUrls(map[string][]string{
		data.VideosDownloadQueueProperty: downloadQueue,
		data.VideosWatchlistProperty:     watchlist,
		data.VideoProgressProperty:       progress,
		data.VideoEndedProperty:          ended,
	})
}

func RemoveUrls(propertyValues map[string][]string) error {

	rva := nod.NewProgress("removing urls...")
	defer rva.End()

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		return rva.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir,
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoProgressProperty,
		data.VideoEndedProperty)
	if err != nil {
		return rva.EndWithError(err)
	}

	rva.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := removePropertyValues(rdx, passthroughUrls, property, values...); err != nil {
			return rva.EndWithError(err)
		}
		rva.Increment()
	}

	rva.EndWithResult("done")

	return nil
}
