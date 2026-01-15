// Package logger provides structured logging for the application
package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/lusoris/venio/internal/config"
)

// Logger wraps slog.Logger
type Logger struct {
	*slog.Logger
}

// New creates a new structured logger based on configuration
func New(cfg *config.AppConfig) *Logger {
	var handler slog.Handler
	var level slog.Level

	// Set log level based on environment
	if cfg.Debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.Debug, // Add source file/line in debug mode
	}

	// Use JSON handler for production, Text for development
	var output io.Writer = os.Stdout
	if cfg.Env == "production" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	return &Logger{Logger: logger}
}

// WithContext returns a logger with context fields
func (l *Logger) WithContext(attrs ...slog.Attr) *Logger {
	return &Logger{
		Logger: l.Logger.With(attrsToAny(attrs)...),
	}
}

// attrsToAny converts []slog.Attr to []any for With() method
func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, attr := range attrs {
		result[i] = attr
	}
	return result
}

// Common helper methods with structured fields

// Info logs an info message with structured fields
func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.Logger.Info(msg, keysAndValues...)
}

// Error logs an error message with structured fields
func (l *Logger) Error(msg string, err error, keysAndValues ...any) {
	args := append([]any{"error", err}, keysAndValues...)
	l.Logger.Error(msg, args...)
}

// Warn logs a warning message with structured fields
func (l *Logger) Warn(msg string, keysAndValues ...any) {
	l.Logger.Warn(msg, keysAndValues...)
}

// Debug logs a debug message with structured fields
func (l *Logger) Debug(msg string, keysAndValues ...any) {
	l.Logger.Debug(msg, keysAndValues...)
}

// HTTP logs HTTP request information
func (l *Logger) HTTP(method, path string, status int, duration int64, keysAndValues ...any) {
	args := append([]any{
		"method", method,
		"path", path,
		"status", status,
		"duration_ms", duration,
	}, keysAndValues...)
	l.Info("HTTP request", args...)
}

// Auth logs authentication events
func (l *Logger) Auth(event string, userID int64, email string, success bool, keysAndValues ...any) {
	args := append([]any{
		"event", event,
		"user_id", userID,
		"email", email,
		"success", success,
	}, keysAndValues...)

	if success {
		l.Info("Authentication event", args...)
	} else {
		l.Warn("Authentication failed", args...)
	}
}

// DB logs database operations
func (l *Logger) DB(operation, table string, duration int64, err error, keysAndValues ...any) {
	args := append([]any{
		"operation", operation,
		"table", table,
		"duration_ms", duration,
	}, keysAndValues...)

	if err != nil {
		args = append(args, "error", err)
		l.Error("Database operation failed", err, args...)
	} else {
		l.Debug("Database operation", args...)
	}
}
