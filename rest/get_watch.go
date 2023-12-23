package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch?v&t

	var err error
	rdx, err = rdx.RefreshWriter()
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

	lastEndedTime := ""
	if et, ok := rdx.GetFirstVal(data.VideoEndedProperty, videoId); ok && et != "" {
		lastEndedTime = et
	}

	if videoUrl == "" || videoTitle == "" {
		videoPage, err := yeti.GetVideoPage(videoId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for p, v := range yeti.ExtractMetadata(videoPage) {
			if err := rdx.AddValues(p, videoId, v...); err != nil {
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

	w.Header().Set("Content-Type", "text/html")

	sb := &strings.Builder{}
	sb.WriteString("<!doctype html>")
	sb.WriteString("<html>")
	sb.WriteString("<head>" +
		"<meta charset='UTF-8'>" +
		"<link rel='icon' href='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🔻</text></svg>' type='image/svg+xml'/>" +
		"<meta name='viewport' content='width=device-width, initial-scale=1.0'>" +
		"<meta name='color-scheme' content='dark light'>" +
		"<title>🔻 " + videoTitle + "</title>" +
		"<style>" +
		"body {background: black; color: white;font-family:sans-serif; margin: 1rem;} " +
		"video {width: 100%; height: 100%; aspect-ratio:16/9} " +
		"h1 {margin-block: 2rem}" +
		"details {margin-block: 1rem}" +
		"details summary {cursor:pointer}" +
		"summary::after {content: '\u2026';flex-shrink: 0}" +
		"summary::-webkit-details-marker {display: none}" +
		"h2 {display: inline}" +
		".videoDescription {margin-block:1rem}" +
		"div.subtle {color:dimgray}" +
		"</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<video controls='controls' preload='metadata' poster='" + videoPoster + "'>")
	sb.WriteString("<source src='" + videoUrl + "' />")
	//for _, vct := range videoCaptionTracks {
	//	sb.WriteString("<track " +
	//		"kind='" + vct.Kind + "' " +
	//		"label='" + vct.TrackName + "' " +
	//		"srclang='" + vct.LanguageCode + "' " +
	//		"src='" + vct.BaseUrl + "'/>")
	//}
	sb.WriteString("</video>")

	server := "origin"
	if playbackType == "local" {
		server = "yet"
	}

	sb.WriteString(fmt.Sprintf("<div class='subtle'>This video will play from the %s server</div>", server))

	if lastEndedTime != "" {
		videoTitle = "☑️ " + videoTitle
	}

	sb.WriteString("<h1 class='videoTitle'>" + videoTitle + "</h1>")

	if videoDescription != "" {
		sb.WriteString("<details>")
		sb.WriteString("<summary><h2>Description</h2></summary>")
		sb.WriteString("<div class='videoDescription'>" + videoDescription + "</div>")
		sb.WriteString("</details>")
	}

	//sb.WriteString("<details>")
	//sb.WriteString("<summary><h1>Tools</h1></summary>")
	//if lastEndedTime != "" {
	//	if lt, err := time.Parse(http.TimeFormat, lastEndedTime); err == nil {
	//		sb.WriteString("<div class='lastEnded'><span>Last watched: ")
	//		sb.WriteString("<time>" + lt.Local().Format(time.RFC1123) + "</time></div>")
	//	}
	//}
	//sb.WriteString("<span class='playbackType'>Video source: " + playbackType + "</span>")
	//sb.WriteString("</details>")

	sb.WriteString("<script>" +
		"let video = document.getElementsByTagName('video')[0];" +
		"</script>")

	// only continue the videos that have not been watched
	if t != "" && lastEndedTime == "" {
		sb.WriteString("<script>video.currentTime = " + t + ";</script>")
	}

	sb.WriteString("<script>" +
		"let lastProgressUpdate = new Date();" +
		"video.addEventListener('timeupdate', (e) => {" +
		"	let now = new Date();" +
		"	let elapsed = now - lastProgressUpdate;" +
		"	if (elapsed > 5000) {" +
		"		fetch('/progress', {" +
		"			method: 'post'," +
		"			headers: {" +
		"				'Content-Type': 'application/json'}," +
		"			body: JSON.stringify({" +
		"				v: '" + videoId + "'," +
		"				t: video.currentTime.toString()})" +
		"		}).then((resp) => { if (resp && !resp.ok) {" +
		"			console.log(resp)}" +
		"		});" +
		"		lastProgressUpdate = now;" +
		"	}});" +
		"</script>")

	sb.WriteString("<script>" +
		"video.addEventListener('ended', (e) => {" +
		"fetch('/ended', {" +
		"		method: 'post'," +
		"		headers: {" +
		"			'Content-Type': 'application/json'}," +
		"		body: JSON.stringify({v: '" + videoId + "'})" +
		"	}).then((resp) => { if (resp && !resp.ok) {" +
		"		console.log(resp)}" +
		"	});});" +
		"</script>")

	sb.WriteString("<script>" +
		"document.body.addEventListener('keydown', (e) => {" +
		"	switch (e.keyCode) {" +
		// ArrowRight
		"		case 39:" +
		"		e.preventDefault();" +
		"		video.currentTime += 15;" +
		"		break;" +
		// ArrowLeft
		"		case 37:" +
		"		e.preventDefault();" +
		"		video.currentTime -= 15;" +
		"		break;" +
		// Space
		"		case 32:" +
		"		e.preventDefault();" +
		"		video.paused ? video.play() : video.pause();" +
		"		break;" +
		"	};" +
		"	});" +
		"</script>")

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
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
