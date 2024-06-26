package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
	"strings"
)

func GetPosterHandler(u *url.URL) error {
	q := u.Query()
	videoIds := strings.Split(q.Get("video-id"), ",")
	force := q.Has("force")
	return GetPoster(force, videoIds...)
}

func GetPoster(force bool, videoIds ...string) error {

	gpa := nod.NewProgress("getting poster(s)...")
	defer gpa.End()

	parsedVideoIds, err := yeti.ParseVideoIds(videoIds...)
	if err != nil {
		return gpa.EndWithError(err)
	}

	gpa.TotalInt(len(parsedVideoIds))

	for _, videoId := range parsedVideoIds {

		if err := yeti.GetPosters(videoId, dolo.DefaultClient, force, youtube_urls.AllThumbnailQualities()...); err != nil {
			gpa.Error(err)
		}

		gpa.Increment()
	}

	gpa.EndWithResult("done")

	return nil
}
