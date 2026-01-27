package config

import (
	"github.com/capcom6/go-project-template/internal/example"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(New),
		fx.Provide(func(cfg Config) fiberfx.Config {
			return fiberfx.Config{
				Address:     cfg.HTTP.Address,
				ProxyHeader: cfg.HTTP.ProxyHeader,
				Proxies:     cfg.HTTP.Proxies,
			}
		}),
		fx.Provide(func(cfg Config) openapi.Config {
			return openapi.Config{
				Enabled:    cfg.HTTP.OpenAPI.Enabled,
				PublicHost: cfg.HTTP.OpenAPI.PublicHost,
				PublicPath: cfg.HTTP.OpenAPI.PublicPath,
			}
		}),
		fx.Provide(func(cfg Config) example.Config {
			return example.Config{
				Example: cfg.Example.Example,
			}
		}),
	)
}
