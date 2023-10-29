package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func TestDependenciesHandler(u *url.URL) error {
	return TestDependencies()
}

func TestDependencies() error {

	tda := nod.Begin("testing dependencies...")
	defer tda.End()

	for _, bn := range yeti.AllBinaries() {
		ba := nod.Begin(" " + string(bn))
		if bp := yeti.GetBinary(bn); bp != "" {
			ba.EndWithResult(bp)
		} else {
			ba.EndWithResult("not found")
		}
	}

	return nil
}
