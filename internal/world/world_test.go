package world_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/msmhq/msm/internal/config"
	"github.com/msmhq/msm/internal/world"
)

var _ = Describe("World", func() {
	var (
		tempDir    string
		serverPath string
		globalCfg  *config.Config
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "msm-world-test")
		Expect(err).NotTo(HaveOccurred())

		serverPath = filepath.Join(tempDir, "survival")
		Expect(os.MkdirAll(serverPath, 0755)).To(Succeed())

		globalCfg = &config.Config{
			RamdiskStoragePath: "/dev/shm/msm",
		}
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("DiscoverAll", func() {
		Context("when no worlds exist", func() {
			It("returns empty slice", func() {
				Expect(os.MkdirAll(filepath.Join(serverPath, "worldstorage"), 0755)).To(Succeed())

				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds).To(BeEmpty())
			})
		})

		Context("when active worlds exist", func() {
			BeforeEach(func() {
				activePath := filepath.Join(serverPath, "worldstorage")
				Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())
				Expect(os.MkdirAll(filepath.Join(activePath, "world_nether"), 0755)).To(Succeed())
				Expect(os.MkdirAll(filepath.Join(activePath, "world_the_end"), 0755)).To(Succeed())
			})

			It("discovers all active worlds", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds).To(HaveLen(3))
			})

			It("marks worlds as active", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())

				for _, w := range worlds {
					Expect(w.Active).To(BeTrue())
				}
			})

			It("sets correct server information", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())

				for _, w := range worlds {
					Expect(w.ServerName).To(Equal("survival"))
					Expect(w.ServerPath).To(Equal(serverPath))
				}
			})
		})

		Context("when inactive worlds exist", func() {
			BeforeEach(func() {
				inactivePath := filepath.Join(serverPath, "worldstorage_inactive")
				Expect(os.MkdirAll(filepath.Join(inactivePath, "old_world"), 0755)).To(Succeed())
			})

			It("discovers inactive worlds", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds).To(HaveLen(1))
			})

			It("marks worlds as inactive", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds[0].Active).To(BeFalse())
			})
		})

		Context("when both active and inactive worlds exist", func() {
			BeforeEach(func() {
				activePath := filepath.Join(serverPath, "worldstorage")
				inactivePath := filepath.Join(serverPath, "worldstorage_inactive")

				Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())
				Expect(os.MkdirAll(filepath.Join(activePath, "world_nether"), 0755)).To(Succeed())
				Expect(os.MkdirAll(filepath.Join(inactivePath, "old_world"), 0755)).To(Succeed())
			})

			It("discovers all worlds", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds).To(HaveLen(3))
			})

			It("correctly categorizes worlds by active status", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())

				activeCount := 0
				inactiveCount := 0
				for _, w := range worlds {
					if w.Active {
						activeCount++
					} else {
						inactiveCount++
					}
				}
				Expect(activeCount).To(Equal(2))
				Expect(inactiveCount).To(Equal(1))
			})
		})

		Context("when world has in_ram flag", func() {
			BeforeEach(func() {
				worldPath := filepath.Join(serverPath, "worldstorage", "world")
				Expect(os.MkdirAll(worldPath, 0755)).To(Succeed())
				Expect(os.WriteFile(filepath.Join(worldPath, "in_ram"), []byte{}, 0644)).To(Succeed())
			})

			It("sets InRAM to true", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds).To(HaveLen(1))
				Expect(worlds[0].InRAM).To(BeTrue())
			})

			It("sets correct RAMPath", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds[0].RAMPath).To(Equal("/dev/shm/msm/survival/world"))
			})
		})

		Context("with absolute paths", func() {
			It("handles absolute world storage paths", func() {
				absPath := filepath.Join(tempDir, "absolute_worlds")
				Expect(os.MkdirAll(filepath.Join(absPath, "world"), 0755)).To(Succeed())

				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, absPath, "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds).To(HaveLen(1))
			})
		})

		Context("ignores non-directory entries", func() {
			BeforeEach(func() {
				activePath := filepath.Join(serverPath, "worldstorage")
				Expect(os.MkdirAll(activePath, 0755)).To(Succeed())
				Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())
				Expect(os.WriteFile(filepath.Join(activePath, "readme.txt"), []byte("ignore"), 0644)).To(Succeed())
			})

			It("only returns directories as worlds", func() {
				worlds, err := world.DiscoverAll(serverPath, "survival", globalCfg, "worldstorage", "worldstorage_inactive")
				Expect(err).NotTo(HaveOccurred())
				Expect(worlds).To(HaveLen(1))
				Expect(worlds[0].Name).To(Equal("world"))
			})
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			activePath := filepath.Join(serverPath, "worldstorage")
			inactivePath := filepath.Join(serverPath, "worldstorage_inactive")

			Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())
			Expect(os.MkdirAll(filepath.Join(inactivePath, "old_world"), 0755)).To(Succeed())
		})

		It("finds active world", func() {
			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(w.Name).To(Equal("world"))
			Expect(w.Active).To(BeTrue())
		})

		It("finds inactive world", func() {
			w, err := world.Get(serverPath, "survival", "old_world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(w.Name).To(Equal("old_world"))
			Expect(w.Active).To(BeFalse())
		})

		It("returns error for non-existent world", func() {
			_, err := world.Get(serverPath, "survival", "nonexistent", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})

		It("detects in_ram flag", func() {
			worldPath := filepath.Join(serverPath, "worldstorage", "world")
			Expect(os.WriteFile(filepath.Join(worldPath, "in_ram"), []byte{}, 0644)).To(Succeed())

			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(w.InRAM).To(BeTrue())
			Expect(w.RAMPath).To(Equal("/dev/shm/msm/survival/world"))
		})
	})

	Describe("Activate", func() {
		BeforeEach(func() {
			inactivePath := filepath.Join(serverPath, "worldstorage_inactive")
			Expect(os.MkdirAll(filepath.Join(inactivePath, "old_world"), 0755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(inactivePath, "old_world", "level.dat"), []byte("data"), 0644)).To(Succeed())
		})

		It("moves world to active directory", func() {
			w, err := world.Get(serverPath, "survival", "old_world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Activate("worldstorage")
			Expect(err).NotTo(HaveOccurred())

			newPath := filepath.Join(serverPath, "worldstorage", "old_world")
			Expect(newPath).To(BeADirectory())
			Expect(filepath.Join(newPath, "level.dat")).To(BeAnExistingFile())
		})

		It("updates world state", func() {
			w, err := world.Get(serverPath, "survival", "old_world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Activate("worldstorage")
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Active).To(BeTrue())
			Expect(w.Path).To(Equal(filepath.Join(serverPath, "worldstorage", "old_world")))
		})

		It("creates active directory if needed", func() {
			w, err := world.Get(serverPath, "survival", "old_world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Activate("worldstorage")
			Expect(err).NotTo(HaveOccurred())

			Expect(filepath.Join(serverPath, "worldstorage")).To(BeADirectory())
		})

		It("returns error if already active", func() {
			activePath := filepath.Join(serverPath, "worldstorage")
			Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())

			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Activate("worldstorage")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already active"))
		})
	})

	Describe("Deactivate", func() {
		BeforeEach(func() {
			activePath := filepath.Join(serverPath, "worldstorage")
			Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(activePath, "world", "level.dat"), []byte("data"), 0644)).To(Succeed())
		})

		It("moves world to inactive directory", func() {
			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Deactivate("worldstorage_inactive")
			Expect(err).NotTo(HaveOccurred())

			newPath := filepath.Join(serverPath, "worldstorage_inactive", "world")
			Expect(newPath).To(BeADirectory())
			Expect(filepath.Join(newPath, "level.dat")).To(BeAnExistingFile())
		})

		It("updates world state", func() {
			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Deactivate("worldstorage_inactive")
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Active).To(BeFalse())
			Expect(w.Path).To(Equal(filepath.Join(serverPath, "worldstorage_inactive", "world")))
		})

		It("creates inactive directory if needed", func() {
			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Deactivate("worldstorage_inactive")
			Expect(err).NotTo(HaveOccurred())

			Expect(filepath.Join(serverPath, "worldstorage_inactive")).To(BeADirectory())
		})

		It("returns error if already inactive", func() {
			inactivePath := filepath.Join(serverPath, "worldstorage_inactive")
			Expect(os.MkdirAll(filepath.Join(inactivePath, "old_world"), 0755)).To(Succeed())

			w, err := world.Get(serverPath, "survival", "old_world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			err = w.Deactivate("worldstorage_inactive")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already inactive"))
		})
	})

	Describe("FlagPath", func() {
		It("returns path to in_ram flag file", func() {
			activePath := filepath.Join(serverPath, "worldstorage")
			Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())

			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			flagPath := w.FlagPath()
			Expect(flagPath).To(Equal(filepath.Join(activePath, "world", "in_ram")))
		})
	})

	Describe("Status", func() {
		BeforeEach(func() {
			activePath := filepath.Join(serverPath, "worldstorage")
			inactivePath := filepath.Join(serverPath, "worldstorage_inactive")
			Expect(os.MkdirAll(filepath.Join(activePath, "world"), 0755)).To(Succeed())
			Expect(os.MkdirAll(filepath.Join(inactivePath, "old_world"), 0755)).To(Succeed())
		})

		It("returns 'active' for active world", func() {
			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Status()).To(Equal("active"))
		})

		It("returns 'inactive' for inactive world", func() {
			w, err := world.Get(serverPath, "survival", "old_world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Status()).To(Equal("inactive"))
		})

		It("returns 'active, in RAM' for active world in RAM", func() {
			worldPath := filepath.Join(serverPath, "worldstorage", "world")
			Expect(os.WriteFile(filepath.Join(worldPath, "in_ram"), []byte{}, 0644)).To(Succeed())

			w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Status()).To(Equal("active, in RAM"))
		})
	})

	Describe("SetupRAMSymlink", func() {
		var (
			worldPath   string
			ramPath     string
			symlinkPath string
		)

		BeforeEach(func() {
			globalCfg.RamdiskStorageEnabled = true
			globalCfg.RamdiskStoragePath = filepath.Join(tempDir, "ramdisk")

			worldPath = filepath.Join(serverPath, "worldstorage", "world")
			Expect(os.MkdirAll(worldPath, 0755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(worldPath, "level.dat"), []byte("data"), 0644)).To(Succeed())

			ramPath = filepath.Join(globalCfg.RamdiskStoragePath, "survival", "world")
			Expect(os.MkdirAll(ramPath, 0755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(ramPath, "level.dat"), []byte("data"), 0644)).To(Succeed())

			symlinkPath = filepath.Join(serverPath, "world")
		})

		Context("when world is in RAM", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(filepath.Join(worldPath, "in_ram"), []byte{}, 0644)).To(Succeed())
			})

			It("creates symlink at server root to RAM path", func() {
				w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(w.InRAM).To(BeTrue())

				err = w.SetupRAMSymlink("worldstorage")
				Expect(err).NotTo(HaveOccurred())

				linkTarget, err := os.Readlink(symlinkPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(linkTarget).To(Equal(ramPath))
			})

			It("is idempotent when symlink already exists", func() {
				w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				err = w.SetupRAMSymlink("worldstorage")
				Expect(err).NotTo(HaveOccurred())

				err = w.SetupRAMSymlink("worldstorage")
				Expect(err).NotTo(HaveOccurred())

				linkTarget, err := os.Readlink(symlinkPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(linkTarget).To(Equal(ramPath))
			})
		})

		Context("when world is not in RAM", func() {
			It("does nothing", func() {
				w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(w.InRAM).To(BeFalse())

				err = w.SetupRAMSymlink("worldstorage")
				Expect(err).NotTo(HaveOccurred())

				_, err = os.Lstat(symlinkPath)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})

		Context("when ramdisk is disabled", func() {
			BeforeEach(func() {
				globalCfg.RamdiskStorageEnabled = false
				Expect(os.WriteFile(filepath.Join(worldPath, "in_ram"), []byte{}, 0644)).To(Succeed())
			})

			It("does nothing", func() {
				w, err := world.Get(serverPath, "survival", "world", "worldstorage", "worldstorage_inactive", globalCfg)
				Expect(err).NotTo(HaveOccurred())

				err = w.SetupRAMSymlink("worldstorage")
				Expect(err).NotTo(HaveOccurred())

				_, err = os.Lstat(symlinkPath)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
	})
})
