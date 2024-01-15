package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/url"
	"strings"
)

func GetCaptionsHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	return GetCaptions(ids)
}

func GetCaptions(ids []string) error {

	gca := nod.NewProgress("getting captions...")
	defer gca.End()

	gca.TotalInt(len(ids))

	dl := dolo.DefaultClient

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		return gca.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir,
		data.VideoCaptionsNamesProperty,
		data.VideoCaptionsKindsProperty,
		data.VideoCaptionsLanguagesProperty)
	if err != nil {
		return gca.EndWithError(err)
	}

	for _, videoId := range ids {

		if err := getVideoPageCaptions(nil, videoId, rdx, dl); err != nil {
			gca.Error(err)
		}

		gca.Increment()
	}

	gca.EndWithResult("done")

	return nil
}

func getVideoPageCaptions(videoPage *yt_urls.InitialPlayerResponse, videoId string, rdx kvas.WriteableRedux, dl *dolo.Client) error {

	gca := nod.Begin(" captions for %s", videoId)
	defer gca.End()

	var err error
	if videoPage == nil {
		videoPage, err = yeti.GetVideoPage(videoId)
		if err != nil {
			return gca.EndWithError(err)
		}
	}

	captionTracks := videoPage.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks
	if err := yeti.GetCaptions(dl, rdx, videoId, captionTracks); err != nil {
		return gca.EndWithError(err)
	}

	gca.EndWithResult("done")

	return nil
}
