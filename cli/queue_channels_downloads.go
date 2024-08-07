package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func QueueChannelsDownloadsHandler(u *url.URL) error {
	return QueueChannelsDownloads(nil)
}

func QueueChannelsDownloads(rdx kevlar.WriteableRedux) error {

	qcda := nod.NewProgress("queueing channels downloads...")
	defer qcda.End()

	var err error
	rdx, err = validateWritableRedux(rdx, data.AllProperties()...)
	if err != nil {
		return qcda.EndWithError(err)
	}

	channelsIds := rdx.Keys(data.ChannelAutoDownloadProperty)
	qcda.TotalInt(len(channelsIds))

	for _, channelId := range channelsIds {

		if err := queueChannelDownloads(rdx, channelId); err != nil {
			return qcda.EndWithError(err)
		}

		qcda.Increment()
	}

	qcda.EndWithResult("done")

	return nil
}

// queueChannelDownloads goes through channel videos according to the download policy,
// skips ended and previously queued videos and queues the rest
func queueChannelDownloads(rdx kevlar.WriteableRedux, channelId string) error {

	queue := make(map[string][]string)

	for _, videoId := range yeti.ChannelNotEndedVideos(channelId, rdx) {
		if rdx.HasKey(data.VideoDownloadQueuedProperty, videoId) {
			continue
		}
		queue[videoId] = []string{yeti.FmtNow()}
	}

	return rdx.BatchAddValues(data.VideoDownloadQueuedProperty, queue)
}
