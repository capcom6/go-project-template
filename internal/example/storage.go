package example

import (
	"context"
	"fmt"
	"time"

	"github.com/go-core-fx/cachefx/cache"
)

// cacheTTL bounds how long the example value stays cached.
const cacheTTL = time.Minute

// exampleValueKey is the cache key of the configured example value.
const exampleValueKey = "example:value"

// Storage wraps the module's private cache for string values.
type Storage struct {
	cache cache.Cache
}

// NewStorage creates a Storage over the module's private cache.
func NewStorage(c cache.Cache) *Storage {
	return &Storage{cache: c}
}

// Get returns the cached example value.
//
// Returns a wrapped cache.ErrKeyNotFound when the value is missing or expired.
func (s *Storage) Get(ctx context.Context) (string, error) {
	data, err := s.cache.Get(ctx, exampleValueKey)
	if err != nil {
		return "", fmt.Errorf("failed to get example value: %w", err)
	}

	return string(data), nil
}

// Set stores the example value with the module TTL.
func (s *Storage) Set(ctx context.Context, value string) error {
	if err := s.cache.Set(ctx, exampleValueKey, []byte(value), cache.WithTTL(cacheTTL)); err != nil {
		return fmt.Errorf("failed to set example value: %w", err)
	}

	return nil
}
