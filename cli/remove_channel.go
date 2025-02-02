package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func RemoveChannelHandler(u *url.URL) error {
	q := u.Query()

	channelId := q.Get("channel-id")
	options := &ChannelOptions{
		AutoRefresh:  q.Has("auto-refresh"),
		AutoDownload: q.Has("auto-download"),
		Expand:       q.Has("expand"),
		Force:        q.Has("force"),
	}

	return RemoveChannel(nil, channelId, options)
}

func RemoveChannel(rdx redux.Writeable, channelId string, opt *ChannelOptions) error {

	rpa := nod.Begin("removing channel %s...", channelId)
	defer rpa.End()

	if opt == nil {
		opt = DefaultChannelOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.PlaylistProperties()...)
	if err != nil {
		return rpa.EndWithError(err)
	}

	propertyKeys := make(map[string]string)

	if opt.AutoRefresh {
		propertyKeys[data.ChannelAutoRefreshProperty] = channelId
	}
	if opt.AutoDownload {
		propertyKeys[data.ChannelAutoDownloadProperty] = channelId
	}
	if opt.DownloadPolicy != data.DefaultDownloadPolicy {
		propertyKeys[data.ChannelDownloadPolicyProperty] = channelId
	}
	if opt.Expand {
		propertyKeys[data.ChannelExpandProperty] = channelId
	}

	for property, key := range propertyKeys {
		if err := rdx.CutKeys(property, key); err != nil {
			return rpa.EndWithError(err)
		}
	}

	rpa.EndWithResult("done")

	return nil
}
