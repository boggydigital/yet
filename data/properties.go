package data

const (
	VideoTitleProperty             = "video-title"
	VideoThumbnailUrlsProperty     = "video-thumbnail-urls"
	VideoExternalChannelIdProperty = "video-external-channelid"
	VideoShortDescriptionProperty  = "video-short-description"
	VideoViewCountProperty         = "video-view-count"
	VideoKeywordsProperty          = "video-keywords"
	VideoOwnerChannelNameProperty  = "video-owner-channel-name"
	VideoOwnerProfileUrlProperty   = "video-owner-profile-url"
	VideoCategoryProperty          = "video-category"
	VideoPublishDateProperty       = "video-publish-date"
	VideoPublishTimeTextProperty   = "video-publish-time-text"
	VideoUploadDateProperty        = "video-upload-date"
	VideoDownloadedDateProperty    = "video-downloaded-date"
	VideoProgressProperty          = "video-progress"
	VideoDurationProperty          = "video-duration"
	VideoEndedProperty             = "video-ended"
	VideoSkippedProperty           = "video-skipped"
	VideoErrorsProperty            = "video-errors"
	VideoFavoriteProperty          = "video-favorite"

	VideoForcedDownloadProperty       = "video-forced-download"
	VideoSingleFormatDownloadProperty = "video-single-format-download"

	VideoDownloadQueuedProperty    = "video-download-queued"
	VideoDownloadStartedProperty   = "video-download-started"
	VideoDownloadCompletedProperty = "video-download-completed"

	VideoWatchlistQueuedProperty   = "video-watchlist-queued"
	VideoWatchlistDequeuedProperty = "video-watchlist-dequeued"

	VideoCaptionsLanguagesProperty = "video-captions-languages"
	VideoCaptionsKindsProperty     = "video-captions-kinds"
	VideoCaptionsNamesProperty     = "video-captions-names"

	PlaylistChannelProperty = "playlist-channel"
	PlaylistTitleProperty   = "playlist-title"
	PlaylistVideosProperty  = "playlist-videos"

	PlaylistRefreshProperty              = "playlist-refresh"
	PlaylistDownloadProperty             = "playlist-download"
	PlaylistSingleFormatDownloadProperty = "playlist-single-format"

	PlaylistQueuedDownloadsProperty = "playlist-queued-downloads"

	ChannelTitleProperty       = "channel-title"
	ChannelDescriptionProperty = "channel-description"
	ChannelPlaylistsProperty   = "channel-playlists"
)

func VideoProperties() []string {
	return []string{
		VideoTitleProperty,
		VideoThumbnailUrlsProperty,
		VideoExternalChannelIdProperty,
		VideoShortDescriptionProperty,
		VideoViewCountProperty,
		VideoKeywordsProperty,
		VideoOwnerChannelNameProperty,
		VideoOwnerProfileUrlProperty,
		VideoCategoryProperty,
		VideoPublishDateProperty,
		VideoPublishTimeTextProperty,
		VideoUploadDateProperty,
		VideoDownloadedDateProperty,
		VideoProgressProperty,
		VideoDurationProperty,
		VideoEndedProperty,
		VideoSkippedProperty,
		VideoErrorsProperty,
		VideoFavoriteProperty,
		VideoForcedDownloadProperty,
		VideoSingleFormatDownloadProperty,
		VideoDownloadQueuedProperty,
		VideoDownloadStartedProperty,
		VideoDownloadCompletedProperty,
		VideoWatchlistQueuedProperty,
		VideoWatchlistDequeuedProperty,
		VideoCaptionsLanguagesProperty,
		VideoCaptionsKindsProperty,
		VideoCaptionsNamesProperty,
	}
}

func PlaylistProperties() []string {
	return []string{
		PlaylistChannelProperty,
		PlaylistTitleProperty,
		PlaylistVideosProperty,
		PlaylistRefreshProperty,
		PlaylistDownloadProperty,
		PlaylistSingleFormatDownloadProperty,
		PlaylistQueuedDownloadsProperty,
	}
}

func ChannelProperties() []string {
	return []string{
		ChannelTitleProperty,
		ChannelDescriptionProperty,
		ChannelPlaylistsProperty,
	}
}

func AllProperties() []string {
	properties := make([]string, 0)
	properties = append(properties, VideoProperties()...)
	properties = append(properties, PlaylistProperties()...)
	properties = append(properties, ChannelProperties()...)
	return properties
}
