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
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    *string   `json:"avatar,omitempty"`
	Password  string    `json:"-"` // Never expose password
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Email     string  `json:"email" binding:"required,email"`
	Username  string  `json:"username" binding:"required,min=3,max=50"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Avatar    *string `json:"avatar,omitempty"`
	Password  string  `json:"password" binding:"required,min=8"`
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
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
