package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"time"
)

func DownloadQueueHandler(u *url.URL) error {
	q := u.Query()

	options := &VideoDownloadOptions{
		VideoOptions: &VideoOptions{
			PreferSingleFormat: q.Has("prefer-single-format"),
			Force:              q.Has("force"),
		},
	}

	return DownloadQueue(nil, options)
}

// DownloadQueue processes download queue using the following rules:
// - download is not already completed
// - download is not in progress since less than 48 hours ago
func DownloadQueue(rdx kvas.WriteableRedux, opt *VideoDownloadOptions) error {

	dqa := nod.NewProgress("downloading queued videos...")
	defer dqa.End()

	if opt == nil {
		opt = DefaultVideoDownloadOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return dqa.EndWithError(err)
	}

	queuedVideoIds := make([]string, 0)

	for _, id := range rdx.Keys(data.VideoDownloadQueuedProperty) {

		vdqTime := ""
		if vdq, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, id); ok && vdq != "" {
			vdqTime = vdq
		}

		// don't re-download videos that have download completed _after_ queue time
		if vdc, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, id); ok && vdc > vdqTime && !opt.Force {
			continue
		}
		// don't re-download videos that started download _after_ queue time and less than 48 hours ago
		if dss, ok := rdx.GetLastVal(data.VideoDownloadStartedProperty, id); ok && dss > vdqTime {
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
		if err := DownloadVideo(rdx, videoId, opt); err != nil {
			return dqa.EndWithError(err)
		}
	}

	dqa.EndWithResult("done")

	return nil
}
