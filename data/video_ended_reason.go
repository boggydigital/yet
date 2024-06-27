package data

type VideoEndedReason string

const (
	Unspecified VideoEndedReason = "unspecified"
	Skipped     VideoEndedReason = "skipped"
	SeenEnough  VideoEndedReason = "seen-enough"
)

func ParseVideoEndedReason(s string) VideoEndedReason {
	switch s {
	case string(Skipped):
		return Skipped
	case string(SeenEnough):
		return SeenEnough
	default:
		return Unspecified
	}
}
