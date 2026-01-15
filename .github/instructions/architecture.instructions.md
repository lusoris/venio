---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/**/*.go"
description: Clean Architecture Best Practices for Venio
---

# Clean Architecture Guidelines

## Project Structure

Venio follows Clean Architecture with clear separation of concerns:

```
internal/
├── api/                    # HTTP Layer (Presentation)
│   ├── handlers/          # HTTP handlers (thin, delegate to services)
│   ├── middleware/        # HTTP middleware (auth, logging, metrics)
│   └── routes.go          # Route setup & dependency injection
├── services/              # Business Logic Layer
│   ├── auth_service.go
│   ├── user_service.go
│   └── ...
├── repositories/          # Data Access Layer
│   ├── user_repository.go
│   ├── role_repository.go
│   └── ...
├── models/                # Domain Models (shared across layers)
├── config/                # Configuration
├── database/              # Database connection pool
├── logger/                # Logging infrastructure
└── redis/                 # Redis client
```

## Layer Responsibilities

### 1. Handlers (Presentation Layer)

**Purpose:** Handle HTTP requests and responses

**Responsibilities:**
- Parse HTTP requests (JSON binding, query params)
- Call service methods
- Return HTTP responses
- Handle HTTP-specific errors (400, 401, 403, 404, 500)
- **NO** business logic
- **NO** database access

**Example:**
```go
// ✅ DO: Thin handler, delegates to service
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid ID",
            Message: "User ID must be a number",
        })
        return
    }

    user, err := h.userService.GetUser(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{
            Error: "Not Found",
            Message: "User not found",
        })
        return
    }

    c.JSON(http.StatusOK, user)
}

// ❌ DON'T: Business logic in handler
func (h *UserHandler) GetUser(c *gin.Context) {
    // ❌ Validation logic belongs in service
    if user.Email == "" || !isValidEmail(user.Email) {
        return errors.New("invalid email")
    }
    
    // ❌ Direct database access belongs in repository
    row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
}
```

**Handler Structure:**
```go
type UserHandler struct {
    userService services.UserService  // ✅ Only depends on service interface
}

func NewUserHandler(userService services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}
```

### 2. Services (Business Logic Layer)

**Purpose:** Implement business logic and orchestrate operations

**Responsibilities:**
- Validation (email format, password strength, etc.)
- Business rules (user must be active, role assignments, etc.)
- Orchestration (call multiple repositories if needed)
- Error handling and wrapping
- **NO** HTTP concerns (status codes, headers)
- **NO** direct SQL queries

**Example:**
```go
// ✅ DO: Business logic in service
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // ✅ Validation
    if err := s.validateEmail(req.Email); err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }
    
    if err := s.validatePassword(req.Password); err != nil {
        return nil, fmt.Errorf("invalid password: %w", err)
    }
    
    // ✅ Check for duplicates (business rule)
    existing, _ := s.repo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }
    
    // ✅ Hash password (business logic)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // ✅ Call repository
    user := &User{
        Email:    req.Email,
        Password: string(hashedPassword),
        IsActive: true,
    }
    
    id, err := s.repo.Create(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    user.ID = id
    return user, nil
}

// ❌ DON'T: HTTP concerns in service
func (s *UserService) CreateUser(ctx context.Context, c *gin.Context) {
    // ❌ Don't accept gin.Context
    // ❌ Don't return HTTP status codes
    return http.StatusCreated, user
}
```

**Service Structure:**
```go
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id int64) (*User, error)
    UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
    DeleteUser(ctx context.Context, id int64) error
    ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
}

type userService struct {
    repo repositories.UserRepository  // ✅ Depends on repository interface
}

func NewUserService(repo repositories.UserRepository) UserService {
    return &userService{repo: repo}
}
```

### 3. Repositories (Data Access Layer)

**Purpose:** Persist and retrieve data

**Responsibilities:**
- SQL queries (SELECT, INSERT, UPDATE, DELETE)
- Database connection management
- Transaction handling
- Error handling for database-specific errors
- **NO** business logic
- **NO** validation (beyond database constraints)

**Example:**
```go
// ✅ DO: Clean data access
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) (int64, error) {
    query := `
        INSERT INTO users (email, username, password, first_name, last_name, is_active)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    
    var id int64
    err := r.pool.QueryRow(ctx, query,
        user.Email,
        user.Username,
        user.Password,
        user.FirstName,
        user.LastName,
        user.IsActive,
    ).Scan(&id)
    
    if err != nil {
        return 0, fmt.Errorf("failed to create user: %w", err)
    }
    
    return id, nil
}

// ❌ DON'T: Business logic in repository
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) (int64, error) {
    // ❌ Validation belongs in service
    if user.Email == "" {
        return 0, errors.New("email is required")
    }
    
    // ❌ Password hashing belongs in service
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
    user.Password = string(hashedPassword)
    
    // Database operation...
}
```

**Repository Structure:**
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) (int64, error)
    FindByID(ctx context.Context, id int64) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, limit, offset int) ([]*User, error)
}

type PostgresUserRepository struct {
    pool *pgxpool.Pool  // ✅ Direct database access
}

func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
    return &PostgresUserRepository{pool: pool}
}
```

## Dependency Flow

**Direction:** Outer layers depend on inner layers

```
Handler → Service → Repository → Database
  ↓         ↓          ↓
 HTTP    Business    Data
        Logic      Access
```

**Rules:**
- Handlers depend on Service interfaces (not implementations)
- Services depend on Repository interfaces (not implementations)
- Repositories depend on database driver (pgxpool)
- **Never:** Inner layers depending on outer layers

**Dependency Injection:**
```go
// routes.go - Assemble dependencies
func SetupRouter(cfg *config.Config, db *database.DB, redis *redis.Client) *gin.Engine {
    // Create repositories (innermost layer)
    userRepo := repositories.NewPostgresUserRepository(db.Pool())
    roleRepo := repositories.NewRoleRepository(db.Pool())
    
    // Create services (middle layer)
    userService := services.NewUserService(userRepo)
    roleService := services.NewRoleService(roleRepo)
    
    // Create handlers (outer layer)
    userHandler := handlers.NewUserHandler(userService)
    roleHandler := handlers.NewRoleHandler(roleService)
    
    // Setup routes
    router := gin.Default()
    api := router.Group("/api/v1")
    api.GET("/users/:id", userHandler.GetUser)
    api.GET("/roles", roleHandler.ListRoles)
    
    return router
}
```

## Interface Usage

**Always code against interfaces, not implementations**

```go
// ✅ DO: Accept interfaces
type UserService interface {
    GetUser(ctx context.Context, id int64) (*User, error)
}

type userService struct {
    repo repositories.UserRepository  // ✅ Interface
}

// ❌ DON'T: Accept concrete types
type userService struct {
    repo *repositories.PostgresUserRepository  // ❌ Concrete type
}
```

**Benefits:**
- Easy to test (mock interfaces)
- Easy to swap implementations (e.g., Redis cache layer)
- Loose coupling

## Error Handling

### Service Layer Errors

**Wrap errors with context:**
```go
// ✅ DO: Wrap errors with context
user, err := s.repo.FindByEmail(ctx, email)
if err != nil {
    return nil, fmt.Errorf("failed to find user by email %s: %w", email, err)
}

// ❌ DON'T: Return raw errors
user, err := s.repo.FindByEmail(ctx, email)
if err != nil {
    return nil, err  // ❌ No context
}
```

### Custom Errors

**Define domain errors:**
```go
// services/errors.go
var (
    ErrUserNotFound       = errors.New("user not found")
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInactiveUser       = errors.New("user account is inactive")
)

// Service usage
func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("database error: %w", err)
    }
    return user, nil
}
```

### Handler Error Mapping

**Map service errors to HTTP status codes:**
```go
// ✅ DO: Map errors to HTTP status
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.userService.GetUser(c.Request.Context(), id)
    if err != nil {
        switch {
        case errors.Is(err, services.ErrUserNotFound):
            c.JSON(http.StatusNotFound, ErrorResponse{
                Error: "Not Found",
                Message: "User not found",
            })
        case errors.Is(err, services.ErrInvalidInput):
            c.JSON(http.StatusBadRequest, ErrorResponse{
                Error: "Bad Request",
                Message: err.Error(),
            })
        default:
            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Error: "Internal Server Error",
                Message: "Something went wrong",
            })
        }
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

## Testing Strategy

### Repository Tests

**Test SQL queries directly:**
```go
func TestUserRepository_Create(t *testing.T) {
    // Use testcontainers for real PostgreSQL
    pool := setupTestDB(t)
    repo := repositories.NewPostgresUserRepository(pool)
    
    user := &models.User{
        Email: "test@example.com",
        Password: "hashed",
    }
    
    id, err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.Greater(t, id, int64(0))
}
```

### Service Tests

**Mock repository layer:**
```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := services.NewUserService(mockRepo)
    
    // Setup mock expectations
    mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(int64(1), nil)
    
    // Test
    user, err := service.CreateUser(context.Background(), &CreateUserRequest{
        Email: "test@example.com",
        Password: "SecurePass123!",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    mockRepo.AssertExpectations(t)
}
```

### Handler Tests

**Mock service layer:**
```go
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, id int64) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestUserHandler_GetUser(t *testing.T) {
    mockService := new(MockUserService)
    handler := handlers.NewUserHandler(mockService)
    
    router := gin.Default()
    router.GET("/users/:id", handler.GetUser)
    
    // Setup mock
    mockService.On("GetUser", mock.Anything, int64(1)).Return(&User{
        ID: 1,
        Email: "test@example.com",
    }, nil)
    
    // Test
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users/1", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}
```

## Context Management

**Always propagate context through all layers:**

```go
// ✅ DO: Pass context through all layers
func (h *UserHandler) GetUser(c *gin.Context) {
    ctx := c.Request.Context()  // ✅ Get context from HTTP request
    user, err := h.userService.GetUser(ctx, id)
    // ...
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)  // ✅ Pass context to repository
    // ...
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    row := r.pool.QueryRow(ctx, query, id)  // ✅ Pass context to database
    // ...
}

// ❌ DON'T: Create new context or ignore it
func (s *UserService) GetUser(id int64) (*User, error) {
    ctx := context.Background()  // ❌ Creates new context, loses cancellation
    user, err := s.repo.FindByID(ctx, id)
    // ...
}
```

## Common Patterns

### Pagination

**Consistent pagination across services:**
```go
// Service interface
type UserService interface {
    ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
    CountUsers(ctx context.Context) (int64, error)
}

// Handler with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    users, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{...})
        return
    }
    
    total, _ := h.userService.CountUsers(c.Request.Context())
    
    c.JSON(http.StatusOK, PaginatedResponse{
        Data:   users,
        Total:  total,
        Limit:  limit,
        Offset: offset,
    })
}
```

### Transactions

**Handle transactions in service layer:**
```go
// Service with transaction
func (s *UserService) CreateUserWithRole(ctx context.Context, req *CreateUserRequest, roleID int64) error {
    // Begin transaction
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)  // Rollback if not committed
    
    // Create user
    userID, err := s.repo.CreateTx(ctx, tx, user)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    // Assign role
    err = s.userRoleRepo.AssignRoleTx(ctx, tx, userID, roleID)
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

## Best Practices Summary

### ✅ DO

**Handlers:**
- Parse HTTP requests
- Call service methods
- Return HTTP responses
- Handle HTTP-specific errors

**Services:**
- Implement business logic
- Validate input
- Orchestrate operations
- Wrap errors with context

**Repositories:**
- Execute SQL queries
- Handle database errors
- Return domain models

**General:**
- Use interfaces for dependencies
- Pass context through all layers
- Test each layer independently
- Keep layers thin and focused

### ❌ DON'T

**Handlers:**
- Implement business logic
- Access database directly
- Perform validation (beyond HTTP parsing)

**Services:**
- Accept gin.Context
- Return HTTP status codes
- Write SQL queries

**Repositories:**
- Implement business logic
- Validate business rules
- Hash passwords or format data

**General:**
- Mix layer responsibilities
- Create circular dependencies
- Ignore context
- Return raw errors without wrapping

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-15  
**Maintained By:** Backend Team
