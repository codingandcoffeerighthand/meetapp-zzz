package app

import "github.com/google/wire"

type App interface {
	Done()
}

var WireSet = wire.NewSet(NewAppInterface)
