package main

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"os"
	"path/filepath"
)

func processList(list map[string][]string) error {

	la := nod.NewProgress("updating %d playlists...", len(list))
	defer la.End()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for dir, ids := range list {

		if len(ids) == 0 {
			continue
		}

		da := nod.Begin("processing %s...", dir)

		if err := os.Chdir(pwd); err != nil {
			return da.EndWithError(err)
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return da.EndWithError(err)
			}
		}

		if err := os.Chdir(filepath.Clean(dir)); err != nil {
			return da.EndWithError(err)
		}

		jar, err := coost.NewJar([]string{yt_urls.YoutubeHost}, "")
		if err != nil {
			return da.EndWithError(err)
		}

		httpClient := jar.NewHttpClient()

		videoIds, err := argsToVideoIds(httpClient, true, ids...)
		if err != nil {
			return da.EndWithError(err)
		}

		if err := DownloadVideos(httpClient, videoIds...); err != nil {
			return da.EndWithError(err)
		}

		la.Increment()
		da.End()
	}
	return nil
}
