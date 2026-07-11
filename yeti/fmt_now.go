package yeti

import "time"

const DefaultDelay = 24 * time.Hour

func FmtNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}
