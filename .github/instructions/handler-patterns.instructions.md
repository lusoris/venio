---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/api/handlers/**/*.go"
description: Handler Layer Patterns & Best Practices
---

# Handler Layer Patterns & Best Practices

## Core Principle

**Handlers are the presentation layer**: They receive HTTP requests, validate inputs, call services, and format responses. They should NOT contain business logic.

## Consistent Error Response Pattern

### ✅ DO: Use Consistent ErrorResponse Struct

```go
// ALWAYS use ErrorResponse struct for errors
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
}

// Handler example
if err != nil {
    c.JSON(http.StatusBadRequest, ErrorResponse{
        Error:   "Invalid input",
        Message: "User ID must be a positive integer",
        Code:    "ERR_VALIDATION_FAILED",
    })
    return
}
```

### ❌ DON'T: Mix gin.H and ErrorResponse

```go
// ❌ WRONG: Inconsistent response formats
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})  // ❌ gin.H
c.JSON(http.StatusNotFound, ErrorResponse{...})                   // ❌ ErrorResponse

// ❌ WRONG: Inconsistent field names
c.JSON(http.StatusOK, gin.H{"value": roles})     // ❌ "value"
c.JSON(http.StatusOK, gin.H{"data": user})       // ❌ "data"
c.JSON(http.StatusOK, gin.H{"message": "ok"})    // ❌ "message"
```

### ✅ CORRECT: Consistent Success Responses

```go
// For single resource
type UserResponse struct {
    Data *models.User `json:"data"`
}

c.JSON(http.StatusOK, UserResponse{Data: user})

// For collections
type UsersResponse struct {
    Data []*models.User `json:"data"`
    Meta *PaginationMeta `json:"meta,omitempty"`
}

c.JSON(http.StatusOK, UsersResponse{Data: users, Meta: meta})

// For operations with no return value
type SuccessResponse struct {
    Message string `json:"message"`
}

c.JSON(http.StatusOK, SuccessResponse{Message: "Role assigned successfully"})
```

## Handler Method Structure

### Standard Flow

```go
func (h *Handler) CreateResource(c *gin.Context) {
    // 1. Parse and validate request
    var req CreateResourceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:   "Invalid request",
            Message: "Request body is malformed",
        })
        return
    }

    // 2. Extract context values (user ID, logger, etc.)
    ctx := c.Request.Context()
    
    // 3. Call service
    resource, err := h.service.Create(ctx, &req)
    if err != nil {
        h.handleError(c, err)  // Centralized error handling
        return
    }

    // 4. Return success response
    c.JSON(http.StatusCreated, ResourceResponse{Data: resource})
}
```

### Parameter Extraction

```go
// ✅ CORRECT: Parse and validate path parameters
func (h *Handler) GetByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:   "Invalid parameter",
            Message: "ID must be a positive integer",
        })
        return
    }

    // ... rest of handler
}

// ❌ WRONG: Don't skip validation
func (h *Handler) GetByID(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)  // ❌ Ignoring error
    // ... no validation if id <= 0
}
```

## Centralized Error Handling

### ✅ DO: Create Error Handler Helper

```go
// handlers/helpers.go
func (h *BaseHandler) handleError(c *gin.Context, err error) {
    // Map service errors to HTTP responses
    switch {
    case errors.Is(err, ErrNotFound):
        c.JSON(http.StatusNotFound, ErrorResponse{
            Error:   "Resource not found",
            Message: "The requested resource does not exist",
        })
    case errors.Is(err, ErrValidation):
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:   "Validation failed",
            Message: err.Error(),
        })
    case errors.Is(err, ErrUnauthorized):
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Error:   "Unauthorized",
            Message: "Authentication required",
        })
    case errors.Is(err, ErrForbidden):
        c.JSON(http.StatusForbidden, ErrorResponse{
            Error:   "Forbidden",
            Message: "Insufficient permissions",
        })
    default:
        // Never expose internal errors
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error:   "Internal server error",
            Message: "An unexpected error occurred",
        })
    }
}

// Usage in handler
if err != nil {
    h.handleError(c, err)
    return
}
```

### ❌ DON'T: Repeat Error Handling Logic

```go
// ❌ WRONG: Duplicated error handling in every handler
if err != nil {
    if errors.Is(err, ErrNotFound) {
        c.JSON(http.StatusNotFound, ...)
    } else {
        c.JSON(http.StatusInternalServerError, ...)
    }
    return
}
```

## Context Values Access

### ✅ DO: Use Type-Safe Context Helpers

```go
// context_keys.go
type contextKey string

const (
    UserIDKey  contextKey = "user_id"
    LoggerKey  contextKey = "logger"
    RequestIDKey contextKey = "request_id"
)

// helpers.go
func GetUserID(ctx context.Context) (int64, bool) {
    userID, ok := ctx.Value(UserIDKey).(int64)
    return userID, ok
}

func GetLogger(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
        return logger
    }
    return slog.Default()
}

// Handler usage
func (h *Handler) MyHandler(c *gin.Context) {
    ctx := c.Request.Context()
    
    userID, ok := GetUserID(ctx)
    if !ok {
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Error: "Unauthorized",
            Message: "User not authenticated",
        })
        return
    }
    
    logger := GetLogger(ctx)
    logger.Info("Processing request", "user_id", userID)
}
```

### ❌ DON'T: Use Unsafe Type Assertions

```go
// ❌ WRONG: Panic if type assertion fails
logger := c.Request.Context().Value("logger").(*slog.Logger)  // ❌ Panic if nil or wrong type

// ❌ WRONG: Ignore if value is missing
userID := c.Request.Context().Value("user_id")  // ❌ Could be nil
```

## Request/Response Models

### Request Structs

```go
// ✅ CORRECT: Validation tags on request
type CreateUserRequest struct {
    Email     string `json:"email" binding:"required,email,max=255"`
    Username  string `json:"username" binding:"required,min=3,max=50"`
    Password  string `json:"password" binding:"required,min=8,max=100"`
    FirstName string `json:"first_name" binding:"required,max=100"`
    LastName  string `json:"last_name" binding:"required,max=100"`
}

// ✅ CORRECT: Pointer fields for optional updates
type UpdateUserRequest struct {
    Email     *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
    Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
    FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=100"`
    LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=100"`
}
```

### Response Structs

```go
// ✅ CORRECT: Wrapper structs for consistent format
type UserResponse struct {
    Data *models.User `json:"data"`
}

type UsersListResponse struct {
    Data []*models.User `json:"data"`
    Meta *PaginationMeta `json:"meta,omitempty"`
}

type PaginationMeta struct {
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    TotalPages int   `json:"total_pages"`
    TotalCount int64 `json:"total_count"`
}

// Usage
c.JSON(http.StatusOK, UsersListResponse{
    Data: users,
    Meta: &PaginationMeta{
        Page:       1,
        PerPage:    20,
        TotalPages: 5,
        TotalCount: 100,
    },
})
```

## Handler Testing Patterns

### ✅ DO: Test HTTP Responses

```go
func TestUserHandler_GetByID_Success(t *testing.T) {
    // Setup
    mockService := new(mocks.MockUserService)
    handler := NewUserHandler(mockService)
    
    user := &models.User{ID: 1, Email: "test@example.com"}
    mockService.On("GetByID", mock.Anything, int64(1)).Return(user, nil)
    
    // Create request
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Request = httptest.NewRequest("GET", "/users/1", nil)
    
    // Execute
    handler.GetByID(c)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response UserResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, int64(1), response.Data.ID)
    
    mockService.AssertExpectations(t)
}
```

## Common Handler Mistakes

### ❌ DON'T Do This

```go
// 1. Don't put business logic in handlers
func (h *Handler) CreateUser(c *gin.Context) {
    // ❌ WRONG: Password hashing belongs in service
    hashedPassword, _ := bcrypt.GenerateFromPassword(...)
    
    // ❌ WRONG: Email validation belongs in service
    if !isValidEmail(req.Email) {
        return
    }
}

// 2. Don't expose internal errors
c.JSON(http.StatusInternalServerError, gin.H{
    "error": err.Error(),  // ❌ May leak DB schema, file paths, etc.
})

// 3. Don't skip parameter validation
id, _ := strconv.ParseInt(c.Param("id"), 10, 64)  // ❌ Ignoring error
user, _ := h.service.GetByID(ctx, id)              // ❌ Ignoring error

// 4. Don't use different response formats
c.JSON(http.StatusOK, gin.H{"value": data})   // ❌ Inconsistent
c.JSON(http.StatusOK, gin.H{"data": data})    // ❌ Inconsistent

// 5. Don't forget to return after errors
if err != nil {
    c.JSON(http.StatusBadRequest, ErrorResponse{...})
    // ❌ MISSING: return
}
// ❌ Continues execution
```

### ✅ DO This

```go
// 1. Keep handlers thin
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.handleError(c, ErrValidation)
        return
    }
    
    // ✅ Service handles business logic
    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, UserResponse{Data: user})
}

// 2. Always sanitize errors
func (h *Handler) handleError(c *gin.Context, err error) {
    // Map to safe HTTP errors
    // Never expose err.Error() directly to clients
}

// 3. Validate all parameters
id, err := strconv.ParseInt(c.Param("id"), 10, 64)
if err != nil || id <= 0 {
    c.JSON(http.StatusBadRequest, ErrorResponse{...})
    return
}

// 4. Use consistent response format
c.JSON(http.StatusOK, UserResponse{Data: user})         // ✅ Always "data"
c.JSON(http.StatusOK, UsersListResponse{Data: users})   // ✅ Always "data"

// 5. Always return after sending response
if err != nil {
    c.JSON(http.StatusBadRequest, ErrorResponse{...})
    return  // ✅ Stop execution
}
```

## Handler Organization

### File Structure

```
handlers/
├── auth_handler.go      # Authentication endpoints
├── user_handler.go      # User CRUD
├── role_handler.go      # Role management
├── permission_handler.go
├── admin_handler.go
├── helpers.go           # Shared helpers (handleError, etc.)
├── context.go           # Context helpers (GetUserID, etc.)
├── responses.go         # Response structs
└── requests.go          # Request structs
```

### Handler Struct Pattern

```go
// ✅ CORRECT: Group related handlers
type UserHandler struct {
    service services.UserService
    logger  *slog.Logger
}

func NewUserHandler(service services.UserService, logger *slog.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  logger,
    }
}

// Methods: GetByID, List, Create, Update, Delete
```

## Checklist for New Handlers

- [ ] Use consistent ErrorResponse struct for all errors
- [ ] Use consistent response wrappers (UserResponse, UsersListResponse)
- [ ] Validate all path parameters
- [ ] Validate all request bodies with binding tags
- [ ] Call service for business logic (don't implement in handler)
- [ ] Use centralized error handler (handleError)
- [ ] Always return after sending response
- [ ] Never expose internal errors to clients
- [ ] Add tests for success and error cases
- [ ] Use type-safe context helpers
- [ ] Log important events with structured logging
- [ ] Document endpoints with comments

## Response Format Standards

### Success Responses

```go
// Single resource: 200 OK / 201 Created
{
    "data": {
        "id": 1,
        "email": "user@example.com"
    }
}

// Collection: 200 OK
{
    "data": [
        {"id": 1, "email": "user1@example.com"},
        {"id": 2, "email": "user2@example.com"}
    ],
    "meta": {
        "page": 1,
        "per_page": 20,
        "total_pages": 5,
        "total_count": 100
    }
}

// Operation success: 200 OK
{
    "message": "Operation completed successfully"
}
```

### Error Responses

```go
// Client error: 400, 401, 403, 404, 409
{
    "error": "Validation failed",
    "message": "Email is required",
    "code": "ERR_VALIDATION_FAILED"
}

// Server error: 500
{
    "error": "Internal server error",
    "message": "An unexpected error occurred"
}
```

## CRITICAL: TODOs in Code

**NEVER leave TODOs in production code:**

```go
// ❌ WRONG: TODO in committed code
c.Request.Context().Value("logger") // TODO: Add proper logger

// ✅ CORRECT: Implement immediately or create GitHub issue
logger := GetLogger(c.Request.Context())
```

If a feature is incomplete:
1. Create GitHub issue with details
2. Reference issue in code comment: `// See issue #123`
3. Implement before merge to main

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial handler patterns guide |

**Remember**: Handlers are the HTTP presentation layer. Keep them thin, consistent, and focused on HTTP concerns. All business logic belongs in services.
