package main

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/cli"
	"github.com/boggydigital/yet/paths"
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
	defer ya.End()

	if err := pathways.Setup(dirOverridesFilename,
		paths.DefaultYetRootDir,
		nil,
		paths.AllAbsDirs...); err != nil {
		_ = ya.EndWithError(err)
		os.Exit(1)
	}

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)

	if err != nil {
		_ = ya.EndWithError(err)
		os.Exit(1)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"add-playlist":               cli.AddPlaylistHandler,
		"add-video":                  cli.AddVideoHandler,
		"backup":                     cli.BackupHandler,
		"cleanup-ended":              cli.CleanupEndedHandler,
		"download-video":             cli.DownloadVideoHandler,
		"download-queue":             cli.DownloadQueueHandler,
		"get-captions":               cli.GetCaptionsHandler,
		"get-channel-metadata":       cli.GetChannelMetadataHandler,
		"get-playlist-metadata":      cli.GetPlaylistMetadataHandler,
		"get-poster":                 cli.GetPosterHandler,
		"get-rutube-video":           cli.GetRuTubeVideoHandler,
		"get-video-metadata":         cli.GetVideoMetadataHandler,
		"migrate-metadata":           cli.MigrateMetadataHandler,
		"queue-playlists-downloads":  cli.QueuePlaylistsDownloadsHandler,
		"remove-playlist":            cli.RemovePlaylistHandler,
		"remove-videos":              cli.RemoveVideosHandler,
		"serve":                      cli.ServeHandler,
		"sync":                       cli.SyncHandler,
		"test-dependencies":          cli.TestDependenciesHandler,
		"refresh-playlists-metadata": cli.RefreshPlaylistsMetadataHandler,
		"version":                    cli.VersionHandler,
	})

	if err := defs.AssertCommandsHaveHandlers(); err != nil {
		_ = ya.EndWithError(err)
		os.Exit(1)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		_ = ya.EndWithError(err)
		os.Exit(1)
	}
}
