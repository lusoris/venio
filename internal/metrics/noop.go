package metrics

import "time"

// NoOpCollector is a metrics collector that does nothing
// Useful for testing or when metrics are disabled
type NoOpCollector struct{}

// NewNoOpCollector creates a new no-op metrics collector
func NewNoOpCollector() *NoOpCollector {
	return &NoOpCollector{}
}

// RecordHTTPRequest does nothing
func (n *NoOpCollector) RecordHTTPRequest(method, path, status string, duration time.Duration, requestSize, responseSize int64) {
}

// RecordDBQuery does nothing
func (n *NoOpCollector) RecordDBQuery(operation, status string, duration time.Duration) {
}

// UpdateDBConnections does nothing
func (n *NoOpCollector) UpdateDBConnections(inUse, idle int) {
}

// RecordRedisCommand does nothing
func (n *NoOpCollector) RecordRedisCommand(command, status string, duration time.Duration) {
}

// RecordAuthAttempt does nothing
func (n *NoOpCollector) RecordAuthAttempt(authType, status string) {
}

// RecordTokenIssued does nothing
func (n *NoOpCollector) RecordTokenIssued() {
}

// RecordRateLimitHit does nothing
func (n *NoOpCollector) RecordRateLimitHit(limiter, status string) {
}

// IncCounter does nothing
func (n *NoOpCollector) IncCounter(name string, labels map[string]string, value float64) {
}

// ObserveHistogram does nothing
func (n *NoOpCollector) ObserveHistogram(name string, labels map[string]string, value float64) {
}

// SetGauge does nothing
func (n *NoOpCollector) SetGauge(name string, labels map[string]string, value float64) {
}
