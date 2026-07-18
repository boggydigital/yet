package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/calc"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/font_sizes"
	"github.com/boggydigital/strom/vars/sizes"
)

func GetPaste(w http.ResponseWriter, r *http.Request) {

	root, body := strom.RootBody("Paste", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))

	topRow.Append(strom.CreateText("h2", "Paste"))

	form := strom.Create("form", atoms.FlexColWrap(sizes.Normal)...).
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

	downloadParameters := strom.Create("ul", atoms.FlexColWrap(sizes.Normal)...)
	form.Append(downloadParameters)

	queueDownload := strom.Create("li", atoms.FlexRowWrap(sizes.Normal)...)
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

	downloadNow := strom.Create("li", atoms.FlexRow(sizes.Normal)...)
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

func roundedButton(title, href string, clr colors.Color) strom.Element {

	return strom.Create("a", atoms.BorderRadiusLarge, atoms.FontSizeNormal, atoms.FontWeightBold).
		SetTextContent(title).
		SetAttribute("href", href).
		SetStyle(buttonStyles(colors.Gray))
}

func navButton(title, href string) strom.Element {
	return roundedButton(title, href, colors.Blue)
}

func actionButton(title, href string) strom.Element {
	return roundedButton(title, href, colors.Purple)
}

func submitButton(value, form string) strom.Element {
	return strom.Create("input", atoms.BorderRadiusLarge, atoms.FontSizeNormal, atoms.FontWeightBold).
		SetAttribute("type", "submit").
		SetAttribute("form", form).
		SetAttribute("value", value).
		SetStyle(map[string]string{"appearance": "none"}).
		SetStyle(buttonStyles(colors.Gray))
}

func buttonStyles(c string) map[string]string {
	return map[string]string{
		"padding-inline":   sizes.Normal,
		"padding-block":    sizes.Small,
		"background-color": c,
		"color":            colors.Background,
		"border":           "none",
		"width":            "max-content",
		"font-size":        font_sizes.XSmall,
	}
}

func textInputStyles() map[string]string {
	return map[string]string{
		"max-width": calc.Mult(sizes.XXXLarge, 1.5),
		"padding":   sizes.Small,
		"font-size": font_sizes.Normal,
	}
}
