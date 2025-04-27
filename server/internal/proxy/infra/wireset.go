package infra

import (
	infras_cloudflare "proxy-srv/internal/proxy/infra/cloudflare"
	"proxy-srv/internal/proxy/infra/smc_infra"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	smc_infra.WireSet,
	infras_cloudflare.WireSet,
)
