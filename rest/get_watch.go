package rest

import (
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
	"time"
)

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch?videoId

	v := r.URL.Query().Get("v")

	if v == "" {
		http.Redirect(w, r, "/new", http.StatusPermanentRedirect)
		return
	}

	// iOS insists on inserting a space on paste
	v = strings.TrimSpace(v)

	videoIds, err := yeti.ArgsToVideoIds(http.DefaultClient, false, v)
	if err != nil {
		http.Error(w, "error extracting videoId", http.StatusBadRequest)
		return
	}

	videoId := ""
	if len(videoIds) > 0 {
		videoId = videoIds[0]
	}

	videoUrl, videoPoster, videoTitle, videoDescription := "", "", "", ""

	if videoId == "" {
		videoId = v
		videoUrl = "/video?file=" + v
		videoTitle = v
	}

	absMetadataDir, err := paths.GetAbsDir(paths.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rxa, err := kvas.ConnectReduxAssets(absMetadataDir, data.AllProperties()...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if title, ok := rxa.GetFirstVal(data.VideoTitleProperty, videoId); ok && title != "" {
		if channel, ok := rxa.GetFirstVal(data.VideoOwnerChannelNameProperty, videoId); ok && channel != "" {
			localVideoFilename := yeti.ChannelTitleVideoIdFilename(channel, title, videoId)
			if absVideosDir, err := paths.GetAbsDir(paths.Videos); err == nil {
				absLocalVideoFilename := filepath.Join(absVideosDir, localVideoFilename)
				if _, err := os.Stat(absLocalVideoFilename); err == nil {
					videoUrl = "/video?file=" + url.QueryEscape(localVideoFilename)
					videoPoster = "/poster?v=" + videoId + "&q=maxresdefault"
					videoTitle = title
					videoDescription, _ = rxa.GetFirstVal(data.VideoShortDescriptionProperty, videoId)
				}
			}
		}
	}

	currentTime := ""
	if ct, ok := rxa.GetFirstVal(data.VideoProgressProperty, videoId); ok && ct != "" {
		currentTime = ct
	}

	lastEndedTime := ""
	if et, ok := rxa.GetFirstVal(data.VideoEndedProperty, videoId); ok && et != "" {
		lastEndedTime = et
	}

	if videoUrl == "" || videoTitle == "" {
		videoPage, playerUrl, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fs := videoPage.Formats()
		var f yt_urls.Format
		for _, ff := range fs {
			if ff.Quality == "hd720" {
				f = ff
			}
		}

		vu, err := decode(f.Url, playerUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		videoUrl = vu.String()
		if len(videoPage.Microformat.PlayerMicroformatRenderer.Thumbnail.Thumbnails) > 0 {
			videoPoster = videoPage.Microformat.PlayerMicroformatRenderer.Thumbnail.Thumbnails[0].Url
		}
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
		"<style>" +
		"body {background: black; color: white;font-family:sans-serif; margin: 1rem;} " +
		"video {width: 100%; height: 100%; object-fit: cover} " +
		"summary.videoTitle {font-size: 2rem; margin: 0.5rem;cursor:pointer}" +
		".videoDescription {margin: 1rem; line-height: 1.2;}" +
		".lastEnded {margin: 0.5rem; font-size: 75%;}" +
		"</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<video controls='controls' preload='metadata' poster='" + videoPoster + "'>")
	sb.WriteString("<source src='" + videoUrl + "' />")
	sb.WriteString("</video>")

	if lastEndedTime != "" {
		if lt, err := time.Parse(http.TimeFormat, lastEndedTime); err == nil {
			sb.WriteString("<div class='lastEnded'><span>Last ended: ")
			sb.WriteString("<time>" + lt.Local().Format(time.RFC1123) + "</time></div>")
		}
	}

	sb.WriteString("<details>")
	sb.WriteString("<summary class='videoTitle'>" + videoTitle + "</summary>")
	sb.WriteString("<div class='videoDescription'>" + videoDescription + "</div>")
	sb.WriteString("</details>")

	sb.WriteString("<script>" +
		"let video = document.getElementsByTagName('video')[0];" +
		"</script>")

	// only continue the videos that have not been watched
	if currentTime != "" && lastEndedTime == "" {
		sb.WriteString("<script>video.currentTime = " + currentTime + ";</script>")
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
		"				videoId: '" + videoId + "'," +
		"				currentTime: video.currentTime.toString()})" +
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
		"		body: JSON.stringify({videoId: '" + videoId + "'})" +
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

func decode(urlStr, playerUrl string) (*url.URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	np := q.Get("n")
	if dnp, err := yeti.DecodeParam(http.DefaultClient, np, playerUrl); err != nil {
		return nil, err
	} else {
		q.Set("n", dnp)
		u.RawQuery = q.Encode()
		return u, nil
	}
}
