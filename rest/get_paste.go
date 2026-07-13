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

	body.SetStyle(map[string]string{
		"padding":        vars.Size(vars.SizeNormal),
		"display":        "flex",
		"flex-direction": "column",
		"row-gap":        vars.Size(vars.SizeLarge),
	})

	body.Append(NavButton("Home", "/"))

	body.Append(strom.CreateText("h2", "Paste YouTube link or video-id"))

	form := strom.Create("form").
		SetAttribute("id", "paste-form").
		SetAttribute("method", "get").
		SetAttribute("action", "/paste_video").
		SetStyle(map[string]string{
			"display":        "flex",
			"flex-direction": "column",
			"row-gap":        vars.Size(vars.SizeNormal),
		})

	body.Append(form)

	form.Append(strom.CreateText("label", "YouTube link or video-id").
		SetAttribute("for", "video-id"))

	form.Append(strom.Create("input").
		SetAttribute("id", "video-id").
		SetAttribute("name", "video-id").
		SetAttribute("type", "text").
		SetAttribute("placeholder", "YouTube link or video-id").
		SetAttribute("autofocus", "").
		SetAttribute("required", "").
		SetStyle(map[string]string{
			"max-width": "calc(1.5 * " + vars.Size(vars.SizeXXXLarge) + ")",
			"padding":   vars.Size(vars.SizeSmall),
			"font-size": vars.FontSize(vars.SizeNormal),
		}))

	downloadParameters := strom.Create("ul").
		SetStyle(map[string]string{
			"display":        "flex",
			"flex-direction": "column",
			"row-gap":        vars.Size(vars.SizeNormal),
		})
	form.Append(downloadParameters)

	queueDownload := strom.Create("li").
		SetStyle(map[string]string{
			"display":        "flex",
			"flex-direction": "row",
			"column-gap":     vars.Size(vars.SizeNormal),
		})
	downloadParameters.Append(queueDownload)

	queueDownload.Append(strom.Create("input").
		SetAttribute("type", "checkbox").
		SetAttribute("switch", "").
		SetAttribute("checked", "").
		SetAttribute("id", "queue-download").
		SetAttribute("name", "queue-download"))

	queueDownload.Append(strom.CreateText("label", "Queue download").
		SetAttribute("for", "queue-download"))

	downloadNow := strom.Create("li").
		SetStyle(map[string]string{
			"display":        "flex",
			"flex-direction": "row",
			"column-gap":     vars.Size(vars.SizeNormal),
		})
	downloadParameters.Append(downloadNow)

	downloadNow.Append(strom.Create("input").
		SetAttribute("type", "checkbox").
		SetAttribute("switch", "").
		SetAttribute("id", "download-now").
		SetAttribute("name", "download-now"))

	downloadNow.Append(strom.CreateText("label", "Download now").
		SetAttribute("for", "download-now"))

	body.Append(strom.Create("input").
		SetAttribute("type", "submit").
		SetAttribute("form", "paste-form").
		SetAttribute("value", "Paste").
		SetStyle(map[string]string{
			"margin-block-start": vars.Size(vars.SizeNormal),
			"appearance":         "none",
			"padding-inline":     vars.Size(vars.SizeNormal),
			"padding-block":      vars.Size(vars.SizeSmall),
			"background-color":   vars.Color(vars.ColorPurple),
			"color":              vars.Color(vars.ColorBackground),
			"border-radius":      vars.Size(vars.SizeLarge),
			"border":             "none",
			"font-size":          vars.FontSize(vars.SizeNormal),
			"font-weight":        vars.FontWeight(vars.WeightBold),
			"width":              "max-content",
		}))

	if err := strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NavButton(title, href string) strom.Element {

	button := strom.Create("a").
		SetTextContent(title).
		SetAttribute("href", href).
		SetStyle(map[string]string{
			"background-color": vars.Color(vars.ColorBlue),
			"color":            vars.Color(vars.ColorBackground),
			"width":            "max-content",
			"padding-inline":   vars.Size(vars.SizeNormal),
			"padding-block":    vars.Size(vars.SizeSmall),
			"border-radius":    vars.Size(vars.SizeLarge),
			"font-size":        vars.FontSize(vars.SizeNormal),
			"font-weight":      vars.FontWeight(vars.WeightBold),
		})

	return button
}
