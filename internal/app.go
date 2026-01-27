package internal

import (
	"context"

	"github.com/capcom6/go-project-template/internal/config"
	"github.com/capcom6/go-project-template/internal/example"
	"github.com/capcom6/go-project-template/internal/server"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		// sqlfx.Module(),
		// goosefx.Module(),
		// bunfx.Module(),
		fiberfx.Module(),
		healthfx.Module(),
		//
		// APP MODULES
		config.Module(),
		// db.Module(),
		server.Module(),
		// bot.Module(),
		//
		// BUSINESS MODULES
		fx.Supply(version),
		example.Module(),
		//
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("app started")
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("app stopped")
					return nil
				},
			})
		}),
	).Run()
}
