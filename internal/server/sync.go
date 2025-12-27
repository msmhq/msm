package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/msmhq/msm/internal/config"
	"github.com/msmhq/msm/internal/logging"
	"github.com/msmhq/msm/pkg/screen"
)

const (
	SyncScreenName   = "msm-sync"
	SyncIntervalSecs = 120
)

var syncSession = screen.New(SyncScreenName)

func IsSyncDaemonRunning() bool {
	return syncSession.IsRunning()
}

func StartSyncDaemon(cfg *config.Config) error {
	if syncSession.IsRunning() {
		logging.Debug("sync daemon already running")
		return nil
	}

	msmBinary, err := os.Executable()
	if err != nil {
		msmBinary = "msm"
	}

	syncScript := fmt.Sprintf(`
while true; do
    sleep %d
    %s worlds todisk --all 2>/dev/null || true
done
`, SyncIntervalSecs, msmBinary)

	fmt.Fprintf(os.Stderr, "Starting background RAM-to-disk sync (every %d seconds)...\n", SyncIntervalSecs)

	cmd := exec.Command("screen", "-dmS", SyncScreenName, "bash", "-c", syncScript)
	cmd.Dir = cfg.ServerStoragePath
	return cmd.Run()
}

func StopSyncDaemon() error {
	if !syncSession.IsRunning() {
		return nil
	}

	fmt.Fprintf(os.Stderr, "Stopping background RAM-to-disk sync...\n")
	return syncSession.Kill()
}

func AnyServersRunning(cfg *config.Config) (bool, error) {
	servers, err := DiscoverAll(cfg)
	if err != nil {
		return false, err
	}

	for _, s := range servers {
		if s.IsRunning() {
			return true, nil
		}
	}

	return false, nil
}

func AnyRAMWorldsConfigured(cfg *config.Config) (bool, error) {
	if !cfg.RamdiskStorageEnabled {
		return false, nil
	}

	servers, err := DiscoverAll(cfg)
	if err != nil {
		return false, err
	}

	for _, s := range servers {
		worldStoragePath := filepath.Join(s.Path, s.Config.WorldStoragePath)
		entries, err := os.ReadDir(worldStoragePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			flagPath := filepath.Join(worldStoragePath, entry.Name(), "in_ram")
			if _, err := os.Stat(flagPath); err == nil {
				return true, nil
			}
		}
	}

	return false, nil
}
