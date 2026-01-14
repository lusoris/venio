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
	Login(email, password string) (string, string, error)
	ValidateToken(tokenString string) (*models.TokenClaims, error)
	RefreshToken(refreshToken string) (string, error)
}

// DefaultAuthService implements AuthService
type DefaultAuthService struct {
	userService UserService
	config      *config.Config
}

// NewDefaultAuthService creates a new auth service
func NewDefaultAuthService(userService UserService, cfg *config.Config) AuthService {
	return &DefaultAuthService{
		userService: userService,
		config:      cfg,
	}
}

// Login authenticates a user and returns access and refresh tokens
func (s *DefaultAuthService) Login(email, password string) (string, string, error) {
	// Get user by email
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
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
func (s *DefaultAuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token validation failed: %w", err)
	}

	// Get user to verify they still exist and are active
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := s.userService.GetUser(ctx, claims.UserID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return "", errors.New("user account is inactive")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// generateAccessToken creates a new JWT access token
func (s *DefaultAuthService) generateAccessToken(user *models.User) (string, error) {
	now := time.Now()

	claims := &models.TokenClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Roles:    []string{}, // TODO: Populate from user roles
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
func (s *DefaultAuthService) generateRefreshToken(user *models.User) (string, error) {
	now := time.Now()
	refreshDays := time.Duration(s.config.JWT.RefreshExpiryDays) * 24 * time.Hour

	claims := &models.TokenClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Roles:    []string{},
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
