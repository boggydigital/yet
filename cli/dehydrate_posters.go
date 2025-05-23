package cli

import (
	"errors"
	"github.com/boggydigital/issa"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
	"os"
)

var (
	ErrVideoHasNoPosterThumbnail = errors.New("video has no poster thumbnails")
)

func DehydratePostersHandler(u *url.URL) error {
	force := u.Query().Has("force")
	return DehydratePosters(force)
}

func DehydratePosters(force bool) error {

	dpa := nod.NewProgress("dehydrating posters...")
	defer dpa.Done()

	metadataDir, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return err
	}

	rdx, err := redux.NewWriter(metadataDir, data.VideoProperties()...)
	if err != nil {
		return err
	}

	dpa.TotalInt(rdx.Len(data.VideoTitleProperty))

	dehydratedPosters := make(map[string][]string)
	dehydratedRepColors := make(map[string][]string)
	dehydratedInputMissing := make(map[string][]string)

	for videoId := range rdx.Keys(data.VideoTitleProperty) {

		if rdx.HasKey(data.VideoDehydratedInputMissingProperty, videoId) && !force {
			dpa.Increment()
			continue
		}

		if rdx.HasKey(data.VideoDehydratedPosterProperty, videoId) && !force {
			dpa.Increment()
			continue
		}

		if dp, rc, err := dehydratePosterImageRepColor(videoId); err == nil {
			dehydratedPosters[videoId] = append(dehydratedPosters[videoId], dp)
			dehydratedRepColors[videoId] = append(dehydratedRepColors[videoId], rc)
		} else if errors.Is(err, ErrVideoHasNoPosterThumbnail) {
			dehydratedInputMissing[videoId] = append(dehydratedInputMissing[videoId], "true")
		} else {
			dpa.Error(err)
		}

		dpa.Increment()
	}

	if err := rdx.BatchReplaceValues(data.VideoDehydratedPosterProperty, dehydratedPosters); err != nil {
		return err
	}

	if err := rdx.BatchReplaceValues(data.VideoDehydratedRepColorProperty, dehydratedRepColors); err != nil {
		return err
	}

	if err := rdx.BatchReplaceValues(data.VideoDehydratedInputMissingProperty, dehydratedInputMissing); err != nil {
		return err
	}

	return nil
}

func dehydratePosterImageRepColor(videoId string) (string, string, error) {

	var absPosterPath string
	var err error

	// find the first existing poster (if any if available at all)
	for _, q := range youtube_urls.AllThumbnailQualities() {
		absPosterPath, err = data.AbsPosterPath(videoId, q)
		if err != nil {
			return "", "", err
		}
		if _, err := os.Stat(absPosterPath); err == nil {
			break
		}
		absPosterPath = ""
	}

	if absPosterPath == "" {
		return "", "", ErrVideoHasNoPosterThumbnail
	}

	return issa.DehydrateImageRepColor(absPosterPath)
}
