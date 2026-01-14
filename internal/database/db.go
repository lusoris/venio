// Package database handles database connections and operations
package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/venio/internal/config"
)

// DB represents a database connection pool
type DB struct {
	pool *pgxpool.Pool
}

// Connect establishes a connection to PostgreSQL
func Connect(ctx context.Context, cfg *config.DatabaseConfig) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("invalid database connection string: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✓ Database connection established")
	return &DB{pool: pool}, nil
}

// Close closes the database connection pool
func (d *DB) Close() error {
	d.pool.Close()
	return nil
}

// Pool returns the underlying pgxpool connection pool
func (d *DB) Pool() *pgxpool.Pool {
	return d.pool
}

// Query executes a query that returns rows
func (d *DB) Query(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	return d.pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that returns a single row
func (d *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) interface{} {
	return d.pool.QueryRow(ctx, sql, args...)
}

// Exec executes a command
func (d *DB) Exec(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	return d.pool.Exec(ctx, sql, args...)
}
