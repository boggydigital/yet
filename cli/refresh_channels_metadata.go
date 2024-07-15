package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RefreshChannelsMetadataHandler(_ *url.URL) error {
	return RefreshChannelsMetadata(nil)
}

func RefreshChannelsMetadata(rdx kevlar.WriteableRedux) error {

	ucma := nod.NewProgress("updating all channels metadata...")
	defer ucma.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.ChannelProperties()...)
	if err != nil {
		return ucma.EndWithError(err)
	}

	// update auto-refresh channels metadata
	channelIds := rdx.Keys(data.ChannelAutoRefreshProperty)
	ucma.TotalInt(len(channelIds))

	refreshOptions := &ChannelOptions{
		Force: true,
	}

	for _, channelId := range channelIds {

		if err := GetChannelsMetadata(rdx, refreshOptions, channelId); err != nil {
			return ucma.EndWithError(err)
		}

		ucma.Increment()
	}

	ucma.EndWithResult("done")

	return nil
}
