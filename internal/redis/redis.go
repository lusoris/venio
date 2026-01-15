// Package redis provides Redis client initialization
package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lusoris/venio/internal/config"
)

// Client wraps the Redis client
type Client struct {
	*redis.Client
}

// Connect establishes a connection to Redis
func Connect(ctx context.Context, cfg *config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Address(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("✓ Redis connection established")
	return &Client{Client: client}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
}
