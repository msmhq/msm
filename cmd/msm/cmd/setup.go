package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up MSM directory structure and permissions",
	Long: `Ensures all MSM directories exist with correct ownership and permissions.

This command must be run as root. It will:
  - Create required directories (/opt/msm/servers, /opt/msm/jars, /opt/msm/fabric, etc.)
  - Set ownership to minecraft:minecraft
  - Set permissions to 2775 (setgid, group-writable)

Run this after installation or to fix permission issues.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if syscall.Getuid() != 0 {
			return fmt.Errorf("setup must be run as root (use sudo)")
		}

		u, err := user.Lookup("minecraft")
		if err != nil {
			return fmt.Errorf("minecraft user not found: %w\n  Hint: Create the user first with 'useradd --system minecraft'", err)
		}

		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			return fmt.Errorf("invalid uid: %w", err)
		}
		gid, err := strconv.Atoi(u.Gid)
		if err != nil {
			return fmt.Errorf("invalid gid: %w", err)
		}

		msmHome := "/opt/msm"
		if cfg != nil && cfg.ServerStoragePath != "" {
			msmHome = filepath.Dir(cfg.ServerStoragePath)
		}

		dirs := []string{
			msmHome,
			filepath.Join(msmHome, "servers"),
			filepath.Join(msmHome, "jars"),
			filepath.Join(msmHome, "jars", "minecraft"),
			filepath.Join(msmHome, "fabric"),
			filepath.Join(msmHome, "fabric", "jars"),
			filepath.Join(msmHome, "versioning"),
			filepath.Join(msmHome, "archives"),
			filepath.Join(msmHome, "archives", "worlds"),
			filepath.Join(msmHome, "archives", "logs"),
			filepath.Join(msmHome, "archives", "backups"),
		}

		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0775); err != nil {
				return fmt.Errorf("failed to create %s: %w", dir, err)
			}

			if err := os.Chown(dir, uid, gid); err != nil {
				return fmt.Errorf("failed to chown %s: %w", dir, err)
			}

			if err := os.Chmod(dir, 0775|os.ModeSetgid); err != nil {
				return fmt.Errorf("failed to chmod %s: %w", dir, err)
			}

			fmt.Printf("  %s (minecraft:minecraft, 2775)\n", dir)
		}

		fmt.Println("Setup complete!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
