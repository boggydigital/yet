package yeti

import (
	"fmt"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"path/filepath"
	"strings"
)

func DefaultFilenameDelegate(videoId string, videoPage *youtube_urls.InitialPlayerResponse) string {
	channel, title := "", ""
	if videoPage != nil {
		title = videoPage.VideoDetails.Title
		channel = videoPage.VideoDetails.Author
	}

	return ChannelTitleVideoIdFilename(channel, title, videoId)
}

// ChannelTitleVideoIdFilename constructs a filename based on video-id and
// optional channel and video title.
// If the channel or video title are available, the filename would be
// "channel/title-video-id.mp4". If the channel, title are not available, the filename would be
// "video-id.mp4". In either case, the resulting filename is sanitized to remove
// characters not suitable for file names.
func ChannelTitleVideoIdFilename(channel, title, videoId string) string {

	if strings.HasSuffix(videoId, youtube_urls.DefaultVideoExt) {
		return videoId
	}

	var fn string
	if title != "" {
		fn = fmt.Sprintf("%s-%s", title, videoId)
	} else {
		fn = fmt.Sprintf("%s", videoId)
	}

	// channel, video titles might contain characters that would be problematic for
	// modern operating system filesystems - removing those
	for _, ch := range []string{"/", ":", "?", "*", "<", ">", "\\", "|", "\"", "\n"} {
		fn = strings.ReplaceAll(fn, ch, "")
		channel = strings.ReplaceAll(channel, ch, "")
	}

	if channel != "" {
		fn = filepath.Join(channel, fn)
	}

	//while unlikely, it's possible for videos to be titled like
	//relative file paths (e.g. "../../title"), cleaning that up
	fn = filepath.Clean(fn)

	return fn + youtube_urls.DefaultVideoExt
}
