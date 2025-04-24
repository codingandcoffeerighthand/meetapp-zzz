package biz

import (
	"context"
	"proxy-srv/pkg/gencode/cloudflare_client"
	"proxy-srv/pkg/gencode/smc_gen"
)

func (b *biz) CreateRoomHandler(evt *smc_gen.MeeetingRoomCreated) error {
	b.grcpCl.EmitRoomCreated(evt.RoomId)
	return nil
}

func (b *biz) JoinRoomHandler(evt *smc_gen.MeeetingParticipantJoined) error {
	ctx := context.Background()

	//1. New session
	sessionId, err := b.clf.NewSession(ctx)
	if err != nil {
		return err
	}

	// 2.1 set session
	err = b.smc.SetParticipantSessionID(ctx, evt.RoomId, evt.Participant.String(), sessionId)
	if err != nil {
		return err
	}
	// 2.2 add local track
	sdpO, err := b.cryt.Decrypt(evt.SdpOffer)
	if err != nil {
		return err
	}
	tracks := make([]cloudflare_client.TrackObject, len(evt.InitialTracks))
	for i := range evt.InitialTracks {
		tracks[i] = cloudflare_client.TrackObject{
			Mid:       &evt.InitialTracks[i].Mid,
			SessionId: &sessionId,
			TrackName: &evt.InitialTracks[i].TrackName,
		}
	}
	sdpAnswer, err := b.PublisbTrack(sessionId, sdpO, tracks)
	if err != nil {
		return err
	}

	err = b.grcpCl.EmitJoinRoom(evt.RoomId, evt.Participant.String(), sessionId, sdpAnswer, "")
	if err != nil {
		return err
	}
	return nil
}
