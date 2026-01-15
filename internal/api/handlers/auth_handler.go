// Package handlers contains HTTP request handlers
package handlers

import (
	"errors"
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

// SuccessResponse represents a generic success message
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required,min=64,max=128" example:"6f0a..."`
}

// ResendVerificationRequest represents resend verification request
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email,max=255" example:"user@example.com"`
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
		// Return generic message to client (detailed error logged in service layer)
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
		// Return generic message to client (detailed error logged in service layer)
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

// VerifyEmail handles email verification
// @Summary Verify email
// @Description Verify a user's email using a verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Verification token"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Invalid verification token",
		})
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		switch {
		case errors.Is(err, services.ErrVerificationTokenExpired):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid token",
				Message: "Verification token has expired. Please request a new email.",
			})
		case errors.Is(err, services.ErrEmailAlreadyVerified):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Email already verified",
				Message: "This email address has already been verified.",
			})
		case errors.Is(err, services.ErrInvalidVerificationToken):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid token",
				Message: "Verification token is invalid or has expired.",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Unable to verify email at this time.",
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Email verified successfully"})
}

// ResendVerificationEmail handles resending verification emails
// @Summary Resend verification email
// @Description Resend an email verification link to the specified email address
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResendVerificationRequest true "Email address"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/resend-verification [post]
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Please provide a valid email address",
		})
		return
	}

	if err := h.authService.ResendVerificationEmail(c.Request.Context(), req.Email); err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyVerified):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Email already verified",
				Message: "This email address has already been verified.",
			})
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "User not found",
				Message: "No account found for the provided email.",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Unable to resend verification email at this time.",
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Verification email sent"})
}

// RefreshTokenResponse represents a token refresh response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}
