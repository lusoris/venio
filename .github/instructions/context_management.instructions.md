---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/**/*.go"
description: Go Context Management Best Practices
---

# Context Management Guidelines

## Core Principle

**Context Propagation:** Always accept `context.Context` from the caller. Never create `context.Background()` in service or repository methods.

## Why Context Matters

Context carries:
1. **Cancellation signals** - Stop work when request is cancelled
2. **Deadlines/Timeouts** - Prevent resource leaks
3. **Request-scoped values** - User ID, trace IDs, etc.

Creating `context.Background()` breaks the chain and loses these signals.

## Architecture Layers

Request flows through layers, each passing context down:

```
HTTP Request → Handler → Service → Repository → Database
     ↓           ↓          ↓           ↓           ↓
  Context    Context    Context     Context    Context
```

## ✅ DO: Accept Context from Caller

### Handler Layer
```go
// ✅ GOOD: Use request context
func (h *AuthHandler) Login(c *gin.Context) {
    // Gin provides c.Request.Context()
    accessToken, err := h.authService.Login(c.Request.Context(), email, password)
    if err != nil {
        // handle error
    }
}
```

### Service Layer
```go
// ✅ GOOD: Accept context parameter
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
    // Add timeout to existing context
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    user, err := s.userService.GetUserByEmail(ctx, email)
    if err != nil {
        return "", err
    }
    // ...
}
```

### Repository Layer
```go
// ✅ GOOD: Pass context to database
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    row := r.pool.QueryRow(ctx,
        `SELECT id, email, username FROM users WHERE email = $1`,
        email,
    )
    // ...
}
```

## ❌ DON'T: Create New Background Context

```go
// ❌ BAD: Creates new context, loses cancellation
func (s *AuthService) Login(email, password string) (string, error) {
    ctx := context.Background()  // ❌ Wrong!
    user, err := s.userService.GetUserByEmail(ctx, email)
    // ...
}

// ❌ BAD: Ignores parent context
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
    newCtx := context.Background()  // ❌ Wrong! Use ctx parameter
    user, err := s.userService.GetUserByEmail(newCtx, email)
    // ...
}
```

## Context Timeout Patterns

### Adding Timeout to Existing Context
```go
// ✅ GOOD: Extend existing context with timeout
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Add 10 second timeout to operation
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Check if email exists (may take time)
    exists, err := s.repo.Exists(ctx, req.Email)
    if err != nil {
        return nil, err
    }

    // Create user (may take time)
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

### Multiple Operations with Same Context
```go
// ✅ GOOD: Reuse same context for related operations
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // All operations share same context
    user, err := s.userService.GetUserByEmail(ctx, email)
    if err != nil {
        return "", "", err
    }

    roles, err := s.userRoleService.GetUserRoles(ctx, user.ID)
    if err != nil {
        return "", "", err
    }

    token, err := s.generateToken(user, roles)
    return token, "", nil
}
```

## Go 1.21+ Enhancement: context.WithoutCancel()

When you need to detach from parent cancellation:

```go
// ✅ GOOD: Go 1.21+ pattern for cleanup operations
func (s *Service) ProcessRequest(ctx context.Context, data string) error {
    // Main operation respects cancellation
    result, err := s.repo.Process(ctx, data)
    if err != nil {
        return err
    }

    // Cleanup must complete even if request cancelled
    cleanupCtx := context.WithoutCancel(ctx)
    go s.cleanup(cleanupCtx, result)

    return nil
}

// ❌ OLD: Pre-1.21 pattern (deprecated)
cleanupCtx := context.Background()  // ❌ Loses request values
```

## Context Values (Use Sparingly)

### ✅ DO: Use for Request-Scoped Data
```go
// Middleware sets user ID
ctx := context.WithValue(r.Context(), "user_id", userID)

// Handler retrieves it
userID := ctx.Value("user_id").(int64)
```

### ❌ DON'T: Use for Function Parameters
```go
// ❌ BAD: Don't pass function parameters via context
ctx = context.WithValue(ctx, "email", email)  // ❌ Wrong!

// ✅ GOOD: Use explicit parameters
func Login(ctx context.Context, email string) error
```

## Testing with Context

### Unit Tests
```go
func TestUserService_CreateUser(t *testing.T) {
    ctx := context.Background()  // ✅ OK in tests

    // Or with timeout for safety
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    user, err := service.CreateUser(ctx, req)
    assert.NoError(t, err)
}
```

### Integration Tests with Cancellation
```go
func TestLongOperation_Cancellation(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    err := service.LongOperation(ctx)

    // Should return context.DeadlineExceeded
    assert.ErrorIs(t, err, context.DeadlineExceeded)
}
```

## Common Mistakes

### 1. Ignoring Context Cancellation
```go
// ❌ BAD: Doesn't check for cancellation
func (s *Service) Process(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        // Long operation without checking ctx
        processItem(i)
    }
    return nil
}

// ✅ GOOD: Check cancellation periodically
func (s *Service) Process(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()  // Return early if cancelled
        default:
            processItem(i)
        }
    }
    return nil
}
```

### 2. Not Deferring Cancel
```go
// ❌ BAD: Cancel not called if error occurs
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
result, err := doWork(ctx)
if err != nil {
    return err  // ❌ Leak: cancel() never called
}
cancel()

// ✅ GOOD: Always defer cancel
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()  // ✅ Always called
result, err := doWork(ctx)
```

### 3. Context in Struct (Anti-Pattern)
```go
// ❌ BAD: Don't store context in struct
type Service struct {
    ctx context.Context  // ❌ Wrong!
}

// ✅ GOOD: Pass context as parameter
type Service struct {
    repo Repository
}

func (s *Service) DoWork(ctx context.Context) error {
    return s.repo.Query(ctx)
}
```

## Migration Guide

### Existing Code Without Context
If you have code without context:

```go
// Before: No context
func (s *Service) CreateUser(email string) error {
    return s.repo.Create(email)
}

// After: Add context parameter
func (s *Service) CreateUser(ctx context.Context, email string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    return s.repo.Create(ctx, email)
}
```

### Update Callers
```go
// Handler update
func (h *Handler) CreateUser(c *gin.Context) {
    // Before
    err := h.service.CreateUser(email)

    // After
    err := h.service.CreateUser(c.Request.Context(), email)
}
```

## Checklist for Every Function

- [ ] Does this function do I/O? (DB, HTTP, file) → Add context parameter
- [ ] Does this function call other services/repos? → Pass context down
- [ ] Am I creating context.Background()? → Should I accept ctx parameter instead?
- [ ] Did I defer cancel() after WithTimeout/WithCancel?
- [ ] Am I checking ctx.Done() in long loops?

---

## Quick Reference

| Scenario | Pattern |
|----------|---------|
| **Service Method** | `func (s *Service) Do(ctx context.Context, ...) error` |
| **Repository Method** | `func (r *Repo) Query(ctx context.Context, ...) error` |
| **Add Timeout** | `ctx, cancel := context.WithTimeout(ctx, 10*time.Second); defer cancel()` |
| **Database Query** | `r.pool.QueryRow(ctx, sql, args...)` |
| **HTTP Request** | `req, _ := http.NewRequestWithContext(ctx, ...)` |
| **Test** | `ctx := context.Background()` or `ctx, cancel := context.WithTimeout(...)` |

---

**Remember:** Context flows from request to database. Never break the chain with `context.Background()` in service/repository layers.
