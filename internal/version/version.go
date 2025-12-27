package version

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Version struct {
	Type     string
	Version  string
	Path     string
	Commands map[string]ConsoleCommand `yaml:"commands"`
	Events   map[string]ConsoleEvent   `yaml:"events"`
}

type ConsoleCommand struct {
	Pattern string   `yaml:"pattern"`
	Output  []string `yaml:"output"`
	Timeout int      `yaml:"timeout"`
}

type ConsoleEvent struct {
	Output  []string `yaml:"output"`
	Timeout int      `yaml:"timeout"`
}

type VersionConfig struct {
	Extends  string                    `yaml:"extends"`
	Commands map[string]ConsoleCommand `yaml:"commands"`
	Events   map[string]ConsoleEvent   `yaml:"events"`
}

func LoadAll(versioningPath string) ([]*Version, error) {
	var versions []*Version

	err := filepath.Walk(versioningPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		relPath, err := filepath.Rel(versioningPath, path)
		if err != nil {
			return err
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) != 2 {
			return nil
		}

		versionType := parts[0]
		versionNum := strings.TrimSuffix(parts[1], ext)

		v, err := Load(path, versionType, versionNum)
		if err != nil {
			return nil
		}

		versions = append(versions, v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return versions, nil
}

func Load(path, versionType, versionNum string) (*Version, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg VersionConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	v := &Version{
		Type:     versionType,
		Version:  versionNum,
		Path:     path,
		Commands: cfg.Commands,
		Events:   cfg.Events,
	}

	if cfg.Extends != "" {
		baseDir := filepath.Dir(filepath.Dir(path))
		basePath := filepath.Join(baseDir, cfg.Extends+".yaml")
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			basePath = filepath.Join(baseDir, cfg.Extends+".yml")
		}

		baseType := filepath.Dir(cfg.Extends)
		baseVersion := filepath.Base(cfg.Extends)

		base, err := Load(basePath, baseType, baseVersion)
		if err == nil {
			for k, c := range base.Commands {
				if _, exists := v.Commands[k]; !exists {
					if v.Commands == nil {
						v.Commands = make(map[string]ConsoleCommand)
					}
					v.Commands[k] = c
				}
			}
			for k, e := range base.Events {
				if _, exists := v.Events[k]; !exists {
					if v.Events == nil {
						v.Events = make(map[string]ConsoleEvent)
					}
					v.Events[k] = e
				}
			}
		}
	}

	return v, nil
}

func FindClosest(versions []*Version, targetType, targetVersion string) *Version {
	var candidates []*Version
	for _, v := range versions {
		if v.Type == targetType {
			candidates = append(candidates, v)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return compareVersions(candidates[i].Version, candidates[j].Version) < 0
	})

	var closest *Version
	for _, v := range candidates {
		if compareVersions(v.Version, targetVersion) <= 0 {
			closest = v
		} else {
			break
		}
	}

	if closest == nil && len(candidates) > 0 {
		closest = candidates[0]
	}

	return closest
}

func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}

		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}

	return 0
}

func (v *Version) GetCommand(name string) *ConsoleCommand {
	if cmd, ok := v.Commands[name]; ok {
		return &cmd
	}
	return nil
}

func (v *Version) GetEvent(name string) *ConsoleEvent {
	if evt, ok := v.Events[name]; ok {
		return &evt
	}
	return nil
}

func (v *Version) String() string {
	return fmt.Sprintf("%s/%s", v.Type, v.Version)
}
