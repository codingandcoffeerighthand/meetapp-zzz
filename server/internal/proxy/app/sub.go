package app

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/event"
)

func (b *app) HandlerJoinRoomEvent(ctx context.Context) (event.Subscription, error) {
	joinRoomChan, joinRoomSub, err := b.smc.SubJoinRoomEvent(ctx)
	if err != nil {
		return nil, err
	}
	handleEventSub(ctx, joinRoomChan, b.JoinRoomHandler, b.errChan)
	return joinRoomSub, err
}

func (b *app) HandlerEventForwardedToBackend(ctx context.Context) (event.Subscription, error) {
	eventForwardedToBackendChan, eventForwardedToBackendSub, err := b.smc.SubEventToBackend(ctx)
	if err != nil {
		return nil, err
	}
	handleEventSub(ctx, eventForwardedToBackendChan, b.EventForwardBackendHandler, b.errChan)
	return eventForwardedToBackendSub, err
}

func (b *app) HandlerEventAddedTracks(ctx context.Context) (event.Subscription, error) {
	eventAddedTracksChan, eventAddedTracksSub, err := b.smc.SubAddTracks(ctx)
	if err != nil {
		return nil, err
	}
	handleEventSub(ctx, eventAddedTracksChan, b.EventAddedTracksHandler, b.errChan)
	return eventAddedTracksSub, err
}
func (b *app) HandlerEventRemovedTrack(ctx context.Context) (event.Subscription, error) {
	eventRemovedTrackChan, eventRemovedTrackSub, err := b.smc.SubRemoveTrack(ctx)
	if err != nil {
		return nil, err
	}
	handleEventSub(ctx, eventRemovedTrackChan, b.EventRemovedTrackHandler, b.errChan)
	return eventRemovedTrackSub, err
}
func (b *app) HandlerLeaveRoom(ctx context.Context) (event.Subscription, error) {
	eventLeaveeRoomChan, eventLeaveeRoomSub, err := b.smc.SubLeaveRoom(ctx)
	if err != nil {
		return nil, err
	}
	handleEventSub(ctx, eventLeaveeRoomChan, b.EventLeaveRoom, b.errChan)
	return eventLeaveeRoomSub, err
}

// run
func (b *app) runSubJoinRoom(ctx context.Context, cleanUp *func()) {
	go func() {
		joinRoomSub, err := b.HandlerJoinRoomEvent(ctx)
		if err != nil {
			b.errChan <- fmt.Errorf("sub join room %v", err)
		}
		cl := func() {
			(*cleanUp)()
			joinRoomSub.Unsubscribe()
		}
		cleanUp = &cl
	}()

}
