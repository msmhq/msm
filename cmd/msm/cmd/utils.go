package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/msmhq/msm/internal/log"
	"github.com/msmhq/msm/internal/server"
	"github.com/spf13/cobra"
)

var logrollCmd = &cobra.Command{
	Use:   "logroll [server]",
	Short: "Roll (archive) server logs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		if all || len(args) == 0 {
			servers, err := server.DiscoverAll(cfg)
			if err != nil {
				return err
			}
			for _, s := range servers {
				if err := log.Roll(s.LogPath(), cfg.LogArchivePath, s.Name); err != nil {
					fmt.Printf("Failed to roll logs for %s: %v\n", s.Name, err)
				} else {
					fmt.Printf("Rolled logs for %s\n", s.Name)
				}
			}
			return nil
		}

		serverName := args[0]
		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		return log.Roll(s.LogPath(), cfg.LogArchivePath, serverName)
	},
}

var serverConfigCmd = &cobra.Command{
	Use:   "config <server> [key]",
	Short: "Show per-server configuration",
	Long: `Show per-server configuration.

Without a key, shows all configuration values.
With a key, shows just that value.
Use 'config set' to change values.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if len(args) == 1 {
			fmt.Printf("Server %q configuration:\n", serverName)
			fmt.Printf("  Path: %s\n", s.Path)
			fmt.Printf("  Username: %s\n", s.Config.Username)
			fmt.Printf("  Screen Name: %s\n", s.Config.ScreenName)
			fmt.Printf("  Jar Path: %s\n", s.Config.JarPath)
			fmt.Printf("  RAM: %d MB\n", s.Config.RAM)
			fmt.Printf("  Stop Delay: %d seconds\n", s.Config.StopDelay)
			fmt.Printf("  Restart Delay: %d seconds\n", s.Config.RestartDelay)
			fmt.Printf("  World Storage: %s\n", s.Config.WorldStoragePath)
			fmt.Printf("  Inactive World Storage: %s\n", s.Config.WorldStorageInactivePath)
			fmt.Printf("  Fabric Enabled: %v\n", s.Config.FabricEnabled)
			if s.Config.FabricEnabled {
				if s.Config.FabricLoaderVersion != "" {
					fmt.Printf("  Fabric Loader Version: %s\n", s.Config.FabricLoaderVersion)
				} else {
					fmt.Printf("  Fabric Loader Version: (latest)\n")
				}
				if s.Config.FabricInstallerVersion != "" {
					fmt.Printf("  Fabric Installer Version: %s\n", s.Config.FabricInstallerVersion)
				}
			}
			return nil
		}

		key := args[1]
		switch key {
		case "username":
			fmt.Println(s.Config.Username)
		case "screen_name":
			fmt.Println(s.Config.ScreenName)
		case "jar_path":
			fmt.Println(s.Config.JarPath)
		case "ram":
			fmt.Println(s.Config.RAM)
		case "stop_delay":
			fmt.Println(s.Config.StopDelay)
		case "restart_delay":
			fmt.Println(s.Config.RestartDelay)
		case "world_storage_path":
			fmt.Println(s.Config.WorldStoragePath)
		case "world_storage_inactive_path":
			fmt.Println(s.Config.WorldStorageInactivePath)
		default:
			return fmt.Errorf("unknown config key: %s\n  Valid keys: username, screen_name, jar_path, ram, stop_delay, restart_delay, world_storage_path, world_storage_inactive_path", key)
		}
		return nil
	},
}

var serverConfigSetCmd = &cobra.Command{
	Use:   "set <server> <key> <value>",
	Short: "Set a per-server configuration value",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		key := args[1]
		value := args[2]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		confPath := filepath.Join(s.Path, "server.conf")

		file, err := os.OpenFile(confPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		configKey := ""
		switch key {
		case "username":
			configKey = "USERNAME"
		case "screen_name":
			configKey = "SCREEN_NAME"
		case "jar_path":
			configKey = "JAR_PATH"
		case "ram":
			configKey = "RAM"
		case "stop_delay":
			configKey = "STOP_DELAY"
		case "restart_delay":
			configKey = "RESTART_DELAY"
		case "world_storage_path":
			configKey = "WORLD_STORAGE_PATH"
		case "world_storage_inactive_path":
			configKey = "WORLD_STORAGE_INACTIVE_PATH"
		default:
			return fmt.Errorf("unknown config key: %s\n  Valid keys: username, screen_name, jar_path, ram, stop_delay, restart_delay, world_storage_path, world_storage_inactive_path", key)
		}

		if _, err := fmt.Fprintf(file, "%s=\"%s\"\n", configKey, value); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Printf("Set %s=%s for server %q\n", key, value, serverName)
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update MSM to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		noScripts, _ := cmd.Flags().GetBool("no-scripts")

		fmt.Println("Checking for updates...")

		if noScripts {
			fmt.Println("Skipping script update (--no-scripts specified)")
			return nil
		}

		versionURL := cfg.UpdateURL + "/versioning/versions.txt"
		resp, err := http.Get(versionURL)
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to check for updates: status %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("Available versions:\n%s\n", string(body))
		fmt.Println("MSM Go rewrite - update functionality is simplified")
		fmt.Println("Download new binaries from the releases page")

		return nil
	},
}

func init() {
	updateCmd.Flags().Bool("no-scripts", false, "Skip updating scripts")
	logrollCmd.Flags().Bool("all", false, "Roll logs for all servers")

	rootCmd.AddCommand(logrollCmd)
	rootCmd.AddCommand(updateCmd)
}
