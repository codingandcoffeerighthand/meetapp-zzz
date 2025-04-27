package smc_infra

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewSMCInfra,
)
