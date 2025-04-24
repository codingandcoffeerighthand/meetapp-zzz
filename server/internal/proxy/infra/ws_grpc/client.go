package ws_grpc

import (
	"context"
	"fmt"
	"proxy-srv/internal/proxy/configs"
	"proxy-srv/pkg/gencode/ws_grpc"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	conn    *grpc.ClientConn
	service ws_grpc.WebsocketServicesClient
}

func NewClient(cfg configs.Config, log *zap.Logger) *client {
	cl := &client{}
	url := cfg.WsService.Url
	var err error
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to connect to WebSocket service: %v", err))
	}
	cl = &client{
		conn:    conn,
		service: ws_grpc.NewWebsocketServicesClient(conn),
	}

	return cl
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) EmitRoomCreated(room string) error {
	_, err := c.service.EmitRoomCreated(context.Background(),
		&ws_grpc.EmitRoomCreatedRequest{RoomId: room})
	return err
}

func (c *client) EmitJoinRoom(room string, participantAddress string, sessionID string, sdpAnswer string, sdpOffer string) error {
	_, err := c.service.EmitJoinRoom(context.Background(),
		&ws_grpc.EmitJoinRoomReequest{
			RoomId:             room,
			ParticipantAddress: participantAddress,
			SessionId:          sessionID,
			SdpAnswer:          sdpAnswer,
			SdpOffer:           sdpOffer,
		})
	return err
}

func (c *client) EmitRequireRenegotiateSession(room string, participantAddrees string, sessionId string, sdpOffer string) error {
	_, err := c.service.EmitRequireRenegotiateSession(context.Background(),
		&ws_grpc.EmitRequireRenegotiateSessionRequest{
			RoomId:             room,
			ParticipantAddress: participantAddrees,
			SessionId:          sessionId,
			SdpOffer:           sdpOffer,
		})
	return err
}
