package config

import (
	"github.com/capcom6/go-project-template/internal/agent"
	"github.com/capcom6/go-project-template/internal/example"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/sqlfx"
	"github.com/go-core-fx/telegofx"
	openrouter "github.com/revrost/go-openrouter"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(New, fx.Private),
		fx.Provide(
			func(cfg Config) fiberfx.Config {
				return fiberfx.Config{
					Address:     cfg.HTTP.Address,
					ProxyHeader: cfg.HTTP.ProxyHeader,
					Proxies:     cfg.HTTP.Proxies,
				}
			},
			func(cfg Config) openapi.Config {
				return openapi.Config{
					Enabled:    cfg.HTTP.OpenAPI.Enabled,
					PublicHost: cfg.HTTP.OpenAPI.PublicHost,
					PublicPath: cfg.HTTP.OpenAPI.PublicPath,
				}
			},
			func(cfg Config) telegofx.Config {
				return telegofx.Config{
					Token: cfg.Telegram.Token,
				}
			},
			func(cfg Config) sqlfx.Config {
				return sqlfx.Config{
					URL:             cfg.Database.URL,
					ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
					ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
					MaxOpenConns:    cfg.Database.MaxOpenConns,
					MaxIdleConns:    cfg.Database.MaxIdleConns,
				}
			},
		),
		fx.Provide(func(cfg Config) example.Config {
			return example.Config{
				Example: cfg.Example.Example,
			}
		}),
		fx.Provide(func(cfg Config) agent.Config {
			return agent.Config{
				Model:        cfg.Agent.Model,
				SystemPrompt: cfg.Agent.SystemPrompt,
				PollInterval: cfg.Agent.PollInterval,
				BatchSize:    cfg.Agent.BatchSize,
			}
		}),
		fx.Provide(func(cfg Config) *openrouter.Client {
			return openrouter.NewClient(cfg.Agent.APIKey)
		}),
	)
}
