---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/api/**/*.go"
description: API Design Standards & REST Best Practices
---

# API Design Standards

## REST API Guidelines

### HTTP Methods

**Use correct HTTP verbs:**

```
GET    - Retrieve resource(s)       - Idempotent, Safe
POST   - Create resource            - Not Idempotent
PUT    - Update entire resource     - Idempotent
PATCH  - Update partial resource    - Not Idempotent
DELETE - Delete resource            - Idempotent
```

### URL Structure

**Follow RESTful conventions:**

```
✅ DO:
GET    /api/v1/users              - List users
GET    /api/v1/users/:id          - Get single user
POST   /api/v1/users              - Create user
PUT    /api/v1/users/:id          - Update user (full)
PATCH  /api/v1/users/:id          - Update user (partial)
DELETE /api/v1/users/:id          - Delete user

GET    /api/v1/users/:id/roles    - Get user roles
POST   /api/v1/users/:id/roles    - Assign role to user
DELETE /api/v1/users/:id/roles/:role_id  - Remove role from user

❌ DON'T:
GET    /api/v1/getUsers           - Verb in URL
POST   /api/v1/user/create        - Verb in URL
GET    /api/v1/users/list         - Verb in URL
PUT    /api/v1/userUpdate         - Inconsistent naming
```

### Response Format

**Standard response structure:**

```go
// Success response
{
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
}

// List response (paginated)
{
    "data": [
        {"id": 1, "email": "user1@example.com"},
        {"id": 2, "email": "user2@example.com"}
    ],
    "total": 100,
    "limit": 10,
    "offset": 0
}

// Error response
{
    "error": "Validation Error",
    "message": "Invalid email format",
    "details": {
        "field": "email",
        "value": "invalid-email"
    }
}
```

### HTTP Status Codes

**Use appropriate status codes:**

```
✅ Success:
200 OK               - Successful GET, PUT, PATCH, DELETE
201 Created          - Successful POST
204 No Content       - Successful DELETE (no response body)

✅ Client Errors:
400 Bad Request      - Invalid input, validation error
401 Unauthorized     - Missing or invalid authentication
403 Forbidden        - Authenticated but not authorized
404 Not Found        - Resource doesn't exist
409 Conflict         - Resource already exists
422 Unprocessable    - Valid syntax but semantic error
429 Too Many Requests - Rate limit exceeded

✅ Server Errors:
500 Internal Server Error - Unhandled error
503 Service Unavailable   - Temporary outage
```

**Example Handler:**
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // 400 - Invalid JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Bad Request",
            Message: "Invalid JSON format",
        })
        return
    }
    
    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        switch {
        // 409 - Duplicate email
        case errors.Is(err, services.ErrEmailAlreadyExists):
            c.JSON(http.StatusConflict, ErrorResponse{
                Error: "Conflict",
                Message: "Email already registered",
            })
        // 422 - Invalid input
        case errors.Is(err, services.ErrInvalidEmail):
            c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
                Error: "Validation Error",
                Message: err.Error(),
            })
        // 500 - Unknown error
        default:
            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Error: "Internal Server Error",
                Message: "Failed to create user",
            })
        }
        return
    }
    
    // 201 - Successfully created
    c.JSON(http.StatusCreated, user)
}
```

## Request Validation

### Input Validation

**Validate all input at handler level (basic) and service level (business rules):**

```go
// Handler - Basic validation (HTTP level)
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // Validate JSON structure
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Bad Request",
            Message: "Invalid JSON format",
        })
        return
    }
    
    // Validate required fields
    if req.Email == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Bad Request",
            Message: "Email is required",
        })
        return
    }
    
    // Delegate to service for business validation
    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    // ...
}

// Service - Business validation
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Validate email format
    if !isValidEmail(req.Email) {
        return nil, ErrInvalidEmail
    }
    
    // Validate password strength
    if len(req.Password) < 8 {
        return nil, errors.New("password must be at least 8 characters")
    }
    
    // Check for duplicate (business rule)
    existing, _ := s.repo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }
    
    // Continue with creation...
}
```

### Request DTOs

**Use dedicated request/response types:**

```go
// ✅ DO: Dedicated request/response types
type CreateUserRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

type UserResponse struct {
    ID        int64     `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Handle error
        return
    }
    
    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        // Handle error
        return
    }
    
    // Convert to response DTO (exclude sensitive fields like password)
    resp := UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
    
    c.JSON(http.StatusCreated, resp)
}

// ❌ DON'T: Reuse domain models directly
func (h *UserHandler) CreateUser(c *gin.Context) {
    var user models.User  // ❌ Exposes all fields including password
    if err := c.ShouldBindJSON(&user); err != nil {
        return
    }
    
    c.JSON(http.StatusCreated, user)  // ❌ Returns password hash
}
```

## Query Parameters

### Pagination

**Standard pagination parameters:**

```go
// ✅ DO: Consistent pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
    // Parse pagination params with defaults
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    // Enforce max limit
    if limit > 100 {
        limit = 100
    }
    if limit < 1 {
        limit = 10
    }
    
    users, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{...})
        return
    }
    
    total, _ := h.userService.CountUsers(c.Request.Context())
    
    c.JSON(http.StatusOK, gin.H{
        "data":   users,
        "total":  total,
        "limit":  limit,
        "offset": offset,
    })
}

// Example request:
// GET /api/v1/users?limit=20&offset=40
```

### Filtering

**Use query parameters for filtering:**

```go
// ✅ DO: Query parameter filtering
func (h *UserHandler) ListUsers(c *gin.Context) {
    filters := UserFilters{
        Email:    c.Query("email"),
        IsActive: c.Query("is_active") == "true",
        Role:     c.Query("role"),
    }
    
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    users, err := h.userService.ListUsers(c.Request.Context(), filters, limit, offset)
    // ...
}

// Example requests:
// GET /api/v1/users?is_active=true
// GET /api/v1/users?role=admin&limit=50
// GET /api/v1/users?email=john@example.com
```

### Sorting

**Use query parameters for sorting:**

```go
// ✅ DO: Query parameter sorting
func (h *UserHandler) ListUsers(c *gin.Context) {
    sortBy := c.DefaultQuery("sort_by", "created_at")
    sortOrder := c.DefaultQuery("sort_order", "desc")
    
    // Validate sort field
    validFields := map[string]bool{
        "id": true, "email": true, "created_at": true, "updated_at": true,
    }
    if !validFields[sortBy] {
        sortBy = "created_at"
    }
    
    // Validate sort order
    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = "desc"
    }
    
    users, err := h.userService.ListUsers(c.Request.Context(), sortBy, sortOrder, limit, offset)
    // ...
}

// Example requests:
// GET /api/v1/users?sort_by=email&sort_order=asc
// GET /api/v1/users?sort_by=created_at&sort_order=desc
```

## Authentication & Authorization

### Token-based Authentication

**JWT token in Authorization header:**

```go
// Middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error: "Unauthorized",
                Message: "Missing authorization header",
            })
            c.Abort()
            return
        }
        
        // Extract token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error: "Unauthorized",
                Message: "Invalid authorization format",
            })
            c.Abort()
            return
        }
        
        // Validate token
        claims, err := validateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error: "Unauthorized",
                Message: "Invalid or expired token",
            })
            c.Abort()
            return
        }
        
        // Store user info in context
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        
        c.Next()
    }
}

// Usage:
// Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Role-based Authorization

**Check permissions in middleware:**

```go
// RBAC Middleware
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt64("user_id")
        if userID == 0 {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error: "Unauthorized",
                Message: "Authentication required",
            })
            c.Abort()
            return
        }
        
        // Check permission
        hasPermission, err := checkUserPermission(c.Request.Context(), userID, permission)
        if err != nil || !hasPermission {
            c.JSON(http.StatusForbidden, ErrorResponse{
                Error: "Forbidden",
                Message: "Insufficient permissions",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Route with permission check
api.DELETE("/users/:id", 
    middleware.AuthMiddleware(),
    middleware.RequirePermission("users.delete"),
    userHandler.DeleteUser,
)
```

## Rate Limiting

### Rate Limit Headers

**Return rate limit info in headers:**

```go
// Rate limit middleware
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.ClientIP()
        
        allowed, remaining, resetTime := limiter.Allow(key)
        
        // Set rate limit headers
        c.Header("X-RateLimit-Limit", "100")
        c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
        
        if !allowed {
            c.JSON(http.StatusTooManyRequests, ErrorResponse{
                Error: "Too Many Requests",
                Message: "Rate limit exceeded. Try again later.",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Response headers:
// X-RateLimit-Limit: 100
// X-RateLimit-Remaining: 42
// X-RateLimit-Reset: 1705320000
```

## API Versioning

### URL Path Versioning

**Version in URL path:**

```go
// ✅ DO: Version in path
func SetupRouter() *gin.Engine {
    router := gin.Default()
    
    // API v1
    v1 := router.Group("/api/v1")
    {
        v1.GET("/users", userHandler.ListUsers)
        v1.POST("/users", userHandler.CreateUser)
    }
    
    // API v2 (future)
    v2 := router.Group("/api/v2")
    {
        v2.GET("/users", userHandlerV2.ListUsers)
        v2.POST("/users", userHandlerV2.CreateUser)
    }
    
    return router
}

// Example URLs:
// GET /api/v1/users
// GET /api/v2/users
```

## Error Responses

### Standard Error Format

**Consistent error response structure:**

```go
type ErrorResponse struct {
    Error   string      `json:"error"`            // Error type
    Message string      `json:"message"`          // Human-readable message
    Details interface{} `json:"details,omitempty"` // Optional details
}

// Usage examples:
// 400 Bad Request
{
    "error": "Validation Error",
    "message": "Invalid email format",
    "details": {
        "field": "email",
        "value": "invalid-email"
    }
}

// 401 Unauthorized
{
    "error": "Unauthorized",
    "message": "Invalid or expired token"
}

// 403 Forbidden
{
    "error": "Forbidden",
    "message": "Insufficient permissions to access this resource"
}

// 404 Not Found
{
    "error": "Not Found",
    "message": "User not found"
}

// 500 Internal Server Error
{
    "error": "Internal Server Error",
    "message": "An unexpected error occurred"
}
```

## CORS Configuration

### CORS Headers

**Configure CORS for frontend:**

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}
```

## Health Checks

### Health Endpoint

**Simple health check endpoint:**

```go
func (h *HealthHandler) Liveness(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
    })
}

func (h *HealthHandler) Readiness(c *gin.Context) {
    // Check dependencies
    dbHealthy := h.checkDatabase()
    redisHealthy := h.checkRedis()
    
    if !dbHealthy || !redisHealthy {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unavailable",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "ready",
    })
}

// Routes:
// GET /health/live   - Always returns 200 if server is running
// GET /health/ready  - Returns 200 if dependencies are healthy
```

## Best Practices Summary

### ✅ DO

- Use proper HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Use correct HTTP status codes (200, 201, 400, 401, 403, 404, 500)
- Version your API (/api/v1/)
- Use pagination for list endpoints
- Validate input at handler and service levels
- Use dedicated request/response DTOs
- Return consistent error responses
- Include rate limit headers
- Document API with Swagger/OpenAPI
- Use JWT tokens in Authorization header

### ❌ DON'T

- Use verbs in URLs (GET /getUsers)
- Return raw errors without mapping to HTTP status
- Expose internal error details in production
- Return sensitive fields (passwords, tokens) in responses
- Use magic numbers for status codes
- Mix authentication and business logic
- Return inconsistent response formats
- Ignore pagination limits (prevent DOS)
- Skip input validation

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-15  
**Maintained By:** Backend Team
