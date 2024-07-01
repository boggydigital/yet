package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/url"
)

func MigrateMetadataHandler(u *url.URL) error {
	return MigrateMetadata(nil)
}

func MigrateMetadata(rdx kvas.WriteableRedux) error {

	mma := nod.Begin("migrating metadata...")
	defer mma.End()

	if err := Backup(); err != nil {
		return mma.EndWithError(err)
	}

	fromToMapping := map[string]string{
		"playlist-download-queue": data.PlaylistAutoDownloadProperty,
		"playlist-watchlist":      data.PlaylistAutoRefreshProperty,
		"playlist-single-format":  data.PlaylistPreferSingleFormatProperty,

		"video-ended":                  data.VideoEndedDateProperty,
		"video-downloaded-date":        data.VideoDownloadCompletedProperty,
		"video-single-format-download": data.VideoPreferSingleFormatProperty,
	}

	properties := make([]string, 0, len(fromToMapping)+2)
	for from, to := range fromToMapping {
		properties = append(properties, from, to)
	}

	properties = append(properties, "video-skipped", "videos-download-queue", data.VideoEndedReasonProperty, data.VideoDownloadQueuedProperty)

	var err error
	rdx, err = validateWritableRedux(rdx, properties...)
	if err != nil {
		return mma.EndWithError(err)
	}

	propertyIdValues := make(map[string]map[string][]string)

	for from, to := range fromToMapping {

		if propertyIdValues[to] == nil {
			propertyIdValues[to] = make(map[string][]string)
		}

		for _, id := range rdx.Keys(from) {
			if values, ok := rdx.GetAllValues(from, id); ok {
				propertyIdValues[to][id] = values
			}
		}
	}

	// video-skipped requires special handling to set the correct values
	property := "video-skipped"
	if propertyIdValues[data.VideoEndedReasonProperty] == nil {
		propertyIdValues[data.VideoEndedReasonProperty] = make(map[string][]string)
	}
	for _, id := range rdx.Keys(property) {
		if vs, ok := rdx.GetLastVal(property, id); ok && vs == "true" {
			propertyIdValues[data.VideoEndedReasonProperty][id] = []string{string(data.Skipped)}
		}
	}

	// videos-download-queue requires special handling to set the correct values
	property = "videos-download-queue"
	if propertyIdValues[data.VideoDownloadQueuedProperty] == nil {
		propertyIdValues[data.VideoDownloadQueuedProperty] = make(map[string][]string)
	}
	for _, id := range rdx.Keys(property) {
		if vs, ok := rdx.GetLastVal(property, id); ok && vs == "true" {
			propertyIdValues[data.VideoDownloadQueuedProperty][id] = []string{yeti.FmtNow()}
		}
	}

	for p, idValues := range propertyIdValues {
		if err := rdx.BatchAddValues(p, idValues); err != nil {
			return mma.EndWithError(err)
		}
	}

	mma.EndWithResult("done")

	return nil
}
