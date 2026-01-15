---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/repositories/**/*.go"
description: Repository Layer Patterns & Interface Design
---

# Repository Layer Patterns

## Core Principle

**Repositories are the ONLY layer that talks to the database.** They return domain models, handle connection pooling, and manage transactions.

## Interface Design Rules

### 1. Return Type Conventions

**CRITICAL**: Return types must be consistent across all repositories.

#### Single Entity Operations

```go
// ✅ Return pointer for single entity
GetByID(ctx context.Context, id int64) (*models.User, error)
GetByEmail(ctx context.Context, email string) (*models.User, error)
Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
Update(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error)
```

**Why pointer?** 
- Can return `nil` to indicate "not found"
- Efficient for large structs
- Standard Go convention for optional returns

#### Collection Operations

```go
// ✅ Return slice of VALUES for collections
List(ctx context.Context, limit, offset int) ([]models.Role, int64, error)
GetAll(ctx context.Context) ([]models.Permission, error)
GetByUserID(ctx context.Context, userID int64) ([]models.Role, error)
```

**Why values not pointers?**
- List operations iterate over data anyway
- Avoid pointer allocation overhead for collections
- Services can convert to pointers if needed

**Important**: Always return `(collection, total, error)` for paginated lists:
```go
// ✅ Paginated list returns total count
List(ctx context.Context, limit, offset int) ([]models.Role, int64, error)

// ✅ Non-paginated list doesn't need count
GetAll(ctx context.Context) ([]models.Permission, error)
```

#### Delete Operations

```go
// ✅ Return only error
Delete(ctx context.Context, id int64) error
RemoveRole(ctx context.Context, userID, roleID int64) error
```

#### Boolean Checks

```go
// ✅ Return bool + error
Exists(ctx context.Context, email string) (bool, error)
HasRole(ctx context.Context, userID int64, roleName string) (bool, error)
HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error)
```

### 2. Standard Interface Pattern

Every repository MUST follow this structure:

```go
// internal/repositories/user_repository.go

package repositories

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/lusoris/venio/internal/models"
)

// UserRepository defines ALL data access operations for users
type UserRepository interface {
    // Single entity operations - return pointer
    GetByID(ctx context.Context, id int64) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByUsername(ctx context.Context, username string) (*models.User, error)
    
    // Create/Update - return pointer
    Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
    Update(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error)
    
    // Delete - return error only
    Delete(ctx context.Context, id int64) error
    
    // List - return slice of values + total count
    List(ctx context.Context, limit, offset int) ([]models.User, int64, error)
    
    // Boolean checks
    Exists(ctx context.Context, email string) (bool, error)
}

// Private implementation struct
type userRepository struct {
    pool *pgxpool.Pool
}

// Constructor returns interface
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
    return &userRepository{pool: pool}
}
```

### 3. Request/Response DTOs

**Create operations**: Accept request DTO, return model

```go
// ✅ CORRECT
Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)

// ❌ WRONG: Don't accept full model
Create(ctx context.Context, user *models.User) (int64, error)
```

**Update operations**: Accept ID + update DTO, return updated model

```go
// ✅ CORRECT
Update(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error)

// ❌ WRONG: Don't return just error
Update(ctx context.Context, user *models.User) error
```

**Why?** 
- Request DTOs enforce validation
- Returning full model avoids extra DB round-trip
- Services get updated entity immediately

### 4. Association Operations

For many-to-many relationships:

```go
// RolePermissionRepository interface
type RolePermissionRepository interface {
    // Get associations
    GetPermissions(ctx context.Context, roleID int64) ([]models.Permission, error)
    GetRoles(ctx context.Context, permissionID int64) ([]models.Role, error)
    
    // Modify associations
    AssignToRole(ctx context.Context, roleID, permissionID int64) error
    RemoveFromRole(ctx context.Context, roleID, permissionID int64) error
}

// UserRoleRepository interface
type UserRoleRepository interface {
    GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error)
    AssignRole(ctx context.Context, userID, roleID int64) error
    RemoveRole(ctx context.Context, userID, roleID int64) error
    HasRole(ctx context.Context, userID int64, roleName string) (bool, error)
    HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error)
}
```

## Implementation Patterns

### 1. Query Execution

```go
func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    var user models.User

    query := `SELECT id, email, username, first_name, last_name, avatar, is_active, created_at, updated_at 
              FROM users WHERE id = $1`
    
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Username,
        &user.FirstName,
        &user.LastName,
        &user.Avatar,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil  // ✅ Not found = nil, nil
        }
        return nil, fmt.Errorf("get user by id: %w", err)
    }

    return &user, nil
}
```

### 2. List Operations

```go
func (r *roleRepository) List(ctx context.Context, limit, offset int) ([]models.Role, int64, error) {
    // 1. Get total count
    var total int64
    countQuery := `SELECT COUNT(*) FROM roles`
    err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("count roles: %w", err)
    }

    // 2. Get paginated results
    query := `SELECT id, name, description, created_at 
              FROM roles 
              ORDER BY name 
              LIMIT $1 OFFSET $2`
    
    rows, err := r.pool.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("list roles: %w", err)
    }
    defer rows.Close()

    // 3. Scan into slice of VALUES
    roles := []models.Role{}  // ✅ Values, not pointers
    for rows.Next() {
        var role models.Role
        err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
        if err != nil {
            return nil, 0, fmt.Errorf("scan role: %w", err)
        }
        roles = append(roles, role)  // ✅ Append value
    }

    if err := rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows error: %w", err)
    }

    return roles, total, nil
}
```

### 3. Create Operations

```go
func (r *permissionRepository) Create(ctx context.Context, req *models.CreatePermissionRequest) (*models.Permission, error) {
    var perm models.Permission

    query := `INSERT INTO permissions (name, description) 
              VALUES ($1, $2) 
              RETURNING id, name, description, created_at`
    
    err := r.pool.QueryRow(ctx, query, req.Name, req.Description).Scan(
        &perm.ID,
        &perm.Name,
        &perm.Description,
        &perm.CreatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("create permission: %w", err)
    }

    return &perm, nil  // ✅ Return created entity immediately
}
```

### 4. Update Operations

```go
func (r *roleRepository) Update(ctx context.Context, id int64, req *models.UpdateRoleRequest) (*models.Role, error) {
    var role models.Role

    // Build dynamic update query based on non-nil fields
    query := `UPDATE roles SET `
    args := []interface{}{}
    argCount := 1

    if req.Name != nil {
        query += fmt.Sprintf("name = $%d, ", argCount)
        args = append(args, *req.Name)
        argCount++
    }

    if req.Description != nil {
        query += fmt.Sprintf("description = $%d, ", argCount)
        args = append(args, *req.Description)
        argCount++
    }

    query += fmt.Sprintf("updated_at = NOW() WHERE id = $%d ", argCount)
    args = append(args, id)

    query += `RETURNING id, name, description, created_at, updated_at`

    err := r.pool.QueryRow(ctx, query, args...).Scan(
        &role.ID,
        &role.Name,
        &role.Description,
        &role.CreatedAt,
        &role.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil  // ✅ Not found
        }
        return nil, fmt.Errorf("update role: %w", err)
    }

    return &role, nil
}
```

### 5. Association Operations

```go
func (ur *userRoleRepository) AssignRole(ctx context.Context, userID, roleID int64) error {
    query := `INSERT INTO user_roles (user_id, role_id) 
              VALUES ($1, $2) 
              ON CONFLICT (user_id, role_id) DO NOTHING`
    
    _, err := ur.pool.Exec(ctx, query, userID, roleID)
    if err != nil {
        return fmt.Errorf("assign role to user: %w", err)
    }

    return nil
}

func (ur *userRoleRepository) HasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
    var exists bool

    query := `SELECT EXISTS(
        SELECT 1 FROM user_roles ur
        INNER JOIN roles r ON ur.role_id = r.id
        WHERE ur.user_id = $1 AND r.name = $2
    )`

    err := ur.pool.QueryRow(ctx, query, userID, roleName).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("check user role: %w", err)
    }

    return exists, nil
}
```

## Error Handling Patterns

### 1. Not Found vs Error

```go
// ✅ CORRECT: Distinguish between "not found" and "error"
func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    // ...
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, nil  // ✅ Not found = (nil, nil)
    }
    return nil, fmt.Errorf("get user: %w", err)  // ✅ Real error = (nil, error)
}

// ❌ WRONG: Return error for "not found"
func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    // ...
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, errors.New("user not found")  // ❌ Service can't distinguish
    }
    return nil, err
}
```

### 2. Wrap Errors with Context

```go
// ✅ CORRECT: Add context to errors
return nil, fmt.Errorf("create user: %w", err)
return nil, fmt.Errorf("list roles with limit %d offset %d: %w", limit, offset, err)

// ❌ WRONG: Return raw error
return nil, err
```

## Batch Operations

For bulk inserts/updates:

```go
func (r *userRepository) CreateBatch(ctx context.Context, users []*models.CreateUserRequest) error {
    batch := &pgx.Batch{}

    for _, req := range users {
        query := `INSERT INTO users (email, username, password_hash, first_name, last_name) 
                  VALUES ($1, $2, $3, $4, $5)`
        batch.Queue(query, req.Email, req.Username, req.Password, req.FirstName, req.LastName)
    }

    br := r.pool.SendBatch(ctx, batch)
    defer br.Close()

    for i := 0; i < len(users); i++ {
        _, err := br.Exec()
        if err != nil {
            return fmt.Errorf("batch insert user %d: %w", i, err)
        }
    }

    return nil
}
```

## Transaction Support

```go
func (r *userRepository) CreateWithRoles(ctx context.Context, req *models.CreateUserRequest, roleIDs []int64) (*models.User, error) {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // 1. Create user
    var user models.User
    query := `INSERT INTO users (email, username, ...) VALUES (...) RETURNING ...`
    err = tx.QueryRow(ctx, query, ...).Scan(...)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }

    // 2. Assign roles
    for _, roleID := range roleIDs {
        _, err = tx.Exec(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, user.ID, roleID)
        if err != nil {
            return nil, fmt.Errorf("assign role %d: %w", roleID, err)
        }
    }

    // 3. Commit
    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("commit transaction: %w", err)
    }

    return &user, nil
}
```

## Repository Factory Pattern

```go
// internal/repositories/factory.go

type Factory struct {
    pool *pgxpool.Pool
}

func NewFactory(pool *pgxpool.Pool) *Factory {
    return &Factory{pool: pool}
}

func (f *Factory) All() (
    UserRepository,
    RoleRepository,
    PermissionRepository,
    UserRoleRepository,
) {
    return NewUserRepository(f.pool),
        NewRoleRepository(f.pool),
        NewPermissionRepository(f.pool),
        NewUserRoleRepository(f.pool)
}
```

## Checklist for New Repositories

- [ ] Interface defined with ALL operations
- [ ] Single entity operations return `(*Model, error)`
- [ ] Collection operations return `([]Model, int64, error)` or `([]Model, error)`
- [ ] Delete operations return `error` only
- [ ] Boolean checks return `(bool, error)`
- [ ] Create/Update accept request DTOs
- [ ] Create/Update return full model
- [ ] Distinguish "not found" (nil, nil) from errors (nil, err)
- [ ] Wrap errors with context
- [ ] Use parameterized queries (prevent SQL injection)
- [ ] Handle pgx.ErrNoRows correctly
- [ ] Add factory constructor

## Common Mistakes

❌ **Don't:**
- Return `int64` from Create (return full model)
- Return `error` from Update (return updated model)
- Mix pointers and values in collection returns
- Return errors for "not found" cases
- Use string concatenation for SQL (SQL injection)
- Ignore `rows.Err()` after iteration

✅ **Do:**
- Return complete models from Create/Update
- Use consistent return types across repositories
- Handle "not found" vs "error" correctly
- Use parameterized queries
- Wrap errors with context
- Check `rows.Err()` after loops

---

**Remember**: Repositories are the truth source for data access patterns. Keep interfaces consistent, return types predictable, and error handling clear.
