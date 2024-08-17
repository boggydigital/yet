package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"golang.org/x/exp/slices"
	"net/url"
)

var preserveVideoProperties = []string{
	data.VideoTitleProperty,             // required for history
	data.VideoOwnerChannelNameProperty,  // required for cleanup (checking empty videos directories)
	data.VideoEndedDateProperty,         // required for cleanup
	data.VideoEndedReasonProperty,       // required for history
	data.VideoFavoriteProperty,          // required for cleanup
	data.VideoDownloadCompletedProperty, // required for cleanup
	data.VideoDownloadCleanedUpProperty, // required for cleanup
}

func ScrubEndedPropertiesHandler(u *url.URL) error {
	return ScrubEndedProperties(nil)
}

func ScrubEndedProperties(rdx kevlar.WriteableRedux) error {

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
