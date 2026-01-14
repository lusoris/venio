package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/services"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	userService       services.UserService
	roleService       services.RoleService
	permissionService services.PermissionService
	userRoleService   services.UserRoleService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	userService services.UserService,
	roleService services.RoleService,
	permissionService services.PermissionService,
	userRoleService services.UserRoleService,
) *AdminHandler {
	return &AdminHandler{
		userService:       userService,
		roleService:       roleService,
		permissionService: permissionService,
		userRoleService:   userRoleService,
	}
}

// ListUsers lists all users (admin only)
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.ListUsers(c.Request.Context(), 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Format response
	var usersData []gin.H
	for _, user := range users {
		usersData = append(usersData, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": usersData,
	})
}

// CreateUser creates a new user with roles
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email     string  `json:"email" binding:"required,email"`
		Username  string  `json:"username" binding:"required,min=3,max=50"`
		FirstName string  `json:"first_name" binding:"required"`
		LastName  string  `json:"last_name" binding:"required"`
		Password  string  `json:"password" binding:"required,min=8"`
		Roles     []int64 `json:"roles"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userReq := &models.CreateUserRequest{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	// Create user
	createdUser, err := h.userService.Register(c.Request.Context(), userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Assign roles if provided
	for _, roleID := range req.Roles {
		if err := h.userRoleService.AssignRole(c.Request.Context(), createdUser.ID, roleID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        createdUser.ID,
		"email":     createdUser.Email,
		"username":  createdUser.Username,
		"firstName": createdUser.FirstName,
		"lastName":  createdUser.LastName,
	})
}

// DeleteUser deletes a user (admin only)
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListRoles lists all roles (admin only)
func (h *AdminHandler) ListRoles(c *gin.Context) {
	roles, _, err := h.roleService.List(c.Request.Context(), 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}

	var rolesData []gin.H
	for _, role := range roles {
		// Get user count for this role
		userCount := 0 // TODO: Implement counting
		rolesData = append(rolesData, gin.H{
			"id":          role.ID,
			"name":        role.Name,
			"description": role.Description,
			"user_count":  userCount,
			"created_at":  role.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": rolesData,
	})
}

// CreateRole creates a new role
func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Permissions []int64 `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleReq := models.CreateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	createdRole, err := h.roleService.Create(c.Request.Context(), roleReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	// Assign permissions if provided
	for _, permID := range req.Permissions {
		if err := h.roleService.AssignPermissionToRole(c.Request.Context(), createdRole.ID, permID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   createdRole.ID,
		"name": createdRole.Name,
	})
}

// DeleteRole deletes a role
func (h *AdminHandler) DeleteRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = h.roleService.Delete(c.Request.Context(), roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// ListPermissions lists all permissions
func (h *AdminHandler) ListPermissions(c *gin.Context) {
	permissions, _, err := h.permissionService.List(c.Request.Context(), 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}

	var permsData []gin.H
	for _, perm := range permissions {
		permsData = append(permsData, gin.H{
			"id":          perm.ID,
			"name":        perm.Name,
			"description": perm.Description,
			"created_at":  perm.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permsData,
	})
}

// ListUserRoles lists all user-role assignments
func (h *AdminHandler) ListUserRoles(c *gin.Context) {
	users, err := h.userService.ListUsers(c.Request.Context(), 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	var assignments []gin.H
	for _, user := range users {
		userRoles, err := h.userRoleService.GetUserRoles(c.Request.Context(), user.ID)
		if err != nil {
			continue
		}

		for _, roleName := range userRoles {
			assignments = append(assignments, gin.H{
				"user_id":     user.ID,
				"user_email":  user.Email,
				"role_name":   roleName,
				"assigned_at": user.CreatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
	})
}

// RemoveUserRole removes a user from a role
func (h *AdminHandler) RemoveUserRole(c *gin.Context) {
	var req struct {
		UserID int64 `json:"user_id" binding:"required"`
		RoleID int64 `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userRoleService.RemoveRole(c.Request.Context(), req.UserID, req.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assignment removed successfully"})
}
