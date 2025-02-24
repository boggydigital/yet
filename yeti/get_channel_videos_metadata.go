package yeti

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
)

func GetChannelVideosMetadata(channelVideosPage *youtube_urls.ChannelVideosInitialData, channelId string, expand bool, rdx redux.Writeable) error {

	gcvma := nod.Begin(" channel videos metadata for %s", channelId)
	defer gcvma.Done()

	if err := rdx.MustHave(
		data.ChannelTitleProperty,
		data.ChannelVideosProperty,
		data.ChannelDescriptionProperty,
		data.VideoTitleProperty,
		data.VideoDurationProperty,
		data.VideoOwnerChannelNameProperty); err != nil {
		return err
	}

	var err error
	if channelVideosPage == nil {
		channelVideosPage, err = youtube_urls.GetChannelVideosPage(http.DefaultClient, channelId)
		if err != nil {
			return err
		}
	}

	if channelTitle := channelVideosPage.Metadata.ChannelMetadataRenderer.Title; channelTitle != "" {
		if err = rdx.ReplaceValues(data.ChannelTitleProperty, channelId, channelTitle); err != nil {
			return err
		}
	}

	if description := channelVideosPage.Metadata.ChannelMetadataRenderer.Description; description != "" {
		if err = rdx.ReplaceValues(data.ChannelDescriptionProperty, channelId, description); err != nil {
			return err
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
				return err
			}
		} else {
			channelVideosPage = nil
		}
	}

	if err = rdx.ReplaceValues(data.ChannelVideosProperty, channelId, channelVideos...); err != nil {
		return err
	}

	if err = rdx.BatchReplaceValues(data.VideoTitleProperty, videoTitles); err != nil {
		return err
	}

	if err = rdx.BatchReplaceValues(data.VideoDurationProperty, videoLengthsSeconds); err != nil {
		return err
	}

	if err = rdx.BatchReplaceValues(data.VideoOwnerChannelNameProperty, videoOwnerChannels); err != nil {
		return err
	}

	return nil
}
