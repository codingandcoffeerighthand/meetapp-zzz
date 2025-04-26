package main

import (
	"fmt"
	"log"
	"proxy-srv/internal/proxy/app"
	"proxy-srv/internal/proxy/configs"
	"proxy-srv/internal/proxy/domain"
	infras_cloudflare "proxy-srv/internal/proxy/infra/cloudflare"
	"proxy-srv/internal/proxy/infra/smc_infra"
	"proxy-srv/internal/utils"

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
		fmt.Println(cfg)

		cleanUp := func() {}
		// clf
		clf, err := infras_cloudflare.NewClient(cfg)
		if err != nil {
			return err
		}
		// smc
		smc, err := smc_infra.NewSMCInfra(&cfg)
		if err != nil {
			return err
		}
		// // logger
		logger, cl, err := utils.NewLogger()
		if err != nil {
			return err
		}
		cleanUp = func() {
			cleanUp()
			cl()
		}
		// // crypt
		crypt, err := domain.NewDumDecryptService()
		if err != nil {
			return err
		}
		app, cl, err := app.NewApp(clf, smc, crypt, logger)
		if err != nil {
			return err
		}
		cleanUp = func() {
			cleanUp()
			cl()
		}
		fmt.Println("starting server...")
		app.Done()

		defer func() {
			cleanUp()
		}()
		// defer cleanUp()
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
