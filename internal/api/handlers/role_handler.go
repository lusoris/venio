// Package handlers contains HTTP request handlers
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/services"
)

// RoleHandler handles role-related HTTP requests
type RoleHandler struct {
	roleService services.RoleService
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// GetRole retrieves a role by ID
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// ListRoles retrieves a paginated list of roles
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	roles, total, err := h.roleService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"value": roles,
		"count": total,
	})
}

// CreateRole creates a new role
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req models.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	role, err := h.roleService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create role",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": role})
}

// UpdateRole updates an existing role
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req models.UpdateRoleRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": bindErr.Error(),
		})
		return
	}

	role, err := h.roleService.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update role",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// DeleteRole deletes a role
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = h.roleService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete role",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetRolePermissions retrieves all permissions for a role
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	permissions, err := h.roleService.GetPermissions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": permissions})
}

// AssignPermissionToRole assigns a permission to a role
func (h *RoleHandler) AssignPermissionToRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req models.AssignPermissionRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": bindErr.Error(),
		})
		return
	}

	err = h.roleService.AssignPermissionToRole(c.Request.Context(), roleID, req.PermissionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to assign permission",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned successfully"})
}

// RemovePermissionFromRole removes a permission from a role
func (h *RoleHandler) RemovePermissionFromRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	permissionID, err := strconv.ParseInt(c.Param("permissionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	err = h.roleService.RemovePermissionFromRole(c.Request.Context(), roleID, permissionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to remove permission",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
