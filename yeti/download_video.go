package yeti

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/boggydigital/camino"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

const ytDlpCookiesFilename = "cookies.txt"

const bgutilBaseUrlParam = "youtubepot-bgutilhttp:base_url"

var defaultYtDlpOptions = map[string]string{
	"-S": "vcodec:h264,res:1080,acodec:m4a",
}

func DownloadVideo(
	videoId string,
	videoPage *youtube_urls.InitialPlayerResponse,
	options *VideoOptions) error {

	var title, channel string
	if videoPage != nil {
		title = videoPage.VideoDetails.Title
		channel = videoPage.VideoDetails.Author
	}

	relFilename := RelLocalVideoFilename(channel, title, videoId)

	absVideosDir := camino.GetAbs(data.Videos)
	absFilename := filepath.Join(absVideosDir, relFilename)

	if _, err := os.Stat(absFilename); err == nil {
		if options.Force {
			if err = os.Remove(absFilename); err != nil {
				return err
			}
		} else {
			//local file already exists - won't attempt to download again
			return nil
		}
	}

	if err := downloadWithYtDlp(videoId, absFilename, options); err != nil {
		return err
	}

	//set file modification time to video publish date to allow OS sorting based on mod time
	if _, err := os.Stat(absFilename); err == nil {
		if err = os.Chtimes(absFilename, videoPage.PublishDate(), videoPage.PublishDate()); err != nil {
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
		if err = os.MkdirAll(absDir, 0755); err != nil {
			return err
		}
	}

	ytDlpDir := camino.GetAbs(data.YtDlp)

	arguments := make([]string, 0)

	absYtDlpPluginsDir := camino.GetRel(data.YtDlpPlugins, data.YtDlp)

	arguments = append(arguments, "--plugin-dirs", absYtDlpPluginsDir)

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
		arguments = append(arguments, "--extractor-args", strings.Join([]string{bgutilBaseUrlParam, options.BgUtilBaseUrl}, "="))
	}

	absYtDlpCookiesPath := filepath.Join(ytDlpDir, ytDlpCookiesFilename)

	if _, err := os.Stat(absYtDlpCookiesPath); err == nil {
		arguments = append(arguments, "--cookies", absYtDlpCookiesPath)
	}

	for flag, value := range defaultYtDlpOptions {
		arguments = append(arguments, flag, value)
	}

	absYtDlpFilename := filepath.Join(ytDlpDir, GetYtDlpBinary())

	cmd := exec.Command(absYtDlpFilename, arguments...)

	if options.Verbose {
		arguments = append(arguments, "-vU")

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}
