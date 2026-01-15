---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/metrics/**/*.go,internal/health/**/*.go,internal/ratelimit/**/*.go"
description: Observability Implementation Patterns
---

# Observability Implementation Patterns

## Core Principle

**Type-Safe Metrics**: Always use interface types for Prometheus metrics, not pointer types. Prometheus client returns interfaces, not pointers.

## Prometheus Metrics Types

### ✅ CORRECT: Use Interface Types

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusCollector struct {
    // ✅ CORRECT: Use interface types
    httpRequestsTotal   *prometheus.CounterVec      // ✅ Pointer to Vec
    httpRequestDuration *prometheus.HistogramVec    // ✅ Pointer to Vec
    
    dbConnectionsInUse  prometheus.Gauge            // ✅ Interface type (not pointer)
    dbConnectionsIdle   prometheus.Gauge            // ✅ Interface type (not pointer)
    authTokensIssued    prometheus.Counter          // ✅ Interface type (not pointer)
}

// Initialization
pc.dbConnectionsInUse = promauto.NewGauge(
    prometheus.GaugeOpts{
        Name: "db_connections_in_use",
        Help: "Number of database connections in use",
    },
)
```

### ❌ WRONG: Using Pointer Types for Gauge/Counter

```go
type PrometheusCollector struct {
    // ❌ WRONG: Pointer types cause compilation errors
    dbConnectionsInUse  *prometheus.Gauge    // ❌ Type error
    authTokensIssued    *prometheus.Counter  // ❌ Type error
}

// This FAILS: promauto.NewGauge returns prometheus.Gauge, not *prometheus.Gauge
pc.dbConnectionsInUse = promauto.NewGauge(...)  // ❌ Type mismatch
```

## Type Rules for Prometheus Metrics

| Metric Type | Return Type | Field Type |
|-------------|-------------|------------|
| `promauto.NewCounter()` | `prometheus.Counter` | `prometheus.Counter` (interface) |
| `promauto.NewGauge()` | `prometheus.Gauge` | `prometheus.Gauge` (interface) |
| `promauto.NewHistogram()` | `prometheus.Histogram` | `prometheus.Histogram` (interface) |
| `promauto.NewSummary()` | `prometheus.Summary` | `prometheus.Summary` (interface) |
| `promauto.NewCounterVec()` | `*prometheus.CounterVec` | `*prometheus.CounterVec` (pointer) |
| `promauto.NewGaugeVec()` | `*prometheus.GaugeVec` | `*prometheus.GaugeVec` (pointer) |
| `promauto.NewHistogramVec()` | `*prometheus.HistogramVec` | `*prometheus.HistogramVec` (pointer) |
| `promauto.NewSummaryVec()` | `*prometheus.SummaryVec` | `*prometheus.SummaryVec` (pointer) |

### Memory Rule

- **Vec types** (CounterVec, GaugeVec, etc.): Use **pointers** (`*prometheus.CounterVec`)
- **Simple types** (Counter, Gauge, etc.): Use **interface** (`prometheus.Counter`)

## Complete Metrics Collector Example

```go
package metrics

import (
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusCollector struct {
    config *Config

    // HTTP metrics (Vec = pointer)
    httpRequestsTotal   *prometheus.CounterVec
    httpRequestDuration *prometheus.HistogramVec
    httpRequestSize     *prometheus.HistogramVec
    httpResponseSize    *prometheus.HistogramVec

    // Database metrics (simple = interface)
    dbConnectionsInUse prometheus.Gauge
    dbConnectionsIdle  prometheus.Gauge
    
    // Database metrics (Vec = pointer)
    dbQueryDuration *prometheus.HistogramVec
    dbQueriesTotal  *prometheus.CounterVec

    // Redis metrics (Vec = pointer)
    redisCommandsTotal   *prometheus.CounterVec
    redisCommandDuration *prometheus.HistogramVec

    // Auth metrics (mixed)
    authAttemptsTotal *prometheus.CounterVec  // Vec = pointer
    authTokensIssued  prometheus.Counter      // Simple = interface

    // Rate limit metrics (Vec = pointer)
    rateLimitHits *prometheus.CounterVec
}

func NewPrometheusCollector(config *Config) (*PrometheusCollector, error) {
    pc := &PrometheusCollector{config: config}

    // ✅ Vec types: Assign to pointer fields
    pc.httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Namespace: config.Namespace,
            Name:      "http_requests_total",
            Help:      "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    // ✅ Simple types: Assign to interface fields
    pc.dbConnectionsInUse = promauto.NewGauge(
        prometheus.GaugeOpts{
            Namespace: config.Namespace,
            Name:      "db_connections_in_use",
            Help:      "Database connections in use",
        },
    )

    pc.authTokensIssued = promauto.NewCounter(
        prometheus.CounterOpts{
            Namespace: config.Namespace,
            Name:      "auth_tokens_issued_total",
            Help:      "Total authentication tokens issued",
        },
    )

    return pc, nil
}

// Methods implementation
func (pc *PrometheusCollector) RecordHTTPRequest(method, path, status string, duration time.Duration) {
    // ✅ Vec types: Use WithLabelValues()
    pc.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
    pc.httpRequestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())
}

func (pc *PrometheusCollector) SetDatabaseConnections(inUse, idle int) {
    // ✅ Simple types: Use direct methods
    pc.dbConnectionsInUse.Set(float64(inUse))
    pc.dbConnectionsIdle.Set(float64(idle))
}

func (pc *PrometheusCollector) IncrementAuthTokens() {
    // ✅ Simple Counter: Direct Inc()
    pc.authTokensIssued.Inc()
}
```

## Health Check Patterns

### Interface Design

```go
type Checker interface {
    Check(ctx context.Context) error
    Name() string
}

// Aggregate health checker
type HealthChecker struct {
    checks []Checker
}

func (hc *HealthChecker) CheckAll(ctx context.Context) map[string]string {
    results := make(map[string]string)
    
    for _, check := range hc.checks {
        if err := check.Check(ctx); err != nil {
            results[check.Name()] = "unhealthy: " + err.Error()
        } else {
            results[check.Name()] = "healthy"
        }
    }
    
    return results
}
```

### Database Health Check

```go
type DatabaseChecker struct {
    db *pgxpool.Pool
}

func (dc *DatabaseChecker) Check(ctx context.Context) error {
    // ✅ Use timeout
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // ✅ Simple ping query
    return dc.db.Ping(ctx)
}

func (dc *DatabaseChecker) Name() string {
    return "database"
}
```

### Redis Health Check

```go
type RedisChecker struct {
    client *redis.Client
}

func (rc *RedisChecker) Check(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    return rc.client.Ping(ctx).Err()
}

func (rc *RedisChecker) Name() string {
    return "redis"
}
```

## Rate Limiter Interface Pattern

### Abstraction Layer

```go
package ratelimit

import (
    "context"
    "time"
)

// Limiter is the interface for rate limiting
type Limiter interface {
    Allow(ctx context.Context, key string) (bool, error)
    Reset(ctx context.Context, key string) error
}

// Config for rate limiter
type Config struct {
    Requests int
    Window   time.Duration
}
```

### Memory Implementation

```go
type MemoryLimiter struct {
    config *Config
    mu     sync.RWMutex
    data   map[string]*bucket
}

type bucket struct {
    count      int
    resetAt    time.Time
}

func (ml *MemoryLimiter) Allow(ctx context.Context, key string) (bool, error) {
    ml.mu.Lock()
    defer ml.mu.Unlock()

    now := time.Now()
    b, exists := ml.data[key]
    
    if !exists || now.After(b.resetAt) {
        ml.data[key] = &bucket{
            count:   1,
            resetAt: now.Add(ml.config.Window),
        }
        return true, nil
    }

    if b.count < ml.config.Requests {
        b.count++
        return true, nil
    }

    return false, nil
}
```

### Redis Implementation

```go
type RedisLimiter struct {
    client *redis.Client
    config *Config
}

func (rl *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
    pipe := rl.client.Pipeline()
    
    // Increment counter
    incrCmd := pipe.Incr(ctx, key)
    
    // Set expiry on first request
    pipe.Expire(ctx, key, rl.config.Window)
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }

    count := incrCmd.Val()
    return count <= int64(rl.config.Requests), nil
}
```

## Factory Pattern for Implementations

### Metrics Factory

```go
package metrics

import "errors"

type CollectorType string

const (
    CollectorTypePrometheus CollectorType = "prometheus"
    CollectorTypeNoop      CollectorType = "noop"
)

func NewCollector(collectorType CollectorType, config *Config) (Collector, error) {
    switch collectorType {
    case CollectorTypePrometheus:
        return NewPrometheusCollector(config)
    case CollectorTypeNoop:
        return &NoopCollector{}, nil
    default:
        return nil, errors.New("unknown collector type")
    }
}
```

### Health Check Factory

```go
package health

import (
    "github.com/jackc/pgxpool"
    "github.com/redis/go-redis/v9"
)

func NewChecker(db *pgxpool.Pool, redisClient *redis.Client) *HealthChecker {
    checks := []Checker{
        &DatabaseChecker{db: db},
        &RedisChecker{client: redisClient},
    }
    
    return &HealthChecker{checks: checks}
}
```

### Rate Limiter Factory

```go
package ratelimit

import (
    "errors"
    
    "github.com/redis/go-redis/v9"
)

type LimiterType string

const (
    LimiterTypeMemory LimiterType = "memory"
    LimiterTypeRedis  LimiterType = "redis"
)

func NewLimiter(limiterType LimiterType, config *Config, redisClient *redis.Client) (Limiter, error) {
    switch limiterType {
    case LimiterTypeMemory:
        return NewMemoryLimiter(config), nil
    case LimiterTypeRedis:
        if redisClient == nil {
            return nil, errors.New("redis client required for redis limiter")
        }
        return NewRedisLimiter(redisClient, config), nil
    default:
        return nil, errors.New("unknown limiter type")
    }
}
```

## Common Mistakes

### ❌ DON'T Do This

```go
// 1. Wrong metric types
type Collector struct {
    counter *prometheus.Counter  // ❌ Should be prometheus.Counter (interface)
    gauge   *prometheus.Gauge    // ❌ Should be prometheus.Gauge (interface)
}

// 2. Missing context timeouts in health checks
func (dc *DatabaseChecker) Check(ctx context.Context) error {
    return dc.db.Ping(ctx)  // ❌ No timeout, could hang forever
}

// 3. Not using factory pattern
limiter := &MemoryLimiter{...}  // ❌ Direct construction, hard to swap implementations

// 4. Exposing implementation details
type MyService struct {
    prometheusCollector *PrometheusCollector  // ❌ Concrete type
}

// 5. Not handling nil in health checks
func (rc *RedisChecker) Check(ctx context.Context) error {
    return rc.client.Ping(ctx).Err()  // ❌ Panic if client is nil
}
```

### ✅ DO This

```go
// 1. Correct metric types
type Collector struct {
    counter prometheus.Counter  // ✅ Interface for simple types
    gauge   prometheus.Gauge    // ✅ Interface for simple types
    counterVec *prometheus.CounterVec  // ✅ Pointer for Vec types
}

// 2. Always use timeouts
func (dc *DatabaseChecker) Check(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    return dc.db.Ping(ctx)  // ✅ Will timeout after 2 seconds
}

// 3. Use factory pattern
limiter, err := ratelimit.NewLimiter(ratelimit.LimiterTypeRedis, config, client)

// 4. Use interface types
type MyService struct {
    metrics metrics.Collector  // ✅ Interface type
}

// 5. Validate inputs
func (rc *RedisChecker) Check(ctx context.Context) error {
    if rc.client == nil {
        return errors.New("redis client not initialized")
    }
    return rc.client.Ping(ctx).Err()  // ✅ Safe
}
```

## Configuration Validation

```go
type Config struct {
    Namespace   string
    Subsystem   string
    HTTPBuckets []float64
}

func (c *Config) Validate() error {
    if c.Namespace == "" {
        return errors.New("namespace is required")
    }
    if len(c.HTTPBuckets) == 0 {
        c.HTTPBuckets = prometheus.DefBuckets  // Default values
    }
    return nil
}

// Always validate before use
func NewPrometheusCollector(config *Config) (*PrometheusCollector, error) {
    if err := config.Validate(); err != nil {
        return nil, err
    }
    // ... initialization
}
```

## Testing Patterns

### Mock Collector

```go
type MockCollector struct {
    mock.Mock
}

func (m *MockCollector) RecordHTTPRequest(method, path, status string, duration time.Duration) {
    m.Called(method, path, status, duration)
}

// Test usage
func TestMetricsRecording(t *testing.T) {
    mockCollector := new(MockCollector)
    mockCollector.On("RecordHTTPRequest", "GET", "/health", "200", mock.Anything).Return()
    
    // ... test code
    
    mockCollector.AssertExpectations(t)
}
```

## Checklist for Observability

- [ ] Use `prometheus.Counter` (interface) for simple Counter metrics
- [ ] Use `prometheus.Gauge` (interface) for simple Gauge metrics
- [ ] Use `*prometheus.CounterVec` (pointer) for Counter with labels
- [ ] Use `*prometheus.HistogramVec` (pointer) for Histogram with labels
- [ ] Add timeouts to all health checks (2-5 seconds)
- [ ] Validate configuration before initialization
- [ ] Use factory pattern for implementation selection
- [ ] Use interface types in service dependencies
- [ ] Handle nil checks in health checkers
- [ ] Test with mock implementations
- [ ] Document metric names and labels
- [ ] Use consistent naming conventions

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial observability patterns guide |

**Remember**: Prometheus metrics use interface types for simple metrics and pointer types for Vec metrics. Always validate configurations and use timeouts in health checks.
