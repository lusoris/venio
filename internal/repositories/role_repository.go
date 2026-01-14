// Package repositories contains data access implementations
package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/venio/internal/models"
)

// RoleRepository defines role data access operations
type RoleRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	Create(ctx context.Context, req *models.CreateRoleRequest) (*models.Role, error)
	Update(ctx context.Context, id int64, req *models.UpdateRoleRequest) (*models.Role, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]models.Role, int64, error)
	GetPermissions(ctx context.Context, roleID int64) ([]models.Permission, error)
}

type roleRepository struct {
	pool *pgxpool.Pool
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(pool *pgxpool.Pool) RoleRepository {
	return &roleRepository{pool: pool}
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(ctx context.Context, id int64) (*models.Role, error) {
	var role models.Role

	query := `SELECT id, name, description, created_at FROM roles WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("get role by id: %w", err)
	}

	return &role, nil
}

// GetByName retrieves a role by name
func (r *roleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role

	query := `SELECT id, name, description, created_at FROM roles WHERE name = $1`
	err := r.pool.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("get role by name: %w", err)
	}

	return &role, nil
}

// Create creates a new role
func (r *roleRepository) Create(ctx context.Context, req *models.CreateRoleRequest) (*models.Role, error) {
	var role models.Role

	query := `
		INSERT INTO roles (name, description, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, name, description, created_at
	`

	err := r.pool.QueryRow(ctx, query, req.Name, req.Description).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	return &role, nil
}

// Update updates an existing role
func (r *roleRepository) Update(ctx context.Context, id int64, req *models.UpdateRoleRequest) (*models.Role, error) {
	var role models.Role
	var name, description string

	// Get current values
	getQuery := `SELECT name, description FROM roles WHERE id = $1`
	err := r.pool.QueryRow(ctx, getQuery, id).Scan(&name, &description)
	if err != nil {
		return nil, fmt.Errorf("get role for update: %w", err)
	}

	// Override with request values if provided
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}

	// Update role
	query := `
		UPDATE roles
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description, created_at
	`

	err = r.pool.QueryRow(ctx, query, name, description, id).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	return &role, nil
}

// Delete deletes a role
func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	// Check if role is assigned to users (prevent deletion of in-use roles)
	checkQuery := `SELECT COUNT(*) FROM user_roles WHERE role_id = $1`
	var count int64
	err := r.pool.QueryRow(ctx, checkQuery, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("check role usage: %w", err)
	}

	if count > 0 {
		return errors.New("cannot delete role that is assigned to users")
	}

	query := `DELETE FROM roles WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("role not found")
	}

	return nil
}

// List retrieves a paginated list of roles
func (r *roleRepository) List(ctx context.Context, limit, offset int) ([]models.Role, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM roles`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count roles: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, name, description, created_at
		FROM roles
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list roles: %w", err)
	}
	defer rows.Close()

	roles := []models.Role{}
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return roles, total, nil
}

// GetPermissions retrieves all permissions for a role
func (r *roleRepository) GetPermissions(ctx context.Context, roleID int64) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.name
	`

	rows, err := r.pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("get role permissions: %w", err)
	}
	defer rows.Close()

	permissions := []models.Permission{}
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return permissions, nil
}
