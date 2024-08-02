package modules

import (
	"access-manager/internal/config"
	"access-manager/internal/db"
	"access-manager/internal/server"

	"go.uber.org/fx"
)

var internalFactories = fx.Provide(
	config.NewConfig,
	server.NewServer,
	db.NewDatabase,
)

var InternalModule = fx.Options(
	internalFactories,
	fx.Invoke(server.StartHTTPServer),
)
