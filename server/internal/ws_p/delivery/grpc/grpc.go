package ws_grpc_server

import (
	"context"
	"fmt"
	"net"
	ws_app "proxy-srv/internal/ws_p/app"
	"proxy-srv/internal/ws_p/config"
	"proxy-srv/pkg/gencode/ws_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AppInterface interface {
	CreateRoom(roomId string) error
	JoinRoom(req ws_app.JoinRoomReq) error
	RequireRenegotiateSession(req ws_app.RequireRenegotiateSessionReq) error
}

type server struct {
	ws_grpc.UnimplementedWebsocketServicesServer
	port int
	app  AppInterface
}

func NewGrpcServer(cfg config.ServerConfig, app AppInterface) *server {
	return &server{
		port: cfg.GrpcPort,
		app:  app,
	}
}

func (s *server) EmitJoinRoom(ctx context.Context, req *ws_grpc.EmitJoinRoomReequest) (*ws_grpc.EmitJoinRoomResponse, error) {
	err := s.app.JoinRoom(ws_app.JoinRoomReq{
		ParticipantAddr: req.ParticipantAddress,
		RoomId:          req.RoomId,
		SessionId:       req.SessionId,
		SdpAnswer:       req.SdpAnswer,
		SdpOffer:        req.SdpOffer,
	})
	return nil, err
}

func (s *server) EmitRequireRenegotiateSession(ctx context.Context, req *ws_grpc.EmitRequireRenegotiateSessionRequest) (*ws_grpc.EmitRequireRenegotiateSessionResponse, error) {
	err := s.app.RequireRenegotiateSession(ws_app.RequireRenegotiateSessionReq{
		ParticipantAddr: req.ParticipantAddress,
		RoomId:          req.RoomId,
		SessionId:       req.SessionId,
		SdpOffer:        req.SdpOffer,
	})
	return nil, err
}

func (s *server) EmitRoomCreated(ctx context.Context, req *ws_grpc.EmitRoomCreatedRequest) (*ws_grpc.EmitRoomCreatedResponse, error) {
	err := s.app.CreateRoom(req.RoomId)
	return &ws_grpc.EmitRoomCreatedResponse{}, err
}

func (s *server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	ws_grpc.RegisterWebsocketServicesServer(srv, s)
	reflection.Register(srv)
	fmt.Printf("### grpc running :%d\n", s.port)
	return srv.Serve(l)
}
