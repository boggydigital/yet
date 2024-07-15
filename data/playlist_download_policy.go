package data

type DownloadPolicy string

const (
	Recent                DownloadPolicy = "recent"
	All                   DownloadPolicy = "all"
	DefaultDownloadPolicy                = Recent
)

const (
	RecentDownloadsLimit = 10
)

func ParseDownloadPolicy(policy string) DownloadPolicy {
	switch policy {
	case string(Recent):
		return Recent
	case string(All):
		return All
	default:
		return DefaultDownloadPolicy
	}
}

func AllDownloadPolicies() []DownloadPolicy {
	return []DownloadPolicy{
		Recent,
		All,
	}
}
