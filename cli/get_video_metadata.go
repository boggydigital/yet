package cli

import (
	"net/url"
	"strings"

	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func GetVideoMetadataHandler(u *url.URL) error {
	q := u.Query()
	videoIds := strings.Split(q.Get("video-id"), ",")
	options := &yeti.VideoOptions{
		Force: q.Has("force"),
	}
	return GetVideoMetadata(nil, options, videoIds...)
}

func GetVideoMetadata(rdx redux.Writeable, opt *yeti.VideoOptions, videoIds ...string) error {
	gvma := nod.NewProgress("getting video metadata...")
	defer gvma.Done()

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)

	parsedVideoIds, err := yeti.ParseVideoIds(videoIds...)
	if err != nil {
		return err
	}

	gvma.TotalInt(len(parsedVideoIds))

	for _, videoId := range parsedVideoIds {

		if rdx.HasKey(data.VideoTitleProperty, videoId) && !opt.Force {
			continue
		}

		if err = yeti.GetVideoPageMetadata(nil, videoId, rdx); err != nil {
			gvma.Error(err)
		}

		gvma.Increment()
	}

	return nil
}
