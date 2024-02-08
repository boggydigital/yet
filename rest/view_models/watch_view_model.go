package view_models

import (
	"fmt"
	"github.com/boggydigital/kvas"
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
)

type WatchViewModel struct {
	VideoId          string
	VideoUrl         string
	VideoPoster      string
	Server           string
	VideoTitle       string
	VideoDescription string
	CurrentTime      string
	LastEndedTime    string
}

func GetWatchViewModel(videoId, currentTime string, rdx kvas.ReadableRedux) (*WatchViewModel, error) {

	videoUrl, videoTitle, videoDescription := "", "", ""
	//var videoCaptionTracks []yt_urls.CaptionTrack
	playbackType := "streaming"

	if strings.HasSuffix(videoId, yt_urls.DefaultVideoExt) {
		// TODO: check if that local file exists first
		playbackType = "local"
		videoUrl = "/video?file=" + videoId
		//videoTitle = videoId
	}

	videoPoster := fmt.Sprintf("/poster?v=%s&q=%s", videoId, yt_urls.ThumbnailQualityMaxRes)

	// if current time is not specified with a query parameter - read it from metadata
	if currentTime == "" {
		if ct, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok {
			currentTime = ct
		}
	}

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && title != "" {
		if channel, ok := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); ok && channel != "" {
			localVideoFilename := yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
			if absVideosDir, err := pasu.GetAbsDir(paths.Videos); err == nil {
				absLocalVideoFilename := filepath.Join(absVideosDir, localVideoFilename)
				if _, err := os.Stat(absLocalVideoFilename); err == nil {
					playbackType = "local"
					videoUrl = "/video?file=" + url.QueryEscape(localVideoFilename)
					videoTitle = title
					videoDescription, _ = rdx.GetLastVal(data.VideoShortDescriptionProperty, videoId)

					//if vct, err := getLocalCaptionTracks(videoId, rdx); err == nil {
					//	videoCaptionTracks = vct
					//}
				}
			}
		}
	}

	if videoUrl == "" || videoTitle == "" {
		videoPage, err := yeti.GetVideoPage(videoId)
		if err != nil {
			return nil, err
		}

		metadataDir, err := pasu.GetAbsDir(paths.Metadata)
		if err != nil {
			return nil, err
		}

		mdRdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
		if err != nil {
			return nil, err
		}

		for p, values := range yeti.ExtractMetadata(videoPage) {
			if err := mdRdx.AddValues(p, videoId, values...); err != nil {
				return nil, err
			}
		}

		vu, err := decode(videoId, videoPage.BestFormat().Url, videoPage.PlayerUrl)
		if err != nil {
			return nil, err
		}

		videoUrl = vu.String()
		videoTitle = videoPage.VideoDetails.Title
		videoDescription = videoPage.VideoDetails.ShortDescription
	}

	lastEndedTime := ""
	if et, ok := rdx.GetLastVal(data.VideoEndedProperty, videoId); ok && et != "" {
		lastEndedTime = et
		titlePrefix := "☑️ "
		if rdx.HasKey(data.VideoSkippedProperty, videoId) {
			titlePrefix = "⏭️ "
		}
		videoTitle = titlePrefix + videoTitle
	}

	return &WatchViewModel{
		VideoId:          videoId,
		VideoUrl:         videoUrl,
		VideoPoster:      videoPoster,
		Server:           playbackType,
		VideoTitle:       videoTitle,
		VideoDescription: videoDescription,
		CurrentTime:      currentTime,
		LastEndedTime:    lastEndedTime,
	}, nil
}

func decode(videoId, urlStr, playerUrl string) (*url.URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	np := q.Get("n")
	if dnp, err := yeti.DecodeParam(http.DefaultClient, videoId, np, playerUrl); err != nil {
		return nil, err
	} else {
		q.Set("n", dnp)
		u.RawQuery = q.Encode()
		return u, nil
	}
}

func getLocalCaptionTracks(videoId string, rdx kvas.ReadableRedux) ([]yt_urls.CaptionTrack, error) {

	if err := rdx.MustHave(
		data.VideoCaptionsNamesProperty,
		data.VideoCaptionsKindsProperty,
		data.VideoCaptionsLanguagesProperty); err != nil {
		return nil, err
	}

	captionsNames, _ := rdx.GetAllValues(data.VideoCaptionsNamesProperty, videoId)
	captionsKinds, _ := rdx.GetAllValues(data.VideoCaptionsKindsProperty, videoId)
	captionsLanguages, _ := rdx.GetAllValues(data.VideoCaptionsLanguagesProperty, videoId)

	cts := make([]yt_urls.CaptionTrack, 0, len(captionsLanguages))
	for i := 0; i < len(captionsLanguages); i++ {

		cn, ck, cl := "", "", captionsLanguages[i]
		if len(captionsNames) >= i {
			cn = captionsNames[i]
		}
		if len(captionsKinds) >= i {
			ck = captionsKinds[i]
		}

		ct := yt_urls.CaptionTrack{
			BaseUrl:      "/captions?v=" + videoId + "&l=" + cl,
			LanguageCode: cl,
			Kind:         ck,
			TrackName:    cn,
		}

		if ct.Kind == "asr" {
			ct.Kind = "subtitles"
		}

		cts = append(cts, ct)
	}

	return cts, nil
}
