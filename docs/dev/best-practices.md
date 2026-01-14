# Venio Best Practices

This document contains framework-specific best practices for all technologies used in Venio.

## Table of Contents

1. [Go Best Practices](#go-best-practices)
2. [Gin Framework](#gin-framework)
3. [pgx Database Driver](#pgx-database-driver)
4. [JWT Authentication](#jwt-authentication)
5. [RBAC Implementation](#rbac-implementation)
6. [Next.js 15](#nextjs-15)
7. [React Best Practices](#react-best-practices)
8. [TypeScript Patterns](#typescript-patterns)
9. [PostgreSQL 18](#postgresql-18)
10. [Redis 8.4](#redis-84)

---

## Go Best Practices

### Use Go 1.25 Features

```go
// ✅ DO: Use sync.WaitGroup.Go() (Go 1.25+)
var wg sync.WaitGroup
wg.Go(func() {
    // Work happens here
    processData()
})
wg.Wait()

// ❌ DON'T: Manual goroutine management
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    processData()
}()
wg.Wait()
```

### Context Management

```go
// ✅ DO: Use context.WithoutCancel() for detached operations (Go 1.21+)
func asyncOperation(parentCtx context.Context) {
    ctx := context.WithoutCancel(parentCtx) // Inherits values, not cancellation
    go func() {
        // Long-running task that shouldn't be cancelled
        performCleanup(ctx)
    }()
}

// ❌ DON'T: Use context.Background() when you have a parent context
func asyncOperation(parentCtx context.Context) {
    go func() {
        ctx := context.Background() // Loses all parent context values!
        performCleanup(ctx)
    }()
}
```

### Error Handling

```go
// ✅ DO: Wrap errors with context using %w
if err := db.Query(ctx, query); err != nil {
    return fmt.Errorf("failed to query users: %w", err)
}

// ✅ DO: Use errors.Is() and errors.As() for error checking
if errors.Is(err, sql.ErrNoRows) {
    return nil, ErrUserNotFound
}

// ❌ DON'T: Return bare errors without context
if err := db.Query(ctx, query); err != nil {
    return err
}
```

### Struct Tags

```go
// ✅ DO: Use consistent struct tags
type User struct {
    ID        int64     `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ❌ DON'T: Inconsistent or missing tags
type User struct {
    ID        int64
    Email     string `json:"email"`
    CreatedAt time.Time `db:"created_at"`
}
```

---

## Gin Framework

### Middleware Order Matters

```go
// ✅ DO: Correct middleware order
router := gin.New()
router.Use(gin.Recovery())        // 1. Panic recovery first
router.Use(middleware.Logger())   // 2. Logging
router.Use(middleware.CORS())     // 3. CORS
router.Use(middleware.RateLimiter()) // 4. Rate limiting
router.Use(middleware.Auth())     // 5. Authentication (when needed)
```

### Context Best Practices

```go
// ✅ DO: Extract context from gin.Context for service calls
func (h *Handler) GetUser(c *gin.Context) {
    ctx := c.Request.Context()
    user, err := h.service.GetUser(ctx, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": user})
}

// ❌ DON'T: Pass gin.Context to service layer
func (h *Handler) GetUser(c *gin.Context) {
    user, err := h.service.GetUser(c, userID) // ❌ Wrong
}
```

### Response Patterns

```go
// ✅ DO: Consistent response structure
type Response struct {
    Data  interface{} `json:"data,omitempty"`
    Error string      `json:"error,omitempty"`
}

func (h *Handler) HandleRequest(c *gin.Context) {
    result, err := h.service.DoWork(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
        return
    }
    c.JSON(http.StatusOK, Response{Data: result})
}
```

### Input Validation

```go
// ✅ DO: Use struct validation tags
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=8"`
}

func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Continue with validated data
}
```

---

## pgx Database Driver

### Connection Pooling

```go
// ✅ DO: Use pgxpool for production
config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
if err != nil {
    return nil, fmt.Errorf("parse config: %w", err)
}

// Configure pool settings
config.MaxConns = 25
config.MinConns = 5
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute
config.HealthCheckPeriod = time.Minute

pool, err := pgxpool.NewWithConfig(context.Background(), config)
```

### Prepared Statements

```go
// ✅ DO: Use named parameters for clarity
const query = `
    SELECT id, email, username
    FROM users
    WHERE email = $1 AND deleted_at IS NULL
`

var user User
err := pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username)

// ✅ DO: Use pgx.NamedArgs for complex queries
const query = `
    UPDATE users
    SET email = @email, updated_at = @updated_at
    WHERE id = @id
`

_, err := pool.Exec(ctx, query, pgx.NamedArgs{
    "email":      newEmail,
    "updated_at": time.Now(),
    "id":         userID,
})
```

### Transactions

```go
// ✅ DO: Use pgx.TxOptions for transaction control
tx, err := pool.BeginTx(ctx, pgx.TxOptions{
    IsoLevel:   pgx.ReadCommitted,
    AccessMode: pgx.ReadWrite,
})
if err != nil {
    return fmt.Errorf("begin tx: %w", err)
}
defer tx.Rollback(ctx) // Safe to call even after commit

// Perform operations
if err := doWork(ctx, tx); err != nil {
    return fmt.Errorf("transaction failed: %w", err)
}

if err := tx.Commit(ctx); err != nil {
    return fmt.Errorf("commit tx: %w", err)
}
```

### Query Performance

```go
// ✅ DO: Use batch operations for multiple inserts
batch := &pgx.Batch{}
for _, user := range users {
    batch.Queue("INSERT INTO users (email, username) VALUES ($1, $2)", user.Email, user.Username)
}
br := pool.SendBatch(ctx, batch)
defer br.Close()

// ❌ DON'T: Loop with individual inserts
for _, user := range users {
    _, err := pool.Exec(ctx, "INSERT INTO users (email, username) VALUES ($1, $2)", user.Email, user.Username)
}
```

---

## JWT Authentication

### Token Generation

```go
// ✅ DO: Use strong signing methods
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id":  user.ID,
    "email":    user.Email,
    "username": user.Username,
    "roles":    user.Roles,
    "iss":      "venio",
    "iat":      time.Now().Unix(),
    "exp":      time.Now().Add(24 * time.Hour).Unix(),
})

tokenString, err := token.SignedString([]byte(jwtSecret))

// ❌ DON'T: Use weak signing methods
token := jwt.NewWithClaims(jwt.SigningMethodNone, claims) // ❌ No security!
```

### Token Validation

```go
// ✅ DO: Validate all JWT claims
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Verify signing method
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(jwtSecret), nil
})

if err != nil {
    return nil, fmt.Errorf("invalid token: %w", err)
}

if !token.Valid {
    return nil, errors.New("token is not valid")
}

// Extract and validate claims
claims, ok := token.Claims.(jwt.MapClaims)
if !ok {
    return nil, errors.New("invalid claims")
}

// Validate issuer
if iss, ok := claims["iss"].(string); !ok || iss != "venio" {
    return nil, errors.New("invalid issuer")
}
```

### Token Storage (Frontend)

```typescript
// ✅ DO: Store in HTTP-only cookies (set by backend)
// Backend sets cookie with httpOnly, secure, sameSite flags

// ✅ DO: Include credentials in fetch requests
const response = await fetch('/api/v1/users', {
  credentials: 'include', // Sends cookies
  headers: {
    'Content-Type': 'application/json',
  },
});

// ❌ DON'T: Store JWT in localStorage (XSS vulnerable)
localStorage.setItem('token', tokenString); // ❌ Security risk!
```

---

## RBAC Implementation

### Permission Checking

```go
// ✅ DO: Check permissions in middleware
func (m *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt64("user_id")

        hasPermission, err := m.rbacService.UserHasPermission(c.Request.Context(), userID, permission)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
            c.Abort()
            return
        }

        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// Usage in routes
authorized := router.Group("/api/v1/admin")
authorized.Use(middleware.RequirePermission("admin:users:write"))
authorized.POST("/users", handler.CreateUser)
```

### Role Hierarchy

```go
// ✅ DO: Cache role permissions
type RBACService struct {
    cache map[int64][]string // userID -> permissions
    mu    sync.RWMutex
}

func (s *RBACService) UserHasPermission(ctx context.Context, userID int64, permission string) (bool, error) {
    s.mu.RLock()
    permissions, cached := s.cache[userID]
    s.mu.RUnlock()

    if !cached {
        // Fetch from DB and cache
        perms, err := s.fetchUserPermissions(ctx, userID)
        if err != nil {
            return false, err
        }
        s.mu.Lock()
        s.cache[userID] = perms
        s.mu.Unlock()
        permissions = perms
    }

    return contains(permissions, permission), nil
}
```

---

## Next.js 15

### App Router Best Practices

```typescript
// ✅ DO: Use Server Components by default
export default async function UsersPage() {
  const users = await getUsers(); // Fetch on server
  return <UserList users={users} />;
}

// ✅ DO: Use 'use client' only when needed
'use client';

import { useState } from 'react';

export default function Counter() {
  const [count, setCount] = useState(0);
  return <button onClick={() => setCount(count + 1)}>{count}</button>;
}
```

### Data Fetching

```typescript
// ✅ DO: Use fetch with Next.js caching
export async function getUser(id: number) {
  const res = await fetch(`${API_URL}/users/${id}`, {
    next: { revalidate: 60 }, // Cache for 60 seconds
  });

  if (!res.ok) {
    throw new Error('Failed to fetch user');
  }

  return res.json();
}

// ✅ DO: Use server actions for mutations
'use server';

export async function createUser(data: FormData) {
  const email = data.get('email');
  // Perform mutation
  revalidatePath('/users');
}
```

### Loading and Error States

```typescript
// ✅ DO: Use loading.tsx for loading states
// app/users/loading.tsx
export default function Loading() {
  return <div>Loading users...</div>;
}

// ✅ DO: Use error.tsx for error boundaries
// app/users/error.tsx
'use client';

export default function Error({ error, reset }: { error: Error; reset: () => void }) {
  return (
    <div>
      <h2>Something went wrong!</h2>
      <button onClick={() => reset()}>Try again</button>
    </div>
  );
}
```

---

## React Best Practices

### Custom Hooks

```typescript
// ✅ DO: Extract reusable logic into custom hooks
function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const data = await getCurrentUser();
        setUser(data);
      } catch (error) {
        console.error('Auth error:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, []);

  return { user, loading };
}

// Usage
function Dashboard() {
  const { user, loading } = useAuth();

  if (loading) return <div>Loading...</div>;
  if (!user) return <div>Not authenticated</div>;

  return <div>Welcome, {user.username}!</div>;
}
```

### Memoization

```typescript
// ✅ DO: Use useMemo for expensive calculations
const sortedUsers = useMemo(() => {
  return users.sort((a, b) => a.name.localeCompare(b.name));
}, [users]);

// ✅ DO: Use useCallback for function props
const handleClick = useCallback((id: number) => {
  console.log('Clicked:', id);
}, []);

// ❌ DON'T: Overuse memoization for simple operations
const doubled = useMemo(() => count * 2, [count]); // ❌ Unnecessary
```

---

## TypeScript Patterns

### Type Safety

```typescript
// ✅ DO: Use strict types
interface User {
  id: number;
  email: string;
  username: string;
  roles: string[];
}

// ✅ DO: Use discriminated unions for variants
type Result<T> =
  | { success: true; data: T }
  | { success: false; error: string };

function handleResult<T>(result: Result<T>) {
  if (result.success) {
    console.log(result.data); // TypeScript knows data exists
  } else {
    console.error(result.error); // TypeScript knows error exists
  }
}

// ❌ DON'T: Use 'any'
function processData(data: any) { // ❌ Loses type safety
  return data.something;
}
```

### API Types

```typescript
// ✅ DO: Define API response types
interface ApiResponse<T> {
  data?: T;
  error?: string;
}

async function fetchUser(id: number): Promise<ApiResponse<User>> {
  const res = await fetch(`/api/users/${id}`);
  return res.json();
}
```

---

## PostgreSQL 18

### Use Modern Features

```sql
-- ✅ DO: Use UUID v7 (timestamp-ordered, PostgreSQL 18+)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ✅ DO: Use SCRAM-SHA-256 authentication (not MD5)
-- In postgresql.conf:
-- password_encryption = 'scram-sha-256'

-- ✅ DO: Use Row-Level Security for multi-tenancy
CREATE POLICY tenant_isolation ON data
    USING (tenant_id = current_setting('app.current_tenant')::int);
```

### Indexing Best Practices

```sql
-- ✅ DO: Use partial indexes for filtered queries
CREATE INDEX idx_active_users ON users (email) WHERE deleted_at IS NULL;

-- ✅ DO: Use BRIN indexes for time-series data
CREATE INDEX idx_logs_created ON logs USING BRIN (created_at);

-- ✅ DO: Use GIN indexes for JSONB
CREATE INDEX idx_user_metadata ON users USING GIN (metadata);
```

---

## Redis 8.4

### Connection Management

```go
// ✅ DO: Use redis/go-redis with connection pooling
rdb := redis.NewClient(&redis.Options{
    Addr:         os.Getenv("REDIS_ADDR"),
    Password:     os.Getenv("REDIS_PASSWORD"),
    DB:           0,
    PoolSize:     10,
    MinIdleConns: 5,
})

// ✅ DO: Use context with timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := rdb.Set(ctx, "key", "value", time.Hour).Err()
```

### Caching Patterns

```go
// ✅ DO: Cache database queries
func (s *UserService) GetUser(ctx context.Context, userID int64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", userID)

    // Try cache first
    cached, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }

    // Cache miss, fetch from DB
    user, err := s.repo.GetUser(ctx, userID)
    if err != nil {
        return nil, err
    }

    // Cache for 1 hour
    data, _ := json.Marshal(user)
    s.redis.Set(ctx, cacheKey, data, time.Hour)

    return user, nil
}
```

---

## Document Revision

| Version | Date | Author | Changes | Source Version |
|---------|------|--------|---------|----------------|
| 1.0.0 | 2026-01-14 | AI Assistant | Initial best practices documentation | - |

## Referenced Documentation

- **Go 1.25:** [Go 1.25 Release Notes](https://go.dev/doc/go1.25) (Released: 2026-01-10)
- **Gin Framework:** [Gin v1.10 Documentation](https://gin-gonic.com/docs/) (Released: 2024-06-20)
- **pgx v5:** [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5) (v5.7.2, Released: 2024-08-05)
- **JWT:** [RFC 7519 - JSON Web Tokens](https://datatracker.ietf.org/doc/html/rfc7519)
- **Next.js 15:** [Next.js 15 Documentation](https://nextjs.org/docs) (Released: 2024-10-21)
- **React 19:** [React 19 Documentation](https://react.dev) (Released: 2024-12-05)
- **TypeScript 5.7:** [TypeScript 5.7 Release Notes](https://www.typescriptlang.org/docs/handbook/release-notes/typescript-5-7.html)
- **PostgreSQL 18:** [PostgreSQL 18 Documentation](https://www.postgresql.org/docs/18/) (Released: 2025-11-14)
- **Redis 8.4:** [Redis 8.4 Commands](https://redis.io/docs/latest/commands/) (Released: 2025-12-15)
