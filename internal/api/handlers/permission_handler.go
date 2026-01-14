// Package handlers contains HTTP request handlers
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/services"
)

// PermissionHandler handles permission-related HTTP requests
type PermissionHandler struct {
	permissionService services.PermissionService
}

// NewPermissionHandler creates a new permission handler
func NewPermissionHandler(permissionService services.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// GetPermission retrieves a permission by ID
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	permission, err := h.permissionService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permission})
}

// ListPermissions retrieves a paginated list of permissions
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
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

	permissions, total, err := h.permissionService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"value": permissions,
		"count": total,
	})
}

// CreatePermission creates a new permission
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req models.CreatePermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	permission, err := h.permissionService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create permission",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": permission})
}

// UpdatePermission updates an existing permission
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	var req models.UpdatePermissionRequest

	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": bindErr.Error(),
		})
		return
	}

	permission, err := h.permissionService.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update permission",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permission})
}

// DeletePermission deletes a permission
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	err = h.permissionService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete permission",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
