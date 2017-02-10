package store

import (
	"time"

	"github.com/ulule/limiter"
)

// Store is the common interface for datastores.
type Store interface {
	Exists(key string) (bool, error)
	Get(key string) (string, error)
	Set(key string, value string, expire int64) error
	ToLimiterStore(prefix string) (limiter.Store, error)
}

// Options are options for store.
type Options struct {
	// Prefix is the prefix to use for the key.
	Prefix string

	// MaxRetry is the maximum number of retry under race conditions.
	MaxRetry int

	// CleanUpInterval is the interval for cleanup.
	CleanUpInterval time.Duration
}
