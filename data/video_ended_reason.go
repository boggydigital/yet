package data

type VideoEndedReason string

const (
	Completed          VideoEndedReason = "completed"
	Skipped            VideoEndedReason = "skipped"
	SeenEnough         VideoEndedReason = "seen-enough"
	DefaultEndedReason                  = Completed
)

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
