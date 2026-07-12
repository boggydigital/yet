package cli

import (
	"net/url"
	"strings"

	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func DownloadVideoHandler(u *url.URL) error {
	q := u.Query()

	videoIds := strings.Split(q.Get("video-id"), ",")

	options := &yeti.VideoOptions{
		BgUtilBaseUrl: q.Get("bgutil-baseurl"),
		Ended:         q.Has("mark-watched"),
		Verbose:       q.Has("verbose"),
		Force:         q.Has("force"),
	}

	return DownloadVideo(nil, options, videoIds...)
}

func DownloadVideo(rdx redux.Writeable, opt *yeti.VideoOptions, videoIds ...string) error {

	da := nod.NewProgress("downloading videos...")
	defer da.Done()

	if opt == nil {
		opt = yeti.DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return err
	}

	da.TotalInt(len(videoIds))

	for _, videoId := range videoIds {

		if err = yeti.DownloadVideoMetadataPoster(da, videoId, opt, rdx); err != nil {
			return err
		}

		da.Increment()
	}

	return nil
}
