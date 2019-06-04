package main

import (
	"fmt"
	"os"

	"github.com/musicmash/artists/internal/commands/search"
	"github.com/musicmash/artists/internal/config"
	"github.com/musicmash/artists/internal/log"
	"github.com/spf13/cobra"
)

func main() {
	var configPath string
	var rootCmd = &cobra.Command{
		Use:           "artistsctl [OPTIONS] COMMAND [ARG...]",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := config.InitConfig(configPath); err != nil {
				return err
			}

			if config.Config.Log.Level == "" {
				config.Config.Log.Level = "info"
			}

			log.SetLogFormatter(&log.DefaultFormatter)
			log.ConfigureStdLogger(config.Config.Log.Level)
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "/etc/musicmash/artists/artists.yaml", "Path to config")
	rootCmd.AddCommand(search.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
