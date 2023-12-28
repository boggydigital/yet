package cli

import "net/url"

func ExportHandler(u *url.URL) error {
	return Export()
}

func Export() error {
	return nil
}
