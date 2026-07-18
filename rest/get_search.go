package rest

import (
	"net/http"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/sizes"
)

func GetSearch(w http.ResponseWriter, r *http.Request) {

	root, body := strom.RootBody("Search", atoms.FlexCol(sizes.Normal)...)

	topRow := strom.Create("ul", atoms.FlexRow(sizes.Small)...).AddAtom(atoms.AlignItemsCenter)
	body.Append(topRow)

	topRow.Append(navButton("Home", "/"))
	topRow.Append(strom.CreateText("h2", "Search"))

	form := strom.Create("form", atoms.FlexColWrap(sizes.Normal)...).
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
