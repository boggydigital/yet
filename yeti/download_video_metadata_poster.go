package yeti

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

func DownloadVideoMetadataPoster(act nod.ActCloser, videoId string, opt *VideoOptions, rdx redux.Writeable) error {

	var err error
	videoId, err = ParseVideoId(videoId)
	if err != nil {
		return err
	}

	// apply video specific options
	opt = ApplyVideoDownloadOptions(opt, videoId, rdx)

	errs := false

	// adding to download queue (if not there already)
	if !rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
		if err = rdx.AddValues(data.VideoDownloadQueuedProperty, videoId, FmtNow()); err != nil {
			return err
		}
	}

	// setting download started timestamp
	if err = rdx.AddValues(data.VideoDownloadStartedProperty, videoId, FmtNow()); err != nil {
		return err
	}

	var videoPage *youtube_urls.InitialPlayerResponse
	videoPage, err = GetVideoPage(videoId)
	if err != nil {
		return err
	}

	if err = GetVideoPageMetadata(videoPage, videoId, rdx); err != nil {
		if act != nil {
			act.Error(err)
		}
		errs = true
	}

	if err = DownloadVideo(videoId, videoPage, opt); err != nil {
		if act != nil {
			act.Error(err)
		}
		errs = true
	}

	if err = GetPosters(videoId, dolo.DefaultClient, opt.Force, youtube_urls.AllThumbnailQualities()...); err != nil {
		if act != nil {
			act.Error(err)
		}
		errs = true
	}

	if !errs {
		// set downloaded date if no errors were encountered
		if err = rdx.AddValues(data.VideoDownloadCompletedProperty, videoId, FmtNow()); err != nil {
			return err
		}
	} else {
		// reset download started if errors were encountered (keeping in download queue)
		if err = rdx.CutKeys(data.VideoDownloadStartedProperty, videoId); err != nil {
			return err
		}
	}

	return nil
}
