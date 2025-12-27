package player_test

import (
	"encoding/json"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/msmhq/msm/internal/player"
)

var _ = Describe("Player", func() {
	var (
		tempDir string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "msm-player-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("Allowlist", func() {
		var allowlistPath string

		BeforeEach(func() {
			allowlistPath = filepath.Join(tempDir, "whitelist.json")
		})

		Describe("LoadAllowlist", func() {
			It("returns empty list when file does not exist", func() {
				entries, err := player.LoadAllowlist(allowlistPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(BeEmpty())
			})

			It("loads existing entries", func() {
				content := `[{"uuid":"abc-123","name":"player1"},{"name":"player2"}]`
				Expect(os.WriteFile(allowlistPath, []byte(content), 0644)).To(Succeed())

				entries, err := player.LoadAllowlist(allowlistPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(2))
				Expect(entries[0].UUID).To(Equal("abc-123"))
				Expect(entries[0].Name).To(Equal("player1"))
				Expect(entries[1].Name).To(Equal("player2"))
			})

			It("returns error for invalid JSON", func() {
				Expect(os.WriteFile(allowlistPath, []byte("not valid json"), 0644)).To(Succeed())

				_, err := player.LoadAllowlist(allowlistPath)
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("SaveAllowlist", func() {
			It("creates file with entries", func() {
				entries := []player.AllowlistEntry{
					{UUID: "uuid-1", Name: "player1"},
					{Name: "player2"},
				}

				err := player.SaveAllowlist(allowlistPath, entries)
				Expect(err).NotTo(HaveOccurred())

				data, err := os.ReadFile(allowlistPath)
				Expect(err).NotTo(HaveOccurred())

				var loaded []player.AllowlistEntry
				Expect(json.Unmarshal(data, &loaded)).To(Succeed())
				Expect(loaded).To(HaveLen(2))
			})

			It("creates properly indented JSON", func() {
				entries := []player.AllowlistEntry{{Name: "player1"}}

				err := player.SaveAllowlist(allowlistPath, entries)
				Expect(err).NotTo(HaveOccurred())

				data, err := os.ReadFile(allowlistPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("  "))
			})
		})

		Describe("AddToAllowlist", func() {
			It("adds player to empty list", func() {
				err := player.AddToAllowlist(allowlistPath, "newplayer")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadAllowlist(allowlistPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Name).To(Equal("newplayer"))
			})

			It("adds player to existing list", func() {
				content := `[{"name":"player1"}]`
				Expect(os.WriteFile(allowlistPath, []byte(content), 0644)).To(Succeed())

				err := player.AddToAllowlist(allowlistPath, "player2")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadAllowlist(allowlistPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(2))
			})

			It("returns error if player already exists", func() {
				content := `[{"name":"player1"}]`
				Expect(os.WriteFile(allowlistPath, []byte(content), 0644)).To(Succeed())

				err := player.AddToAllowlist(allowlistPath, "player1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already in the allowlist"))
			})
		})

		Describe("RemoveFromAllowlist", func() {
			It("removes player from list", func() {
				content := `[{"name":"player1"},{"name":"player2"}]`
				Expect(os.WriteFile(allowlistPath, []byte(content), 0644)).To(Succeed())

				err := player.RemoveFromAllowlist(allowlistPath, "player1")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadAllowlist(allowlistPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Name).To(Equal("player2"))
			})

			It("returns error if player not in list", func() {
				content := `[{"name":"player1"}]`
				Expect(os.WriteFile(allowlistPath, []byte(content), 0644)).To(Succeed())

				err := player.RemoveFromAllowlist(allowlistPath, "nonexistent")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not in the allowlist"))
			})
		})
	})

	Describe("Ops", func() {
		var opsPath string

		BeforeEach(func() {
			opsPath = filepath.Join(tempDir, "ops.json")
		})

		Describe("LoadOps", func() {
			It("returns empty list when file does not exist", func() {
				entries, err := player.LoadOps(opsPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(BeEmpty())
			})

			It("loads existing entries with all fields", func() {
				content := `[{"uuid":"abc-123","name":"admin","level":4,"bypassesPlayerLimit":true}]`
				Expect(os.WriteFile(opsPath, []byte(content), 0644)).To(Succeed())

				entries, err := player.LoadOps(opsPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].UUID).To(Equal("abc-123"))
				Expect(entries[0].Name).To(Equal("admin"))
				Expect(entries[0].Level).To(Equal(4))
				Expect(entries[0].BypassesPlayerLimit).To(BeTrue())
			})
		})

		Describe("AddOp", func() {
			It("adds operator to empty list", func() {
				err := player.AddOp(opsPath, "admin", 4)
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadOps(opsPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Name).To(Equal("admin"))
				Expect(entries[0].Level).To(Equal(4))
				Expect(entries[0].BypassesPlayerLimit).To(BeFalse())
			})

			It("updates level if operator already exists", func() {
				content := `[{"name":"admin","level":2,"bypassesPlayerLimit":false}]`
				Expect(os.WriteFile(opsPath, []byte(content), 0644)).To(Succeed())

				err := player.AddOp(opsPath, "admin", 4)
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadOps(opsPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Level).To(Equal(4))
			})
		})

		Describe("RemoveOp", func() {
			It("removes operator from list", func() {
				content := `[{"name":"admin","level":4},{"name":"mod","level":2}]`
				Expect(os.WriteFile(opsPath, []byte(content), 0644)).To(Succeed())

				err := player.RemoveOp(opsPath, "admin")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadOps(opsPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Name).To(Equal("mod"))
			})

			It("returns error if player not an operator", func() {
				content := `[{"name":"admin","level":4}]`
				Expect(os.WriteFile(opsPath, []byte(content), 0644)).To(Succeed())

				err := player.RemoveOp(opsPath, "nonexistent")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not an operator"))
			})
		})
	})

	Describe("BannedPlayers", func() {
		var bannedPath string

		BeforeEach(func() {
			bannedPath = filepath.Join(tempDir, "banned-players.json")
		})

		Describe("LoadBannedPlayers", func() {
			It("returns empty list when file does not exist", func() {
				entries, err := player.LoadBannedPlayers(bannedPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(BeEmpty())
			})

			It("loads existing entries with all fields", func() {
				content := `[{"uuid":"abc-123","name":"griefer","created":"2024-01-01","source":"admin","expires":"forever","reason":"griefing"}]`
				Expect(os.WriteFile(bannedPath, []byte(content), 0644)).To(Succeed())

				entries, err := player.LoadBannedPlayers(bannedPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Name).To(Equal("griefer"))
				Expect(entries[0].Reason).To(Equal("griefing"))
				Expect(entries[0].Source).To(Equal("admin"))
			})
		})

		Describe("BanPlayer", func() {
			It("bans player with reason", func() {
				err := player.BanPlayer(bannedPath, "griefer", "destroying builds")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadBannedPlayers(bannedPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Name).To(Equal("griefer"))
				Expect(entries[0].Reason).To(Equal("destroying builds"))
			})

			It("returns error if player already banned", func() {
				content := `[{"name":"griefer","reason":"griefing"}]`
				Expect(os.WriteFile(bannedPath, []byte(content), 0644)).To(Succeed())

				err := player.BanPlayer(bannedPath, "griefer", "more griefing")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already banned"))
			})
		})

		Describe("UnbanPlayer", func() {
			It("unbans player", func() {
				content := `[{"name":"griefer","reason":"griefing"},{"name":"spammer","reason":"spam"}]`
				Expect(os.WriteFile(bannedPath, []byte(content), 0644)).To(Succeed())

				err := player.UnbanPlayer(bannedPath, "griefer")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadBannedPlayers(bannedPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Name).To(Equal("spammer"))
			})

			It("returns error if player not banned", func() {
				content := `[{"name":"griefer","reason":"griefing"}]`
				Expect(os.WriteFile(bannedPath, []byte(content), 0644)).To(Succeed())

				err := player.UnbanPlayer(bannedPath, "innocent")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not banned"))
			})
		})
	})

	Describe("BannedIPs", func() {
		var bannedIPPath string

		BeforeEach(func() {
			bannedIPPath = filepath.Join(tempDir, "banned-ips.json")
		})

		Describe("LoadBannedIPs", func() {
			It("returns empty list when file does not exist", func() {
				entries, err := player.LoadBannedIPs(bannedIPPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(BeEmpty())
			})

			It("loads existing entries", func() {
				content := `[{"ip":"192.168.1.100","reason":"malicious activity"}]`
				Expect(os.WriteFile(bannedIPPath, []byte(content), 0644)).To(Succeed())

				entries, err := player.LoadBannedIPs(bannedIPPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].IP).To(Equal("192.168.1.100"))
				Expect(entries[0].Reason).To(Equal("malicious activity"))
			})
		})

		Describe("BanIP", func() {
			It("bans IP with reason", func() {
				err := player.BanIP(bannedIPPath, "10.0.0.50", "bot network")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadBannedIPs(bannedIPPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].IP).To(Equal("10.0.0.50"))
				Expect(entries[0].Reason).To(Equal("bot network"))
			})

			It("returns error if IP already banned", func() {
				content := `[{"ip":"192.168.1.100","reason":"spam"}]`
				Expect(os.WriteFile(bannedIPPath, []byte(content), 0644)).To(Succeed())

				err := player.BanIP(bannedIPPath, "192.168.1.100", "more spam")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already banned"))
			})
		})

		Describe("UnbanIP", func() {
			It("unbans IP", func() {
				content := `[{"ip":"192.168.1.100"},{"ip":"10.0.0.50"}]`
				Expect(os.WriteFile(bannedIPPath, []byte(content), 0644)).To(Succeed())

				err := player.UnbanIP(bannedIPPath, "192.168.1.100")
				Expect(err).NotTo(HaveOccurred())

				entries, err := player.LoadBannedIPs(bannedIPPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].IP).To(Equal("10.0.0.50"))
			})

			It("returns error if IP not banned", func() {
				content := `[{"ip":"192.168.1.100"}]`
				Expect(os.WriteFile(bannedIPPath, []byte(content), 0644)).To(Succeed())

				err := player.UnbanIP(bannedIPPath, "10.0.0.1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not banned"))
			})
		})
	})
})
