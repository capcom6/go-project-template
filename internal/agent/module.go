package agent

import (
	"github.com/go-core-fx/fxutil"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"agent",
		logger.WithNamedLogger("agent"),

		fx.Provide(NewMetrics, fx.Private),
		fx.Provide(NewRepository, fx.Private),
		fx.Provide(New),

		fx.Invoke(fxutil.RegisterRunnable[*Service]()),
	)
}
