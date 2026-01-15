// Package middleware contains HTTP middleware functions
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lusoris/venio/internal/logger"
)

// LoggingMiddleware logs HTTP requests using structured logging
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Milliseconds()
		status := c.Writer.Status()

		// Get user context if available
		userID, userExists := c.Get("user_id")
		email, _ := c.Get("email")

		// Log request
		if userExists {
			log.HTTP(method, path, status, duration,
				"user_id", userID,
				"email", email,
				"ip", c.ClientIP(),
			)
		} else {
			log.HTTP(method, path, status, duration,
				"ip", c.ClientIP(),
			)
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Error("Request error", err.Err,
					"method", method,
					"path", path,
					"type", err.Type,
				)
			}
		}
	}
}
