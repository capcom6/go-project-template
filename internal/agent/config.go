package agent

import "time"

type Config struct {
	Model        string        `koanf:"model"`
	SystemPrompt string        `koanf:"system_prompt"`
	PollInterval time.Duration `koanf:"poll_interval"`
	BatchSize    int           `koanf:"batch_size"`
}
