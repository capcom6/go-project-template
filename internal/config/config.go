package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-core-fx/config"
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

type telegram struct {
	Token string `koanf:"token"`
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

type agentConfig struct {
	APIKey       string        `koanf:"api_key"`
	Model        string        `koanf:"model"`
	SystemPrompt string        `koanf:"system_prompt"`
	PollInterval time.Duration `koanf:"poll_interval"`
	BatchSize    int           `koanf:"batch_size"`
}

type Config struct {
	HTTP     http           `koanf:"http"`
	Telegram telegram       `koanf:"telegram"`
	Database databaseConfig `koanf:"database"`

	Example exampleConfig `koanf:"example"`
	Agent   agentConfig   `koanf:"agent"`
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
		Telegram: telegram{
			Token: "",
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

		Agent: agentConfig{
			APIKey:       "",
			Model:        "openai/gpt-4o-mini",
			SystemPrompt: "You are a helpful AI assistant.",
			PollInterval: 5 * time.Second,
			BatchSize:    1,
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
