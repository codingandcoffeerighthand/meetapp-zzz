package main

import (
	"fmt"
	"log"
	command "proxy-srv/cmd/server/command"
	"proxy-srv/internal/proxy/biz"
	"proxy-srv/internal/proxy/configs"
	"proxy-srv/internal/proxy/grpc"
	infras_cloudflare "proxy-srv/internal/proxy/infra/cloudflare"
	"proxy-srv/internal/proxy/infra/smc_infra"
	"proxy-srv/internal/proxy/infra/ws_grpc"
	crypt_srv "proxy-srv/internal/proxy/service/crypt"
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
		// logger
		logger, cl, err := utils.NewLogger()
		if err != nil {
			return err
		}
		cleanUp = func() {
			cleanUp()
			cl()
		}
		// ws grpc client
		grpcCl := ws_grpc.NewClient(cfg, logger)
		// crypt
		crypt, err := crypt_srv.NewDumDecryptService()
		if err != nil {
			return err
		}

		biz, cl, err := biz.NewBiz(clf, smc, grpcCl, crypt, logger)
		if err != nil {
			return err
		}
		cleanUp = func() {
			cleanUp()
			cl()
		}
		app := grpc.NewServerGrpc(biz, cfg.ServerConfig)
		defer cleanUp()
		return app.Listen()

	},
}

func main() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s-%s", version, hash),
	}
	runCommand.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file.")
	rootCommand.AddCommand(runCommand)
	rootCommand.AddCommand(command.WsCommand())
	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
