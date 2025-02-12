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
	defer sevpa.Done()

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return err
	}

	sevpa.TotalInt(rdx.Len(data.VideoEndedDateProperty))

	for videoId := range rdx.Keys(data.VideoEndedDateProperty) {

		for _, vp := range data.VideoProperties() {
			if slices.Contains(preserveVideoProperties, vp) {
				continue
			}

			if rdx.HasKey(vp, videoId) {
				if err := rdx.CutKeys(vp, videoId); err != nil {
					return err
				}
			}
		}

		sevpa.Increment()
	}

	return nil
}
