package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
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

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch?v&t

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := r.URL.Query().Get("v")
	t := r.URL.Query().Get("t")

	if v == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	// iOS insists on inserting a space on paste
	v = strings.TrimSpace(v)

	videoId := ""
	videoUrl, videoTitle, videoDescription := "", "", ""
	//var videoCaptionTracks []yt_urls.CaptionTrack
	playbackType := "streaming"

	if videoIds, err := yeti.ParseVideoIds(v); err == nil && len(videoIds) > 0 {
		videoId = videoIds[0]
	} else {
		// TODO: check if that local file exists first
		playbackType = "local"
		videoId = v
		videoUrl = "/video?file=" + v
		videoTitle = v
	}

	videoPoster := fmt.Sprintf("/poster?v=%s&q=%s", videoId, yt_urls.ThumbnailQualityMaxRes)

	// if current time is not specified with a query parameter - read it from metadata
	if t == "" {
		if ct, ok := rdx.GetFirstVal(data.VideoProgressProperty, videoId); ok {
			t = ct
		}
	}

	if title, ok := rdx.GetFirstVal(data.VideoTitleProperty, videoId); ok && title != "" {
		if channel, ok := rdx.GetFirstVal(data.VideoOwnerChannelNameProperty, videoId); ok && channel != "" {
			localVideoFilename := yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
			if absVideosDir, err := paths.GetAbsDir(paths.Videos); err == nil {
				absLocalVideoFilename := filepath.Join(absVideosDir, localVideoFilename)
				if _, err := os.Stat(absLocalVideoFilename); err == nil {
					playbackType = "local"
					videoUrl = "/video?file=" + url.QueryEscape(localVideoFilename)
					videoTitle = title
					videoDescription, _ = rdx.GetFirstVal(data.VideoShortDescriptionProperty, videoId)

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metadataDir, err := paths.GetAbsDir(paths.Metadata)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mdRdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for p, values := range yeti.ExtractMetadata(videoPage) {
			if err := mdRdx.AddValues(p, videoId, values...); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		fs := videoPage.Formats()
		var f yt_urls.Format
		for _, ff := range fs {
			if ff.Quality == "hd720" {
				f = ff
			}
		}

		vu, err := decode(videoId, f.Url, videoPage.PlayerUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		videoUrl = vu.String()
		videoTitle = videoPage.VideoDetails.Title
		videoDescription = videoPage.VideoDetails.ShortDescription
	}

	lastEndedTime := ""
	if et, ok := rdx.GetFirstVal(data.VideoEndedProperty, videoId); ok && et != "" {
		lastEndedTime = et
		videoTitle = "☑️ " + videoTitle
	}

	w.Header().Set("Content-Type", "text/html")

	wvm := &WatchViewModel{
		VideoId:          videoId,
		VideoUrl:         videoUrl,
		VideoPoster:      videoPoster,
		Server:           playbackType,
		VideoTitle:       videoTitle,
		VideoDescription: videoDescription,
		CurrentTime:      t,
		LastEndedTime:    lastEndedTime,
	}

	if err := tmpl.ExecuteTemplate(w, "watch", wvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
