package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"net/url"
	"slices"
)

var preserveVideoProperties = []string{
	data.VideoTitleProperty,             // required for history
	data.VideoOwnerChannelNameProperty,  // required for cleanup (checking empty videos directories)
	data.VideoEndedDateProperty,         // required for cleanup
	data.VideoEndedReasonProperty,       // required for history
	data.VideoFavoriteProperty,          // required for cleanup
	data.VideoDownloadCompletedProperty, // required for cleanup
	data.VideoDownloadCleanedUpProperty, // required for cleanup
	data.VideoExternalChannelIdProperty, // required for watch
}

func ScrubEndedPropertiesHandler(_ *url.URL) error {
	return ScrubEndedProperties(nil)
}

// ScrubEndedProperties will remove all non-preserved properties for ended videos.
// Preserved properties are required for core functionality - history, cleanup, etc.
func ScrubEndedProperties(rdx redux.Writeable) error {
	sevpa := nod.NewProgress("scrubbing ended videos properties...")
	defer sevpa.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return sevpa.EndWithError(err)
	}

	endedVideos := rdx.Keys(data.VideoEndedDateProperty)

	sevpa.TotalInt(len(endedVideos))

	for _, videoId := range endedVideos {

		for _, vp := range data.VideoProperties() {
			if slices.Contains(preserveVideoProperties, vp) {
				continue
			}

			if rdx.HasKey(vp, videoId) {
				if err := rdx.CutKeys(vp, videoId); err != nil {
					return sevpa.EndWithError(err)
				}
			}
		}

		sevpa.Increment()
	}

	sevpa.EndWithResult("done")

	return nil
}
