// Package handlers contains HTTP request handlers
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/services"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"Invalid input data"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	User         *models.User `json:"user"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService services.AuthService
	userService services.UserService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService services.AuthService, userService services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "Registration request"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Please check your input and try again",
		})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		// Log detailed error server-side
		c.Request.Context().Value("logger") // TODO: Add proper logger

		// Return generic message to client
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Registration failed",
			Message: "Unable to create account. Email may already be registered.",
		})
		return
	}

	// Don't expose password
	user.Password = ""

	c.JSON(http.StatusCreated, user)
}

// Login handles user login
// @Summary Login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Please provide valid email and password",
		})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Authentication failed",
			Message: "Invalid email or password",
		})
		return
	}

	// Get user info
	user, err := h.userService.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// Log detailed error server-side
		// TODO: Use logger from context

		// Generic message to client
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Authentication failed",
			Message: "Unable to process login request",
		})
		return
	}

	// Don't expose password
	user.Password = ""

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Invalid refresh token format",
		})
		return
	}

	accessToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token refresh failed",
			Message: "Refresh token expired or invalid. Please login again.",
		})
		return
	}

	c.JSON(http.StatusOK, RefreshTokenResponse{
		AccessToken: accessToken,
	})
}

// RefreshTokenResponse represents a token refresh response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}
