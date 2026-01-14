// Package services contains business logic
package services

import (
	"context"
	"errors"

	"github.com/lusoris/venio/internal/repositories"
)

// UserRoleService handles business logic for user-role assignments
type UserRoleService interface {
	GetUserRoles(ctx context.Context, userID int64) ([]string, error)
	AssignRole(ctx context.Context, userID, roleID int64) error
	RemoveRole(ctx context.Context, userID, roleID int64) error
	HasRole(ctx context.Context, userID int64, roleName string) (bool, error)
	HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error)
}

type userRoleService struct {
	userRoleRepository repositories.UserRoleRepository
}

// NewUserRoleService creates a new user-role service
func NewUserRoleService(userRoleRepository repositories.UserRoleRepository) UserRoleService {
	return &userRoleService{
		userRoleRepository: userRoleRepository,
	}
}

// GetUserRoles retrieves all roles for a user
func (s *userRoleService) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	roles, err := s.userRoleRepository.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Extract role names
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return roleNames, nil
}

// AssignRole assigns a role to a user
func (s *userRoleService) AssignRole(ctx context.Context, userID, roleID int64) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	if roleID <= 0 {
		return errors.New("invalid role ID")
	}

	return s.userRoleRepository.AssignRole(ctx, userID, roleID)
}

// RemoveRole removes a role from a user
func (s *userRoleService) RemoveRole(ctx context.Context, userID, roleID int64) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	if roleID <= 0 {
		return errors.New("invalid role ID")
	}

	return s.userRoleRepository.RemoveRole(ctx, userID, roleID)
}

// HasRole checks if a user has a specific role
func (s *userRoleService) HasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	if userID <= 0 {
		return false, errors.New("invalid user ID")
	}

	if roleName == "" {
		return false, errors.New("role name cannot be empty")
	}

	return s.userRoleRepository.HasRole(ctx, userID, roleName)
}

// HasPermission checks if a user has a specific permission through roles
func (s *userRoleService) HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	if userID <= 0 {
		return false, errors.New("invalid user ID")
	}

	if permissionName == "" {
		return false, errors.New("permission name cannot be empty")
	}

	return s.userRoleRepository.HasPermission(ctx, userID, permissionName)
}
