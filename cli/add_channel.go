package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"net/url"
)

func AddChannelHandler(u *url.URL) error {
	q := u.Query()

	channelId := q.Get("channel-id")
	options := &ChannelOptions{
		AutoRefresh:    q.Has("auto-refresh"),
		AutoDownload:   q.Has("auto-download"),
		DownloadPolicy: data.ParseDownloadPolicy(q.Get("download-policy")),
		Expand:         q.Has("expand"),
		Force:          q.Has("force"),
	}

	return AddChannel(nil, channelId, options)
}

func AddChannel(rdx redux.Writeable, channelId string, opt *ChannelOptions) error {

	aca := nod.Begin("adding channel %s...", channelId)
	defer aca.Done()

	if opt == nil {
		opt = DefaultChannelOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return err
	}

	propertyValues := make(map[string]map[string][]string)

	if opt.AutoRefresh {
		propertyValues[data.ChannelAutoRefreshProperty] = map[string][]string{
			channelId: {data.TrueValue},
		}
	}
	if opt.AutoDownload {
		propertyValues[data.ChannelAutoDownloadProperty] = map[string][]string{
			channelId: {data.TrueValue},
		}
	}
	if opt.DownloadPolicy != data.DefaultDownloadPolicy {
		propertyValues[data.ChannelDownloadPolicyProperty] = map[string][]string{
			channelId: {string(opt.DownloadPolicy)},
		}
	}
	if opt.Expand {
		propertyValues[data.ChannelExpandProperty] = map[string][]string{
			channelId: {data.TrueValue},
		}
	}

	for property, idValues := range propertyValues {
		if err := rdx.BatchAddValues(property, idValues); err != nil {
			return err
		}
	}

	// TODO: add GetChannelMetadata when ready
	//if err := GetPlaylistMetadata(rdx, opt, playlistId); err != nil {
	//	return err
	//}

	return nil
}
