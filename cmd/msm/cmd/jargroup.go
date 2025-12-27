package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/msmhq/msm/internal/fabric"
	"github.com/msmhq/msm/internal/jar"
	"github.com/msmhq/msm/internal/mojang"
	"github.com/msmhq/msm/internal/server"
	"github.com/spf13/cobra"
)

var jargroupCmd = &cobra.Command{
	Use:   "jargroup",
	Short: "Jar group management commands",
}

var jargroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all jar groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		groups, err := jar.DiscoverAll(cfg)
		if err != nil {
			return err
		}

		if len(groups) == 0 {
			fmt.Println("No jar groups found.")
			return nil
		}

		for _, g := range groups {
			fmt.Printf("%s", g.Name)
			if g.URL != "" {
				fmt.Printf(" (%s)", g.URL)
			}
			fmt.Println()
			for _, f := range g.Files {
				fmt.Printf("  - %s\n", f)
			}
		}
		return nil
	},
}

var jargroupCreateCmd = &cobra.Command{
	Use:   "create <name> <url>",
	Short: "Create a new jar group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, url := args[0], args[1]
		g, err := jar.Create(name, url, cfg)
		if err != nil {
			return err
		}
		fmt.Printf("Created jar group %q with URL %s\n", g.Name, g.URL)
		return nil
	},
}

var jargroupDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a jar group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := jar.Delete(name, cfg); err != nil {
			return err
		}
		fmt.Printf("Deleted jar group %q\n", name)
		return nil
	},
}

var jargroupRenameCmd = &cobra.Command{
	Use:   "rename <old-name> <new-name>",
	Short: "Rename a jar group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldName, newName := args[0], args[1]
		if err := jar.Rename(oldName, newName, cfg); err != nil {
			return err
		}
		fmt.Printf("Renamed jar group %q to %q\n", oldName, newName)
		return nil
	},
}

var jargroupChangeURLCmd = &cobra.Command{
	Use:   "changeurl <name> <url>",
	Short: "Change the URL for a jar group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, url := args[0], args[1]
		g, err := jar.Get(name, cfg)
		if err != nil {
			return err
		}
		if err := g.ChangeURL(url); err != nil {
			return err
		}
		fmt.Printf("Changed URL for jar group %q to %s\n", name, url)
		return nil
	},
}

var jargroupGetLatestCmd = &cobra.Command{
	Use:   "getlatest <name>",
	Short: "Download the latest jar for a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		g, err := jar.Get(name, cfg)
		if err != nil {
			return err
		}
		filename, err := g.GetLatest()
		if err != nil {
			return err
		}
		fmt.Printf("Downloaded %s to jar group %q\n", filename, name)
		return nil
	},
}

var jarCmd = &cobra.Command{
	Use:   "jar <server>",
	Short: "Show or manage server jar",
	Long: `Show or manage server jar.

Without a subcommand, shows the current jar path.
Use 'jar link' to link to a jar group, or 'jar download' for vanilla.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		fmt.Printf("Server %q jar: %s\n", serverName, s.Config.JarPath)
		return nil
	},
}

var jarLinkCmd = &cobra.Command{
	Use:   "link <server> <jargroup> [file]",
	Short: "Link a server to a jar from a jar group",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		jargroupName := args[1]
		var jarFile string
		if len(args) > 2 {
			jarFile = args[2]
		}

		force, _ := cmd.Flags().GetBool("force")

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if s.Config.FabricEnabled && !force {
			g, err := jar.Get(jargroupName, cfg)
			if err != nil {
				return err
			}

			targetFile := jarFile
			if targetFile == "" {
				targetFile = g.LatestFile()
			}

			newVersion := extractMCVersionFromJarName(targetFile)
			if newVersion != "" {
				client, err := fabric.NewClient(cfg.FabricStoragePath, cfg.FabricCacheTTL)
				if err != nil {
					return fmt.Errorf("failed to check fabric compatibility: %w", err)
				}

				supported, err := client.SupportsVersion(newVersion)
				if err != nil {
					return fmt.Errorf("failed to check fabric compatibility: %w", err)
				}

				if !supported {
					return fmt.Errorf("fabric does not yet support minecraft %s - upgrade blocked\n  Hint: Use --force to override (may cause issues)", newVersion)
				}
			}
		}

		if err := jar.LinkJar(s.Path, s.Config.JarPath, jargroupName, jarFile, cfg); err != nil {
			return err
		}

		if jarFile != "" {
			fmt.Printf("Linked server %q to jar %s from group %q\n", serverName, jarFile, jargroupName)
		} else {
			fmt.Printf("Linked server %q to latest jar from group %q\n", serverName, jargroupName)
		}
		return nil
	},
}

func extractMCVersionFromJarName(filename string) string {
	patterns := []struct {
		prefix string
		suffix string
	}{
		{"minecraft_server.", ".jar"},
		{"server-", ".jar"},
		{"paper-", ".jar"},
		{"spigot-", ".jar"},
		{"craftbukkit-", ".jar"},
		{"purpur-", ".jar"},
	}

	for _, p := range patterns {
		if len(filename) > len(p.prefix)+len(p.suffix) {
			if filename[:len(p.prefix)] == p.prefix {
				rest := filename[len(p.prefix):]
				if idx := len(rest) - len(p.suffix); idx > 0 && rest[idx:] == p.suffix {
					version := rest[:idx]
					if idx := indexOf(version, '-'); idx != -1 {
						version = version[:idx]
					}
					return version
				}
			}
		}
	}

	return ""
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

var jarDownloadCmd = &cobra.Command{
	Use:   "download <server> [version]",
	Short: "Download vanilla Minecraft server jar",
	Long: `Download the official Minecraft server jar directly to a server.

If no version is specified, downloads the latest release.
This is the simplest way to set up a new server - no jar groups needed.

Jars are cached centrally in the jar storage path, so multiple servers
using the same version share a single downloaded file.

Examples:
  msm jar download survival           # Latest release
  msm jar download survival 1.21.4    # Specific version`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		var version string
		if len(args) > 1 {
			version = args[1]
		} else {
			fmt.Print("Fetching latest Minecraft version... ")
			latest, err := mojang.GetLatestRelease()
			if err != nil {
				return err
			}
			version = latest
			fmt.Println(version)
		}

		cachedPath := mojang.CachedJarPath(cfg.JarStoragePath, version)
		cached := false
		if _, err := os.Stat(cachedPath); err == nil {
			cached = true
			fmt.Printf("Using cached Minecraft %s server jar\n", version)
		} else {
			fmt.Printf("Downloading Minecraft %s server... ", version)
		}

		jarPath, err := mojang.EnsureCached(cfg.JarStoragePath, version)
		if err != nil {
			return err
		}

		if !cached {
			fmt.Println("Done.")
		}

		destPath := filepath.Join(s.Path, s.Config.JarPath)
		if _, err := os.Lstat(destPath); err == nil {
			if err := os.Remove(destPath); err != nil {
				return fmt.Errorf("failed to remove existing jar: %w", err)
			}
		}

		if err := os.Symlink(jarPath, destPath); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}

		fmt.Printf("Server %q is now set up with Minecraft %s\n", serverName, version)
		return nil
	},
}

func init() {
	jargroupCmd.AddCommand(jargroupListCmd)
	jargroupCmd.AddCommand(jargroupCreateCmd)
	jargroupCmd.AddCommand(jargroupDeleteCmd)
	jargroupCmd.AddCommand(jargroupRenameCmd)
	jargroupCmd.AddCommand(jargroupChangeURLCmd)
	jargroupCmd.AddCommand(jargroupGetLatestCmd)

	jarLinkCmd.Flags().Bool("force", false, "Force jar change even if Fabric doesn't support the new version")
	jarCmd.AddCommand(jarLinkCmd)
	jarCmd.AddCommand(jarDownloadCmd)

	rootCmd.AddCommand(jargroupCmd)
	rootCmd.AddCommand(jarCmd)
}
