package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/sizes"
)

const (
	paramVideoId       = "video-id"
	paramQueueDownload = "queue-download"
	paramDownloadVideo = "download-video"
)

func GetPaste(w http.ResponseWriter, r *http.Request) {

	root, body := strom.RootBody("Paste", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
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

	downloadParameters := strom.Create("ul", atoms.FlexColWrap(sizes.Normal)...)
	form.Append(downloadParameters)

	queueDownload := strom.Create("li", atoms.FlexRowWrap(sizes.Normal)...)
	downloadParameters.Append(queueDownload)

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
	downloadParameters.Append(downloadNow)

	downloadNow.Append(strom.Create("input").
		SetAttribute("id", "name", paramDownloadVideo).
		SetAttribute("type", "checkbox").
		SetAttribute("switch"))

	downloadNow.Append(strom.CreateText("label", "Download now").
		SetAttribute("for", paramDownloadVideo))

	body.Append(submitButton("Paste", form.GetAttribute("id")))

	if err := strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
