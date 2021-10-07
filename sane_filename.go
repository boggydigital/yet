package main

import (
	"fmt"
	"strings"
)

func saneFilename(title, videoId string) string {
	fn := fmt.Sprintf("%s-%s", title, videoId)
	if title == "" {
		fn = fmt.Sprintf("%s", videoId)
	}

	fn = strings.ReplaceAll(fn, "/", "")
	fn = strings.ReplaceAll(fn, ":", "")
	fn = strings.ReplaceAll(fn, ".", "")

	fn += ".mp4"

	return fn
}
