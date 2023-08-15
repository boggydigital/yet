package yeti

import (
	"errors"
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	videoExt = ".video"
	videoSfx = " (video)"
	audioExt = ".audio"
	audioSfx = " (audio)"
	fastEnv  = "YET_FAST"
)

type FilenameDelegate func(videoId string, videoPage *yt_urls.InitialPlayerResponse) string

func DownloadVideos(
	httpClient *http.Client,
	filenameDelegate FilenameDelegate,
	binaries *Binaries,
	videoIds ...string) error {

	if len(videoIds) == 0 {
		return nil
	}

	if filenameDelegate == nil {
		return errors.New("filename delegate is nil")
	}

	dvtpw := nod.NewProgress(fmt.Sprintf("downloading %d video(s)", len(videoIds)))
	defer dvtpw.End()

	dvtpw.Total(uint64(len(videoIds)))

	dl := dolo.NewClient(httpClient, dolo.Defaults())

	for _, videoId := range videoIds {

		gv := nod.Begin("video-id: " + videoId)

		videoPage, playerUrl, err := yt_urls.GetVideoPage(httpClient, videoId)
		if err != nil {
			_ = gv.EndWithError(err)
			dvtpw.Increment()
			continue
		}

		fn := filenameDelegate(videoId, videoPage)

		start := time.Now()

		if err := downloadVideo(dl, fn, binaries, videoPage, playerUrl); err != nil {
			gv.Error(err)
		}

		elapsed := time.Since(start)

		gv.EndWithResult("done in %.1fs", elapsed.Seconds())
		dvtpw.Increment()
	}

	return nil
}

func downloadVideo(
	dl *dolo.Client,
	fn string,
	binaries *Binaries,
	videoPage *yt_urls.InitialPlayerResponse,
	playerUrl string) error {

	if _, err := os.Stat(fn); err == nil {
		//local file already exists - won't attempt to download again
		return nil
	}

	vt := videoPage.Title()

	if binaries.FFMpeg == "" {
		if err := downloadSingleFormat(dl, binaries, vt, fn, videoPage.Formats(), playerUrl); err != nil {
			return err
		}
	} else {
		if err := downloadAdaptiveFormat(
			dl,
			binaries,
			vt,
			fn,
			videoPage.AdaptiveVideoFormats(),
			videoPage.AdaptiveAudioFormats(),
			playerUrl); err != nil {
			return err
		}
	}

	//set file modification time to video publish date to allow OS sorting based on mod time
	if _, err := os.Stat(fn); err == nil {
		if err := os.Chtimes(fn, videoPage.PublishDate(), videoPage.PublishDate()); err != nil {
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

func downloadSingleFormat(dl *dolo.Client, binaries *Binaries, title, filename string, formats yt_urls.Formats, playerUrl string) error {

	for _, format := range formats {

		if format.Url == "" {
			continue
		}

		tpw := nod.NewProgress("title: " + title)

		u, err := url.Parse(format.Url)
		if err != nil {
			_ = tpw.EndWithError(err)
			continue
		}

		fast := os.Getenv(fastEnv) != ""

		if binaries.NodeJS != "" || binaries.Deno != "" || fast {
			q := u.Query()
			np := q.Get("n")
			if dnp, err := decodeParam(http.DefaultClient, binaries, np, playerUrl); err != nil {
				return tpw.EndWithError(err)
			} else {
				q.Set("n", dnp)
				u.RawQuery = q.Encode()
			}
		}

		if err := dl.Download(u, tpw, filename); err != nil {
			_ = tpw.EndWithError(err)
			continue
		}

		tpw.EndWithResult("done")

		//yt_urls.StreamingUrls returns bitrate sorted video urls,
		//so we can stop, if we've successfully got the best streaming quality
		break
	}

	return nil
}

func downloadAdaptiveFormat(dl *dolo.Client, binaries *Binaries, title, filename string, videoFormats, audioFormats yt_urls.Formats, playerUrl string) error {

	ext := filepath.Ext(filename)
	fse := strings.TrimSuffix(filename, ext)

	//download video format
	videoTitle := title + videoSfx
	videoFilename := fse + videoExt
	if err := downloadSingleFormat(dl, binaries, videoTitle, videoFilename, videoFormats, playerUrl); err != nil {
		return err
	}

	//download audio format
	audioTitle := title + audioSfx
	audioFilename := fse + audioExt
	if err := downloadSingleFormat(dl, binaries, audioTitle, audioFilename, audioFormats, playerUrl); err != nil {
		return err
	}

	//merge streams into a single file
	//since yt_urls filters to mp4 formats only, we don't need to do any transcoding
	//and can quickly merge by copying streams:
	//ffmpeg -i video.mp4 -i audio.wav -c copy output.mp4
	ma := nod.Begin("merging streams: %s...", title)
	args := []string{"-i", videoFilename, "-i", audioFilename, "-c", "copy", filename}
	cmd := exec.Command(binaries.FFMpeg, args...)
	if err := cmd.Run(); err != nil {
		return ma.EndWithError(err)
	}

	//cleanup separate streams after successful merge
	if err := os.Remove(videoFilename); err != nil {
		return ma.EndWithError(err)
	}
	if err := os.Remove(audioFilename); err != nil {
		return ma.EndWithError(err)
	}

	ma.EndWithResult("done")
	return nil
}
