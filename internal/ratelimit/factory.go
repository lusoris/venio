package ratelimit

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Type represents the rate limiter type
type Type string

const (
	// TypeMemory uses in-memory storage
	TypeMemory Type = "memory"

	// TypeRedis uses Redis storage
	TypeRedis Type = "redis"
)

// FactoryConfig holds configuration for the factory
type FactoryConfig struct {
	Type        Type
	RedisClient *redis.Client
}

// Factory creates rate limiters based on configuration
type Factory struct {
	config *FactoryConfig
}

// NewFactory creates a new rate limiter factory
func NewFactory(config *FactoryConfig) *Factory {
	return &Factory{config: config}
}

// NewLimiter creates a new rate limiter with the given configuration
func (f *Factory) NewLimiter(config *Config) (Limiter, error) {
	switch f.config.Type {
	case TypeMemory:
		return NewMemoryLimiter(config)
	case TypeRedis:
		return NewRedisLimiter(config, f.config.RedisClient)
	default:
		return NewMemoryLimiter(config)
	}
}

// NewAuthLimiter creates a rate limiter for authentication endpoints
// Default: 5 attempts per minute
func (f *Factory) NewAuthLimiter() (Limiter, error) {
	return f.NewLimiter(&Config{
		MaxRequests: 5,
		Window:      1 * time.Minute,
	})
}

// NewGeneralLimiter creates a rate limiter for general API endpoints
// Default: 100 requests per minute
func (f *Factory) NewGeneralLimiter() (Limiter, error) {
	return f.NewLimiter(&Config{
		MaxRequests: 100,
		Window:      1 * time.Minute,
	})
}

// NewAdminLimiter creates a rate limiter for admin endpoints
// Default: 200 requests per minute (higher limit)
func (f *Factory) NewAdminLimiter() (Limiter, error) {
	return f.NewLimiter(&Config{
		MaxRequests: 200,
		Window:      1 * time.Minute,
	})
}

// NewStrictLimiter creates a strict rate limiter
// Default: 3 attempts per 5 minutes
func (f *Factory) NewStrictLimiter() (Limiter, error) {
	return f.NewLimiter(&Config{
		MaxRequests: 3,
		Window:      5 * time.Minute,
	})
}
