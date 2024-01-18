package cli

import (
	"github.com/boggydigital/nod"
	"net/url"
)

func GetUrlDataHandler(u *url.URL) error {
	id := u.Query().Get("id")
	videoId := u.Query().Get("video-id")
	return GetUrlData(id, videoId)
}

func GetUrlData(id, videoId string) error {

	guda := nod.Begin("getting url data...")
	defer guda.End()

	if err := GetVideoMetadata(id, true, videoId); err != nil {
		return guda.EndWithError(err)
	}

	if err := GetPoster(id, videoId); err != nil {
		return guda.EndWithError(err)
	}

	guda.EndWithResult("done")

	return nil
}
