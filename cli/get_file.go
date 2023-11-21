package cli

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func GetFileHandler(u *url.URL) error {
	urls := strings.Split(u.Query().Get("url"), ",")
	return GetFile(urls...)
}

func GetFile(urls ...string) error {

	if len(urls) == 0 {
		return nil
	}

	gfa := nod.NewProgress(fmt.Sprintf("downloading %d file(s)", len(urls)))
	defer gfa.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return gfa.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.AllProperties()...)
	if err != nil {
		return gfa.EndWithError(err)
	}

	gfa.Total(uint64(len(urls)))

	dl := dolo.DefaultClient

	for _, rawUrl := range urls {

		u, err := url.Parse(rawUrl)
		if err != nil {
			return gfa.EndWithError(err)
		}

		_, filename := filepath.Split(u.Path)

		gv := nod.NewProgress("file: " + filename)

		start := time.Now()

		// add to the download queue
		if err := rxa.AddValues(data.VideosDownloadQueueProperty, filename, data.TrueValue); err != nil {
			return gfa.EndWithError(err)
		}

		absVideosDir, err := paths.GetAbsDir(paths.Videos)
		if err != nil {
			return gfa.EndWithError(err)
		}

		if err := dl.Download(u, gv, absVideosDir, filename); err != nil {
			return gfa.EndWithError(err)
		}

		// clear from the queue upon successful download
		if err := rxa.CutVal(data.VideosDownloadQueueProperty, filename, data.TrueValue); err != nil {
			return gfa.EndWithError(err)
		}

		// add to the watchlist upon successful download
		if err := rxa.AddValues(data.VideosWatchlistProperty, filename, data.TrueValue); err != nil {
			return gfa.EndWithError(err)
		}

		elapsed := time.Since(start)

		gv.EndWithResult("done in %.1fs", elapsed.Seconds())
		gfa.Increment()
	}

	return nil
}
