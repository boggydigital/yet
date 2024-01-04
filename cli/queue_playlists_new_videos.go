package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathology"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/paths"
	"net/url"
)

func QueuePlaylistsNewVideosHandler(u *url.URL) error {
	return QueuePlaylistsNewVideos()
}

func QueuePlaylistsNewVideos() error {

	qpnva := nod.NewProgress("queueing playlists new videos...")
	defer qpnva.End()

	metadataDir, err := pathology.GetAbsDir(paths.Metadata)
	if err != nil {
		return qpnva.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(metadataDir, data.AllProperties()...)
	if err != nil {
		return qpnva.EndWithError(err)
	}

	notNewIndicatorProperties := []string{
		data.VideosDownloadQueueProperty,
		data.VideosWatchlistProperty,
		data.VideoEndedProperty,
		data.VideoProgressProperty}

	for _, pdq := range rdx.Keys(data.PlaylistDownloadQueueProperty) {
		if newVideos, ok := rdx.GetAllValues(data.PlaylistNewVideosProperty, pdq); ok && len(newVideos) > 0 {
			for _, videoId := range newVideos {

				skipVideo := false
				// don't add videos already in download queue, watch, ended, in progress
				for _, nnip := range notNewIndicatorProperties {
					if rdx.HasKey(nnip, videoId) {
						skipVideo = true
						break
					}
				}

				if skipVideo {
					break
				}

				if err := rdx.AddValues(data.VideosDownloadQueueProperty, videoId, data.TrueValue); err != nil {
					return qpnva.EndWithError(err)
				}
			}
		}
	}

	qpnva.EndWithResult("done")

	return nil
}
