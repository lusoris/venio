// Package handlers contains HTTP request handlers
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/services"
)

// UserRoleHandler handles user-role assignment HTTP requests
type UserRoleHandler struct {
	userRoleService services.UserRoleService
}

// NewUserRoleHandler creates a new user-role handler
func NewUserRoleHandler(userRoleService services.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{
		userRoleService: userRoleService,
	}
}

// GetUserRoles retrieves all roles for a user
func (h *UserRoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roles, err := h.userRoleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": roles})
}

// AssignRoleToUser assigns a role to a user
func (h *UserRoleHandler) AssignRoleToUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.AssignRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	err = h.userRoleService.AssignRole(c.Request.Context(), userID, req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to assign role",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// RemoveRoleFromUser removes a role from a user
func (h *UserRoleHandler) RemoveRoleFromUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = h.userRoleService.RemoveRole(c.Request.Context(), userID, roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to remove role",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
