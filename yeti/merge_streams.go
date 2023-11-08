package yeti

import (
	"errors"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"os"
	"os/exec"
	"path/filepath"
)

func mergeStreams(relFilename string) error {
	absVideosDir, err := paths.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	relVideoFilename, relAudioFilename := videoAudioFilenames(relFilename)

	absVideoFilename := filepath.Join(absVideosDir, relVideoFilename)
	absAudioFilename := filepath.Join(absVideosDir, relAudioFilename)
	absFilename := filepath.Join(absVideosDir, relFilename)

	ffmb := GetBinary(FFMpegBin)
	if ffmb == "" {
		return errors.New("ffmpeg not available")
	}

	//merge streams into a single file
	//since yt_urls filters to mp4 formats only, we don't need to do any transcoding
	//and can quickly merge by copying streams:
	//ffmpeg -i video.mp4 -i audio.wav -c copy output.mp4
	ma := nod.Begin("merging streams for %s...", relFilename)
	args := []string{"-i", absVideoFilename, "-i", absAudioFilename, "-c", "copy", absFilename}
	cmd := exec.Command(ffmb, args...)
	if err := cmd.Run(); err != nil {
		return ma.EndWithError(err)
	}

	//cleanup separate streams after successful merge
	if err := os.Remove(absVideoFilename); err != nil {
		return ma.EndWithError(err)
	}
	if err := os.Remove(absAudioFilename); err != nil {
		return ma.EndWithError(err)
	}

	return nil
}