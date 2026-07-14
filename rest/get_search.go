package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
)

func GetSearch(w http.ResponseWriter, r *http.Request) {

	root := strom.Page("Search")

	var body strom.Element
	for body = range root.GetElementsByTagName("body") {
		break
	}

	body.AddClass("d-f", "fd-c", "rg-l")

	body.Append(navButton("Home", "/"))

	body.Append(strom.CreateText("h2", "Search YouTube videos"))

	form := strom.Create("form", "d-f", "fd-c", "rg-n").
		SetAttribute("id", "search-form").
		SetAttribute("method", "get").
		SetAttribute("action", "/results")

	body.Append(form)

	form.Append(strom.CreateText("label", "Search terms").
		SetAttribute("for", "search-query"))

	form.Append(strom.Create("input").
		SetAttribute("id", "name", "search-query").
		SetAttribute("type", "search").
		SetAttribute("placeholder", "Search terms").
		SetAttribute("autofocus").
		SetAttribute("required").
		SetStyle(textInputStyles()))

	body.Append(submitButton("Search", form.GetAttribute("id")))

	if err := strom.WriteResponse(w, root); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
