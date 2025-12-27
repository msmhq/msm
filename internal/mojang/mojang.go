package mojang

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

const VersionManifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type VersionManifest struct {
	Latest   LatestVersions `json:"latest"`
	Versions []Version      `json:"versions"`
}

type LatestVersions struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type Version struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type VersionDetails struct {
	ID        string    `json:"id"`
	Downloads Downloads `json:"downloads"`
}

type Downloads struct {
	Server DownloadInfo `json:"server"`
}

type DownloadInfo struct {
	SHA1 string `json:"sha1"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func GetVersionManifest() (*VersionManifest, error) {
	resp, err := http.Get(VersionManifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch version manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("version manifest request failed with status %d", resp.StatusCode)
	}

	var manifest VersionManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse version manifest: %w", err)
	}

	return &manifest, nil
}

func GetLatestRelease() (string, error) {
	manifest, err := GetVersionManifest()
	if err != nil {
		return "", err
	}
	return manifest.Latest.Release, nil
}

func GetVersionDetails(versionID string) (*VersionDetails, error) {
	manifest, err := GetVersionManifest()
	if err != nil {
		return nil, err
	}

	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == versionID {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		return nil, fmt.Errorf("version %q not found", versionID)
	}

	resp, err := http.Get(versionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch version details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("version details request failed with status %d", resp.StatusCode)
	}

	var details VersionDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to parse version details: %w", err)
	}

	return &details, nil
}

func JarFilename(version string) string {
	return fmt.Sprintf("minecraft_server.%s.jar", version)
}

func CachedJarPath(cacheDir, version string) string {
	return filepath.Join(cacheDir, "minecraft", JarFilename(version))
}

func EnsureCached(cacheDir, version string) (string, error) {
	jarPath := CachedJarPath(cacheDir, version)

	if _, err := os.Stat(jarPath); err == nil {
		return jarPath, nil
	}

	details, err := GetVersionDetails(version)
	if err != nil {
		return "", err
	}

	if details.Downloads.Server.URL == "" {
		return "", fmt.Errorf("no server download available for version %q", version)
	}

	if err := mkdirAllWithOwner(filepath.Dir(jarPath)); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	resp, err := http.Get(details.Downloads.Server.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download server jar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server download failed with status %d", resp.StatusCode)
	}

	file, err := os.Create(jarPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(jarPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	chownToDefaultUser(jarPath)

	return jarPath, nil
}

func mkdirAllWithOwner(path string) error {
	if err := os.MkdirAll(path, 0775); err != nil {
		return err
	}
	chownToDefaultUser(path)
	parentDir := filepath.Dir(path)
	if parentDir != path && parentDir != "/" {
		chownToDefaultUser(parentDir)
	}
	return nil
}

func chownToDefaultUser(path string) {
	if syscall.Getuid() != 0 {
		return
	}

	u, err := user.Lookup("minecraft")
	if err != nil {
		return
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return
	}

	os.Chown(path, uid, gid)
}
