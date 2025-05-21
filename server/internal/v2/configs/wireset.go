package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "ServerConfig"),
	wire.FieldsOf(new(Config), "CloudflareConfig"),
	wire.FieldsOf(new(Config), "Web3Config"),
	wire.FieldsOf(new(Config), "WsService"),
)
