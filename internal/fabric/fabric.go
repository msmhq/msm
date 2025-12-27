package fabric

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type GameVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

type LoaderVersion struct {
	Version string `json:"version"`
	Build   int    `json:"build"`
	Maven   string `json:"maven"`
	Stable  bool   `json:"stable"`
}

type InstallerVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	URL     string `json:"url"`
	Maven   string `json:"maven"`
}

type Cache struct {
	mu sync.RWMutex

	GameVersions   []GameVersion         `json:"game_versions"`
	GameVersionsAt time.Time             `json:"game_versions_at"`
	LoaderVersions map[string]cacheEntry `json:"loader_versions"`
	Installers     []InstallerVersion    `json:"installers"`
	InstallersAt   time.Time             `json:"installers_at"`
}

type cacheEntry struct {
	Versions  []LoaderVersion `json:"versions"`
	FetchedAt time.Time       `json:"fetched_at"`
}

func NewCache() *Cache {
	return &Cache{
		LoaderVersions: make(map[string]cacheEntry),
	}
}

func LoadCache(path string) (*Cache, error) {
	cache := NewCache()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cache, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return cache, nil
	}

	if cache.LoaderVersions == nil {
		cache.LoaderVersions = make(map[string]cacheEntry)
	}

	return cache, nil
}

func (c *Cache) Save(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if err := mkdirAllWithOwner(filepath.Dir(path)); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0664); err != nil {
		return err
	}

	chownToDefaultUser(path)
	return nil
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

func (c *Cache) GetGameVersions(ttl time.Duration) ([]GameVersion, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Since(c.GameVersionsAt) > ttl {
		return nil, false
	}
	return c.GameVersions, true
}

func (c *Cache) SetGameVersions(versions []GameVersion) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.GameVersions = versions
	c.GameVersionsAt = time.Now()
}

func (c *Cache) GetLoaderVersions(mcVersion string, ttl time.Duration) ([]LoaderVersion, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.LoaderVersions[mcVersion]
	if !ok || time.Since(entry.FetchedAt) > ttl {
		return nil, false
	}
	return entry.Versions, true
}

func (c *Cache) SetLoaderVersions(mcVersion string, versions []LoaderVersion) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.LoaderVersions[mcVersion] = cacheEntry{
		Versions:  versions,
		FetchedAt: time.Now(),
	}
}

func (c *Cache) GetInstallers(ttl time.Duration) ([]InstallerVersion, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Since(c.InstallersAt) > ttl {
		return nil, false
	}
	return c.Installers, true
}

func (c *Cache) SetInstallers(installers []InstallerVersion) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Installers = installers
	c.InstallersAt = time.Now()
}

func JarFilename(mcVersion, loaderVersion, installerVersion string) string {
	return "fabric-server-mc." + mcVersion + "-loader." + loaderVersion + "-launcher." + installerVersion + ".jar"
}

func JarPath(storagePath, mcVersion, loaderVersion, installerVersion string) string {
	return filepath.Join(storagePath, "jars", JarFilename(mcVersion, loaderVersion, installerVersion))
}

func CachePath(storagePath string) string {
	return filepath.Join(storagePath, "cache.json")
}
