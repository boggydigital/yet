package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/sizes"
)

const (
	paramVideoId       = "video-id"
	paramQueueDownload = "queue-download"
	paramDownloadVideo = "download-video"
	paramWithProgress  = "with-progress"
)

func GetPaste(w http.ResponseWriter, r *http.Request) {

	root, body := strom.RootBody("Paste", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Normal)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))

	topRow.Append(strom.CreateText("h2", "Paste"))

	form := strom.Create("form", atoms.FlexColWrap(sizes.Normal)...).
		SetAttribute("id", "paste-form").
		SetAttribute("method", "post").
		SetAttribute("action", "/paste")
	body.Append(form)

	form.Append(
		strom.Create("input").
			SetAttribute("id", "name", paramVideoId).
			SetAttribute("type", "text").
			SetAttribute("placeholder", "YouTube link or video-id").
			SetAttribute("autofocus").
			SetAttribute("required").
			SetStyle(textInputStyles()...))

	pasteParameters := strom.Create("ul", atoms.FlexColWrap(sizes.Normal)...)
	form.Append(pasteParameters)

	queueDownload := strom.Create("li", atoms.FlexRowWrap(sizes.Normal)...)
	pasteParameters.Append(queueDownload)

	queueDownload.Append(
		strom.Create("input").
			SetAttribute("id", "name", paramQueueDownload).
			SetAttribute("type", "checkbox").
			SetAttribute("switch").
			SetAttribute("checked"))

	queueDownload.Append(
		strom.CreateText("label", "Queue download").
			SetAttribute("for", paramQueueDownload))

	downloadNow := strom.Create("li", atoms.FlexRow(sizes.Normal)...)
	pasteParameters.Append(downloadNow)

	downloadNow.Append(strom.Create("input").
		SetAttribute("id", "name", paramDownloadVideo).
		SetAttribute("type", "checkbox").
		SetAttribute("switch"))

	downloadNow.Append(strom.CreateText("label", "Download now").
		SetAttribute("for", paramDownloadVideo))

	saveProgress := strom.Create("li", atoms.FlexRow(sizes.Normal)...)
	pasteParameters.Append(saveProgress)

	saveProgress.Append(strom.Create("input").
		SetAttribute("id", "name", paramWithProgress).
		SetAttribute("type", "checkbox").
		SetAttribute("switch").
		SetAttribute("checked"))

	saveProgress.Append(strom.CreateText("label", "With progress (when provided)").
		SetAttribute("for", paramWithProgress))

	body.Append(submitButton("Paste", form.GetAttribute("id"), colors.Green))

	if err := strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
