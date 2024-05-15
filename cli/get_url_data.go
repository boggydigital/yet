package cli

import (
	"errors"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/url"
	"strings"
)

func GetUrlDataHandler(u *url.URL) error {
	id := u.Query().Get("id")
	videoId := u.Query().Get("video-id")
	lastDownloaded := u.Query().Has("last-downloaded")
	force := u.Query().Has("force")
	return GetUrlData(id, lastDownloaded, force, videoId)
}

func GetUrlData(id string, lastDownloaded, force bool, videoId string) error {

	// set id to the last downloaded URL
	if lastDownloaded {
		var err error
		if id, err = lastDownloadedId(); err != nil {
			return err
		}
	}

	guda := nod.Begin("getting url data for %s...", id)
	defer guda.End()

	if id == "" {
		return guda.EndWithError(errors.New("id is required"))
	}

	if err := GetVideoMetadata(id, true, videoId); err != nil {
		return guda.EndWithError(err)
	}

	if err := GetPoster(id, force, videoId); err != nil {
		return guda.EndWithError(err)
	}

	guda.EndWithResult("done")

	return nil
}

func lastDownloadedId() (string, error) {
	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		return "", err
	}

	rdx, err := kvas.NewReduxReader(metadataDir, data.VideoDownloadedDateProperty)
	if err != nil {
		return "", err
	}

	ids, err := rdx.Sort(rdx.Keys(data.VideoDownloadedDateProperty), true, data.VideoDownloadedDateProperty)
	if err != nil {
		return "", err
	}

	for _, sid := range ids {
		if strings.HasSuffix(sid, yt_urls.DefaultVideoExt) {
			return sid, nil
		}
	}

	return "", nil
}
