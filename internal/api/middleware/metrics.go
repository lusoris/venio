package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "venio_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "venio_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "venio_http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B to 10MB
		},
		[]string{"method", "path"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "venio_http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path", "status"},
	)

	// Database metrics
	dbConnectionsInUse = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "venio_db_connections_in_use",
			Help: "Number of database connections currently in use",
		},
	)

	dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "venio_db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "venio_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"operation"},
	)

	dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "venio_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	// Redis metrics
	redisCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "venio_redis_commands_total",
			Help: "Total number of Redis commands",
		},
		[]string{"command", "status"},
	)

	redisCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "venio_redis_command_duration_seconds",
			Help:    "Redis command duration in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1, .25, .5},
		},
		[]string{"command"},
	)

	// Auth metrics
	authAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "venio_auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"type", "status"}, // type: login, refresh; status: success, failure
	)

	authTokensIssued = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "venio_auth_tokens_issued_total",
			Help: "Total number of JWT tokens issued",
		},
	)

	// Rate limit metrics
	rateLimitHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "venio_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"limiter", "status"}, // status: allowed, denied
	)
)

// PrometheusMiddleware records HTTP request metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath() // Template path (e.g., /api/v1/users/:id)
		if path == "" {
			path = c.Request.URL.Path // Fallback for unmatched routes
		}

		// Record request size
		if c.Request.ContentLength > 0 {
			httpRequestSize.WithLabelValues(c.Request.Method, path).Observe(float64(c.Request.ContentLength))
		}

		// Process request
		c.Next()

		// Record metrics after request processing
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
		httpResponseSize.WithLabelValues(c.Request.Method, path, status).Observe(float64(c.Writer.Size()))
	}
}

// RecordDBMetrics updates database connection pool metrics
func RecordDBMetrics(inUse, idle int) {
	dbConnectionsInUse.Set(float64(inUse))
	dbConnectionsIdle.Set(float64(idle))
}

// RecordDBQuery records a database query metric
func RecordDBQuery(operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	dbQueriesTotal.WithLabelValues(operation, status).Inc()
	dbQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordRedisCommand records a Redis command metric
func RecordRedisCommand(command string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	redisCommandsTotal.WithLabelValues(command, status).Inc()
	redisCommandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

// RecordAuthAttempt records an authentication attempt
func RecordAuthAttempt(authType string, success bool) {
	status := "failure"
	if success {
		status = "success"
	}

	authAttemptsTotal.WithLabelValues(authType, status).Inc()
}

// RecordTokenIssued increments the tokens issued counter
func RecordTokenIssued() {
	authTokensIssued.Inc()
}

// RecordRateLimitHit records a rate limit event
func RecordRateLimitHit(limiter string, allowed bool) {
	status := "denied"
	if allowed {
		status = "allowed"
	}

	rateLimitHits.WithLabelValues(limiter, status).Inc()
}
