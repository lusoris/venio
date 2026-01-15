// Package metrics provides abstraction for application metrics collection
package metrics

import (
	"time"
)

// Collector defines the interface for metrics collection
type Collector interface {
	// HTTP Metrics
	RecordHTTPRequest(method, path, status string, duration time.Duration, requestSize, responseSize int64)

	// Database Metrics
	RecordDBQuery(operation, status string, duration time.Duration)
	UpdateDBConnections(inUse, idle int)

	// Redis Metrics
	RecordRedisCommand(command, status string, duration time.Duration)

	// Auth Metrics
	RecordAuthAttempt(authType, status string)
	RecordTokenIssued()

	// Rate Limit Metrics
	RecordRateLimitHit(limiter, status string)

	// Custom Metrics
	IncCounter(name string, labels map[string]string, value float64)
	ObserveHistogram(name string, labels map[string]string, value float64)
	SetGauge(name string, labels map[string]string, value float64)
}

// Config holds metrics collector configuration
type Config struct {
	// Enabled controls whether metrics collection is active
	Enabled bool

	// Namespace is the metrics namespace prefix
	Namespace string

	// Subsystem is the metrics subsystem
	Subsystem string

	// HTTPBuckets are histogram buckets for HTTP request duration
	HTTPBuckets []float64

	// DBBuckets are histogram buckets for database query duration
	DBBuckets []float64

	// RedisBuckets are histogram buckets for Redis command duration
	RedisBuckets []float64
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:      true,
		Namespace:    "venio",
		Subsystem:    "",
		HTTPBuckets:  []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		DBBuckets:    []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		RedisBuckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1, .25, .5},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Namespace == "" {
		c.Namespace = "venio"
	}
	if len(c.HTTPBuckets) == 0 {
		c.HTTPBuckets = DefaultConfig().HTTPBuckets
	}
	if len(c.DBBuckets) == 0 {
		c.DBBuckets = DefaultConfig().DBBuckets
	}
	if len(c.RedisBuckets) == 0 {
		c.RedisBuckets = DefaultConfig().RedisBuckets
	}
	return nil
}
