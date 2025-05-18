package app

import "go.uber.org/zap"

type app struct {
	clf         CloudFlareInfraV2
	meet        MeetInfra
	log         *zap.Logger
	errChan     chan error
	unsubscribe func()
}

func NewApp(cloudFlareInfa CloudFlareInfraV2, meetInfra MeetInfra, log *zap.Logger) *app {
	return &app{
		clf:     cloudFlareInfa,
		meet:    meetInfra,
		log:     log,
		errChan: make(chan error, 20),
	}
}

func (a *app) Err() <-chan error {
	return a.errChan
}

func (a *app) Log() {
	for err := range a.errChan {
		a.log.Error("error", zap.Error(err))
	}
}
