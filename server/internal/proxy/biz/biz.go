package biz

import (
	"context"
	"errors"
	"fmt"
	"proxy-srv/pkg/gencode/cloudflare_client"
	"proxy-srv/pkg/gencode/smc_gen"
	"slices"

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
	CheckAuthorized(addressStr string) (bool, error)
	SubCreateRoomEvent(ctx context.Context) (<-chan *smc_gen.MeeetingRoomCreated, event.Subscription, error)
	SetParticipantSessionID(ctx context.Context, room_id string, addr string, session_id string) error
	SubJoinRoomEvent(ctx context.Context) (<-chan *smc_gen.MeeetingParticipantJoined, event.Subscription, error)
	GetParticipantsAndTracksOfRoom(roomId string) (smc_gen.GetParticipantOfRoomOutput, error)
}
type CryptionService interface {
	// Encrypt(data []byte) ([]byte, error)
	Decrypt(data string) (string, error)
}

type GrpcClient interface {
	EmitRoomCreated(room string) error
	EmitJoinRoom(room string, participantAddress string, sessionID string, sdpAnswer string, sdpOffer string) error
	EmitRequireRenegotiateSession(room string, participantAddrees string, sessionId string, sdpOffer string) error
}

type biz struct {
	clf     CloudFlareInfra
	smc     SMCInfra
	grcpCl  GrpcClient
	cryt    CryptionService
	log     *zap.Logger
	errChan chan error
}

func NewBiz(clf CloudFlareInfra, smc SMCInfra, grpcCl GrpcClient, crypt CryptionService, log *zap.Logger) (*biz, func(), error) {

	b := &biz{
		clf:     clf,
		smc:     smc,
		grcpCl:  grpcCl,
		cryt:    crypt,
		log:     log,
		errChan: make(chan error),
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

func (b *biz) HandlerSubCreateRoomEvent(ctx context.Context) (event.Subscription, error) {
	createRoomChan, creatRoomSub, err := b.smc.SubCreateRoomEvent(ctx)
	if err != nil {
		return nil, err
	}
	handleEventSub(ctx, createRoomChan, b.CreateRoomHandler, b.errChan)
	return creatRoomSub, err
}

func (b *biz) HandlerJoinRoomEvent(ctx context.Context) (event.Subscription, error) {
	joinRoomChan, joinRoomSub, err := b.smc.SubJoinRoomEvent(ctx)
	if err != nil {
		return nil, err
	}
	handleEventSub(ctx, joinRoomChan, b.JoinRoomHandler, b.errChan)
	return joinRoomSub, err
}
func (b *biz) runSubCreatedRoom(ctx context.Context, cleanUp *func()) {
	go func() {
		createRoomSub, err := b.HandlerSubCreateRoomEvent(ctx)
		if err != nil {
			b.errChan <- fmt.Errorf("sub created room %v", err)
		}
		cl := func() {
			(*cleanUp)()
			createRoomSub.Unsubscribe()
		}
		cleanUp = &cl
	}()
}

func (b *biz) runSubJoinRoom(ctx context.Context, cleanUp *func()) {
	go func() {
		joinRoomSub, err := b.HandlerJoinRoomEvent(ctx)
		if err != nil {
			b.errChan <- fmt.Errorf("sub join room %v", err)
			b.errChan <- err
		}
		cl := func() {
			(*cleanUp)()
			joinRoomSub.Unsubscribe()
		}
		cleanUp = &cl
	}()

}
func (b *biz) Run(ctx context.Context) (func(), error) {
	cleanUp := func() {
		close(b.errChan)
	}
	defer func() {
		if r := recover(); r != nil {
			b.log.Error((r.(error)).Error())
		}
	}()

	b.runSubCreatedRoom(ctx, &cleanUp)
	b.runSubJoinRoom(ctx, &cleanUp)

	go func() {
		for err := range b.errChan {
			b.log.Debug(err.Error())
			b.log.Error(err.Error())
		}
	}()

	return cleanUp, nil
}

func (b *biz) RenegatiateSession(ctx context.Context, session string, sdpAnswer string) error {
	_, err := b.clf.RenegatiateSession(ctx, session, sdpAnswer)
	if err != nil {
		b.log.Error(err.Error())
	}
	return err
}

func (b *biz) PublisbTrack(
	sessionId string, sdpOffer string,
	tracks []cloudflare_client.TrackObject) (
	sdpAnswer string, err error) {
	resp, err := b.clf.AddLocalTrack(context.Background(), sessionId, sdpOffer, tracks)
	return *resp.SessionDescription.Sdp, err
}

func (b *biz) RoomPull(roomId string) error {
	smcResp, err := b.smc.GetParticipantsAndTracksOfRoom(roomId)
	if err != nil {
		return err
	}
	ps := smcResp.Arg0
	track := smcResp.Arg1
	for _, p := range ps {
		go b.PullTrack(roomId, p.SessionID, ps, track)
	}
	return nil
}

func (b *biz) PullTrack(roomId string, sessionId string, ps []smc_gen.DAppMeetingParticipant, tracks []smc_gen.DAppMeetingTrack) (sdpOffer string, err error) {
	remoteSession, err := b.clf.NewSession(context.Background())
	if err != nil {
		return
	}
	idx := slices.IndexFunc(ps, func(p smc_gen.DAppMeetingParticipant) bool {
		return p.SessionID == sessionId
	})
	if idx == -1 {
		return "", errors.New("participant not found")
	}
	resp, err := b.clf.GetStatusSession(sessionId)
	if err != nil {
		return
	}
	var remoteTracks = make([]cloudflare_client.TrackObject, 0)
	var remoteLocalTrackValue = cloudflare_client.TrackObjectLocationRemote
	for _, rTrack := range tracks {
		if rTrack.SessionId == sessionId {
			continue
		}
		exist := false
		for _, sTrack := range *resp.Tracks {
			if sTrack.SessionId == nil {
				continue
			}
			if rTrack.Mid == *sTrack.Mid &&
				rTrack.SessionId != *sTrack.SessionId &&
				rTrack.TrackName != *sTrack.TrackName {
				exist = true
				break
			}
		}
		if !exist {
			remoteTracks = append(remoteTracks, cloudflare_client.TrackObject{
				Mid:       &rTrack.Mid,
				SessionId: &rTrack.SessionId,
				TrackName: &rTrack.TrackName,
				Location:  &remoteLocalTrackValue,
			})
		}
	}

	if len(remoteTracks) > 0 {
		resp, err := b.clf.AddRemoteTrack(context.Background(), remoteSession, remoteTracks)
		if err != nil {
			return "", err
		}
		if !*resp.RequiresImmediateRenegotiation {
			return "", nil
		}
		defer func() {
			if r := recover(); r != nil {
				b.log.Error((r.(error)).Error())
			}
		}()
		go func(roomId string, pAddr string, sessionId string, sdpOffer string) {
			err := b.grcpCl.EmitRequireRenegotiateSession(
				roomId, pAddr, sessionId, sdpOffer,
			)
			if err != nil {
				b.log.Error(err.Error())
			}
		}(roomId, ps[idx].WalletAddress.String(), remoteSession, *resp.SessionDescription.Sdp)

		return *resp.SessionDescription.Sdp, err
	}
	return "", fmt.Errorf("no remote track")
}
