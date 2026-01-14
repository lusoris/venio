// Package repositories contains data access implementations
package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/venio/internal/models"
)

// UserRoleRepository defines user-role assignment operations
type UserRoleRepository interface {
	GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error)
	AssignRole(ctx context.Context, userID, roleID int64) error
	RemoveRole(ctx context.Context, userID, roleID int64) error
	HasRole(ctx context.Context, userID int64, roleName string) (bool, error)
	HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error)
}

type userRoleRepository struct {
	pool *pgxpool.Pool
}

// NewUserRoleRepository creates a new user-role repository
func NewUserRoleRepository(pool *pgxpool.Pool) UserRoleRepository {
	return &userRoleRepository{pool: pool}
}

// GetUserRoles retrieves all roles for a user
func (ur *userRoleRepository) GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at
		FROM roles r
		INNER JOIN user_roles urt ON r.id = urt.role_id
		WHERE urt.user_id = $1
		ORDER BY r.name
	`

	rows, err := ur.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	defer rows.Close()

	roles := []models.Role{}
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return roles, nil
}

// AssignRole assigns a role to a user
func (ur *userRoleRepository) AssignRole(ctx context.Context, userID, roleID int64) error {
	// Check if user and role exist
	userQuery := `SELECT id FROM users WHERE id = $1`
	roleQuery := `SELECT id FROM roles WHERE id = $1`

	var uid, rid int64
	err := ur.pool.QueryRow(ctx, userQuery, userID).Scan(&uid)
	if err != nil {
		return errors.New("user not found")
	}

	err = ur.pool.QueryRow(ctx, roleQuery, roleID).Scan(&rid)
	if err != nil {
		return errors.New("role not found")
	}

	// Assign role
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT DO NOTHING
	`

	_, err = ur.pool.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("assign role to user: %w", err)
	}

	return nil
}

// RemoveRole removes a role from a user
func (ur *userRoleRepository) RemoveRole(ctx context.Context, userID, roleID int64) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	result, err := ur.pool.Exec(ctx, query, userID, roleID)

	if err != nil {
		return fmt.Errorf("remove role from user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("role not assigned to user")
	}

	return nil
}

// HasRole checks if a user has a specific role
func (ur *userRoleRepository) HasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM user_roles urt
		INNER JOIN roles r ON urt.role_id = r.id
		WHERE urt.user_id = $1 AND r.name = $2
	`

	var hasRole bool
	err := ur.pool.QueryRow(ctx, query, userID, roleName).Scan(&hasRole)
	if err != nil {
		return false, fmt.Errorf("check user role: %w", err)
	}

	return hasRole, nil
}

// HasPermission checks if a user has a specific permission (through roles)
func (ur *userRoleRepository) HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN user_roles urt ON r.id = urt.role_id
		WHERE urt.user_id = $1 AND p.name = $2
	`

	var hasPermission bool
	err := ur.pool.QueryRow(ctx, query, userID, permissionName).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("check user permission: %w", err)
	}

	return hasPermission, nil
}
