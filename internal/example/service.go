package example

import (
	"context"
	"errors"

	"github.com/go-core-fx/cachefx/cache"
	"go.uber.org/zap"
)

// Service implements the business logic for the example module.
//
// This struct represents the service layer in a clean architecture pattern.
// The service is responsible for implementing business rules, orchestrating
// operations between different components, and exposing functionality to
// the rest of the application.
//
// The Service depends on:
//   - Config: For configuration parameters
//   - Repository: For data access operations
//   - Storage: For caching the example value
//   - Metrics: For collecting and reporting metrics
//   - Logger: For structured logging
//
// In a real application, this struct would contain methods for:
//   - Business operations and workflows
//   - Coordinating between repositories and other services
//   - Enforcing business rules and validation
//   - Handling errors and logging
type Service struct {
	// config holds the configuration for the service
	config Config

	// examples provides access to data operations
	examples *Repository

	// storage caches the example value (cache-aside)
	storage *Storage

	// metrics is used for collecting and reporting metrics
	metrics *Metrics

	// logger is used for structured logging
	logger *zap.Logger
}

// New creates and initializes a new Service instance.
//
// This function serves as a constructor for the Service struct, accepting
// all its dependencies as parameters. This approach, known as dependency
// injection, makes the code more testable and maintainable.
//
// Parameters:
//   - config: Configuration for the service
//   - examples: Repository for data access
//   - storage: Cache storage for the example value
//   - metrics: Metrics collector for monitoring
//   - logger: Logger for structured logging
//
// Returns:
//   - *Service: A pointer to the newly created Service instance
//
// Example:
//
//	config := example.Config{Example: "demo"}
//	repo := example.NewRepository()
//	storage := example.NewStorage(cache.NewMemory(time.Hour))
//	metrics := example.NewMetrics()
//	logger, _ := zap.NewProduction()
//
//	service := example.New(config, repo, storage, metrics, logger)
func New(config Config, examples *Repository, storage *Storage, metrics *Metrics, logger *zap.Logger) *Service {
	return &Service{
		config:   config,
		examples: examples,
		storage:  storage,
		metrics:  metrics,
		logger:   logger,
	}
}

// Run starts the service and blocks until the context is canceled.
// It is called automatically by fxutil.RegisterRunnable in a managed goroutine.
func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("example service started")

	<-ctx.Done()

	s.logger.Info("example service stopped")

	return nil
}

// Example returns the example value, serving it from the module cache when
// possible (cache-aside): a miss computes the value from config and populates
// the cache; cache failures are logged and never fail the call.
func (s *Service) Example(ctx context.Context) string {
	value, err := s.storage.Get(ctx)
	switch {
	case err == nil:
		s.logger.Debug("example value served from cache")

		return value
	case errors.Is(err, cache.ErrKeyNotFound):
		s.logger.Debug("example value cache miss")
	default:
		s.logger.Warn("failed to read example value from cache", zap.Error(err))
	}

	value = s.config.Example
	if err = s.storage.Set(ctx, value); err != nil {
		s.logger.Warn("failed to write example value to cache", zap.Error(err))
	}

	return value
}
