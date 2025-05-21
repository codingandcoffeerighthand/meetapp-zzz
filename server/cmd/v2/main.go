package main

import (
	"context"
	"fmt"
	"log"
	"proxy-srv/cmd/v2/initial"
	"proxy-srv/internal/v2/configs"

	"github.com/spf13/cobra"
)

var (
	version = "v0.0.1"
	hash    = "0x0"
)

const (
	flagConfigFilePath = "config-file-path"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "run",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString(flagConfigFilePath)
		if err != nil {
			return err
		}
		cfg, err := configs.NewConfig(configs.ConfigPath(configPath))
		if err != nil {
			return err
		}
		app, _, err := initial.NewApp(&cfg)
		if err != nil {
			return err
		}
		app.Run(context.Background())
		app.Done()
		return nil
	},
}

func main() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s-%s", version, hash),
	}
	runCommand.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file.")
	rootCommand.AddCommand(runCommand)
	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
