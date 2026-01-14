// Package repositories contains data access implementations
package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/venio/internal/models"
)

// PermissionRepository defines permission data access operations
type PermissionRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Permission, error)
	GetByName(ctx context.Context, name string) (*models.Permission, error)
	Create(ctx context.Context, req *models.CreatePermissionRequest) (*models.Permission, error)
	Update(ctx context.Context, id int64, req *models.UpdatePermissionRequest) (*models.Permission, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]models.Permission, int64, error)
	GetByUserID(ctx context.Context, userID int64) ([]models.Permission, error)
	AssignToRole(ctx context.Context, roleID, permissionID int64) error
	RemoveFromRole(ctx context.Context, roleID, permissionID int64) error
}

type permissionRepository struct {
	pool *pgxpool.Pool
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(pool *pgxpool.Pool) PermissionRepository {
	return &permissionRepository{pool: pool}
}

// GetByID retrieves a permission by ID
func (p *permissionRepository) GetByID(ctx context.Context, id int64) (*models.Permission, error) {
	var perm models.Permission

	query := `SELECT id, name, description, created_at FROM permissions WHERE id = $1`
	err := p.pool.QueryRow(ctx, query, id).Scan(&perm.ID, &perm.Name, &perm.Description, &perm.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("get permission by id: %w", err)
	}

	return &perm, nil
}

// GetByName retrieves a permission by name
func (p *permissionRepository) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	var perm models.Permission

	query := `SELECT id, name, description, created_at FROM permissions WHERE name = $1`
	err := p.pool.QueryRow(ctx, query, name).Scan(&perm.ID, &perm.Name, &perm.Description, &perm.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("get permission by name: %w", err)
	}

	return &perm, nil
}

// Create creates a new permission
func (p *permissionRepository) Create(ctx context.Context, req *models.CreatePermissionRequest) (*models.Permission, error) {
	var perm models.Permission

	query := `
		INSERT INTO permissions (name, description, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, name, description, created_at
	`

	err := p.pool.QueryRow(ctx, query, req.Name, req.Description).
		Scan(&perm.ID, &perm.Name, &perm.Description, &perm.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create permission: %w", err)
	}

	return &perm, nil
}

// Update updates an existing permission
func (p *permissionRepository) Update(ctx context.Context, id int64, req *models.UpdatePermissionRequest) (*models.Permission, error) {
	var perm models.Permission
	var name, description string

	// Get current values
	getQuery := `SELECT name, description FROM permissions WHERE id = $1`
	err := p.pool.QueryRow(ctx, getQuery, id).Scan(&name, &description)
	if err != nil {
		return nil, fmt.Errorf("get permission for update: %w", err)
	}

	// Override with request values if provided
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}

	// Update permission
	query := `
		UPDATE permissions
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description, created_at
	`

	err = p.pool.QueryRow(ctx, query, name, description, id).
		Scan(&perm.ID, &perm.Name, &perm.Description, &perm.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("update permission: %w", err)
	}

	return &perm, nil
}

// Delete deletes a permission
func (p *permissionRepository) Delete(ctx context.Context, id int64) error {
	// Check if permission is assigned to roles (prevent deletion of in-use permissions)
	checkQuery := `SELECT COUNT(*) FROM role_permissions WHERE permission_id = $1`
	var count int64
	err := p.pool.QueryRow(ctx, checkQuery, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("check permission usage: %w", err)
	}

	if count > 0 {
		return errors.New("cannot delete permission that is assigned to roles")
	}

	query := `DELETE FROM permissions WHERE id = $1`
	result, err := p.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete permission: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("permission not found")
	}

	return nil
}

// List retrieves a paginated list of permissions
func (p *permissionRepository) List(ctx context.Context, limit, offset int) ([]models.Permission, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM permissions`
	var total int64
	err := p.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count permissions: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, name, description, created_at
		FROM permissions
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := p.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list permissions: %w", err)
	}
	defer rows.Close()

	permissions := []models.Permission{}
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return permissions, total, nil
}

// GetByUserID retrieves all permissions for a user (through roles)
func (p *permissionRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Permission, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY p.name
	`

	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
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

// AssignToRole assigns a permission to a role
func (p *permissionRepository) AssignToRole(ctx context.Context, roleID, permissionID int64) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id, assigned_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT DO NOTHING
	`

	_, err := p.pool.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("assign permission to role: %w", err)
	}

	return nil
}

// RemoveFromRole removes a permission from a role
func (p *permissionRepository) RemoveFromRole(ctx context.Context, roleID, permissionID int64) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	result, err := p.pool.Exec(ctx, query, roleID, permissionID)

	if err != nil {
		return fmt.Errorf("remove permission from role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("permission not assigned to role")
	}

	return nil
}
