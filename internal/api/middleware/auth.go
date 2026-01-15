// Package middleware contains HTTP middleware functions
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/services"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(ErrMissingAuthHeader.StatusCode(), ErrMissingAuthHeader)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(ErrInvalidAuthFormat.StatusCode(), ErrInvalidAuthFormat)
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(ErrInvalidToken.StatusCode(), ErrInvalidToken)
			c.Abort()
			return
		}

		// Store claims in context using type-safe helpers
		SetUserID(c, claims.UserID)
		SetEmail(c, claims.Email)
		SetUsername(c, claims.Username)
		SetRoles(c, claims.Roles)

		c.Next()
	}
}
