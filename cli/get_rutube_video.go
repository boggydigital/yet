package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boggydigital/busan"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/rutube_urls"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func GetRuTubeVideoHandler(u *url.URL) error {
	urls := strings.Split(u.Query().Get("url"), ",")
	force := u.Query().Has("Force")
	return GetRuTubeVideo(force, urls...)
}

func GetRuTubeVideo(force bool, urls ...string) error {

	grva := nod.NewProgress("getting Rutube videos...")
	defer grva.End()

	grva.TotalInt(len(urls))

	dc := dolo.DefaultClient

	for _, u := range urls {

		if err := getRuTubeVideo(dc, u, force); err != nil {
			grva.Error(err)
		}

		grva.Increment()
	}

	grva.EndWithResult("done")

	return nil
}

func getRuTubeVideo(dc *dolo.Client, u string, force bool) error {

	ru, err := url.Parse(u)
	if err != nil {
		return err
	}

	videoId, p := path.Base(ru.Path), ru.Query().Get("p")

	grtva := nod.Begin("getting %s...", videoId)
	defer grtva.End()

	playOptions, err := getPlayOptions(videoId, p)
	if err != nil {
		return err
	}

	formats, err := getVideoBalancerFormats(videoId, playOptions)
	if err != nil {
		return err
	}

	segments, err := getVideoSegmentsPlaylist(videoId, formats)
	if err != nil {
		return err
	}

	err = getVideoSegments(playOptions, formats, segments, dc, force)
	if err != nil {
		return err
	}

	err = generateMergeManifest(playOptions, segments, force)
	if err != nil {
		return err
	}

	err = mergeVideoSegments(playOptions, force)
	if err != nil {
		return err
	}

	err = removeManifestSegments(playOptions, segments)
	if err != nil {
		return err
	}

	return nil
}

func getPlayOptions(videoId, p string) (*rutube_urls.PlayOptions, error) {

	gpo := nod.Begin(" getting play options for %s...", videoId)
	defer gpo.End()

	pou := rutube_urls.PlayOptionsUrl(videoId, p)

	resp, err := http.DefaultClient.Get(pou.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var playOptions rutube_urls.PlayOptions

	if err := json.NewDecoder(resp.Body).Decode(&playOptions); err != nil {
		return nil, err
	}

	summary := map[string][]string{
		"author": {playOptions.Author.Name},
		"title":  {strings.Replace(playOptions.Title, "\n", " ", -1)},
	}

	gpo.EndWithSummary("details:", summary)

	return &playOptions, nil
}

func getVideoBalancerFormats(videoId string, playOptions *rutube_urls.PlayOptions) ([]string, error) {

	gvbfa := nod.Begin(" getting video balancer formats for %s...", videoId)
	defer gvbfa.End()

	resp, err := http.DefaultClient.Get(playOptions.VideoBalancer.Default)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	formats := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			formats = append(formats, line)
		}
	}

	if len(formats) == 0 {
		return nil, errors.New("no formats found")
	}

	gvbfa.EndWithResult("done")

	return formats, nil
}

func getVideoSegmentsPlaylist(videoId string, formats []string) ([]string, error) {

	gvspa := nod.Begin(" getting video segments playlist for %s...", videoId)
	defer gvspa.End()

	resp, err := http.DefaultClient.Get(formats[len(formats)-1])
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	segments := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			segments = append(segments, line)
		}
	}

	gvspa.EndWithResult("done")

	return segments, nil
}

func getVideoSegments(
	playOptions *rutube_urls.PlayOptions,
	formats []string,
	segments []string,
	dc *dolo.Client,
	force bool) error {

	videoId := playOptions.VideoId
	channel := busan.Sanitize(playOptions.Author.Name)

	gvsa := nod.NewProgress(" getting video segments for %s...", videoId)
	defer gvsa.End()

	absVideosDir, err := pathways.GetAbsDir(data.Videos)
	if err != nil {
		return err
	}

	gvsa.TotalInt(len(segments))

	lastFormat := formats[len(formats)-1]
	lfBase := path.Base(lastFormat)

	for _, segment := range segments {

		suStr := strings.Replace(lastFormat, lfBase, segment, 1)
		su, err := url.Parse(suStr)
		if err != nil {
			gvsa.Error(err)
			continue
		}

		if err := dc.Download(su, force, nil, absVideosDir, channel, segment); err != nil {
			gvsa.Error(err)
			gvsa.Increment()
			continue
		}

		gvsa.Increment()
	}

	return nil
}

func relManifestFilename(videoId string) string {
	return busan.Sanitize(videoId) + ".txt"
}

func absManifestFilename(playOptions *rutube_urls.PlayOptions) (string, error) {

	videoId := playOptions.VideoId
	channel := busan.Sanitize(playOptions.Author.Name)

	absVideosDir, err := pathways.GetAbsDir(data.Videos)
	if err != nil {
		return "", err
	}

	return filepath.Join(absVideosDir, channel, relManifestFilename(videoId)), nil
}

func generateMergeManifest(playOptions *rutube_urls.PlayOptions, segments []string, force bool) error {

	videoId := playOptions.VideoId

	gmma := nod.Begin(" generating merge manifest for %s...", videoId)
	defer gmma.End()

	amf, err := absManifestFilename(playOptions)
	if err != nil {
		return err
	}

	if _, err := os.Stat(amf); err == nil && force {
		if err := os.Remove(amf); err != nil {
			return err
		}
	}

	manifestFile, err := os.Create(amf)
	if err != nil {
		return err
	}
	defer manifestFile.Close()

	for _, segment := range segments {
		line := fmt.Sprintf("file '%s'\n", segment)
		if _, err := manifestFile.WriteString(line); err != nil {
			return err
		}
	}

	gmma.EndWithResult("done")

	return nil
}

func mergeVideoSegments(playOptions *rutube_urls.PlayOptions, force bool) error {

	videoId := playOptions.VideoId
	channel := busan.Sanitize(playOptions.Author.Name)
	title := playOptions.Title

	title = strings.Replace(title, "\n", " ", -1)

	mvsa := nod.Begin(" merging video segments for %s, this can take a while...", videoId)
	defer mvsa.End()

	absVideosDir, err := pathways.GetAbsDir(data.Videos)
	if err != nil {
		return err
	}

	absOutputDir := filepath.Join(absVideosDir, channel)

	relOutputFilename := yeti.RelLocalVideoFilename("", title, videoId)
	absOutputFilename := filepath.Join(absVideosDir, channel, relOutputFilename)

	if _, err := os.Stat(absOutputFilename); err == nil && force {
		if err := os.Remove(absOutputFilename); err != nil {
			return err
		}
	}

	ffmb, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}

	args := []string{
		"-f", "concat",
		"-i", relManifestFilename(videoId),
		"-c", "copy", relOutputFilename}

	cmd := exec.Command(ffmb, args...)
	cmd.Dir = absOutputDir
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func removeManifestSegments(playOptions *rutube_urls.PlayOptions, segments []string) error {

	videoId := playOptions.VideoId
	channel := busan.Sanitize(playOptions.Author.Name)

	rmsa := nod.NewProgress(" removing manifest, segments for %s...", videoId)
	defer rmsa.End()

	absVideosDir, err := pathways.GetAbsDir(data.Videos)
	if err != nil {
		return err
	}

	absOutputDir := filepath.Join(absVideosDir, channel)

	amf, err := absManifestFilename(playOptions)
	if err != nil {
		return err
	}

	if err := os.Remove(amf); err != nil {
		return err
	}

	rmsa.TotalInt(len(segments))

	dir := ""

	for _, segment := range segments {

		dir = path.Dir(segment)

		absSegmentFilename := filepath.Join(absOutputDir, segment)
		if err := os.Remove(absSegmentFilename); err != nil {
			return err
		}

		rmsa.Increment()
	}

	if err := os.Remove(filepath.Join(absOutputDir, dir)); err != nil {
		return err
	}

	return nil
}
