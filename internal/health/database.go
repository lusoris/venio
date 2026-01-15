package health

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseChecker checks database health
type DatabaseChecker struct {
	pool *pgxpool.Pool
	name string
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(pool *pgxpool.Pool) *DatabaseChecker {
	return &DatabaseChecker{
		pool: pool,
		name: "database",
	}
}

// Name returns the checker name
func (d *DatabaseChecker) Name() string {
	return d.name
}

// Check performs the database health check
func (d *DatabaseChecker) Check(ctx context.Context) Check {
	start := time.Now()

	check := Check{
		Name:      d.name,
		Timestamp: start,
	}

	// Ping database
	if err := d.pool.Ping(ctx); err != nil {
		check.Status = StatusUnhealthy
		check.Message = "Database ping failed: " + err.Error()
		check.ResponseTime = time.Since(start)
		return check
	}

	// Get pool stats
	stats := d.pool.Stat()
	check.Metadata = map[string]interface{}{
		"total_connections":        stats.TotalConns(),
		"idle_connections":         stats.IdleConns(),
		"acquired_connections":     stats.AcquiredConns(),
		"max_connections":          stats.MaxConns(),
		"constructing_connections": stats.ConstructingConns(),
	}

	// Check if pool is healthy
	if stats.TotalConns() >= stats.MaxConns() {
		check.Status = StatusDegraded
		check.Message = "Database connection pool at maximum capacity"
	} else {
		check.Status = StatusHealthy
		check.Message = "Database connection successful"
	}

	check.ResponseTime = time.Since(start)
	return check
}
