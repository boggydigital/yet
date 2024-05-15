package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"strings"
)

func AddUrlsHandler(u *url.URL) error {
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

func AddUrls(propertyValues map[string][]string) error {

	aua := nod.NewProgress("adding urls...")
	defer aua.End()

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		return aua.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir,
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoEndedProperty)
	if err != nil {
		return aua.EndWithError(err)
	}

	aua.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := addPropertyValues(rdx, passthroughUrls, property, values...); err != nil {
			return aua.EndWithError(err)
		}
		aua.Increment()
	}

	aua.EndWithResult("done")

	return nil
}

func passthroughUrls(args ...string) ([]string, error) {
	return args, nil
}
