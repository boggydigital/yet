package yeti

import (
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yt_urls"
)

func ExtractMetadata(ipr *yt_urls.InitialPlayerResponse) map[string][]string {
	vpm := make(map[string][]string)

	vpm[data.VideoTitleProperty] = []string{ipr.VideoDetails.Title}
	vpm[data.VideoThumbnailUrlsProperty] = make([]string, 0, len(ipr.VideoDetails.Thumbnail.Thumbnails))
	for _, t := range ipr.VideoDetails.Thumbnail.Thumbnails {
		vpm[data.VideoThumbnailUrlsProperty] = append(vpm[data.VideoThumbnailUrlsProperty], t.Url)
	}
	vpm[data.VideoDurationProperty] = []string{ipr.VideoDetails.LengthSeconds}
	vpm[data.VideoExternalChannelIdProperty] = []string{ipr.VideoDetails.ChannelId}
	vpm[data.VideoShortDescriptionProperty] = []string{ipr.VideoDetails.ShortDescription}
	vpm[data.VideoViewCountProperty] = []string{ipr.VideoDetails.ViewCount}
	vpm[data.VideoKeywordsProperty] = ipr.VideoDetails.Keywords

	vpm[data.VideoOwnerChannelNameProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.OwnerChannelName}
	vpm[data.VideoOwnerProfileUrlProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.OwnerProfileUrl}
	vpm[data.VideoCategoryProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.Category}
	vpm[data.VideoPublishDateProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.PublishDate}
	vpm[data.VideoUploadDateProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.UploadDate}

	return vpm
}
