package app

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type App struct {
	clf         CloudFlareInfraV2
	meet        MeetInfra
	log         *zap.Logger
	errChan     chan error
	unsubscribe func()
	done        chan any
	mu          sync.Mutex
}

func NewApp(cloudFlareInfa CloudFlareInfraV2, meetInfra MeetInfra, log *zap.Logger) *App {
	return &App{
		clf:     cloudFlareInfa,
		meet:    meetInfra,
		log:     log,
		errChan: make(chan error, 20),
		done:    make(chan any),
	}
}

func (a *App) Err() <-chan error {
	return a.errChan
}

func (a *App) Log() {
	for err := range a.errChan {
		a.log.Error("error", zap.Error(err))
	}
}

func (s *App) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			s.errChan <- fmt.Errorf("run app %v", r)
		}
	}()

	go s.Log()
	go s.SubJoinRoom(ctx)
	go s.SubLeaveRoom(ctx)
	go s.SubAddTracks(ctx)
	go s.SubRemoveTracks(ctx)
	go s.SubBackend(ctx)

	return nil
}

func (s *App) Done() {
	<-s.done
}
