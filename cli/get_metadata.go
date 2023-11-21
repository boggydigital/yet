package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"strings"
)

func GetMetadataHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	return GetMetadata(ids...)
}

func GetMetadata(ids ...string) error {
	gma := nod.NewProgress("getting metadata...")
	defer gma.End()

	gma.TotalInt(len(ids))

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return gma.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.AllProperties()...)
	if err != nil {
		return gma.EndWithError(err)
	}

	for _, videoId := range ids {

		videoPage, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			gma.Error(err)
			gma.Increment()
			continue
		}

		for p, v := range yeti.ExtractMetadata(videoPage) {
			if err := rxa.AddValues(p, videoId, v...); err != nil {
				return gma.EndWithError(err)
			}
		}

		gma.Increment()
	}

	gma.EndWithResult("done")

	return nil
}
