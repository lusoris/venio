// Package middleware contains HTTP middleware functions
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers to responses
// Headers added:
// - X-Content-Type-Options: nosniff (prevents MIME type sniffing)
// - X-Frame-Options: DENY (prevents clickjacking)
// - X-XSS-Protection: 1; mode=block (legacy XSS protection)
// - Strict-Transport-Security: HSTS (forces HTTPS in production)
// - Content-Security-Policy: Restricts resource loading
// - Referrer-Policy: no-referrer (privacy protection)
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking attacks
		c.Header("X-Frame-Options", "DENY")

		// Legacy XSS protection (modern browsers ignore this)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Prevent loading resources from external origins
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none';")

		// Control referrer information
		c.Header("Referrer-Policy", "no-referrer")

		// Disable browser features that can be exploited
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// In production, enable HSTS
		// Uncomment when deploying to production with HTTPS:
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		c.Next()
	}
}

// StrictSecurityHeaders is a more strict version for production
// Includes HSTS and stricter CSP
func StrictSecurityHeaders(hstsMaxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking attacks
		c.Header("X-Frame-Options", "DENY")

		// Legacy XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Strict Content Security Policy for production
		c.Header("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self';")

		// Control referrer information
		c.Header("Referrer-Policy", "strict-no-referrer")

		// Disable browser features
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Force HTTPS for a year (production only)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		c.Next()
	}
}
