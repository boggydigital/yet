package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"net/url"
	"strings"
)

func GetPlaylistMetadataHandler(u *url.URL) error {
	q := u.Query()
	ids := strings.Split(q.Get("id"), ",")
	allVideos := q.Has("all-videos")
	force := q.Has("force")
	return GetPlaylistMetadata(nil, allVideos, force, ids...)
}

func GetPlaylistMetadata(rdx kvas.WriteableRedux, allVideos, force bool, ids ...string) error {
	gpma := nod.NewProgress("getting playlist metadata...")
	defer gpma.End()

	playlistIds, err := yeti.ParsePlaylistIds(ids...)
	if err != nil {
		return gpma.EndWithError(err)
	}

	gpma.TotalInt(len(playlistIds))

	if rdx == nil {

		metadataDir, err := pasu.GetAbsDir(paths.Metadata)
		if err != nil {
			return gpma.EndWithError(err)
		}

		rdx, err = kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
		if err != nil {
			return gpma.EndWithError(err)
		}
	}

	for _, playlistId := range playlistIds {

		if rdx.HasKey(data.PlaylistTitleProperty, playlistId) && !force {
			continue
		}

		if err := yeti.GetPlaylistPageMetadata(nil, playlistId, allVideos, rdx); err != nil {
			gpma.Error(err)
		}

		gpma.Increment()
	}

	gpma.EndWithResult("done")

	return nil
}
