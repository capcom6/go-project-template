// Package agent provides an example AI agent module with LLM integration via OpenRouter.
//
// The agent module demonstrates a complete AI agent pattern:
//   - config.go: Configuration for the agent (model, system prompt, poll interval)
//   - domain.go: Core domain entities (Task)
//   - errors.go: Module-specific error types
//   - metrics.go: Prometheus metrics for monitoring
//   - models.go: Internal data models
//   - module.go: FX module definition for dependency injection
//   - repository.go: Data access layer for tasks
//   - service.go: Agent loop with LLM integration
//
// Architecture:
//
//	The agent runs a background loop that polls for pending tasks, processes them
//	using the OpenRouter LLM API, and stores results. Tasks can be enqueued via
//	HTTP API or Telegram bot commands.
//
// Dependencies:
//   - github.com/revrost/go-openrouter: OpenRouter LLM client
//   - go.uber.org/fx: Dependency injection
//   - go.uber.org/zap: Structured logging
//   - github.com/prometheus/client_golang: Metrics collection
//
// Usage:
//
//	app := fx.New(
//	    agent.Module(),
//	    // other modules...
//	)
//	app.Run()
package agent
