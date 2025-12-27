package world

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/msmhq/msm/internal/config"
)

type World struct {
	Name       string
	Path       string
	ServerName string
	ServerPath string
	Active     bool
	InRAM      bool
	RAMPath    string
	GlobalCfg  *config.Config
}

func DiscoverAll(serverPath, serverName string, cfg *config.Config, worldStoragePath, worldStorageInactivePath string) ([]*World, error) {
	var worlds []*World

	activeWorlds, err := discoverInDir(serverPath, serverName, worldStoragePath, true, cfg)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	worlds = append(worlds, activeWorlds...)

	inactiveWorlds, err := discoverInDir(serverPath, serverName, worldStorageInactivePath, false, cfg)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	worlds = append(worlds, inactiveWorlds...)

	return worlds, nil
}

func discoverInDir(serverPath, serverName, storagePath string, active bool, cfg *config.Config) ([]*World, error) {
	fullPath := storagePath
	if !filepath.IsAbs(storagePath) {
		fullPath = filepath.Join(serverPath, storagePath)
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var worlds []*World
	for _, entry := range entries {
		worldPath := filepath.Join(fullPath, entry.Name())

		info, err := os.Stat(worldPath)
		if err != nil || !info.IsDir() {
			continue
		}
		world := &World{
			Name:       entry.Name(),
			Path:       worldPath,
			ServerName: serverName,
			ServerPath: serverPath,
			Active:     active,
			GlobalCfg:  cfg,
		}

		flagPath := filepath.Join(worldPath, "in_ram")
		if _, err := os.Stat(flagPath); err == nil {
			world.InRAM = true
			world.RAMPath = filepath.Join(cfg.RamdiskStoragePath, serverName, entry.Name())
		}

		worlds = append(worlds, world)
	}

	return worlds, nil
}

func Get(serverPath, serverName, worldName, worldStoragePath, worldStorageInactivePath string, cfg *config.Config) (*World, error) {
	activeFullPath := worldStoragePath
	if !filepath.IsAbs(worldStoragePath) {
		activeFullPath = filepath.Join(serverPath, worldStoragePath)
	}

	worldPath := filepath.Join(activeFullPath, worldName)
	if info, err := os.Stat(worldPath); err == nil && info.IsDir() {
		world := &World{
			Name:       worldName,
			Path:       worldPath,
			ServerName: serverName,
			ServerPath: serverPath,
			Active:     true,
			GlobalCfg:  cfg,
		}

		flagPath := filepath.Join(worldPath, "in_ram")
		if _, err := os.Stat(flagPath); err == nil {
			world.InRAM = true
			world.RAMPath = filepath.Join(cfg.RamdiskStoragePath, serverName, worldName)
		}

		return world, nil
	}

	inactiveFullPath := worldStorageInactivePath
	if !filepath.IsAbs(worldStorageInactivePath) {
		inactiveFullPath = filepath.Join(serverPath, worldStorageInactivePath)
	}

	worldPath = filepath.Join(inactiveFullPath, worldName)
	if info, err := os.Stat(worldPath); err == nil && info.IsDir() {
		return &World{
			Name:       worldName,
			Path:       worldPath,
			ServerName: serverName,
			ServerPath: serverPath,
			Active:     false,
			GlobalCfg:  cfg,
		}, nil
	}

	return nil, fmt.Errorf("world %q not found", worldName)
}

func (w *World) FlagPath() string {
	return filepath.Join(w.Path, "in_ram")
}

func (w *World) Activate(worldStoragePath string) error {
	if w.Active {
		return fmt.Errorf("world %q is already active", w.Name)
	}

	activePath := worldStoragePath
	if !filepath.IsAbs(worldStoragePath) {
		activePath = filepath.Join(w.ServerPath, worldStoragePath)
	}

	newPath := filepath.Join(activePath, w.Name)

	if err := os.MkdirAll(activePath, 0755); err != nil {
		return fmt.Errorf("failed to create world storage directory: %w", err)
	}

	if err := os.Rename(w.Path, newPath); err != nil {
		return fmt.Errorf("failed to move world: %w", err)
	}

	w.Path = newPath
	w.Active = true
	return nil
}

func (w *World) Deactivate(worldStorageInactivePath string) error {
	if !w.Active {
		return fmt.Errorf("world %q is already inactive", w.Name)
	}

	inactivePath := worldStorageInactivePath
	if !filepath.IsAbs(worldStorageInactivePath) {
		inactivePath = filepath.Join(w.ServerPath, worldStorageInactivePath)
	}

	newPath := filepath.Join(inactivePath, w.Name)

	if err := os.MkdirAll(inactivePath, 0755); err != nil {
		return fmt.Errorf("failed to create inactive world storage directory: %w", err)
	}

	if err := os.Rename(w.Path, newPath); err != nil {
		return fmt.Errorf("failed to move world: %w", err)
	}

	w.Path = newPath
	w.Active = false
	return nil
}

func (w *World) Status() string {
	status := "inactive"
	if w.Active {
		status = "active"
	}
	if w.InRAM {
		status += ", in RAM"
	}
	return status
}
