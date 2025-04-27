package infras_cloudflare

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewClient,
)
