package infra_smc

import (
	"context"
	"fmt"
	"proxy-srv/internal/v2/domain"
	"proxy-srv/pkg/gencode/smc_gen/meet_smc"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/event"
)

type Sub struct {
	event.Subscription
	cleanUp func()
}

func NewSub(sub event.Subscription, cleanUp func()) *Sub {
	return &Sub{
		Subscription: sub,
		cleanUp:      cleanUp,
	}
}

func (s *Sub) Err() <-chan error {
	return s.Subscription.Err()
}

func (s *Sub) Unsubscribe() {
	s.Subscription.Unsubscribe()
	s.cleanUp()
}

func (s *smcInfra) SubJoinRoomEvent(ctx context.Context) (
	<-chan *domain.JoinRoomEvent, domain.Subscription, error) {
	sink := make(chan *meet_smc.MeetJoinRoomEvent)
	eventSub, err := bind.WatchEvents(
		s.BoundContract,
		&bind.WatchOpts{Context: ctx},
		s.contract.UnpackJoinRoomEventEvent,
		sink,
	)
	sinkDomain := make(chan *domain.JoinRoomEvent)
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("sub join room event %v", r)
		}
	}()
	go func() {
		for evt := range sink {
			evtDomain := &domain.JoinRoomEvent{
				RoomID:    evt.RoomId.Hex(),
				SessionID: evt.SessionId.Hex(),
				Tracks:    make([]domain.Track, len(evt.Tracks)),
				SdpOffer:  evt.SdpOffer,
			}
			for i, track := range evt.Tracks {
				evtDomain.Tracks[i] = domain.Track{
					Mid:       track.Mid,
					TrackName: track.TrackName,
					SessionID: track.SessionId,
					Location:  track.Location,
				}
			}
			sinkDomain <- evtDomain
		}
	}()

	return sinkDomain, NewSub(eventSub, func() {
		close(sinkDomain)
	}), err
}

func (s *smcInfra) SubAddTracksEvent(ctx context.Context) (
	<-chan *domain.AddTracksEvent, domain.Subscription, error) {
	sink := make(chan *meet_smc.MeetAddTracksEvent)
	eventSub, err := bind.WatchEvents(
		s.BoundContract,
		&bind.WatchOpts{Context: ctx},
		s.contract.UnpackAddTracksEventEvent,
		sink,
	)
	sinkDomain := make(chan *domain.AddTracksEvent)
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("sub add tracks event %v", r)
		}
	}()
	go func() {
		for evt := range sink {
			evtDomain := &domain.AddTracksEvent{
				RoomID:    evt.RoomId.Hex(),
				SessionID: evt.SessionId.Hex(),
				Tracks:    make([]domain.Track, len(evt.Tracks)),
				SdpOffer:  evt.SdpOffer,
			}
			for i, track := range evt.Tracks {
				evtDomain.Tracks[i] = domain.Track{
					Mid:       track.Mid,
					TrackName: track.TrackName,
					SessionID: track.SessionId,
					Location:  track.Location,
				}
			}
			sinkDomain <- evtDomain
		}
	}()

	return sinkDomain, NewSub(eventSub, func() {
		close(sinkDomain)
	}), err
}

func (s *smcInfra) SubRemoveTrack(ctx context.Context) (
	<-chan *domain.RemoveTracksEvent, domain.Subscription, error) {
	sink := make(chan *meet_smc.MeetRemoveTracksEvent)
	evtSub, err := bind.WatchEvents(
		s.BoundContract,
		&bind.WatchOpts{Context: ctx},
		s.contract.UnpackRemoveTracksEventEvent,
		sink)
	if err != nil {
		return nil, NewSub(evtSub, func() {}), err
	}
	sinkDomain := make(chan *domain.RemoveTracksEvent)
	defer func() {
		if e := recover(); e != nil {
			s.errChan <- fmt.Errorf("sub remove tracks event %v", e)
		}
	}()

	go func() {
		for evt := range sink {
			evtDomain := &domain.RemoveTracksEvent{
				RoomID:    evt.RoomId.Hex(),
				SessionID: evt.SessionId.Hex(),
				SdpOffer:  evt.SdpOffer,
			}
			sinkDomain <- evtDomain
		}
	}()
	return sinkDomain, NewSub(evtSub, func() { close(sinkDomain) }), err
}

func (s *smcInfra) SubLeaveRoomEvent(ctx context.Context) (
	<-chan *domain.LeaveRoomEvent, domain.Subscription, error) {
	sink := make(chan *meet_smc.MeetLeftRoomEvent)
	evtSub, err := bind.WatchEvents(
		s.BoundContract,
		&bind.WatchOpts{Context: ctx},
		s.contract.UnpackLeftRoomEventEvent,
		sink)
	if err != nil {
		return nil, NewSub(evtSub, func() {}), err
	}
	sinkDomain := make(chan *domain.LeaveRoomEvent)
	defer func() {
		if e := recover(); e != nil {
			s.errChan <- fmt.Errorf("sub leave room event %v", e)
		}
	}()
	go func() {
		for evt := range sink {
			evtDomain := &domain.LeaveRoomEvent{
				RoomID:    evt.RoomId.Hex(),
				SessionID: evt.SessionId.Hex(),
			}
			sinkDomain <- evtDomain
		}
	}()
	return sinkDomain, NewSub(evtSub, func() { close(sinkDomain) }), err
}

func (s *smcInfra) SubBackendEvent(ctx context.Context) (
	<-chan *domain.BackendEvent, domain.Subscription, error) {
	sink := make(chan *meet_smc.MeetBackendEvent)
	evtSub, err := bind.WatchEvents(
		s.BoundContract,
		&bind.WatchOpts{Context: ctx},
		s.contract.UnpackBackendEventEvent,
		sink)
	if err != nil {
		return nil, NewSub(evtSub, func() {}), err
	}
	sinkDomain := make(chan *domain.BackendEvent)
	defer func() {
		if e := recover(); e != nil {
			s.errChan <- fmt.Errorf("sub backend event %v", e)
		}
	}()
	go func() {
		for evt := range sink {
			evtDomain := &domain.BackendEvent{
				RoomID:    evt.RoomId.Hex(),
				SessionID: evt.SessionId.Hex(),
				EventType: evt.EventType,
				Data:      evt.Data,
			}
			sinkDomain <- evtDomain

		}
	}()
	return sinkDomain, NewSub(evtSub, func() { close(sinkDomain) }), err
}
