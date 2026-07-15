package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars"
)

func GetPaste(w http.ResponseWriter, r *http.Request) {

	root := strom.Page("Paste")

	var body strom.Element
	for body = range root.GetElementsByTagName("body") {
		break
	}

	body.AddClass("d-f", "fd-c", "rg-l")

	body.Append(navButton("Home", "/"))

	body.Append(
		strom.CreateText("h1", "Paste YouTube link or video-id"))

	form := strom.Create("form", "d-f", "fd-c", "rg-n").
		SetAttribute("id", "paste-form").
		SetAttribute("method", "get").
		SetAttribute("action", "/paste_video")
	body.Append(form)

	form.Append(
		strom.Create("input").
			SetAttribute("id", "name", "video-id").
			SetAttribute("type", "text").
			SetAttribute("placeholder", "YouTube link or video-id").
			SetAttribute("autofocus").
			SetAttribute("required").
			SetStyle(textInputStyles()))

	downloadParameters := strom.Create("ul").AddClass("d-f", "fd-c", "rg-n")
	form.Append(downloadParameters)

	queueDownload := strom.Create("li").AddClass("d-f", "fd-r", "cg-n")
	downloadParameters.Append(queueDownload)

	queueDownload.Append(
		strom.Create("input").
			SetAttribute("id", "name", "queue-download").
			SetAttribute("type", "checkbox").
			SetAttribute("switch").
			SetAttribute("checked"))

	queueDownload.Append(
		strom.CreateText("label", "Queue download").
			SetAttribute("for", "queue-download"))

	downloadNow := strom.Create("li").AddClass("d-f", "fd-r", "cg-n")
	downloadParameters.Append(downloadNow)

	downloadNow.Append(strom.Create("input").
		SetAttribute("id", "name", "download-now").
		SetAttribute("type", "checkbox").
		SetAttribute("switch"))

	downloadNow.Append(strom.CreateText("label", "Download now").
		SetAttribute("for", "download-now"))

	body.Append(submitButton("Paste", form.GetAttribute("id")))

	if err := strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func navButton(title, href string) strom.Element {

	button := strom.Create("a", "br-l", "fs-n", "fw-b").
		SetTextContent(title).
		SetAttribute("href", href).
		SetStyle(map[string]string{
			"background-color": vars.Color(vars.ColorBlue),
			"color":            vars.Color(vars.ColorBackground),
			"width":            "max-content",
			"padding-inline":   vars.Size(vars.SizeNormal),
			"padding-block":    vars.Size(vars.SizeSmall),
		})

	return button
}

func submitButton(value, form string) strom.Element {
	return strom.Create("input", "br-l", "fs-n", "fw-b").
		SetAttribute("type", "submit").
		SetAttribute("form", form).
		SetAttribute("value", value).
		SetStyle(map[string]string{
			"margin-block-start": vars.Size(vars.SizeNormal),
			"appearance":         "none",
			"padding-inline":     vars.Size(vars.SizeNormal),
			"padding-block":      vars.Size(vars.SizeSmall),
			"background-color":   vars.Color(vars.ColorPurple),
			"color":              vars.Color(vars.ColorBackground),
			"border":             "none",
			"width":              "max-content",
		})
}

func textInputStyles() map[string]string {
	return map[string]string{
		"max-width": "calc(1.5 * " + vars.Size(vars.SizeXXXLarge) + ")",
		"padding":   vars.Size(vars.SizeSmall),
		"font-size": vars.FontSize(vars.SizeNormal),
	}
}
