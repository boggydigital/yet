package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yet/yeti"
	"github.com/boggydigital/yet_urls/rutube_urls"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetRutubeVideoHandler(u *url.URL) error {
	urls := strings.Split(u.Query().Get("url"), ",")
	force := u.Query().Has("force")
	return GetRutubeVideo(force, urls...)
}

func GetRutubeVideo(force bool, urls ...string) error {

	grva := nod.NewProgress("getting Rutube videos...")
	defer grva.End()

	grva.TotalInt(len(urls))

	absVideosDir, err := pathways.GetAbsDir(paths.Videos)
	if err != nil {
		return grva.EndWithError(err)
	}

	dc := dolo.DefaultClient

	for _, u := range urls {

		if err := downloadRutubeVideo(dc, u, absVideosDir, force); err != nil {
			grva.Error(err)
		}

		grva.Increment()
	}

	grva.EndWithResult("done")

	return nil
}

func downloadRutubeVideo(dc *dolo.Client, u string, videosDir string, force bool) error {

	ru, err := url.Parse(u)
	if err != nil {
		return err
	}

	videoId, p := path.Base(ru.Path), ru.Query().Get("p")
	po := rutube_urls.PlayOptionsUrl(videoId, p)

	playOptions, err := getPlayOptions(po)
	if err != nil {
		return err
	}

	channel := playOptions.Author.Name
	title := playOptions.Title
	title = strings.Replace(title, "\n", " ", -1)

	drva := nod.NewProgress(" %s", title)
	defer drva.End()

	formats, err := getVideoBalancerFormats(playOptions.VideoBalancer.Default)
	if err != nil {
		return err
	}

	if len(formats) == 0 {
		return fmt.Errorf("no formats found")
	}

	fu, err := url.Parse(formats[len(formats)-1])
	if err != nil {
		return err
	}

	segments, err := getVideoSegments(fu.String())
	if err != nil {
		return err
	}

	drva.TotalInt(len(segments))

	fuBase := path.Base(fu.Path)

	for _, segment := range segments {

		suStr := strings.Replace(fu.String(), fuBase, segment, 1)
		su, err := url.Parse(suStr)
		if err != nil {
			drva.Error(err)
			continue
		}

		fn := path.Base(su.Path)

		if err := dc.Download(su, force, nil, videosDir, channel, fn); err != nil {
			drva.Error(err)
			continue
		}

		drva.Increment()
	}

	outputDirectory := filepath.Join(videosDir, channel)
	outputFilename := filepath.Join(
		videosDir,
		yeti.ChannelTitleVideoIdFilename(channel, title, videoId))

	tempOutputFilename, err := yeti.MergeSegments(videoId, outputDirectory, segments...)
	if err != nil {
		return err
	}

	if err := os.Rename(tempOutputFilename, outputFilename); err != nil {
		if strings.Contains(err.Error(), "cross-device link") {

			if src, err := os.Open(tempOutputFilename); err == nil {
				defer src.Close()
				if dst, err := os.Create(outputFilename); err == nil {
					defer dst.Close()
					if _, err := io.Copy(dst, src); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func getPlayOptions(playerOptionsUrl *url.URL) (*rutube_urls.PlayOptions, error) {
	resp, err := http.DefaultClient.Get(playerOptionsUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var playOptions rutube_urls.PlayOptions

	if err := json.NewDecoder(resp.Body).Decode(&playOptions); err != nil {
		return nil, err
	}

	return &playOptions, nil
}

func getVideoBalancerFormats(videoBalancerUrl string) ([]string, error) {
	resp, err := http.DefaultClient.Get(videoBalancerUrl)
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

	return formats, nil
}

func getVideoSegments(videoSegmentsUrl string) ([]string, error) {
	resp, err := http.DefaultClient.Get(videoSegmentsUrl)
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

	return segments, nil
}
