package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const updateStateFileName = "pitboardz_update.json"
const githubReleasesURL = "https://api.github.com/repos/MaximusTrue/Pitboardz-Public/releases?per_page=10"
const installerAssetName = "pitboardz.exe"
const maxInstallerDownloadBytes int64 = 100 * 1024 * 1024

type UpdateState struct {
	CurrentVersion string `json:"currentVersion"`
}

type githubRelease struct {
	TagName    string        `json:"tag_name"`
	Draft      bool          `json:"draft"`
	Prerelease bool          `json:"prerelease"`
	Assets     []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// UpdateRelease identifies the installer for a newer published version.
type UpdateRelease struct {
	Version      string
	InstallerURL string
}

func WriteUpdateState(savePath, version string, logger Logger) {
	if savePath == "" {
		return
	}

	state := UpdateState{
		CurrentVersion: version,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		logger("initUpdateState: Failed to encode update state: %v", err)
		return
	}
	data = append(data, '\n')

	updateStatePath := filepath.Join(savePath, updateStateFileName)
	if err := os.WriteFile(updateStatePath, data, 0644); err != nil {
		logger("initUpdateState: Failed to write update state: %v", err)
		return
	}

	logger("initUpdateState: Wrote update state to %s", updateStatePath)
}

func LogLatestVersionComparison(version string, logger Logger) *UpdateRelease {
	client := &http.Client{Timeout: 5 * time.Second}
	request, err := http.NewRequest(http.MethodGet, githubReleasesURL, nil)
	if err != nil {
		logger("Version check: failed to create GitHub request: %v", err)
		return nil
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "PitBoardzz/"+version)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := client.Do(request)
	if err != nil {
		logger("Version check: failed to fetch GitHub releases: %v", err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger("Version check: GitHub returned %s", response.Status)
		return nil
	}

	var releases []githubRelease
	if err := json.NewDecoder(response.Body).Decode(&releases); err != nil {
		logger("Version check: failed to decode GitHub releases: %v", err)
		return nil
	}

	for _, release := range releases {
		if release.Draft {
			continue
		}

		installedVersion := normalizeVersion(version)
		latestVersion := normalizeVersion(release.TagName)
		versionsMatch := installedVersion == latestVersion
		logger(
			"Version check: Installed version: %s | Latest version: %s | Versions match: %t | Latest is prerelease: %t",
			version,
			latestVersion,
			versionsMatch,
			release.Prerelease,
		)
		if versionsMatch {
			return nil
		}

		for _, asset := range release.Assets {
			if strings.EqualFold(asset.Name, installerAssetName) && asset.BrowserDownloadURL != "" {
				return &UpdateRelease{
					Version:      latestVersion,
					InstallerURL: asset.BrowserDownloadURL,
				}
			}
		}
		logger("Version check: %s does not include %s", release.TagName, installerAssetName)
		return nil
	}

	logger("Version check: no published GitHub releases found")
	return nil
}

// DownloadAndLaunchUpdate downloads the release installer and starts it as a
// detached process that waits for MX Bikes to exit before replacing the plugin.
func DownloadAndLaunchUpdate(installerURL string) error {
	if err := validateInstallerURL(installerURL); err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 2 * time.Minute,
		CheckRedirect: func(request *http.Request, previous []*http.Request) error {
			if !isTrustedDownloadHost(request.URL) {
				return fmt.Errorf("untrusted update redirect host %q", request.URL.Hostname())
			}
			return nil
		},
	}
	request, err := http.NewRequest(http.MethodGet, installerURL, nil)
	if err != nil {
		return fmt.Errorf("create installer request: %w", err)
	}
	request.Header.Set("Accept", "application/octet-stream")
	request.Header.Set("User-Agent", "PitBoardzz-Updater")

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("download installer: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download installer: GitHub returned %s", response.Status)
	}
	if response.ContentLength > maxInstallerDownloadBytes {
		return fmt.Errorf("download installer: asset exceeds %d bytes", maxInstallerDownloadBytes)
	}

	installerFile, err := os.CreateTemp("", "pitboardz-update-*.exe")
	if err != nil {
		return fmt.Errorf("create temporary installer: %w", err)
	}
	installerPath := installerFile.Name()
	downloadSucceeded := false
	defer func() {
		_ = installerFile.Close()
		if !downloadSucceeded {
			_ = os.Remove(installerPath)
		}
	}()

	writtenBytes, err := io.Copy(installerFile, io.LimitReader(response.Body, maxInstallerDownloadBytes+1))
	if err != nil {
		return fmt.Errorf("save installer: %w", err)
	}
	if writtenBytes > maxInstallerDownloadBytes {
		return fmt.Errorf("save installer: asset exceeds %d bytes", maxInstallerDownloadBytes)
	}
	if err := installerFile.Close(); err != nil {
		return fmt.Errorf("close installer: %w", err)
	}

	command := exec.Command(installerPath, "-wait-pid", strconv.Itoa(os.Getpid()))
	if err := command.Start(); err != nil {
		return fmt.Errorf("launch installer: %w", err)
	}
	if err := command.Process.Release(); err != nil {
		return fmt.Errorf("detach installer: %w", err)
	}
	downloadSucceeded = true
	return nil
}

func validateInstallerURL(installerURL string) error {
	parsedURL, err := url.Parse(installerURL)
	if err != nil {
		return fmt.Errorf("parse installer URL: %w", err)
	}
	if parsedURL.Scheme != "https" || !isTrustedDownloadHost(parsedURL) {
		return fmt.Errorf("untrusted installer URL %q", installerURL)
	}
	if !strings.EqualFold(filepath.Base(parsedURL.Path), installerAssetName) {
		return fmt.Errorf("unexpected installer asset %q", filepath.Base(parsedURL.Path))
	}
	return nil
}

func isTrustedDownloadHost(downloadURL *url.URL) bool {
	hostname := strings.ToLower(downloadURL.Hostname())
	return downloadURL.Scheme == "https" &&
		(hostname == "github.com" || strings.HasSuffix(hostname, ".githubusercontent.com"))
}

func normalizeVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")
	return version
}
