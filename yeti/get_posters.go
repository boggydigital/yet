package yeti

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"path/filepath"
	"strings"
)

func GetPosters(videoId string, dl *dolo.Client, qualities ...yt_urls.ThumbnailQuality) error {

	gpa := nod.NewProgress(" posters for %s", videoId)
	defer gpa.End()

	gpa.TotalInt(len(qualities))

	for _, q := range qualities {

		u := yt_urls.ThumbnailUrl(videoId, q)

		_, fnse := filepath.Split(u.Path)
		fnse = strings.TrimSuffix(fnse, filepath.Ext(fnse))

		if absFilename, err := paths.AbsPosterPath(videoId, q); err == nil {
			if err := dl.Download(u, nil, absFilename); err != nil {
				if lq := yt_urls.LowerQuality(q); lq != yt_urls.ThumbnailQualityUnknown {
					return GetPosters(videoId, dl, lq)
				} else {
					return gpa.EndWithError(err)
				}
			}
		} else {
			return gpa.EndWithError(err)
		}

		gpa.Increment()
	}

	gpa.EndWithResult("done")

	return nil
}
