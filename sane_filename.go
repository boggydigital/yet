package main

import (
	"fmt"
	"strings"
)

func saneFilename(title, videoId string) string {
	unsafeChars := []string{"/", ":", "?", "*"}

	fn := fmt.Sprintf("%s-%s", title, videoId)
	if title == "" {
		fn = fmt.Sprintf("%s", videoId)
	}

	for _, ch := range unsafeChars {
		fn = strings.ReplaceAll(fn, ch, "")
	}

	fn += ".mp4"

	return fn
}
