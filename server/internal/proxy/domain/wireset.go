package domain

import (
	"github.com/google/wire"
)

type Ns struct {
	ns
}

var WireSet = wire.NewSet(NewDumDecryptService,
	wire.Bind(new(Ns), new(ns)))
