package data

type PlaylistDownloadPolicy string

const (
	Unset  PlaylistDownloadPolicy = "unset"
	Recent PlaylistDownloadPolicy = "recent"
	All    PlaylistDownloadPolicy = "all"
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
		return Unset
	}
}
