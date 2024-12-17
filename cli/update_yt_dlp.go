package cli

import (
	"encoding/json"
	"errors"
	"github.com/arelate/southern_light/github_integration"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/yet/data"
	"github.com/boggydigital/yet/yeti"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	ytDlpOwnerRepo     = "yt-dlp/yt-dlp"
	relYtDlpPluginsDir = ".yt-dlp/plugins"
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
	defer uyda.EndWithResult("done")

	metadataDir, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return uyda.EndWithError(err)
	}

	rdx, err := kevlar.NewReduxWriter(metadataDir, data.YtDlpLatestDownloadedVersionProperty)
	if err != nil {
		return uyda.EndWithError(err)
	}

	ytDlpDir, err := pathways.GetAbsDir(data.YtDlp)
	if err != nil {
		return uyda.EndWithError(err)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return uyda.EndWithError(err)
	}

	ytDlpPluginsDir := filepath.Join(userHomeDir, relYtDlpPluginsDir)

	dc := dolo.DefaultClient

	// update yt-dlp
	ytDlpAsset := yeti.GetYtDlpBinary()
	if err := getAsset(ytDlpOwnerRepo, ytDlpAsset, ytDlpDir, dc, rdx, force); err != nil {
		return uyda.EndWithError(err)
	}

	ytDlpBinaryFilename := filepath.Join(ytDlpDir, ytDlpAsset)
	if err := os.Chmod(ytDlpBinaryFilename, 0555); err != nil {
		return uyda.EndWithError(err)
	}

	// update yt-dlp-get-pot
	if err := getAsset(ytDlpGetPotOwnerRepo, ytDlpGetPotAsset, ytDlpDir, dc, rdx, force); err != nil {
		return uyda.EndWithError(err)
	}

	if err := copyYtDlpPlugin(ytDlpDir, ytDlpPluginsDir, ytDlpGetPotAsset, force); err != nil {
		return err
	}

	// update bgutil-ytdlp-pot-provider
	if err := getAsset(ytDlpPotProviderOwnerRepo, ytDlpPotProviderAsset, ytDlpDir, dc, rdx, force); err != nil {
		return uyda.EndWithError(err)
	}

	if err := copyYtDlpPlugin(ytDlpDir, ytDlpPluginsDir, ytDlpPotProviderAsset, force); err != nil {
		return err
	}

	return nil
}

func getAsset(ownerRepo, asset string, downloadDir string, dc *dolo.Client, rdx kevlar.WriteableRedux, force bool) error {

	gaa := nod.Begin(" getting %s asset...", ownerRepo)
	defer gaa.EndWithResult("done")

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
	defer gra.EndWithResult("done")

	owner, repo, ok := strings.Cut(ownerRepo, "/")
	if !ok {
		return nil, gra.EndWithError(errors.New("invalid owner/repo " + ownerRepo))
	}

	ytDlpReleasesUrl := github_integration.ReleasesUrl(owner, repo)

	resp, err := http.Get(ytDlpReleasesUrl.String())
	if err != nil {
		return nil, gra.EndWithError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, gra.EndWithError(errors.New(resp.Status))
	}

	var releases []github_integration.GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, gra.EndWithError(err)
	}

	if len(releases) > 0 {
		latestRelease := &releases[0]
		gra.EndWithResult("found %s", latestRelease.Name)
		return latestRelease, nil
	}

	return nil, gra.EndWithError(errors.New("latest release not found for " + ownerRepo))
}

func downloadAsset(dstDir string, release *github_integration.GitHubRelease, assetName string, dc *dolo.Client, force bool) error {

	daa := nod.NewProgress(" downloading %s...", assetName)
	defer daa.EndWithResult("done")

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
		return daa.EndWithError(errors.New("cannot locate asset in the provided release"))
	}

	assetUrl, err := url.Parse(desiredAsset.BrowserDownloadUrl)
	if err != nil {
		return daa.EndWithError(err)
	}

	return dc.Download(assetUrl, force, daa, dstDir)
}

func copyYtDlpPlugin(srcDir, dstDir, pluginFilename string, force bool) error {

	cpa := nod.Begin(" copying yt-dlp plugin %s...", pluginFilename)
	defer cpa.EndWithResult("done")

	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return cpa.EndWithError(err)
		}
	}

	dstFilename := filepath.Join(dstDir, pluginFilename)

	if _, err := os.Stat(dstFilename); err == nil {
		if force {
			if err := os.Remove(dstFilename); err != nil {
				return cpa.EndWithError(err)
			}
		} else {
			cpa.EndWithResult("already exists")
			return nil
		}
	}

	srcFilename := filepath.Join(srcDir, pluginFilename)
	srcFile, err := os.Open(srcFilename)
	if err != nil {
		return cpa.EndWithError(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstFilename)
	if err != nil {
		return cpa.EndWithError(err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return cpa.EndWithError(err)
	}

	return nil
}
