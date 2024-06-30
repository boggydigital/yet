package cli

type VideoDownloadOptions struct {
	*VideoOptions
	Source string
}

func DefaultVideoDownloadOptions() *VideoDownloadOptions {
	return &VideoDownloadOptions{
		VideoOptions: &VideoOptions{
			PreferSingleFormat: true,
			Force:              false,
		},
		Source: "",
	}
}
