package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
	"strings"
)

func GetCaptionsHandler(u *url.URL) error {
	videoIds := strings.Split(u.Query().Get("video-id"), ",")
	force := u.Query().Has("force")
	return GetCaptions(force, videoIds...)
}

func GetCaptions(force bool, videoIds ...string) error {

	gca := nod.NewProgress("getting captions...")
	defer gca.End()

	gca.TotalInt(len(videoIds))

	dl := dolo.DefaultClient

	metadataDir, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return gca.EndWithError(err)
	}

	rdx, err := redux.NewWriter(metadataDir,
		data.VideoCaptionsNamesProperty,
		data.VideoCaptionsKindsProperty,
		data.VideoCaptionsLanguagesProperty)
	if err != nil {
		return gca.EndWithError(err)
	}

	for _, videoId := range videoIds {

		if err := getVideoPageCaptions(nil, videoId, rdx, dl, force); err != nil {
			gca.Error(err)
		}

		gca.Increment()
	}

	gca.EndWithResult("done")

	return nil
}

func getVideoPageCaptions(
	videoPage *youtube_urls.InitialPlayerResponse,
	videoId string,
	rdx redux.Writeable,
	dl *dolo.Client,
	force bool) error {

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
	if err := yeti.GetCaptions(dl, rdx, videoId, captionTracks, force); err != nil {
		return gca.EndWithError(err)
	}

	gca.EndWithResult("done")

	return nil
}
