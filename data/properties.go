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
	VideoUploadDateProperty        = "video-upload-date"
	VideoDownloadedDateProperty    = "video-downloaded-date"
	VideoProgressProperty          = "video-progress"
	VideoEndedProperty             = "video-ended"
	VideoErrorsProperty            = "video-errors"

	VideosDownloadQueueProperty = "videos-download-queue"
	VideosWatchlistProperty     = "videos-watchlist"

	VideoCaptionsLanguages = "video-captions-languages"
	VideoCaptionsKinds     = "video-captions-kinds"
	VideoCaptionsNames     = "video-captions-names"
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
		VideoUploadDateProperty,
		VideoDownloadedDateProperty,
		VideoProgressProperty,
		VideoEndedProperty,
		VideoErrorsProperty,
		VideosDownloadQueueProperty,
		VideosWatchlistProperty,
		VideoCaptionsLanguages,
		VideoCaptionsKinds,
		VideoCaptionsNames,
	}
}
