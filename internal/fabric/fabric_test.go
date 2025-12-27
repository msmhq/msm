package fabric_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/msmhq/msm/internal/fabric"
)

var _ = Describe("Fabric", func() {
	var (
		tempDir string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "msm-fabric-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("Cache", func() {
		It("creates a new empty cache", func() {
			cache := fabric.NewCache()
			Expect(cache).NotTo(BeNil())
		})

		It("loads cache from non-existent file", func() {
			cachePath := filepath.Join(tempDir, "cache.json")
			cache, err := fabric.LoadCache(cachePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(cache).NotTo(BeNil())
		})

		It("saves and loads cache", func() {
			cachePath := filepath.Join(tempDir, "cache.json")
			cache := fabric.NewCache()

			cache.SetGameVersions([]fabric.GameVersion{
				{Version: "1.21.4", Stable: true},
				{Version: "1.21.3", Stable: true},
			})

			cache.SetLoaderVersions("1.21.4", []fabric.LoaderVersion{
				{Version: "0.16.10", Stable: true},
				{Version: "0.16.9", Stable: false},
			})

			cache.SetInstallers([]fabric.InstallerVersion{
				{Version: "1.1.0", Stable: true},
			})

			err := cache.Save(cachePath)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := fabric.LoadCache(cachePath)
			Expect(err).NotTo(HaveOccurred())

			versions, ok := loaded.GetGameVersions(time.Hour)
			Expect(ok).To(BeTrue())
			Expect(versions).To(HaveLen(2))
			Expect(versions[0].Version).To(Equal("1.21.4"))

			loaders, ok := loaded.GetLoaderVersions("1.21.4", time.Hour)
			Expect(ok).To(BeTrue())
			Expect(loaders).To(HaveLen(2))

			installers, ok := loaded.GetInstallers(time.Hour)
			Expect(ok).To(BeTrue())
			Expect(installers).To(HaveLen(1))
		})

		It("respects TTL for game versions", func() {
			cache := fabric.NewCache()
			cache.SetGameVersions([]fabric.GameVersion{
				{Version: "1.21.4", Stable: true},
			})

			versions, ok := cache.GetGameVersions(time.Hour)
			Expect(ok).To(BeTrue())
			Expect(versions).To(HaveLen(1))

			versions, ok = cache.GetGameVersions(0)
			Expect(ok).To(BeFalse())
		})

		It("respects TTL for loader versions", func() {
			cache := fabric.NewCache()
			cache.SetLoaderVersions("1.21.4", []fabric.LoaderVersion{
				{Version: "0.16.10", Stable: true},
			})

			loaders, ok := cache.GetLoaderVersions("1.21.4", time.Hour)
			Expect(ok).To(BeTrue())
			Expect(loaders).To(HaveLen(1))

			loaders, ok = cache.GetLoaderVersions("1.21.4", 0)
			Expect(ok).To(BeFalse())
		})
	})

	Describe("JarFilename", func() {
		It("generates correct filename", func() {
			filename := fabric.JarFilename("1.21.4", "0.16.10", "1.1.0")
			Expect(filename).To(Equal("fabric-server-mc.1.21.4-loader.0.16.10-launcher.1.1.0.jar"))
		})
	})

	Describe("JarPath", func() {
		It("generates correct path", func() {
			path := fabric.JarPath("/opt/msm/fabric", "1.21.4", "0.16.10", "1.1.0")
			Expect(path).To(Equal("/opt/msm/fabric/jars/fabric-server-mc.1.21.4-loader.0.16.10-launcher.1.1.0.jar"))
		})
	})

	Describe("CachePath", func() {
		It("generates correct path", func() {
			path := fabric.CachePath("/opt/msm/fabric")
			Expect(path).To(Equal("/opt/msm/fabric/cache.json"))
		})
	})

	Describe("Client", func() {
		var (
			server *httptest.Server
			client *fabric.Client
		)

		BeforeEach(func() {
			mux := http.NewServeMux()

			mux.HandleFunc("/v2/versions/game", func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode([]fabric.GameVersion{
					{Version: "1.21.4", Stable: true},
					{Version: "1.21.3", Stable: true},
					{Version: "1.21.3-pre1", Stable: false},
				})
			})

			mux.HandleFunc("/v2/versions/loader/1.21.4", func(w http.ResponseWriter, r *http.Request) {
				response := []struct {
					Loader fabric.LoaderVersion `json:"loader"`
				}{
					{Loader: fabric.LoaderVersion{Version: "0.16.10", Stable: true}},
					{Loader: fabric.LoaderVersion{Version: "0.16.9", Stable: false}},
				}
				json.NewEncoder(w).Encode(response)
			})

			mux.HandleFunc("/v2/versions/loader/99.99.99", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})

			mux.HandleFunc("/v2/versions/installer", func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode([]fabric.InstallerVersion{
					{Version: "1.1.0", Stable: true, URL: "https://example.com/installer.jar"},
					{Version: "1.0.1", Stable: false, URL: "https://example.com/installer-old.jar"},
				})
			})

			mux.HandleFunc("/v2/versions/loader/1.21.4/0.16.10/1.1.0/server/jar", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("fake jar content"))
			})

			server = httptest.NewServer(mux)
		})

		AfterEach(func() {
			server.Close()
		})

		Describe("with mock server", func() {
			It("fetches and caches game versions", func() {
				Skip("requires BaseURL override")
			})

			It("fetches and caches loader versions", func() {
				Skip("requires BaseURL override")
			})

			It("fetches and caches installer versions", func() {
				Skip("requires BaseURL override")
			})
		})

		Describe("NewClient", func() {
			It("creates client with new cache if none exists", func() {
				var err error
				client, err = fabric.NewClient(tempDir, 60)
				Expect(err).NotTo(HaveOccurred())
				Expect(client).NotTo(BeNil())
			})

			It("creates client with existing cache", func() {
				cachePath := fabric.CachePath(tempDir)
				cache := fabric.NewCache()
				cache.SetGameVersions([]fabric.GameVersion{
					{Version: "1.21.4", Stable: true},
				})
				err := cache.Save(cachePath)
				Expect(err).NotTo(HaveOccurred())

				client, err = fabric.NewClient(tempDir, 60)
				Expect(err).NotTo(HaveOccurred())
				Expect(client).NotTo(BeNil())
			})
		})
	})
})
