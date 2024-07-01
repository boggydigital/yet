package clo

import "github.com/boggydigital/yet/data"

func PlaylistDownloadPolicies() []string {
	return []string{
		string(data.Recent),
		string(data.All),
	}
}

func VideoEndedReasons() []string {
	return []string{
		string(data.Skipped),
		string(data.SeenEnough),
	}
}

func Values() map[string]func() []string {
	return map[string]func() []string{
		"playlist-download-policies": PlaylistDownloadPolicies,
		"video-ended-reasons":        VideoEndedReasons,
	}
}
