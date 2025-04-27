package app

import (
	"context"
	"errors"
	"fmt"
	"proxy-srv/internal/proxy/domain"
	"proxy-srv/pkg/gencode/cloudflare_client"
	"proxy-srv/pkg/gencode/smc_gen"
	"sync"
)

func (b *app) PublisbTrack(
	sessionId string, sdpOffer string,
	tracks []cloudflare_client.TrackObject) (
	sdpAnswer string, err error) {
	resp, err := b.clf.AddLocalTrack(context.Background(), sessionId, sdpOffer, tracks)
	return *resp.SessionDescription.Sdp, err
}

var pullRoomMutex = sync.Mutex{}

func (b *app) RoomPull(roomId string) error {
	pullRoomMutex.Lock()
	defer pullRoomMutex.Unlock()
	smcResp, err := b.smc.GetParticipantsAndTracksOfRoom(roomId)
	if err != nil {
		return err
	}
	ps := smcResp.Arg0
	track := smcResp.Arg1
	defer func() {
		if r := recover(); r != nil {
			b.errChan <- fmt.Errorf("room pull %v", r)
		}
	}()
	for _, p := range ps {
		_, err := b.PullTrack(roomId, p.WalletAddress.String(), p.SessionID, ps, track)
		if err != nil {
			b.errChan <- err
		}
	}
	return nil
}

func (b *app) PullTrack(roomId string, pAddr string, sessionId string, ps []smc_gen.DAppMeetingParticipant, tracks []smc_gen.DAppMeetingTrack) (sdpOffer string, err error) {
	remoteSession, err := b.clf.NewSession(context.Background())
	if err != nil {
		return
	}
	resp, err := b.clf.GetStatusSession(sessionId)
	if err != nil {
		return
	}
	var remoteTracks = make([]cloudflare_client.TrackObject, 0)
	var remoteLocalTrackValue = cloudflare_client.TrackObjectLocationRemote
	for _, rTrack := range tracks {
		if rTrack.SessionId == sessionId || !rTrack.IsPublished {
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
			return "", fmt.Errorf("pullTrack err %v", err)
		}
		if !*resp.RequiresImmediateRenegotiation {
			return "", fmt.Errorf("pullTrack err not sdp answer %v", err)
		}
		if resp.SessionDescription == nil {
			return "", errors.New("no remote track")
		}
		defer func() {
			if r := recover(); r != nil {
				b.log.Error((r.(error)).Error())
			}
		}()
		go func(roomId string, addr string) {
			err := b.smc.EmitEventToFrontend(roomId, addr, domain.EventPullTrack{
				EventName:     domain.EventPullTrackName,
				SdpOffer:      *resp.SessionDescription.Sdp,
				RemoteSession: remoteSession,
			})
			if err != nil {
				b.errChan <- fmt.Errorf("emit pulltrack %v", err)
			}
		}(roomId, pAddr)
		return *resp.SessionDescription.Sdp, err
	}
	err = b.smc.EmitEventToFrontend(roomId, pAddr, domain.EventPullTrack{
		EventName:     domain.EventPullTrackName,
		SdpOffer:      "",
		RemoteSession: remoteSession,
	})
	return "", nil
}

func (b *app) RenegatiateSession(ctx context.Context, session string, sdpAnswer string, roomId string, addr string) error {
	_, err := b.clf.RenegatiateSession(ctx, session, sdpAnswer)
	defer func() {
		if r := recover(); r != nil {
			b.errChan <- (r.(error))
		}
	}()
	go func() {
		err := b.EmmitRenegoiateSucces(session, roomId, addr)
		if err != nil {
			b.errChan <- err
		}
	}()
	return err
}

func (b *app) EmmitRenegoiateSucces(sessionId string,
	roomId string, addr string) error {
	resp, err := b.clf.GetStatusSession(sessionId)
	if err != nil {
		return err
	}
	if resp == nil || resp.Tracks == nil {
		return fmt.Errorf("status session %s nil", sessionId)
	}
	tracks := make([]domain.Track, len(*resp.Tracks))
	for i, v := range *resp.Tracks {
		tracks[i] = domain.Track{
			Mid:       *v.Mid,
			SessionId: *v.SessionId,
			TrackName: *v.TrackName,
		}
	}
	evt := domain.EventRenegoiateSuccess{
		EventName: domain.EventRenegoiateSuccessName,
		Tracks:    tracks,
	}
	b.smc.EmitEventToFrontend(roomId, addr, evt)
	return nil
}
