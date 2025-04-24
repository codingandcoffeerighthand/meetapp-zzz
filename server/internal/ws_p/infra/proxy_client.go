package ws_infra

import (
	"context"
	"proxy-srv/internal/ws_p/config"
	"proxy-srv/pkg/gencode/proxy_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	conn    *grpc.ClientConn
	service proxy_grpc.ProxyServiceClient
}

func NewClient(cfg config.ProxyConfig) (*client, error) {
	conn, err := grpc.NewClient(cfg.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		conn:    conn,
		service: proxy_grpc.NewProxyServiceClient(conn),
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) RenegotiateSession(sessionId string, sdpAnswer string) error {
	_, err := c.service.PUTRenegotiateSession(context.Background(), &proxy_grpc.PUTRenegotiateSessionRequest{
		SessionId: sessionId,
		SdpAnswer: sdpAnswer,
	})
	return err
}
func (c *client) RoomPull(roomId string) (err error) {
	_, err = c.service.EmitParicipantJoinedToRom(context.Background(), &proxy_grpc.EmitParicipantJoinedToRomRequest{
		RoomId: roomId,
	})
	return err
}
