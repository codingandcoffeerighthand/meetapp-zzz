package app

import (
	"context"
	"proxy-srv/internal/v2/domain"
)

// CloudFlareInfraV2 là interface cho infra cloudflare v2
// Đảm bảo các method và kiểu dữ liệu phù hợp với client.go
// server/internal/v2/infra/cloudflare/client.go
type CloudFlareInfraV2 interface {
	NewSession(ctx context.Context) (string, error)
	AddLocalTrack(ctx context.Context, sessionId string, sdpOffer string, tracks []domain.Track) (string, error)
	AddRemoteTrack(ctx context.Context, sessionId string, tracks []domain.Track) (string, error)
	RenegatiateSession(ctx context.Context, session string, sdpAnswer string) (string, error)
	GetStatusSession(sessionId string) ([]domain.Track, error)
}

type MeetInfra interface {
	SetNewSession(ctx context.Context, roomID string, oldSessionID string, sessionID string) error
	EmitFrontEndEvent(ctx context.Context, roomID string, sessionID string, eventType string, data []byte) error
	GetRoomInfo(ctx context.Context, roomID string) (domain.Room, error)
	SubJoinRoomEvent(ctx context.Context) (chan *domain.JoinRoomEvent, domain.Subscription, error)
	SubLeaveRoomEvent(ctx context.Context) (chan *domain.LeaveRoomEvent, domain.Subscription, error)
	SubAddTracksEvent(ctx context.Context) (chan *domain.AddTracksEvent, domain.Subscription, error)
	SubRemoveTrack(ctx context.Context) (chan *domain.RemoveTracksEvent, domain.Subscription, error)
	SubBackendEvent(ctx context.Context) (chan *domain.BackendEvent, domain.Subscription, error)
}
