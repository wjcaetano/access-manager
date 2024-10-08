package modules

import (
	"access-manager/internal/config"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	options := []fx.Option{
		InternalModule,
	}

	if !config.IsLocalScope() {
		options = append(options, fx.NopLogger)
	}

	return fx.New(
		fx.Options(options...),
	)
}
