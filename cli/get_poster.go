package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetPosterHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	forId := u.Query().Get("for-id")
	return GetPoster(ids, forId)
}

func GetPoster(ids []string, forId string) error {

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

		if err := yeti.GetPosters(dl, videoId, videoPage.VideoDetails.Thumbnail.Thumbnails); err != nil {
			gpa.Error(err)
		} else {
			if err := renamePosters(videoId, forId); err != nil {
				return gpa.EndWithError(err)
			}
		}

		gpa.Increment()
	}

	gpa.EndWithResult("done")

	return nil
}

func renamePosters(videoId, forId string) error {

	if forId == "" || forId == videoId {
		return nil
	}

	for _, q := range []string{paths.PosterQualityHigh, paths.PosterQualityMax} {
		app, err := paths.AbsPosterPath(videoId, q)
		if err != nil {
			return err
		}

		napp, err := paths.AbsPosterPath(forId, q)
		if err != nil {
			return err
		}

		if err := os.Rename(app, napp); err != nil {
			return err
		}
	}

	return nil
}
