// Package health provides abstraction for health checks
package health

import (
	"context"
	"time"
)

// Status represents the health status of a component
type Status string

const (
	// StatusHealthy indicates the component is healthy
	StatusHealthy Status = "healthy"

	// StatusUnhealthy indicates the component is unhealthy
	StatusUnhealthy Status = "unhealthy"

	// StatusDegraded indicates the component is degraded but functional
	StatusDegraded Status = "degraded"

	// StatusUnknown indicates the health status is unknown
	StatusUnknown Status = "unknown"
)

// Check represents a health check result
type Check struct {
	// Name is the name of the check
	Name string

	// Status is the health status
	Status Status

	// Message provides additional information
	Message string

	// Timestamp is when the check was performed
	Timestamp time.Time

	// ResponseTime is how long the check took
	ResponseTime time.Duration

	// Metadata contains additional check-specific data
	Metadata map[string]interface{}
}

// Checker defines the interface for health checks
type Checker interface {
	// Check performs a health check
	Check(ctx context.Context) Check

	// Name returns the name of the checker
	Name() string
}

// Result represents the overall health check result
type Result struct {
	// Status is the overall health status
	Status Status

	// Timestamp is when the check was performed
	Timestamp time.Time

	// Version is the application version
	Version string

	// Checks contains individual check results
	Checks []Check
}

// IsHealthy returns true if all checks are healthy
func (r *Result) IsHealthy() bool {
	if r.Status == StatusHealthy {
		return true
	}

	for _, check := range r.Checks {
		if check.Status != StatusHealthy {
			return false
		}
	}

	return true
}

// Manager manages multiple health checkers
type Manager struct {
	checkers []Checker
	version  string
}

// NewManager creates a new health check manager
func NewManager(version string) *Manager {
	return &Manager{
		checkers: make([]Checker, 0),
		version:  version,
	}
}

// Register registers a new health checker
func (m *Manager) Register(checker Checker) {
	m.checkers = append(m.checkers, checker)
}

// CheckAll performs all registered health checks
func (m *Manager) CheckAll(ctx context.Context) Result {
	result := Result{
		Timestamp: time.Now().UTC(),
		Version:   m.version,
		Checks:    make([]Check, 0, len(m.checkers)),
	}

	allHealthy := true
	for _, checker := range m.checkers {
		check := checker.Check(ctx)
		result.Checks = append(result.Checks, check)

		if check.Status != StatusHealthy {
			allHealthy = false
		}
	}

	if allHealthy {
		result.Status = StatusHealthy
	} else {
		result.Status = StatusUnhealthy
	}

	return result
}

// CheckOne performs a specific health check by name
func (m *Manager) CheckOne(ctx context.Context, name string) (Check, bool) {
	for _, checker := range m.checkers {
		if checker.Name() == name {
			return checker.Check(ctx), true
		}
	}

	return Check{
		Name:      name,
		Status:    StatusUnknown,
		Message:   "Checker not found",
		Timestamp: time.Now().UTC(),
	}, false
}
