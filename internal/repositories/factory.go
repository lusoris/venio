// Package repositories provides a factory for creating repositories
package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Factory creates repository instances
type Factory struct {
	pool *pgxpool.Pool
}

// NewFactory creates a new repository factory
func NewFactory(pool *pgxpool.Pool) *Factory {
	return &Factory{pool: pool}
}

// User creates a new User repository
func (f *Factory) User() UserRepository {
	return NewPostgresUserRepository(f.pool)
}

// Role creates a new Role repository
func (f *Factory) Role() RoleRepository {
	return NewRoleRepository(f.pool)
}

// Permission creates a new Permission repository
func (f *Factory) Permission() PermissionRepository {
	return NewPermissionRepository(f.pool)
}

// UserRole creates a new UserRole repository
func (f *Factory) UserRole() UserRoleRepository {
	return NewUserRoleRepository(f.pool)
}

// All creates all repositories at once
// Returns: userRepo, roleRepo, permissionRepo, userRoleRepo
func (f *Factory) All() (UserRepository, RoleRepository, PermissionRepository, UserRoleRepository) {
	return f.User(), f.Role(), f.Permission(), f.UserRole()
}
