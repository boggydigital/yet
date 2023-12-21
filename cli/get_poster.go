package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
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

	videoIds, err := yeti.ParseVideoIds(ids...)
	if err != nil {
		return gpa.EndWithError(err)
	}

	gpa.TotalInt(len(ids))

	for _, videoId := range videoIds {

		if err := yeti.GetPosters(videoId, dolo.DefaultClient, yt_urls.AllThumbnailQualities()...); err != nil {
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

	for _, q := range yt_urls.AllThumbnailQualities() {
		app, err := paths.AbsPosterPath(videoId, q)
		if err != nil {
			return err
		}

		if _, err := os.Stat(app); os.IsNotExist(err) {
			continue
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
