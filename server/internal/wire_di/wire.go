//go:build wireinject
// +build wireinject

package wiredi

import (
	"proxy-srv/internal/proxy/domain"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	domain.NewDumDecryptService,
)

func Init() (*domain.Ns, error) {
	wire.Build(WireSet)
	_ = domain.Ns{}
	return nil, nil
}
