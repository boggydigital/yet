package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"strings"
)

func GetPosterHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	return GetPoster(ids)
}

func GetPoster(ids []string) error {

	gpa := nod.NewProgress("getting poster(s)...")
	defer gpa.End()

	gpa.TotalInt(len(ids))

	dl := dolo.DefaultClient

	for _, videoId := range ids {

		videoPage, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			gpa.Error(err)
			gpa.Increment()
			continue
		}

		if err := yeti.GetThumbnails(dl, videoId, videoPage.VideoDetails.Thumbnail.Thumbnails); err != nil {
			gpa.Error(err)
		}

		gpa.Increment()
	}

	gpa.EndWithResult("done")

	return nil
}
