package main

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pasu"
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

	if err := pasu.Setup(dirOverridesFilename,
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
		"add-playlists":               cli.AddPlaylistsHandler,
		"add-videos":                  cli.AddVideosHandler,
		"add-urls":                    cli.AddUrlsHandler,
		"backup":                      cli.BackupHandler,
		"cleanup-ended":               cli.CleanupEndedHandler,
		"download":                    cli.DownloadHandler,
		"get-captions":                cli.GetCaptionsHandler,
		"get-channel-metadata":        cli.GetChannelMetadataHandler,
		"get-url":                     cli.GetUrlHandler,
		"get-playlist-metadata":       cli.GetPlaylistMetadataHandler,
		"get-poster":                  cli.GetPosterHandler,
		"get-video-metadata":          cli.GetVideoMetadataHandler,
		"get-video-file":              cli.GetVideoFileHandler,
		"queue-playlists-new-videos":  cli.QueuePlaylistsNewVideosHandler,
		"remove-playlists":            cli.RemovePlaylistsHandler,
		"remove-videos":               cli.RemoveVideosHandler,
		"remove-urls":                 cli.RemoveUrlsHandler,
		"serve":                       cli.ServeHandler,
		"sync":                        cli.SyncHandler,
		"test-dependencies":           cli.TestDependenciesHandler,
		"update-playlists-metadata":   cli.UpdatePlaylistsMetadataHandler,
		"update-playlists-new-videos": cli.UpdatePlaylistsNewVideosHandler,
		"version":                     cli.VersionHandler,
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
