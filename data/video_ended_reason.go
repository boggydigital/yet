package data

import (
	"maps"
	"slices"
)

type VideoEndedReason int

const (
	Completed VideoEndedReason = iota
	Skipped
	SeenEnough

	DefaultEndedReason = Completed
)

var videoEndedReasonNames = map[VideoEndedReason]string{
	Completed:  "completed",
	Skipped:    "skipped",
	SeenEnough: "seen-enough",
}

func ParseVideoEndedReason(s string) VideoEndedReason {
	for ver, name := range videoEndedReasonNames {
		if s == name {
			return ver
		}
	}
	return DefaultEndedReason
}

func AllVideoEndedReasons() []VideoEndedReason {
	return slices.Collect(maps.Keys(videoEndedReasonNames))
}

func (ver VideoEndedReason) String() string {
	return videoEndedReasonNames[ver]
}
