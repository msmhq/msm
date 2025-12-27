package server_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/msmhq/msm/internal/config"
	"github.com/msmhq/msm/internal/server"
)

var _ = Describe("Server", func() {
	var (
		tempDir   string
		globalCfg *config.Config
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "msm-server-test")
		Expect(err).NotTo(HaveOccurred())

		globalCfg = &config.Config{
			ServerStoragePath:                 tempDir,
			DefaultUsername:                   "minecraft",
			DefaultScreenName:                 "msm-{SERVER_NAME}",
			DefaultJarPath:                    "server.jar",
			DefaultRAM:                        1024,
			DefaultInvocation:                 "java -Xmx{RAM}M -jar {JAR}",
			DefaultStopDelay:                  10,
			DefaultRestartDelay:               10,
			DefaultWorldStoragePath:           "worldstorage",
			DefaultWorldStorageInactivePath:   "worldstorage_inactive",
			DefaultLogPath:                    "logs/latest.log",
			DefaultPropertiesPath:             "server.properties",
			DefaultAllowlistPath:              "whitelist.json",
			DefaultBannedPlayersPath:          "banned-players.json",
			DefaultBannedIPsPath:              "banned-ips.json",
			DefaultOpsPath:                    "ops.json",
			DefaultMessageStop:                "Stopping server",
			DefaultMessageStopAbort:           "Stop aborted",
			DefaultMessageRestart:             "Restarting server",
			DefaultMessageRestartAbort:        "Restart aborted",
			DefaultMessageWorldBackupStarted:  "Backing up",
			DefaultMessageWorldBackupFinished: "Backup complete",
		}
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("DiscoverAll", func() {
		Context("when server storage does not exist", func() {
			It("returns nil without error", func() {
				globalCfg.ServerStoragePath = "/nonexistent/path"

				servers, err := server.DiscoverAll(globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(servers).To(BeNil())
			})
		})

		Context("when server storage is empty", func() {
			It("returns empty slice", func() {
				servers, err := server.DiscoverAll(globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(servers).To(BeEmpty())
			})
		})

		Context("when servers exist", func() {
			BeforeEach(func() {
				Expect(os.MkdirAll(filepath.Join(tempDir, "survival"), 0755)).To(Succeed())
				Expect(os.MkdirAll(filepath.Join(tempDir, "creative"), 0755)).To(Succeed())
				Expect(os.MkdirAll(filepath.Join(tempDir, "minigames"), 0755)).To(Succeed())
			})

			It("discovers all server directories", func() {
				servers, err := server.DiscoverAll(globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(servers).To(HaveLen(3))
			})

			It("sets correct server names", func() {
				servers, err := server.DiscoverAll(globalCfg)
				Expect(err).NotTo(HaveOccurred())

				names := make([]string, len(servers))
				for i, s := range servers {
					names[i] = s.Name
				}
				Expect(names).To(ContainElements("survival", "creative", "minigames"))
			})

			It("ignores files in server storage", func() {
				Expect(os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("ignore"), 0644)).To(Succeed())

				servers, err := server.DiscoverAll(globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(servers).To(HaveLen(3))
			})
		})
	})

	Describe("Load", func() {
		var serverPath string

		BeforeEach(func() {
			serverPath = filepath.Join(tempDir, "survival")
			Expect(os.MkdirAll(serverPath, 0755)).To(Succeed())
		})

		Context("without server.conf", func() {
			It("uses global defaults", func() {
				s, err := server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(s.Name).To(Equal("survival"))
				Expect(s.Path).To(Equal(serverPath))
				Expect(s.Config.Username).To(Equal("minecraft"))
				Expect(s.Config.ScreenName).To(Equal("msm-survival"))
				Expect(s.Config.RAM).To(Equal(1024))
				Expect(s.Config.JarPath).To(Equal("server.jar"))
			})

			It("replaces {SERVER_NAME} in screen name", func() {
				s, err := server.Load(serverPath, "myserver", globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Config.ScreenName).To(Equal("msm-myserver"))
			})
		})

		Context("with server.conf", func() {
			It("overrides global defaults", func() {
				confContent := `USERNAME="serveruser"
RAM="2048"
JAR_PATH="custom.jar"
STOP_DELAY="30"
`
				Expect(os.WriteFile(filepath.Join(serverPath, "server.conf"), []byte(confContent), 0644)).To(Succeed())

				s, err := server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(s.Config.Username).To(Equal("serveruser"))
				Expect(s.Config.RAM).To(Equal(2048))
				Expect(s.Config.JarPath).To(Equal("custom.jar"))
				Expect(s.Config.StopDelay).To(Equal(30))
			})

			It("replaces {SERVER_NAME} in custom screen name", func() {
				confContent := `SCREEN_NAME="mc-{SERVER_NAME}-screen"
`
				Expect(os.WriteFile(filepath.Join(serverPath, "server.conf"), []byte(confContent), 0644)).To(Succeed())

				s, err := server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Config.ScreenName).To(Equal("mc-survival-screen"))
			})

			It("parses all path configurations", func() {
				confContent := `WORLD_STORAGE_PATH="worlds"
WORLD_STORAGE_INACTIVE_PATH="worlds_off"
LOG_PATH="server.log"
PROPERTIES_PATH="config/server.properties"
WHITELIST_PATH="data/whitelist.json"
BANNED_PLAYERS_PATH="data/bans.json"
BANNED_IPS_PATH="data/ip-bans.json"
OPS_PATH="data/ops.json"
`
				Expect(os.WriteFile(filepath.Join(serverPath, "server.conf"), []byte(confContent), 0644)).To(Succeed())

				s, err := server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(s.Config.WorldStoragePath).To(Equal("worlds"))
				Expect(s.Config.WorldStorageInactivePath).To(Equal("worlds_off"))
				Expect(s.Config.LogPath).To(Equal("server.log"))
				Expect(s.Config.PropertiesPath).To(Equal("config/server.properties"))
				Expect(s.Config.AllowlistPath).To(Equal("data/whitelist.json"))
				Expect(s.Config.BannedPlayersPath).To(Equal("data/bans.json"))
				Expect(s.Config.BannedIPsPath).To(Equal("data/ip-bans.json"))
				Expect(s.Config.OpsPath).To(Equal("data/ops.json"))
			})

			It("parses all message configurations", func() {
				confContent := `MESSAGE_STOP="Shutting down"
MESSAGE_STOP_ABORT="Shutdown cancelled"
MESSAGE_RESTART="Rebooting"
MESSAGE_RESTART_ABORT="Reboot cancelled"
MESSAGE_WORLD_BACKUP_STARTED="World backup started"
MESSAGE_WORLD_BACKUP_FINISHED="World backup done"
MESSAGE_COMPLETE_BACKUP_STARTED="Full backup started"
MESSAGE_COMPLETE_BACKUP_FINISHED="Full backup done"
`
				Expect(os.WriteFile(filepath.Join(serverPath, "server.conf"), []byte(confContent), 0644)).To(Succeed())

				s, err := server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(s.Config.MessageStop).To(Equal("Shutting down"))
				Expect(s.Config.MessageStopAbort).To(Equal("Shutdown cancelled"))
				Expect(s.Config.MessageRestart).To(Equal("Rebooting"))
				Expect(s.Config.MessageRestartAbort).To(Equal("Reboot cancelled"))
				Expect(s.Config.MessageWorldBackupStarted).To(Equal("World backup started"))
				Expect(s.Config.MessageWorldBackupFinished).To(Equal("World backup done"))
				Expect(s.Config.MessageCompleteBackupStarted).To(Equal("Full backup started"))
				Expect(s.Config.MessageCompleteBackupFinished).To(Equal("Full backup done"))
			})

			It("ignores comments and empty lines", func() {
				confContent := `# This is a comment
USERNAME="testuser"

# Another comment
RAM="512"
`
				Expect(os.WriteFile(filepath.Join(serverPath, "server.conf"), []byte(confContent), 0644)).To(Succeed())

				s, err := server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(s.Config.Username).To(Equal("testuser"))
				Expect(s.Config.RAM).To(Equal(512))
			})

			It("handles values with and without quotes", func() {
				confContent := `USERNAME=noquotes
RAM="512"
`
				Expect(os.WriteFile(filepath.Join(serverPath, "server.conf"), []byte(confContent), 0644)).To(Succeed())

				s, err := server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(s.Config.Username).To(Equal("noquotes"))
				Expect(s.Config.RAM).To(Equal(512))
			})
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			Expect(os.MkdirAll(filepath.Join(tempDir, "survival"), 0755)).To(Succeed())
		})

		It("returns server when it exists", func() {
			s, err := server.Get("survival", globalCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.Name).To(Equal("survival"))
		})

		It("returns error when server does not exist", func() {
			_, err := server.Get("nonexistent", globalCfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("Path helpers", func() {
		var s *server.Server

		BeforeEach(func() {
			serverPath := filepath.Join(tempDir, "survival")
			Expect(os.MkdirAll(serverPath, 0755)).To(Succeed())

			var err error
			s, err = server.Load(serverPath, "survival", globalCfg)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("FullPath", func() {
			It("joins relative paths with server path", func() {
				result := s.FullPath("logs/server.log")
				Expect(result).To(Equal(filepath.Join(s.Path, "logs/server.log")))
			})

			It("returns absolute paths unchanged", func() {
				result := s.FullPath("/absolute/path/file.txt")
				Expect(result).To(Equal("/absolute/path/file.txt"))
			})
		})

		Describe("WorldStoragePath", func() {
			It("returns full path to world storage", func() {
				result := s.WorldStoragePath()
				Expect(result).To(Equal(filepath.Join(s.Path, "worldstorage")))
			})
		})

		Describe("WorldStorageInactivePath", func() {
			It("returns full path to inactive world storage", func() {
				result := s.WorldStorageInactivePath()
				Expect(result).To(Equal(filepath.Join(s.Path, "worldstorage_inactive")))
			})
		})

		Describe("JarPath", func() {
			It("returns full path to server jar", func() {
				result := s.JarPath()
				Expect(result).To(Equal(filepath.Join(s.Path, "server.jar")))
			})
		})

		Describe("LogPath", func() {
			It("returns full path to log file", func() {
				result := s.LogPath()
				Expect(result).To(Equal(filepath.Join(s.Path, "logs/latest.log")))
			})
		})

		Describe("PropertiesPath", func() {
			It("returns full path to server properties", func() {
				result := s.PropertiesPath()
				Expect(result).To(Equal(filepath.Join(s.Path, "server.properties")))
			})
		})

		Describe("AllowlistPath", func() {
			It("returns full path to whitelist file", func() {
				result := s.AllowlistPath()
				Expect(result).To(Equal(filepath.Join(s.Path, "whitelist.json")))
			})
		})

		Describe("BannedPlayersPath", func() {
			It("returns full path to banned players file", func() {
				result := s.BannedPlayersPath()
				Expect(result).To(Equal(filepath.Join(s.Path, "banned-players.json")))
			})
		})

		Describe("BannedIPsPath", func() {
			It("returns full path to banned IPs file", func() {
				result := s.BannedIPsPath()
				Expect(result).To(Equal(filepath.Join(s.Path, "banned-ips.json")))
			})
		})

		Describe("OpsPath", func() {
			It("returns full path to ops file", func() {
				result := s.OpsPath()
				Expect(result).To(Equal(filepath.Join(s.Path, "ops.json")))
			})
		})
	})

	Describe("AnyRAMWorldsConfigured", func() {
		var serverPath string

		BeforeEach(func() {
			serverPath = filepath.Join(tempDir, "survival")
			Expect(os.MkdirAll(filepath.Join(serverPath, "worldstorage", "world"), 0755)).To(Succeed())
		})

		Context("when ramdisk is disabled", func() {
			BeforeEach(func() {
				globalCfg.RamdiskStorageEnabled = false
			})

			It("returns false", func() {
				hasRAM, err := server.AnyRAMWorldsConfigured(globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRAM).To(BeFalse())
			})
		})

		Context("when no worlds have in_ram flag", func() {
			BeforeEach(func() {
				globalCfg.RamdiskStorageEnabled = true
			})

			It("returns false", func() {
				hasRAM, err := server.AnyRAMWorldsConfigured(globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRAM).To(BeFalse())
			})
		})

		Context("when a world has in_ram flag", func() {
			BeforeEach(func() {
				globalCfg.RamdiskStorageEnabled = true
				worldPath := filepath.Join(serverPath, "worldstorage", "world")
				Expect(os.WriteFile(filepath.Join(worldPath, "in_ram"), []byte{}, 0644)).To(Succeed())
			})

			It("returns true", func() {
				hasRAM, err := server.AnyRAMWorldsConfigured(globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRAM).To(BeTrue())
			})
		})
	})

	Describe("SetupRAMWorlds", func() {
		var (
			s           *server.Server
			serverPath  string
			ramdiskPath string
		)

		BeforeEach(func() {
			serverPath = filepath.Join(tempDir, "survival")
			ramdiskPath = filepath.Join(tempDir, "ramdisk")

			Expect(os.MkdirAll(serverPath, 0755)).To(Succeed())
			Expect(os.MkdirAll(filepath.Join(serverPath, "worldstorage", "world"), 0755)).To(Succeed())

			globalCfg.RamdiskStorageEnabled = true
			globalCfg.RamdiskStoragePath = ramdiskPath

			var err error
			s, err = server.Load(serverPath, "survival", globalCfg)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when ramdisk is disabled", func() {
			BeforeEach(func() {
				globalCfg.RamdiskStorageEnabled = false
				var err error
				s, err = server.Load(serverPath, "survival", globalCfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("does nothing", func() {
				err := s.SetupRAMWorlds()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when world is flagged for RAM", func() {
			BeforeEach(func() {
				worldPath := filepath.Join(serverPath, "worldstorage", "world")
				Expect(os.WriteFile(filepath.Join(worldPath, "in_ram"), []byte{}, 0644)).To(Succeed())

				ramWorldPath := filepath.Join(ramdiskPath, "survival", "world")
				Expect(os.MkdirAll(ramWorldPath, 0755)).To(Succeed())
			})

			It("sets up symlink at server root to RAM path", func() {
				err := s.SetupRAMWorlds()
				Expect(err).NotTo(HaveOccurred())

				symlinkPath := filepath.Join(serverPath, "world")
				linkTarget, err := os.Readlink(symlinkPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(linkTarget).To(Equal(filepath.Join(ramdiskPath, "survival", "world")))
			})
		})

		Context("when world is not flagged for RAM", func() {
			It("leaves world directory unchanged", func() {
				err := s.SetupRAMWorlds()
				Expect(err).NotTo(HaveOccurred())

				worldPath := filepath.Join(serverPath, "worldstorage", "world")
				info, err := os.Lstat(worldPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.Mode() & os.ModeSymlink).To(BeZero())
			})
		})
	})

	Describe("Create", func() {
		It("generates server.properties with force-gamemode enabled", func() {
			s, err := server.Create("testserver", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			propsPath := filepath.Join(s.Path, "server.properties")
			content, err := os.ReadFile(propsPath)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(content)).To(ContainSubstring("force-gamemode=true"))
		})
	})

	Describe("Init", func() {
		It("adds force-gamemode to existing server.properties if missing", func() {
			serverPath := filepath.Join(tempDir, "existingserver")
			Expect(os.MkdirAll(serverPath, 0755)).To(Succeed())

			serverConf := fmt.Sprintf(`USERNAME="%s"`, os.Getenv("USER"))
			Expect(os.WriteFile(filepath.Join(serverPath, "server.conf"), []byte(serverConf), 0644)).To(Succeed())

			propsContent := `server-port=25565
gamemode=survival
`
			propsPath := filepath.Join(serverPath, "server.properties")
			Expect(os.WriteFile(propsPath, []byte(propsContent), 0644)).To(Succeed())

			created, err := server.Init("existingserver", globalCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(created).To(ContainElement(ContainSubstring("force-gamemode")))

			content, err := os.ReadFile(propsPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("force-gamemode=true"))
		})
	})
})
