package yeti

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"path/filepath"
	"strings"
)

func GetPosters(videoId string, dl *dolo.Client, force bool, qualities ...youtube_urls.ThumbnailQuality) error {

	gpa := nod.NewProgress(" posters for %s", videoId)
	defer gpa.End()

	gpa.TotalInt(len(qualities))

	for _, q := range qualities {

		u := youtube_urls.ThumbnailUrl(videoId, q)

		_, fnse := filepath.Split(u.Path)
		fnse = strings.TrimSuffix(fnse, filepath.Ext(fnse))

		if absFilename, err := paths.AbsPosterPath(videoId, q); err == nil {
			if err := dl.Download(u, force, nil, absFilename); err != nil {
				if lq := youtube_urls.LowerQuality(q); lq != youtube_urls.ThumbnailQualityUnknown {
					return GetPosters(videoId, dl, force, lq)
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
