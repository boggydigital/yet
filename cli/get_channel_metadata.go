package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func GetChannelMetadataHandler(u *url.URL) error {
	q := u.Query()
	channelIds := strings.Split(q.Get("channel-id"), ",")
	options := &ChannelOptions{Force: q.Has("force")}
	return GetChannelMetadata(nil, options, channelIds...)
}

func GetChannelMetadata(rdx kvas.WriteableRedux, opt *ChannelOptions, channelIds ...string) error {
	gchma := nod.NewProgress("getting channel metadata...")
	defer gchma.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return gchma.EndWithError(err)
	}

	gchma.TotalInt(len(channelIds))

	for _, channelId := range channelIds {

		if rdx.HasKey(data.ChannelTitleProperty, channelId) && !opt.Force {
			continue
		}

		if err := yeti.GetChannelPageMetadata(nil, channelId, rdx); err != nil {
			gchma.Error(err)
		}

		gchma.Increment()
	}

	gchma.EndWithResult("done")

	return nil
}
