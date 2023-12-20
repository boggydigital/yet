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

func GetVideoMetadataHandler(u *url.URL) error {
	q := u.Query()
	ids := strings.Split(q.Get("id"), ",")
	force := q.Has("force")
	return GetVideoMetadata(force, ids...)
}

func GetVideoMetadata(force bool, ids ...string) error {
	gvma := nod.NewProgress("getting video metadata...")
	defer gvma.End()

	gvma.TotalInt(len(ids))

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return gvma.EndWithError(err)
	}

	rdx, err := kvas.ReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return gvma.EndWithError(err)
	}

	for _, videoId := range ids {

		if rdx.HasKey(data.VideoTitleProperty, videoId) && !force {
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

func getVideoPageMetadata(videoPage *yt_urls.InitialPlayerResponse, videoId string, rdx kvas.WriteableRedux) error {

	gvpma := nod.Begin(" metadata for %s", videoId)
	defer gvpma.End()

	var err error
	if videoPage == nil {
		videoPage, err = yt_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			return gvpma.EndWithError(err)
		}
	}

	for p, v := range yeti.ExtractMetadata(videoPage) {
		if err := rdx.AddValues(p, videoId, v...); err != nil {
			return gvpma.EndWithError(err)
		}
	}

	gvpma.EndWithResult("done")

	return nil
}
