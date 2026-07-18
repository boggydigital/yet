package rest

import (
	_ "embed"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boggydigital/camino"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/calc"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

//go:embed "scripts/watch.js"
var scriptWatch string

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch/{videoId}?t

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	t := r.URL.Query().Get("t")

	if t == "" {

		var rem, dur int64

		if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" {
			var duri int64
			if duri, err = strconv.ParseInt(durs, 10, 64); err == nil {
				dur = duri
			}
		}

		var ct int64
		if cts, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && cts != "" {
			var cti int64
			if cti, err = strconv.ParseInt(cts, 10, 64); err == nil {
				ct = cti
			}
		}
		rem = dur - ct

		t = strconv.FormatInt(dur-rem, 10)
	}

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	// iOS insists on inserting a space on paste
	videoId = strings.TrimSpace(videoId)

	var videoIds []string
	if videoIds, err = yeti.ParseVideoIds(videoId); err == nil && len(videoIds) > 0 {
		videoId = videoIds[0]
	}

	var videoTitle string
	if vt, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && vt != "" {
		videoTitle = vt
	}

	root, body := strom.RootBody(videoTitle, atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(
		navButton("Home", "/"),
		navButton("Paste", "/paste"))

	var absLocalVideoFilename string

	if title, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && title != "" {
		if channelId, sure := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, videoId); sure && channelId != "" {
			videosDir := camino.GetAbs(data.Videos)
			absLocalVideoFilename = filepath.Join(videosDir, yeti.RelLocalVideoFilename(channelId, title, videoId))
		}
	}

	if absLocalVideoFilename == "" {
		absLocalVideoFilename, err = yeti.LocateLocalVideo(videoId)
		if os.IsNotExist(err) {
			// do nothing
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	videoPosterUrl := "/poster?v=" + videoId + "&q=" + youtube_urls.ThumbnailQualityMaxRes.String()

	var mediaElement strom.Element

	if absLocalVideoFilename != "" {
		if _, err = os.Stat(absLocalVideoFilename); err == nil {
			videosDir := camino.GetAbs(data.Videos)

			var relLocalVideoFilename string
			relLocalVideoFilename, err = filepath.Rel(videosDir, absLocalVideoFilename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			videoUrl := "/video?file=" + url.QueryEscape(relLocalVideoFilename)
			//videoDescription, _ = rdx.GetLastVal(data.VideoShortDescriptionProperty, videoId)

			mediaElement = strom.Create("video").
				SetAttribute("src", videoUrl).
				SetAttribute("poster", videoPosterUrl).
				SetAttribute("controls", "controls").
				SetAttribute("preload", "none")

		} else {
			topRow.Append(navButton("Download", path.Join("/download_video", videoId)))
			mediaElement = strom.Create("img").SetAttribute("src", videoPosterUrl)
		}
	}

	topRow.Append(strom.CreateText("h2", videoTitle))

	mediaElement.SetStyle(map[string]string{
		"max-width":     calc.Mult(sizes.XXXLarge, 4),
		"border-radius": sizes.XSmall,
	})

	body.Append(mediaElement)

	if channelId, ok := rdx.GetLastVal(data.VideoExternalChannelIdProperty, videoId); ok && channelId != "" {
		body.Append(channelTile(channelId, rdx))
	}

	videoNavButtonsRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...)
	body.Append(videoNavButtonsRow)

	videoNavButtonsRow.Append(
		navButton("Manage", "/manage_video?v="+videoId),
		actionButton("Seen enough", "/end/"+videoId+"/seen-enough"),
		actionButton("Skip", "/end/"+videoId+"/skipped"),
	)

	if absLocalVideoFilename == "" {
		videoNavButtonsRow.Append(actionButton("Queue download", "/queue_download/"+videoId))
	}

	if vd, ok := rdx.GetLastVal(data.VideoShortDescriptionProperty, videoId); ok && vd != "" {
		body.Append(
			strom.CreateText("h3", "Description"),
			strom.CreateText("pre", vd).SetStyle(map[string]string{
				"white-space": "pre-wrap",
				"word-break":  "break-word",
				"color":       colors.Gray,
				"max-width":   calc.Mult(sizes.XXXLarge, 4),
			}))
	}

	// must be a new string per video otherwise global will be rewritten for all
	videoScriptWatch := strings.Replace(scriptWatch, "{currentTime}", t, -1)
	videoScriptWatch = strings.Replace(videoScriptWatch, "{videoId}", videoId, -1)

	if err = strom.WriteResponse(w, root, []byte(videoScriptWatch)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
