---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/api/middleware/**/*.go"
description: Middleware Patterns & Best Practices
---

# Middleware Patterns & Best Practices

## Core Principle

**Middleware is the cross-cutting concern layer**: Authentication, authorization, logging, metrics, rate limiting, CORS. They should be composable, reusable, and order-aware.

## Consistent Error Response Pattern

### ✅ DO: Use ErrorResponse Struct

```go
// middleware/responses.go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
}

// In middleware
if authHeader == "" {
    c.JSON(http.StatusUnauthorized, ErrorResponse{
        Error:   "Unauthorized",
        Message: "Missing authorization header",
        Code:    "ERR_AUTH_MISSING",
    })
    c.Abort()
    return
}
```

### ❌ DON'T: Use gin.H in Middleware

```go
// ❌ WRONG: Inconsistent with handler responses
c.JSON(http.StatusUnauthorized, gin.H{
    "error":   "Unauthorized",
    "message": "Missing authorization header",
})
```

## Context Management

### ✅ DO: Use Type-Safe Context Helpers

```go
// middleware/context.go
type contextKey string

const (
    UserIDKey    contextKey = "user_id"
    EmailKey     contextKey = "email"
    UsernameKey  contextKey = "username"
    RolesKey     contextKey = "roles"
    RequestIDKey contextKey = "request_id"
)

// Type-safe setters
func SetUserID(c *gin.Context, userID int64) {
    c.Set(string(UserIDKey), userID)
}

func SetRoles(c *gin.Context, roles []string) {
    c.Set(string(RolesKey), roles)
}

// Type-safe getters
func GetUserID(c *gin.Context) (int64, bool) {
    value, exists := c.Get(string(UserIDKey))
    if !exists {
        return 0, false
    }
    userID, ok := value.(int64)
    return userID, ok
}

func GetRoles(c *gin.Context) ([]string, bool) {
    value, exists := c.Get(string(RolesKey))
    if !exists {
        return nil, false
    }
    roles, ok := value.([]string)
    return roles, ok
}
```

### ❌ DON'T: Use String Keys Directly

```go
// ❌ WRONG: Type-unsafe, error-prone
c.Set("user_id", claims.UserID)
c.Set("roles", claims.Roles)

// ❌ WRONG: No type safety, can panic
userID := c.GetInt64("user_id")  // Panic if wrong type
roles := c.Get("roles").([]string)  // Panic if nil or wrong type
```

## Authentication Middleware Pattern

### Complete Example

```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "github.com/lusoris/venio/internal/services"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract token from header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error:   "Unauthorized",
                Message: "Missing authorization header",
                Code:    "ERR_AUTH_MISSING",
            })
            c.Abort()
            return
        }

        // 2. Validate Bearer format
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error:   "Unauthorized",
                Message: "Invalid authorization header format",
                Code:    "ERR_AUTH_FORMAT",
            })
            c.Abort()
            return
        }

        token := parts[1]

        // 3. Validate token
        claims, err := authService.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error:   "Unauthorized",
                Message: "Invalid or expired token",
                Code:    "ERR_AUTH_INVALID",
            })
            c.Abort()
            return
        }

        // 4. Store claims in context using type-safe helpers
        SetUserID(c, claims.UserID)
        SetEmail(c, claims.Email)
        SetUsername(c, claims.Username)
        SetRoles(c, claims.Roles)

        c.Next()
    }
}
```

## Authorization (RBAC) Middleware Pattern

```go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user roles from context
        userRoles, exists := GetRoles(c)
        if !exists {
            c.JSON(http.StatusUnauthorized, ErrorResponse{
                Error:   "Unauthorized",
                Message: "Authentication required",
                Code:    "ERR_AUTH_REQUIRED",
            })
            c.Abort()
            return
        }

        // Check if user has any of the required roles
        hasRole := false
        for _, required := range roles {
            for _, userRole := range userRoles {
                if userRole == required {
                    hasRole = true
                    break
                }
            }
            if hasRole {
                break
            }
        }

        if !hasRole {
            c.JSON(http.StatusForbidden, ErrorResponse{
                Error:   "Forbidden",
                Message: "Insufficient permissions",
                Code:    "ERR_PERMISSION_DENIED",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// Usage
router.GET("/admin/users", 
    AuthMiddleware(authService),
    RequireRole("admin", "super_admin"),
    handler.ListUsers,
)
```

## Rate Limiting Middleware Pattern

### Using Abstraction

```go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/lusoris/venio/internal/ratelimit"
)

func RateLimitMiddleware(limiter ratelimit.Limiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := keyFunc(c)
        
        allowed, err := limiter.Allow(c.Request.Context(), key)
        if err != nil {
            // Log error but don't fail request
            c.Next()
            return
        }

        if !allowed {
            c.JSON(http.StatusTooManyRequests, ErrorResponse{
                Error:   "Rate limit exceeded",
                Message: "Too many requests. Please try again later.",
                Code:    "ERR_RATE_LIMIT",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// Key functions
func IPBasedKey(c *gin.Context) string {
    return "ratelimit:ip:" + c.ClientIP()
}

func UserBasedKey(c *gin.Context) string {
    userID, exists := GetUserID(c)
    if !exists {
        return "ratelimit:ip:" + c.ClientIP()
    }
    return fmt.Sprintf("ratelimit:user:%d", userID)
}

// Usage
authRateLimiter := ratelimit.NewLimiter(ratelimit.LimiterTypeRedis, 
    &ratelimit.Config{Requests: 5, Window: time.Minute}, 
    redisClient,
)

router.POST("/auth/login", 
    RateLimitMiddleware(authRateLimiter, IPBasedKey),
    handler.Login,
)
```

## CORS Middleware Pattern

```go
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Check if origin is allowed
        allowed := false
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                allowed = true
                break
            }
        }

        if allowed {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Credentials", "true")
            c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
            c.Header("Access-Control-Max-Age", "86400")
        }

        // Handle preflight
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

// Usage with configuration
allowedOrigins := []string{
    "http://localhost:3000",  // Development
    "https://app.venio.io",   // Production
}

router.Use(CORSMiddleware(allowedOrigins))
```

## Security Headers Middleware

```go
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Prevent clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // Prevent MIME type sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // Enable browser XSS protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy
        c.Header("Content-Security-Policy", 
            "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'")
        
        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions Policy
        c.Header("Permissions-Policy", 
            "geolocation=(), microphone=(), camera=()")

        c.Next()
    }
}

// Strict version for production
func StrictSecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // HSTS - Force HTTPS (only in production!)
        c.Header("Strict-Transport-Security", 
            "max-age=31536000; includeSubDomains; preload")
        
        // All other headers from SecurityHeadersMiddleware
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Content-Security-Policy", 
            "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' https:; font-src 'self'; connect-src 'self'")
        
        c.Next()
    }
}
```

## Logging Middleware Pattern

```go
func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        // Process request
        c.Next()

        // Calculate latency
        latency := time.Since(start)
        status := c.Writer.Status()

        // Get user context if available
        userID, _ := GetUserID(c)

        // Log request
        logger.Info("HTTP request",
            "method", method,
            "path", path,
            "status", status,
            "latency_ms", latency.Milliseconds(),
            "ip", c.ClientIP(),
            "user_id", userID,
        )
    }
}
```

## Metrics Middleware Pattern

```go
func MetricsMiddleware(collector metrics.Collector) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        // Process request
        c.Next()

        // Record metrics
        duration := time.Since(start)
        status := fmt.Sprintf("%d", c.Writer.Status())
        
        collector.RecordHTTPRequest(method, path, status, duration)
    }
}
```

## Middleware Order (CRITICAL)

### ✅ CORRECT Order

```go
func SetupRouter(
    authService services.AuthService,
    userRoleService services.UserRoleService,
    logger *slog.Logger,
    metrics metrics.Collector,
    rateLimiter ratelimit.Limiter,
) *gin.Engine {
    router := gin.New()

    // 1. Recovery - Must be first to catch panics
    router.Use(gin.Recovery())

    // 2. Logging - Log all requests
    router.Use(LoggingMiddleware(logger))

    // 3. Metrics - Record all requests
    router.Use(MetricsMiddleware(metrics))

    // 4. CORS - Early to handle preflight
    router.Use(CORSMiddleware(allowedOrigins))

    // 5. Security Headers - Apply to all responses
    router.Use(SecurityHeadersMiddleware())

    // 6. Rate Limiting (Global) - Optional global limiter
    router.Use(RateLimitMiddleware(globalLimiter, IPBasedKey))

    // 7. Authentication - Route-specific, not global
    // Applied per route group

    // 8. Authorization - Route-specific, after auth
    // Applied per route or handler

    return router
}

// Route-specific middleware
api := router.Group("/api/v1")
{
    // Public routes
    auth := api.Group("/auth")
    {
        auth.POST("/login", 
            RateLimitMiddleware(authLimiter, IPBasedKey),  // Stricter rate limit
            handler.Login,
        )
        auth.POST("/register", handler.Register)
    }

    // Protected routes
    users := api.Group("/users")
    users.Use(AuthMiddleware(authService))  // Auth required
    {
        users.GET("", handler.ListUsers)  // All authenticated users
        users.GET("/:id", handler.GetUser)

        // Admin-only routes
        users.DELETE("/:id", 
            RequireRole("admin"),  // Authorization after auth
            handler.DeleteUser,
        )
    }
}
```

### ❌ WRONG Order

```go
// ❌ Bad: Auth before recovery - panics won't be caught
router.Use(AuthMiddleware(authService))
router.Use(gin.Recovery())

// ❌ Bad: Metrics after auth - won't record failed auth attempts
router.Use(AuthMiddleware(authService))
router.Use(MetricsMiddleware(metrics))

// ❌ Bad: CORS after handlers - preflight will fail
router.Use(handler)
router.Use(CORSMiddleware(allowedOrigins))
```

## Middleware Testing Patterns

### Testing Authentication Middleware

```go
func TestAuthMiddleware_Success(t *testing.T) {
    // Setup
    mockAuthService := new(mocks.MockAuthService)
    middleware := AuthMiddleware(mockAuthService)

    claims := &models.TokenClaims{
        UserID:   1,
        Email:    "test@example.com",
        Username: "testuser",
        Roles:    []string{"user"},
    }
    mockAuthService.On("ValidateToken", "valid-token").Return(claims, nil)

    // Create test context
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/", nil)
    c.Request.Header.Set("Authorization", "Bearer valid-token")

    // Execute
    middleware(c)

    // Assert
    assert.False(t, c.IsAborted())
    
    userID, exists := GetUserID(c)
    assert.True(t, exists)
    assert.Equal(t, int64(1), userID)
    
    mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
    mockAuthService := new(mocks.MockAuthService)
    middleware := AuthMiddleware(mockAuthService)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/", nil)
    // No Authorization header

    middleware(c)

    assert.True(t, c.IsAborted())
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
```

## Common Middleware Mistakes

### ❌ DON'T Do This

```go
// 1. Don't forget c.Abort() after error response
if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
    // ❌ MISSING: c.Abort() - request will continue!
}

// 2. Don't use gin.H for errors
c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})  // ❌ Inconsistent

// 3. Don't use string keys for context
c.Set("user_id", userID)  // ❌ Type-unsafe
userID := c.Get("user_id").(int64)  // ❌ Can panic

// 4. Don't apply auth globally
router.Use(AuthMiddleware(authService))  // ❌ Breaks /health, /metrics, /auth/login

// 5. Don't forget to call c.Next()
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Do stuff
        // ❌ MISSING: c.Next() - request ends here!
    }
}

// 6. Don't modify request after c.Next()
func BadMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        c.Request.Header.Set("X-Custom", "value")  // ❌ Too late!
    }
}
```

### ✅ DO This

```go
// 1. Always abort after error response
if err != nil {
    c.JSON(http.StatusUnauthorized, ErrorResponse{...})
    c.Abort()  // ✅ Stop processing
    return
}

// 2. Use consistent error response
c.JSON(http.StatusBadRequest, ErrorResponse{...})  // ✅ Consistent

// 3. Use type-safe context helpers
SetUserID(c, userID)  // ✅ Type-safe
userID, exists := GetUserID(c)  // ✅ Can't panic

// 4. Apply auth per route or group
users := api.Group("/users")
users.Use(AuthMiddleware(authService))  // ✅ Protected routes only

// 5. Always call c.Next()
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Pre-processing
        c.Next()  // ✅ Continue chain
        // Post-processing (logging, metrics, etc.)
    }
}

// 6. Modify before c.Next()
func GoodMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Request.Header.Set("X-Custom", "value")  // ✅ Before c.Next()
        c.Next()
    }
}
```

## Middleware Configuration Pattern

```go
// middleware/config.go
type Config struct {
    RateLimit struct {
        Auth struct {
            Requests int
            Window   time.Duration
        }
        API struct {
            Requests int
            Window   time.Duration
        }
    }
    CORS struct {
        AllowedOrigins []string
    }
    Security struct {
        EnableHSTS bool
    }
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    // Load from environment or config file
    cfg.RateLimit.Auth.Requests = getEnvInt("RATE_LIMIT_AUTH_REQUESTS", 5)
    cfg.RateLimit.Auth.Window = getEnvDuration("RATE_LIMIT_AUTH_WINDOW", time.Minute)
    
    cfg.CORS.AllowedOrigins = getEnvStringSlice("CORS_ALLOWED_ORIGINS", 
        []string{"http://localhost:3000"})
    
    cfg.Security.EnableHSTS = getEnvBool("SECURITY_ENABLE_HSTS", false)
    
    return cfg, nil
}
```

## Checklist for New Middleware

- [ ] Use `ErrorResponse` struct for all errors
- [ ] Call `c.Abort()` after error responses
- [ ] Call `c.Next()` to continue chain
- [ ] Use type-safe context helpers
- [ ] Handle missing context values gracefully
- [ ] Consider middleware order
- [ ] Add tests for success and error cases
- [ ] Document what the middleware does
- [ ] Validate configuration at startup
- [ ] Log important events (auth failures, rate limits, etc.)
- [ ] Record metrics where appropriate
- [ ] Never expose internal errors to clients

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial middleware patterns guide |

**Remember**: Middleware order matters! Recovery → Logging → Metrics → CORS → Security Headers → Rate Limiting → Auth → Authorization.
