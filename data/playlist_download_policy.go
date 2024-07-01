package data

type PlaylistDownloadPolicy string

const (
	Recent                PlaylistDownloadPolicy = "recent"
	All                   PlaylistDownloadPolicy = "all"
	DefaultDownloadPolicy                        = Recent
)

const (
	RecentDownloadsLimit = 10
)

func ParsePlaylistDownloadPolicy(policy string) PlaylistDownloadPolicy {
	switch policy {
	case string(Recent):
		return Recent
	case string(All):
		return All
	default:
		return DefaultDownloadPolicy
	}
}

func AllPlaylistDownloadPolicies() []PlaylistDownloadPolicy {
	return []PlaylistDownloadPolicy{
		Recent,
		All,
	}
}
