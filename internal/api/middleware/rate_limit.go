// Package middleware contains HTTP middleware functions
package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter stores rate limit configuration
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	maxReqs  int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
// maxReqs: maximum number of requests allowed
// window: time window for counting requests
func NewRateLimiter(maxReqs int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		maxReqs:  maxReqs,
		window:   window,
	}

	// Cleanup old requests every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// Middleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
				"message": fmt.Sprintf(
					"Rate limit exceeded: %d requests per %v",
					rl.maxReqs, rl.window,
				),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow checks if the client IP is within rate limit
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get or create request list for this IP
	requests, exists := rl.requests[ip]
	if !exists {
		requests = []time.Time{now}
		rl.requests[ip] = requests
		return true
	}

	// Remove requests outside the window
	validRequests := []time.Time{}
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if within limit
	if len(validRequests) >= rl.maxReqs {
		rl.requests[ip] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests
	return true
}

// cleanup removes old entries
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for ip, requests := range rl.requests {
		validRequests := []time.Time{}
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = validRequests
		}
	}
}

// AuthRateLimiter is a rate limiter specifically for authentication endpoints
// Default: 5 attempts per minute
func AuthRateLimiter() *RateLimiter {
	return NewRateLimiter(5, 1*time.Minute)
}

// GeneralRateLimiter is a rate limiter for general API endpoints
// Default: 100 requests per minute
func GeneralRateLimiter() *RateLimiter {
	return NewRateLimiter(100, 1*time.Minute)
}
