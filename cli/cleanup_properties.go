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

func CleanupPropertiesHandler(u *url.URL) error {
	return CleanupProperties(nil)
}

func CleanupProperties(rdx kevlar.WriteableRedux) error {

	cpa := nod.NewProgress("cleaning up video properties...")
	defer cpa.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return cpa.EndWithError(err)
	}

	endedVideos := rdx.Keys(data.VideoEndedDateProperty)

	cpa.TotalInt(len(endedVideos))

	for _, videoId := range endedVideos {

		for _, vp := range data.VideoProperties() {
			if slices.Contains(preserveVideoProperties, vp) {
				continue
			}

			if rdx.HasKey(vp, videoId) {
				if err := rdx.CutKeys(vp, videoId); err != nil {
					return cpa.EndWithError(err)
				}
			}
		}

		cpa.Increment()
	}

	cpa.EndWithResult("done")

	return nil
}
