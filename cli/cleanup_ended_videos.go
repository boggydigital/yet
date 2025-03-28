package cli

import (
	"github.com/boggydigital/busan"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func CleanupEndedVideosHandler(u *url.URL) error {
	now := u.Query().Has("now")
	return CleanupEndedVideos(now, nil)
}

// CleanupEndedVideos removes downloads for Ended videos that match the following conditions:
// - video download has not been downloaded earlier
// - at least 48 hours have passed since the ended date (unless no-delay was set)
// Additionally CleanupEndedVideos will remove video properties (except for title) for ended videos
func CleanupEndedVideos(now bool, rdx redux.Writeable) error {

	cea := nod.NewProgress("cleaning up ended media...")
	defer cea.Done()

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return err
	}

	absVideosDir, err := pathways.GetAbsDir(data.Videos)
	if err != nil {
		return err
	}

	cleanupVideoIds := make([]string, 0)

	for id := range rdx.Keys(data.VideoEndedDateProperty) {

		// don't cleanup favorite videos
		if rdx.HasKey(data.VideoFavoriteProperty, id) {
			continue
		}

		dcTime := ""
		if dct, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, id); ok && dct != "" {
			dcTime = dct
		}

		// skip video that have been cleaned up _after_ the latest download
		if dcut, ok := rdx.GetLastVal(data.VideoDownloadCleanedUpProperty, id); ok && dcut > dcTime {
			continue
		}

		if !now {
			if eds, ok := rdx.GetLastVal(data.VideoEndedDateProperty, id); ok {
				if ed, err := time.Parse(time.RFC3339, eds); err == nil {
					dur := time.Now().Sub(ed)
					if dur < yeti.DefaultDelay {
						continue
					}
				} else {
					return err
				}
			}
		}

		cleanupVideoIds = append(cleanupVideoIds, id)
	}

	cea.TotalInt(len(cleanupVideoIds))

	for _, videoId := range cleanupVideoIds {

		if err := removeVideoFile(videoId, rdx); err != nil {
			return err
		}
		if err := removePosters(videoId); err != nil {
			return err
		}

		if err := rdx.AddValues(data.VideoDownloadCleanedUpProperty, videoId, time.Now().Format(time.RFC3339)); err != nil {
			return err
		}

		// checking and removing empty directories
		if channelTitle, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && channelTitle != "" {
			channelTitle = busan.Sanitize(channelTitle)
			absDirName := filepath.Join(absVideosDir, channelTitle)
			if ok, err := directoryIsEmpty(absDirName); ok && err == nil {
				if err := os.Remove(absDirName); err != nil {
					return err
				}
			}
		}

		cea.Increment()
	}

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

func removeVideoFile(videoId string, rdx redux.Readable) error {
	title, _ := rdx.GetLastVal(data.VideoTitleProperty, videoId)
	channel, _ := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId)

	if title == "" || channel == "" {
		return nil
	}

	var absVideoFilename string

	if afv, err := yeti.LocateLocalVideo(videoId); os.IsNotExist(err) {
		// do nothing - file doesn't exist
	} else if err != nil {
		return err
	} else {
		absVideoFilename = afv
	}

	if absVideoFilename == "" {
		return nil
	}

	if _, err := os.Stat(absVideoFilename); err == nil {
		if err = os.Remove(absVideoFilename); err != nil {
			return err
		}
	}

	return nil
}

func removePosters(videoId string) error {

	for _, tq := range youtube_urls.AllThumbnailQualities() {
		if app, err := data.AbsPosterPath(videoId, tq); err == nil {
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
