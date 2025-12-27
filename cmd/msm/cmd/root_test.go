package cmd_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/msmhq/msm/cmd/msm/cmd"
)

var _ = Describe("Root Commands", func() {
	var (
		tempDir    string
		configFile string
		origArgs   []string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "msm-cmd-test")
		Expect(err).NotTo(HaveOccurred())

		serverDir := filepath.Join(tempDir, "servers")
		Expect(os.MkdirAll(serverDir, 0755)).To(Succeed())

		configFile = filepath.Join(tempDir, "msm.conf")
		configContent := `
SERVER_STORAGE_PATH="` + serverDir + `"
USERNAME="` + os.Getenv("USER") + `"
`
		Expect(os.WriteFile(configFile, []byte(configContent), 0644)).To(Succeed())

		origArgs = os.Args
	})

	AfterEach(func() {
		os.Args = origArgs
		os.RemoveAll(tempDir)
	})

	Describe("start command with server argument", func() {
		Context("when specifying a non-existent server", func() {
			It("returns an error about the specific server not being found", func() {
				os.Args = []string{"msm", "--config", configFile, "start", "nonexistent-server"}

				err := cmd.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nonexistent-server"))
			})
		})

		Context("when specifying one server among multiple", func() {
			BeforeEach(func() {
				serverDir := filepath.Join(tempDir, "servers")

				server1 := filepath.Join(serverDir, "server1")
				Expect(os.MkdirAll(server1, 0755)).To(Succeed())
				Expect(os.WriteFile(filepath.Join(server1, "server.conf"), []byte(`USERNAME="`+os.Getenv("USER")+`"`), 0644)).To(Succeed())

				server2 := filepath.Join(serverDir, "server2")
				Expect(os.MkdirAll(server2, 0755)).To(Succeed())
				Expect(os.WriteFile(filepath.Join(server2, "server.conf"), []byte(`USERNAME="`+os.Getenv("USER")+`"`), 0644)).To(Succeed())
			})

			It("only attempts to start the specified server, not all servers", func() {
				os.Args = []string{"msm", "--config", configFile, "start", "server1"}

				err := cmd.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("server1"))
				Expect(err.Error()).NotTo(ContainSubstring("server2"))
			})
		})
	})

	Describe("stop command with server argument", func() {
		Context("when specifying a non-existent server", func() {
			It("returns an error about the specific server not being found", func() {
				os.Args = []string{"msm", "--config", configFile, "stop", "nonexistent-server"}

				err := cmd.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nonexistent-server"))
			})
		})
	})

	Describe("restart command with server argument", func() {
		Context("when specifying a non-existent server", func() {
			It("returns an error about the specific server not being found", func() {
				os.Args = []string{"msm", "--config", configFile, "restart", "nonexistent-server"}

				err := cmd.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nonexistent-server"))
			})
		})
	})
})
