package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/msmhq/msm/internal/config"
)

var _ = Describe("Config", func() {
	var (
		tempDir    string
		configPath string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "msm-config-test")
		Expect(err).NotTo(HaveOccurred())
		configPath = filepath.Join(tempDir, "msm.conf")
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("Load", func() {
		Context("when config file does not exist", func() {
			It("returns default configuration", func() {
				cfg, err := config.Load("/nonexistent/path/msm.conf")
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())
				Expect(cfg.Username).To(Equal("minecraft"))
				Expect(cfg.ServerStoragePath).To(Equal("/opt/msm/servers"))
				Expect(cfg.DefaultRAM).To(Equal(1024))
			})
		})

		Context("when config file exists", func() {
			It("parses simple key=value pairs", func() {
				content := `USERNAME="testuser"
SERVER_STORAGE_PATH="/custom/servers"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("testuser"))
				Expect(cfg.ServerStoragePath).To(Equal("/custom/servers"))
			})

			It("parses values without quotes", func() {
				content := `USERNAME=noquotes
SERVER_STORAGE_PATH=/no/quotes/path
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("noquotes"))
				Expect(cfg.ServerStoragePath).To(Equal("/no/quotes/path"))
			})

			It("ignores empty lines", func() {
				content := `USERNAME="testuser"

SERVER_STORAGE_PATH="/custom/servers"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("testuser"))
				Expect(cfg.ServerStoragePath).To(Equal("/custom/servers"))
			})

			It("ignores comment lines", func() {
				content := `# This is a comment
USERNAME="testuser"
# Another comment
SERVER_STORAGE_PATH="/custom/servers"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("testuser"))
				Expect(cfg.ServerStoragePath).To(Equal("/custom/servers"))
			})

			It("ignores malformed lines", func() {
				content := `USERNAME="testuser"
this line has no equals sign
SERVER_STORAGE_PATH="/custom/servers"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("testuser"))
				Expect(cfg.ServerStoragePath).To(Equal("/custom/servers"))
			})

			It("handles whitespace around keys and values", func() {
				content := `  USERNAME  =  "testuser"
	SERVER_STORAGE_PATH	=	"/custom/servers"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("testuser"))
				Expect(cfg.ServerStoragePath).To(Equal("/custom/servers"))
			})
		})

		Context("when parsing boolean values", func() {
			It("parses 'true' as true", func() {
				content := `RAMDISK_STORAGE_ENABLED="true"
RDIFF_BACKUP_ENABLED="true"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.RamdiskStorageEnabled).To(BeTrue())
				Expect(cfg.RdiffBackupEnabled).To(BeTrue())
			})

			It("parses anything other than 'true' as false", func() {
				content := `RAMDISK_STORAGE_ENABLED="false"
RDIFF_BACKUP_ENABLED="yes"
RSYNC_BACKUP_ENABLED="1"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.RamdiskStorageEnabled).To(BeFalse())
				Expect(cfg.RdiffBackupEnabled).To(BeFalse())
				Expect(cfg.RsyncBackupEnabled).To(BeFalse())
			})
		})

		Context("when parsing integer values", func() {
			It("parses valid integers", func() {
				content := `DEFAULT_RAM="2048"
DEFAULT_STOP_DELAY="30"
DEFAULT_RESTART_DELAY="15"
RDIFF_BACKUP_ROTATION="14"
RDIFF_BACKUP_NICE="10"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.DefaultRAM).To(Equal(2048))
				Expect(cfg.DefaultStopDelay).To(Equal(30))
				Expect(cfg.DefaultRestartDelay).To(Equal(15))
				Expect(cfg.RdiffBackupRotation).To(Equal(14))
				Expect(cfg.RdiffBackupNice).To(Equal(10))
			})

			It("keeps default value for invalid integers", func() {
				content := `DEFAULT_RAM="notanumber"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.DefaultRAM).To(Equal(1024))
			})
		})

		Context("when parsing all storage paths", func() {
			It("parses all path configurations", func() {
				content := `SERVER_STORAGE_PATH="/servers"
JAR_STORAGE_PATH="/jars"
VERSIONING_STORAGE_PATH="/versioning"
RAMDISK_STORAGE_PATH="/ramdisk"
WORLD_ARCHIVE_PATH="/archives/worlds"
LOG_ARCHIVE_PATH="/archives/logs"
BACKUP_ARCHIVE_PATH="/archives/backups"
WORLD_RDIFF_PATH="/rdiff/worlds"
WORLD_RSYNC_PATH="/rsync/worlds"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.ServerStoragePath).To(Equal("/servers"))
				Expect(cfg.JarStoragePath).To(Equal("/jars"))
				Expect(cfg.VersioningStoragePath).To(Equal("/versioning"))
				Expect(cfg.RamdiskStoragePath).To(Equal("/ramdisk"))
				Expect(cfg.WorldArchivePath).To(Equal("/archives/worlds"))
				Expect(cfg.LogArchivePath).To(Equal("/archives/logs"))
				Expect(cfg.BackupArchivePath).To(Equal("/archives/backups"))
				Expect(cfg.WorldRdiffPath).To(Equal("/rdiff/worlds"))
				Expect(cfg.WorldRsyncPath).To(Equal("/rsync/worlds"))
			})
		})

		Context("when parsing default server settings", func() {
			It("parses default file paths", func() {
				content := `DEFAULT_LOG_PATH="logs/server.log"
DEFAULT_PROPERTIES_PATH="config/server.properties"
DEFAULT_WHITELIST_PATH="data/whitelist.json"
DEFAULT_BANNED_PLAYERS_PATH="data/banned-players.json"
DEFAULT_BANNED_IPS_PATH="data/banned-ips.json"
DEFAULT_OPS_PATH="data/ops.json"
DEFAULT_JAR_PATH="minecraft.jar"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.DefaultLogPath).To(Equal("logs/server.log"))
				Expect(cfg.DefaultPropertiesPath).To(Equal("config/server.properties"))
				Expect(cfg.DefaultAllowlistPath).To(Equal("data/whitelist.json"))
				Expect(cfg.DefaultBannedPlayersPath).To(Equal("data/banned-players.json"))
				Expect(cfg.DefaultBannedIPsPath).To(Equal("data/banned-ips.json"))
				Expect(cfg.DefaultOpsPath).To(Equal("data/ops.json"))
				Expect(cfg.DefaultJarPath).To(Equal("minecraft.jar"))
			})

			It("parses world storage paths", func() {
				content := `DEFAULT_WORLD_STORAGE_PATH="worlds"
DEFAULT_WORLD_STORAGE_INACTIVE_PATH="worlds_inactive"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.DefaultWorldStoragePath).To(Equal("worlds"))
				Expect(cfg.DefaultWorldStorageInactivePath).To(Equal("worlds_inactive"))
			})
		})

		Context("when parsing default messages", func() {
			It("parses all message templates", func() {
				content := `DEFAULT_MESSAGE_STOP="Stopping in {DELAY}s"
DEFAULT_MESSAGE_STOP_ABORT="Stop cancelled"
DEFAULT_MESSAGE_RESTART="Restarting in {DELAY}s"
DEFAULT_MESSAGE_RESTART_ABORT="Restart cancelled"
DEFAULT_MESSAGE_WORLD_BACKUP_STARTED="World backup starting"
DEFAULT_MESSAGE_WORLD_BACKUP_FINISHED="World backup done"
DEFAULT_MESSAGE_COMPLETE_BACKUP_STARTED="Full backup starting"
DEFAULT_MESSAGE_COMPLETE_BACKUP_FINISHED="Full backup done"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.DefaultMessageStop).To(Equal("Stopping in {DELAY}s"))
				Expect(cfg.DefaultMessageStopAbort).To(Equal("Stop cancelled"))
				Expect(cfg.DefaultMessageRestart).To(Equal("Restarting in {DELAY}s"))
				Expect(cfg.DefaultMessageRestartAbort).To(Equal("Restart cancelled"))
				Expect(cfg.DefaultMessageWorldBackupStarted).To(Equal("World backup starting"))
				Expect(cfg.DefaultMessageWorldBackupFinished).To(Equal("World backup done"))
				Expect(cfg.DefaultMessageCompleteBackupStarted).To(Equal("Full backup starting"))
				Expect(cfg.DefaultMessageCompleteBackupFinished).To(Equal("Full backup done"))
			})
		})

		Context("when parsing other settings", func() {
			It("parses screen name and invocation", func() {
				content := `DEFAULT_SCREEN_NAME="mc-{SERVER_NAME}"
DEFAULT_INVOCATION="java -Xms{RAM}M -Xmx{RAM}M -jar {JAR}"
DEFAULT_OPS_LIST="admin,moderator"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.DefaultScreenName).To(Equal("mc-{SERVER_NAME}"))
				Expect(cfg.DefaultInvocation).To(Equal("java -Xms{RAM}M -Xmx{RAM}M -jar {JAR}"))
				Expect(cfg.DefaultOpsList).To(Equal("admin,moderator"))
			})

			It("parses update URL", func() {
				content := `UPDATE_URL="https://custom.url/msm"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.UpdateURL).To(Equal("https://custom.url/msm"))
			})

			It("parses complete backup follow symlinks", func() {
				content := `DEFAULT_COMPLETE_BACKUP_FOLLOW_SYMLINKS="true"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())

				cfg, err := config.Load(configPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.DefaultCompleteBackupFollowSymlinks).To(BeTrue())
			})
		})

		Context("when using MSM_CONF environment variable", func() {
			var originalEnv string

			BeforeEach(func() {
				originalEnv = os.Getenv("MSM_CONF")
			})

			AfterEach(func() {
				if originalEnv == "" {
					os.Unsetenv("MSM_CONF")
				} else {
					os.Setenv("MSM_CONF", originalEnv)
				}
			})

			It("uses MSM_CONF when path is empty", func() {
				content := `USERNAME="envuser"
`
				Expect(os.WriteFile(configPath, []byte(content), 0644)).To(Succeed())
				os.Setenv("MSM_CONF", configPath)

				cfg, err := config.Load("")
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("envuser"))
			})
		})
	})

	Describe("Default Values", func() {
		It("has correct default values when no config file exists", func() {
			cfg, err := config.Load("/nonexistent/config")
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.Username).To(Equal("minecraft"))
			Expect(cfg.ServerStoragePath).To(Equal("/opt/msm/servers"))
			Expect(cfg.JarStoragePath).To(Equal("/opt/msm/jars"))
			Expect(cfg.VersioningStoragePath).To(Equal("/opt/msm/versioning"))
			Expect(cfg.RamdiskStorageEnabled).To(BeTrue())
			Expect(cfg.RamdiskStoragePath).To(Equal("/dev/shm/msm"))
			Expect(cfg.WorldArchiveEnabled).To(BeTrue())
			Expect(cfg.WorldArchivePath).To(Equal("/opt/msm/archives/worlds"))
			Expect(cfg.LogArchivePath).To(Equal("/opt/msm/archives/logs"))
			Expect(cfg.BackupArchivePath).To(Equal("/opt/msm/archives/backups"))
			Expect(cfg.RdiffBackupEnabled).To(BeFalse())
			Expect(cfg.RdiffBackupRotation).To(Equal(7))
			Expect(cfg.RdiffBackupNice).To(Equal(19))
			Expect(cfg.RsyncBackupEnabled).To(BeFalse())
			Expect(cfg.DefaultRAM).To(Equal(1024))
			Expect(cfg.DefaultStopDelay).To(Equal(10))
			Expect(cfg.DefaultRestartDelay).To(Equal(10))
			Expect(cfg.DefaultAllowlistPath).To(Equal("whitelist.json"))
			Expect(cfg.DefaultLogPath).To(Equal("logs/latest.log"))
			Expect(cfg.UpdateURL).To(Equal("https://raw.githubusercontent.com/msmhq/msm/master"))
		})
	})
})
