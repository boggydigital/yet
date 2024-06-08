package cli

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func MigrateHandler(u *url.URL) error {
	return Migrate()
}

func Migrate() error {

	rdr, err := pathways.GetAbsDir(paths.Metadata)
	if err != nil {
		return err
	}

	vdr, err := pathways.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	rdx, err := kvas.NewReduxReader(rdr, data.AllProperties()...)
	if err != nil {
		return err
	}

	for _, id := range rdx.Keys(data.VideoDownloadedDateProperty) {

		if title, ok := rdx.GetLastVal(data.VideoTitleProperty, id); ok {
			if channel, sure := rdx.GetLastVal(data.VideoOwnerChannelNameProperty, id); sure {

				relOldFn := yeti.OldChannelTitleVideoIdFilename(channel, title, id)
				relNewFn := yeti.ChannelTitleVideoIdFilename(channel, title, id)
				absOldFn := path.Join(vdr, relOldFn)
				absNewFn := path.Join(vdr, relNewFn)

				fmt.Println(relOldFn)

				if _, err := os.Stat(absOldFn); err == nil {
					oldFile, err := os.Open(absOldFn)
					if err != nil {
						return err
					}

					dir, _ := filepath.Split(absNewFn)
					if err := os.MkdirAll(dir, 0755); err != nil {
						return err
					}

					newFile, err := os.Create(absNewFn)
					if err != nil {
						return err
					}

					if _, err := io.Copy(newFile, oldFile); err != nil {
						return err
					}

					oldFile.Close()
					newFile.Close()

					if err := os.Remove(absOldFn); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil

}
