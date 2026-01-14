// Package services contains business logic implementations
package services

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/repositories"
)

// UserService defines the interface for user business logic
type UserService interface {
	Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit int, offset int) ([]*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// DefaultUserService implements UserService
type DefaultUserService struct {
	repo repositories.UserRepository
}

// NewDefaultUserService creates a new user service
func NewDefaultUserService(repo repositories.UserRepository) UserService {
	return &DefaultUserService{repo: repo}
}

// Register creates a new user with validation
func (s *DefaultUserService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if email already exists
	exists, err := s.repo.Exists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Avatar:    req.Avatar,
		Password:  string(hashedPassword),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = id
	return user, nil
}

// GetUser retrieves a user by ID
func (s *DefaultUserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser modifies a user
func (s *DefaultUserService) UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields if provided
	if req.Email != nil {
		if !isValidEmail(*req.Email) {
			return nil, fmt.Errorf("invalid email format")
		}
		user.Email = *req.Email
	}

	if req.Username != nil {
		if len(*req.Username) < 3 || len(*req.Username) > 50 {
			return nil, fmt.Errorf("username must be between 3 and 50 characters")
		}
		user.Username = *req.Username
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedAt = time.Now()

	// Update in database
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser removes a user
func (s *DefaultUserService) DeleteUser(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers retrieves paginated user list
func (s *DefaultUserService) ListUsers(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// GetUserByEmail retrieves a user by email
func (s *DefaultUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *DefaultUserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

// isValidEmail checks if email format is valid
func isValidEmail(email string) bool {
	// Simple email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
