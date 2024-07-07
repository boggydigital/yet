package rest

import (
	"embed"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"html/template"
)

var (
	rdx  kevlar.WriteableRedux
	tmpl *template.Template
	//go:embed "templates/*.gohtml"
	templates embed.FS
)

func Init() error {

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		return err
	}

	if rdx, err = kevlar.NewReduxWriter(metadataDir, data.AllProperties()...); err != nil {
		return err
	}

	tmpl = template.Must(
		template.
			New("").
			ParseFS(templates, "templates/*.gohtml"))

	return err
}
