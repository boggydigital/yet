package rest

import (
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/font_sizes"
	"github.com/boggydigital/strom/vars/sizes"
	"github.com/boggydigital/yet/data"
)

func GetManageVideo(w http.ResponseWriter, r *http.Request) {

	// GET /manage_video/{videoId}

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

	root, body := strom.RootBody("Manage video", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))
	topRow.Append(strom.CreateText("h2", "Manage video"))

	if videoTitle, ok := rdx.GetLastVal(data.VideoTitleProperty, videoId); ok && videoTitle != "" {
		body.Append(strom.CreateText("h3", videoTitle))
	}

	body.Append(strom.Create("span").Append(
		strom.CreateText("span", "VideoId: ").SetStyle("color:"+colors.Gray),
		strom.CreateText("span", videoId)))

	originRow := strom.Create("ul", atoms.FlexRowWrap(sizes.Small)...).
		AddAtom(atoms.AlignItemsCenter)
	body.Append(originRow)

	var originTitle, originUrl string

	switch vt.currentTime {
	case 0:
		originTitle = "Watch at origin"
		originUrl = "https://www.youtube.com/watch?v=" + videoId
	default:
		originTitle = "Continue watching at origin"
		originUrl = "https://www.youtube.com/watch?v=" + videoId + "&t=" + strconv.FormatInt(vt.currentTime, 10)
	}

	originRow.Append(navButton(originTitle, originUrl))

	form := strom.Create("form", atoms.FlexColWrap(sizes.Normal)...).
		SetAttribute("id", "manage-video").
		SetAttribute("method", "get").
		SetAttribute("action", path.Join("/update_video", videoId))
	body.Append(form)

	favorite := rdx.HasKey(data.VideoFavoriteProperty, videoId)
	form.Append(switchTitleSubtitle(favorite, "favorite", "Favorite", "On: Prevent video cleanup. Off: cleanup after ended."))

	progress := vt.currentTime > 0
	form.Append(switchTitleSubtitle(progress, "progress", "Progress", "Cannot be set here, will be set during video playback. Off: clear progress."))

	ended := rdx.HasKey(data.VideoEndedDateProperty, videoId)
	form.Append(switchTitleSubtitle(ended, "ended", "Ended", "On: mark as ended. Off: mark as new."))

	if videoEndedDate, ok := rdx.GetLastVal(data.VideoEndedDateProperty, videoId); ok && videoEndedDate != "" {
		endedDateTime := videoEndedDate
		var et time.Time
		if et, err = time.Parse(time.RFC3339, videoEndedDate); err == nil {
			endedDateTime = et.Local().Format(time.DateTime)
		}
		form.Append(strom.CreateText("span", "Last ended: "+endedDateTime).
			SetStyle("color:"+colors.Gray, "font-size:"+font_sizes.Small))
	}

	var videoEndedReason data.VideoEndedReason
	if vers, ok := rdx.GetLastVal(data.VideoEndedReasonProperty, videoId); ok && vers != "" {
		videoEndedReason = data.ParseVideoEndedReason(vers)
	}
	form.Append(videoEndedReasonSelect(videoEndedReason))

	downloadQueued := false
	if dqs, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, videoId); ok {
		downloadQueued = true
		if dcs, sure := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); sure {
			if dqs < dcs {
				downloadQueued = false
			}
		}
	}

	form.Append(switchTitleSubtitle(downloadQueued, "download-queued", "Download queued", "On: add to download queue. Off: remove from download queue."))

	if dqs, ok := rdx.GetLastVal(data.VideoDownloadQueuedProperty, videoId); ok {

		dqDateTime := dqs
		var dqt time.Time
		if dqt, err = time.Parse(time.RFC3339, dqs); err == nil {
			dqDateTime = dqt.Local().Format(time.DateTime)
		}

		form.Append(strom.CreateText("span", "Last download queued: "+dqDateTime).
			SetStyle("color:"+colors.Gray, "font-size:"+font_sizes.Small))

	}

	if videoDownloadCompleted, ok := rdx.GetLastVal(data.VideoDownloadCompletedProperty, videoId); ok && videoDownloadCompleted != "" {
		downloadCompletedTime := videoDownloadCompleted
		var dct time.Time
		if dct, err = time.Parse(time.RFC3339, videoDownloadCompleted); err == nil {
			downloadCompletedTime = dct.Local().Format(time.DateTime)
		}
		form.Append(strom.CreateText("span", "Last download completed: "+downloadCompletedTime).
			SetStyle("color:"+colors.Gray, "font-size:"+font_sizes.Small))
	}

	forcedDownload := rdx.HasKey(data.VideoForcedDownloadProperty, videoId)
	form.Append(switchTitleSubtitle(forcedDownload, "forced-download", "Forced download", "On: re-download if file exists. Off: skip re-downloading."))

	body.Append(submitButton("Update", "manage-video"))

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func videoEndedReasonSelect(videoEndedReason data.VideoEndedReason) strom.Element {

	row := strom.Create("ul", atoms.FlexRow(sizes.Normal)...).
		AddAtom(atoms.AlignItemsCenter)

	ers := strom.Create("select").
		SetAttribute("name", "ended-reason").
		AddAtom(atoms.FontSizeNormal)

	for _, ver := range data.AllVideoEndedReasons() {
		opt := strom.CreateText("option", ver.String())
		if ver == videoEndedReason {
			opt.SetAttribute("selected")
		}
		ers.Append(opt)
	}

	row.Append(ers)

	row.Append(titleSubtitle("ended-reason", "Ended reason", "Optional explanation why the video has ended."))

	return row
}
