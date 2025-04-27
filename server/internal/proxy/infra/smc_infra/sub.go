package smc_infra

import (
	"context"
	"proxy-srv/pkg/gencode/smc_gen"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/event"
)

func (s *smcInfra) SubJoinRoomEvent(ctx context.Context) (<-chan *smc_gen.MeeetingParticipantJoined, event.Subscription, error) {
	sink := make(chan *smc_gen.MeeetingParticipantJoined)
	eventSub, err := bind.WatchEvents(s.BoundContract, &bind.WatchOpts{Context: ctx, Start: block}, s.contract.UnpackParticipantJoinedEvent, sink)
	return sink, eventSub, err
}

func (s *smcInfra) SubCreateRoomEvent(ctx context.Context) (<-chan *smc_gen.MeeetingRoomCreated, event.Subscription, error) {
	sink := make(chan *smc_gen.MeeetingRoomCreated)
	eventSub, err := bind.WatchEvents(s.BoundContract, &bind.WatchOpts{Context: ctx, Start: block}, s.contract.UnpackRoomCreatedEvent, sink)
	return sink, eventSub, err
}

func (s *smcInfra) SubEventToBackend(ctx context.Context) (<-chan *smc_gen.MeeetingEventForwardedToBackend, event.Subscription, error) {
	sink := make(chan *smc_gen.MeeetingEventForwardedToBackend)
	eventSub, err := bind.WatchEvents(s.BoundContract, &bind.WatchOpts{Context: ctx, Start: block}, s.contract.UnpackEventForwardedToBackendEvent, sink)
	return sink, eventSub, err
}

func (s *smcInfra) SubAddTracks(ctx context.Context) (<-chan *smc_gen.MeeetingTrackAdded, event.Subscription, error) {
	sink := make(chan *smc_gen.MeeetingTrackAdded)
	eventSub, err := bind.WatchEvents(s.BoundContract, &bind.WatchOpts{Context: ctx, Start: block}, s.contract.UnpackTrackAddedEvent, sink)
	return sink, eventSub, err
}

func (s *smcInfra) SubRemoveTrack(ctx context.Context) (<-chan *smc_gen.MeeetingRemoveTracks, event.Subscription, error) {
	sink := make(chan *smc_gen.MeeetingRemoveTracks)
	eventSub, err := bind.WatchEvents(s.BoundContract, &bind.WatchOpts{Context: ctx, Start: block}, s.contract.UnpackRemoveTracksEvent, sink)
	return sink, eventSub, err
}
