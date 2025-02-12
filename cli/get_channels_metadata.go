package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func GetChannelsMetadataHandler(u *url.URL) error {
	q := u.Query()
	channelIds := strings.Split(q.Get("channel-id"), ",")
	options := &ChannelOptions{
		Playlists: q.Has("playlists"),
		Expand:    q.Has("expand"),
		Force:     q.Has("force"),
	}
	return GetChannelsMetadata(nil, options, channelIds...)
}

func GetChannelsMetadata(rdx redux.Writeable, opt *ChannelOptions, channelIds ...string) error {
	gchma := nod.NewProgress("getting channel metadata...")
	defer gchma.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return err
	}

	gchma.TotalInt(len(channelIds))

	for _, channelId := range channelIds {

		if rdx.HasKey(data.ChannelTitleProperty, channelId) && !opt.Force {
			continue
		}

		expand := opt.Expand
		if ce, ok := rdx.GetLastVal(data.ChannelExpandProperty, channelId); ok && ce == data.TrueValue {
			expand = true
		}

		if err := yeti.GetChannelVideosMetadata(nil, channelId, expand, rdx); err != nil {
			gchma.Error(err)
		}

		if opt.Playlists {
			if err := yeti.GetChannelPlaylistsMetadata(nil, channelId, rdx); err != nil {
				gchma.Error(err)
			}
		}

		gchma.Increment()
	}

	gchma.EndWithResult("done")

	return nil
}
