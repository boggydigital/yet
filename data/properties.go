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
	VideoWatchedProperty           = "video-watched"
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
		VideoWatchedProperty,
	}
}
