package cli

import (
	"errors"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func DownloadVideoHandler(u *url.URL) error {
	q := u.Query()

	videoId := q.Get("video-id")
	options := &VideoOptions{
		PreferSingleFormat: q.Has("prefer-single-format"),
		Force:              q.Has("force"),
		Source:             q.Get("source"),
	}

	return DownloadVideo(nil, videoId, options)
}

func DownloadVideo(rdx kevlar.WriteableRedux, videoId string, opt *VideoOptions) error {

	da := nod.Begin("downloading video %s...", videoId)
	defer da.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return da.EndWithError(err)
	}

	videoId, err = yeti.ParseVideoId(videoId)
	if err != nil {
		return da.EndWithError(err)
	}

	if opt == nil {
		opt = DefaultVideoOptions()
	}
	// apply video specific options
	opt = ApplyVideoDownloadOptions(opt, videoId, rdx)

	errs := false

	// adding to download queue (if not there already)
	if !rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
		if err := rdx.AddValues(data.VideoDownloadQueuedProperty, videoId, yeti.FmtNow()); err != nil {
			return da.EndWithError(err)
		}
	}

	// setting download started timestamp
	if err := rdx.AddValues(data.VideoDownloadStartedProperty, videoId, yeti.FmtNow()); err != nil {
		return da.EndWithError(err)
	}

	videoPage, err := yeti.GetVideoPage(videoId)
	if err != nil {
		return da.EndWithError(err)
	}

	if err := yeti.DecodeSignatureCiphers(http.DefaultClient, videoPage); err != nil {
		return da.EndWithError(err)
	}

	if err := getVideoPageMetadata(videoPage, videoId, rdx); err != nil {
		da.Error(err)
		errs = true
	}

	if err := downloadVideo(dolo.DefaultClient, videoId, videoPage, opt); err != nil {
		da.Error(err)
		errs = true
	}

	if err := yeti.GetPosters(videoId, dolo.DefaultClient, opt.Force, youtube_urls.AllThumbnailQualities()...); err != nil {
		da.Error(err)
		errs = true
	}

	if err := getVideoPageCaptions(videoPage, videoId, rdx, dolo.DefaultClient, opt.Force); err != nil {
		da.Error(err)
		errs = true
	}

	if !errs {
		// set downloaded date if no errors were encountered
		if err := rdx.AddValues(data.VideoDownloadCompletedProperty, videoId, yeti.FmtNow()); err != nil {
			return da.EndWithError(err)
		}
	} else {
		// reset download started if errors were encountered (keeping in download queue)
		if err := rdx.CutKeys(data.VideoDownloadStartedProperty, videoId); err != nil {
			return da.EndWithError(err)
		}
	}

	da.EndWithResult("done")

	return nil
}

func downloadVideo(
	dl *dolo.Client,
	videoId string,
	videoPage *youtube_urls.InitialPlayerResponse,
	options *VideoOptions) error {

	relFilename := yeti.DefaultFilenameDelegate(videoId, videoPage)

	absVideosDir, err := pathways.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	absFilename := filepath.Join(absVideosDir, relFilename)

	if _, err := os.Stat(absFilename); !options.Force && err == nil {
		//local file already exists - won't attempt to download again
		return nil
	}

	if options.Source != "" {
		// download file by URL using single format download method and faking format
		sourceFormat := &youtube_urls.Format{Url: options.Source}
		if err := downloadSingleFormat(dl, relFilename, sourceFormat, videoPage.PlayerUrl, options.Force); err != nil {
			return err
		}
	} else if yeti.GetBinary(yeti.FFMpegBin) == "" || options.PreferSingleFormat {
		if err := downloadSingleFormat(dl, relFilename, videoPage.BestFormat(), videoPage.PlayerUrl, options.Force); err != nil {
			return err
		}
	} else {
		if err := downloadAdaptiveFormat(dl, videoId, relFilename, videoPage, options.Force); err != nil {
			return err
		}
	}

	//set file modification time to video publish date to allow OS sorting based on mod time
	if _, err := os.Stat(absFilename); err == nil {
		if err := os.Chtimes(absFilename, videoPage.PublishDate(), videoPage.PublishDate()); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		//there was an error downloading this version, but we have a partial file
		//we can try resuming next application session
	} else {
		return err
	}

	return nil
}

func downloadSingleFormat(
	dl *dolo.Client,
	relFilename string,
	format *youtube_urls.Format,
	playerUrl string,
	force bool) error {

	if format.Url == "" {
		return errors.New("stream format needs url")
	}

	tpw := nod.NewProgress("file: " + relFilename)

	u, err := url.Parse(format.Url)
	if err != nil {
		return tpw.EndWithError(err)
	}

	if yeti.HasBinary(yeti.NodeBin) {
		q := u.Query()
		if np := q.Get("n"); np != "" {
			if dnp, err := yeti.DecodeNParam(np, playerUrl); err != nil {
				return tpw.EndWithError(err)
			} else {
				if dnp != "" {
					q.Set("n", dnp)
					u.RawQuery = q.Encode()
				}
			}
		}
	}

	absVideosDir, err := pathways.GetAbsDir(paths.Videos)
	if err != nil {
		return tpw.EndWithError(err)
	}

	if force {
		absFilename := filepath.Join(absVideosDir, relFilename)
		if _, err := os.Stat(absFilename); err == nil {
			if err := os.Remove(absFilename); err != nil {
				return tpw.EndWithError(err)
			}
		}
	}

	if err := dl.Download(u, force, tpw, absVideosDir, relFilename); err != nil {
		return tpw.EndWithError(err)
	}

	tpw.EndWithResult("done")

	return nil
}

func downloadAdaptiveFormat(dl *dolo.Client, videoId, relFilename string, vp *youtube_urls.InitialPlayerResponse, force bool) error {

	rvfn, rafn := yeti.VideoAudioFilenames(relFilename)

	//download video format
	if err := downloadSingleFormat(dl, rvfn, vp.BestAdaptiveVideoFormat(), vp.PlayerUrl, force); err != nil {
		return err
	}

	//download audio format
	if err := downloadSingleFormat(dl, rafn, vp.BestAdaptiveAudioFormat(), vp.PlayerUrl, force); err != nil {
		return err
	}

	if err := yeti.MergeStreams(relFilename, force); err != nil {
		return err
	}

	return nil
}
