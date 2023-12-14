package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"io"
	"net/http"
	"strings"
)

func GetList(w http.ResponseWriter, r *http.Request) {

	// GET /list

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
		"a.video {display:block;color:white;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"a.video img {border-radius:0.25rem;width:200px}" +
		"a.video span {font-size:1rem}" +
		"a.video.ended {filter:grayscale(1.0)}" +
		"a.highlight {color:gold; margin-block:2rem}" +
		"details {margin-block:0.5rem; content-visibility: auto}" +
		"summary h1 {display: inline; cursor: pointer}" +
		"a.playlist {display:block;color:deeppink;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<a class='video highlight' href='/new'>Something else</a>")

	// continue watching
	// videos watchlist
	// videos download queue

	cwKeys := rxa.Keys(data.VideoProgressProperty)
	if len(cwKeys) > 0 {
		sb.WriteString("<details open><summary><h1>Continue</h1></summary>")
		for _, id := range cwKeys {
			if ended, ok := rxa.GetFirstVal(data.VideoEndedProperty, id); !ok || ended == "" {
				writeVideo(id, rxa, sb)
			}
		}
		sb.WriteString("</details>")
	}

	wlKeys := rxa.Keys(data.VideosWatchlistProperty)
	if len(wlKeys) > 0 {
		sb.WriteString("<details><summary><h1>Watchlist</h1></summary>")
		for _, id := range wlKeys {
			if le, ok := rxa.GetFirstVal(data.VideoEndedProperty, id); ok && le != "" {
				continue
			}
			if ct, ok := rxa.GetFirstVal(data.VideoProgressProperty, id); ok || ct != "" {
				continue
			}
			writeVideo(id, rxa, sb)
		}
		sb.WriteString("</details>")
	}

	plKeys := rxa.Keys(data.PlaylistWatchlistProperty)
	if len(plKeys) > 0 {
		sb.WriteString("<details open><summary><h1>Playlists</h1></summary>")
		sb.WriteString("<ul>")
		for _, id := range plKeys {
			if plt, ok := rxa.GetFirstVal(data.PlaylistTitleProperty, id); ok && plt != "" {

				if plc, ok := rxa.GetFirstVal(data.PlaylistChannelProperty, id); ok && plc != "" && !strings.Contains(plt, plc) {
					plt = fmt.Sprintf("%s - %s", plc, plt)
				}

				sb.WriteString("<li><a class='playlist' href='/playlist?id=" + id + "'>" +
					plt +
					"</a></li>")
			}
		}
		sb.WriteString("</ul>")
		sb.WriteString("</details>")
	}

	dqKeys := rxa.Keys(data.VideosDownloadQueueProperty)
	if len(dqKeys) > 0 {
		sb.WriteString("<details><summary><h1>Download queue</h1></summary>")
		for _, id := range dqKeys {
			writeVideo(id, rxa, sb)
		}
		sb.WriteString("</details>")
	}

	whKeys := rxa.Keys(data.VideoEndedProperty)
	if len(whKeys) > 0 {
		sb.WriteString("<details><summary><h1>Watch history</h1></summary>")
		for _, id := range whKeys {
			writeVideo(id, rxa, sb)
		}
		sb.WriteString("</details>")
	}

	sb.WriteString("<a class='video highlight' href='/new'>Something else</a>")

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func writeVideo(videoId string, rxa kvas.ReduxAssets, sb *strings.Builder) {

	videoTitle := videoId
	if title, ok := rxa.GetFirstVal(data.VideoTitleProperty, videoId); ok && title != "" {
		videoTitle = title
	}

	ended := false
	if et, ok := rxa.GetFirstVal(data.VideoEndedProperty, videoId); ok && et != "" {
		ended = true
	}

	//progress := false
	//if pt, ok := rxa.GetFirstVal(data.VideoProgressProperty, videoId); ok && pt != "" {
	//	progress = true
	//}

	videoUrl := "/watch?"
	if videoId != "" {
		videoUrl += "v=" + videoId
	}

	class := "video"
	if ended {
		videoTitle = "☑️ " + videoTitle
		class += " ended"
	}

	sb.WriteString("<a class='" + class + "' href='" + videoUrl + "'>" +
		"<img src='/poster?v=" + videoId + "&q=hqdefault' />" +
		"<br/>" +
		"<span>" + videoTitle + "</span>" +
		"</a>")

}
