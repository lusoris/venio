// Package handlers contains HTTP request handlers
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/services"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	// Don't expose password
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// ListUsers retrieves a paginated list of users
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.User
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	limit := 10
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(c.DefaultQuery("offset", "0")); err == nil {
		offset = o
	}

	users, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list users",
			Message: err.Error(),
		})
		return
	}

	// Don't expose passwords
	for _, user := range users {
		user.Password = ""
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser updates a user
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "Update request"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	var req models.UpdateUserRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: bindErr.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Update failed",
			Message: err.Error(),
		})
		return
	}

	// Don't expose password
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Delete failed",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
