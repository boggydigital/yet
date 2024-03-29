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

	VideosDownloadQueueProperty       = "videos-download-queue"
	VideoForcedDownloadProperty       = "video-forced-download"
	VideoSingleFormatDownloadProperty = "video-single-format-download"
	VideosWatchlistProperty           = "videos-watchlist"

	VideoCaptionsLanguagesProperty = "video-captions-languages"
	VideoCaptionsKindsProperty     = "video-captions-kinds"
	VideoCaptionsNamesProperty     = "video-captions-names"

	PlaylistWatchlistProperty            = "playlist-watchlist"
	PlaylistDownloadQueueProperty        = "playlist-download-queue"
	PlaylistNewVideosProperty            = "playlist-new-videos"
	PlaylistChannelProperty              = "playlist-channel"
	PlaylistTitleProperty                = "playlist-title"
	PlaylistVideosProperty               = "playlist-videos"
	PlaylistSingleFormatDownloadProperty = "playlist-single-format"

	ChannelTitleProperty       = "channel-title"
	ChannelDescriptionProperty = "channel-description"
	ChannelPlaylistsProperty   = "channel-playlists"
)

func AllProperties() []string {
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
		VideoForcedDownloadProperty,
		VideoSingleFormatDownloadProperty,
		VideoProgressProperty,
		VideoDurationProperty,
		VideoEndedProperty,
		VideoSkippedProperty,
		VideoErrorsProperty,
		VideosDownloadQueueProperty,
		VideosWatchlistProperty,
		VideoCaptionsLanguagesProperty,
		VideoCaptionsKindsProperty,
		VideoCaptionsNamesProperty,

		PlaylistWatchlistProperty,
		PlaylistDownloadQueueProperty,
		PlaylistSingleFormatDownloadProperty,
		PlaylistNewVideosProperty,
		PlaylistChannelProperty,
		PlaylistTitleProperty,
		PlaylistVideosProperty,

		ChannelTitleProperty,
		ChannelDescriptionProperty,
		ChannelPlaylistsProperty,
	}
}
