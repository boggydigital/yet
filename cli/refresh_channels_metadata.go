package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RefreshChannelsMetadataHandler(_ *url.URL) error {
	return RefreshChannelsMetadata(nil)
}

func RefreshChannelsMetadata(rdx redux.Writeable) error {

	ucma := nod.NewProgress("updating all channels metadata...")
	defer ucma.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.ChannelProperties()...)
	if err != nil {
		return ucma.EndWithError(err)
	}

	// update auto-refresh channels metadata
	ucma.TotalInt(rdx.Len(data.ChannelAutoRefreshProperty))

	refreshOptions := &ChannelOptions{
		Force: true,
	}

	for channelId := range rdx.Keys(data.ChannelAutoRefreshProperty) {

		if err := GetChannelsMetadata(rdx, refreshOptions, channelId); err != nil {
			return ucma.EndWithError(err)
		}

		ucma.Increment()
	}

	ucma.EndWithResult("done")

	return nil
}
