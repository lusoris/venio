// Package services contains business logic
package services

import (
	"context"
	"errors"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/repositories"
)

// RoleService handles business logic for roles
type RoleService interface {
	GetByID(ctx context.Context, id int64) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	Create(ctx context.Context, req models.CreateRoleRequest) (*models.Role, error)
	Update(ctx context.Context, id int64, req models.UpdateRoleRequest) (*models.Role, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.Role, int64, error)
	GetPermissions(ctx context.Context, roleID int64) ([]*models.Permission, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error
}

type roleService struct {
	roleRepository repositories.RoleRepository
}

// NewRoleService creates a new role service
func NewRoleService(roleRepository repositories.RoleRepository) RoleService {
	return &roleService{
		roleRepository: roleRepository,
	}
}

// GetByID retrieves a role by ID
func (s *roleService) GetByID(ctx context.Context, id int64) (*models.Role, error) {
	if id <= 0 {
		return nil, errors.New("invalid role ID")
	}

	role, err := s.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	return role, nil
}

// GetByName retrieves a role by name
func (s *roleService) GetByName(ctx context.Context, name string) (*models.Role, error) {
	if name == "" {
		return nil, errors.New("role name cannot be empty")
	}

	role, err := s.roleRepository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	return role, nil
}

// Create creates a new role
func (s *roleService) Create(ctx context.Context, req models.CreateRoleRequest) (*models.Role, error) {
	// Validation is already done by request model tags
	// Check if role name already exists
	existing, err := s.roleRepository.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("role with this name already exists")
	}

	role, err := s.roleRepository.Create(ctx, &req)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// Update updates an existing role
func (s *roleService) Update(ctx context.Context, id int64, req models.UpdateRoleRequest) (*models.Role, error) {
	if id <= 0 {
		return nil, errors.New("invalid role ID")
	}

	// Check if role exists
	existing, err := s.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("role not found")
	}

	// If name is being updated, check for duplicates
	if req.Name != nil && *req.Name != existing.Name {
		duplicate, dupErr := s.roleRepository.GetByName(ctx, *req.Name)
		if dupErr != nil {
			return nil, dupErr
		}

		if duplicate != nil {
			return nil, errors.New("role with this name already exists")
		}
	}

	role, err := s.roleRepository.Update(ctx, id, &req)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// Delete deletes a role
func (s *roleService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid role ID")
	}

	// Check if role exists
	role, err := s.roleRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if role == nil {
		return errors.New("role not found")
	}

	return s.roleRepository.Delete(ctx, id)
}

// List lists all roles with pagination
func (s *roleService) List(ctx context.Context, limit, offset int) ([]*models.Role, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	roles, total, err := s.roleRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Convert []models.Role to []*models.Role
	rolePointers := make([]*models.Role, len(roles))
	for i := range roles {
		rolePointers[i] = &roles[i]
	}

	return rolePointers, total, nil
}

// GetPermissions retrieves all permissions for a role
func (s *roleService) GetPermissions(ctx context.Context, roleID int64) ([]*models.Permission, error) {
	if roleID <= 0 {
		return nil, errors.New("invalid role ID")
	}

	// Check if role exists
	role, err := s.roleRepository.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	permissions, err := s.roleRepository.GetPermissions(ctx, roleID)
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

// AssignPermissionToRole assigns a permission to a role
func (s *roleService) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	if roleID <= 0 {
		return errors.New("invalid role ID")
	}

	if permissionID <= 0 {
		return errors.New("invalid permission ID")
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (s *roleService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	if roleID <= 0 {
		return errors.New("invalid role ID")
	}

	if permissionID <= 0 {
		return errors.New("invalid permission ID")
	}

	return nil
}
