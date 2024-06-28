package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"time"
)

func DownloadQueueHandler(u *url.URL) error {
	q := u.Query()

	options := &DownloadVideoOptions{
		PreferSingleFormat: q.Has("prefer-single-format"),
		Force:              q.Has("force"),
	}

	return DownloadQueue(nil, options)
}

// DownloadQueue processes download queue using the following rules:
// - download is not already completed
// - download is not in progress since less than 48 hours ago
func DownloadQueue(rdx kvas.WriteableRedux, options *DownloadVideoOptions) error {

	dqa := nod.NewProgress("downloading queued videos...")
	defer dqa.End()

	if options == nil {
		options = DefaultDownloadVideoOptions()
	}

	if rdx == nil {
		metadataDir, err := pathways.GetAbsDir(paths.Metadata)
		if err != nil {
			return dqa.EndWithError(err)
		}

		rdx, err = kvas.NewReduxWriter(metadataDir, data.VideoProperties()...)
		if err != nil {
			return dqa.EndWithError(err)
		}
	} else if err := rdx.MustHave(data.VideoProperties()...); err != nil {
		return dqa.EndWithError(err)
	}

	queuedVideoIds := make([]string, 0)

	for _, id := range rdx.Keys(data.VideoDownloadQueuedProperty) {
		// don't re-download completed videos
		if rdx.HasKey(data.VideoDownloadCompletedProperty, id) && !options.Force {
			continue
		}
		// don't re-download videos that started download less than 48 hours ago
		if dss, ok := rdx.GetLastVal(data.VideoDownloadStartedProperty, id); ok {
			if ds, err := time.Parse(time.RFC3339, dss); err == nil {
				dur := time.Now().Sub(ds)
				if dur < yeti.DefaultDelay {
					continue
				}
			} else {
				return dqa.EndWithError(err)
			}
		}
		queuedVideoIds = append(queuedVideoIds, id)
	}

	for _, videoId := range queuedVideoIds {
		if err := DownloadVideo(rdx, videoId, options); err != nil {
			return dqa.EndWithError(err)
		}
	}

	dqa.EndWithResult("done")

	return nil
}
