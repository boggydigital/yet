package cli

import (
	"net/url"

	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
)

func QueueChannelsDownloadsHandler(_ *url.URL) error {
	return QueueChannelsDownloads(nil)
}

func QueueChannelsDownloads(rdx redux.Writeable) error {

	qcda := nod.NewProgress("queueing channels downloads...")
	defer qcda.Done()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return err
	}

	qcda.TotalInt(rdx.Len(data.ChannelAutoDownloadProperty))

	for channelId := range rdx.Keys(data.ChannelAutoDownloadProperty) {

		if err = yeti.QueueChannelDownloads(rdx, channelId); err != nil {
			return err
		}

		qcda.Increment()
	}

	return nil
}
