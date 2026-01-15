// Package services contains business logic implementations
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lusoris/venio/internal/config"
	"github.com/lusoris/venio/internal/models"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	ValidateToken(tokenString string) (*models.TokenClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

// DefaultAuthService implements AuthService
type DefaultAuthService struct {
	userService     UserService
	userRoleService UserRoleService
	config          *config.Config
}

// NewDefaultAuthService creates a new auth service
func NewDefaultAuthService(userService UserService, userRoleService UserRoleService, cfg *config.Config) AuthService {
	return &DefaultAuthService{
		userService:     userService,
		userRoleService: userRoleService,
		config:          cfg,
	}
}

// Login authenticates a user and returns access and refresh tokens
func (s *DefaultAuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	// Add timeout to existing context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("authentication failed: %w", err)
	}

	// Verify password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return "", "", errors.New("user account is inactive")
	}

	// Get user roles for JWT
	roles, err := s.userRoleService.GetUserRoles(ctx, user.ID)
	if err != nil {
		// Log error but don't fail login if roles can't be fetched
		// User will have empty roles array
		roles = []string{}
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user, roles)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user, roles)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates and parses a JWT token
func (s *DefaultAuthService) ValidateToken(tokenString string) (*models.TokenClaims, error) {
	claims := &models.TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *DefaultAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token validation failed: %w", err)
	}

	// Add timeout to existing context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	user, err := s.userService.GetUser(ctx, claims.UserID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return "", errors.New("user account is inactive")
	}

	// Get fresh roles for new token
	roles, err := s.userRoleService.GetUserRoles(ctx, user.ID)
	if err != nil {
		roles = []string{} // Fail gracefully
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user, roles)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// generateAccessToken creates a new JWT access token
func (s *DefaultAuthService) generateAccessToken(user *models.User, roles []string) (string, error) {
	now := time.Now()

	claims := &models.TokenClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWT.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "venio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("token signing failed: %w", err)
	}

	return tokenString, nil
}

// generateRefreshToken creates a new JWT refresh token with longer expiration
func (s *DefaultAuthService) generateRefreshToken(user *models.User, roles []string) (string, error) {
	now := time.Now()
	refreshDays := time.Duration(s.config.JWT.RefreshExpiryDays) * 24 * time.Hour

	claims := &models.TokenClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshDays)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "venio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("token signing failed: %w", err)
	}

	return tokenString, nil
}
