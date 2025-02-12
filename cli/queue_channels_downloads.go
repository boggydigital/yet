package cli

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func QueueChannelsDownloadsHandler(u *url.URL) error {
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

		if err := queueChannelDownloads(rdx, channelId); err != nil {
			return err
		}

		qcda.Increment()
	}

	return nil
}

// queueChannelDownloads goes through channel videos according to the download policy,
// skips ended and previously queued videos and queues the rest
func queueChannelDownloads(rdx redux.Writeable, channelId string) error {

	queue := make(map[string][]string)

	for _, videoId := range yeti.ChannelNotEndedVideos(channelId, rdx) {
		if rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			continue
		}
		queue[videoId] = []string{yeti.FmtNow()}
	}

	return rdx.BatchAddValues(data.VideoDownloadQueuedProperty, queue)
}
