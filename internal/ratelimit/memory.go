package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryLimiter implements Limiter using in-memory storage
// Suitable for single-instance deployments
type MemoryLimiter struct {
	config   *Config
	requests map[string][]time.Time
	mu       sync.Mutex
	stopCh   chan struct{}
}

// NewMemoryLimiter creates a new memory-based rate limiter
func NewMemoryLimiter(config *Config) (*MemoryLimiter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	ml := &MemoryLimiter{
		config:   config,
		requests: make(map[string][]time.Time),
		stopCh:   make(chan struct{}),
	}

	// Start cleanup goroutine
	go ml.cleanupLoop()

	return ml, nil
}

// Allow checks if a request is allowed
func (ml *MemoryLimiter) Allow(ctx context.Context, key string) (bool, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-ml.config.Window)

	// Get or create request list
	requests, exists := ml.requests[key]
	if !exists {
		ml.requests[key] = []time.Time{now}
		return true, nil
	}

	// Filter requests within window
	validRequests := make([]time.Time, 0, len(requests))
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	// Check limit
	if len(validRequests) >= ml.config.MaxRequests {
		ml.requests[key] = validRequests
		return false, nil
	}

	// Allow request
	validRequests = append(validRequests, now)
	ml.requests[key] = validRequests
	return true, nil
}

// Reset clears the rate limit for a key
func (ml *MemoryLimiter) Reset(ctx context.Context, key string) error {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	delete(ml.requests, key)
	return nil
}

// Close stops the cleanup goroutine
func (ml *MemoryLimiter) Close() error {
	close(ml.stopCh)
	return nil
}

// cleanupLoop periodically removes old entries
func (ml *MemoryLimiter) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ml.cleanup()
		case <-ml.stopCh:
			return
		}
	}
}

// cleanup removes expired entries
func (ml *MemoryLimiter) cleanup() {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-ml.config.Window)

	for key, requests := range ml.requests {
		validRequests := make([]time.Time, 0, len(requests))
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}

		if len(validRequests) == 0 {
			delete(ml.requests, key)
		} else {
			ml.requests[key] = validRequests
		}
	}
}

// Info returns information about the limiter
func (ml *MemoryLimiter) Info() string {
	return fmt.Sprintf("MemoryLimiter(maxReq=%d, window=%s)",
		ml.config.MaxRequests, ml.config.Window)
}
