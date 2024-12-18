package data

type VideoEndedReason string

const (
	Completed          VideoEndedReason = "completed"
	Skipped            VideoEndedReason = "skipped"
	SeenEnough         VideoEndedReason = "seen-enough"
	DefaultEndedReason                  = Completed
)

var videoEndedReasonStrings = map[VideoEndedReason]string{
	Completed:  "Completed",
	Skipped:    "Skipped",
	SeenEnough: "Seen enough",
}

func ParseVideoEndedReason(s string) VideoEndedReason {
	switch s {
	case string(Completed):
		return Completed
	case string(Skipped):
		return Skipped
	case string(SeenEnough):
		return SeenEnough
	default:
		return DefaultEndedReason
	}
}

func AllVideoEndedReasons() []VideoEndedReason {
	return []VideoEndedReason{
		Completed,
		Skipped,
		SeenEnough,
	}
}

func (ver VideoEndedReason) String() string {
	return videoEndedReasonStrings[ver]
}
