package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
	"strings"
)

func GetVideoMetadataHandler(u *url.URL) error {
	q := u.Query()
	videoIds := strings.Split(q.Get("video-id"), ",")
	options := &VideoOptions{
		Force: q.Has("force"),
	}
	return GetVideoMetadata(nil, options, videoIds...)
}

func GetVideoMetadata(rdx redux.Writeable, opt *VideoOptions, videoIds ...string) error {
	gvma := nod.NewProgress("getting video metadata...")
	defer gvma.End()

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

		if err := getVideoPageMetadata(nil, videoId, rdx); err != nil {
			gvma.Error(err)
		}

		gvma.Increment()
	}

	gvma.EndWithResult("done")

	return nil
}

func getVideoPageMetadata(videoPage *youtube_urls.InitialPlayerResponse, videoId string, rdx redux.Writeable) error {

	gvpma := nod.Begin(" metadata for %s", videoId)
	defer gvpma.End()

	var err error
	if videoPage == nil {
		videoPage, err = yeti.GetVideoPage(videoId)
		if err != nil {
			return err
		}
	}

	for p, v := range yeti.ExtractMetadata(videoPage) {
		if err := rdx.AddValues(p, videoId, v...); err != nil {
			return err
		}
	}

	gvpma.EndWithResult("done")

	return nil
}

func copyMetadata(videoId, forId string, rdx redux.Writeable) error {

	if forId == "" || forId == videoId {
		return nil
	}

	for _, property := range data.AllProperties() {
		if values, ok := rdx.GetAllValues(property, videoId); ok && len(values) > 0 {
			if err := rdx.ReplaceValues(property, forId, values...); err != nil {
				return err
			}
		}
	}

	return nil
}
