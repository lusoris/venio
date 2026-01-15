// Package middleware contains HTTP middleware based on the new ratelimit abstractions
package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/ratelimit"
)

// RateLimitMiddleware creates a rate limiting middleware using the limiter interface
func RateLimitMiddleware(limiter ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		allowed, err := limiter.Allow(c.Request.Context(), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": "Rate limit check failed",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddlewareWithCustomKey creates a rate limiting middleware with custom key extractor
func RateLimitMiddlewareWithCustomKey(limiter ratelimit.Limiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)

		allowed, err := limiter.Allow(c.Request.Context(), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": "Rate limit check failed",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimitKey extracts user ID from context for rate limiting
func UserRateLimitKey(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return c.ClientIP()
	}
	return fmt.Sprintf("user:%v", userID)
}

// EndpointRateLimitKey creates a rate limit key based on endpoint and user
func EndpointRateLimitKey(endpoint string) func(*gin.Context) string {
	return func(c *gin.Context) string {
		userID, exists := c.Get("user_id")
		if exists {
			return fmt.Sprintf("endpoint:%s:user:%v", endpoint, userID)
		}
		return fmt.Sprintf("endpoint:%s:ip:%s", endpoint, c.ClientIP())
	}
}
