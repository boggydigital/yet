package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/url"
	"strings"
)

func GetVideoMetadataHandler(u *url.URL) error {
	q := u.Query()
	ids := strings.Split(q.Get("id"), ",")
	forId := q.Get("for-id")
	force := q.Has("force")
	return GetVideoMetadata(forId, force, ids...)
}

func GetVideoMetadata(forId string, force bool, ids ...string) error {
	gvma := nod.NewProgress("getting video metadata...")
	defer gvma.End()

	videoIds, err := yeti.ParseVideoIds(ids...)
	if err != nil {
		return gvma.EndWithError(err)
	}

	gvma.TotalInt(len(videoIds))

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		return gvma.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return gvma.EndWithError(err)
	}

	for _, videoId := range videoIds {

		if rdx.HasKey(data.VideoTitleProperty, videoId) && !force {
			continue
		}

		if err := getVideoPageMetadata(nil, videoId, rdx); err != nil {
			gvma.Error(err)
		} else {
			if err := copyMetadata(videoId, forId, rdx); err != nil {
				return gvma.EndWithError(err)
			}
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
		videoPage, err = yeti.GetVideoPage(videoId)
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

func copyMetadata(videoId, forId string, rdx kvas.WriteableRedux) error {

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
