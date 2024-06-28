package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

type addVideoOptions struct {
	downloadQueue      bool
	ended              bool
	reason             data.VideoEndedReason
	source             string
	preferSingleFormat bool
	force              bool
}

func defaultAddVideoOptions() *addVideoOptions {
	return &addVideoOptions{
		downloadQueue:      true,
		ended:              false,
		reason:             data.Unspecified,
		source:             "",
		preferSingleFormat: true,
		force:              false,
	}
}

func AddVideoHandler(u *url.URL) error {
	q := u.Query()

	videoId := q.Get("video-id")
	options := &addVideoOptions{
		downloadQueue:      q.Has("download-queue"),
		ended:              q.Has("ended"),
		reason:             data.ParseVideoEndedReason(q.Get("reason")),
		source:             q.Get("source"),
		preferSingleFormat: q.Has("prefer-single-format"),
		force:              q.Has("force"),
	}

	return AddVideo(nil, videoId, options)
}

func AddVideo(rdx kvas.WriteableRedux, videoId string, options *addVideoOptions) error {

	ava := nod.Begin("adding video %s...", videoId)
	defer ava.End()

	if options == nil {
		options = defaultAddVideoOptions()
	}

	if rdx == nil {
		metadataDir, err := pathways.GetAbsDir(paths.Metadata)
		if err != nil {
			return ava.EndWithError(err)
		}

		rdx, err = kvas.NewReduxWriter(metadataDir, data.VideoProperties()...)
		if err != nil {
			return ava.EndWithError(err)
		}

	} else if err := rdx.MustHave(data.VideoProperties()...); err != nil {
		return ava.EndWithError(err)
	}

	var err error
	videoId, err = yeti.ParseVideoId(videoId)
	if err != nil {
		return ava.EndWithError(err)
	}

	propertyValues := make(map[string]map[string][]string)

	if options.downloadQueue {
		propertyValues[data.VideoDownloadQueuedProperty] = map[string][]string{
			videoId: {yeti.FmtNow()},
		}
	}
	if options.ended {
		propertyValues[data.VideoEndedDateProperty] = map[string][]string{
			videoId: {yeti.FmtNow()},
		}
	}
	if options.reason != data.Unspecified {
		propertyValues[data.VideoEndedReasonProperty] = map[string][]string{
			videoId: {string(options.reason)},
		}
	}
	if options.source != "" {
		propertyValues[data.VideoSourceProperty] = map[string][]string{
			videoId: {options.source},
		}
	}
	if options.preferSingleFormat {
		propertyValues[data.VideoPreferSingleFormatProperty] = map[string][]string{
			videoId: {data.TrueValue},
		}
	}
	if options.force {
		propertyValues[data.VideoForcedDownloadProperty] = map[string][]string{
			videoId: {data.TrueValue},
		}
	}

	for property, idValues := range propertyValues {
		if err := rdx.BatchAddValues(property, idValues); err != nil {
			return ava.EndWithError(err)
		}
	}

	if err := GetVideoMetadata(options.force, videoId); err != nil {
		return ava.EndWithError(err)
	}

	ava.EndWithResult("done")

	return nil
}
