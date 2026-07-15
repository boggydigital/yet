package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
)

func GetSearch(w http.ResponseWriter, r *http.Request) {

	root, body := strom.RootBody("Search")

	body.AddClass("d-f", "fd-c", "rg-n")

	body.Append(navButton("Home", "/"))

	body.Append(strom.CreateText("h1", "Search YouTube videos"))

	form := strom.Create("form", "d-f", "fd-c", "rg-n").
		SetAttribute("id", "search-form").
		SetAttribute("method", "get").
		SetAttribute("action", "/results")

	body.Append(form)

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
