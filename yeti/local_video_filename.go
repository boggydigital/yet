package yeti

import (
	"errors"
	"fmt"
	"github.com/boggydigital/busan"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"os"
	"path/filepath"
	"strings"
)

const (
	globTemplate = "*/*-{video-id}" + youtube_urls.DefaultVideoExt
)

// RelLocalVideoFilename constructs a filename based on video-id and
// optional channel and video title.
// If the channel or video title are available, the filename would be
// "channel/title-video-id.mp4". If the channel, title are not available,
// the filename would be "video-id.mp4". In either case, the resulting
// filename is sanitized to remove characters not suitable for file names.
func RelLocalVideoFilename(channel, title, videoId string) string {

	// channel, video titles might contain characters that would be problematic for
	// modern operating system filesystems - removing those
	channel = busan.Sanitize(channel)
	title = busan.Sanitize(title)

	fn := videoId
	if title != "" {
		fn = fmt.Sprintf("%s-%s", title, videoId)
	}
	if channel != "" {
		fn = filepath.Join(channel, fn)
	}

	return fn + youtube_urls.DefaultVideoExt
}

// LocateLocalVideo looks for local video files for a given video-id
// previously we'd use RelLocalVideoFilename func above and check that name,
// however as it turns out - it's not uncommon for videos to change the title
// which leads to existing local files seemingly missing - as we've downloaded
// them under different title previously. Using LocateLocalVideo should
// mitigate this problem.
func LocateLocalVideo(videoId string) (string, error) {

	videosDir, err := pathways.GetAbsDir(data.Videos)
	if err != nil {
		return "", err
	}

	pattern := strings.Replace(globTemplate, "{video-id}", videoId, 1)
	pattern = filepath.Join(videosDir, pattern)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 1 {
		return matches[0], nil
	} else if len(matches) == 0 {
		return "", os.ErrNotExist
	} else {
		return "", errors.New("several local files match video-id " + videoId)
	}
}
