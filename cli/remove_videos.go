package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func RemoveVideosHandler(u *url.URL) error {
	q := u.Query()

	downloadQueue := strings.Split(q.Get("download-queue"), ",")
	watchlist := strings.Split(q.Get("watchlist"), ",")
	progress := strings.Split(q.Get("progress"), ",")
	ended := strings.Split(q.Get("ended"), ",")
	skipped := strings.Split(q.Get("skipped"), ",")
	forcedDownload := strings.Split(q.Get("forced-download"), ",")
	singleFormat := strings.Split(q.Get("single-format"), ",")

	return RemoveVideos(map[string][]string{
		data.VideosDownloadQueueProperty:       downloadQueue,
		data.VideosWatchlistProperty:           watchlist,
		data.VideoProgressProperty:             progress,
		data.VideoEndedProperty:                ended,
		data.VideoSkippedProperty:              skipped,
		data.VideoForcedDownloadProperty:       forcedDownload,
		data.VideoSingleFormatDownloadProperty: singleFormat,
	})
}

func RemoveVideos(propertyValues map[string][]string) error {

	rva := nod.NewProgress("removing videos...")
	defer rva.End()

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		return rva.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir,
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoProgressProperty,
		data.VideoEndedProperty,
		data.VideoSkippedProperty,
		data.VideoForcedDownloadProperty,
		data.VideoSingleFormatDownloadProperty)
	if err != nil {
		return rva.EndWithError(err)
	}

	rva.TotalInt(len(propertyValues))

	for property, values := range propertyValues {
		if err := removePropertyValues(rdx, yeti.ParseVideoIds, property, values...); err != nil {
			return rva.EndWithError(err)
		}
		rva.Increment()
	}

	rva.EndWithResult("done")

	return nil
}

func removePropertyValues(rdx kvas.WriteableRedux, parseDelegate func(...string) ([]string, error), property string, values ...string) error {
	rpva := nod.Begin(" %s", property)
	defer rpva.End()

	var err error
	if values, err = parseDelegate(values...); err != nil {
		return rpva.EndWithError(err)
	}

	if err := rdx.CutKeys(property, values...); err != nil {
		return rpva.EndWithError(err)
	}

	result := "done "
	if len(values) > 0 {
		result += strings.Join(values, ",")
	}
	rpva.EndWithResult(result)
	return nil
}
