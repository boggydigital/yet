package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func GetChannelMetadataHandler(u *url.URL) error {
	q := u.Query()
	ids := strings.Split(q.Get("id"), ",")
	force := q.Has("force")
	return GetChannelMetadata(force, ids...)
}

func GetChannelMetadata(force bool, ids ...string) error {
	gchma := nod.NewProgress("getting channel metadata...")
	defer gchma.End()

	gchma.TotalInt(len(ids))

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return gchma.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return gchma.EndWithError(err)
	}

	for _, channelId := range ids {

		if rdx.HasKey(data.ChannelTitleProperty, channelId) && !force {
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
