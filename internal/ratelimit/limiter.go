// Package ratelimit provides rate limiting abstractions and implementations
package ratelimit

import (
	"context"
	"time"
)

// Limiter defines the interface for rate limiting
type Limiter interface {
	// Allow checks if a request from the given key is allowed
	// Returns true if allowed, false if rate limit exceeded
	Allow(ctx context.Context, key string) (bool, error)

	// Reset clears the rate limit counter for the given key
	Reset(ctx context.Context, key string) error

	// Close releases any resources held by the limiter
	Close() error
}

// Config holds rate limiter configuration
type Config struct {
	// MaxRequests is the maximum number of requests allowed in the window
	MaxRequests int

	// Window is the time window for counting requests
	Window time.Duration

	// BurstSize is the maximum burst size (optional, defaults to MaxRequests)
	BurstSize int
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.MaxRequests <= 0 {
		return ErrInvalidMaxRequests
	}
	if c.Window <= 0 {
		return ErrInvalidWindow
	}
	if c.BurstSize == 0 {
		c.BurstSize = c.MaxRequests
	}
	return nil
}
