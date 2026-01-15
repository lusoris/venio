---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/**/*.go"
description: Observability Best Practices for Venio
---

# Observability Best Practices

This document provides guidelines for implementing observability features (metrics, logging, health checks) in Venio.

## Prometheus Metrics

### When to Add Metrics

Add metrics for:
- ✅ **HTTP endpoints** - Request rate, latency, errors
- ✅ **Database operations** - Query duration, connection pool
- ✅ **External API calls** - Success rate, duration
- ✅ **Background jobs** - Execution time, success/failure
- ✅ **Business events** - User registrations, logins, actions
- ✅ **Resource usage** - Connection pools, caches, queues

Do NOT add metrics for:
- ❌ Internal helper functions (unless performance-critical)
- ❌ One-time initialization code
- ❌ Unit test code

### Metric Naming Conventions

Follow Prometheus naming best practices:

```go
// ✅ DO: Use descriptive names with units
venio_http_request_duration_seconds
venio_db_query_duration_seconds
venio_cache_hits_total
venio_user_registrations_total

// ❌ DON'T: Generic names or missing units
request_time
db_query
cache
users
```

**Pattern:** `{namespace}_{subsystem}_{metric}_{unit}`

- **namespace:** `venio`
- **subsystem:** `http`, `db`, `redis`, `auth`, etc.
- **metric:** Descriptive name
- **unit:** `seconds`, `bytes`, `total`, etc.

### Metric Types

**Counter** - Monotonically increasing value:
```go
// ✅ Use for: Request counts, error counts, events
authAttemptsTotal = promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "venio_auth_attempts_total",
        Help: "Total number of authentication attempts",
    },
    []string{"type", "status"},
)

// Record
authAttemptsTotal.WithLabelValues("login", "success").Inc()
```

**Gauge** - Value that can go up and down:
```go
// ✅ Use for: Current connections, queue size, active users
dbConnectionsInUse = promauto.NewGauge(
    prometheus.GaugeOpts{
        Name: "venio_db_connections_in_use",
        Help: "Number of database connections currently in use",
    },
)

// Record
dbConnectionsInUse.Set(float64(connections))
```

**Histogram** - Distribution of values:
```go
// ✅ Use for: Request duration, response size
httpRequestDuration = promauto.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "venio_http_request_duration_seconds",
        Help:    "HTTP request latency in seconds",
        Buckets: prometheus.DefBuckets, // or custom buckets
    },
    []string{"method", "path", "status"},
)

// Record
start := time.Now()
// ... handle request ...
httpRequestDuration.WithLabelValues("GET", "/users", "200").Observe(time.Since(start).Seconds())
```

### Label Best Practices

**✅ DO:**

```go
// Low-cardinality labels (few unique values)
httpRequestsTotal.WithLabelValues(
    "GET",              // method (limited values)
    "/api/v1/users",    // path template (limited values)
    "200",              // status code (limited values)
).Inc()
```

**❌ DON'T:**

```go
// High-cardinality labels (many unique values)
httpRequestsTotal.WithLabelValues(
    "GET",
    "/api/v1/users/12345",  // ❌ User ID - infinite cardinality!
    "200",
    "192.168.1.100",        // ❌ IP address - high cardinality!
).Inc()
```

**Cardinality Limits:**
- Keep total label combinations < 10,000
- Each label should have < 100 unique values
- Use labels for **filtering**, not **identification**

### Example: Adding Database Metrics

```go
// internal/repositories/user_repository.go

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
    start := time.Now()
    
    // Execute query
    user, err := r.queryUser(ctx, id)
    
    // Record metrics
    duration := time.Since(start)
    status := "success"
    if err != nil {
        status = "error"
    }
    
    middleware.RecordDBQuery("GetUserByID", duration, err)
    
    return user, err
}
```

## Structured Logging

### When to Log

**✅ DO Log:**
- Application startup/shutdown
- Authentication events (login, logout, failures)
- Authorization failures
- Database connection errors
- External API failures
- Background job execution
- Critical errors

**❌ DON'T Log:**
- Every function entry/exit (use tracing instead)
- Successful operations in hot paths (use metrics)
- Debug info in production (use DEBUG level)

### Log Levels

```go
// ERROR - Something failed and needs attention
logger.Error("Database connection failed",
    "host", cfg.Database.Host,
    "error", err)

// WARN - Something unexpected but handled
logger.Warn("Rate limit exceeded",
    "ip", clientIP,
    "endpoint", path)

// INFO - Important business events
logger.Info("User logged in",
    "user_id", user.ID,
    "email", user.Email)

// DEBUG - Detailed diagnostic information (disabled in production)
logger.Debug("Processing request",
    "method", r.Method,
    "path", r.URL.Path,
    "headers", r.Header)
```

### Contextual Logging

Use helper methods for domain-specific logs:

```go
// ✅ DO: Use domain-specific loggers
logger.HTTP().Info("Request completed",
    "method", method,
    "path", path,
    "status", status,
    "duration", duration)

logger.Auth().Warn("Failed login attempt",
    "email", email,
    "ip", clientIP,
    "reason", "invalid password")

logger.DB().Error("Query failed",
    "operation", "SELECT",
    "table", "users",
    "error", err)
```

### Sensitive Data Protection

**❌ NEVER Log:**
- Passwords (plain or hashed)
- JWT tokens
- API keys
- Credit card numbers
- Personal identifiable information (PII) without explicit need

```go
// ❌ BAD
logger.Info("User login", "email", email, "password", password)

// ✅ GOOD
logger.Info("User login", "email", email)

// ✅ GOOD (with sanitization)
logger.Info("User updated", "user", sanitizeUser(user))
```

### Example: Structured Logging in Handler

```go
func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.HTTP().Warn("Invalid login request",
            "ip", c.ClientIP(),
            "error", err)
        
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request",
            Message: "Please check your input",
        })
        return
    }
    
    accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        logger.Auth().Warn("Login failed",
            "email", req.Email,
            "ip", c.ClientIP(),
            "error", err)
        
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Error: "Authentication failed",
            Message: "Invalid email or password",
        })
        return
    }
    
    logger.Auth().Info("User logged in",
        "email", req.Email,
        "ip", c.ClientIP())
    
    // ... return tokens
}
```

## Health Checks

### When to Add Health Checks

Add health checks for:
- ✅ **Critical dependencies** - Database, Redis, message queues
- ✅ **External services** - APIs, authentication providers
- ✅ **Filesystem** - Required directories, disk space

Do NOT add for:
- ❌ Internal services (they're part of the app)
- ❌ Optional features
- ❌ Non-critical dependencies

### Liveness vs Readiness

**Liveness Probe** - Is the application alive?
```go
// ✅ Simple check - just confirm process is running
func (h *HealthHandler) Liveness(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    })
}
```

**Readiness Probe** - Can the application serve traffic?
```go
// ✅ Check all critical dependencies
func (h *HealthHandler) Readiness(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    allHealthy := true
    
    // Check database
    if err := h.db.Ping(ctx); err != nil {
        allHealthy = false
    }
    
    // Check Redis
    if err := h.redis.Ping(ctx).Err(); err != nil {
        allHealthy = false
    }
    
    if !allHealthy {
        c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
```

### Health Check Best Practices

**✅ DO:**
- Set reasonable timeouts (5 seconds max)
- Fail fast on critical dependencies
- Return 200 OK for healthy, 503 Service Unavailable for unhealthy
- Cache results for 10-30 seconds to avoid overhead
- Minimize detailed information in production

**❌ DON'T:**
- Make expensive operations in health checks
- Expose internal architecture details
- Check non-critical dependencies
- Return 500 errors (use 503 for unhealthy)

### Production-Safe Health Checks

```go
// ✅ DO: Minimal info in production
if h.env == "production" {
    if allHealthy {
        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    } else {
        c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
    }
    return
}

// Detailed info only in development
c.JSON(statusCode, HealthResponse{
    Status: status,
    Services: serviceStatuses,
})
```

## Alert Rules

### When to Create Alerts

Create alerts for:
- ✅ **Error rates** - High 5xx responses
- ✅ **Latency** - P95/P99 above SLA
- ✅ **Resource exhaustion** - Connection pool, memory
- ✅ **Service down** - Critical dependencies unavailable
- ✅ **Security events** - High auth failure rate

### Alert Rule Template

```yaml
- alert: DescriptiveAlertName
  expr: |
    # PromQL query
    rate(venio_http_requests_total{status=~"5.."}[5m]) > 0.05
  for: 5m  # Evaluation period (avoid flapping)
  labels:
    severity: critical|warning|info
    component: api|database|redis|auth
  annotations:
    summary: "Brief description"
    description: "Detailed message with {{ $value }}"
```

### Alert Best Practices

**✅ DO:**
- Choose appropriate thresholds based on SLAs
- Include context in annotations
- Set evaluation period to avoid flapping
- Categorize by severity
- Test alerts in staging first

**❌ DON'T:**
- Alert on everything (alert fatigue)
- Use overly sensitive thresholds
- Forget to document resolution steps
- Alert without actionable information

## Dashboard Design

### Essential Panels

**System Overview:**
- Request rate (RPS)
- Error rate (%)
- P95/P99 latency
- Active users/connections

**Resource Monitoring:**
- CPU usage
- Memory usage
- Disk I/O
- Network traffic

**Application Metrics:**
- Database connection pool
- Redis command rate
- Rate limit hits
- Authentication attempts

### Dashboard Best Practices

**✅ DO:**
- Group related metrics
- Use appropriate visualization (graph, gauge, table)
- Set meaningful Y-axis ranges
- Include thresholds/baselines
- Add descriptions to panels

**❌ DON'T:**
- Overcrowd dashboards
- Use default titles ("Panel 1")
- Mix unrelated metrics
- Forget to add time range selector

## Performance Considerations

### Metrics Overhead

- **Counter:** ~50ns per increment
- **Histogram:** ~200ns per observation
- **Gauge:** ~50ns per set

**Best Practice:** Metrics have minimal overhead. Use liberally!

### Label Cardinality Impact

```
Memory per metric ≈ 3KB × (label combinations)

Example:
- 10 methods × 50 paths × 10 status codes = 5,000 combinations
- Memory usage ≈ 3KB × 5,000 = 15MB ✅ OK

- 10 methods × ∞ user IDs × 10 status = ∞ combinations  
- Memory usage ≈ ∞ ❌ BAD!
```

## Testing Observability

### Testing Metrics

```go
func TestMetricsRecorded(t *testing.T) {
    // Setup
    router := setupRouter()
    
    // Execute request
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/health", nil)
    router.ServeHTTP(w, req)
    
    // Verify metrics (pseudo-code)
    assert.Equal(t, 1, getMetricValue("venio_http_requests_total"))
}
```

### Testing Health Checks

```go
func TestReadinessHealthy(t *testing.T) {
    // Setup with healthy dependencies
    handler := NewHealthHandler(db, redis, "1.0.0", "test")
    
    // Execute
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    handler.Readiness(c)
    
    // Verify
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "healthy")
}
```

## Summary Checklist

Before adding new code, ensure:

- [ ] Critical paths have metrics
- [ ] Metrics have appropriate types (Counter/Gauge/Histogram)
- [ ] Labels have low cardinality
- [ ] Errors are logged with context
- [ ] No sensitive data in logs/metrics
- [ ] Health checks for new dependencies
- [ ] Alerts for critical failures
- [ ] Dashboard panels for new metrics

## References

- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [Grafana Dashboard Best Practices](https://grafana.com/docs/grafana/latest/best-practices/best-practices-for-creating-dashboards/)
- [Go Structured Logging (slog)](https://pkg.go.dev/log/slog)
- [Venio Observability Guide](../docs/observability.md)
