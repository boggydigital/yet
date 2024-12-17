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
	VideoProgressProperty          = "video-progress"
	VideoDurationProperty          = "video-duration"
	VideoEndedDateProperty         = "video-ended-date"
	VideoEndedReasonProperty       = "video-ended-reason"
	VideoErrorsProperty            = "video-errors"
	VideoFavoriteProperty          = "video-favorite"

	VideoDehydratedPosterProperty       = "video-dehydrated-poster"
	VideoDehydratedRepColorProperty     = "video-dehydrated-rep-color"
	VideoDehydratedInputMissingProperty = "video-dehydrated-input-missing"

	VideoForcedDownloadProperty = "video-forced-download"

	VideoDownloadQueuedProperty    = "video-download-queued"
	VideoDownloadStartedProperty   = "video-download-started"
	VideoDownloadCompletedProperty = "video-download-completed"
	VideoDownloadCleanedUpProperty = "video-download-cleaned-up"

	VideoCaptionsLanguagesProperty = "video-captions-languages"
	VideoCaptionsKindsProperty     = "video-captions-kinds"
	VideoCaptionsNamesProperty     = "video-captions-names"

	PlaylistTitleProperty   = "playlist-title"
	PlaylistChannelProperty = "playlist-channel"
	PlaylistVideosProperty  = "playlist-videos"

	PlaylistAutoRefreshProperty    = "playlist-auto-refresh"
	PlaylistAutoDownloadProperty   = "playlist-auto-download"
	PlaylistDownloadPolicyProperty = "playlist-download-policy"
	PlaylistExpandProperty         = "playlist-expand"

	ChannelTitleProperty       = "channel-title"
	ChannelDescriptionProperty = "channel-description"
	ChannelVideosProperty      = "channel-videos"
	ChannelPlaylistsProperty   = "channel-playlists"

	ChannelAutoRefreshProperty    = "channel-auto-refresh"
	ChannelAutoDownloadProperty   = "channel-auto-download"
	ChannelDownloadPolicyProperty = "channel-download-policy"
	ChannelExpandProperty         = "channel-expand"

	YtDlpLatestDownloadedVersionProperty = "yt-dlp-latest-downloaded-version"
)

func VideoProperties() []string {
	return []string{
		VideoTitleProperty,
		VideoThumbnailUrlsProperty,
		VideoDehydratedPosterProperty,
		VideoDehydratedRepColorProperty,
		VideoDehydratedInputMissingProperty,
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
		VideoProgressProperty,
		VideoDurationProperty,
		VideoEndedDateProperty,
		VideoEndedReasonProperty,
		VideoErrorsProperty,
		VideoFavoriteProperty,
		VideoForcedDownloadProperty,
		VideoDownloadQueuedProperty,
		VideoDownloadStartedProperty,
		VideoDownloadCompletedProperty,
		VideoDownloadCleanedUpProperty,
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
		PlaylistAutoRefreshProperty,
		PlaylistAutoDownloadProperty,
		PlaylistDownloadPolicyProperty,
		PlaylistExpandProperty,
	}
}

func ChannelProperties() []string {
	return []string{
		ChannelTitleProperty,
		ChannelDescriptionProperty,
		ChannelVideosProperty,
		ChannelPlaylistsProperty,
		ChannelAutoRefreshProperty,
		ChannelAutoDownloadProperty,
		ChannelDownloadPolicyProperty,
		ChannelExpandProperty,
	}
}

func YtDlpProperties() []string {
	return []string{
		YtDlpLatestDownloadedVersionProperty,
	}
}

func AllProperties() []string {
	properties := make([]string, 0)
	properties = append(properties, VideoProperties()...)
	properties = append(properties, PlaylistProperties()...)
	properties = append(properties, ChannelProperties()...)
	properties = append(properties, YtDlpProperties()...)
	return properties
}
