package version_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/msmhq/msm/internal/version"
)

var _ = Describe("Version", func() {
	var (
		tempDir string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "msm-version-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("Load", func() {
		Context("with a simple version file", func() {
			It("loads commands and events", func() {
				mcDir := filepath.Join(tempDir, "minecraft")
				Expect(os.MkdirAll(mcDir, 0755)).To(Succeed())

				content := `commands:
  KICK:
    pattern: "kick <player>"
    output:
      - "Kicked <player> from the game"
      - "That player cannot be found"
  SAY:
    pattern: "say <message>"
events:
  START:
    output:
      - "Done"
    timeout: 30
`
				versionPath := filepath.Join(mcDir, "1.0.0.yaml")
				Expect(os.WriteFile(versionPath, []byte(content), 0644)).To(Succeed())

				v, err := version.Load(versionPath, "minecraft", "1.0.0")
				Expect(err).NotTo(HaveOccurred())
				Expect(v.Type).To(Equal("minecraft"))
				Expect(v.Version).To(Equal("1.0.0"))
				Expect(v.Commands).To(HaveLen(2))
				Expect(v.Commands["KICK"].Pattern).To(Equal("kick <player>"))
				Expect(v.Commands["KICK"].Output).To(HaveLen(2))
				Expect(v.Commands["SAY"].Pattern).To(Equal("say <message>"))
				Expect(v.Events).To(HaveLen(1))
				Expect(v.Events["START"].Timeout).To(Equal(30))
			})
		})

		Context("with inheritance", func() {
			BeforeEach(func() {
				mcDir := filepath.Join(tempDir, "minecraft")
				Expect(os.MkdirAll(mcDir, 0755)).To(Succeed())

				baseContent := `commands:
  KICK:
    pattern: "kick <player>"
    output:
      - "Kicked <player>"
  SAY:
    pattern: "say <message>"
events:
  START:
    output:
      - "Done"
    timeout: 30
`
				basePath := filepath.Join(mcDir, "1.0.0.yaml")
				Expect(os.WriteFile(basePath, []byte(baseContent), 0644)).To(Succeed())

				childContent := `extends: minecraft/1.0.0

commands:
  KICK:
    pattern: "kick <player> [reason]"
    output:
      - "Kicked <player> from the game"
      - "That player cannot be found"
  BAN:
    pattern: "ban <player>"
`
				childPath := filepath.Join(mcDir, "1.1.0.yaml")
				Expect(os.WriteFile(childPath, []byte(childContent), 0644)).To(Succeed())
			})

			It("inherits commands from parent", func() {
				childPath := filepath.Join(tempDir, "minecraft", "1.1.0.yaml")
				v, err := version.Load(childPath, "minecraft", "1.1.0")
				Expect(err).NotTo(HaveOccurred())

				Expect(v.Commands["SAY"].Pattern).To(Equal("say <message>"))
			})

			It("overrides parent commands with child commands", func() {
				childPath := filepath.Join(tempDir, "minecraft", "1.1.0.yaml")
				v, err := version.Load(childPath, "minecraft", "1.1.0")
				Expect(err).NotTo(HaveOccurred())

				Expect(v.Commands["KICK"].Pattern).To(Equal("kick <player> [reason]"))
				Expect(v.Commands["KICK"].Output).To(HaveLen(2))
			})

			It("adds new commands from child", func() {
				childPath := filepath.Join(tempDir, "minecraft", "1.1.0.yaml")
				v, err := version.Load(childPath, "minecraft", "1.1.0")
				Expect(err).NotTo(HaveOccurred())

				Expect(v.Commands).To(HaveKey("BAN"))
				Expect(v.Commands["BAN"].Pattern).To(Equal("ban <player>"))
			})

			It("inherits events from parent", func() {
				childPath := filepath.Join(tempDir, "minecraft", "1.1.0.yaml")
				v, err := version.Load(childPath, "minecraft", "1.1.0")
				Expect(err).NotTo(HaveOccurred())

				Expect(v.Events).To(HaveKey("START"))
				Expect(v.Events["START"].Timeout).To(Equal(30))
			})
		})

		Context("with cross-type inheritance", func() {
			BeforeEach(func() {
				mcDir := filepath.Join(tempDir, "minecraft")
				cbDir := filepath.Join(tempDir, "craftbukkit")
				Expect(os.MkdirAll(mcDir, 0755)).To(Succeed())
				Expect(os.MkdirAll(cbDir, 0755)).To(Succeed())

				baseContent := `commands:
  KICK:
    pattern: "kick <player>"
events:
  START:
    output:
      - "Done"
`
				basePath := filepath.Join(mcDir, "1.0.0.yaml")
				Expect(os.WriteFile(basePath, []byte(baseContent), 0644)).To(Succeed())

				childContent := `extends: minecraft/1.0.0

commands:
  RELOAD:
    pattern: "reload"
`
				childPath := filepath.Join(cbDir, "1.0.0.yaml")
				Expect(os.WriteFile(childPath, []byte(childContent), 0644)).To(Succeed())
			})

			It("inherits from different type", func() {
				childPath := filepath.Join(tempDir, "craftbukkit", "1.0.0.yaml")
				v, err := version.Load(childPath, "craftbukkit", "1.0.0")
				Expect(err).NotTo(HaveOccurred())

				Expect(v.Type).To(Equal("craftbukkit"))
				Expect(v.Commands).To(HaveKey("KICK"))
				Expect(v.Commands).To(HaveKey("RELOAD"))
			})
		})

		Context("with empty version file", func() {
			It("loads without error", func() {
				mcDir := filepath.Join(tempDir, "minecraft")
				Expect(os.MkdirAll(mcDir, 0755)).To(Succeed())

				versionPath := filepath.Join(mcDir, "1.0.0.yaml")
				Expect(os.WriteFile(versionPath, []byte(""), 0644)).To(Succeed())

				v, err := version.Load(versionPath, "minecraft", "1.0.0")
				Expect(err).NotTo(HaveOccurred())
				Expect(v.Commands).To(BeNil())
				Expect(v.Events).To(BeNil())
			})
		})
	})

	Describe("LoadAll", func() {
		BeforeEach(func() {
			mcDir := filepath.Join(tempDir, "minecraft")
			cbDir := filepath.Join(tempDir, "craftbukkit")
			Expect(os.MkdirAll(mcDir, 0755)).To(Succeed())
			Expect(os.MkdirAll(cbDir, 0755)).To(Succeed())

			for _, ver := range []string{"1.0.0", "1.1.0", "1.2.0"} {
				content := "commands: {}\n"
				Expect(os.WriteFile(filepath.Join(mcDir, ver+".yaml"), []byte(content), 0644)).To(Succeed())
			}

			cbContent := `extends: minecraft/1.0.0
`
			Expect(os.WriteFile(filepath.Join(cbDir, "1.0.0.yaml"), []byte(cbContent), 0644)).To(Succeed())
		})

		It("loads all version files", func() {
			versions, err := version.LoadAll(tempDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(versions).To(HaveLen(4))
		})

		It("loads versions of different types", func() {
			versions, err := version.LoadAll(tempDir)
			Expect(err).NotTo(HaveOccurred())

			types := make(map[string]int)
			for _, v := range versions {
				types[v.Type]++
			}
			Expect(types["minecraft"]).To(Equal(3))
			Expect(types["craftbukkit"]).To(Equal(1))
		})

		It("ignores non-YAML files", func() {
			Expect(os.WriteFile(filepath.Join(tempDir, "minecraft", "readme.txt"), []byte("ignore me"), 0644)).To(Succeed())

			versions, err := version.LoadAll(tempDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(versions).To(HaveLen(4))
		})

		It("ignores files in wrong directory structure", func() {
			Expect(os.WriteFile(filepath.Join(tempDir, "invalid.yaml"), []byte("commands: {}"), 0644)).To(Succeed())

			versions, err := version.LoadAll(tempDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(versions).To(HaveLen(4))
		})
	})

	Describe("FindClosest", func() {
		var versions []*version.Version

		BeforeEach(func() {
			versions = []*version.Version{
				{Type: "minecraft", Version: "1.0.0"},
				{Type: "minecraft", Version: "1.2.0"},
				{Type: "minecraft", Version: "1.5.0"},
				{Type: "minecraft", Version: "1.10.0"},
				{Type: "craftbukkit", Version: "1.0.0"},
				{Type: "craftbukkit", Version: "1.5.0"},
			}
		})

		It("finds exact match", func() {
			v := version.FindClosest(versions, "minecraft", "1.2.0")
			Expect(v).NotTo(BeNil())
			Expect(v.Version).To(Equal("1.2.0"))
		})

		It("finds closest lower version when exact match not available", func() {
			v := version.FindClosest(versions, "minecraft", "1.3.0")
			Expect(v).NotTo(BeNil())
			Expect(v.Version).To(Equal("1.2.0"))
		})

		It("finds closest for version between available versions", func() {
			v := version.FindClosest(versions, "minecraft", "1.7.0")
			Expect(v).NotTo(BeNil())
			Expect(v.Version).To(Equal("1.5.0"))
		})

		It("returns first version when target is lower than all", func() {
			v := version.FindClosest(versions, "minecraft", "0.5.0")
			Expect(v).NotTo(BeNil())
			Expect(v.Version).To(Equal("1.0.0"))
		})

		It("returns highest version when target is higher than all", func() {
			v := version.FindClosest(versions, "minecraft", "2.0.0")
			Expect(v).NotTo(BeNil())
			Expect(v.Version).To(Equal("1.10.0"))
		})

		It("returns nil when type not found", func() {
			v := version.FindClosest(versions, "spigot", "1.0.0")
			Expect(v).To(BeNil())
		})

		It("filters by type correctly", func() {
			v := version.FindClosest(versions, "craftbukkit", "1.3.0")
			Expect(v).NotTo(BeNil())
			Expect(v.Type).To(Equal("craftbukkit"))
			Expect(v.Version).To(Equal("1.0.0"))
		})

		It("handles version comparison with different segment counts", func() {
			versions = append(versions, &version.Version{Type: "minecraft", Version: "1.5"})
			v := version.FindClosest(versions, "minecraft", "1.5.1")
			Expect(v).NotTo(BeNil())
			Expect(v.Version).To(Equal("1.5"))
		})

		It("handles version 1.10 vs 1.9 correctly", func() {
			versions = append(versions, &version.Version{Type: "minecraft", Version: "1.9.0"})
			v := version.FindClosest(versions, "minecraft", "1.11.0")
			Expect(v).NotTo(BeNil())
			Expect(v.Version).To(Equal("1.10.0"))
		})

		Context("with year-based versions (25.x, 26.x)", func() {
			It("correctly orders year-based versions after 1.x versions", func() {
				v := version.FindClosest(versions, "minecraft", "25.11")
				Expect(v).NotTo(BeNil())
				Expect(v.Version).To(Equal("1.10.0"))
			})

			It("finds closest when mix of old and year-based versions exist", func() {
				versions = append(versions, &version.Version{Type: "minecraft", Version: "25.1"})
				versions = append(versions, &version.Version{Type: "minecraft", Version: "25.11"})

				v := version.FindClosest(versions, "minecraft", "25.5")
				Expect(v).NotTo(BeNil())
				Expect(v.Version).To(Equal("25.1"))
			})

			It("finds year-based version when target is newer year", func() {
				versions = append(versions, &version.Version{Type: "minecraft", Version: "25.1"})
				versions = append(versions, &version.Version{Type: "minecraft", Version: "25.11"})

				v := version.FindClosest(versions, "minecraft", "26.1")
				Expect(v).NotTo(BeNil())
				Expect(v.Version).To(Equal("25.11"))
			})

			It("correctly compares 1.21.11 (old) vs 25.11 (new)", func() {
				versions = append(versions, &version.Version{Type: "minecraft", Version: "1.21.11"})

				v := version.FindClosest(versions, "minecraft", "25.11")
				Expect(v).NotTo(BeNil())
				Expect(v.Version).To(Equal("1.21.11"))
			})

			It("correctly compares multiple year-based versions", func() {
				versions = []*version.Version{
					{Type: "minecraft", Version: "1.7.0"},
					{Type: "minecraft", Version: "25.1"},
					{Type: "minecraft", Version: "25.11"},
					{Type: "minecraft", Version: "26.1"},
					{Type: "minecraft", Version: "26.2"},
				}

				v := version.FindClosest(versions, "minecraft", "26.5")
				Expect(v).NotTo(BeNil())
				Expect(v.Version).To(Equal("26.2"))
			})
		})
	})

	Describe("GetCommand", func() {
		It("returns command when it exists", func() {
			v := &version.Version{
				Commands: map[string]version.ConsoleCommand{
					"KICK": {Pattern: "kick <player>"},
				},
			}

			cmd := v.GetCommand("KICK")
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Pattern).To(Equal("kick <player>"))
		})

		It("returns nil when command does not exist", func() {
			v := &version.Version{
				Commands: map[string]version.ConsoleCommand{},
			}

			cmd := v.GetCommand("NONEXISTENT")
			Expect(cmd).To(BeNil())
		})
	})

	Describe("GetEvent", func() {
		It("returns event when it exists", func() {
			v := &version.Version{
				Events: map[string]version.ConsoleEvent{
					"START": {Timeout: 30},
				},
			}

			evt := v.GetEvent("START")
			Expect(evt).NotTo(BeNil())
			Expect(evt.Timeout).To(Equal(30))
		})

		It("returns nil when event does not exist", func() {
			v := &version.Version{
				Events: map[string]version.ConsoleEvent{},
			}

			evt := v.GetEvent("NONEXISTENT")
			Expect(evt).To(BeNil())
		})
	})

	Describe("String", func() {
		It("formats as type/version", func() {
			v := &version.Version{
				Type:    "minecraft",
				Version: "1.20.4",
			}

			Expect(v.String()).To(Equal("minecraft/1.20.4"))
		})
	})
})
