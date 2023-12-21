package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	maxPlaylistVideosWatchlist = 3
)

func GetList(w http.ResponseWriter, r *http.Request) {

	// GET /list

	var err error
	rdx, err = rdx.RefreshWriter()
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
		"<title>🔻 Watch list</title>" +
		"<style>" +
		"body {background: black; color: white;font-family:sans-serif; margin: 1rem;} " +
		"a.video {display:block;color:white;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"a.video img {border-radius:0.25rem;width:200px;aspect-ratio:16/9;background:dimgray}" +
		"a.video span {font-size:1rem}" +
		"a.video.ended {filter:grayscale(1.0)}" +
		"a.highlight {color:gold; margin-block:2rem}" +
		"details {margin-block:2rem; content-visibility: auto}" +
		"summary {margin-block-end: 2rem}" +
		"summary h1 {display: inline; cursor: pointer; margin-inline-start: 0.5rem;color:turquoise}" +
		"a.playlist {display:block;color:deeppink;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:0.5rem;margin-block-end: 1rem}" +
		"a.playlist.ended {color:dimgray}" +
		"div.subtle {color: dimgray}" +
		"ul {list-style:none; padding-inline-start: 1rem}" +
		"</style></head>")
	sb.WriteString("<body>")

	sb.WriteString("<a class='video highlight' href='/new'>Watch new</a>")

	// continue watching
	// videos watchlist
	// videos download queue

	cwKeys := rdx.Keys(data.VideoProgressProperty)
	if len(cwKeys) > 0 {
		cwKeys, err = rdx.Sort(cwKeys, false, data.VideoTitleProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sb.WriteString("<details open><summary><h1>Continue</h1></summary>")
		for _, id := range cwKeys {
			if ended, ok := rdx.GetFirstVal(data.VideoEndedProperty, id); !ok || ended == "" {
				writeVideo(id, true, rdx, sb)
			}
		}
		sb.WriteString("</details>")
	}

	plnv := rdx.Keys(data.PlaylistNewVideosProperty)
	newPlaylistVideos := make([]string, 0, len(plnv))

	for _, pl := range plnv {
		if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, pl); ok {
			newPlaylistVideos = append(newPlaylistVideos, nv...)
		}
	}

	wlKeys := rdx.Keys(data.VideosWatchlistProperty)
	if len(wlKeys) > 0 {

		wlKeys, err = rdx.Sort(wlKeys, false, data.VideoTitleProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sb.WriteString("<details><summary><h1>Watchlist</h1></summary>")

		for _, id := range wlKeys {
			if slices.Contains(newPlaylistVideos, id) {
				continue
			}
			if le, ok := rdx.GetFirstVal(data.VideoEndedProperty, id); ok && le != "" {
				continue
			}
			if ct, ok := rdx.GetFirstVal(data.VideoProgressProperty, id); ok || ct != "" {
				continue
			}
			writeVideo(id, true, rdx, sb)
		}

		sb.WriteString("</details>")
	}

	plKeys := rdx.Keys(data.PlaylistWatchlistProperty)
	if len(plKeys) > 0 {

		plKeys, err = rdx.Sort(plKeys, false, data.PlaylistTitleProperty, data.PlaylistChannelProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		openOrClose := ""
		if len(plKeys) < 7 {
			openOrClose = "open"
		}

		sb.WriteString("<details " + openOrClose + "><summary><h1>Playlists</h1></summary>")
		sb.WriteString("<ul>")
		for _, id := range plKeys {

			nvc := 0

			if nv, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, id); ok {
				nvc = len(nv)
			}

			pt := playlistTitle(id, rdx)
			if nvc > 0 {
				pt += " (" + strconv.Itoa(nvc) + " new)"
			}

			pc := "playlist"
			if nvc == 0 {
				pc += " ended"
			}

			sb.WriteString("<li><a class='" + pc + "' href='/playlist?list=" + id + "'>" + pt + "</a></li>")

		}
		sb.WriteString("</ul>")
		sb.WriteString("</details>")
	}

	dqKeys := rdx.Keys(data.VideosDownloadQueueProperty)
	if len(dqKeys) > 0 {

		dqKeys, err = rdx.Sort(dqKeys, false, data.VideoTitleProperty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sb.WriteString("<details><summary><h1>Downloads</h1></summary>")
		for _, id := range dqKeys {
			writeVideo(id, true, rdx, sb)
		}
		sb.WriteString("</details>")
	}

	whKeys := rdx.Keys(data.VideoEndedProperty)
	if len(whKeys) > 0 {
		sb.WriteString("<details><summary><h1>History</h1></summary>")
		sb.WriteString("<a class='video' href='/history'>See all watch history</a>")
		sb.WriteString("</details>")
	}

	sb.WriteString("</body>")
	sb.WriteString("</html>")

	if _, err := io.WriteString(w, sb.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func writeVideo(videoId string, showImage bool, rdx kvas.ReadableRedux, sb *strings.Builder) {

	videoTitle := videoId
	if title, ok := rdx.GetFirstVal(data.VideoTitleProperty, videoId); ok && title != "" {
		videoTitle = title
	}

	ended := false
	if et, ok := rdx.GetFirstVal(data.VideoEndedProperty, videoId); ok && et != "" {
		ended = true
	}

	//progress := false
	//if pt, ok := rdx.GetFirstVal(data.VideoProgressProperty, videoId); ok && pt != "" {
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

	imageContent := ""
	if showImage {
		imageContent = "<img src='/poster?v=" + videoId + "&q=mqdefault' loading='lazy'/>"
	}

	sb.WriteString("<a class='" + class + "' href='" + videoUrl + "'>" +
		imageContent +
		"<br/>" +
		"<span>" + videoTitle + "</span>" +
		"</a>")

}

func playlistTitle(playlistId string, rdx kvas.ReadableRedux) string {
	if plt, ok := rdx.GetFirstVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {

		if plc, ok := rdx.GetFirstVal(data.PlaylistChannelProperty, playlistId); ok && plc != "" && !strings.Contains(plt, plc) {
			return fmt.Sprintf("%s - %s", plc, plt)
		}

		return plt
	}

	return playlistId
}
