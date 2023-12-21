package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/url"
	"os"
	"path/filepath"
)

func CleanupEndedHandler(u *url.URL) error {
	return CleanupEnded()
}

func CleanupEnded() error {

	cea := nod.NewProgress("cleaning up ended media...")
	defer cea.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return cea.EndWithError(err)
	}

	absVideosDir, err := paths.GetAbsDir(paths.Videos)
	if err != nil {
		return cea.EndWithError(err)
	}

	rdx, err := kvas.ReduxReader(metadataDir,
		data.VideoEndedProperty,
		data.VideoTitleProperty,
		data.VideoOwnerChannelNameProperty)
	if err != nil {
		return cea.EndWithError(err)
	}

	videoIds := rdx.Keys(data.VideoEndedProperty)

	cea.TotalInt(len(videoIds))

	for _, videoId := range videoIds {
		if err := removeVideoFile(videoId, absVideosDir, rdx); err != nil {
			return cea.EndWithError(err)
		}
		if err := removePosters(videoId); err != nil {
			return cea.EndWithError(err)
		}
		cea.Increment()
	}

	cea.EndWithResult("done")

	return nil
}

func removeVideoFile(videoId, absVideosDir string, rdx kvas.ReadableRedux) error {
	title, _ := rdx.GetFirstVal(data.VideoTitleProperty, videoId)
	channel, _ := rdx.GetFirstVal(data.VideoOwnerChannelNameProperty, videoId)

	if title == "" || channel == "" {
		return nil
	}

	relVideoFilename := yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
	absVideoFilename := filepath.Join(absVideosDir, relVideoFilename)

	if _, err := os.Stat(absVideoFilename); err == nil {
		if err = os.Remove(absVideoFilename); err != nil {
			return err
		}
	}

	return nil
}

func removePosters(videoId string) error {

	for _, tq := range yt_urls.AllThumbnailQualities() {
		if app, err := paths.AbsPosterPath(videoId, tq); err == nil {
			if _, err := os.Stat(app); err == nil {
				if err = os.Remove(app); err != nil {
					return err
				}
			}
		} else {
			return err

		}
	}
	return nil
}
