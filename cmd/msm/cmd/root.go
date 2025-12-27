package cmd

import (
	"fmt"

	"github.com/msmhq/msm/internal/config"
	"github.com/msmhq/msm/internal/logging"
	"github.com/msmhq/msm/internal/server"
	"github.com/spf13/cobra"
)

const Version = "0.12.0"

var (
	cfg     *config.Config
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "msm",
	Short: "Minecraft Server Manager",
	Long: `MSM (Minecraft Server Manager) is a comprehensive management tool
for running multiple Minecraft servers on a single machine.

It provides unified management including server lifecycle control,
world backups, RAM disk caching, and player management.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logging.Init(verbose)

		if cmd.Name() == "version" || cmd.Name() == "help" {
			return nil
		}

		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("MSM version %s\n", Version)
	},
}

var startCmd = &cobra.Command{
	Use:   "start [server]",
	Short: "Start a server, or all servers if none specified",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			s, err := server.Get(args[0], cfg)
			if err != nil {
				return err
			}
			if err := s.Start(); err != nil {
				return err
			}
			if port := s.Port(); port > 0 {
				fmt.Printf("Started server %q on port %d\n", s.Name, port)
			} else {
				fmt.Printf("Started server %q\n", s.Name)
			}
			return nil
		}

		servers, err := server.DiscoverAll(cfg)
		if err != nil {
			return err
		}

		for _, s := range servers {
			if err := s.Start(); err != nil {
				logging.Error("failed to start server", "server", s.Name, "error", err)
			} else {
				port := s.Port()
				if port > 0 {
					logging.Info("started server", "server", s.Name, "port", port)
				} else {
					logging.Info("started server", "server", s.Name)
				}
			}
		}
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [server]",
	Short: "Stop a server, or all servers if none specified",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		now, _ := cmd.Flags().GetBool("now")

		if len(args) == 1 {
			s, err := server.Get(args[0], cfg)
			if err != nil {
				return err
			}
			if err := s.Stop(now); err != nil {
				return err
			}
			fmt.Printf("Stopped server %q\n", s.Name)
			return nil
		}

		servers, err := server.DiscoverAll(cfg)
		if err != nil {
			return err
		}

		for _, s := range servers {
			if err := s.Stop(now); err != nil {
				logging.Error("failed to stop server", "server", s.Name, "error", err)
			} else {
				logging.Info("stopped server", "server", s.Name)
			}
		}
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart [server]",
	Short: "Restart a server, or all servers if none specified",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		now, _ := cmd.Flags().GetBool("now")

		if len(args) == 1 {
			s, err := server.Get(args[0], cfg)
			if err != nil {
				return err
			}
			if err := s.Restart(now); err != nil {
				return err
			}
			fmt.Printf("Restarted server %q\n", s.Name)
			return nil
		}

		servers, err := server.DiscoverAll(cfg)
		if err != nil {
			return err
		}

		for _, s := range servers {
			if err := s.Restart(now); err != nil {
				logging.Error("failed to restart server", "server", s.Name, "error", err)
			} else {
				logging.Info("restarted server", "server", s.Name)
			}
		}
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display global configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.Print()
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/msm.conf)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose/debug output")

	stopCmd.Flags().Bool("now", false, "Stop immediately without warning players")
	restartCmd.Flags().Bool("now", false, "Restart immediately without warning players")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(configCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func GetConfig() *config.Config {
	return cfg
}
