package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisLimiter implements Limiter using Redis
// Suitable for distributed deployments
type RedisLimiter struct {
	config *Config
	client *redis.Client
}

// NewRedisLimiter creates a new Redis-based rate limiter
func NewRedisLimiter(config *Config, client *redis.Client) (*RedisLimiter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if client == nil {
		return nil, ErrRedisConnectionRequired
	}

	return &RedisLimiter{
		config: config,
		client: client,
	}, nil
}

// Allow checks if a request is allowed using Redis
func (rl *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	// Use Redis pipeline for atomic operations
	pipe := rl.client.Pipeline()

	// Add current timestamp
	now := time.Now().UnixNano()
	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// Remove old entries outside the window
	windowStart := time.Now().Add(-rl.config.Window).UnixNano()
	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStart))

	// Count entries in window
	pipe.ZCard(ctx, redisKey)

	// Set expiration
	pipe.Expire(ctx, redisKey, rl.config.Window)

	// Execute pipeline
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("redis pipeline error: %w", err)
	}

	// Get count result (3rd command)
	count := cmds[2].(*redis.IntCmd).Val()

	// Check if within limit
	return count <= int64(rl.config.MaxRequests), nil
}

// Reset clears the rate limit for a key
func (rl *RedisLimiter) Reset(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	return rl.client.Del(ctx, redisKey).Err()
}

// Close does nothing for Redis limiter (connection managed externally)
func (rl *RedisLimiter) Close() error {
	return nil
}

// Info returns information about the limiter
func (rl *RedisLimiter) Info() string {
	return fmt.Sprintf("RedisLimiter(maxReq=%d, window=%s)",
		rl.config.MaxRequests, rl.config.Window)
}
