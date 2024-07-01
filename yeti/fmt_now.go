package yeti

import "time"

const DefaultDelay time.Duration = 48 * time.Hour

func FmtNow() string {
	return time.Now().Format(time.RFC3339)
}
