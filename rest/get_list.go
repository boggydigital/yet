package rest

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/yet/data"
	"io"
	"net/http"
	"slices"
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
		"<style>")

	writeSharedStyles(sb)

	// list specific styles
	sb.WriteString(
		"a.playlist {display:flex;flex-direction:column;color:deeppink;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:2rem;}" +
			"a.playlist.ended {color:dimgray}" +
			"a.playlist .subtitle {color: inherit}" +
			"ul {list-style:none; padding-inline-start: 0rem}")

	sb.WriteString("</style></head>")
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
				writeVideo(id, rdx, sb, ShowPoster, ShowPublishedDate)
			}
		}
		sb.WriteString("</details>")
	}

	pldq := rdx.Keys(data.PlaylistDownloadQueueProperty)
	newPlaylistVideos := make([]string, 0, len(pldq))

	for _, pl := range pldq {
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
			writeVideo(id, rdx, sb, ShowPoster, ShowPublishedDate)
		}

		if len(newPlaylistVideos) > 0 {
			sb.WriteString("<div class='subtle'>Looking for more? New videos are available in the Playlists</div>")
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
				pt += fmt.Sprintf("<span class='subtitle'>%d new</span>", nvc)
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
			writeVideo(id, rdx, sb, ShowPoster, ShowPublishedDate)
		}
		sb.WriteString("</details>")
	}

	whKeys := rdx.Keys(data.VideoEndedProperty)
	if len(whKeys) > 0 {
		sb.WriteString("<details open><summary><h1>History</h1></summary>")
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

func playlistTitle(playlistId string, rdx kvas.ReadableRedux) string {
	if plt, ok := rdx.GetFirstVal(data.PlaylistTitleProperty, playlistId); ok && plt != "" {

		if plc, ok := rdx.GetFirstVal(data.PlaylistChannelProperty, playlistId); ok && plc != "" && !strings.Contains(plt, plc) {
			return fmt.Sprintf("<span class='playlistTitle'>%s - %s</span>", plc, plt)
		}

		return plt
	}

	return playlistId
}
