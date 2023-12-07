package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func DownloadHandler(u *url.URL) error {

	ids := strings.Split(u.Query().Get("id"), ",")
	force := u.Query().Has("force")
	return Download(ids, force)
}

func Download(ids []string, force bool) error {

	da := nod.Begin("downloading videos...")
	defer da.End()

	videoIds, err := yeti.ParseVideoIds(ids...)
	if err != nil {
		return da.EndWithError(err)
	}

	if err := GetVideo(force, videoIds...); err != nil {
		return da.EndWithError(err)
	}

	if err := GetPoster(videoIds, ""); err != nil {
		return da.EndWithError(err)
	}

	da.EndWithResult("done")

	return nil
}
