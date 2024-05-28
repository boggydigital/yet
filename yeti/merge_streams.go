package yeti

import (
	"errors"
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet_urls/youtube_urls"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func MergeStreams(relFilename string, force bool) error {
	absVideosDir, err := pathways.GetAbsDir(paths.Videos)
	if err != nil {
		return err
	}

	relVideoFilename, relAudioFilename := VideoAudioFilenames(relFilename)

	absVideoFilename := filepath.Join(absVideosDir, relVideoFilename)
	absAudioFilename := filepath.Join(absVideosDir, relAudioFilename)
	absFilename := filepath.Join(absVideosDir, relFilename)

	if _, err := os.Stat(absFilename); err == nil && force {
		if err := os.Remove(absFilename); err != nil {
			return err
		}
	}

	ffmb := GetBinary(FFMpegBin)
	if ffmb == "" {
		return errors.New("ffmpeg not available")
	}

	//merge streams into a single file
	//since youtube_urls filters to mp4 formats only, we don't need to do any transcoding
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

func MergeSegments(id, dir string, segments ...string) (string, error) {

	msa := nod.Begin("merging segments for %s...", id)
	defer msa.End()

	tempFilename := filepath.Join(os.TempDir(), id+".txt")
	tempFile, err := os.Create(tempFilename)
	if err != nil {
		return "", err
	}

	for _, segment := range segments {
		line := fmt.Sprintf("file '%s'\n", filepath.Join(dir, path.Base(segment)))
		if _, err := tempFile.WriteString(line); err != nil {
			return "", err
		}
	}

	ffmb := GetBinary(FFMpegBin)
	if ffmb == "" {
		return "", errors.New("ffmpeg not available")
	}

	tempOutputFilename := filepath.Join(os.TempDir(), id+youtube_urls.DefaultVideoExt)

	args := []string{"-f", "concat", "-safe", "0", "-i", tempFilename, "-c", "copy", tempOutputFilename}

	cmd := exec.Command(ffmb, args...)
	if err := cmd.Run(); err != nil {
		return "", msa.EndWithError(err)
	}

	return tempOutputFilename, nil
}
