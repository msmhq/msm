package cmd

import (
	"fmt"

	"github.com/msmhq/msm/internal/logging"
	"github.com/msmhq/msm/internal/player"
	"github.com/msmhq/msm/internal/server"
	"github.com/spf13/cobra"
)

var allowlistCmd = &cobra.Command{
	Use:     "allowlist <server>",
	Aliases: []string{"al"},
	Short:   "Allowlist management commands",
	Args:    cobra.MinimumNArgs(1),
}

var allowlistSetCmd = &cobra.Command{
	Use:   "set <server> <on|off>",
	Short: "Enable or disable the allowlist (whitelist)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, state := args[0], args[1]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if !s.IsRunning() {
			return fmt.Errorf("server %q is not running", serverName)
		}

		switch state {
		case "on":
			if err := s.SendCommand("whitelist on"); err != nil {
				return err
			}
			fmt.Printf("Allowlist enabled for server %q\n", serverName)
		case "off":
			if err := s.SendCommand("whitelist off"); err != nil {
				return err
			}
			fmt.Printf("Allowlist disabled for server %q\n", serverName)
		default:
			return fmt.Errorf("invalid state %q: must be 'on' or 'off'", state)
		}
		return nil
	},
}

var allowlistAddCmd = &cobra.Command{
	Use:   "add <server> <player>...",
	Short: "Add players to the allowlist",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		players := args[1:]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		for _, p := range players {
			if err := player.AddToAllowlist(s.AllowlistPath(), p); err != nil {
				logging.Error("failed to add player to allowlist", "player", p, "error", err)
				continue
			}

			if s.IsRunning() {
				if err := s.SendCommand(fmt.Sprintf("whitelist add %s", p)); err != nil {
					logging.Warn("failed to send whitelist command", "player", p, "error", err)
				}
			}

			logging.Info("added to allowlist", "player", p)
		}

		return nil
	},
}

var allowlistRemoveCmd = &cobra.Command{
	Use:   "remove <server> <player>...",
	Short: "Remove players from the allowlist",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		players := args[1:]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		for _, p := range players {
			if err := player.RemoveFromAllowlist(s.AllowlistPath(), p); err != nil {
				logging.Error("failed to remove player from allowlist", "player", p, "error", err)
				continue
			}

			if s.IsRunning() {
				if err := s.SendCommand(fmt.Sprintf("whitelist remove %s", p)); err != nil {
					logging.Warn("failed to send whitelist command", "player", p, "error", err)
				}
			}

			logging.Info("removed from allowlist", "player", p)
		}

		return nil
	},
}

var allowlistListCmd = &cobra.Command{
	Use:   "list <server>",
	Short: "List all players on the allowlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		entries, err := player.LoadAllowlist(s.AllowlistPath())
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			fmt.Println("Allowlist is empty.")
			return nil
		}

		fmt.Printf("Allowlist for server %q:\n", serverName)
		for _, e := range entries {
			fmt.Printf("  %s\n", e.Name)
		}

		return nil
	},
}

var opCmd = &cobra.Command{
	Use:   "op",
	Short: "Operator management commands",
}

var opAddCmd = &cobra.Command{
	Use:   "add <server> <player>",
	Short: "Make a player an operator",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, playerName := args[0], args[1]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if err := player.AddOp(s.OpsPath(), playerName, 4); err != nil {
			return err
		}

		if s.IsRunning() {
			if err := s.SendCommand(fmt.Sprintf("op %s", playerName)); err != nil {
				logging.Warn("failed to send op command", "player", playerName, "error", err)
			}
		}

		logging.Info("made player an operator", "player", playerName, "server", serverName)
		return nil
	},
}

var opRemoveCmd = &cobra.Command{
	Use:   "remove <server> <player>",
	Short: "Remove operator status from a player",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName, playerName := args[0], args[1]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		if err := player.RemoveOp(s.OpsPath(), playerName); err != nil {
			return err
		}

		if s.IsRunning() {
			if err := s.SendCommand(fmt.Sprintf("deop %s", playerName)); err != nil {
				logging.Warn("failed to send deop command", "player", playerName, "error", err)
			}
		}

		logging.Info("removed operator status", "player", playerName, "server", serverName)
		return nil
	},
}

var opListCmd = &cobra.Command{
	Use:   "list <server>",
	Short: "List all operators for a server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		entries, err := player.LoadOps(s.OpsPath())
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			fmt.Println("No operators.")
			return nil
		}

		fmt.Printf("Operators for server %q:\n", serverName)
		for _, e := range entries {
			fmt.Printf("  %s (level %d)\n", e.Name, e.Level)
		}

		return nil
	},
}

func init() {
	allowlistCmd.AddCommand(allowlistSetCmd)
	allowlistCmd.AddCommand(allowlistAddCmd)
	allowlistCmd.AddCommand(allowlistRemoveCmd)
	allowlistCmd.AddCommand(allowlistListCmd)

	opCmd.AddCommand(opAddCmd)
	opCmd.AddCommand(opRemoveCmd)
	opCmd.AddCommand(opListCmd)

	rootCmd.AddCommand(allowlistCmd)
	rootCmd.AddCommand(opCmd)
}
