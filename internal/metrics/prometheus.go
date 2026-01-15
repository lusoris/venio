package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusCollector implements Collector using Prometheus
type PrometheusCollector struct {
	config *Config

	// HTTP metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestSize     *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec

	// Database metrics
	dbConnectionsInUse prometheus.Gauge
	dbConnectionsIdle  prometheus.Gauge
	dbQueryDuration    *prometheus.HistogramVec
	dbQueriesTotal     *prometheus.CounterVec

	// Redis metrics
	redisCommandsTotal   *prometheus.CounterVec
	redisCommandDuration *prometheus.HistogramVec

	// Auth metrics
	authAttemptsTotal *prometheus.CounterVec
	authTokensIssued  prometheus.Counter

	// Rate limit metrics
	rateLimitHits *prometheus.CounterVec
}

// NewPrometheusCollector creates a new Prometheus metrics collector
func NewPrometheusCollector(config *Config) (*PrometheusCollector, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	pc := &PrometheusCollector{
		config: config,
	}

	// Initialize HTTP metrics
	pc.httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	pc.httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latency in seconds",
			Buckets:   config.HTTPBuckets,
		},
		[]string{"method", "path", "status"},
	)

	pc.httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "http_request_size_bytes",
			Help:      "HTTP request size in bytes",
			Buckets:   prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	pc.httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "http_response_size_bytes",
			Help:      "HTTP response size in bytes",
			Buckets:   prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path", "status"},
	)

	// Initialize Database metrics
	pc.dbConnectionsInUse = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "db_connections_in_use",
			Help:      "Number of database connections currently in use",
		},
	)

	pc.dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "db_connections_idle",
			Help:      "Number of idle database connections",
		},
	)

	pc.dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "db_query_duration_seconds",
			Help:      "Database query duration in seconds",
			Buckets:   config.DBBuckets,
		},
		[]string{"operation"},
	)

	pc.dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "db_queries_total",
			Help:      "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	// Initialize Redis metrics
	pc.redisCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "redis_commands_total",
			Help:      "Total number of Redis commands",
		},
		[]string{"command", "status"},
	)

	pc.redisCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "redis_command_duration_seconds",
			Help:      "Redis command duration in seconds",
			Buckets:   config.RedisBuckets,
		},
		[]string{"command"},
	)

	// Initialize Auth metrics
	pc.authAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "auth_attempts_total",
			Help:      "Total number of authentication attempts",
		},
		[]string{"type", "status"},
	)

	pc.authTokensIssued = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "auth_tokens_issued_total",
			Help:      "Total number of JWT tokens issued",
		},
	)

	// Initialize Rate limit metrics
	pc.rateLimitHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      "rate_limit_hits_total",
			Help:      "Total number of rate limit hits",
		},
		[]string{"limiter", "status"},
	)

	return pc, nil
}

// RecordHTTPRequest records HTTP request metrics
func (pc *PrometheusCollector) RecordHTTPRequest(method, path, status string, duration time.Duration, requestSize, responseSize int64) {
	pc.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	pc.httpRequestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())

	if requestSize > 0 {
		pc.httpRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	}
	if responseSize > 0 {
		pc.httpResponseSize.WithLabelValues(method, path, status).Observe(float64(responseSize))
	}
}

// RecordDBQuery records database query metrics
func (pc *PrometheusCollector) RecordDBQuery(operation, status string, duration time.Duration) {
	pc.dbQueriesTotal.WithLabelValues(operation, status).Inc()
	pc.dbQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// UpdateDBConnections updates database connection metrics
func (pc *PrometheusCollector) UpdateDBConnections(inUse, idle int) {
	pc.dbConnectionsInUse.Set(float64(inUse))
	pc.dbConnectionsIdle.Set(float64(idle))
}

// RecordRedisCommand records Redis command metrics
func (pc *PrometheusCollector) RecordRedisCommand(command, status string, duration time.Duration) {
	pc.redisCommandsTotal.WithLabelValues(command, status).Inc()
	pc.redisCommandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

// RecordAuthAttempt records authentication attempt metrics
func (pc *PrometheusCollector) RecordAuthAttempt(authType, status string) {
	pc.authAttemptsTotal.WithLabelValues(authType, status).Inc()
}

// RecordTokenIssued records token issuance
func (pc *PrometheusCollector) RecordTokenIssued() {
	pc.authTokensIssued.Inc()
}

// RecordRateLimitHit records rate limit hit metrics
func (pc *PrometheusCollector) RecordRateLimitHit(limiter, status string) {
	pc.rateLimitHits.WithLabelValues(limiter, status).Inc()
}

// IncCounter increments a custom counter
func (pc *PrometheusCollector) IncCounter(name string, labels map[string]string, value float64) {
	// Implementation for custom counters
	// Note: This requires dynamic metric registration which is advanced
	// For now, this is a placeholder
}

// ObserveHistogram observes a value in a custom histogram
func (pc *PrometheusCollector) ObserveHistogram(name string, labels map[string]string, value float64) {
	// Implementation for custom histograms
	// Note: This requires dynamic metric registration which is advanced
	// For now, this is a placeholder
}

// SetGauge sets a custom gauge value
func (pc *PrometheusCollector) SetGauge(name string, labels map[string]string, value float64) {
	// Implementation for custom gauges
	// Note: This requires dynamic metric registration which is advanced
	// For now, this is a placeholder
}
