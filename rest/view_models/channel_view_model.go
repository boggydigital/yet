package view_models

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
)

type ChannelViewModel struct {
	ChannelId    string
	ChannelTitle string
}

func GetChannelViewModel(channelId string, rdx kvas.ReadableRedux) *ChannelViewModel {
	channelTitle := channelId
	if ct, ok := rdx.GetFirstVal(data.ChannelTitleProperty, channelId); ok && ct != "" {
		channelTitle = ct
	}

	return &ChannelViewModel{
		ChannelId:    channelId,
		ChannelTitle: channelTitle,
	}
}
