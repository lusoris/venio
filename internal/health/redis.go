package health

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisChecker checks Redis health
type RedisChecker struct {
	client *redis.Client
	name   string
}

// NewRedisChecker creates a new Redis health checker
func NewRedisChecker(client *redis.Client) *RedisChecker {
	return &RedisChecker{
		client: client,
		name:   "redis",
	}
}

// Name returns the checker name
func (r *RedisChecker) Name() string {
	return r.name
}

// Check performs the Redis health check
func (r *RedisChecker) Check(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:      r.name,
		Timestamp: start,
	}

	// Ping Redis
	if err := r.client.Ping(ctx).Err(); err != nil {
		check.Status = StatusUnhealthy
		check.Message = "Redis ping failed: " + err.Error()
		check.ResponseTime = time.Since(start)
		return check
	}

	// Get Redis info
	info, err := r.client.Info(ctx, "server", "memory").Result()
	if err == nil {
		check.Metadata = map[string]interface{}{
			"info": info,
		}
	}

	// Get pool stats
	stats := r.client.PoolStats()
	check.Metadata = map[string]interface{}{
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"timeouts":    stats.Timeouts,
		"total_conns": stats.TotalConns,
		"idle_conns":  stats.IdleConns,
		"stale_conns": stats.StaleConns,
	}

	// Check if connection pool is healthy
	if stats.Timeouts > 0 {
		check.Status = StatusDegraded
		check.Message = "Redis connection timeouts detected"
	} else {
		check.Status = StatusHealthy
		check.Message = "Redis connection successful"
	}

	check.ResponseTime = time.Since(start)
	return check
}
