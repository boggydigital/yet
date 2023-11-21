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

func GetVideoHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	force := u.Query().Has("force")
	return GetVideo(force, ids...)
}

func GetVideo(force bool, videoIds ...string) error {

	if len(videoIds) == 0 {
		return nil
	}

	gva := nod.NewProgress(fmt.Sprintf("downloading %d video(s)", len(videoIds)))
	defer gva.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return gva.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.AllProperties()...)
	if err != nil {
		return gva.EndWithError(err)
	}

	videosDir, err := paths.GetAbsDir(paths.Videos)
	if err != nil {
		return gva.EndWithError(err)
	}

	gva.Total(uint64(len(videoIds)))

	httpClient := http.DefaultClient

	dl := dolo.NewClient(httpClient, dolo.Defaults())

	for _, videoId := range videoIds {

		gv := nod.Begin("video-id: " + videoId)

		// check known errors before doing anything else
		if !force {
			if knownError, ok := rxa.GetFirstVal(data.VideoErrorsProperty, videoId); ok && knownError != "" {
				if err := completeVideo(rxa, videoId, gva, gv, knownError); err != nil {
					return gva.EndWithError(err)
				}
				continue
			}
		}

		// check if the video file matching videoId is already available locally
		if !force && videoExistsLocally(rxa, videosDir, videoId) {
			if err := completeVideo(rxa, videoId, gva, gv, "already exists"); err != nil {
				return gva.EndWithError(err)
			}
			continue
		}

		// adding to the queue before attempting to download
		if err := rxa.AddValues(data.VideosDownloadQueueProperty, videoId, data.TrueValue); err != nil {
			return gv.EndWithError(err)
		}

		videoPage, err := yt_urls.GetVideoPage(httpClient, videoId)
		if err != nil {
			if rerr := rxa.ReplaceValues(data.VideoErrorsProperty, videoId, err.Error()); rerr != nil {
				return gva.EndWithError(rerr)
			}
			if err := completeVideo(rxa, videoId, gva, gv, err.Error()); err != nil {
				return gva.EndWithError(err)
			}
			continue
		}

		for p, v := range yeti.ExtractMetadata(videoPage) {
			if err := rxa.AddValues(p, videoId, v...); err != nil {
				return gv.EndWithError(err)
			}
		}

		thumbnails := videoPage.VideoDetails.Thumbnail.Thumbnails
		if err := yeti.GetPosters(dl, videoId, thumbnails); err != nil {
			return gv.EndWithError(err)
		}

		captionTracks := videoPage.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks
		if err := yeti.GetCaptions(dl, rxa, videoId, captionTracks); err != nil {
			return gv.EndWithError(err)
		}

		relFilename := yeti.DefaultFilenameDelegate(videoId, videoPage)

		start := time.Now()

		if err := downloadVideo(dl, relFilename, videoPage); err != nil {
			gv.Error(err)
		}

		// remove from the queue upon successful download
		if err := rxa.CutVal(data.VideosDownloadQueueProperty, videoId, data.TrueValue); err != nil {
			return gv.EndWithError(err)
		}

		// add to watchlist upon successful download
		if err := rxa.AddValues(data.VideosWatchlistProperty, videoId, data.TrueValue); err != nil {
			return gv.EndWithError(err)
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

func videoExistsLocally(rxa kvas.ReduxAssets, videosDir, videoId string) bool {
	// check if the video file matching videoId is already available locally
	if title, ok := rxa.GetFirstVal(data.VideoTitleProperty, videoId); ok {
		if channel, ok := rxa.GetFirstVal(data.VideoOwnerChannelNameProperty, videoId); ok {
			relVideoFilename := yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
			absVideoFilename := filepath.Join(videosDir, relVideoFilename)
			if _, err := os.Stat(absVideoFilename); err == nil {
				return true
			}
		}
	}
	return false
}

func completeVideo(rxa kvas.ReduxAssets, videoId string, cmd nod.TotalProgressWriter, video nod.ActCloser, result string) error {
	video.EndWithResult(result)
	cmd.Increment()
	return rxa.CutVal(data.VideosDownloadQueueProperty, videoId, data.TrueValue)
}
