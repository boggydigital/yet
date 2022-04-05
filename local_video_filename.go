package main

import (
	"fmt"
	"github.com/boggydigital/yt_urls"
	"path/filepath"
	"strings"
)

const mp4Ext = ".mp4"

//localVideoFilename constructs a filename based on video-id and
//optional video title. If the title is available, the filename would be
//"title-video-id.mp4". If the title is not available, the filename would be
//"video-id.mp4". In either case, the resulting filename is sanitized to remove
//characters not suitable for file names.
func filenameDelegate(videoId string, videoPage *yt_urls.InitialPlayerResponse) string {

	title := ""
	if videoPage != nil {
		title = videoPage.Title()
	}

	return titleVideoIdFilename(title, videoId)

}

func titleVideoIdFilename(title, videoId string) string {
	var fn string
	if title != "" {
		fn = fmt.Sprintf("%s-%s", title, videoId)
	} else {
		fn = fmt.Sprintf("%s", videoId)
	}

	//while unlikely, it's possible for videos to be titled like
	//relative file paths (e.g. "../../title"), cleaning that up
	fn = filepath.Clean(fn)

	//video titles commonly contain characters that would be problematic for
	//modern operating system filesystems - removing those
	for _, ch := range []string{"/", ":", "?", "*", "<", ">", "\\", "|"} {
		fn = strings.ReplaceAll(fn, ch, "")
	}

	return fn + yt_urls.DefaultExt
}
