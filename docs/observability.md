# Observability

This document describes the observability features of Venio, including metrics, logging, and health checks.

## Table of Contents

- [Metrics](#metrics)
- [Health Checks](#health-checks)
- [Logging](#logging)
- [Dashboards](#dashboards)

## Metrics

Venio exposes Prometheus metrics at the `/metrics` endpoint.

### Available Metrics

#### HTTP Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `venio_http_requests_total` | Counter | Total number of HTTP requests | method, path, status |
| `venio_http_request_duration_seconds` | Histogram | HTTP request latency | method, path, status |
| `venio_http_request_size_bytes` | Histogram | HTTP request size | method, path |
| `venio_http_response_size_bytes` | Histogram | HTTP response size | method, path, status |

#### Database Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `venio_db_connections_in_use` | Gauge | Database connections in use | - |
| `venio_db_connections_idle` | Gauge | Idle database connections | - |
| `venio_db_queries_total` | Counter | Total database queries | operation, status |
| `venio_db_query_duration_seconds` | Histogram | Database query duration | operation |

#### Redis Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `venio_redis_commands_total` | Counter | Total Redis commands | command, status |
| `venio_redis_command_duration_seconds` | Histogram | Redis command duration | command |

#### Authentication Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `venio_auth_attempts_total` | Counter | Authentication attempts | type, status |
| `venio_auth_tokens_issued_total` | Counter | JWT tokens issued | - |

#### Rate Limiting Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `venio_rate_limit_hits_total` | Counter | Rate limit hits | limiter, status |

### Querying Metrics

#### PromQL Examples

```promql
# Request rate (requests per second)
rate(venio_http_requests_total[5m])

# Average request duration
rate(venio_http_request_duration_seconds_sum[5m]) / rate(venio_http_request_duration_seconds_count[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(venio_http_request_duration_seconds_bucket[5m]))

# Error rate (5xx responses)
sum(rate(venio_http_requests_total{status=~"5.."}[5m])) / sum(rate(venio_http_requests_total[5m]))

# Database connection pool utilization
venio_db_connections_in_use / (venio_db_connections_in_use + venio_db_connections_idle)

# Failed authentication rate
rate(venio_auth_attempts_total{status="failure"}[5m])
```

## Health Checks

Venio provides two health check endpoints for Kubernetes-style probes.

### Liveness Probe

**Endpoint:** `GET /health/live`

Checks if the application is running. This endpoint always returns 200 OK if the process is alive.

```bash
curl http://localhost:3690/health/live
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-01-15T10:30:00Z"
}
```

### Readiness Probe

**Endpoint:** `GET /health/ready`

Checks if the application is ready to serve traffic. Validates connectivity to all required services.

```bash
curl http://localhost:3690/health/ready
```

**Success Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2026-01-15T10:30:00Z",
  "version": "1.0.0",
  "services": {
    "database": {
      "status": "healthy"
    },
    "redis": {
      "status": "healthy"
    }
  }
}
```

**Failure Response (503 Service Unavailable):**
```json
{
  "status": "unhealthy",
  "timestamp": "2026-01-15T10:30:00Z",
  "version": "1.0.0",
  "services": {
    "database": {
      "status": "unhealthy",
      "message": "Database connection failed"
    },
    "redis": {
      "status": "healthy"
    }
  }
}
```

### Kubernetes Configuration

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: venio
    image: venio:latest
    livenessProbe:
      httpGet:
        path: /health/live
        port: 3690
      initialDelaySeconds: 5
      periodSeconds: 10
      timeoutSeconds: 2
      failureThreshold: 3
    readinessProbe:
      httpGet:
        path: /health/ready
        port: 3690
      initialDelaySeconds: 10
      periodSeconds: 5
      timeoutSeconds: 5
      failureThreshold: 3
```

## Logging

Venio uses structured logging with Go's `log/slog` package.

### Log Levels

- **DEBUG:** Detailed diagnostic information (disabled in production)
- **INFO:** General informational messages
- **WARN:** Warning messages for potentially harmful situations
- **ERROR:** Error messages for failures

### Log Format

**Development (text):**
```
time=2026-01-15T10:30:00.000Z level=INFO msg="HTTP request" method=GET path=/api/v1/users status=200 duration=45ms
```

**Production (JSON):**
```json
{
  "time": "2026-01-15T10:30:00.000Z",
  "level": "INFO",
  "msg": "HTTP request",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "duration": "45ms",
  "user_id": 123,
  "ip": "192.168.1.1"
}
```

### Contextual Logging

Use logger helper methods for domain-specific logs:

```go
// HTTP request logging
logger.HTTP().Info("processing request",
    "method", r.Method,
    "path", r.URL.Path)

// Authentication logging
logger.Auth().Warn("failed login attempt",
    "email", email,
    "ip", clientIP)

// Database logging
logger.DB().Error("query failed",
    "operation", "SELECT",
    "error", err)
```

## Dashboards

### Prometheus Configuration

Add Venio to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'venio'
    static_configs:
      - targets: ['venio:3690']
    scrape_interval: 15s
    metrics_path: /metrics
```

### Grafana Dashboards

#### Import Pre-built Dashboard

1. Open Grafana
2. Navigate to **Dashboards > Import**
3. Upload `deployments/grafana/venio-dashboard.json`

#### Key Panels

**System Overview:**
- Request rate (RPS)
- Error rate
- P50/P95/P99 latency
- Active users

**HTTP Metrics:**
- Requests by endpoint
- Response status codes
- Request/response size
- Slowest endpoints

**Database:**
- Connection pool utilization
- Query duration
- Query rate by operation
- Failed queries

**Redis:**
- Command rate
- Command duration
- Rate limit hits
- Cache hit/miss ratio

**Authentication:**
- Login attempts
- Failed authentications
- Tokens issued
- Active sessions

### Example Alerts

```yaml
# Prometheus alerting rules
groups:
  - name: venio
    interval: 30s
    rules:
      - alert: HighErrorRate
        expr: |
          sum(rate(venio_http_requests_total{status=~"5.."}[5m]))
          /
          sum(rate(venio_http_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }}"

      - alert: HighLatency
        expr: |
          histogram_quantile(0.95,
            rate(venio_http_request_duration_seconds_bucket[5m])
          ) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High request latency"
          description: "P95 latency is {{ $value }}s"

      - alert: DatabaseDown
        expr: venio_db_connections_in_use == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database connection lost"

      - alert: HighAuthFailureRate
        expr: |
          rate(venio_auth_attempts_total{status="failure"}[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High authentication failure rate"
          description: "{{ $value }} failed auth attempts per second"
```

## Best Practices

### Metrics

1. **Use labels sparingly** - High cardinality kills performance
2. **Don't include user IDs in labels** - Use low-cardinality dimensions
3. **Prefer histograms for latency** - More accurate than summaries
4. **Set appropriate bucket sizes** - Match your SLOs

### Health Checks

1. **Keep liveness simple** - Just check if process is alive
2. **Readiness should be fast** - 5-second timeout maximum
3. **Don't fail on temporary issues** - Use retries with exponential backoff
4. **Check all dependencies** - Database, Redis, external APIs

### Logging

1. **Use structured logging** - JSON in production
2. **Add context** - Request ID, user ID, correlation ID
3. **Don't log secrets** - Sanitize sensitive data
4. **Log at appropriate levels** - ERROR for failures, INFO for events
5. **Include duration** - Time expensive operations

## Troubleshooting

### Metrics not appearing

1. Check if Prometheus can reach the `/metrics` endpoint
2. Verify firewall rules allow port 3690
3. Check Prometheus configuration and reload
4. View metrics directly: `curl http://localhost:3690/metrics`

### Health checks failing

1. Check service connectivity: `docker compose ps`
2. Verify environment variables are set correctly
3. Check logs: `docker compose logs venio`
4. Test manually: `curl http://localhost:3690/health/ready`

### High memory usage

Prometheus metrics can consume memory with high cardinality labels. Review:
1. Number of unique label combinations
2. Histogram bucket configuration
3. Metric retention period

## References

- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Go Prometheus Client](https://github.com/prometheus/client_golang)
- [Kubernetes Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
