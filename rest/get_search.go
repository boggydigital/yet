package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars"
)

func GetSearch(w http.ResponseWriter, r *http.Request) {

	root := strom.Page("Search")

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

	body.Append(strom.CreateText("h2", "Search YouTube videos"))

	form := strom.Create("form").
		SetAttribute("id", "search-form").
		SetAttribute("method", "get").
		SetAttribute("action", "/results").
		SetStyle(map[string]string{
			"display":        "flex",
			"flex-direction": "column",
			"row-gap":        vars.Size(vars.SizeNormal),
		})

	body.Append(form)

	form.Append(strom.CreateText("label", "Search terms").
		SetAttribute("for", "search-query"))

	form.Append(strom.Create("input").
		SetAttribute("id", "search-query").
		SetAttribute("name", "search-query").
		SetAttribute("type", "text").
		SetAttribute("placeholder", "Search terms").
		SetAttribute("autofocus", "").
		SetAttribute("required", "").
		SetStyle(map[string]string{
			"max-width": "calc(1.5 * " + vars.Size(vars.SizeXXXLarge) + ")",
			"padding":   vars.Size(vars.SizeSmall),
			"font-size": vars.FontSize(vars.SizeNormal),
		}))

	body.Append(strom.Create("input").
		SetAttribute("type", "submit").
		SetAttribute("form", "search-form").
		SetAttribute("value", "Search").
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
