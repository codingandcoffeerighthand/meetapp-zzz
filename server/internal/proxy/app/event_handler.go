package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"proxy-srv/internal/proxy/domain"
	"proxy-srv/pkg/gencode/cloudflare_client"
	"proxy-srv/pkg/gencode/smc_gen"
)

/*
Join room flow
 1. New session
 2. Set session
 3. publish track
 4. emit event hasi new participant joined
 5. emit event Resp for client
*/
func (b *app) JoinRoomHandler(evt *smc_gen.MeeetingParticipantJoined) error {
	defer func() {
		if r := recover(); r != nil {
			b.errChan <- fmt.Errorf("join room %v", r)
		}
	}()
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

	go func() {
		err := b.smc.EmitEventToFrontend(evt.RoomId, evt.Participant.String(), domain.EventJoinedRoom{
			EventName: domain.EventJoinedRoomName,
			SdpAnswer: sdpAnswer,
		})
		if err != nil {
			b.errChan <- fmt.Errorf("emit event joined room %v", err)
		}
	}()
	go func() {
		err := b.EmitNewParticipantJoined(evt.RoomId, evt.Participant.String())
		if err != nil {
			b.errChan <- err
		}
	}()
	return nil
}

func (b *app) EventForwardBackendHandler(evt *smc_gen.MeeetingEventForwardedToBackend) error {
	evtData := domain.EventForwardedToBackend{}
	err := json.Unmarshal(evt.EventData, &evtData)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			b.errChan <- fmt.Errorf("event forwarded to backend %v", r)
		}
	}()
	switch evtData.EventName {
	case domain.EventLocalPeerConnectionSuccess:
		err := b.RoomPull(evt.RoomId)
		if err != nil {
			b.errChan <- fmt.Errorf("room pull %v", err)
		}
	case domain.EventRenegoiateSession:
		if evtData.Data == nil {
			return fmt.Errorf("invalid data")
		}
		data, ok := (*evtData.Data).(map[string]any)
		if !ok {
			return fmt.Errorf("invalid data")
		}
		sdpAnswer64 := data["sdp_answer"].(string)
		sdpAnswer, err := b.cryt.Decrypt(sdpAnswer64)
		if err != nil {
			return errors.New("decrypt sdp")
		}
		remote_session := data["remote_session"].(string)

		go func() {
			err := b.RenegatiateSession(context.Background(), remote_session, sdpAnswer)
			if err != nil {
				b.errChan <- fmt.Errorf("room pull %v", err)
			}
		}()

	default:
		return fmt.Errorf("unknown event name %s", evtData.EventName)
	}
	return nil
}
