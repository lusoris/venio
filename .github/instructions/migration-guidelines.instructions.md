---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "migrations/**/*.sql"
description: Database Migration Guidelines
---

# Database Migration Guidelines

## Core Principle

**Migrations are immutable, sequential, and reversible**: Once deployed, never modify. Always add new migration. Every up must have a down.

## Naming Convention

### ✅ CORRECT Format

```
XXX_descriptive_name.up.sql
XXX_descriptive_name.down.sql
```

Where:
- `XXX` = Sequential number (001, 002, 003, ...)
- `descriptive_name` = Snake_case description of change
- `.up.sql` = Forward migration (apply change)
- `.down.sql` = Reverse migration (undo change)

### Examples

```
001_initial_schema.up.sql
001_initial_schema.down.sql

002_seed_roles_and_permissions.up.sql
002_seed_roles_and_permissions.down.sql

003_add_user_email_verification.up.sql
003_add_user_email_verification.down.sql

004_create_oauth_tokens_table.up.sql
004_create_oauth_tokens_table.down.sql

005_add_user_last_login_at.up.sql
005_add_user_last_login_at.down.sql
```

### ❌ WRONG Names

```
add_column.sql                    # ❌ No number, no up/down
001_migration.up.sql              # ❌ Not descriptive
AddUserTable.up.sql               # ❌ CamelCase instead of snake_case
001-add-user-table.up.sql         # ❌ Hyphens instead of underscores
```

## Migration Structure

### Table Creation (.up.sql)

```sql
-- ✅ CORRECT: Complete table creation with constraints
CREATE TABLE IF NOT EXISTS users (
    -- Primary Key
    id BIGSERIAL PRIMARY KEY,
    
    -- Unique Constraints
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    
    -- Required Fields
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    
    -- Optional Fields
    phone VARCHAR(20),
    avatar_url TEXT,
    
    -- Status Fields
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    
    -- Timestamps (REQUIRED)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes (AFTER table creation)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- Comments (for documentation)
COMMENT ON TABLE users IS 'Core user accounts table';
COMMENT ON COLUMN users.is_email_verified IS 'Tracks whether user has verified their email address';
```

### Table Deletion (.down.sql)

```sql
-- ✅ CORRECT: Clean deletion with CASCADE if needed
DROP TABLE IF EXISTS users CASCADE;

-- Note: Indexes are automatically dropped with the table
```

### Adding Column (.up.sql)

```sql
-- ✅ CORRECT: Add column with constraints
ALTER TABLE users
ADD COLUMN last_login_at TIMESTAMPTZ;

-- Add index if needed for queries
CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON users(last_login_at DESC);

-- Add comment
COMMENT ON COLUMN users.last_login_at IS 'Timestamp of last successful login';
```

### Removing Column (.down.sql)

```sql
-- ✅ CORRECT: Remove column
ALTER TABLE users
DROP COLUMN IF EXISTS last_login_at;

-- Note: Index is automatically dropped
```

### Modifying Column (.up.sql)

```sql
-- ✅ CORRECT: Modify column type with migration strategy
-- Step 1: Add new column
ALTER TABLE users
ADD COLUMN email_new VARCHAR(320);  -- New length

-- Step 2: Copy data
UPDATE users
SET email_new = email;

-- Step 3: Drop old column
ALTER TABLE users
DROP COLUMN email;

-- Step 4: Rename new column
ALTER TABLE users
RENAME COLUMN email_new TO email;

-- Step 5: Add back constraints
ALTER TABLE users
ADD CONSTRAINT users_email_unique UNIQUE (email);

-- Step 6: Recreate index
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

### Modifying Column (.down.sql)

```sql
-- ✅ CORRECT: Reverse the change
ALTER TABLE users
ADD COLUMN email_old VARCHAR(255);

UPDATE users
SET email_old = email;

ALTER TABLE users
DROP COLUMN email;

ALTER TABLE users
RENAME COLUMN email_old TO email;

ALTER TABLE users
ADD CONSTRAINT users_email_unique UNIQUE (email);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

## Foreign Keys

### ✅ CORRECT: With CASCADE Rules

```sql
-- Many-to-Many relationship
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    
    -- Foreign Keys with CASCADE
    CONSTRAINT fk_user_roles_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,  -- Delete user_roles when user deleted
        
    CONSTRAINT fk_user_roles_role 
        FOREIGN KEY (role_id) 
        REFERENCES roles(id) 
        ON DELETE CASCADE,  -- Delete user_roles when role deleted
    
    -- Unique constraint (user can have role only once)
    CONSTRAINT unique_user_role UNIQUE (user_id, role_id),
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes on foreign keys (IMPORTANT for joins)
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
```

### Cascade Options

```sql
-- ON DELETE CASCADE: Delete child records when parent deleted
-- Use for: junction tables, user-owned data

-- ON DELETE RESTRICT: Prevent deletion if children exist
-- Use for: lookup tables, prevent accidental data loss
CONSTRAINT fk_users_country 
    FOREIGN KEY (country_id) 
    REFERENCES countries(id) 
    ON DELETE RESTRICT

-- ON DELETE SET NULL: Set foreign key to NULL when parent deleted
-- Use for: optional relationships
ALTER TABLE posts
ADD CONSTRAINT fk_posts_author 
    FOREIGN KEY (author_id) 
    REFERENCES users(id) 
    ON DELETE SET NULL;
```

## Indexes

### When to Add Indexes

```sql
-- ✅ ALWAYS index:
-- 1. Primary keys (automatic)
-- 2. Foreign keys (for JOIN performance)
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);

-- 3. Unique constraints (automatic if using UNIQUE keyword)
-- But if separate: CREATE UNIQUE INDEX idx_users_email ON users(email);

-- 4. Columns frequently in WHERE clauses
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- 5. Columns frequently in ORDER BY
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- 6. Composite indexes for multi-column queries
CREATE INDEX IF NOT EXISTS idx_users_active_created 
    ON users(is_active, created_at DESC)
    WHERE is_active = TRUE;
```

### Partial Indexes

```sql
-- ✅ CORRECT: Partial index for common queries
-- Only index active users (saves space, faster queries)
CREATE INDEX IF NOT EXISTS idx_users_active_email 
    ON users(email) 
    WHERE is_active = TRUE;

-- Only index unverified emails (for email verification queries)
CREATE INDEX IF NOT EXISTS idx_users_unverified_email 
    ON users(email) 
    WHERE is_email_verified = FALSE;
```

### When NOT to Index

```sql
-- ❌ DON'T index:
-- 1. Small tables (< 1000 rows) - table scan is faster
-- 2. Columns with low cardinality (few distinct values)
--    Example: boolean fields (unless using partial index)
CREATE INDEX idx_users_is_active ON users(is_active);  -- ❌ Only 2 values

-- 3. Columns rarely queried
-- 4. TEXT/BLOB columns (use full-text search instead)
```

## Data Types

### ✅ CORRECT: Use Appropriate Types

```sql
-- IDs: BIGSERIAL (auto-incrementing 64-bit integer)
id BIGSERIAL PRIMARY KEY,

-- Strings: VARCHAR with reasonable max length
email VARCHAR(320) NOT NULL,        -- Email max length: 320
username VARCHAR(50) NOT NULL,      -- Reasonable limit
first_name VARCHAR(100) NOT NULL,
description TEXT,                   -- No length limit

-- Numbers
age INT,                            -- Small integers
price DECIMAL(10, 2),               -- Money: 10 digits, 2 decimal places
quantity BIGINT,                    -- Large integers

-- Booleans
is_active BOOLEAN NOT NULL DEFAULT TRUE,

-- Dates/Times: TIMESTAMPTZ (with timezone)
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
date_of_birth DATE,                 -- Date only (no time)

-- JSON
metadata JSONB,                     -- Binary JSON (faster, indexable)
settings JSON,                      -- Text JSON (if you need exact input)

-- UUIDs (for distributed systems)
uuid UUID NOT NULL DEFAULT gen_random_uuid(),

-- Arrays
tags TEXT[],
role_ids BIGINT[],
```

### ❌ WRONG: Poor Type Choices

```sql
-- ❌ Using SERIAL for IDs (32-bit, can overflow)
id SERIAL PRIMARY KEY,  -- Use BIGSERIAL instead

-- ❌ VARCHAR without length
name VARCHAR,  -- Use TEXT or VARCHAR(n)

-- ❌ TIMESTAMP without timezone
created_at TIMESTAMP,  -- Use TIMESTAMPTZ

-- ❌ JSON instead of JSONB for querying
settings JSON,  -- Use JSONB for indexing/querying

-- ❌ Storing money as FLOAT
price FLOAT,  -- Use DECIMAL(10, 2)
```

## Seed Data

### ✅ CORRECT: Idempotent Seeds

```sql
-- 002_seed_roles_and_permissions.up.sql

-- Insert roles (idempotent with ON CONFLICT)
INSERT INTO roles (name, description, created_at, updated_at)
VALUES 
    ('super_admin', 'Super administrator with all permissions', NOW(), NOW()),
    ('admin', 'Administrator with management permissions', NOW(), NOW()),
    ('user', 'Regular user with basic permissions', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;  -- ✅ Don't fail if already exists

-- Insert permissions
INSERT INTO permissions (name, description, created_at, updated_at)
VALUES 
    ('users:create', 'Create new users', NOW(), NOW()),
    ('users:read', 'View user information', NOW(), NOW()),
    ('users:update', 'Update user information', NOW(), NOW()),
    ('users:delete', 'Delete users', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Link roles to permissions (get IDs from tables)
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT 
    r.id,
    p.id,
    NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'  -- Super admin gets all permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;
```

### Seed Data (.down.sql)

```sql
-- 002_seed_roles_and_permissions.down.sql

-- Delete in reverse order (respect foreign keys)
DELETE FROM role_permissions 
WHERE role_id IN (SELECT id FROM roles WHERE name IN ('super_admin', 'admin', 'user'));

DELETE FROM permissions 
WHERE name LIKE 'users:%';

DELETE FROM roles 
WHERE name IN ('super_admin', 'admin', 'user');
```

## Common Patterns

### Created/Updated Timestamps

```sql
-- Every table should have these
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

-- Trigger for auto-updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### Soft Deletes

```sql
-- Add soft delete columns
ALTER TABLE users
ADD COLUMN deleted_at TIMESTAMPTZ,
ADD COLUMN deleted_by BIGINT REFERENCES users(id);

-- Index for querying non-deleted records
CREATE INDEX IF NOT EXISTS idx_users_not_deleted 
    ON users(id) 
    WHERE deleted_at IS NULL;

-- Queries should filter: WHERE deleted_at IS NULL
```

## Migration Testing

### Before Committing

```bash
# 1. Test up migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/venio_test?sslmode=disable" up

# 2. Test down migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/venio_test?sslmode=disable" down 1

# 3. Test up again (ensure repeatable)
migrate -path migrations -database "postgres://user:pass@localhost:5432/venio_test?sslmode=disable" up

# 4. Check schema
\d users  -- In psql
```

## Migration Checklist

- [ ] **Naming**: Sequential number + descriptive snake_case name
- [ ] **Pairs**: Every .up.sql has matching .down.sql
- [ ] **Idempotent**: Use IF NOT EXISTS, IF EXISTS, ON CONFLICT
- [ ] **Data Types**: BIGSERIAL for IDs, TIMESTAMPTZ for dates, VARCHAR(n) for strings
- [ ] **Constraints**: NOT NULL, UNIQUE, CHECK constraints defined
- [ ] **Foreign Keys**: With appropriate CASCADE/RESTRICT rules
- [ ] **Indexes**: On foreign keys, WHERE clauses, ORDER BY columns
- [ ] **Timestamps**: created_at, updated_at on all tables
- [ ] **Comments**: Document purpose of table/columns
- [ ] **Tested**: Run up, down, up locally before commit
- [ ] **Reversible**: Down migration completely undoes up migration

## Common Mistakes

### ❌ DON'T Do This

```sql
-- 1. Modifying existing migrations
-- ❌ Never edit 001_initial_schema.up.sql after deployed
-- ✅ Create 005_fix_user_email_length.up.sql

-- 2. Missing IF NOT EXISTS
CREATE TABLE users (...);  -- ❌ Fails if table exists
CREATE TABLE IF NOT EXISTS users (...);  -- ✅ Idempotent

-- 3. Missing indexes on foreign keys
FOREIGN KEY (user_id) REFERENCES users(id)  -- ❌ No index
-- ✅ Always add: CREATE INDEX idx_table_user_id ON table(user_id);

-- 4. No timestamps
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    title TEXT
);  -- ❌ Missing created_at, updated_at

-- 5. Wrong timestamp type
created_at TIMESTAMP  -- ❌ No timezone
created_at TIMESTAMPTZ  -- ✅ With timezone

-- 6. Incomplete down migration
-- up: CREATE TABLE + indexes + constraints + seeds
-- down: DROP TABLE  -- ❌ Should match complexity of up

-- 7. Non-reversible migrations
ALTER TABLE users DROP COLUMN password;  -- ❌ Can't restore data in down
-- ✅ Consider soft deletes or data migration strategy
```

## Complex Migration Example

### Phase 1: Add New Column (.up.sql)

```sql
-- 006_add_user_email_verification.up.sql

-- Add columns
ALTER TABLE users
ADD COLUMN is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN email_verification_token VARCHAR(255),
ADD COLUMN email_verification_token_expires_at TIMESTAMPTZ,
ADD COLUMN email_verified_at TIMESTAMPTZ;

-- Add index for verification queries
CREATE INDEX IF NOT EXISTS idx_users_email_verification_token 
    ON users(email_verification_token) 
    WHERE email_verification_token IS NOT NULL;

-- Add partial index for unverified users
CREATE INDEX IF NOT EXISTS idx_users_unverified 
    ON users(email) 
    WHERE is_email_verified = FALSE;

-- Comments
COMMENT ON COLUMN users.is_email_verified IS 'Whether user has verified their email address';
COMMENT ON COLUMN users.email_verification_token IS 'Token sent via email for verification';
COMMENT ON COLUMN users.email_verification_token_expires_at IS 'Expiration timestamp for verification token';
COMMENT ON COLUMN users.email_verified_at IS 'Timestamp when email was verified';
```

### Phase 1: Reverse (.down.sql)

```sql
-- 006_add_user_email_verification.down.sql

-- Drop columns (indexes dropped automatically)
ALTER TABLE users
DROP COLUMN IF EXISTS is_email_verified,
DROP COLUMN IF EXISTS email_verification_token,
DROP COLUMN IF EXISTS email_verification_token_expires_at,
DROP COLUMN IF EXISTS email_verified_at;
```

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial migration guidelines |

**Remember**: Migrations are code. They must be reviewed, tested, and versioned. Never modify deployed migrations - always create new ones.
