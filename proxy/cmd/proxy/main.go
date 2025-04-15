package main

import (
	"context"
	"fmt"
	"log"
	"proxy-srv/internal/configs"
	infras_cloudflare "proxy-srv/internal/infra/client/cloudflare"

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
	Use:   "proxy",
	Short: "proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString(flagConfigFilePath)
		if err != nil {
			return err
		}

		cfg, err := configs.NewConfig(configs.ConfigPath(configPath))
		if err != nil {
			return err
		}
		// fmt.Println(cfg)
		cl, error := infras_cloudflare.NewClient(cfg.CloudflareConfig)
		if error != nil {
			return error
		}
		resp, err := cl.PostAppsAppIdSessionsNewWithResponse(
			context.Background(),
			cfg.CloudflareConfig.AppId)
		if err != nil {
			return err
		}
		fmt.Println(resp.HTTPResponse.Status)
		fmt.Println(string(resp.Body))
		fmt.Println(*resp.JSON201.SessionId)
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
