package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func GetPlaylistsMetadataHandler(u *url.URL) error {
	q := u.Query()
	playlistIds := strings.Split(q.Get("playlist-id"), ",")
	options := &PlaylistOptions{
		Expand: q.Has("expand"),
		Force:  q.Has("force"),
	}
	return GetPlaylistsMetadata(nil, options, playlistIds...)
}

func GetPlaylistsMetadata(rdx kevlar.WriteableRedux, opt *PlaylistOptions, playlistIds ...string) error {
	gpma := nod.NewProgress("getting playlist metadata...")
	defer gpma.End()

	if opt == nil {
		opt = DefaultPlaylistOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return gpma.EndWithError(err)
	}

	parsedPlaylistIds, err := yeti.ParsePlaylistIds(playlistIds...)
	if err != nil {
		return gpma.EndWithError(err)
	}

	gpma.TotalInt(len(parsedPlaylistIds))

	for _, playlistId := range parsedPlaylistIds {

		if rdx.HasKey(data.PlaylistTitleProperty, playlistId) && !opt.Force {
			continue
		}

		expand := opt.Expand
		if re, ok := rdx.GetLastVal(data.PlaylistExpandProperty, playlistId); ok && re == data.TrueValue {
			expand = true
		}

		if err := yeti.GetPlaylistMetadata(nil, playlistId, expand, rdx); err != nil {
			gpma.Error(err)
		}

		gpma.Increment()
	}

	gpma.EndWithResult("done")

	return nil
}
