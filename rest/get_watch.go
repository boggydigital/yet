package rest

import (
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch?v

	v := r.URL.Query().Get("v")
	if v == "" {
		http.Error(w, "missing video-id (v)", http.StatusBadRequest)
		return
	}

	videoPage, playerUrl, err := yt_urls.GetVideoPage(http.DefaultClient, v)
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

	posterUrl := yt_urls.ThumbnailUrl(v, yt_urls.ThumbnailQualityHQ).String()

	w.Header().Set("Content-Type", "text/html")

	sb := &strings.Builder{}
	sb.WriteString("<!doctype html>")
	sb.WriteString("<html>")
	sb.WriteString("<head><style>" +
		"body {background: black; color: white;font-family:sans-serif} " +
		"video {width: 100%; height: 100%; object-fit: cover} " +
		".videoTitle {font-size: 2rem; margin: 0.5rem;}" +
		"</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<video id='video" + v + "' controls='controls' poster='" + posterUrl + "' preload='metadata'>")
	sb.WriteString("<source src='" + vu.String() + "' type='" + f.MIMEType + "' />")
	sb.WriteString("</video>")

	sb.WriteString("<div class='videoTitle'>" + videoPage.Microformat.PlayerMicroformatRenderer.Title.SimpleText + "</div>")
	sb.WriteString("<div class='viewCount'>" + videoPage.Microformat.PlayerMicroformatRenderer.ViewCount + "</div>")

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
