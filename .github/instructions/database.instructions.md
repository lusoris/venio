---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "migrations/**/*.sql,internal/repositories/**/*.go"
description: Database Best Practices & Migration Guidelines
---

# Database Guidelines

## PostgreSQL Best Practices

### Connection Management

**Use pgxpool for connection pooling:**

```go
// ✅ DO: Use pgxpool for connection pooling
import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(ctx context.Context, connString string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, fmt.Errorf("failed to parse connection string: %w", err)
    }
    
    // Configure pool
    config.MaxConns = 25
    config.MinConns = 5
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = 30 * time.Minute
    config.HealthCheckPeriod = 1 * time.Minute
    
    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }
    
    // Test connection
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return pool, nil
}

// ❌ DON'T: Create new connection per request
func GetUser(id int64) (*User, error) {
    conn, _ := pgx.Connect(context.Background(), connString)  // ❌ New connection
    defer conn.Close(context.Background())
    
    // Query...
}
```

### Query Patterns

**Use parameterized queries to prevent SQL injection:**

```go
// ✅ DO: Parameterized queries
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    query := `SELECT id, email, password FROM users WHERE email = $1`
    
    var user User
    err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("query failed: %w", err)
    }
    
    return &user, nil
}

// ❌ DON'T: String concatenation (SQL injection risk)
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)  // ❌ SQL injection
    // Query...
}
```

### Batch Operations

**Use pgx.Batch for multiple operations:**

```go
// ✅ DO: Batch operations for efficiency
func (r *UserRepository) CreateMultiple(ctx context.Context, users []*User) error {
    batch := &pgx.Batch{}
    
    query := `INSERT INTO users (email, password) VALUES ($1, $2)`
    for _, user := range users {
        batch.Queue(query, user.Email, user.Password)
    }
    
    results := r.pool.SendBatch(ctx, batch)
    defer results.Close()
    
    for i := 0; i < len(users); i++ {
        _, err := results.Exec()
        if err != nil {
            return fmt.Errorf("failed to insert user %d: %w", i, err)
        }
    }
    
    return nil
}

// ❌ DON'T: Individual queries in loop (N+1 problem)
func (r *UserRepository) CreateMultiple(ctx context.Context, users []*User) error {
    for _, user := range users {
        _, err := r.pool.Exec(ctx, query, user.Email, user.Password)  // ❌ N queries
        if err != nil {
            return err
        }
    }
    return nil
}
```

### Transactions

**Handle transactions properly:**

```go
// ✅ DO: Proper transaction handling
func (r *UserRepository) CreateUserWithRole(ctx context.Context, user *User, roleID int64) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    // Always rollback on panic or error
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback(ctx)
            panic(p)
        } else if err != nil {
            tx.Rollback(ctx)
        }
    }()
    
    // Insert user
    var userID int64
    err = tx.QueryRow(ctx,
        `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
        user.Email, user.Password,
    ).Scan(&userID)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    
    // Assign role
    _, err = tx.Exec(ctx,
        `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`,
        userID, roleID,
    )
    if err != nil {
        return fmt.Errorf("failed to assign role: %w", err)
    }
    
    // Commit transaction
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

## Migration Guidelines

### Migration Files

**Use numbered migrations with up/down:**

```
migrations/
├── 001_initial_schema.up.sql
├── 001_initial_schema.down.sql
├── 002_seed_roles_and_permissions.up.sql
├── 002_seed_roles_and_permissions.down.sql
├── 003_add_user_preferences.up.sql
└── 003_add_user_preferences.down.sql
```

### Migration Naming

**Format:** `{number}_{description}.{up|down}.sql`

```
✅ DO:
001_create_users_table.up.sql
002_add_email_index.up.sql
003_add_phone_column.up.sql

❌ DON'T:
create_users.sql              (no number, no direction)
2023_01_15_users.sql          (date-based, hard to order)
001_migration.up.sql          (vague description)
```

### Schema Migrations

**Create tables with proper constraints:**

```sql
-- ✅ DO: Proper table definition
-- migrations/001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(50) UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_active ON users(is_active) WHERE deleted_at IS NULL;

-- Add check constraints
ALTER TABLE users ADD CONSTRAINT check_email_format 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Add comment
COMMENT ON TABLE users IS 'Core user accounts';
COMMENT ON COLUMN users.email IS 'Unique email address for authentication';
COMMENT ON COLUMN users.is_active IS 'Account status - inactive accounts cannot log in';

-- migrations/001_create_users_table.down.sql
DROP TABLE IF EXISTS users CASCADE;
```

### Data Migrations

**Separate schema and data migrations:**

```sql
-- ✅ DO: Data migration
-- migrations/002_seed_default_roles.up.sql
INSERT INTO roles (name, description) VALUES
    ('admin', 'System administrator with full access'),
    ('user', 'Regular user with basic access'),
    ('moderator', 'Content moderator')
ON CONFLICT (name) DO NOTHING;

-- migrations/002_seed_default_roles.down.sql
DELETE FROM roles WHERE name IN ('admin', 'user', 'moderator');
```

### Adding Columns

**Add columns with defaults for existing data:**

```sql
-- ✅ DO: Add column with default
-- migrations/003_add_phone_column.up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users ADD COLUMN phone_verified BOOLEAN NOT NULL DEFAULT false;

-- Create index after adding data (if table is large)
CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL;

-- migrations/003_add_phone_column.down.sql
DROP INDEX IF EXISTS idx_users_phone;
ALTER TABLE users DROP COLUMN IF EXISTS phone_verified;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### Modifying Columns

**Handle data carefully when modifying columns:**

```sql
-- ✅ DO: Safe column modification
-- migrations/004_expand_email_length.up.sql

-- 1. Add new column
ALTER TABLE users ADD COLUMN email_new VARCHAR(320);

-- 2. Copy data
UPDATE users SET email_new = email;

-- 3. Drop old column
ALTER TABLE users DROP COLUMN email;

-- 4. Rename new column
ALTER TABLE users RENAME COLUMN email_new TO email;

-- 5. Add constraints
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);

-- migrations/004_expand_email_length.down.sql
-- Reverse the process
ALTER TABLE users ADD COLUMN email_old VARCHAR(255);
UPDATE users SET email_old = email;
ALTER TABLE users DROP COLUMN email;
ALTER TABLE users RENAME COLUMN email_old TO email;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);
```

### Foreign Keys

**Use foreign keys with proper cascade behavior:**

```sql
-- ✅ DO: Foreign keys with cascade
CREATE TABLE user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by BIGINT,
    
    -- Foreign keys with appropriate cascade
    CONSTRAINT fk_user_roles_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,  -- Delete user_roles when user is deleted
    
    CONSTRAINT fk_user_roles_role 
        FOREIGN KEY (role_id) 
        REFERENCES roles(id) 
        ON DELETE RESTRICT,  -- Prevent role deletion if in use
    
    CONSTRAINT fk_user_roles_assigned_by 
        FOREIGN KEY (assigned_by) 
        REFERENCES users(id) 
        ON DELETE SET NULL,  -- Keep record but clear assigned_by
    
    -- Prevent duplicate role assignments
    CONSTRAINT unique_user_role UNIQUE (user_id, role_id)
);

-- Create indexes for foreign keys
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
```

## Indexing Strategy

### Index Types

**Choose appropriate index types:**

```sql
-- B-tree index (default, most common)
CREATE INDEX idx_users_email ON users(email);

-- Partial index (filter on condition)
CREATE INDEX idx_active_users_email ON users(email) WHERE is_active = true AND deleted_at IS NULL;

-- Composite index (multiple columns)
CREATE INDEX idx_users_name_created ON users(last_name, first_name, created_at);

-- Unique index
CREATE UNIQUE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;

-- GIN index (for JSONB, arrays, full-text search)
CREATE INDEX idx_users_metadata ON users USING GIN (metadata);

-- Hash index (equality checks only)
CREATE INDEX idx_users_email_hash ON users USING HASH (email);
```

### When to Index

**✅ Index these:**
- Primary keys (automatic)
- Foreign keys
- Columns used in WHERE clauses
- Columns used in JOIN conditions
- Columns used in ORDER BY
- Columns used for UNIQUE constraints

**❌ Don't index these:**
- Small tables (< 1000 rows)
- Columns with low cardinality (e.g., boolean with 50/50 split)
- Columns rarely used in queries
- Columns with high write frequency and low read frequency

### Index Monitoring

**Monitor index usage:**

```sql
-- Find unused indexes
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY schemaname, tablename;

-- Find missing indexes (high seq scans)
SELECT
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    seq_tup_read / seq_scan AS avg_seq_tup_read
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC
LIMIT 20;
```

## Query Optimization

### Use EXPLAIN ANALYZE

**Analyze query performance:**

```sql
-- ✅ DO: Use EXPLAIN ANALYZE to optimize queries
EXPLAIN ANALYZE
SELECT u.id, u.email, r.name AS role_name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.is_active = true
ORDER BY u.created_at DESC
LIMIT 10;

-- Look for:
-- - Seq Scan → Add index
-- - High cost → Optimize query
-- - Nested Loop with large tables → Ensure proper indexes
```

### Avoid N+1 Queries

**Use JOINs instead of multiple queries:**

```go
// ✅ DO: Single query with JOIN
func (r *UserRepository) GetUsersWithRoles(ctx context.Context) ([]*UserWithRoles, error) {
    query := `
        SELECT 
            u.id, u.email, u.first_name, u.last_name,
            COALESCE(json_agg(
                json_build_object('id', r.id, 'name', r.name)
            ) FILTER (WHERE r.id IS NOT NULL), '[]') AS roles
        FROM users u
        LEFT JOIN user_roles ur ON u.id = ur.user_id
        LEFT JOIN roles r ON ur.role_id = r.id
        WHERE u.is_active = true
        GROUP BY u.id
        ORDER BY u.created_at DESC
    `
    
    rows, err := r.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*UserWithRoles
    for rows.Next() {
        var user UserWithRoles
        var rolesJSON []byte
        
        err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &rolesJSON)
        if err != nil {
            return nil, err
        }
        
        json.Unmarshal(rolesJSON, &user.Roles)
        users = append(users, &user)
    }
    
    return users, nil
}

// ❌ DON'T: N+1 queries
func (r *UserRepository) GetUsersWithRoles(ctx context.Context) ([]*UserWithRoles, error) {
    users, _ := r.GetAllUsers(ctx)  // 1 query
    
    for _, user := range users {
        roles, _ := r.GetUserRoles(ctx, user.ID)  // N queries
        user.Roles = roles
    }
    
    return users, nil
}
```

## Soft Deletes

**Implement soft deletes with deleted_at:**

```sql
-- Migration
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ;
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

-- Repository
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
    
    result, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    if result.RowsAffected() == 0 {
        return ErrUserNotFound
    }
    
    return nil
}

// Include WHERE deleted_at IS NULL in all queries
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
    // ...
}
```

## Best Practices Summary

### ✅ DO

**Connections:**
- Use pgxpool for connection pooling
- Configure max connections appropriately
- Close connections properly
- Use context for cancellation

**Queries:**
- Use parameterized queries ($1, $2)
- Use EXPLAIN ANALYZE to optimize
- Use JOINs to avoid N+1
- Use batch operations for multiple inserts
- Use transactions for multi-step operations

**Migrations:**
- Number migrations sequentially
- Write both up and down migrations
- Add indexes after data is loaded
- Use constraints for data integrity
- Document migrations with comments

**Indexes:**
- Index foreign keys
- Index WHERE clause columns
- Use partial indexes when appropriate
- Monitor index usage
- Remove unused indexes

### ❌ DON'T

**Queries:**
- Use string concatenation for queries (SQL injection)
- Create new connection per request
- Forget to close connections
- Use SELECT * in production
- Ignore query performance

**Migrations:**
- Edit existing migrations after deployment
- Use date-based migration names
- Skip down migrations
- Modify data without backup
- Add NOT NULL without default (on large tables)

**Indexes:**
- Index everything blindly
- Create duplicate indexes
- Skip index monitoring
- Add indexes without EXPLAIN ANALYZE
- Use indexes on small tables

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-15  
**Maintained By:** Backend Team
