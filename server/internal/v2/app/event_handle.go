package app

import (
	"context"
	"encoding/json"
	"fmt"
	"proxy-srv/internal/v2/domain"
)

func (a *app) JoinRoomHandler(ctx context.Context, evt *domain.JoinRoomEvent) error {
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("join room %v", r)
		}
	}()

	//1. New session
	newSession, err := a.clf.NewSession(ctx)
	if err != nil {
		return err
	}

	// 2.1 set session
	err = a.meet.SetNewSession(
		ctx,
		evt.RoomID,
		evt.SessionID,
		newSession,
	)
	if err != nil {
		return err
	}

	tracks := make([]domain.Track, len(evt.Tracks))
	for i := range evt.Tracks {
		tracks[i] = domain.Track{
			Mid:       evt.Tracks[i].Mid,
			TrackName: evt.Tracks[i].TrackName,
			SessionID: evt.Tracks[i].SessionID,
			Location:  evt.Tracks[i].Location,
		}
	}
	sdpAnswer, err := a.clf.AddLocalTrack(ctx, newSession, evt.SdpOffer, tracks)
	if err != nil {
		return err
	}
	go func() {
		err := a.meet.EmitFrontEndEvent(ctx, evt.RoomID, newSession, domain.EventJoinedRoomName, []byte(sdpAnswer))
		if err != nil {
			a.errChan <- err
		}
	}()
	return nil
}

/*
	after local connected, pull tracks from meet
*/

func (a *app) PullTracksForRoom(ctx context.Context, roomID string) error {
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("pull tracks for participant %v", r)
		}
	}()
	roomInfo, err := a.meet.GetRoomInfo(ctx, roomID)
	if err != nil {
		return err
	}

	tracksOfRoom := roomInfo.GetAllTracks()
	for _, participant := range roomInfo.Participants {
		err := a.PullTracksForParticipant(ctx, roomID, participant.SessionID, tracksOfRoom)
		if err != nil {
			a.errChan <- err
		}
	}
	return nil
}

func (a *app) PullTracksForParticipant(ctx context.Context,
	roomID string, sessionID string,
	tracksOfRoom []domain.Track) error {
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("pull tracks for participant %v", r)
		}
	}()

	remoteSession, err := a.clf.NewSession(ctx)
	if err != nil {
		return err
	}

	tracksOfParticipant := make([]domain.Track, 0)
	for _, track := range tracksOfRoom {
		if track.SessionID != sessionID {
			tracksOfParticipant = append(tracksOfParticipant, track)
		}
	}

	sdpoffer, err := a.clf.AddRemoteTrack(ctx, remoteSession, tracksOfParticipant)
	if err != nil {
		return err
	}

	go func() {
		evtRemoteConnect := domain.RemoteInfoEvent{
			RoomID:          roomID,
			RemoteSessionID: remoteSession,
			SdpOffer:        sdpoffer,
		}
		evtJson, _ := json.Marshal(evtRemoteConnect)
		err := a.meet.EmitFrontEndEvent(ctx, roomID, sessionID, domain.RemoteInfoEventName, evtJson)
		if err != nil {
			a.errChan <- err
		}
	}()
	return nil
}

func (s *app) AddTracksHandler(ctx context.Context, evt *domain.AddTracksEvent) error {
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("add tracks %v", r)
		}
	}()

	sdpAnswer, err := s.clf.AddLocalTrack(ctx, evt.SessionID, evt.SdpOffer, evt.Tracks)
	if err != nil {
		return err
	}

	go func() {
		err := s.meet.EmitFrontEndEvent(ctx, evt.RoomID, evt.SessionID, domain.EventAddTracksName, []byte(sdpAnswer))
		if err != nil {
			s.errChan <- err
		}
	}()
	return nil
}

func (s *app) RemoveTracksHandler(ctx context.Context, evt *domain.RemoveTracksEvent) error {
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("remove tracks %v", r)
		}
	}()
	err := s.PullTracksForRoom(ctx, evt.RoomID)
	if err != nil {
		return err
	}
	return nil
}

func (s *app) LocalConnectedHandler(ctx context.Context, evt *domain.LocalConnectedEvent) error {
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("local connected %v", r)
		}
	}()
	err := s.PullTracksForRoom(ctx, evt.RoomID)
	if err != nil {
		return err
	}
	return nil
}

func (s *app) RemoteConnectedHandler(ctx context.Context, evt *domain.RemoteConnectedEvent) error {
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("remote connected %v", r)
		}
	}()
	_, err := s.clf.RenegatiateSession(ctx, evt.RemoteSessionID, evt.SdpAnswer)
	if err != nil {
		return err
	}
	go func() {
		tracks, err := s.clf.GetStatusSession(evt.RemoteSessionID)
		if err != nil {
			s.errChan <- err
			return
		}
		room, err := s.meet.GetRoomInfo(ctx, evt.RoomID)
		if err != nil {
			s.errChan <- err
			return
		}
		evtDomain := domain.RemoteConectSuccessEvent{
			RoomID: room,
			Tracks: tracks,
		}
		evtJson, _ := json.Marshal(evtDomain)
		err = s.meet.EmitFrontEndEvent(ctx, evt.RoomID, evt.RemoteSessionID, domain.EventRemoteConectSuccessName, evtJson)
		if err != nil {
			s.errChan <- err
		}
	}()
	return nil
}

func (s *app) LeaveRoomHandler(ctx context.Context, evt *domain.LeaveRoomEvent) error {
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("leave room %v", r)
		}
	}()
	err := s.PullTracksForRoom(ctx, evt.RoomID)
	return err
}

func (s *app) BackendHandler(ctx context.Context, evt *domain.BackendEvent) error {
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("backend %v", r)
		}
	}()
	switch evt.EventType {
	case domain.EventLocalConnectedName:
		localConnectEvent := domain.LocalConnectedEvent{}
		err := json.Unmarshal(evt.Data, &localConnectEvent)
		if err != nil {
			s.errChan <- err
			return err
		}
		return s.LocalConnectedHandler(ctx, &localConnectEvent)
	case domain.EventRemoteConnectedName:
		remoteConnectEvent := domain.RemoteConnectedEvent{}
		err := json.Unmarshal(evt.Data, &remoteConnectEvent)
		if err != nil {
			s.errChan <- err
			return err
		}
		return s.RemoteConnectedHandler(ctx, &remoteConnectEvent)
	default:
		err := fmt.Errorf("unknown event type %s", evt.EventType)
		s.errChan <- err
		return err
	}
}
