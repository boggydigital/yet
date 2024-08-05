package cli

import (
	"fmt"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"time"
)

const maxAttempts = 3

func ProcessQueueHandler(u *url.URL) error {
	q := u.Query()

	options := &VideoOptions{
		PreferSingleFormat: q.Has("prefer-single-format"),
		Force:              q.Has("force"),
	}

	return ProcessQueue(nil, options)
}

// ProcessQueue processes download queue using the following rules:
// - download has not been completed after queue time
// - download is not in progress since queue time and less than 48 hours ago
func ProcessQueue(rdx kevlar.WriteableRedux, opt *VideoOptions) error {

	dqa := nod.NewProgress("processing videos queued for download...")
	defer dqa.End()

	if opt == nil {
		opt = DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return dqa.EndWithError(err)
	}

	processedVideoIds := make(map[string]int)

	for {
		videoId, err := getNextQueuedDownload(rdx, opt.Force)
		if err != nil {
			return dqa.EndWithError(err)
		}
		if videoId == "" {
			break
		}
		// this will serve as the final line of defence:
		// for some reason that would indicate that we're getting
		// the same videoId as earlier - which might happen (e.g. origin
		// server refused to serve video initially)
		// we'll attempt up to maxAttempts and break to avoid infinite loop
		// returning error to allow to get to the root cause if that happens
		if att, ok := processedVideoIds[videoId]; ok && att >= maxAttempts {
			return dqa.EndWithError(fmt.Errorf("exceeded max attempts for video %s", videoId))
		}
		processedVideoIds[videoId]++
		if err := DownloadVideo(rdx, videoId, opt); err != nil {
			return dqa.EndWithError(err)
		}
	}

	dqa.EndWithResult("done")

	return nil
}

// getNextQueuedDownload goes through queued downloads and returns the first one that:
// - was not completed after download was queued (earlier is fine, means it was added again)
// - has not started within the last 24 hours (allegedly in progress)
func getNextQueuedDownload(rdx kevlar.WriteableRedux, force bool) (string, error) {

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		return "", err
	}

	for _, id := range rdx.Keys(data.VideoDownloadQueuedProperty) {

		vdqTime := ""
		if vdq, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, id); ok && vdq != "" {
			vdqTime = vdq
		}

		// don't download videos that have ended _after_ queue time
		// instead, remove them from the queue
		if ved, ok := rdx.GetLastVal(data.VideoEndedDateProperty, id); ok && ved > vdqTime && !force {
			if err := rdx.CutKeys(data.VideoDownloadQueuedProperty, id); err != nil {
				return "", err
			}
			continue
		}
		// don't re-download videos that have download completed _after_ queue time
		if vdc, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, id); ok && vdc > vdqTime && !force {
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
				return "", err
			}
		}
		return id, nil
	}

	return "", nil
}
