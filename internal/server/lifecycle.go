package server

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/msmhq/msm/internal/config"
	"github.com/msmhq/msm/internal/fabric"
	"github.com/msmhq/msm/internal/logging"
	"github.com/msmhq/msm/internal/world"
	"github.com/msmhq/msm/pkg/screen"
)

func (s *Server) IsRunning() bool {
	return s.Screen.IsRunningAsUser(s.Config.Username)
}

func (s *Server) Start() error {
	if err := s.CheckPermission(); err != nil {
		return err
	}

	if s.IsRunning() {
		return fmt.Errorf("server %q is already running", s.Name)
	}

	fmt.Fprintf(os.Stderr, "Starting server %q...\n", s.Name)

	var jarPath string
	var jarForInvocation string

	if s.Config.FabricEnabled {
		fabricJar, err := s.ResolveFabricJar()
		if err != nil {
			return fmt.Errorf("fabric jar resolution failed: %w", err)
		}
		jarPath = fabricJar
		jarForInvocation = fabricJar
	} else {
		jarPath = s.JarPath()
		jarForInvocation = s.Config.JarPath
	}

	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		return fmt.Errorf("jar file not found: %s", jarPath)
	}

	if err := s.SetupWorldSymlinks(); err != nil {
		return fmt.Errorf("failed to set up world symlinks: %w", err)
	}

	if err := s.SetupRAMWorlds(); err != nil {
		return fmt.Errorf("failed to set up RAM worlds: %w", err)
	}

	invocation := s.Config.Invocation
	invocation = strings.ReplaceAll(invocation, "{RAM}", strconv.Itoa(s.Config.RAM))
	invocation = strings.ReplaceAll(invocation, "{JAR}", jarForInvocation)

	if err := s.Screen.Start(s.Path, invocation, s.Config.Username); err != nil {
		return err
	}

	if hasRAM, _ := AnyRAMWorldsConfigured(s.GlobalCfg); hasRAM {
		if err := StartSyncDaemon(s.GlobalCfg); err != nil {
			logging.Warn("failed to start sync daemon", "error", err)
		}
	}

	return nil
}

func (s *Server) ResolveFabricJar() (string, error) {
	mcVersion, err := s.DetectMCVersion()
	if err != nil {
		return "", err
	}

	client, err := fabric.NewClient(s.GlobalCfg.FabricStoragePath, s.GlobalCfg.FabricCacheTTL)
	if err != nil {
		return "", fmt.Errorf("failed to create fabric client: %w", err)
	}

	loaderVersion := s.Config.FabricLoaderVersion
	installerVersion := s.Config.FabricInstallerVersion

	loaderVersion, installerVersion, err = client.ResolveVersions(mcVersion, loaderVersion, installerVersion)
	if err != nil {
		return "", err
	}

	jarPath := fabric.JarPath(s.GlobalCfg.FabricStoragePath, mcVersion, loaderVersion, installerVersion)

	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		logging.Info("downloading fabric launcher",
			"mcVersion", mcVersion,
			"loaderVersion", loaderVersion,
			"installerVersion", installerVersion)
		jarPath, err = client.DownloadServerJar(mcVersion, loaderVersion, installerVersion)
		if err != nil {
			return "", err
		}
	}

	return jarPath, nil
}

func (s *Server) SetupWorldSymlinks() error {
	levelName := s.getLevelName()
	if levelName == "" {
		levelName = "world"
	}

	worldStoragePath := s.WorldStoragePath()
	worldInStorage := filepath.Join(worldStoragePath, levelName)
	symlinkPath := filepath.Join(s.Path, levelName)

	storageInfo, err := os.Stat(worldInStorage)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !storageInfo.IsDir() {
		return nil
	}

	symlinkInfo, err := os.Lstat(symlinkPath)
	if err == nil {
		if symlinkInfo.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(symlinkPath)
			if err == nil && target == worldInStorage {
				return nil
			}
			if err := os.Remove(symlinkPath); err != nil {
				return fmt.Errorf("failed to remove existing symlink: %w", err)
			}
		} else if symlinkInfo.IsDir() {
			entries, _ := os.ReadDir(symlinkPath)
			if len(entries) == 0 {
				if err := os.Remove(symlinkPath); err != nil {
					return fmt.Errorf("failed to remove empty world directory: %w", err)
				}
			} else {
				logging.Warn("world directory exists in server root with data, not creating symlink",
					"path", symlinkPath, "hint", "move contents to worldstorage/"+levelName)
				return nil
			}
		}
	}

	logging.Debug("creating world symlink", "from", symlinkPath, "to", worldInStorage)
	if err := os.Symlink(worldInStorage, symlinkPath); err != nil {
		return fmt.Errorf("failed to create world symlink: %w", err)
	}

	return nil
}

func (s *Server) getLevelName() string {
	propsPath := s.PropertiesPath()
	file, err := os.Open(propsPath)
	if err != nil {
		return "world"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "level-name=") {
			return strings.TrimPrefix(line, "level-name=")
		}
	}
	return "world"
}

func (s *Server) SetupRAMWorlds() error {
	if !s.GlobalCfg.RamdiskStorageEnabled {
		return nil
	}

	worlds, err := world.DiscoverAll(
		s.Path,
		s.Name,
		s.GlobalCfg,
		s.Config.WorldStoragePath,
		s.Config.WorldStorageInactivePath,
	)
	if err != nil {
		return err
	}

	for _, w := range worlds {
		if w.InRAM && w.Active {
			logging.Debug("setting up RAM symlink", "world", w.Name, "server", s.Name)
			if err := w.SetupRAMSymlink(s.Config.WorldStoragePath); err != nil {
				return fmt.Errorf("failed to set up RAM symlink for world %q: %w", w.Name, err)
			}
		}
	}

	return nil
}

func (s *Server) SyncRAMWorldsToDisk() error {
	if !s.GlobalCfg.RamdiskStorageEnabled {
		return nil
	}

	worlds, err := world.DiscoverAll(
		s.Path,
		s.Name,
		s.GlobalCfg,
		s.Config.WorldStoragePath,
		s.Config.WorldStorageInactivePath,
	)
	if err != nil {
		return err
	}

	hasRAMWorlds := false
	for _, w := range worlds {
		if w.InRAM && w.Active {
			hasRAMWorlds = true
			break
		}
	}

	if !hasRAMWorlds {
		return nil
	}

	fmt.Fprintf(os.Stderr, "Syncing RAM worlds to disk for server %q...\n", s.Name)

	if s.IsRunning() {
		s.SaveOff()
		s.SaveAll()
		time.Sleep(1 * time.Second)
	}

	for _, w := range worlds {
		if w.InRAM && w.Active {
			logging.Debug("syncing world to disk", "world", w.Name, "server", s.Name)
			if err := w.ToDisk(s.Config.Username); err != nil {
				logging.Warn("failed to sync world", "world", w.Name, "error", err)
			}
		}
	}

	return nil
}

func (s *Server) Stop(immediate bool) error {
	if err := s.CheckPermission(); err != nil {
		return err
	}

	if !s.IsRunning() {
		return nil
	}

	fmt.Fprintf(os.Stderr, "Stopping server %q...\n", s.Name)

	if err := s.SyncRAMWorldsToDisk(); err != nil {
		logging.Warn("failed to sync RAM worlds before stop", "server", s.Name, "error", err)
	}

	if !immediate {
		delay := s.Config.StopDelay
		msg := strings.ReplaceAll(s.Config.MessageStop, "{DELAY}", strconv.Itoa(delay))

		fmt.Fprintf(os.Stderr, "Warning players (%d second delay)...\n", delay)
		s.Say(msg)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	if err := s.SendCommand("stop"); err != nil {
		return err
	}

	for i := 0; i < 30; i++ {
		if !s.IsRunning() {
			s.maybeStopSyncDaemon()
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	if err := s.Screen.Kill(); err != nil {
		return err
	}

	s.maybeStopSyncDaemon()
	return nil
}

func (s *Server) maybeStopSyncDaemon() {
	if !IsSyncDaemonRunning() {
		return
	}

	running, err := AnyServersRunning(s.GlobalCfg)
	if err != nil {
		logging.Warn("failed to check running servers", "error", err)
		return
	}

	if !running {
		if err := StopSyncDaemon(); err != nil {
			logging.Warn("failed to stop sync daemon", "error", err)
		}
	}
}

func (s *Server) Restart(immediate bool) error {
	if err := s.CheckPermission(); err != nil {
		return err
	}

	if s.IsRunning() {
		if !immediate {
			delay := s.Config.RestartDelay
			msg := strings.ReplaceAll(s.Config.MessageRestart, "{DELAY}", strconv.Itoa(delay))

			fmt.Fprintf(os.Stderr, "Warning players (%d second delay)...\n", delay)
			s.Say(msg)
			time.Sleep(time.Duration(delay) * time.Second)
		}

		if err := s.Stop(true); err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	return s.Start()
}

func (s *Server) Status() string {
	if s.IsRunning() {
		return "running"
	}
	return "stopped"
}

func (s *Server) Console() error {
	if err := s.CheckPermission(); err != nil {
		return err
	}

	if !s.IsRunning() {
		return fmt.Errorf("server %q is not running", s.Name)
	}

	return s.Screen.AttachAsUser(s.Config.Username)
}

func (s *Server) SendCommand(cmd string) error {
	if err := s.CheckPermission(); err != nil {
		return err
	}
	return s.Screen.SendCommandAsUser(cmd, s.Config.Username)
}

func (s *Server) Say(message string) error {
	return s.SendCommand(fmt.Sprintf("say %s", message))
}

func (s *Server) Kick(player, reason string) error {
	if reason != "" {
		return s.SendCommand(fmt.Sprintf("kick %s %s", player, reason))
	}
	return s.SendCommand(fmt.Sprintf("kick %s", player))
}

func (s *Server) SaveAll() error {
	return s.SendCommand("save-all")
}

func (s *Server) SaveOff() error {
	return s.SendCommand("save-off")
}

func (s *Server) SaveOn() error {
	return s.SendCommand("save-on")
}

func (s *Server) ConnectedPlayers() ([]string, error) {
	logPath := s.LogPath()

	initialSize, err := getFileSize(logPath)
	if err != nil {
		initialSize = 0
	}

	if err := s.SendCommand("list"); err != nil {
		return nil, err
	}

	time.Sleep(200 * time.Millisecond)

	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	if initialSize > 0 {
		file.Seek(initialSize, 0)
	}

	var players []string
	listPattern := regexp.MustCompile(`players online: (.*)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := listPattern.FindStringSubmatch(line); len(matches) > 1 {
			playerList := strings.TrimSpace(matches[1])
			if playerList != "" {
				for _, p := range strings.Split(playerList, ", ") {
					p = strings.TrimSpace(p)
					if p != "" {
						players = append(players, p)
					}
				}
			}
			break
		}
	}

	return players, scanner.Err()
}

func getFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func Create(name string, cfg *config.Config) (*Server, error) {
	serverPath := filepath.Join(cfg.ServerStoragePath, name)

	if _, err := os.Stat(serverPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("server %q already exists", name)
	}

	ownerUsername := cfg.DefaultUsername
	if !screen.IsRoot() {
		ownerUsername = screen.CurrentUser()
	}

	if err := os.MkdirAll(serverPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create server directory: %w", err)
	}

	worldStoragePath := filepath.Join(serverPath, cfg.DefaultWorldStoragePath)
	if err := os.MkdirAll(worldStoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create world storage directory: %w", err)
	}

	worldStorageInactivePath := filepath.Join(serverPath, cfg.DefaultWorldStorageInactivePath)
	if err := os.MkdirAll(worldStorageInactivePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create inactive world storage directory: %w", err)
	}

	confPath := filepath.Join(serverPath, "server.conf")
	confContent := fmt.Sprintf("USERNAME=%q\n", ownerUsername)
	if err := os.WriteFile(confPath, []byte(confContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to create server.conf: %w", err)
	}

	port := findNextAvailablePort(cfg)
	if err := writeServerProperties(serverPath, port, cfg); err != nil {
		return nil, fmt.Errorf("failed to create server.properties: %w", err)
	}
	logging.Info("assigned port", "server", name, "port", port)

	if err := writeEULA(serverPath); err != nil {
		return nil, fmt.Errorf("failed to create eula.txt: %w", err)
	}

	if err := setOwnership(serverPath, ownerUsername); err != nil {
		logging.Warn("failed to set ownership", "path", serverPath, "user", ownerUsername, "error", err)
	}

	return Load(serverPath, name, cfg)
}

func findNextAvailablePort(cfg *config.Config) int {
	usedPorts := make(map[int]bool)

	entries, err := os.ReadDir(cfg.ServerStoragePath)
	if err != nil {
		return cfg.DefaultServerPort
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		propsPath := filepath.Join(cfg.ServerStoragePath, entry.Name(), cfg.DefaultPropertiesPath)
		port := readPortFromProperties(propsPath)
		if port > 0 {
			usedPorts[port] = true
		}
	}

	port := cfg.DefaultServerPort
	for usedPorts[port] {
		port++
	}
	return port
}

func readPortFromProperties(path string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "server-port=") {
			portStr := strings.TrimPrefix(line, "server-port=")
			port, err := strconv.Atoi(portStr)
			if err == nil {
				return port
			}
		}
	}
	return 0
}

func writeServerProperties(serverPath string, port int, cfg *config.Config) error {
	propsPath := filepath.Join(serverPath, cfg.DefaultPropertiesPath)

	content := fmt.Sprintf(`# Minecraft server properties
# Generated by MSM

server-port=%d
view-distance=%d
max-players=%d
difficulty=%s
gamemode=%s
force-gamemode=true
motd=%s
level-name=world
enable-command-block=false
spawn-protection=16
online-mode=true
white-list=false
`, port, cfg.DefaultRenderDistance, cfg.DefaultMaxPlayers, cfg.DefaultDifficulty, cfg.DefaultGamemode, cfg.DefaultMOTD)

	return os.WriteFile(propsPath, []byte(content), 0644)
}

func mergeServerProperties(serverPath string, port int, cfg *config.Config) ([]string, error) {
	propsPath := filepath.Join(serverPath, cfg.DefaultPropertiesPath)

	existing := make(map[string]string)
	file, err := os.Open(propsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if idx := strings.Index(trimmed, "="); idx != -1 {
			key := trimmed[:idx]
			existing[key] = trimmed[idx+1:]
		}
	}

	defaults := map[string]string{
		"server-port":    strconv.Itoa(port),
		"view-distance":  strconv.Itoa(cfg.DefaultRenderDistance),
		"max-players":    strconv.Itoa(cfg.DefaultMaxPlayers),
		"difficulty":     cfg.DefaultDifficulty,
		"gamemode":       cfg.DefaultGamemode,
		"force-gamemode": "true",
		"level-name":     "world",
	}

	var added []string
	for key, value := range defaults {
		if _, exists := existing[key]; !exists {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
			added = append(added, key)
		}
	}

	if len(added) > 0 {
		content := strings.Join(lines, "\n") + "\n"
		if err := os.WriteFile(propsPath, []byte(content), 0644); err != nil {
			return nil, err
		}
	}

	return added, nil
}

func writeEULA(serverPath string) error {
	eulaPath := filepath.Join(serverPath, "eula.txt")
	content := `# By using MSM to manage Minecraft servers, you agree to the Minecraft EULA
# https://aka.ms/MinecraftEULA
eula=true
`
	return os.WriteFile(eulaPath, []byte(content), 0644)
}

func setOwnership(path, username string) error {
	u, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("user %q not found: %w", username, err)
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return fmt.Errorf("invalid uid: %w", err)
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return fmt.Errorf("invalid gid: %w", err)
	}

	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chown(p, uid, gid)
	})
}

func Delete(name string, cfg *config.Config) error {
	serverPath := filepath.Join(cfg.ServerStoragePath, name)

	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		return fmt.Errorf("server %q not found", name)
	}

	server, err := Load(serverPath, name, cfg)
	if err != nil {
		return err
	}

	if err := server.CheckPermission(); err != nil {
		return err
	}

	if server.IsRunning() {
		return fmt.Errorf("server %q is running; stop it first", name)
	}

	return os.RemoveAll(serverPath)
}

func Init(name string, cfg *config.Config) ([]string, error) {
	serverPath := filepath.Join(cfg.ServerStoragePath, name)

	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("server %q not found", name)
	}

	var created []string

	ownerUsername := cfg.DefaultUsername
	if !screen.IsRoot() {
		ownerUsername = screen.CurrentUser()
	}

	confPath := filepath.Join(serverPath, "server.conf")
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		confContent := fmt.Sprintf("USERNAME=%q\n", ownerUsername)
		if err := os.WriteFile(confPath, []byte(confContent), 0644); err != nil {
			return created, fmt.Errorf("failed to create server.conf: %w", err)
		}
		created = append(created, "server.conf")
	}

	propsPath := filepath.Join(serverPath, cfg.DefaultPropertiesPath)
	if _, err := os.Stat(propsPath); os.IsNotExist(err) {
		port := findNextAvailablePort(cfg)
		if err := writeServerProperties(serverPath, port, cfg); err != nil {
			return created, fmt.Errorf("failed to create server.properties: %w", err)
		}
		created = append(created, fmt.Sprintf("server.properties (port %d)", port))
	} else {
		port := findNextAvailablePort(cfg)
		added, err := mergeServerProperties(serverPath, port, cfg)
		if err != nil {
			logging.Warn("failed to merge server.properties", "error", err)
		} else if len(added) > 0 {
			created = append(created, fmt.Sprintf("server.properties (added: %s)", strings.Join(added, ", ")))
		}
	}

	eulaPath := filepath.Join(serverPath, "eula.txt")
	eulaNeeded := false
	if _, err := os.Stat(eulaPath); os.IsNotExist(err) {
		eulaNeeded = true
	} else {
		content, err := os.ReadFile(eulaPath)
		if err != nil || !strings.Contains(string(content), "eula=true") {
			eulaNeeded = true
		}
	}
	if eulaNeeded {
		if err := writeEULA(serverPath); err != nil {
			return created, fmt.Errorf("failed to create eula.txt: %w", err)
		}
		created = append(created, "eula.txt (accepted)")
	}

	worldStoragePath := filepath.Join(serverPath, cfg.DefaultWorldStoragePath)
	if _, err := os.Stat(worldStoragePath); os.IsNotExist(err) {
		if err := os.MkdirAll(worldStoragePath, 0755); err != nil {
			return created, fmt.Errorf("failed to create world storage directory: %w", err)
		}
		created = append(created, cfg.DefaultWorldStoragePath+"/")
	}

	worldStorageInactivePath := filepath.Join(serverPath, cfg.DefaultWorldStorageInactivePath)
	if _, err := os.Stat(worldStorageInactivePath); os.IsNotExist(err) {
		if err := os.MkdirAll(worldStorageInactivePath, 0755); err != nil {
			return created, fmt.Errorf("failed to create inactive world storage directory: %w", err)
		}
		created = append(created, cfg.DefaultWorldStorageInactivePath+"/")
	}

	if len(created) > 0 {
		if err := setOwnership(serverPath, ownerUsername); err != nil {
			logging.Warn("failed to set ownership", "path", serverPath, "user", ownerUsername, "error", err)
		}
	}

	return created, nil
}

func Rename(oldName, newName string, cfg *config.Config) error {
	oldPath := filepath.Join(cfg.ServerStoragePath, oldName)
	newPath := filepath.Join(cfg.ServerStoragePath, newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("server %q not found", oldName)
	}

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return fmt.Errorf("server %q already exists", newName)
	}

	server, err := Load(oldPath, oldName, cfg)
	if err != nil {
		return err
	}

	if err := server.CheckPermission(); err != nil {
		return err
	}

	if server.IsRunning() {
		return fmt.Errorf("server %q is running; stop it first", oldName)
	}

	return os.Rename(oldPath, newPath)
}
