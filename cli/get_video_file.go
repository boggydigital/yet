package cli

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	fastEnv = "YET_FAST"
)

func GetVideoFileHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	force := u.Query().Has("force")
	return GetVideoFile(force, ids...)
}

func GetVideoFile(force bool, ids ...string) error {

	if len(ids) == 0 {
		return nil
	}

	gva := nod.NewProgress(fmt.Sprintf("getting %d video(s)", len(ids)))
	defer gva.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return gva.EndWithError(err)
	}

	rdx, err := kvas.ReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return gva.EndWithError(err)
	}

	videosDir, err := paths.GetAbsDir(paths.Videos)
	if err != nil {
		return gva.EndWithError(err)
	}

	videoIds, err := yeti.ParseVideoIds(ids...)
	if err != nil {
		return gva.EndWithError(err)
	}

	gva.Total(uint64(len(ids)))

	for _, videoId := range videoIds {

		gv := nod.Begin("video-id: " + videoId)

		// check known errors before doing anything else
		if !force {
			if knownError, ok := rdx.GetFirstVal(data.VideoErrorsProperty, videoId); ok && knownError != "" {
				if err := completeVideo(rdx, videoId, gva, gv, knownError); err != nil {
					return gva.EndWithError(err)
				}
				continue
			}
		}

		// check if the video file matching videoId is already available locally
		if !force && videoExistsLocally(rdx, videosDir, videoId) {
			if err := completeVideo(rdx, videoId, gva, gv, "already exists"); err != nil {
				return gva.EndWithError(err)
			}
			continue
		}

		videoPage, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			if rerr := rdx.ReplaceValues(data.VideoErrorsProperty, videoId, err.Error()); rerr != nil {
				return gva.EndWithError(rerr)
			}
			if err := completeVideo(rdx, videoId, gva, gv, err.Error()); err != nil {
				return gva.EndWithError(err)
			}
			continue
		}

		relFilename := yeti.DefaultFilenameDelegate(videoId, videoPage)

		start := time.Now()

		if err := downloadVideo(dolo.DefaultClient, relFilename, videoPage); err != nil {
			gv.Error(err)
		}

		elapsed := time.Since(start)

		gv.EndWithResult("done in %.1fs", elapsed.Seconds())
		gva.Increment()
	}

	return nil
}

func downloadVideo(
	dl *dolo.Client,
	relFilename string,
	videoPage *yt_urls.InitialPlayerResponse) error {

	absVideosDir, err := paths.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	absFilename := filepath.Join(absVideosDir, relFilename)

	if _, err := os.Stat(absFilename); err == nil {
		//local file already exists - won't attempt to download again
		return nil
	}

	if yeti.GetBinary(yeti.FFMpegBin) == "" {
		if err := downloadSingleFormat(dl, relFilename, videoPage.Formats(), videoPage.PlayerUrl); err != nil {
			return err
		}
	} else {
		if err := downloadAdaptiveFormat(dl, relFilename, videoPage); err != nil {
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

func downloadSingleFormat(dl *dolo.Client, relFilename string, formats yt_urls.Formats, playerUrl string) error {

	for _, format := range formats {

		if format.Url == "" {
			continue
		}

		tpw := nod.NewProgress("file: " + relFilename)

		u, err := url.Parse(format.Url)
		if err != nil {
			_ = tpw.EndWithError(err)
			continue
		}

		fast := os.Getenv(fastEnv) != ""

		if yeti.IsJSBinaryAvailable() || fast {
			q := u.Query()
			np := q.Get("n")
			if dnp, err := yeti.DecodeParam(http.DefaultClient, np, playerUrl); err != nil {
				return tpw.EndWithError(err)
			} else {
				q.Set("n", dnp)
				u.RawQuery = q.Encode()
			}
		}

		absVideosDir, err := paths.GetAbsDir(paths.Videos)
		if err != nil {
			return tpw.EndWithError(err)
		}

		if err := dl.Download(u, tpw, absVideosDir, relFilename); err != nil {
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

func downloadAdaptiveFormat(dl *dolo.Client, relFilename string, videoPage *yt_urls.InitialPlayerResponse) error {

	relVideoFilename, relAudioFilename := yeti.VideoAudioFilenames(relFilename)

	//download video format
	if err := downloadSingleFormat(dl, relVideoFilename, videoPage.AdaptiveVideoFormats(), videoPage.PlayerUrl); err != nil {
		return err
	}

	//download audio format
	if err := downloadSingleFormat(dl, relAudioFilename, videoPage.AdaptiveAudioFormats(), videoPage.PlayerUrl); err != nil {
		return err
	}

	if err := yeti.MergeStreams(relFilename); err != nil {
		return err
	}

	return nil
}

func videoExistsLocally(rdx kvas.ReadableRedux, videosDir, videoId string) bool {
	// check if the video file matching videoId is already available locally
	if title, ok := rdx.GetFirstVal(data.VideoTitleProperty, videoId); ok {
		if channel, ok := rdx.GetFirstVal(data.VideoOwnerChannelNameProperty, videoId); ok {
			relVideoFilename := yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
			absVideoFilename := filepath.Join(videosDir, relVideoFilename)
			if _, err := os.Stat(absVideoFilename); err == nil {
				return true
			}
		}
	}
	return false
}

func completeVideo(rdx kvas.WriteableRedux, videoId string, cmd nod.TotalProgressWriter, video nod.ActCloser, result string) error {
	video.EndWithResult(result)
	cmd.Increment()
	return rdx.CutValues(data.VideosDownloadQueueProperty, videoId, data.TrueValue)
}
