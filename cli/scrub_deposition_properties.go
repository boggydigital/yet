package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"net/url"
	"slices"
)

func ScrubDepositionPropertiesHandler(u *url.URL) error {
	return ScrubDepositionProperties(nil)
}

// ScrubDepositionProperties will remove all accumulated property depositions:
// - search results
// - older channel and playlist videos properties
// To do that we start by identifying all critical videos:
// - part of current channel, playlist data
// - downloaded, not-ended videos
// Then we iterate over all non-preserved properties and remove data for all non-critical videos
func ScrubDepositionProperties(rdx kevlar.WriteableRedux) error {
	sdpa := nod.NewProgress("scrubbing deposition properties...")
	defer sdpa.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return sdpa.EndWithError(err)
	}

	currentVideoIds := make(map[string]any)

	for _, videoId := range rdx.Keys(data.VideoDownloadCompletedProperty) {
		if rdx.HasKey(data.VideoEndedDateProperty, videoId) {
			continue
		}
		currentVideoIds[videoId] = nil
	}

	for _, channelId := range rdx.Keys(data.ChannelAutoRefreshProperty) {
		if videos, ok := rdx.GetAllValues(data.ChannelVideosProperty, channelId); ok {
			for _, videoId := range videos {
				currentVideoIds[videoId] = nil
			}
		}
	}

	for _, playlistId := range rdx.Keys(data.PlaylistAutoRefreshProperty) {
		if videos, ok := rdx.GetAllValues(data.PlaylistVideosProperty, playlistId); ok {
			for _, videoId := range videos {
				currentVideoIds[videoId] = nil
			}
		}
	}

	properties := data.VideoProperties()

	sdpa.TotalInt(len(properties))

	for _, vp := range properties {
		if slices.Contains(preserveVideoProperties, vp) {
			sdpa.Increment()
			continue
		}
		videos := rdx.Keys(vp)
		for _, videoId := range videos {
			if _, ok := currentVideoIds[videoId]; !ok {
				if rdx.HasKey(vp, videoId) {
					if err := rdx.CutKeys(vp, videoId); err != nil {
						return sdpa.EndWithError(err)
					}
				}
			}
		}
		sdpa.Increment()
	}

	sdpa.EndWithResult("done")

	return nil
}
