package yeti

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
)

func GetChannelVideosMetadata(channelVideosPage *youtube_urls.ChannelVideosInitialData, channelId string, expand bool, rdx kevlar.WriteableRedux) error {

	gcvma := nod.Begin(" metadata for %s", channelId)
	defer gcvma.End()

	if err := rdx.MustHave(
		data.ChannelTitleProperty,
		data.ChannelVideosProperty,
		data.VideoTitleProperty,
		data.VideoDurationProperty,
		data.VideoOwnerChannelNameProperty); err != nil {
		return gcvma.EndWithError(err)
	}

	var err error
	if channelVideosPage == nil {
		channelVideosPage, err = youtube_urls.GetChannelVideosPage(http.DefaultClient, channelId)
		if err != nil {
			return gcvma.EndWithError(err)
		}
	}

	if channelTitle := channelVideosPage.Metadata.ChannelMetadataRenderer.Title; channelTitle != "" {
		if err := rdx.AddValues(data.ChannelTitleProperty, channelId, channelTitle); err != nil {
			return gcvma.EndWithError(err)
		}
	}

	channelVideos := make([]string, 0)
	videoTitles := make(map[string][]string)
	videoLengthsSeconds := make(map[string][]string)
	videoOwnerChannels := make(map[string][]string)

	for channelVideosPage != nil &&
		len(channelVideosPage.Videos()) > 0 {

		for _, video := range channelVideosPage.Videos() {
			videoId := video.VideoId
			channelVideos = append(channelVideos, videoId)
			videoTitles[videoId] = []string{video.Title}
			videoOwnerChannels[videoId] = []string{channelVideosPage.Metadata.ChannelMetadataRenderer.Title}
			videoLengthsSeconds[videoId] = []string{video.LengthSeconds}
		}

		if expand && channelVideosPage.HasContinuation() {
			if err = channelVideosPage.Continue(http.DefaultClient); err != nil {
				return gcvma.EndWithError(err)
			}
		} else {
			channelVideosPage = nil
		}
	}

	if err := rdx.ReplaceValues(data.ChannelVideosProperty, channelId, channelVideos...); err != nil {
		return gcvma.EndWithError(err)
	}

	if err := rdx.BatchReplaceValues(data.VideoTitleProperty, videoTitles); err != nil {
		return gcvma.EndWithError(err)
	}

	if err := rdx.BatchReplaceValues(data.VideoDurationProperty, videoLengthsSeconds); err != nil {
		return gcvma.EndWithError(err)
	}

	if err := rdx.BatchReplaceValues(data.VideoOwnerChannelNameProperty, videoOwnerChannels); err != nil {
		return gcvma.EndWithError(err)
	}

	gcvma.EndWithResult("done")

	return nil
}
