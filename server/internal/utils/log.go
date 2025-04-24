package utils

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, func(), error) {
	log, err := zap.NewDevelopment()
	cl := func() { log.Sync() }
	return log, cl, err
}
