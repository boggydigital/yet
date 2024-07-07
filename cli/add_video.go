package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func AddVideoHandler(u *url.URL) error {
	q := u.Query()

	videoId := q.Get("video-id")
	options := &VideoOptions{
		Favorite:           q.Has("favorite"),
		DownloadQueue:      q.Has("download-queue"),
		Ended:              q.Has("ended"),
		Reason:             data.ParseVideoEndedReason(q.Get("reason")),
		Source:             q.Get("source"),
		PreferSingleFormat: q.Has("prefer-single-format"),
		Force:              q.Has("force"),
	}

	return AddVideo(nil, videoId, options)
}

func AddVideo(rdx kevlar.WriteableRedux, videoId string, opt *VideoOptions) error {

	ava := nod.Begin("adding video %s...", videoId)
	defer ava.End()

	if opt == nil {
		opt = DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return ava.EndWithError(err)
	}

	videoId, err = yeti.ParseVideoId(videoId)
	if err != nil {
		return ava.EndWithError(err)
	}

	propertyValues := make(map[string]map[string][]string)

	if opt.Favorite {
		propertyValues[data.VideoFavoriteProperty] = map[string][]string{
			videoId: {data.TrueValue},
		}
	}
	if opt.DownloadQueue {
		propertyValues[data.VideoDownloadQueuedProperty] = map[string][]string{
			videoId: {yeti.FmtNow()},
		}
	}
	if opt.Ended {
		propertyValues[data.VideoEndedDateProperty] = map[string][]string{
			videoId: {yeti.FmtNow()},
		}
	}
	if opt.Reason != data.DefaultEndedReason {
		propertyValues[data.VideoEndedReasonProperty] = map[string][]string{
			videoId: {string(opt.Reason)},
		}
	}
	if opt.Source != "" {
		propertyValues[data.VideoSourceProperty] = map[string][]string{
			videoId: {opt.Source},
		}
	}
	if opt.PreferSingleFormat {
		propertyValues[data.VideoPreferSingleFormatProperty] = map[string][]string{
			videoId: {data.TrueValue},
		}
	}
	if opt.Force {
		propertyValues[data.VideoForcedDownloadProperty] = map[string][]string{
			videoId: {data.TrueValue},
		}
	}

	for property, idValues := range propertyValues {
		if err := rdx.BatchAddValues(property, idValues); err != nil {
			return ava.EndWithError(err)
		}
	}

	if err := GetVideoMetadata(rdx, opt, videoId); err != nil {
		return ava.EndWithError(err)
	}

	ava.EndWithResult("done")

	return nil
}
