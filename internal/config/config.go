package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-core-fx/config"
	"github.com/go-core-fx/telegofx"
)

type http struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`

	OpenAPI openAPIConfig `koanf:"openapi"`
}

type openAPIConfig struct {
	Enabled    bool   `koanf:"enabled"`
	PublicHost string `koanf:"public_host"`
	PublicPath string `koanf:"public_path"`
}

type telegramConfig struct {
	Token    string        `koanf:"token"`
	ProxyURL string        `koanf:"proxy_url"`
	Timeout  time.Duration `koanf:"timeout"`
	Mode     telegofx.Mode `koanf:"mode"`
}

type databaseConfig struct {
	URL             string        `koanf:"url"`
	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
}

type exampleConfig struct {
	Example string `koanf:"example"`
}

type Config struct {
	HTTP     http           `koanf:"http"`
	Telegram telegramConfig `koanf:"telegram"`
	Database databaseConfig `koanf:"database"`

	Example exampleConfig `koanf:"example"`
}

func Default() Config {
	//nolint:gosec // default values
	return Config{
		HTTP: http{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
			OpenAPI: openAPIConfig{
				Enabled:    true,
				PublicHost: "",
				PublicPath: "",
			},
		},
		Telegram: telegramConfig{
			Token:    "",
			ProxyURL: "",
			Timeout:  time.Minute,
			Mode:     telegofx.ModePolling,
		},
		Database: databaseConfig{
			URL:             "mariadb://example:example@127.0.0.1:3306/example?charset=utf8mb4&parseTime=True&loc=UTC",
			ConnMaxIdleTime: 0,
			ConnMaxLifetime: 0,
			MaxOpenConns:    0,
			MaxIdleConns:    0,
		},

		Example: exampleConfig{
			Example: "example",
		},
	}
}

func New() (Config, error) {
	cfg := Default()

	options := []config.Option{}
	if yamlPath := os.Getenv("CONFIG_PATH"); yamlPath != "" {
		options = append(options, config.WithLocalYAML(yamlPath))
	}

	if err := config.Load(&cfg, options...); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
