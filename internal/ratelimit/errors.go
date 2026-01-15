package ratelimit

import "errors"

var (
	// ErrRateLimitExceeded is returned when the rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrInvalidMaxRequests is returned when MaxRequests is invalid
	ErrInvalidMaxRequests = errors.New("max requests must be greater than 0")

	// ErrInvalidWindow is returned when Window is invalid
	ErrInvalidWindow = errors.New("window must be greater than 0")

	// ErrRedisConnectionRequired is returned when Redis is required but not configured
	ErrRedisConnectionRequired = errors.New("redis connection required for distributed rate limiting")
)
