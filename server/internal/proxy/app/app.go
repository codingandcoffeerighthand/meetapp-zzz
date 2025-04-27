package app

import (
	"context"
	"fmt"
	"proxy-srv/pkg/gencode/cloudflare_client"
	"proxy-srv/pkg/gencode/smc_gen"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/event"
	"go.uber.org/zap"
)

type CloudFlareInfra interface {
	NewSession(ctx context.Context) (string, error)
	AddLocalTrack(ctx context.Context, sessionId string, sdpOffer string, tracks []cloudflare_client.TrackObject) (cloudflare_client.TracksResponse, error)
	AddRemoteTrack(ctx context.Context, sessionId string, tracks []cloudflare_client.TrackObject) (cloudflare_client.TracksResponse, error)
	RenegatiateSession(ctx context.Context, session string, sdpAnswer string) (cloudflare_client.SessionDescription, error)
	GetStatusSession(sessionId string) (*cloudflare_client.GetSessionStateResponse, error)
}

type SMCInfra interface {
	SubCreateRoomEvent(ctx context.Context) (<-chan *smc_gen.MeeetingRoomCreated, event.Subscription, error)
	SetParticipantSessionID(ctx context.Context, room_id string, addr string, session_id string) error
	SubJoinRoomEvent(ctx context.Context) (<-chan *smc_gen.MeeetingParticipantJoined, event.Subscription, error)
	GetParticipantsAndTracksOfRoom(roomId string) (smc_gen.GetParticipantOfRoomOutput, error)
	EmitEventToFrontend(rooId string, addrStr string, data any) error
	SubEventToBackend(ctx context.Context) (<-chan *smc_gen.MeeetingEventForwardedToBackend, event.Subscription, error)
	SubAddTracks(ctx context.Context) (<-chan *smc_gen.MeeetingTrackAdded, event.Subscription, error)
	SubRemoveTrack(ctx context.Context) (<-chan *smc_gen.MeeetingRemoveTracks, event.Subscription, error)
}
type CryptionService interface {
	// Encrypt(data []byte) ([]byte, error)
	Decrypt(data string) (string, error)
}

type app struct {
	clf     CloudFlareInfra
	smc     SMCInfra
	cryt    CryptionService
	log     *zap.Logger
	errChan chan error
	done    chan any
}

func NewApp(clf CloudFlareInfra, smc SMCInfra, crypt CryptionService, log *zap.Logger) (*app, func(), error) {

	b := &app{
		clf:     clf,
		smc:     smc,
		cryt:    crypt,
		log:     log,
		errChan: make(chan error),
		done:    make(chan any),
	}

	cl, err := b.Run(context.Background())
	return b, cl, err
}

type EventHandler[Ev bind.ContractEvent] func(*Ev) error

func handleEventSub[Ev bind.ContractEvent](ctx context.Context, sink <-chan *Ev, handler EventHandler[Ev], errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-sink:
			if err := handler(ev); err != nil {
				errChan <- err
			}
		}
	}
}

func (b *app) Run(ctx context.Context) (func(), error) {
	cleanUp := func() {
		close(b.errChan)
		close(b.done)
	}
	defer func() {
		if r := recover(); r != nil {
			b.log.Error((r.(error)).Error())
		}
	}()

	b.runSubJoinRoom(ctx, &cleanUp)
	b.runSub(ctx, &cleanUp, b.HandlerEventForwardedToBackend)
	b.runSub(ctx, &cleanUp, b.HandlerEventAddedTracks)
	b.runSub(ctx, &cleanUp, b.HandlerEventRemovedTrack)

	go func() {
		for err := range b.errChan {
			if err != nil {
				b.log.Error(err.Error())
			}
		}
	}()

	return cleanUp, nil
}

func (a *app) Stop() {
	a.done <- struct{}{}
}

func (a *app) Done() {
	<-a.done
}

type SubEvtHandler func(context.Context) (event.Subscription, error)

func (b *app) runSub(ctx context.Context, cleanUp *func(), hdl SubEvtHandler) {
	go func() {
		sub, err := hdl(ctx)
		if err != nil {
			b.errChan <- fmt.Errorf("sub %v", err)
		}
		cl := func() {
			(*cleanUp)()
			sub.Unsubscribe()
		}
		cleanUp = &cl
	}()
}
