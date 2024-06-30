package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RemoveVideosHandler(u *url.URL) error {
	q := u.Query()

	videoId := q.Get("video-id")
	options := &VideoOptions{
		DownloadQueue:      q.Has("download-queue"),
		Progress:           q.Has("progress"),
		Ended:              q.Has("ended"),
		PreferSingleFormat: q.Has("prefer-single-format"),
		Force:              q.Has("force"),
	}

	return RemoveVideos(nil, videoId, options)
}

func RemoveVideos(rdx kvas.WriteableRedux, videoId string, opt *VideoOptions) error {

	rva := nod.Begin("removing video %s...", videoId)
	defer rva.End()

	if opt == nil {
		opt = DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return rva.EndWithError(err)
	}

	propertyKeys := make(map[string]string)

	if opt.DownloadQueue {
		propertyKeys[data.VideoDownloadQueuedProperty] = videoId
	}
	if opt.Progress {
		propertyKeys[data.VideoProgressProperty] = videoId
	}
	if opt.Ended {
		propertyKeys[data.VideoEndedDateProperty] = videoId
	}
	if opt.Reason != data.Unspecified {
		propertyKeys[data.VideoEndedReasonProperty] = videoId
	}
	if opt.Source != "" {
		propertyKeys[data.VideoSourceProperty] = videoId
	}
	if opt.PreferSingleFormat {
		propertyKeys[data.VideoPreferSingleFormatProperty] = videoId
	}
	if opt.Force {
		propertyKeys[data.VideoForcedDownloadProperty] = videoId
	}

	for property, key := range propertyKeys {
		if err := rdx.CutKeys(property, key); err != nil {
			return rva.EndWithError(err)
		}
	}

	rva.EndWithResult("done")

	return nil
}
