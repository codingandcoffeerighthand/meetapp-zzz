package command

import (
	"proxy-srv/internal/utils"
	ws_app "proxy-srv/internal/ws_p/app"
	"proxy-srv/internal/ws_p/config"
	ws_grpc_server "proxy-srv/internal/ws_p/delivery/grpc"
	ws_server "proxy-srv/internal/ws_p/delivery/ws"
	ws_infra "proxy-srv/internal/ws_p/infra"

	"github.com/spf13/cobra"
)

const (
	flagConfigFileWsPath = "config-file-ws-path"
)

func WsCommand() *cobra.Command {
	wsCmd := &cobra.Command{
		Use:   "ws",
		Short: "websocket server",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigFileWsPath)
			if err != nil {
				return err
			}

			cfg, err := config.NewConfig(config.ConfigPath(configPath))
			if err != nil {
				return err
			}

			cleanUp := func() {}

			gCl, err := ws_infra.NewClient(cfg.ProxyConfig)
			if err != nil {
				return err
			}
			logger, cl, err := utils.NewLogger()
			if err != nil {
				return err
			}
			cleanUp = func() {
				cleanUp()
				cl()
			}

			app := ws_app.NewApp(gCl, logger)

			ws := ws_server.NewServer(cfg.ServerConfig, app)
			grpcServere := ws_grpc_server.NewGrpcServer(cfg.ServerConfig, app)

			errChan := make(chan error)

			go func() {
				err := grpcServere.Listen()
				errChan <- err
			}()
			err = ws.Listen()
			errChan <- err
			defer cleanUp()
			return <-errChan
		},
	}
	wsCmd.Flags().String(flagConfigFileWsPath, "", "If provided, will use the provided config file.	")
	return wsCmd
}
