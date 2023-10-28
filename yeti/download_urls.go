package yeti

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"net/http"
	"net/url"
	"time"
)

func DownloadUrls(httpClient *http.Client, urls ...string) error {

	if len(urls) == 0 {
		return nil
	}

	dftpw := nod.NewProgress(fmt.Sprintf("downloading %d file(s)", len(urls)))
	defer dftpw.End()

	dftpw.Total(uint64(len(urls)))

	dl := dolo.NewClient(httpClient, dolo.Defaults())

	for _, rawUrl := range urls {

		u, err := url.Parse(rawUrl)
		if err != nil {
			return dftpw.EndWithError(err)
		}

		gv := nod.NewProgress("file: " + u.Path)

		start := time.Now()

		absVideosDir, err := paths.GetAbsDir(paths.Videos)
		if err != nil {
			return dftpw.EndWithError(err)
		}

		if err := dl.Download(u, gv, absVideosDir); err != nil {
			return dftpw.EndWithError(err)
		}

		elapsed := time.Since(start)

		gv.EndWithResult("done in %.1fs", elapsed.Seconds())
		dftpw.Increment()
	}

	return nil
}
