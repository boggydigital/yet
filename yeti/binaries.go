package yeti

import (
	"golang.org/x/exp/maps"
	"os"
	"os/exec"
)

type Binary string

const (
	FFMpegBin Binary = "ffmpeg"
	NodeBin   Binary = "node"
	DenoBin   Binary = "deno"
)

var cmdEnv = map[Binary]string{
	FFMpegBin: "YET_FFMPEG_CMD",
	NodeBin:   "YET_NODE_CMD",
	DenoBin:   "YET_DENO_CMD",
}

func AllBinaries() []Binary {
	return maps.Keys(cmdEnv)
}

func GetBinary(name Binary) string {

	if bin := os.Getenv(cmdEnv[name]); bin != "" {
		return bin
	} else {
		if path, err := exec.LookPath(string(name)); err == nil {
			return path
		}
	}

	return ""
}

func IsJSBinaryAvailable() bool {
	return GetBinary(NodeBin) != "" || GetBinary(DenoBin) != ""
}
