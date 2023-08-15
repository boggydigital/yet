package yeti

import (
	"os"
	"os/exec"
)

const (
	ffmpegCmdEnv = "YET_FFMPEG_CMD"
	nodeCmdEnv   = "YET_NODE_CMD"
	denoCmdEnv   = "YET_DENO_CMD"
)

type Binaries struct {
	FFMpeg string
	NodeJS string
	Deno   string
}

func NewBinaries() *Binaries {

	bins := &Binaries{
		FFMpeg: os.Getenv(ffmpegCmdEnv),
		NodeJS: os.Getenv(nodeCmdEnv),
		Deno:   os.Getenv(denoCmdEnv),
	}

	if bins.FFMpeg == "" {
		if path, err := exec.LookPath("ffmpeg"); err == nil {
			bins.FFMpeg = path
		}
	}

	if bins.NodeJS == "" {
		if path, err := exec.LookPath("node"); err == nil {
			bins.NodeJS = path
		}
	}

	if bins.Deno == "" {
		if path, err := exec.LookPath("deno"); err == nil {
			bins.Deno = path
		}
	}

	return bins
}
