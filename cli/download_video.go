package cli

import (
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var defaultYtDlpOptions = map[string]string{
	"-S": "vcodec:h264,res:1080,acodec:m4a",
}

func DownloadVideoHandler(u *url.URL) error {
	q := u.Query()

	videoId := q.Get("video-id")
	options := &VideoOptions{
		Force: q.Has("force"),
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

	if err := getVideoPageMetadata(videoPage, videoId, rdx); err != nil {
		da.Error(err)
		errs = true
	}

	if err := downloadVideo(videoId, videoPage, opt); err != nil {
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
	videoId string,
	videoPage *youtube_urls.InitialPlayerResponse,
	options *VideoOptions) error {

	relFilename := yeti.DefaultFilenameDelegate(videoId, videoPage)

	absVideosDir, err := pathways.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	absFilename := filepath.Join(absVideosDir, relFilename)

	if _, err := os.Stat(absFilename); err == nil {
		if options.Force {
			if err := os.Remove(absFilename); err != nil {
				return err
			}
		} else {
			//local file already exists - won't attempt to download again
			return nil
		}
	}

	if err := downloadWithYtDlp(videoId, absFilename); err != nil {
		return err
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

func downloadWithYtDlp(videoId, absFilename string) error {

	dyda := nod.Begin(" downloading %s with yt-dlp, please wait...", videoId)
	defer dyda.EndWithResult("done")

	absDir, _ := path.Split(absFilename)
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		if err := os.MkdirAll(absDir, 0755); err != nil {
			return dyda.EndWithError(err)
		}
	}

	options := make([]string, 0)

	options = append(options, videoId)
	options = append(options, "-o", absFilename)

	for flag, value := range defaultYtDlpOptions {
		options = append(options, flag, value)
	}

	ytDlpDir, err := pathways.GetAbsDir(paths.YtDlp)
	if err != nil {
		return dyda.EndWithError(err)
	}

	absYtDlpFilename := filepath.Join(ytDlpDir, yeti.GetYtDlpBinary())

	cmd := exec.Command(absYtDlpFilename, options...)
	return cmd.Run()
}
