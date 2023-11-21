package yeti

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/url"
	"path/filepath"
	"strings"
)

func GetPosters(dl *dolo.Client, videoId string, thumbnails []yt_urls.Thumbnail) error {

	remains := map[string]bool{
		paths.PosterQualityMax:  true,
		paths.PosterQualityHigh: true,
	}

	for ii := len(thumbnails) - 1; ii >= 0; ii-- {

		thumbnail := thumbnails[ii]

		u, err := url.Parse(thumbnail.Url)
		if err != nil {
			return err
		}

		_, fnse := filepath.Split(u.Path)
		fnse = strings.TrimSuffix(fnse, filepath.Ext(fnse))

		if absFilename, err := paths.AbsPosterPath(videoId, fnse); err == nil {
			if err := dl.Download(u, nil, absFilename); err != nil {
				return err
			} else {
				remains[fnse] = false
				moreDownloadsRemain := false
				for _, v := range remains {
					moreDownloadsRemain = moreDownloadsRemain || v
				}
				if !moreDownloadsRemain {
					return nil
				}
			}
		} else {
			return err
		}
	}

	return nil
}
