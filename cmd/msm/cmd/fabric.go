package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/msmhq/msm/internal/fabric"
	"github.com/msmhq/msm/internal/server"
	"github.com/spf13/cobra"
)

var fabricCmd = &cobra.Command{
	Use:   "fabric",
	Short: "Fabric mod loader management",
}

var fabricStatusCmd = &cobra.Command{
	Use:   "status <server>",
	Short: "Show Fabric status for a server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		fmt.Printf("Server: %s\n", s.Name)
		fmt.Printf("Fabric Enabled: %v\n", s.Config.FabricEnabled)

		if s.Config.FabricEnabled {
			mcVersion, err := s.DetectMCVersion()
			if err != nil {
				fmt.Printf("Minecraft Version: (unknown - %s)\n", err)
			} else {
				fmt.Printf("Minecraft Version: %s\n", mcVersion)
			}

			if s.Config.FabricLoaderVersion != "" {
				fmt.Printf("Loader Version: %s (pinned)\n", s.Config.FabricLoaderVersion)
			} else {
				fmt.Printf("Loader Version: (latest stable)\n")
			}

			if s.Config.FabricInstallerVersion != "" {
				fmt.Printf("Installer Version: %s (pinned)\n", s.Config.FabricInstallerVersion)
			} else {
				fmt.Printf("Installer Version: (latest stable)\n")
			}

			if mcVersion != "" {
				client, err := fabric.NewClient(cfg.FabricStoragePath, cfg.FabricCacheTTL)
				if err == nil {
					loaderVersion, installerVersion, err := client.ResolveVersions(
						mcVersion,
						s.Config.FabricLoaderVersion,
						s.Config.FabricInstallerVersion,
					)
					if err == nil {
						jarPath := fabric.JarPath(cfg.FabricStoragePath, mcVersion, loaderVersion, installerVersion)
						if _, err := os.Stat(jarPath); err == nil {
							fmt.Printf("Fabric JAR: %s (cached)\n", filepath.Base(jarPath))
						} else {
							fmt.Printf("Fabric JAR: (will download on start)\n")
						}
					} else {
						fmt.Printf("Fabric JAR: (error resolving versions: %s)\n", err)
					}
				} else {
					fmt.Printf("Fabric JAR: (error: %s)\n", err)
				}
			}
		}

		return nil
	},
}

var fabricSetCmd = &cobra.Command{
	Use:   "set <server> <on|off>",
	Short: "Enable or disable Fabric for a server",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		state := args[1]

		if state != "on" && state != "off" {
			return fmt.Errorf("invalid state %q: must be 'on' or 'off'", state)
		}

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if state == "on" {
			if s.Config.FabricEnabled {
				fmt.Printf("Fabric is already enabled for server %q\n", serverName)
				return nil
			}

			mcVersion, err := s.DetectMCVersion()
			if err != nil {
				return fmt.Errorf("cannot enable fabric: %w", err)
			}

			client, err := fabric.NewClient(cfg.FabricStoragePath, cfg.FabricCacheTTL)
			if err != nil {
				return err
			}

			supported, err := client.SupportsVersion(mcVersion)
			if err != nil {
				return fmt.Errorf("failed to check fabric compatibility: %w", err)
			}

			if !supported {
				return fmt.Errorf("fabric does not support minecraft %s", mcVersion)
			}

			if err := setServerConfigValue(s.Path, "FABRIC_ENABLED", "true"); err != nil {
				return err
			}

			fmt.Printf("Enabled Fabric for server %q (Minecraft %s)\n", serverName, mcVersion)
		} else {
			if !s.Config.FabricEnabled {
				fmt.Printf("Fabric is already disabled for server %q\n", serverName)
				return nil
			}

			if err := setServerConfigValue(s.Path, "FABRIC_ENABLED", "false"); err != nil {
				return err
			}

			fmt.Printf("Disabled Fabric for server %q\n", serverName)
		}

		return nil
	},
}

var fabricVersionsCmd = &cobra.Command{
	Use:   "versions [minecraft-version]",
	Short: "List available Fabric versions",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := fabric.NewClient(cfg.FabricStoragePath, cfg.FabricCacheTTL)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			versions, err := client.FetchGameVersions()
			if err != nil {
				return err
			}

			fmt.Println("Supported Minecraft versions:")
			count := 0
			for _, v := range versions {
				if v.Stable {
					fmt.Printf("  %s\n", v.Version)
					count++
					if count >= 20 {
						fmt.Printf("  ... and %d more (use 'msm fabric versions <mc-version>' for loader versions)\n", len(versions)-count)
						break
					}
				}
			}
			return nil
		}

		mcVersion := args[0]
		loaders, err := client.FetchLoaderVersions(mcVersion)
		if err != nil {
			return err
		}

		fmt.Printf("Fabric loader versions for Minecraft %s:\n", mcVersion)
		for _, l := range loaders {
			stability := ""
			if l.Stable {
				stability = " (stable)"
			}
			fmt.Printf("  %s%s\n", l.Version, stability)
		}

		return nil
	},
}

var fabricUpdateCmd = &cobra.Command{
	Use:   "update <server>",
	Short: "Check for and download newer Fabric loader",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if !s.Config.FabricEnabled {
			return fmt.Errorf("fabric is not enabled for server %q", serverName)
		}

		mcVersion, err := s.DetectMCVersion()
		if err != nil {
			return err
		}

		client, err := fabric.NewClient(cfg.FabricStoragePath, cfg.FabricCacheTTL)
		if err != nil {
			return err
		}

		loader, err := client.GetLatestStableLoader(mcVersion)
		if err != nil {
			return err
		}

		installer, err := client.GetLatestStableInstaller()
		if err != nil {
			return err
		}

		jarPath := fabric.JarPath(cfg.FabricStoragePath, mcVersion, loader.Version, installer.Version)

		if _, err := os.Stat(jarPath); err == nil {
			fmt.Printf("Latest Fabric loader %s for Minecraft %s is already cached\n", loader.Version, mcVersion)
			return nil
		}

		fmt.Printf("Downloading Fabric loader %s for Minecraft %s...\n", loader.Version, mcVersion)
		if _, err := client.DownloadServerJar(mcVersion, loader.Version, installer.Version); err != nil {
			return err
		}

		fmt.Printf("Downloaded Fabric loader %s\n", loader.Version)
		return nil
	},
}

var fabricSetLoaderCmd = &cobra.Command{
	Use:   "set-loader <server> <version>",
	Short: "Pin a specific Fabric loader version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, version := args[0], args[1]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if err := setServerConfigValue(s.Path, "FABRIC_LOADER_VERSION", version); err != nil {
			return err
		}

		fmt.Printf("Pinned Fabric loader version to %s for server %q\n", version, serverName)
		return nil
	},
}

var fabricSetInstallerCmd = &cobra.Command{
	Use:   "set-installer <server> <version>",
	Short: "Pin a specific Fabric installer version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, version := args[0], args[1]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if err := setServerConfigValue(s.Path, "FABRIC_INSTALLER_VERSION", version); err != nil {
			return err
		}

		fmt.Printf("Pinned Fabric installer version to %s for server %q\n", version, serverName)
		return nil
	},
}

func setServerConfigValue(serverPath, key, value string) error {
	confPath := filepath.Join(serverPath, "server.conf")

	var lines []string
	found := false

	file, err := os.Open(confPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if file != nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			if strings.HasPrefix(trimmed, key+"=") {
				lines = append(lines, fmt.Sprintf("%s=\"%s\"", key, value))
				found = true
			} else {
				lines = append(lines, line)
			}
		}
		file.Close()

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	return os.WriteFile(confPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func init() {
	fabricCmd.AddCommand(fabricStatusCmd)
	fabricCmd.AddCommand(fabricSetCmd)
	fabricCmd.AddCommand(fabricVersionsCmd)
	fabricCmd.AddCommand(fabricUpdateCmd)
	fabricCmd.AddCommand(fabricSetLoaderCmd)
	fabricCmd.AddCommand(fabricSetInstallerCmd)

	rootCmd.AddCommand(fabricCmd)
}
