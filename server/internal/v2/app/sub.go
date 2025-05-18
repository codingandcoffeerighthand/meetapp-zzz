package app

import (
	"context"
	"fmt"
)

func (a *app) SubJoinRoom(ctx context.Context) error {
	c, sub, err := a.meet.SubJoinRoomEvent(ctx)
	if err != nil {
		return err
	}
	a.unsubscribe = func() {
		sub.Unsubscribe()
		a.unsubscribe()
	}
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("sub join room %v", r)
		}
	}()
	go func() {
		for evt := range c {
			err := a.JoinRoomHandler(ctx, evt)
			if err != nil {
				a.errChan <- err
			}
		}
	}()
	return nil
}

func (a *app) SubLeaveRoom(ctx context.Context) error {
	c, sub, err := a.meet.SubLeaveRoomEvent(ctx)
	if err != nil {
		return err
	}
	a.unsubscribe = func() {
		sub.Unsubscribe()
		a.unsubscribe()
	}
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("sub leave room %v", r)
		}
	}()
	go func() {
		for evt := range c {
			err := a.LeaveRoomHandler(ctx, evt)
			if err != nil {
				a.errChan <- err
			}
		}
	}()
	return nil
}

func (a *app) SubAddTracks(ctx context.Context) error {
	c, sub, err := a.meet.SubAddTracksEvent(ctx)
	if err != nil {
		return err
	}
	a.unsubscribe = func() {
		sub.Unsubscribe()
		a.unsubscribe()
	}
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("sub add tracks %v", r)
		}
	}()
	go func() {
		for evt := range c {
			err := a.AddTracksHandler(ctx, evt)
			if err != nil {
				a.errChan <- err
			}
		}
	}()
	return nil
}

func (a *app) SubRemoveTracks(ctx context.Context) error {
	c, sub, err := a.meet.SubRemoveTrack(ctx)
	if err != nil {
		return err
	}
	a.unsubscribe = func() {
		sub.Unsubscribe()
		a.unsubscribe()
	}
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("sub remove tracks %v", r)
		}
	}()
	go func() {
		for evt := range c {
			err := a.RemoveTracksHandler(ctx, evt)
			if err != nil {
				a.errChan <- err
			}
		}
	}()
	return nil
}

func (a *app) SubBackend(ctx context.Context) error {
	c, sub, err := a.meet.SubBackendEvent(ctx)
	if err != nil {
		return err
	}
	a.unsubscribe = func() {
		sub.Unsubscribe()
		a.unsubscribe()
	}
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("sub backend %v", r)
		}
	}()
	go func() {
		for evt := range c {
			err := a.BackendHandler(ctx, evt)
			if err != nil {
				a.errChan <- err
			}
		}
	}()
	return nil
}

func (a *app) Unsubscribe() {
	a.unsubscribe()
}
