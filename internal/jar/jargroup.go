package jar

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/msmhq/msm/internal/config"
)

type JarGroup struct {
	Name      string
	Path      string
	URL       string
	Files     []string
	GlobalCfg *config.Config
}

func DiscoverAll(cfg *config.Config) ([]*JarGroup, error) {
	entries, err := os.ReadDir(cfg.JarStoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var groups []*JarGroup
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		groupPath := filepath.Join(cfg.JarStoragePath, entry.Name())
		group, err := Load(groupPath, entry.Name(), cfg)
		if err != nil {
			continue
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func Load(path, name string, cfg *config.Config) (*JarGroup, error) {
	group := &JarGroup{
		Name:      name,
		Path:      path,
		GlobalCfg: cfg,
	}

	urlPath := filepath.Join(path, "url.txt")
	if data, err := os.ReadFile(urlPath); err == nil {
		group.URL = strings.TrimSpace(string(data))
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") {
			group.Files = append(group.Files, entry.Name())
		}
	}

	return group, nil
}

func Get(name string, cfg *config.Config) (*JarGroup, error) {
	groupPath := filepath.Join(cfg.JarStoragePath, name)
	if _, err := os.Stat(groupPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("jar group %q not found", name)
	}
	return Load(groupPath, name, cfg)
}

func Create(name, url string, cfg *config.Config) (*JarGroup, error) {
	groupPath := filepath.Join(cfg.JarStoragePath, name)

	if _, err := os.Stat(groupPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("jar group %q already exists", name)
	}

	if err := os.MkdirAll(groupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create jar group directory: %w", err)
	}

	urlPath := filepath.Join(groupPath, "url.txt")
	if err := os.WriteFile(urlPath, []byte(url+"\n"), 0644); err != nil {
		os.RemoveAll(groupPath)
		return nil, fmt.Errorf("failed to write URL file: %w", err)
	}

	return &JarGroup{
		Name:      name,
		Path:      groupPath,
		URL:       url,
		GlobalCfg: cfg,
	}, nil
}

func Delete(name string, cfg *config.Config) error {
	groupPath := filepath.Join(cfg.JarStoragePath, name)

	if _, err := os.Stat(groupPath); os.IsNotExist(err) {
		return fmt.Errorf("jar group %q not found", name)
	}

	return os.RemoveAll(groupPath)
}

func Rename(oldName, newName string, cfg *config.Config) error {
	oldPath := filepath.Join(cfg.JarStoragePath, oldName)
	newPath := filepath.Join(cfg.JarStoragePath, newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("jar group %q not found", oldName)
	}

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return fmt.Errorf("jar group %q already exists", newName)
	}

	return os.Rename(oldPath, newPath)
}

func (g *JarGroup) ChangeURL(url string) error {
	urlPath := filepath.Join(g.Path, "url.txt")
	if err := os.WriteFile(urlPath, []byte(url+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write URL file: %w", err)
	}
	g.URL = url
	return nil
}

func (g *JarGroup) GetLatest() (string, error) {
	if g.URL == "" {
		return "", fmt.Errorf("no URL configured for jar group %q", g.Name)
	}

	fmt.Printf("Downloading from %s... ", g.URL)

	resp, err := http.Get(g.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	filename := getFilenameFromResponse(resp, g.URL)
	destPath := filepath.Join(g.Path, filename)

	file, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(destPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("Done.")
	return filename, nil
}

func getFilenameFromResponse(resp *http.Response, url string) string {
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if strings.Contains(cd, "filename=") {
			parts := strings.Split(cd, "filename=")
			if len(parts) > 1 {
				return strings.Trim(parts[1], "\"")
			}
		}
	}

	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1]

	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}

	if filename == "" || !strings.HasSuffix(filename, ".jar") {
		filename = "server.jar"
	}

	return filename
}

func (g *JarGroup) LatestFile() string {
	if len(g.Files) == 0 {
		return ""
	}
	return g.Files[len(g.Files)-1]
}

func LinkJar(serverPath, jarPath, jarGroupName, jarFile string, cfg *config.Config) error {
	group, err := Get(jarGroupName, cfg)
	if err != nil {
		return err
	}

	var sourceFile string
	if jarFile != "" {
		found := false
		for _, f := range group.Files {
			if f == jarFile {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("jar file %q not found in group %q", jarFile, jarGroupName)
		}
		sourceFile = jarFile
	} else {
		sourceFile = group.LatestFile()
		if sourceFile == "" {
			return fmt.Errorf("no jar files in group %q", jarGroupName)
		}
	}

	sourcePath := filepath.Join(group.Path, sourceFile)
	destPath := jarPath
	if !filepath.IsAbs(destPath) {
		destPath = filepath.Join(serverPath, jarPath)
	}

	if _, err := os.Lstat(destPath); err == nil {
		if err := os.Remove(destPath); err != nil {
			return fmt.Errorf("failed to remove existing jar: %w", err)
		}
	}

	if err := os.Symlink(sourcePath, destPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}
