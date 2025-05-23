package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RemoveVideosHandler(u *url.URL) error {
	q := u.Query()

	videoId := q.Get("video-id")
	options := &VideoOptions{
		DownloadQueue: q.Has("download-queue"),
		Progress:      q.Has("progress"),
		Ended:         q.Has("ended"),
		Force:         q.Has("force"),
	}

	return RemoveVideos(nil, videoId, options)
}

func RemoveVideos(rdx redux.Writeable, videoId string, opt *VideoOptions) error {

	rva := nod.Begin("removing video %s...", videoId)
	defer rva.Done()

	if opt == nil {
		opt = DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return err
	}

	propertyKeys := make(map[string]string)

	if opt.Favorite {
		propertyKeys[data.VideoFavoriteProperty] = videoId
	}
	if opt.DownloadQueue {
		propertyKeys[data.VideoDownloadQueuedProperty] = videoId
	}
	if opt.Progress {
		propertyKeys[data.VideoProgressProperty] = videoId
	}
	if opt.Ended {
		propertyKeys[data.VideoEndedDateProperty] = videoId
	}
	if opt.Reason != data.DefaultEndedReason {
		propertyKeys[data.VideoEndedReasonProperty] = videoId
	}
	if opt.Force {
		propertyKeys[data.VideoForcedDownloadProperty] = videoId
	}

	for property, key := range propertyKeys {
		if err := rdx.CutKeys(property, key); err != nil {
			return err
		}
	}

	return nil
}
