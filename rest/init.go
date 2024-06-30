package rest

import (
	"embed"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"html/template"
)

var (
	rdx         kvas.ReadableRedux
	progressRdx kvas.WriteableRedux
	tmpl        *template.Template
	//go:embed "templates/*.gohtml"
	templates embed.FS
)

func Init() error {

	metadataDir, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		return err
	}

	if rdx, err = kvas.NewReduxReader(metadataDir, data.AllProperties()...); err != nil {
		return err
	}

	if progressRdx, err = kvas.NewReduxWriter(metadataDir,
		data.VideoProgressProperty,
		data.VideoEndedDateProperty,
	); err != nil {
		return err
	}

	tmpl = template.Must(
		template.
			New("").
			ParseFS(templates, "templates/*.gohtml"))

	return err
}
