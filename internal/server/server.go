package server

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/msmhq/msm/internal/config"
	"github.com/msmhq/msm/internal/logging"
	"github.com/msmhq/msm/pkg/screen"
)

type Server struct {
	Name      string
	Path      string
	Config    *ServerConfig
	GlobalCfg *config.Config
	Screen    *screen.Session
}

type ServerConfig struct {
	Username     string
	ScreenName   string
	JarPath      string
	RAM          int
	Invocation   string
	StopDelay    int
	RestartDelay int

	WorldStoragePath         string
	WorldStorageInactivePath string
	LogPath                  string
	PropertiesPath           string
	AllowlistPath            string
	BannedPlayersPath        string
	BannedIPsPath            string
	OpsPath                  string

	MessageStop                   string
	MessageStopAbort              string
	MessageRestart                string
	MessageRestartAbort           string
	MessageWorldBackupStarted     string
	MessageWorldBackupFinished    string
	MessageCompleteBackupStarted  string
	MessageCompleteBackupFinished string

	FabricEnabled          bool
	FabricLoaderVersion    string
	FabricInstallerVersion string
	MCVersion              string
}

func DiscoverAll(cfg *config.Config) ([]*Server, error) {
	entries, err := os.ReadDir(cfg.ServerStoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read server storage: %w", err)
	}

	var servers []*Server
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		serverPath := filepath.Join(cfg.ServerStoragePath, entry.Name())
		server, err := Load(serverPath, entry.Name(), cfg)
		if err != nil {
			continue
		}
		servers = append(servers, server)
	}

	return servers, nil
}

func Load(path, name string, globalCfg *config.Config) (*Server, error) {
	serverCfg := &ServerConfig{
		Username:                      globalCfg.DefaultUsername,
		ScreenName:                    strings.ReplaceAll(globalCfg.DefaultScreenName, "{SERVER_NAME}", name),
		JarPath:                       globalCfg.DefaultJarPath,
		RAM:                           globalCfg.DefaultRAM,
		Invocation:                    globalCfg.DefaultInvocation,
		StopDelay:                     globalCfg.DefaultStopDelay,
		RestartDelay:                  globalCfg.DefaultRestartDelay,
		WorldStoragePath:              globalCfg.DefaultWorldStoragePath,
		WorldStorageInactivePath:      globalCfg.DefaultWorldStorageInactivePath,
		LogPath:                       globalCfg.DefaultLogPath,
		PropertiesPath:                globalCfg.DefaultPropertiesPath,
		AllowlistPath:                 globalCfg.DefaultAllowlistPath,
		BannedPlayersPath:             globalCfg.DefaultBannedPlayersPath,
		BannedIPsPath:                 globalCfg.DefaultBannedIPsPath,
		OpsPath:                       globalCfg.DefaultOpsPath,
		MessageStop:                   globalCfg.DefaultMessageStop,
		MessageStopAbort:              globalCfg.DefaultMessageStopAbort,
		MessageRestart:                globalCfg.DefaultMessageRestart,
		MessageRestartAbort:           globalCfg.DefaultMessageRestartAbort,
		MessageWorldBackupStarted:     globalCfg.DefaultMessageWorldBackupStarted,
		MessageWorldBackupFinished:    globalCfg.DefaultMessageWorldBackupFinished,
		MessageCompleteBackupStarted:  globalCfg.DefaultMessageCompleteBackupStarted,
		MessageCompleteBackupFinished: globalCfg.DefaultMessageCompleteBackupFinished,
	}

	confPath := filepath.Join(path, "server.conf")
	if file, err := os.Open(confPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
			serverCfg.set(key, value, name)
		}
	}

	return &Server{
		Name:      name,
		Path:      path,
		Config:    serverCfg,
		GlobalCfg: globalCfg,
		Screen:    screen.New(serverCfg.ScreenName),
	}, nil
}

func (c *ServerConfig) set(key, value, serverName string) {
	switch key {
	case "USERNAME":
		c.Username = value
	case "SCREEN_NAME":
		c.ScreenName = strings.ReplaceAll(value, "{SERVER_NAME}", serverName)
	case "JAR_PATH":
		c.JarPath = value
	case "RAM":
		if v, err := strconv.Atoi(value); err != nil {
			logging.Warn("invalid RAM value in server config", "server", serverName, "value", value)
		} else {
			c.RAM = v
		}
	case "INVOCATION":
		c.Invocation = value
	case "STOP_DELAY":
		if v, err := strconv.Atoi(value); err != nil {
			logging.Warn("invalid STOP_DELAY value in server config", "server", serverName, "value", value)
		} else {
			c.StopDelay = v
		}
	case "RESTART_DELAY":
		if v, err := strconv.Atoi(value); err != nil {
			logging.Warn("invalid RESTART_DELAY value in server config", "server", serverName, "value", value)
		} else {
			c.RestartDelay = v
		}
	case "WORLD_STORAGE_PATH":
		c.WorldStoragePath = value
	case "WORLD_STORAGE_INACTIVE_PATH":
		c.WorldStorageInactivePath = value
	case "LOG_PATH":
		c.LogPath = value
	case "PROPERTIES_PATH":
		c.PropertiesPath = value
	case "WHITELIST_PATH":
		c.AllowlistPath = value
	case "BANNED_PLAYERS_PATH":
		c.BannedPlayersPath = value
	case "BANNED_IPS_PATH":
		c.BannedIPsPath = value
	case "OPS_PATH":
		c.OpsPath = value
	case "MESSAGE_STOP":
		c.MessageStop = value
	case "MESSAGE_STOP_ABORT":
		c.MessageStopAbort = value
	case "MESSAGE_RESTART":
		c.MessageRestart = value
	case "MESSAGE_RESTART_ABORT":
		c.MessageRestartAbort = value
	case "MESSAGE_WORLD_BACKUP_STARTED":
		c.MessageWorldBackupStarted = value
	case "MESSAGE_WORLD_BACKUP_FINISHED":
		c.MessageWorldBackupFinished = value
	case "MESSAGE_COMPLETE_BACKUP_STARTED":
		c.MessageCompleteBackupStarted = value
	case "MESSAGE_COMPLETE_BACKUP_FINISHED":
		c.MessageCompleteBackupFinished = value
	case "FABRIC_ENABLED":
		c.FabricEnabled = value == "true"
	case "FABRIC_LOADER_VERSION":
		c.FabricLoaderVersion = value
	case "FABRIC_INSTALLER_VERSION":
		c.FabricInstallerVersion = value
	case "MC_VERSION":
		c.MCVersion = value
	}
}

func Get(name string, cfg *config.Config) (*Server, error) {
	serverPath := filepath.Join(cfg.ServerStoragePath, name)
	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("server %q not found", name)
	}
	return Load(serverPath, name, cfg)
}

func (s *Server) FullPath(relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(s.Path, relativePath)
}

func (s *Server) WorldStoragePath() string {
	return s.FullPath(s.Config.WorldStoragePath)
}

func (s *Server) WorldStorageInactivePath() string {
	return s.FullPath(s.Config.WorldStorageInactivePath)
}

func (s *Server) JarPath() string {
	return s.FullPath(s.Config.JarPath)
}

func (s *Server) LogPath() string {
	return s.FullPath(s.Config.LogPath)
}

func (s *Server) PropertiesPath() string {
	return s.FullPath(s.Config.PropertiesPath)
}

func (s *Server) AllowlistPath() string {
	return s.FullPath(s.Config.AllowlistPath)
}

func (s *Server) BannedPlayersPath() string {
	return s.FullPath(s.Config.BannedPlayersPath)
}

func (s *Server) BannedIPsPath() string {
	return s.FullPath(s.Config.BannedIPsPath)
}

func (s *Server) OpsPath() string {
	return s.FullPath(s.Config.OpsPath)
}

func (s *Server) Port() int {
	return readPortFromProperties(s.PropertiesPath())
}

func (s *Server) CanManage() bool {
	return screen.CanManageUser(s.Config.Username)
}

func (s *Server) CheckPermission() error {
	if s.CanManage() {
		return nil
	}
	currentUser := screen.CurrentUser()
	return fmt.Errorf("permission denied: server %q is owned by user %q\n  You are running as %q\n  Hint: Run with sudo, or as user %q",
		s.Name, s.Config.Username, currentUser, s.Config.Username)
}

func (s *Server) DetectMCVersion() (string, error) {
	if s.Config.MCVersion != "" {
		return s.Config.MCVersion, nil
	}

	jarPath := s.JarPath()
	jarName := filepath.Base(jarPath)

	if info, err := os.Lstat(jarPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(jarPath)
		if err == nil {
			jarName = filepath.Base(target)
		}
	}

	if version := extractVersionFromFilename(jarName); version != "" {
		return version, nil
	}

	return "", fmt.Errorf("could not detect minecraft version for server %q; set MC_VERSION in server.conf", s.Name)
}

func extractVersionFromFilename(filename string) string {
	patterns := []struct {
		prefix string
		suffix string
	}{
		{"minecraft_server.", ".jar"},
		{"server-", ".jar"},
		{"paper-", ".jar"},
		{"spigot-", ".jar"},
		{"craftbukkit-", ".jar"},
		{"purpur-", ".jar"},
		{"fabric-server-mc.", "-loader"},
	}

	for _, p := range patterns {
		if idx := strings.Index(filename, p.prefix); idx != -1 {
			start := idx + len(p.prefix)
			rest := filename[start:]
			if endIdx := strings.Index(rest, p.suffix); endIdx != -1 {
				version := rest[:endIdx]
				version = strings.TrimSuffix(version, "-SNAPSHOT")
				if idx := strings.Index(version, "-"); idx != -1 {
					version = version[:idx]
				}
				if isValidVersion(version) {
					return version
				}
			}
		}
	}

	return ""
}

func isValidVersion(s string) bool {
	if len(s) < 3 {
		return false
	}

	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}

	for _, p := range parts {
		for _, c := range p {
			if c < '0' || c > '9' {
				return false
			}
		}
	}

	return true
}
