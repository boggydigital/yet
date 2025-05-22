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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const ytDlpCookiesFilename = "cookies.txt"

var defaultYtDlpOptions = map[string]string{
	"-S": "vcodec:h264,res:1080,acodec:m4a",
}

func DownloadVideoHandler(u *url.URL) error {
	q := u.Query()

	videoIds := strings.Split(q.Get("video-id"), ",")

	options := &VideoOptions{
		BgUtilBaseUrl: q.Get("bgutil-baseurl"),
		Ended:         q.Has("mark-watched"),
		Verbose:       q.Has("verbose"),
		Force:         q.Has("force"),
	}

	return DownloadVideo(nil, options, videoIds...)
}

func DownloadVideo(rdx redux.Writeable, opt *VideoOptions, videoIds ...string) error {

	da := nod.NewProgress("downloading videos...")
	defer da.Done()

	if opt == nil {
		opt = DefaultVideoOptions()
	}

	var err error
	rdx, err = validateWritableRedux(rdx, data.VideoProperties()...)
	if err != nil {
		return err
	}

	da.TotalInt(len(videoIds))

	for _, videoId := range videoIds {

		videoId, err = yeti.ParseVideoId(videoId)
		if err != nil {
			return err
		}

		// apply video specific options
		opt = ApplyVideoDownloadOptions(opt, videoId, rdx)

		errs := false

		// adding to download queue (if not there already)
		if !rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			if err := rdx.AddValues(data.VideoDownloadQueuedProperty, videoId, yeti.FmtNow()); err != nil {
				return err
			}
		}

		// setting download started timestamp
		if err := rdx.AddValues(data.VideoDownloadStartedProperty, videoId, yeti.FmtNow()); err != nil {
			return err
		}

		videoPage, err := yeti.GetVideoPage(videoId)
		if err != nil {
			return err
		}

		if err = getVideoPageMetadata(videoPage, videoId, rdx); err != nil {
			da.Error(err)
			errs = true
		}

		if err = downloadVideo(videoId, videoPage, opt); err != nil {
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
				return err
			}
		} else {
			// reset download started if errors were encountered (keeping in download queue)
			if err := rdx.CutKeys(data.VideoDownloadStartedProperty, videoId); err != nil {
				return err
			}
		}

		da.Increment()
	}

	return nil
}

func downloadVideo(
	videoId string,
	videoPage *youtube_urls.InitialPlayerResponse,
	options *VideoOptions) error {

	var title, channel string
	if videoPage != nil {
		title = videoPage.VideoDetails.Title
		channel = videoPage.VideoDetails.Author
	}

	relFilename := yeti.RelLocalVideoFilename(channel, title, videoId)

	absVideosDir, err := pathways.GetAbsDir(data.Videos)
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

	if err = downloadWithYtDlp(videoId, absFilename, options); err != nil {
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

func downloadWithYtDlp(videoId, absFilename string, options *VideoOptions) error {

	dyda := nod.Begin(" downloading %s with yt-dlp, please wait...", videoId)
	defer dyda.Done()

	absDir, _ := path.Split(absFilename)
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		if err := os.MkdirAll(absDir, 0755); err != nil {
			return err
		}
	}

	ytDlpDir, err := pathways.GetAbsDir(data.YtDlp)
	if err != nil {
		return err
	}

	arguments := make([]string, 0)

	if strings.HasPrefix(videoId, "-") {
		arguments = append(arguments, youtube_urls.VideoUrl(videoId).String())
	} else {
		arguments = append(arguments, videoId)
	}

	arguments = append(arguments, "-o", absFilename)

	if options.Ended {
		arguments = append(arguments, "--mark-watched")
	}

	if options.BgUtilBaseUrl != "" {
		arguments = append(arguments, "--extractor-args", "youtube:getpot_bgutil_baseurl="+options.BgUtilBaseUrl)
	}

	absYtDlpCookiesPath := filepath.Join(ytDlpDir, ytDlpCookiesFilename)

	if _, err = os.Stat(absYtDlpCookiesPath); err == nil {
		arguments = append(arguments, "--cookies", absYtDlpCookiesPath)
	}

	for flag, value := range defaultYtDlpOptions {
		arguments = append(arguments, flag, value)
	}

	absYtDlpFilename := filepath.Join(ytDlpDir, yeti.GetYtDlpBinary())

	cmd := exec.Command(absYtDlpFilename, arguments...)

	if options.Verbose {
		arguments = append(arguments, "-vU")

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}
