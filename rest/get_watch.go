package rest

import (
	_ "embed"
	"iter"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/boggydigital/camino"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/styles"
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

type videoTimings struct {
	duration    int64
	remaining   int64
	currentTime int64
}

func GetWatch(w http.ResponseWriter, r *http.Request) {

	// GET /watch/{videoId}?t

	var err error
	rdx, err = rdx.RefreshWriter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoId := r.PathValue("videoId")

	if videoId == "" {
		http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
		return
	}

	vt := getVideoTimings(videoId, rdx)

	t := r.URL.Query().Get("t")
	if t == "" {
		t = strconv.FormatInt(vt.currentTime, 10)
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
	mediaElement = strom.Create("img").SetAttribute("src", videoPosterUrl)

	videoNavButtonsRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		AddAtom(atoms.AlignItemsCenter)
	videoNavButtonsRow.Append(
		navButton("Manage", path.Join("/manage_video", videoId)),
		navButton("Seen enough", path.Join("/end", videoId, "seen-enough")),
		navButton("Skip", path.Join("/end", videoId, "skipped")),
	)

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
			addQueueDownloadAction(videoId, videoNavButtonsRow, rdx)
		}
	} else {
		topRow.Append(navButton("Download", path.Join("/download_video", videoId)))
		addQueueDownloadAction(videoId, videoNavButtonsRow, rdx)
	}

	topRow.Append(strom.CreateText("h2", videoTitle))

	mediaElement.SetStyle(
		styles.Decl("max-width", calc.Mult(sizes.XXXLarge, 4)),
		styles.Decl("border-radius", sizes.XSmall))

	body.Append(mediaElement)

	pct := new(playlistChannelTile{videoId: videoId, rdx: rdx})
	body.Append(strom.OnDemand(pct.getPlaylistChannelTile))

	body.Append(videoNavButtonsRow)

	body.Append(strom.CreateText("h3", "Description"))
	body.Append(navButton("Refresh", path.Join("/refresh_video", videoId)))

	if vd, ok := rdx.GetLastVal(data.VideoShortDescriptionProperty, videoId); ok && vd != "" {
		body.Append(
			strom.CreateText("pre", vd).
				SetStyle(
					"white-space:pre-wrap",
					"word-break:break-word",
					styles.Decl("color", colors.Gray),
					styles.Decl("max-width", calc.Mult(sizes.XXXLarge, 4))))
	}

	// must be a new string per video otherwise global will be rewritten for all
	videoScriptWatch := strings.Replace(scriptWatch, "{currentTime}", t, -1)
	videoScriptWatch = strings.Replace(videoScriptWatch, "{videoId}", videoId, -1)

	if err = strom.WriteResponse(w, root, []byte(videoScriptWatch)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type playlistChannelTile struct {
	videoId string
	rdx     redux.Readable
}

func (pct *playlistChannelTile) getPlaylistChannelTile() iter.Seq[strom.Element] {
	return func(yield func(element strom.Element) bool) {

		allPlaylistsWithVideo := rdx.MatchAsset(data.PlaylistVideosProperty, []string{pct.videoId}, nil)
		var playlistId string
		for pid := range allPlaylistsWithVideo {
			if rdx.HasKey(data.PlaylistAutoRefreshProperty, pid) {
				playlistId = pid
				break
			}
		}

		if playlistId == "" {
			for pid := range allPlaylistsWithVideo {
				playlistId = pid
				break
			}
		}

		if playlistId != "" {
			if rdx.HasKey(data.PlaylistAutoRefreshProperty, playlistId) {
				if !yield(playlistTile(playlistId, pct.rdx)) {
					return
				}
				return
			}
		}

		if channelId, ok := rdx.GetLastVal(data.VideoExternalChannelIdProperty, pct.videoId); ok && channelId != "" {
			if !yield(channelTile(channelId, rdx)) {
				return
			}
		}
	}
}

func getVideoTimings(videoId string, rdx redux.Readable) *videoTimings {

	vt := new(videoTimings)

	if durs, sure := rdx.GetLastVal(data.VideoDurationProperty, videoId); sure && durs != "" {
		if duri, err := strconv.ParseInt(durs, 10, 64); err == nil {
			vt.duration = duri
		}
	}

	if cts, ok := rdx.GetLastVal(data.VideoProgressProperty, videoId); ok && cts != "" {
		if cti, err := strconv.ParseInt(cts, 10, 64); err == nil {
			vt.currentTime = cti
		}
	}

	vt.remaining = vt.duration - vt.currentTime

	return vt
}

func addQueueDownloadAction(videoId string, container strom.Element, rdx redux.Readable) {
	if dqs, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, videoId); ok && dqs != "" {
		downloadQueued := true
		if dcs, sure := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); sure {
			if dqs < dcs {
				downloadQueued = false
			}
		}
		if downloadQueued {
			dqDateTime := dqs
			if dqt, err := time.Parse(time.RFC3339, dqs); err == nil {
				dqDateTime = dqt.Format(time.DateTime)
			}

			container.Append(strom.CreateText("span", "Download queued: "+dqDateTime).
				SetStyle(styles.Decl("color", colors.Gray)))
		}
	} else {
		container.Append(navButton("Queue download", path.Join("/queue_download", videoId)))
	}

}
