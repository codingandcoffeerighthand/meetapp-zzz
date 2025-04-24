package grpc

import (
	"context"
	"fmt"
	"net"
	"proxy-srv/internal/proxy/configs"
	"proxy-srv/pkg/gencode/proxy_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Biz interface {
	Run(ctx context.Context) (func(), error)
	RenegatiateSession(ctx context.Context, session string, sdpAnswer string) error
	RoomPull(roomId string) error
}
type server struct {
	proxy_grpc.UnimplementedProxyServiceServer
	biz  Biz
	port int
}

func NewServerGrpc(biz Biz, cfg configs.ServerConfig) *server {
	return &server{
		biz:  biz,
		port: cfg.Port,
	}
}

func (s *server) PUTRenegotiateSession(ctx context.Context, req *proxy_grpc.PUTRenegotiateSessionRequest) (*proxy_grpc.PUTRenegotiateSessionResponse, error) {
	err := s.biz.RenegatiateSession(ctx, req.SessionId, req.SdpAnswer)
	if err != nil {
		return nil, err
	}
	return &proxy_grpc.PUTRenegotiateSessionResponse{}, nil
}

func (s *server) EmitParicipantJoinedToRom(ctx context.Context, req *proxy_grpc.EmitParicipantJoinedToRomRequest) (*proxy_grpc.EmitParicipantJoinedToRomResponse, error) {
	err := s.biz.RoomPull(req.RoomId)
	if err != nil {
		return nil, err
	}
	return &proxy_grpc.EmitParicipantJoinedToRomResponse{}, nil
}

func (s *server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	proxy_grpc.RegisterProxyServiceServer(srv, s)
	reflection.Register(srv)
	return srv.Serve(l)
}
