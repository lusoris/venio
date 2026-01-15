// Package models contains data models for the application
package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User validation errors
var (
	ErrEmptyEmail      = errors.New("email is required")
	ErrInvalidUsername = errors.New("username must be between 3 and 50 characters")
	ErrEmptyFirstName  = errors.New("first name is required")
	ErrEmptyLastName   = errors.New("last name is required")
	ErrWeakPassword    = errors.New("password must be at least 8 characters")
)

// User represents a user in the system
type User struct {
	ID        int64   `json:"id" example:"1"`
	Email     string  `json:"email" example:"user@example.com"`
	Username  string  `json:"username" example:"johndoe"`
	FirstName string  `json:"first_name" example:"John"`
	LastName  string  `json:"last_name" example:"Doe"`
	Avatar    *string `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"`
	Password  string  `json:"-"` // Never expose password
	IsActive  bool    `json:"is_active" example:"true"`

	// Email verification fields
	IsEmailVerified              bool       `json:"is_email_verified" example:"false"`
	EmailVerificationToken       *string    `json:"-"` // Never expose token
	EmailVerificationTokenExpiry *time.Time `json:"-"` // Never expose expiry
	EmailVerifiedAt              *time.Time `json:"email_verified_at,omitempty" example:"2026-01-15T10:30:00Z"`

	CreatedAt time.Time `json:"created_at" example:"2026-01-15T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-01-15T10:30:00Z"`
}

// Role represents a user role
type Role struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"` // admin, moderator, user, etc.
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Permission represents a permission
type Permission struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"` // e.g., "users.read", "users.write"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRole represents the junction between users and roles
type UserRole struct {
	UserID     int64     `json:"user_id"`
	RoleID     int64     `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

// RolePermission represents the junction between roles and permissions
type RolePermission struct {
	RoleID       int64     `json:"role_id"`
	PermissionID int64     `json:"permission_id"`
	AssignedAt   time.Time `json:"assigned_at"`
}

// CreateUserRequest is the request body for creating a user
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email,max=255" example:"user@example.com"`
	Username  string  `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	FirstName string  `json:"first_name" binding:"required,max=100" example:"John"`
	LastName  string  `json:"last_name" binding:"required,max=100" example:"Doe"`
	Avatar    *string `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"`
	Password  string  `json:"password" binding:"required,min=8,max=128" example:"SecurePass123!"`
}

// Validate checks if the CreateUserRequest is valid
func (r *CreateUserRequest) Validate() error {
	if r.Email == "" {
		return ErrEmptyEmail
	}
	if r.Username == "" || len(r.Username) < 3 || len(r.Username) > 50 {
		return ErrInvalidUsername
	}
	if r.FirstName == "" {
		return ErrEmptyFirstName
	}
	if r.LastName == "" {
		return ErrEmptyLastName
	}
	if r.Password == "" || len(r.Password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

// UpdateUserRequest is the request body for updating a user
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=100"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=100"`
	Avatar    *string `json:"avatar,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

// LoginResponse is the response for login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID   int64    `json:"user_id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// CreateRoleRequest is the request body for creating a role
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"required,min=10,max=255"`
}

// UpdateRoleRequest is the request body for updating a role
type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=50"`
	Description *string `json:"description,omitempty" binding:"omitempty,min=10,max=255"`
}

// CreatePermissionRequest is the request body for creating a permission
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required,min=10,max=255"`
}

// UpdatePermissionRequest is the request body for updating a permission
type UpdatePermissionRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,min=10,max=255"`
}

// AssignRoleRequest is the request body for assigning a role to a user
type AssignRoleRequest struct {
	RoleID int64 `json:"role_id" binding:"required"`
}

// AssignPermissionRequest is the request body for assigning a permission to a role
type AssignPermissionRequest struct {
	PermissionID int64 `json:"permission_id" binding:"required"`
}
