package bot

import (
	"github.com/capcom6/go-project-template/internal/bot/handler"
	"github.com/capcom6/go-project-template/internal/bot/handlers/start"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"go.uber.org/fx"
)

// Module wires the bot's business logic: the command handlers and their
// registration on the shared telegofx router. The Telegram transport (HTTP
// client and telegofx options) lives in the telegram module.
func Module() fx.Option {
	return fx.Module(
		"bot",
		logger.WithNamedLogger("bot"),
		fx.Provide(
			fx.Annotate(start.New, fx.ResultTags(`group:"handlers"`)),
		),
		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, r *telegofx.Router) {
					for _, h := range handlers {
						h.Register(r)
					}
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
