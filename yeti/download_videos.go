package yeti

import (
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	fastEnv = "YET_FAST"
)

func DownloadVideos(
	httpClient *http.Client,
	videoIds ...string) error {

	if len(videoIds) == 0 {
		return nil
	}

	dvtpw := nod.NewProgress(fmt.Sprintf("downloading %d video(s)", len(videoIds)))
	defer dvtpw.End()

	metadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		return dvtpw.EndWithError(err)
	}

	videosDir, err := paths.GetAbsDir(paths.Videos)
	if err != nil {
		return dvtpw.EndWithError(err)
	}

	rxa, err := kvas.ConnectReduxAssets(metadataDir, data.AllProperties()...)
	if err != nil {
		return dvtpw.EndWithError(err)
	}

	dvtpw.Total(uint64(len(videoIds)))

	dl := dolo.NewClient(httpClient, dolo.Defaults())

	for _, videoId := range videoIds {

		gv := nod.Begin("video-id: " + videoId)

		// check known errors before doing anything else
		if knownError, ok := rxa.GetFirstVal(data.VideoErrorsProperty, videoId); ok && knownError != "" {
			gv.EndWithResult(knownError)
			dvtpw.Increment()
			continue
		}

		// check if the video file matching videoId is already available locally
		if title, ok := rxa.GetFirstVal(data.VideoTitleProperty, videoId); ok {
			if channel, ok := rxa.GetFirstVal(data.VideoOwnerChannelNameProperty, videoId); ok {
				relVideoFilename := ChannelTitleVideoIdFilename(channel, title, videoId)
				absVideoFilename := filepath.Join(videosDir, relVideoFilename)
				if _, err := os.Stat(absVideoFilename); err == nil {
					gv.EndWithResult("already exists")
					dvtpw.Increment()
					continue
				}
			}
		}

		videoPage, playerUrl, err := yt_urls.GetVideoPage(httpClient, videoId)
		if err != nil {
			if rerr := rxa.ReplaceValues(data.VideoErrorsProperty, videoId, err.Error()); rerr != nil {
				return dvtpw.EndWithError(rerr)
			}
			_ = gv.EndWithError(err)
			dvtpw.Increment()
			continue
		}

		for p, v := range videoPageMetadata(videoPage) {
			if err := rxa.ReplaceValues(p, videoId, v...); err != nil {
				return gv.EndWithError(err)
			}
		}

		if err := GetThumbnails(dl, videoId, videoPage.VideoDetails.Thumbnail.Thumbnails); err != nil {
			return gv.EndWithError(err)
		}

		relFilename := DefaultFilenameDelegate(videoId, videoPage)

		start := time.Now()

		if err := downloadVideo(dl, relFilename, videoPage, playerUrl); err != nil {
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
	relFilename string,
	videoPage *yt_urls.InitialPlayerResponse,
	playerUrl string) error {

	absVideosDir, err := paths.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	absFilename := filepath.Join(absVideosDir, relFilename)

	if _, err := os.Stat(absFilename); err == nil {
		//local file already exists - won't attempt to download again
		return nil
	}

	if GetBinary(FFMpegBin) == "" {
		if err := downloadSingleFormat(dl, relFilename, videoPage.Formats(), playerUrl); err != nil {
			return err
		}
	} else {
		if err := downloadAdaptiveFormat(
			dl,
			relFilename,
			videoPage.AdaptiveVideoFormats(),
			videoPage.AdaptiveAudioFormats(),
			playerUrl); err != nil {
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

		if IsJSBinaryAvailable() || fast {
			q := u.Query()
			np := q.Get("n")
			if dnp, err := DecodeParam(http.DefaultClient, np, playerUrl); err != nil {
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

func downloadAdaptiveFormat(dl *dolo.Client, relFilename string, videoFormats, audioFormats yt_urls.Formats, playerUrl string) error {

	relVideoFilename, relAudioFilename := videoAudioFilenames(relFilename)

	//download video format
	if err := downloadSingleFormat(dl, relVideoFilename, videoFormats, playerUrl); err != nil {
		return err
	}

	//download audio format
	if err := downloadSingleFormat(dl, relAudioFilename, audioFormats, playerUrl); err != nil {
		return err
	}

	if err := mergeStreams(relFilename); err != nil {
		return err
	}

	return nil
}

func videoPageMetadata(ipr *yt_urls.InitialPlayerResponse) map[string][]string {
	vpm := make(map[string][]string)

	vpm[data.VideoTitleProperty] = []string{ipr.VideoDetails.Title}
	vpm[data.VideoThumbnailUrlsProperty] = make([]string, 0, len(ipr.VideoDetails.Thumbnail.Thumbnails))
	for _, t := range ipr.VideoDetails.Thumbnail.Thumbnails {
		vpm[data.VideoThumbnailUrlsProperty] = append(vpm[data.VideoThumbnailUrlsProperty], t.Url)
	}
	vpm[data.VideoExternalChannelIdProperty] = []string{ipr.VideoDetails.ChannelId}
	vpm[data.VideoShortDescriptionProperty] = []string{ipr.VideoDetails.ShortDescription}
	vpm[data.VideoViewCountProperty] = []string{ipr.VideoDetails.ViewCount}
	vpm[data.VideoKeywordsProperty] = ipr.VideoDetails.Keywords

	vpm[data.VideoOwnerChannelNameProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.OwnerChannelName}
	vpm[data.VideoOwnerProfileUrlProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.OwnerProfileUrl}
	vpm[data.VideoCategoryProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.Category}
	vpm[data.VideoPublishDateProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.PublishDate}
	vpm[data.VideoUploadDateProperty] = []string{ipr.Microformat.PlayerMicroformatRenderer.UploadDate}

	return vpm
}
