// Package middleware contains HTTP middleware functions
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/services"
)

// RBACMiddleware provides role-based access control
type RBACMiddleware struct {
	userRoleService services.UserRoleService
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(userRoleService services.UserRoleService) *RBACMiddleware {
	return &RBACMiddleware{
		userRoleService: userRoleService,
	}
}

// RequireRole is middleware that checks if user has the specified role
func (m *RBACMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userIDInt, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			c.Abort()
			return
		}

		// Check if user has the required role
		hasRole, err := m.userRoleService.HasRole(c.Request.Context(), userIDInt, requiredRole)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role"})
			c.Abort()
			return
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "User does not have required role"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission is middleware that checks if user has the specified permission
func (m *RBACMiddleware) RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userIDInt, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			c.Abort()
			return
		}

		// Check if user has the required permission
		hasPermission, err := m.userRoleService.HasPermission(c.Request.Context(), userIDInt, requiredPermission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "User does not have required permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole is middleware that checks if user has any of the specified roles
func (m *RBACMiddleware) RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userIDInt, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasAnyRole := false
		for _, role := range requiredRoles {
			has, err := m.userRoleService.HasRole(c.Request.Context(), userIDInt, role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role"})
				c.Abort()
				return
			}

			if has {
				hasAnyRole = true
				break
			}
		}

		if !hasAnyRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "User does not have any of the required roles"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission is middleware that checks if user has any of the specified permissions
func (m *RBACMiddleware) RequireAnyPermission(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userIDInt, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			c.Abort()
			return
		}

		// Check if user has any of the required permissions
		hasAnyPermission := false
		for _, permission := range requiredPermissions {
			has, err := m.userRoleService.HasPermission(c.Request.Context(), userIDInt, permission)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
				c.Abort()
				return
			}

			if has {
				hasAnyPermission = true
				break
			}
		}

		if !hasAnyPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "User does not have any of the required permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
