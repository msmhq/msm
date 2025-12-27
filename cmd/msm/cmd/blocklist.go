package cmd

import (
	"fmt"

	"github.com/msmhq/msm/internal/logging"
	"github.com/msmhq/msm/internal/player"
	"github.com/msmhq/msm/internal/server"
	"github.com/spf13/cobra"
)

var blocklistCmd = &cobra.Command{
	Use:     "blocklist <server>",
	Aliases: []string{"bl"},
	Short:   "Blocklist (ban) management commands",
	Args:    cobra.MinimumNArgs(1),
}

var blocklistPlayerCmd = &cobra.Command{
	Use:   "player",
	Short: "Player ban management",
}

var blocklistPlayerAddCmd = &cobra.Command{
	Use:   "add <server> <player>...",
	Short: "Ban players",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		players := args[1:]
		reason, _ := cmd.Flags().GetString("reason")

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		for _, p := range players {
			if err := player.BanPlayer(s.BannedPlayersPath(), p, reason); err != nil {
				logging.Error("failed to ban player", "player", p, "error", err)
				continue
			}

			if s.IsRunning() {
				var banCmd string
				if reason != "" {
					banCmd = fmt.Sprintf("ban %s %s", p, reason)
				} else {
					banCmd = fmt.Sprintf("ban %s", p)
				}
				if err := s.SendCommand(banCmd); err != nil {
					logging.Warn("failed to send ban command", "player", p, "error", err)
				}
			}

			logging.Info("banned player", "player", p)
		}

		return nil
	},
}

var blocklistPlayerRemoveCmd = &cobra.Command{
	Use:   "remove <server> <player>...",
	Short: "Unban players",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		players := args[1:]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		for _, p := range players {
			if err := player.UnbanPlayer(s.BannedPlayersPath(), p); err != nil {
				logging.Error("failed to unban player", "player", p, "error", err)
				continue
			}

			if s.IsRunning() {
				if err := s.SendCommand(fmt.Sprintf("pardon %s", p)); err != nil {
					logging.Warn("failed to send pardon command", "player", p, "error", err)
				}
			}

			logging.Info("unbanned player", "player", p)
		}

		return nil
	},
}

var blocklistPlayerListCmd = &cobra.Command{
	Use:   "list <server>",
	Short: "List banned players",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		entries, err := player.LoadBannedPlayers(s.BannedPlayersPath())
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			fmt.Println("No banned players.")
			return nil
		}

		fmt.Printf("Banned players for server %q:\n", serverName)
		for _, e := range entries {
			if e.Reason != "" {
				fmt.Printf("  %s (reason: %s)\n", e.Name, e.Reason)
			} else {
				fmt.Printf("  %s\n", e.Name)
			}
		}

		return nil
	},
}

var blocklistIPCmd = &cobra.Command{
	Use:   "ip",
	Short: "IP ban management",
}

var blocklistIPAddCmd = &cobra.Command{
	Use:   "add <server> <ip>...",
	Short: "Ban IP addresses",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		ips := args[1:]
		reason, _ := cmd.Flags().GetString("reason")

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		for _, ip := range ips {
			if err := player.BanIP(s.BannedIPsPath(), ip, reason); err != nil {
				logging.Error("failed to ban IP", "ip", ip, "error", err)
				continue
			}

			if s.IsRunning() {
				var banCmd string
				if reason != "" {
					banCmd = fmt.Sprintf("ban-ip %s %s", ip, reason)
				} else {
					banCmd = fmt.Sprintf("ban-ip %s", ip)
				}
				if err := s.SendCommand(banCmd); err != nil {
					logging.Warn("failed to send ban-ip command", "ip", ip, "error", err)
				}
			}

			logging.Info("banned IP", "ip", ip)
		}

		return nil
	},
}

var blocklistIPRemoveCmd = &cobra.Command{
	Use:   "remove <server> <ip>...",
	Short: "Unban IP addresses",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		ips := args[1:]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		for _, ip := range ips {
			if err := player.UnbanIP(s.BannedIPsPath(), ip); err != nil {
				logging.Error("failed to unban IP", "ip", ip, "error", err)
				continue
			}

			if s.IsRunning() {
				if err := s.SendCommand(fmt.Sprintf("pardon-ip %s", ip)); err != nil {
					logging.Warn("failed to send pardon-ip command", "ip", ip, "error", err)
				}
			}

			logging.Info("unbanned IP", "ip", ip)
		}

		return nil
	},
}

var blocklistIPListCmd = &cobra.Command{
	Use:   "list <server>",
	Short: "List banned IP addresses",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		s, err := server.Get(serverName, cfg)
		if err != nil {
			return err
		}

		entries, err := player.LoadBannedIPs(s.BannedIPsPath())
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			fmt.Println("No banned IPs.")
			return nil
		}

		fmt.Printf("Banned IPs for server %q:\n", serverName)
		for _, e := range entries {
			if e.Reason != "" {
				fmt.Printf("  %s (reason: %s)\n", e.IP, e.Reason)
			} else {
				fmt.Printf("  %s\n", e.IP)
			}
		}

		return nil
	},
}

func init() {
	blocklistPlayerAddCmd.Flags().String("reason", "", "Ban reason")
	blocklistIPAddCmd.Flags().String("reason", "", "Ban reason")

	blocklistPlayerCmd.AddCommand(blocklistPlayerAddCmd)
	blocklistPlayerCmd.AddCommand(blocklistPlayerRemoveCmd)
	blocklistPlayerCmd.AddCommand(blocklistPlayerListCmd)

	blocklistIPCmd.AddCommand(blocklistIPAddCmd)
	blocklistIPCmd.AddCommand(blocklistIPRemoveCmd)
	blocklistIPCmd.AddCommand(blocklistIPListCmd)

	blocklistCmd.AddCommand(blocklistPlayerCmd)
	blocklistCmd.AddCommand(blocklistIPCmd)

	rootCmd.AddCommand(blocklistCmd)
}
