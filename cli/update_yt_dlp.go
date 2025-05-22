package cli

import (
	"encoding/json"
	"errors"
	"github.com/arelate/southern_light/github_integration"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	ytDlpOwnerRepo = "yt-dlp/yt-dlp"
)

const (
	ytDlpGetPotOwnerRepo = "coletdjnz/yt-dlp-get-pot"
	ytDlpGetPotAsset     = "yt-dlp-get-pot.zip"
)

const (
	ytDlpPotProviderOwnerRepo = "Brainicism/bgutil-ytdlp-pot-provider"
	ytDlpPotProviderAsset     = "bgutil-ytdlp-pot-provider.zip"
)

func UpdateYtDlpHandler(u *url.URL) error {

	force := u.Query().Has("force")

	return UpdateYtDlp(force)
}

func UpdateYtDlp(force bool) error {

	uyda := nod.Begin("updating yt-dlp and plugins...")
	defer uyda.Done()

	metadataDir, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return err
	}

	rdx, err := redux.NewWriter(metadataDir, data.YtDlpLatestDownloadedVersionProperty)
	if err != nil {
		return err
	}

	ytDlpDir, err := pathways.GetAbsDir(data.YtDlp)
	if err != nil {
		return err
	}

	absYtDlpPluginsDir, err := pathways.GetAbsRelDir(data.YtDlpPlugins)
	if err != nil {
		return err
	}

	dc := dolo.DefaultClient

	// update yt-dlp
	ytDlpAsset := yeti.GetYtDlpBinary()
	if err := getAsset(ytDlpOwnerRepo, ytDlpAsset, ytDlpDir, dc, rdx, force); err != nil {
		return err
	}

	ytDlpBinaryFilename := filepath.Join(ytDlpDir, ytDlpAsset)
	if err := os.Chmod(ytDlpBinaryFilename, 0555); err != nil {
		return err
	}

	// update yt-dlp-get-pot
	if err := getAsset(ytDlpGetPotOwnerRepo, ytDlpGetPotAsset, absYtDlpPluginsDir, dc, rdx, force); err != nil {
		return err
	}

	//if err := copyYtDlpPlugin(ytDlpDir, absYtDlpPluginsDir, ytDlpGetPotAsset, force); err != nil {
	//	return err
	//}

	// update bgutil-ytdlp-pot-provider
	if err := getAsset(ytDlpPotProviderOwnerRepo, ytDlpPotProviderAsset, absYtDlpPluginsDir, dc, rdx, force); err != nil {
		return err
	}

	//if err := copyYtDlpPlugin(ytDlpDir, absYtDlpPluginsDir, ytDlpPotProviderAsset, force); err != nil {
	//	return err
	//}

	return nil
}

func getAsset(ownerRepo, asset string, downloadDir string, dc *dolo.Client, rdx redux.Writeable, force bool) error {

	gaa := nod.Begin(" getting %s asset...", ownerRepo)
	defer gaa.Done()

	latestRelease, err := getLatestGitHubRelease(ownerRepo)
	if err != nil {
		return err
	}

	updateAsset := false

	if ldv, ok := rdx.GetLastVal(data.YtDlpLatestDownloadedVersionProperty, ownerRepo); ok {
		if ldv != latestRelease.TagName || force {
			updateAsset = true
		}
	} else {
		updateAsset = true
	}

	if updateAsset {
		if err := downloadAsset(downloadDir, latestRelease, asset, dc, updateAsset); err != nil {
			return err
		}

		if err := rdx.ReplaceValues(data.YtDlpLatestDownloadedVersionProperty, ownerRepo, latestRelease.TagName); err != nil {
			return err
		}
	} else {
		gaa.EndWithResult("already got the latest version")
	}

	return nil
}

func getLatestGitHubRelease(ownerRepo string) (*github_integration.GitHubRelease, error) {

	gra := nod.Begin(" getting the latest GitHub release for %s...", ownerRepo)
	defer gra.Done()

	owner, repo, ok := strings.Cut(ownerRepo, "/")
	if !ok {
		return nil, errors.New("invalid owner/repo " + ownerRepo)
	}

	ytDlpReleasesUrl := github_integration.ReleasesUrl(owner, repo)

	resp, err := http.Get(ytDlpReleasesUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(resp.Status)
	}

	var releases []github_integration.GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	if len(releases) > 0 {
		latestRelease := &releases[0]
		gra.EndWithResult("found %s", latestRelease.Name)
		return latestRelease, nil
	}

	return nil, errors.New("latest release not found for " + ownerRepo)
}

func downloadAsset(dstDir string, release *github_integration.GitHubRelease, assetName string, dc *dolo.Client, force bool) error {

	daa := nod.NewProgress(" downloading %s...", assetName)
	defer daa.Done()

	dstAssetFilename := filepath.Join(dstDir, assetName)
	if _, err := os.Stat(dstAssetFilename); err == nil && !force {
		daa.EndWithResult("already exists")
		return nil
	}

	var desiredAsset *github_integration.GitHubAsset

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			desiredAsset = &asset
			break
		}
	}

	if desiredAsset == nil {
		return errors.New("cannot locate asset in the provided release")
	}

	assetUrl, err := url.Parse(desiredAsset.BrowserDownloadUrl)
	if err != nil {
		return err
	}

	return dc.Download(assetUrl, force, daa, dstDir)
}

//func copyYtDlpPlugin(srcDir, dstDir, pluginFilename string, force bool) error {
//
//	cpa := nod.Begin(" copying yt-dlp plugin %s...", pluginFilename)
//	defer cpa.Done()
//
//	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
//		if err := os.MkdirAll(dstDir, 0755); err != nil {
//			return err
//		}
//	}
//
//	dstFilename := filepath.Join(dstDir, pluginFilename)
//
//	if _, err := os.Stat(dstFilename); err == nil {
//		if force {
//			if err := os.Remove(dstFilename); err != nil {
//				return err
//			}
//		} else {
//			cpa.EndWithResult("already exists")
//			return nil
//		}
//	}
//
//	srcFilename := filepath.Join(srcDir, pluginFilename)
//	srcFile, err := os.Open(srcFilename)
//	if err != nil {
//		return err
//	}
//	defer srcFile.Close()
//
//	dstFile, err := os.Create(dstFilename)
//	if err != nil {
//		return err
//	}
//	defer dstFile.Close()
//
//	if _, err := io.Copy(dstFile, srcFile); err != nil {
//		return err
//	}
//
//	return nil
//}
