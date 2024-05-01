package cli

import (
	"errors"
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
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

func GetVideoFileHandler(u *url.URL) error {
	ids := strings.Split(u.Query().Get("id"), ",")
	force := u.Query().Has("force")
	singleFormat := u.Query().Has("single-format")
	return GetVideoFile(force, singleFormat, ids...)
}

func GetVideoFile(force, singleFormat bool, ids ...string) error {

	if len(ids) == 0 {
		return nil
	}

	gva := nod.NewProgress(fmt.Sprintf("getting %d video(s)", len(ids)))
	defer gva.End()

	metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	if err != nil {
		return gva.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return gva.EndWithError(err)
	}

	videosDir, err := pasu.GetAbsDir(paths.Videos)
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
			if knownError, ok := rdx.GetLastVal(data.VideoErrorsProperty, videoId); ok && knownError != "" {
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

		videoPage, err := yeti.GetVideoPage(videoId)
		if err != nil {
			if rerr := rdx.ReplaceValues(data.VideoErrorsProperty, videoId, err.Error()); rerr != nil {
				return gva.EndWithError(rerr)
			}
			if err := completeVideo(rdx, videoId, gva, gv, err.Error()); err != nil {
				return gva.EndWithError(err)
			}
			continue
		}

		if err := yeti.DecodeSignatureCipher(http.DefaultClient, videoPage); err != nil {
			return gva.EndWithError(err)
		}

		relFilename := yeti.DefaultFilenameDelegate(videoId, videoPage)

		start := time.Now()

		if err := downloadVideo(dolo.DefaultClient, relFilename, videoPage, force, singleFormat); err != nil {
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
	videoId string,
	videoPage *yt_urls.InitialPlayerResponse,
	force bool,
	singleFormat bool) error {

	relFilename := yeti.DefaultFilenameDelegate(videoId, videoPage)

	absVideosDir, err := pasu.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	absFilename := filepath.Join(absVideosDir, relFilename)

	if _, err := os.Stat(absFilename); !force && err == nil {
		//local file already exists - won't attempt to download again
		return nil
	}

	if yeti.GetBinary(yeti.FFMpegBin) == "" || singleFormat {

		if err := downloadSingleFormat(dl, relFilename, videoPage.BestFormat(), videoPage.PlayerUrl, force); err != nil {
			return err
		}
	} else {
		if err := downloadAdaptiveFormat(dl, videoId, relFilename, videoPage, force); err != nil {
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
	format *yt_urls.Format,
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
		np := q.Get("n")
		if dnp, err := yeti.DecodeNParam(np, playerUrl); err != nil {
			return tpw.EndWithError(err)
		} else {
			q.Set("n", dnp)
			u.RawQuery = q.Encode()
		}
	}

	absVideosDir, err := pasu.GetAbsDir(paths.Videos)
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

func downloadAdaptiveFormat(dl *dolo.Client, videoId, relFilename string, vp *yt_urls.InitialPlayerResponse, force bool) error {

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

func videoExistsLocally(rdx kvas.ReadableRedux, videosDir, videoId string) bool {
	// check if the video file matching videoId is already available locally
	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok {
		if channel, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok {
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
