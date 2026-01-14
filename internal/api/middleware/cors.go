// Package middleware contains HTTP middleware functions
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware configured for the application
// AllowOrigins: Only allow requests from the frontend domain (not *)
// AllowMethods: GET, POST, PUT, DELETE, OPTIONS
// AllowHeaders: Authorization, Content-Type, Accept, Origin
// AllowCredentials: true (for cookies/auth headers)
// MaxAge: 12 hours
func CORS(frontendURL string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "Origin", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           43200, // 12 hours in seconds
	})
}

// CORSDevelopment returns a CORS middleware for development (allows all origins)
// ONLY use this in development mode!
func CORSDevelopment() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "Origin", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           43200,
	})
}
