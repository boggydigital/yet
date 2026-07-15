package rest

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boggydigital/camino"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/youtube_urls"
)

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

	root := strom.Page(videoTitle)

	var body strom.Element
	for body = range root.GetElementsByTagName("body") {
		break
	}

	body.AddClass("d-f", "fd-c", "rg-l")

	topNavButtons := strom.Create("ul", "d-f", "cg-n", "rg-n").
		SetStyle(map[string]string{
			"flex-flow": "row wrap",
		})

	body.Append(topNavButtons)

	topNavButtons.Append(
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

	videoPosterUrl := "/poster?v=" + videoId + "&q" + youtube_urls.ThumbnailQualityMaxRes.String()

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

			body.Append(strom.Create("video", "br-s").
				SetAttribute("src", videoUrl).
				SetAttribute("poster", videoPosterUrl).
				SetAttribute("controls", "controls").
				SetAttribute("preload", "none").
				SetStyle(map[string]string{
					"max-width": "calc(4 * " + vars.Size(vars.SizeXXXLarge) + ")"}))

		} else {
			body.Append(strom.Create("img", "br-s").
				SetAttribute("src", videoPosterUrl).
				SetStyle(map[string]string{
					"max-width": "calc(4 * " + vars.Size(vars.SizeXXXLarge) + ")"}))
		}
	}

	body.Append(strom.CreateText("h2", videoTitle))

	if channelId, ok := rdx.GetLastVal(data.VideoExternalChannelIdProperty, videoId); ok && channelId != "" {
		body.Append(channelTile(channelId, rdx))
	}

	videoNavButtonsRow := strom.Create("ul", "d-f", "cg-n", "rg-n").
		SetStyle(map[string]string{
			"flex-flow": "row wrap",
		})
	body.Append(videoNavButtonsRow)

	videoNavButtonsRow.Append(
		navButton("Manage video", "/manage_video?v="+videoId),
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
				"color":       vars.Color(vars.ColorGray),
			}))
	}

	body.Append(strom.CreateText("script", "let video = document.getElementsByTagName('video')[0];"))
	body.Append(strom.CreateText("script", "video.currentTime = "+t+";"))
	body.Append(strom.CreateText("script",
		"let lastProgressUpdate = new Date();\n        video.addEventListener('timeupdate', (e) => {\n            let now = new Date();\n            let elapsed = now - lastProgressUpdate;\n            if (elapsed > 5000) {\n                fetch('/progress', {\n                    method: 'post',\n                    headers: {\n                        'Content-Type': 'application/json'},\n                    body: JSON.stringify({\n                        v: '"+videoId+"',\n                        t: video.currentTime.toString()})\n                }).then((resp) => { if (resp && !resp.ok) {\n                    console.log(resp)}\n                });\n                lastProgressUpdate = now;\n            }});"))
	body.Append(strom.CreateText("script", "video.addEventListener('ended', (e) => {\n        fetch('/end/{{.VideoId}}/completed', {\n                method: 'get',\n            }).then((resp) => { if (resp && !resp.ok) {\n                console.log(resp)}\n            });\n        if (prg) {prg.value = prg.max}\n        });"))
	body.Append(strom.CreateText("script", "document.body.addEventListener('keydown', (e) => {\n            switch (e.keyCode) {\n        // ArrowRight\n                case 39:\n                e.preventDefault();\n                video.currentTime += 15;\n                break;\n        // ArrowLeft\n                case 37:\n                e.preventDefault();\n                video.currentTime -= 15;\n                break;\n        // Space\n                case 32:\n                e.preventDefault();\n                video.paused ? video.play() : video.pause();\n                break;\n            };\n            });"))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
