package main

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/cli"
	"github.com/boggydigital/yet/data"
	"log"
	"os"
)

var (
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

const (
	dirOverridesFilename = "directories.txt"
)

func main() {
	nod.EnableStdOutPresenter()

	ya := nod.Begin("yet is serving your videos needs")
	defer ya.Done()

	if err := pathways.Setup(dirOverridesFilename,
		data.DefaultYetRootDir,
		nil,
		data.AllAbsDirs...); err != nil {
		log.Fatalln(err)
	}

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)

	if err != nil {
		log.Fatalln(err)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"add-channel":                 cli.AddChannelHandler,
		"add-playlist":                cli.AddPlaylistHandler,
		"add-video":                   cli.AddVideoHandler,
		"backup":                      cli.BackupHandler,
		"cleanup-ended-videos":        cli.CleanupEndedVideosHandler,
		"dehydrate-posters":           cli.DehydratePostersHandler,
		"download-video":              cli.DownloadVideoHandler,
		"get-captions":                cli.GetCaptionsHandler,
		"get-channels-metadata":       cli.GetChannelsMetadataHandler,
		"get-playlists-metadata":      cli.GetPlaylistsMetadataHandler,
		"get-poster":                  cli.GetPosterHandler,
		"get-rutube-video":            cli.GetRuTubeVideoHandler,
		"get-video-metadata":          cli.GetVideoMetadataHandler,
		"process-queue":               cli.ProcessQueueHandler,
		"migrate":                     cli.MigrateHandler,
		"queue-channels-downloads":    cli.QueueChannelsDownloadsHandler,
		"queue-playlists-downloads":   cli.QueuePlaylistsDownloadsHandler,
		"refresh-channels-metadata":   cli.RefreshChannelsMetadataHandler,
		"refresh-playlists-metadata":  cli.RefreshPlaylistsMetadataHandler,
		"remove-channel":              cli.RemoveChannelHandler,
		"remove-playlist":             cli.RemovePlaylistHandler,
		"remove-videos":               cli.RemoveVideosHandler,
		"scrub-deposition-properties": cli.ScrubDepositionPropertiesHandler,
		"scrub-ended-properties":      cli.ScrubEndedPropertiesHandler,
		"serve":                       cli.ServeHandler,
		"sync":                        cli.SyncHandler,
		"update-yt-dlp":               cli.UpdateYtDlpHandler,
		"version":                     cli.VersionHandler,
	})

	if err = defs.AssertCommandsHaveHandlers(); err != nil {
		log.Fatalln(err)
	}

	if err = defs.Serve(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}
