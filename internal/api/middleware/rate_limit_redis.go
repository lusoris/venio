// Package middleware contains HTTP middleware functions
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter implements distributed rate limiting using Redis
type RedisRateLimiter struct {
	client  *redis.Client
	maxReqs int
	window  time.Duration
	prefix  string
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
// maxReqs: maximum number of requests allowed
// window: time window for counting requests
// redisClient: Redis client connection
func NewRedisRateLimiter(redisClient *redis.Client, maxReqs int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:  redisClient,
		maxReqs: maxReqs,
		window:  window,
		prefix:  "ratelimit:",
	}
}

// Middleware returns a Gin middleware for Redis-based rate limiting
func (rl *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		allowed, remaining, resetTime, err := rl.Allow(c.Request.Context(), ip)
		if err != nil {
			// Log error but fail open (allow request if Redis is down)
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.maxReqs))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			c.Header("Retry-After", strconv.FormatInt(int64(time.Until(resetTime).Seconds()), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
				"message": fmt.Sprintf(
					"Rate limit exceeded: %d requests per %v. Try again in %v",
					rl.maxReqs, rl.window, time.Until(resetTime).Round(time.Second),
				),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow checks if the client IP is within rate limit using Redis
// Returns: allowed, remaining requests, reset time, error
func (rl *RedisRateLimiter) Allow(ctx context.Context, ip string) (bool, int, time.Time, error) {
	key := rl.prefix + ip
	now := time.Now()

	// Use Redis pipeline for atomic operations
	pipe := rl.client.Pipeline()

	// Increment counter
	incrCmd := pipe.Incr(ctx, key)

	// Set expiration on first request
	pipe.Expire(ctx, key, rl.window)

	// Get TTL to calculate reset time
	ttlCmd := pipe.TTL(ctx, key)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, now, fmt.Errorf("redis pipeline failed: %w", err)
	}

	// Get results
	count, err := incrCmd.Result()
	if err != nil {
		return false, 0, now, fmt.Errorf("failed to get count: %w", err)
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return false, 0, now, fmt.Errorf("failed to get TTL: %w", err)
	}

	// Calculate remaining and reset time
	remaining := rl.maxReqs - int(count)
	if remaining < 0 {
		remaining = 0
	}

	resetTime := now.Add(ttl)

	// Check if limit exceeded
	allowed := count <= int64(rl.maxReqs)

	return allowed, remaining, resetTime, nil
}

// RedisAuthRateLimiter creates a stricter rate limiter for auth endpoints
func RedisAuthRateLimiter(client *redis.Client) *RedisRateLimiter {
	return NewRedisRateLimiter(client, 5, 1*time.Minute)
}

// RedisGeneralRateLimiter creates a general rate limiter for API endpoints
func RedisGeneralRateLimiter(client *redis.Client) *RedisRateLimiter {
	return NewRedisRateLimiter(client, 100, 1*time.Minute)
}
