package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/msmhq/msm/internal/logging"
	"github.com/msmhq/msm/internal/server"
	"github.com/msmhq/msm/internal/world"
	"github.com/spf13/cobra"
)

var worldsCmd = &cobra.Command{
	Use:   "worlds <server>",
	Short: "World management commands",
	Args:  cobra.MinimumNArgs(1),
}

var worldsListCmd = &cobra.Command{
	Use:   "list <server>",
	Short: "List all worlds for a server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		worlds, err := world.DiscoverAll(s.Path, s.Name, cfg, s.Config.WorldStoragePath, s.Config.WorldStorageInactivePath)
		if err != nil {
			return err
		}

		if len(worlds) == 0 {
			fmt.Printf("No worlds found for server %q.\n", serverName)
			return nil
		}

		mcVersion := "-"
		if v, err := s.DetectMCVersion(); err == nil {
			mcVersion = v
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "NAME\tSERVER\tLOCATION\tVERSION")

		for _, w := range worlds {
			location := "disk"
			if w.InRAM {
				location = "RAM"
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", w.Name, serverName, location, mcVersion)
		}
		tw.Flush()
		return nil
	},
}

var worldsSetCmd = &cobra.Command{
	Use:   "set <server> <world> <on|off>",
	Short: "Activate or deactivate a world",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, worldName, state := args[0], args[1], args[2]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		w, err := world.Get(s.Path, s.Name, worldName, s.Config.WorldStoragePath, s.Config.WorldStorageInactivePath, cfg)
		if err != nil {
			return err
		}

		switch state {
		case "on":
			if err := w.Activate(s.Config.WorldStoragePath); err != nil {
				return err
			}
			fmt.Printf("Activated world %q for server %q\n", worldName, serverName)
		case "off":
			if err := w.Deactivate(s.Config.WorldStorageInactivePath); err != nil {
				return err
			}
			fmt.Printf("Deactivated world %q for server %q\n", worldName, serverName)
		default:
			return fmt.Errorf("invalid state %q: must be 'on' or 'off'", state)
		}
		return nil
	},
}

var worldsRAMCmd = &cobra.Command{
	Use:   "ram <server> <world>",
	Short: "Show or manage RAM disk state for a world",
	Long: `Show or manage RAM disk state for a world.

Without a subcommand, shows the current RAM status.
Use 'ram set <server> <world> on|off' to change the state.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, worldName := args[0], args[1]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		w, err := world.Get(s.Path, s.Name, worldName, s.Config.WorldStoragePath, s.Config.WorldStorageInactivePath, cfg)
		if err != nil {
			return err
		}

		if w.InRAM {
			fmt.Printf("World %q is in RAM\n", w.Name)
			fmt.Printf("  RAM path: %s\n", w.RAMPath)
		} else {
			fmt.Printf("World %q is not in RAM\n", w.Name)
		}
		return nil
	},
}

var worldsRAMSetCmd = &cobra.Command{
	Use:   "set <server> <world> <on|off>",
	Short: "Enable or disable RAM disk for a world",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, worldName, state := args[0], args[1], args[2]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		w, err := world.Get(s.Path, s.Name, worldName, s.Config.WorldStoragePath, s.Config.WorldStorageInactivePath, cfg)
		if err != nil {
			return err
		}

		switch state {
		case "on":
			if err := w.EnableRAM(s.Config.Username); err != nil {
				return err
			}
			fmt.Printf("Enabled RAM disk for world %q\n", worldName)
		case "off":
			if err := w.DisableRAM(s.Config.Username); err != nil {
				return err
			}
			fmt.Printf("Disabled RAM disk for world %q\n", worldName)
		default:
			return fmt.Errorf("invalid state %q: must be 'on' or 'off'", state)
		}
		return nil
	},
}

var worldsToDiskCmd = &cobra.Command{
	Use:   "todisk [server]",
	Short: "Sync all RAM worlds to disk",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		if all || len(args) == 0 {
			return syncAllServers()
		}

		serverName := args[0]
		return syncServerToDisk(serverName)
	},
}

func syncServerToDisk(serverName string) error {
	s, err := server.Get(serverName, cfg)
	if err != nil {
		return err
	}

	if s.IsRunning() {
		s.SaveOff()
		s.SaveAll()
		defer s.SaveOn()
	}

	worlds, err := world.DiscoverAll(s.Path, s.Name, cfg, s.Config.WorldStoragePath, s.Config.WorldStorageInactivePath)
	if err != nil {
		return err
	}

	synced := 0
	for _, w := range worlds {
		if w.InRAM {
			fmt.Fprintf(os.Stderr, "Syncing world %q to disk...\n", w.Name)
			if err := w.ToDisk(s.Config.Username); err != nil {
				logging.Error("failed to sync world", "world", w.Name, "error", err)
				continue
			}
			synced++
		}
	}

	if synced == 0 {
		fmt.Fprintln(os.Stderr, "No RAM worlds to sync.")
	}
	return nil
}

func syncAllServers() error {
	servers, err := server.DiscoverAll(cfg)
	if err != nil {
		return err
	}

	for _, s := range servers {
		if s.IsRunning() {
			if err := syncServerToDisk(s.Name); err != nil {
				logging.Warn("failed to sync server", "server", s.Name, "error", err)
			}
		}
	}
	return nil
}

var worldsBackupCmd = &cobra.Command{
	Use:   "backup [server]",
	Short: "Backup all worlds for a server",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		if all || len(args) == 0 {
			return backupAllServers()
		}

		serverName := args[0]
		return backupServerWorlds(serverName)
	},
}

func backupServerWorlds(serverName string) error {
	s, err := server.Get(serverName, cfg)
	if err != nil {
		return err
	}

	if s.IsRunning() {
		s.Say(s.Config.MessageWorldBackupStarted)
		s.SaveOff()
		s.SaveAll()
		defer func() {
			s.SaveOn()
			s.Say(s.Config.MessageWorldBackupFinished)
		}()
	}

	worlds, err := world.DiscoverAll(s.Path, s.Name, cfg, s.Config.WorldStoragePath, s.Config.WorldStorageInactivePath)
	if err != nil {
		return err
	}

	for _, w := range worlds {
		if w.InRAM {
			if err := w.ToDisk(s.Config.Username); err != nil {
				logging.Warn("failed to sync world to disk", "world", w.Name, "error", err)
			}
		}

		if err := w.Backup("", s.Config.Username); err != nil {
			logging.Error("failed to backup world", "world", w.Name, "error", err)
		}
	}

	return nil
}

func backupAllServers() error {
	servers, err := server.DiscoverAll(cfg)
	if err != nil {
		return err
	}

	for _, s := range servers {
		logging.Info("backing up worlds", "server", s.Name)
		if err := backupServerWorlds(s.Name); err != nil {
			logging.Warn("failed to backup server", "server", s.Name, "error", err)
		}
	}
	return nil
}

var serverBackupCmd = &cobra.Command{
	Use:   "backup <server>",
	Short: "Backup entire server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if s.IsRunning() {
			s.Say(s.Config.MessageCompleteBackupStarted)
			s.SaveOff()
			s.SaveAll()
			defer func() {
				s.SaveOn()
				s.Say(s.Config.MessageCompleteBackupFinished)
			}()
		}

		worlds, err := world.DiscoverAll(s.Path, s.Name, cfg, s.Config.WorldStoragePath, s.Config.WorldStorageInactivePath)
		if err != nil {
			return err
		}

		for _, w := range worlds {
			if w.InRAM {
				if err := w.ToDisk(s.Config.Username); err != nil {
					logging.Warn("failed to sync world to disk", "world", w.Name, "error", err)
				}
			}
		}

		return world.BackupServer(s.Path, s.Name, cfg.BackupArchivePath, s.Config.Username, false)
	},
}

func init() {
	worldsToDiskCmd.Flags().Bool("all", false, "Sync all running servers")
	worldsBackupCmd.Flags().Bool("all", false, "Backup all servers")

	worldsRAMCmd.AddCommand(worldsRAMSetCmd)

	worldsCmd.AddCommand(worldsListCmd)
	worldsCmd.AddCommand(worldsSetCmd)
	worldsCmd.AddCommand(worldsRAMCmd)
	worldsCmd.AddCommand(worldsToDiskCmd)
	worldsCmd.AddCommand(worldsBackupCmd)

	rootCmd.AddCommand(worldsCmd)
	rootCmd.AddCommand(serverBackupCmd)
}
