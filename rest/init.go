package rest

import (
	"embed"
	"html/template"

	"github.com/boggydigital/camino"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
)

var (
	rdx  redux.Writeable
	tmpl *template.Template
	//go:embed "templates/*.gohtml"
	templates embed.FS
)

func Init() error {

	metadataDir := camino.GetAbs(data.Metadata)

	var err error
	if rdx, err = redux.NewWriter(metadataDir, data.AllProperties()...); err != nil {
		return err
	}

	tmpl = template.Must(
		template.
			New("").
			ParseFS(templates, "templates/*.gohtml"))

	return err
}
