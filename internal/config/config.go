package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/msmhq/msm/internal/logging"
)

type Config struct {
	Username string

	ServerStoragePath     string
	JarStoragePath        string
	VersioningStoragePath string

	RamdiskStorageEnabled bool
	RamdiskStoragePath    string

	WorldArchiveEnabled bool
	WorldArchivePath    string
	LogArchivePath      string
	BackupArchivePath   string

	RdiffBackupEnabled  bool
	RdiffBackupRotation int
	RdiffBackupNice     int
	WorldRdiffPath      string

	RsyncBackupEnabled bool
	WorldRsyncPath     string

	DefaultUsername                     string
	DefaultScreenName                   string
	DefaultWorldStoragePath             string
	DefaultWorldStorageInactivePath     string
	DefaultCompleteBackupFollowSymlinks bool

	DefaultLogPath           string
	DefaultPropertiesPath    string
	DefaultAllowlistPath     string
	DefaultBannedPlayersPath string
	DefaultBannedIPsPath     string
	DefaultOpsPath           string
	DefaultOpsList           string

	DefaultJarPath      string
	DefaultRAM          int
	DefaultInvocation   string
	DefaultStopDelay    int
	DefaultRestartDelay int

	DefaultMessageStop                   string
	DefaultMessageStopAbort              string
	DefaultMessageRestart                string
	DefaultMessageRestartAbort           string
	DefaultMessageWorldBackupStarted     string
	DefaultMessageWorldBackupFinished    string
	DefaultMessageCompleteBackupStarted  string
	DefaultMessageCompleteBackupFinished string

	UpdateURL string

	CronMSMBinary            string
	CronMaintenanceHour      int
	CronArchiveRetentionDays int

	FabricStoragePath string
	FabricCacheTTL    int

	DefaultServerPort     int
	DefaultRenderDistance int
	DefaultMaxPlayers     int
	DefaultDifficulty     string
	DefaultGamemode       string
	DefaultMOTD           string
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = os.Getenv("MSM_CONF")
		if path == "" {
			path = "/etc/msm.conf"
		}
	}

	cfg := &Config{
		Username:                             "minecraft",
		ServerStoragePath:                    "/opt/msm/servers",
		JarStoragePath:                       "/opt/msm/jars",
		VersioningStoragePath:                "/opt/msm/versioning",
		RamdiskStorageEnabled:                true,
		RamdiskStoragePath:                   "/dev/shm/msm",
		WorldArchiveEnabled:                  true,
		WorldArchivePath:                     "/opt/msm/archives/worlds",
		LogArchivePath:                       "/opt/msm/archives/logs",
		BackupArchivePath:                    "/opt/msm/archives/backups",
		RdiffBackupEnabled:                   false,
		RdiffBackupRotation:                  7,
		RdiffBackupNice:                      19,
		WorldRdiffPath:                       "/opt/msm/rdiff-backup/worlds",
		RsyncBackupEnabled:                   false,
		WorldRsyncPath:                       "/opt/msm/rsync/worlds",
		DefaultUsername:                      "minecraft",
		DefaultScreenName:                    "msm-{SERVER_NAME}",
		DefaultWorldStoragePath:              "worldstorage",
		DefaultWorldStorageInactivePath:      "worldstorage_inactive",
		DefaultCompleteBackupFollowSymlinks:  false,
		DefaultLogPath:                       "logs/latest.log",
		DefaultPropertiesPath:                "server.properties",
		DefaultAllowlistPath:                 "whitelist.json",
		DefaultBannedPlayersPath:             "banned-players.json",
		DefaultBannedIPsPath:                 "banned-ips.json",
		DefaultOpsPath:                       "ops.json",
		DefaultOpsList:                       "",
		DefaultJarPath:                       "server.jar",
		DefaultRAM:                           1024,
		DefaultInvocation:                    "java -Xms{RAM}M -Xmx{RAM}M -jar {JAR} nogui",
		DefaultStopDelay:                     10,
		DefaultRestartDelay:                  10,
		DefaultMessageStop:                   "SERVER SHUTTING DOWN IN {DELAY} SECONDS!",
		DefaultMessageStopAbort:              "Server shut down aborted.",
		DefaultMessageRestart:                "SERVER REBOOT IN {DELAY} SECONDS!",
		DefaultMessageRestartAbort:           "Server reboot aborted.",
		DefaultMessageWorldBackupStarted:     "Backing up world.",
		DefaultMessageWorldBackupFinished:    "Backup complete.",
		DefaultMessageCompleteBackupStarted:  "Backing up entire server.",
		DefaultMessageCompleteBackupFinished: "Backup complete.",
		UpdateURL:                            "https://raw.githubusercontent.com/msmhq/msm/master",
		CronMSMBinary:                        "/usr/local/bin/msm",
		CronMaintenanceHour:                  5,
		CronArchiveRetentionDays:             30,
		FabricStoragePath:                    "/opt/msm/fabric",
		FabricCacheTTL:                       60,
		DefaultServerPort:                    25565,
		DefaultRenderDistance:                10,
		DefaultMaxPlayers:                    20,
		DefaultDifficulty:                    "normal",
		DefaultGamemode:                      "survival",
		DefaultMOTD:                          "A Minecraft Server",
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
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
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")

		cfg.set(key, value)
	}

	return cfg, scanner.Err()
}

func (c *Config) set(key, value string) {
	switch key {
	case "USERNAME":
		c.Username = value
	case "SERVER_STORAGE_PATH":
		c.ServerStoragePath = value
	case "JAR_STORAGE_PATH":
		c.JarStoragePath = value
	case "VERSIONING_STORAGE_PATH":
		c.VersioningStoragePath = value
	case "RAMDISK_STORAGE_ENABLED":
		c.RamdiskStorageEnabled = value == "true"
	case "RAMDISK_STORAGE_PATH":
		c.RamdiskStoragePath = value
	case "WORLD_ARCHIVE_ENABLED":
		c.WorldArchiveEnabled = value == "true"
	case "WORLD_ARCHIVE_PATH":
		c.WorldArchivePath = value
	case "LOG_ARCHIVE_PATH":
		c.LogArchivePath = value
	case "BACKUP_ARCHIVE_PATH":
		c.BackupArchivePath = value
	case "RDIFF_BACKUP_ENABLED":
		c.RdiffBackupEnabled = value == "true"
	case "RDIFF_BACKUP_ROTATION":
		parseIntConfig(key, value, &c.RdiffBackupRotation)
	case "RDIFF_BACKUP_NICE":
		parseIntConfig(key, value, &c.RdiffBackupNice)
	case "WORLD_RDIFF_PATH":
		c.WorldRdiffPath = value
	case "RSYNC_BACKUP_ENABLED":
		c.RsyncBackupEnabled = value == "true"
	case "WORLD_RSYNC_PATH":
		c.WorldRsyncPath = value
	case "DEFAULT_USERNAME":
		c.DefaultUsername = value
	case "DEFAULT_SCREEN_NAME":
		c.DefaultScreenName = value
	case "DEFAULT_WORLD_STORAGE_PATH":
		c.DefaultWorldStoragePath = value
	case "DEFAULT_WORLD_STORAGE_INACTIVE_PATH":
		c.DefaultWorldStorageInactivePath = value
	case "DEFAULT_COMPLETE_BACKUP_FOLLOW_SYMLINKS":
		c.DefaultCompleteBackupFollowSymlinks = value == "true"
	case "DEFAULT_LOG_PATH":
		c.DefaultLogPath = value
	case "DEFAULT_PROPERTIES_PATH":
		c.DefaultPropertiesPath = value
	case "DEFAULT_WHITELIST_PATH":
		c.DefaultAllowlistPath = value
	case "DEFAULT_BANNED_PLAYERS_PATH":
		c.DefaultBannedPlayersPath = value
	case "DEFAULT_BANNED_IPS_PATH":
		c.DefaultBannedIPsPath = value
	case "DEFAULT_OPS_PATH":
		c.DefaultOpsPath = value
	case "DEFAULT_OPS_LIST":
		c.DefaultOpsList = value
	case "DEFAULT_JAR_PATH":
		c.DefaultJarPath = value
	case "DEFAULT_RAM":
		parseIntConfig(key, value, &c.DefaultRAM)
	case "DEFAULT_INVOCATION":
		c.DefaultInvocation = value
	case "DEFAULT_STOP_DELAY":
		parseIntConfig(key, value, &c.DefaultStopDelay)
	case "DEFAULT_RESTART_DELAY":
		parseIntConfig(key, value, &c.DefaultRestartDelay)
	case "DEFAULT_MESSAGE_STOP":
		c.DefaultMessageStop = value
	case "DEFAULT_MESSAGE_STOP_ABORT":
		c.DefaultMessageStopAbort = value
	case "DEFAULT_MESSAGE_RESTART":
		c.DefaultMessageRestart = value
	case "DEFAULT_MESSAGE_RESTART_ABORT":
		c.DefaultMessageRestartAbort = value
	case "DEFAULT_MESSAGE_WORLD_BACKUP_STARTED":
		c.DefaultMessageWorldBackupStarted = value
	case "DEFAULT_MESSAGE_WORLD_BACKUP_FINISHED":
		c.DefaultMessageWorldBackupFinished = value
	case "DEFAULT_MESSAGE_COMPLETE_BACKUP_STARTED":
		c.DefaultMessageCompleteBackupStarted = value
	case "DEFAULT_MESSAGE_COMPLETE_BACKUP_FINISHED":
		c.DefaultMessageCompleteBackupFinished = value
	case "UPDATE_URL":
		c.UpdateURL = value
	case "CRON_MSM_BINARY":
		c.CronMSMBinary = value
	case "CRON_MAINTENANCE_HOUR":
		parseIntConfig(key, value, &c.CronMaintenanceHour)
	case "CRON_ARCHIVE_RETENTION_DAYS":
		parseIntConfig(key, value, &c.CronArchiveRetentionDays)
	case "FABRIC_STORAGE_PATH":
		c.FabricStoragePath = value
	case "FABRIC_CACHE_TTL_MINUTES":
		parseIntConfig(key, value, &c.FabricCacheTTL)
	case "DEFAULT_SERVER_PORT":
		parseIntConfig(key, value, &c.DefaultServerPort)
	case "DEFAULT_RENDER_DISTANCE":
		parseIntConfig(key, value, &c.DefaultRenderDistance)
	case "DEFAULT_MAX_PLAYERS":
		parseIntConfig(key, value, &c.DefaultMaxPlayers)
	case "DEFAULT_DIFFICULTY":
		c.DefaultDifficulty = value
	case "DEFAULT_GAMEMODE":
		c.DefaultGamemode = value
	case "DEFAULT_MOTD":
		c.DefaultMOTD = value
	}
}

func (c *Config) Print() {
	fmt.Println("MSM Configuration:")
	fmt.Println("==================")
	fmt.Printf("Username: %s\n", c.Username)
	fmt.Printf("Server Storage Path: %s\n", c.ServerStoragePath)
	fmt.Printf("Jar Storage Path: %s\n", c.JarStoragePath)
	fmt.Printf("Versioning Storage Path: %s\n", c.VersioningStoragePath)
	fmt.Printf("Ramdisk Enabled: %v\n", c.RamdiskStorageEnabled)
	fmt.Printf("Ramdisk Path: %s\n", c.RamdiskStoragePath)
	fmt.Printf("World Archive Enabled: %v\n", c.WorldArchiveEnabled)
	fmt.Printf("World Archive Path: %s\n", c.WorldArchivePath)
	fmt.Printf("Log Archive Path: %s\n", c.LogArchivePath)
	fmt.Printf("Backup Archive Path: %s\n", c.BackupArchivePath)
	fmt.Printf("Rdiff Backup Enabled: %v\n", c.RdiffBackupEnabled)
	fmt.Printf("Rsync Backup Enabled: %v\n", c.RsyncBackupEnabled)
	fmt.Printf("Default RAM: %d MB\n", c.DefaultRAM)
	fmt.Printf("Default Stop Delay: %d seconds\n", c.DefaultStopDelay)
	fmt.Printf("Default Restart Delay: %d seconds\n", c.DefaultRestartDelay)
	fmt.Printf("Fabric Storage Path: %s\n", c.FabricStoragePath)
	fmt.Printf("Fabric Cache TTL: %d minutes\n", c.FabricCacheTTL)
}

func parseIntConfig(key, value string, target *int) {
	if v, err := strconv.Atoi(value); err != nil {
		logging.Warn("invalid integer value in config", "key", key, "value", value)
	} else {
		*target = v
	}
}
