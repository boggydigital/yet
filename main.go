package main

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
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
	userDirsFilename = "directories.txt"
)

func main() {
	nod.EnableStdOutPresenter()

	ya := nod.Begin("yet is serving your videos needs")
	defer ya.End()

	if err := chRoot(userDirsFilename, paths.DefaultDirs); err != nil {
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
		"clear-ended":       cli.ClearEndedHandler,
		"clear-progress":    cli.ClearProgressHandler,
		"download":          cli.DownloadHandler,
		"get-poster":        cli.GetPosterHandler,
		"serve":             cli.ServeHandler,
		"set-ended":         cli.SetEndedHandler,
		"test-dependencies": cli.TestDependenciesHandler,
		"version":           cli.VersionHandler,
		"watchlist-add":     cli.WatchlistAddHandler,
		"watchlist-remove":  cli.WatchlistRemoveHandler,
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

func chRoot(userDirsFilename string, defaultDirs map[string]string) error {

	var userDirs map[string]string

	if _, err := os.Stat(userDirsFilename); err == nil {
		udFile, err := os.Open(userDirsFilename)
		if err != nil {
			return err
		}

		userDirs, err = wits.ReadKeyValue(udFile)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		userDirs = defaultDirs
	} else {
		return err
	}

	return paths.SetAbsDirs(userDirs)
}
