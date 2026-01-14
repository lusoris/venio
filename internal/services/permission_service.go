// Package services contains business logic
package services

import (
	"context"
	"errors"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/repositories"
)

// PermissionService handles business logic for permissions
type PermissionService interface {
	GetByID(ctx context.Context, id int64) (*models.Permission, error)
	GetByName(ctx context.Context, name string) (*models.Permission, error)
	Create(ctx context.Context, req models.CreatePermissionRequest) (*models.Permission, error)
	Update(ctx context.Context, id int64, req models.UpdatePermissionRequest) (*models.Permission, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Permission, int64, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Permission, error)
}

type permissionService struct {
	permissionRepository repositories.PermissionRepository
}

// NewPermissionService creates a new permission service
func NewPermissionService(permissionRepository repositories.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepository: permissionRepository,
	}
}

// GetByID retrieves a permission by ID
func (s *permissionService) GetByID(ctx context.Context, id int64) (*models.Permission, error) {
	if id <= 0 {
		return nil, errors.New("invalid permission ID")
	}

	permission, err := s.permissionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if permission == nil {
		return nil, errors.New("permission not found")
	}

	return permission, nil
}

// GetByName retrieves a permission by name
func (s *permissionService) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	if name == "" {
		return nil, errors.New("permission name cannot be empty")
	}

	permission, err := s.permissionRepository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if permission == nil {
		return nil, errors.New("permission not found")
	}

	return permission, nil
}

// Create creates a new permission
func (s *permissionService) Create(ctx context.Context, req models.CreatePermissionRequest) (*models.Permission, error) {
	// Validation is already done by request model tags
	// Check if permission name already exists
	existing, err := s.permissionRepository.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("permission with this name already exists")
	}

	permission, err := s.permissionRepository.Create(ctx, &req)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// Update updates an existing permission
func (s *permissionService) Update(ctx context.Context, id int64, req models.UpdatePermissionRequest) (*models.Permission, error) {
	if id <= 0 {
		return nil, errors.New("invalid permission ID")
	}

	// Check if permission exists
	existing, err := s.permissionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("permission not found")
	}

	// If name is being updated, check for duplicates
	if req.Name != nil && *req.Name != existing.Name {
		duplicate, dupErr := s.permissionRepository.GetByName(ctx, *req.Name)
		if dupErr != nil {
			return nil, dupErr
		}

		if duplicate != nil {
			return nil, errors.New("permission with this name already exists")
		}
	}

	permission, err := s.permissionRepository.Update(ctx, id, &req)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// Delete deletes a permission
func (s *permissionService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid permission ID")
	}

	// Check if permission exists
	permission, err := s.permissionRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if permission == nil {
		return errors.New("permission not found")
	}

	return s.permissionRepository.Delete(ctx, id)
}

// List lists all permissions with pagination
func (s *permissionService) List(ctx context.Context, limit, offset int) ([]*models.Permission, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	permissions, total, err := s.permissionRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Convert []models.Permission to []*models.Permission
	permPointers := make([]*models.Permission, len(permissions))
	for i := range permissions {
		permPointers[i] = &permissions[i]
	}

	return permPointers, total, nil
}

// GetByUserID retrieves all permissions for a user through roles
func (s *permissionService) GetByUserID(ctx context.Context, userID int64) ([]*models.Permission, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	permissions, err := s.permissionRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert []models.Permission to []*models.Permission
	permPointers := make([]*models.Permission, len(permissions))
	for i := range permissions {
		permPointers[i] = &permissions[i]
	}

	return permPointers, nil
}
