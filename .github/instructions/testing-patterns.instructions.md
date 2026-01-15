---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "**/*_test.go"
description: Testing Patterns & Mock Creation Guidelines
---

# Testing Patterns & Best Practices

## Core Principle

**Write tests BEFORE finding bugs, not after.** Every new service, repository, or handler must have comprehensive unit tests with mocks.

## Mock Creation Rules

### 1. Mock Repository Interfaces

**CRITICAL**: Mocks MUST match the EXACT interface signature from the repository file.

#### ❌ WRONG: Guessing the signature
```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.User), args.Error(1)  // ❌ Assumes pointer return
}
```

#### ✅ CORRECT: Read the actual interface first
```go
// Step 1: Read internal/repositories/user_repository.go to see actual signature
// GetByID(ctx context.Context, id int64) (*models.User, error)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}
```

### 2. Return Type Mapping

**Repository interfaces return different types - ALWAYS check the actual file:**

| Pattern | Example | Mock Return Type |
|---------|---------|------------------|
| Single pointer | `GetByID(...) (*Model, error)` | `*models.User` |
| Slice of values | `List(...) ([]Model, int64, error)` | `[]models.Role` |
| Slice of pointers | `GetAll(...) ([]*Model, error)` | `[]*models.Permission` |

#### Example: Slice of Values (not pointers)
```go
// Repository interface:
// List(ctx context.Context, limit, offset int) ([]models.Role, int64, error)

func (m *MockRoleRepository) List(ctx context.Context, limit, offset int) ([]models.Role, int64, error) {
    args := m.Called(ctx, limit, offset)
    return args.Get(0).([]models.Role), args.Get(1).(int64), args.Error(2)
    // ✅ []models.Role, NOT []*models.Role
}
```

### 3. Complete Interface Implementation

**Every method in the interface MUST be mocked, even if unused in current tests.**

#### ❌ WRONG: Partial mock
```go
type MockPermissionRepository struct {
    mock.Mock
}

func (m *MockPermissionRepository) GetByID(...) {...}
func (m *MockPermissionRepository) Create(...) {...}
// ❌ Missing: AssignToRole, RemoveFromRole
```

#### ✅ CORRECT: Full interface
```go
// Step 1: Read PermissionRepository interface
// Step 2: Implement ALL 9 methods

type MockPermissionRepository struct {
    mock.Mock
}

func (m *MockPermissionRepository) GetByID(ctx context.Context, id int64) (*models.Permission, error) {...}
func (m *MockPermissionRepository) GetByName(ctx context.Context, name string) (*models.Permission, error) {...}
func (m *MockPermissionRepository) Create(ctx context.Context, req *models.CreatePermissionRequest) (*models.Permission, error) {...}
func (m *MockPermissionRepository) Update(ctx context.Context, id int64, req *models.UpdatePermissionRequest) (*models.Permission, error) {...}
func (m *MockPermissionRepository) Delete(ctx context.Context, id int64) error {...}
func (m *MockPermissionRepository) List(ctx context.Context, limit, offset int) ([]models.Permission, int64, error) {...}
func (m *MockPermissionRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Permission, error) {...}
func (m *MockPermissionRepository) AssignToRole(ctx context.Context, roleID, permissionID int64) error {...}
func (m *MockPermissionRepository) RemoveFromRole(ctx context.Context, roleID, permissionID int64) error {...}
```

## Service Test Patterns

### 1. Understanding Service Call Chains

**CRITICAL**: Services often call multiple repository methods. Mock ALL of them.

#### Example: Delete Service Pattern
```go
// Service code (user_service.go):
func (s *userService) Delete(ctx context.Context, id int64) error {
    // 1. Check if user exists
    user, err := s.userRepository.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    // 2. Actually delete
    return s.userRepository.Delete(ctx, id)
}

// ❌ WRONG: Only mock Delete
mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

// ✅ CORRECT: Mock the ENTIRE call chain
mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&models.User{ID: 1}, nil)
mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)
```

### 2. Test Structure Template

```go
// Test file naming: {service_name}_test.go
// Example: user_service_test.go

package services

import (
    "context"
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/lusoris/venio/internal/models"
)

// Step 1: Create mock repository (implement ALL interface methods)
type MockUserRepository struct {
    mock.Mock
}

// Step 2: Implement interface methods
func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

// ... (implement ALL other methods)

// Step 3: Write tests for each service method
// Pattern: Test{ServiceName}_{MethodName}_{Scenario}

// ✅ Success path
func TestUserService_GetByID_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    expectedUser := &models.User{
        ID:    1,
        Email: "test@example.com",
    }

    mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

    user, err := service.GetByID(context.Background(), 1)

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, expectedUser.ID, user.ID)
    mockRepo.AssertExpectations(t)
}

// ✅ Error path
func TestUserService_GetByID_NotFound(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

    user, err := service.GetByID(context.Background(), 999)

    assert.Error(t, err)
    assert.Nil(t, user)
    mockRepo.AssertExpectations(t)
}

// ✅ Validation error path
func TestUserService_GetByID_InvalidID(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    user, err := service.GetByID(context.Background(), 0)

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "invalid")
    // No repository calls should happen for validation errors
    mockRepo.AssertExpectations(t)
}
```

### 3. Test Coverage Requirements

For each service method, write AT LEAST:
- ✅ 1 success test
- ✅ 1 validation error test (if applicable)
- ✅ 1 not found test (for Get/Update/Delete operations)
- ✅ 1 duplicate test (for Create operations)
- ✅ 1 repository error test

Example for `Create` method:
```go
func TestUserService_Create_Success(t *testing.T) {...}
func TestUserService_Create_InvalidEmail(t *testing.T) {...}
func TestUserService_Create_DuplicateEmail(t *testing.T) {...}
func TestUserService_Create_RepositoryError(t *testing.T) {...}
```

### 4. Request Object Validation

**CRITICAL**: When testing services that accept request objects, provide ALL required fields.

#### ❌ WRONG: Missing required fields
```go
req := &models.CreateUserRequest{
    Email:    "test@example.com",
    Password: "SecurePass123!",
}
// ❌ Missing: Username, FirstName, LastName
```

#### ✅ CORRECT: Complete request
```go
req := &models.CreateUserRequest{
    Email:     "test@example.com",
    Username:  "testuser",
    FirstName: "Test",
    LastName:  "User",
    Password:  "SecurePass123!",
}
// ✅ All required fields present
```

### 5. Mock Setup Order Matters

```go
// ✅ CORRECT: Setup mocks in the order they're called
mockRepo.On("Exists", mock.Anything, req.Email).Return(false, nil)  // Called first
mockRepo.On("Create", mock.Anything, mock.Anything).Return(int64(1), nil)  // Called second

// ❌ WRONG: Out of order or missing
mockRepo.On("Create", mock.Anything, mock.Anything).Return(int64(1), nil)
// Missing Exists() call setup - test will fail
```

## Handler Test Patterns

### 1. HTTP Test Structure

```go
func TestAuthHandler_Register_Success(t *testing.T) {
    // Setup
    mockService := new(MockAuthService)
    handler := NewAuthHandler(mockService)

    // Create request
    reqBody := `{"email":"test@example.com","username":"testuser","first_name":"Test","last_name":"User","password":"SecurePass123!"}`
    req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    // Mock service response
    mockService.On("Register", mock.Anything, mock.Anything).Return(&models.User{ID: 1}, nil)

    // Execute
    w := httptest.NewRecorder()
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/register", handler.Register)
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

### 2. Response Validation

```go
// Validate response body structure
var response map[string]interface{}
err := json.Unmarshal(w.Body.Bytes(), &response)
assert.NoError(t, err)
assert.Contains(t, response, "user")
assert.Contains(t, response, "token")
```

## Middleware Test Patterns

```go
func TestRateLimiter_Middleware_Success(t *testing.T) {
    // Setup limiter
    limiter := ratelimit.NewMemoryLimiter(ratelimit.Config{
        RequestsPerWindow: 10,
        WindowDuration:    time.Minute,
    })

    // Create test router
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.Use(middleware.RateLimitMiddleware(limiter))
    router.GET("/test", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })

    // Execute
    req, _ := http.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Common Pitfalls

### ❌ Don't Do This
```go
// 1. Don't mock concrete types
var userService UserService  // ❌ Interface
userService = &userServiceImpl{}  // ❌ Concrete implementation

// 2. Don't skip error checks in tests
user, _ := service.GetByID(ctx, 1)  // ❌ Ignoring error

// 3. Don't use hardcoded values without context
mockRepo.On("GetByID", mock.Anything, int64(1)).Return(user, nil)  // ❌ What is 1?

// 4. Don't forget to check mock expectations
// Missing: mockRepo.AssertExpectations(t)  // ❌
```

### ✅ Do This
```go
// 1. Use interfaces for testability
var service UserService = NewUserService(mockRepo)

// 2. Always check errors
user, err := service.GetByID(ctx, 1)
assert.NoError(t, err)

// 3. Use descriptive constants
const testUserID = int64(1)
mockRepo.On("GetByID", mock.Anything, testUserID).Return(user, nil)

// 4. Always assert expectations
defer mockRepo.AssertExpectations(t)
```

## Test Naming Conventions

```
Test{Type}_{Method}_{Scenario}

Examples:
- TestUserService_Register_Success
- TestUserService_Register_DuplicateEmail
- TestAuthHandler_Login_InvalidCredentials
- TestRateLimiter_Allow_ExceedsLimit
```

## Running Tests

```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./internal/services/... -v

# Run with coverage
go test ./internal/services/... -v -cover

# Run single test
go test ./internal/services/... -v -run TestUserService_Register_Success
```

## Test Coverage Goals

| Package | Minimum Coverage | Target Coverage |
|---------|-----------------|-----------------|
| Services | 60% | 70%+ |
| Handlers | 50% | 60%+ |
| Middleware | 40% | 50%+ |
| Repositories | Optional | Integration tests preferred |

## Checklist for New Tests

- [ ] Read actual interface from repository file
- [ ] Mock ALL interface methods
- [ ] Match exact return types (pointers vs values)
- [ ] Test success path
- [ ] Test validation errors
- [ ] Test not found scenarios
- [ ] Test duplicate/conflict scenarios
- [ ] Test repository errors
- [ ] Mock complete call chains
- [ ] Assert expectations at end
- [ ] Use descriptive test names
- [ ] Provide complete request objects

## Example: Complete Test File

See `internal/services/user_service_test.go` for a complete example with:
- 52+ test functions covering all services
- Complete mock implementations for all repository interfaces (9+ methods each)
- Success, error, and edge case coverage
- Proper mock setup and assertions

## Test Credential Security & Prevention Pattern

### ❌ The Hardcoded Credential Problem

**CRITICAL SECURITY ISSUE**: Hardcoded secrets in tests are flagged by security scanners and violate best practices.

```go
// ❌ WRONG: Hardcoded credentials in tests (Snyk will flag these)
func TestAuthService_Login_Success(t *testing.T) {
    cfg := &config.Config{
        JWT: config.JWTConfig{Secret: "SecurePass123!"},  // ❌ Hardcoded!
    }
    authService := NewAuthService(mockUserService, cfg)
    
    // ...
}

// ❌ WRONG: Hardcoded passwords in fixtures
func TestUserService_Create_Success(t *testing.T) {
    req := &models.CreateUserRequest{
        Email:    "test@example.com",
        Password: "SecurePass123!",  // ❌ Hardcoded!
    }
    
    // ...
}

// ❌ WRONG: Hardcoded tokens in handlers
func TestAuthHandler_VerifyEmail_Success(t *testing.T) {
    token := "abcd1234abcd1234abcd1234abcd1234"  // ❌ Hardcoded!
    
    // ...
}
```

**Snyk Impact**: Each hardcoded credential in a test file is flagged as "Use of Hardcoded Credentials" (low-severity, but accumulates). When you have multiple test files with hardcoded secrets, you can quickly hit 20+ findings.

### ✅ Solution: Dynamic Credential Generation

Use time-based, dynamically generated credentials that:
- Avoid hardcoding ANY secrets
- Change every test run (time-based uniqueness)
- Meet minimum length requirements
- Are consistent within a test function

#### Pattern 1: JWT Secret Helper (config level)

```go
// internal/services/auth_service_test.go

// ✅ CORRECT: Dynamic config generator
func newTestConfig() *config.Config {
    return &config.Config{
        JWT: config.JWTConfig{
            // Generate 32+ character secret using time-based uniqueness
            Secret: fmt.Sprintf("auth-test-secret-%d-must-be-32chars", time.Now().UnixNano()),
        },
    }
}

// Usage in tests
func TestAuthService_GenerateEmailVerificationToken_Success(t *testing.T) {
    cfg := newTestConfig()  // ✅ Dynamic, unique secret
    authService := NewAuthService(mockUserService, cfg)
    
    token, err := authService.GenerateEmailVerificationToken(context.Background(), 1)
    
    assert.NoError(t, err)
    assert.Len(t, token, 64)  // Secure random 32 bytes = 64 hex chars
}

// ✅ CORRECT: Multiple tests can use same helper
func TestAuthService_ValidateToken_Success(t *testing.T) {
    cfg := newTestConfig()  // ✅ Different secret for different test
    authService := NewAuthService(mockUserService, cfg)
    
    // ...
}
```

#### Pattern 2: Password Helper

```go
// internal/services/user_service_test.go

// ✅ CORRECT: Dynamic password generator
func testPassword() string {
    return fmt.Sprintf("pw-secure-%d", time.Now().UnixNano())
}

// Usage in tests
func TestUserService_Create_Success(t *testing.T) {
    req := &models.CreateUserRequest{
        Email:     "test@example.com",
        Username:  "testuser",
        FirstName: "Test",
        LastName:  "User",
        Password:  testPassword(),  // ✅ Unique password per test run
    }
    
    user, err := service.Create(context.Background(), req)
    
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```

#### Pattern 3: JWT Token Helper (handler level)

```go
// internal/api/handlers/auth_handler_test.go

// ✅ CORRECT: Dynamic token generator for handler tests
func testSecret() string {
    // Must be 32+ characters for JWT
    return fmt.Sprintf("access/refresh-%d-test-secret-for-handlers", time.Now().UnixNano())
}

// Usage for authorization headers
func TestAuthHandler_GetUser_Success(t *testing.T) {
    mockService := new(MockAuthService)
    handler := NewAuthHandler(mockService)
    
    secret := testSecret()
    token := generateTestJWT(secret, int64(1))  // Helper to create token
    
    req, _ := http.NewRequest("GET", "/profile", nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    
    // ...
}
```

### ✅ Email Verification Credential Pattern

**Important**: Email verification tokens are secure random 32-byte values (64 hex chars minimum):

```go
// ✅ CORRECT: Proper verification token length
func TestAuthHandler_VerifyEmail_Success(t *testing.T) {
    mockService := new(MockAuthService)
    handler := NewAuthHandler(mockService)
    
    // Generate proper 64-character token (32 bytes in hex)
    token := generateSecureRandomToken(32)  // Produces 64-char hex string
    
    reqBody := fmt.Sprintf(`{"token":"%s"}`, token)
    req, _ := http.NewRequest("POST", "/verify-email", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    mockService.On("VerifyEmail", mock.Anything, token).Return(nil)
    
    // Execute
    w := httptest.NewRecorder()
    handler.VerifyEmail(&gin.Context{Writer: w, Request: req})
    
    assert.Equal(t, http.StatusOK, w.Code)
}

// ❌ WRONG: Token too short
func TestAuthHandler_VerifyEmail_InvalidToken(t *testing.T) {
    token := "short-token"  // ❌ Only 11 chars, fails min=64 validation
    
    // This will fail validation before reaching service
}

// ✅ CORRECT: Test validation properly
func TestAuthHandler_VerifyEmail_InvalidTokenLength(t *testing.T) {
    mockService := new(MockAuthService)
    handler := NewAuthHandler(mockService)
    
    token := "short"  // ❌ Intentionally too short for negative test
    
    reqBody := fmt.Sprintf(`{"token":"%s"}`, token)
    req, _ := http.NewRequest("POST", "/verify-email", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    handler.VerifyEmail(&gin.Context{Writer: w, Request: req})
    
    assert.Equal(t, http.StatusBadRequest, w.Code)
    // Service should NOT be called - validation catches it first
    mockService.AssertNotCalled(t, "VerifyEmail")
}
```

### Test Credential Best Practices

| Type | Pattern | Min Length | Example |
|------|---------|-----------|---------|
| JWT Secret | `fmt.Sprintf("secret-%d", time.Now().UnixNano())` | 32 chars | `secret-1673456789123456789` |
| Password | `fmt.Sprintf("pw-%d", time.Now().UnixNano())` | 8 chars | `pw-1673456789123456789` |
| Verification Token | `generateSecureRandomToken(32)` | 64 chars | `a1b2c3d4e5f6...` (64 hex) |
| JWT Token | Generated with `jwt.NewWithClaims()` | 32 chars | `eyJhbGc...` (standard JWT) |

### Snyk Integration & Verification

**After updating tests with dynamic credentials:**

```bash
# Scan the workspace for hardcoded credentials
snyk code scan .

# Expected improvement:
# BEFORE: 27 issues (multiple hardcoded secrets in test files)
# AFTER: ~5-10 issues (only legitimate test fixtures like emails/usernames)
```

**Acceptable Low-Severity Findings** (after cleanup):
- Test email addresses like "test@example.com" (not secrets, just data)
- Test usernames like "testuser" (not secrets, just data)

**UNACCEPTABLE Findings** (must fix):
- Any hardcoded JWT secrets
- Any hardcoded passwords
- Any hardcoded API keys
- Any hardcoded tokens

### Mock Repository Update Pattern

**CRITICAL**: When adding new repository methods, update ALL test mocks immediately:

```go
// internal/services/user_service_test.go
func (m *MockUserRepository) GetByVerificationToken(ctx context.Context, token string) (*models.User, error) {
    args := m.Called(ctx, token)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

// internal/api/handlers/auth_handler_test.go  
func (m *MockUserServiceForHandler) GetByVerificationToken(ctx context.Context, token string) (*models.User, error) {
    args := m.Called(ctx, token)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}
```

---

**Remember**: The goal is to catch bugs BEFORE they reach production, not document them after. Write comprehensive tests for every new feature.
