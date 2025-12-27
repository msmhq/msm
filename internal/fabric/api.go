package fabric

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/msmhq/msm/internal/logging"
)

const BaseURL = "https://meta.fabricmc.net/v2"

type Client struct {
	cache       *Cache
	storagePath string
	ttl         time.Duration
	httpClient  *http.Client
}

func NewClient(storagePath string, cacheTTLMinutes int) (*Client, error) {
	cachePath := CachePath(storagePath)
	cache, err := LoadCache(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load fabric cache: %w", err)
	}

	return &Client{
		cache:       cache,
		storagePath: storagePath,
		ttl:         time.Duration(cacheTTLMinutes) * time.Minute,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) SaveCache() error {
	return c.cache.Save(CachePath(c.storagePath))
}

func (c *Client) FetchGameVersions() ([]GameVersion, error) {
	if versions, ok := c.cache.GetGameVersions(c.ttl); ok {
		return versions, nil
	}

	logging.Debug("fetching game versions from Fabric API")

	resp, err := c.httpClient.Get(BaseURL + "/versions/game")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fabric API returned status %d", resp.StatusCode)
	}

	var versions []GameVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode game versions: %w", err)
	}

	c.cache.SetGameVersions(versions)
	if err := c.SaveCache(); err != nil {
		logging.Warn("failed to save fabric cache", "error", err)
	}

	return versions, nil
}

func (c *Client) FetchLoaderVersions(mcVersion string) ([]LoaderVersion, error) {
	if versions, ok := c.cache.GetLoaderVersions(mcVersion, c.ttl); ok {
		return versions, nil
	}

	logging.Debug("fetching loader versions from Fabric API", "mcVersion", mcVersion)

	url := fmt.Sprintf("%s/versions/loader/%s", BaseURL, mcVersion)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch loader versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("minecraft version %s is not supported by Fabric", mcVersion)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fabric API returned status %d", resp.StatusCode)
	}

	var loaderResp []struct {
		Loader LoaderVersion `json:"loader"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loaderResp); err != nil {
		return nil, fmt.Errorf("failed to decode loader versions: %w", err)
	}

	versions := make([]LoaderVersion, len(loaderResp))
	for i, lr := range loaderResp {
		versions[i] = lr.Loader
	}

	c.cache.SetLoaderVersions(mcVersion, versions)
	if err := c.SaveCache(); err != nil {
		logging.Warn("failed to save fabric cache", "error", err)
	}

	return versions, nil
}

func (c *Client) FetchInstallers() ([]InstallerVersion, error) {
	if installers, ok := c.cache.GetInstallers(c.ttl); ok {
		return installers, nil
	}

	logging.Debug("fetching installer versions from Fabric API")

	resp, err := c.httpClient.Get(BaseURL + "/versions/installer")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch installer versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fabric API returned status %d", resp.StatusCode)
	}

	var installers []InstallerVersion
	if err := json.NewDecoder(resp.Body).Decode(&installers); err != nil {
		return nil, fmt.Errorf("failed to decode installer versions: %w", err)
	}

	c.cache.SetInstallers(installers)
	if err := c.SaveCache(); err != nil {
		logging.Warn("failed to save fabric cache", "error", err)
	}

	return installers, nil
}

func (c *Client) SupportsVersion(mcVersion string) (bool, error) {
	versions, err := c.FetchGameVersions()
	if err != nil {
		return false, err
	}

	for _, v := range versions {
		if v.Version == mcVersion {
			return true, nil
		}
	}

	return false, nil
}

func (c *Client) GetLatestStableLoader(mcVersion string) (*LoaderVersion, error) {
	versions, err := c.FetchLoaderVersions(mcVersion)
	if err != nil {
		return nil, err
	}

	for _, v := range versions {
		if v.Stable {
			return &v, nil
		}
	}

	if len(versions) > 0 {
		return &versions[0], nil
	}

	return nil, fmt.Errorf("no loader versions available for minecraft %s", mcVersion)
}

func (c *Client) GetLatestStableInstaller() (*InstallerVersion, error) {
	installers, err := c.FetchInstallers()
	if err != nil {
		return nil, err
	}

	for _, i := range installers {
		if i.Stable {
			return &i, nil
		}
	}

	if len(installers) > 0 {
		return &installers[0], nil
	}

	return nil, fmt.Errorf("no installer versions available")
}

func (c *Client) DownloadServerJar(mcVersion, loaderVersion, installerVersion string) (string, error) {
	jarPath := JarPath(c.storagePath, mcVersion, loaderVersion, installerVersion)

	if _, err := os.Stat(jarPath); err == nil {
		logging.Debug("fabric jar already exists", "path", jarPath)
		return jarPath, nil
	}

	jarsDir := filepath.Dir(jarPath)
	if err := mkdirAllWithOwner(jarsDir); err != nil {
		return "", fmt.Errorf("failed to create jars directory: %w", err)
	}

	url := fmt.Sprintf("%s/versions/loader/%s/%s/%s/server/jar",
		BaseURL, mcVersion, loaderVersion, installerVersion)

	logging.Info("downloading fabric server jar",
		"mcVersion", mcVersion,
		"loaderVersion", loaderVersion,
		"installerVersion", installerVersion)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download fabric jar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fabric API returned status %d when downloading jar", resp.StatusCode)
	}

	tmpPath := jarPath + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to create jar file: %w", err)
	}

	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to write jar file: %w", err)
	}

	if err := os.Rename(tmpPath, jarPath); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to finalize jar file: %w", err)
	}

	chownToDefaultUser(jarPath)

	logging.Info("fabric server jar downloaded", "path", jarPath)
	return jarPath, nil
}

func (c *Client) ResolveVersions(mcVersion, loaderVersion, installerVersion string) (string, string, error) {
	if loaderVersion == "" {
		loader, err := c.GetLatestStableLoader(mcVersion)
		if err != nil {
			return "", "", err
		}
		loaderVersion = loader.Version
	}

	if installerVersion == "" {
		installer, err := c.GetLatestStableInstaller()
		if err != nil {
			return "", "", err
		}
		installerVersion = installer.Version
	}

	return loaderVersion, installerVersion, nil
}

func (c *Client) GetOrDownloadJar(mcVersion, loaderVersion, installerVersion string) (string, error) {
	loaderVersion, installerVersion, err := c.ResolveVersions(mcVersion, loaderVersion, installerVersion)
	if err != nil {
		return "", err
	}

	jarPath := JarPath(c.storagePath, mcVersion, loaderVersion, installerVersion)

	if _, err := os.Stat(jarPath); err == nil {
		return jarPath, nil
	}

	return c.DownloadServerJar(mcVersion, loaderVersion, installerVersion)
}
