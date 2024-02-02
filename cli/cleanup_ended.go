package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func CleanupEndedHandler(u *url.URL) error {
	return CleanupEnded()
}

func CleanupEnded() error {

	cea := nod.NewProgress("cleaning up ended media...")
	defer cea.End()

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		return cea.EndWithError(err)
	}

	absVideosDir, err := pasu.GetAbsDir(paths.Videos)
	if err != nil {
		return cea.EndWithError(err)
	}

	rdx, err := kvas.NewReduxReader(metadataDir,
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

		// checking and removing empty directories
		if channelTitle, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && channelTitle != "" {
			absDirName := filepath.Join(absVideosDir, channelTitle)
			if ok, err := directoryIsEmpty(absDirName); ok && err == nil {
				if err := os.Remove(absDirName); err != nil {
					return err
				}
			}
		}
	}

	cea.EndWithResult("done")

	return nil
}

// https://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
func directoryIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func removeVideoFile(videoId, absVideosDir string, rdx kvas.ReadableRedux) error {
	title, _ := rdx.GetLastVal(data.VideoTitleProperty, videoId)
	channel, _ := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId)

	if title == "" || channel == "" {
		return nil
	}

	relVideoFilename := ""

	if strings.HasSuffix(videoId, yt_urls.DefaultVideoExt) {
		relVideoFilename = videoId
	} else {
		relVideoFilename = yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
	}

	if relVideoFilename == "" {
		return nil
	}

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
