---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "**/*.go,**/*.sql,**/*.yml,**/*.yaml"
description: Security Best Practices
---

# Security Best Practices

## Core Principle

**Security by Design, Defense in Depth**: Never trust user input. Always validate, sanitize, and encrypt. Assume breach and minimize impact.

## JWT Security

### Secret Management

```go
// ✅ CORRECT: Load from environment
type JWTConfig struct {
    Secret            string        // From env: JWT_SECRET
    ExpirationTime    time.Duration // From env: JWT_EXPIRATION
    RefreshExpiryDays int          // From env: JWT_REFRESH_EXPIRY_DAYS
}

func Load() (*Config, error) {
    viper.SetEnvPrefix("VENIO")
    viper.AutomaticEnv()
    
    jwtSecret := viper.GetString("JWT_SECRET")
    if jwtSecret == "" {
        return nil, errors.New("JWT_SECRET environment variable is required")
    }
    
    if len(jwtSecret) < 32 {
        return nil, errors.New("JWT_SECRET must be at least 32 characters")
    }
    
    // ... load other config
}

// ❌ WRONG: Hardcoded secrets
const jwtSecret = "my-secret-key"  // ❌ NEVER!
```

### Token Generation Best Practices

```go
// ✅ CORRECT: Secure token generation
func (s *DefaultAuthService) generateAccessToken(user *models.User, roles []string) (string, error) {
    claims := models.TokenClaims{
        UserID:   user.ID,
        Email:    user.Email,
        Username: user.Username,
        Roles:    roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.ExpirationTime)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "venio-api",
            Subject:   fmt.Sprintf("%d", user.ID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.JWT.Secret))
}

// ❌ WRONG: No expiration
claims := jwt.MapClaims{
    "user_id": user.ID,
    // ❌ Missing: ExpiresAt, IssuedAt, NotBefore
}

// ❌ WRONG: Weak signing method
token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)  // ❌ No signature!
```

### Token Validation

```go
// ✅ CORRECT: Strict validation
func (s *DefaultAuthService) ValidateToken(tokenString string) (*models.TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.config.JWT.Secret), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    claims, ok := token.Claims.(*models.TokenClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }

    return claims, nil
}

// ❌ WRONG: No signature verification
token, _ := jwt.Parse(tokenString, nil)  // ❌ No key verification!
```

## Password Security

### Hashing with bcrypt

```go
// ✅ CORRECT: Use bcrypt with appropriate cost
const bcryptCost = 12  // Configurable via environment

func hashPassword(password string) (string, error) {
    // Validate password strength first
    if len(password) < 8 {
        return "", errors.New("password must be at least 8 characters")
    }
    
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    
    return string(hash), nil
}

func comparePassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ❌ WRONG: Store plain text passwords
user.Password = password  // ❌ NEVER store plain text!

// ❌ WRONG: Use weak hashing
hash := md5.Sum([]byte(password))  // ❌ MD5 is broken!
hash := sha1.Sum([]byte(password))  // ❌ SHA1 is weak!
```

### Password Requirements

```go
// ✅ CORRECT: Enforce password policy
type PasswordPolicy struct {
    MinLength          int
    RequireUppercase   bool
    RequireLowercase   bool
    RequireDigit       bool
    RequireSpecialChar bool
}

func ValidatePassword(password string, policy PasswordPolicy) error {
    if len(password) < policy.MinLength {
        return fmt.Errorf("password must be at least %d characters", policy.MinLength)
    }

    if policy.RequireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return errors.New("password must contain at least one uppercase letter")
    }

    if policy.RequireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
        return errors.New("password must contain at least one lowercase letter")
    }

    if policy.RequireDigit && !regexp.MustCompile(`[0-9]`).MatchString(password) {
        return errors.New("password must contain at least one digit")
    }

    if policy.RequireSpecialChar && !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};:'",.<>?]`).MatchString(password) {
        return errors.New("password must contain at least one special character")
    }

    return nil
}
```

## SQL Injection Prevention

### ✅ CORRECT: Parameterized Queries

```go
// ✅ Always use parameterized queries
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `
        SELECT id, email, username, first_name, last_name, password, is_active, created_at, updated_at
        FROM users
        WHERE email = $1
    `
    
    var user models.User
    err := r.db.QueryRow(ctx, query, email).Scan(
        &user.ID,
        &user.Email,
        // ... other fields
    )
    
    return &user, err
}

// ✅ Batch operations with parameters
func (r *PostgresUserRepository) GetByIDs(ctx context.Context, ids []int64) ([]*models.User, error) {
    query := `
        SELECT id, email, username, first_name, last_name, password, is_active, created_at, updated_at
        FROM users
        WHERE id = ANY($1)
    `
    
    rows, err := r.db.Query(ctx, query, ids)
    // ... process rows
}
```

### ❌ WRONG: String Concatenation

```go
// ❌ NEVER concatenate user input into SQL
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := "SELECT * FROM users WHERE email = '" + email + "'"  // ❌ SQL INJECTION!
    
    // Attacker input: "admin@example.com' OR '1'='1"
    // Result: "SELECT * FROM users WHERE email = 'admin@example.com' OR '1'='1'"
    // Returns all users!
}

// ❌ Even with fmt.Sprintf
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)  // ❌ Still vulnerable!
```

## CORS Security

### ✅ CORRECT: Whitelist Origins

```go
// config/cors.go
type CORSConfig struct {
    AllowedOrigins []string  // Load from env: CORS_ALLOWED_ORIGINS
}

func LoadCORSConfig() (*CORSConfig, error) {
    originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
    if originsStr == "" {
        return nil, errors.New("CORS_ALLOWED_ORIGINS is required")
    }
    
    origins := strings.Split(originsStr, ",")
    
    // Validate each origin
    for _, origin := range origins {
        if !isValidOrigin(origin) {
            return nil, fmt.Errorf("invalid origin: %s", origin)
        }
    }
    
    return &CORSConfig{AllowedOrigins: origins}, nil
}

// middleware/cors.go
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Only allow whitelisted origins
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
        }

        c.Next()
    }
}

// ❌ WRONG: Allow all origins
c.Header("Access-Control-Allow-Origin", "*")  // ❌ Allows any origin!
c.Header("Access-Control-Allow-Credentials", "true")  // ❌ Security risk!

// ❌ WRONG: Reflect origin without validation
origin := c.Request.Header.Get("Origin")
c.Header("Access-Control-Allow-Origin", origin)  // ❌ No validation!
```

## Rate Limiting

### Authentication Endpoints

```go
// ✅ CORRECT: Strict rate limiting on auth endpoints
authConfig := &ratelimit.Config{
    Requests: 5,                // 5 requests
    Window:   time.Minute,      // per minute
}

authLimiter := ratelimit.NewLimiter(
    ratelimit.LimiterTypeRedis,
    authConfig,
    redisClient,
)

router.POST("/auth/login", 
    RateLimitMiddleware(authLimiter, IPBasedKey),  // IP-based for login
    handler.Login,
)

router.POST("/auth/register", 
    RateLimitMiddleware(authLimiter, IPBasedKey),  // IP-based for register
    handler.Register,
)

// ❌ WRONG: No rate limiting on auth
router.POST("/auth/login", handler.Login)  // ❌ Brute force possible!
```

### API Endpoints

```go
// ✅ CORRECT: Reasonable limits for API
apiConfig := &ratelimit.Config{
    Requests: 100,              // 100 requests
    Window:   time.Minute,      // per minute
}

apiLimiter := ratelimit.NewLimiter(
    ratelimit.LimiterTypeRedis,
    apiConfig,
    redisClient,
)

api := router.Group("/api/v1")
api.Use(RateLimitMiddleware(apiLimiter, UserBasedKey))  // User-based for API
```

## Session Management

### Token Expiration

```go
// ✅ CORRECT: Short-lived access tokens, longer refresh tokens
type JWTConfig struct {
    AccessTokenExpiry  time.Duration // 15 minutes to 1 hour
    RefreshTokenExpiry time.Duration // 7 to 30 days
}

// Load from environment with sensible defaults
cfg := &JWTConfig{
    AccessTokenExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 24*time.Hour),
    RefreshTokenExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
}

// ❌ WRONG: Long-lived access tokens
AccessTokenExpiry: 365 * 24 * time.Hour  // ❌ 1 year is too long!

// ❌ WRONG: No expiration
// If token never expires, can't revoke access!
```

### Token Refresh Strategy

```go
// ✅ CORRECT: Validate refresh token thoroughly
func (s *DefaultAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
    // Validate refresh token
    claims, err := s.ValidateToken(refreshToken)
    if err != nil {
        return "", errors.New("invalid refresh token")
    }

    // Check token type (if you embed it in claims)
    if claims.Type != "refresh" {
        return "", errors.New("not a refresh token")
    }

    // Get user and verify still active
    user, err := s.userService.GetByID(ctx, claims.UserID)
    if err != nil {
        return "", errors.New("user not found")
    }

    if !user.IsActive {
        return "", errors.New("user account is inactive")
    }

    // Get fresh roles
    roles, _ := s.userRoleService.GetUserRoles(ctx, user.ID)

    // Generate new access token
    return s.generateAccessToken(user, roles)
}
```

## Environment Variables & Secrets

### ✅ CORRECT: Secret Management

```go
// Load secrets from environment only
func LoadConfig() (*Config, error) {
    // Required secrets
    requiredSecrets := []string{
        "JWT_SECRET",
        "DATABASE_PASSWORD",
        "REDIS_PASSWORD",
    }

    for _, secret := range requiredSecrets {
        value := os.Getenv(secret)
        if value == "" {
            return nil, fmt.Errorf("required secret %s is not set", secret)
        }
    }

    // Validate secret strength
    jwtSecret := os.Getenv("JWT_SECRET")
    if len(jwtSecret) < 32 {
        return nil, errors.New("JWT_SECRET must be at least 32 characters")
    }

    return &Config{
        JWT: JWTConfig{
            Secret: jwtSecret,
        },
        Database: DatabaseConfig{
            Password: os.Getenv("DATABASE_PASSWORD"),
        },
    }, nil
}

// ❌ WRONG: Hardcoded secrets
const (
    jwtSecret      = "my-secret-key"           // ❌ Exposed in code!
    dbPassword     = "postgres123"             // ❌ In version control!
    redisPassword  = "redis-password"          // ❌ Never!
)

// ❌ WRONG: Secrets in config files
// config.yml
jwt:
  secret: "my-secret-key"  # ❌ Committed to git!
```

### Environment File (.env)

```bash
# ✅ CORRECT: .env.example (template, commit this)
JWT_SECRET=generate-random-32-char-string
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=venio
DATABASE_PASSWORD=change-me-in-production
DATABASE_NAME=venio

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=change-me-in-production

CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.venio.io

# .env (actual secrets, NEVER commit)
JWT_SECRET=7Jk9mN2pQ5rS8tV1wX4yZ6aB3cD0eF7gH
DATABASE_PASSWORD=K8mP2qT5vY9zA3cF6iL0oR4sU7xB1dE
REDIS_PASSWORD=nQ9sV2yB5eH8kN1pT4wZ7aC0fI3lO6rU

# ❌ WRONG: Real secrets in .env.example
JWT_SECRET=7Jk9mN2pQ5rS8tV1wX4yZ6aB3cD0eF7gH  # ❌ Real secret in template!
```

### .gitignore

```gitignore
# ✅ ALWAYS ignore secret files
.env
.env.local
.env.production
*.key
*.pem
secrets/
```

## Security Headers

### ✅ CORRECT: Comprehensive Headers

```go
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Prevent clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // Prevent MIME type sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // Enable browser XSS protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy (strict)
        c.Header("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self'; "+
            "style-src 'self' 'unsafe-inline'; "+  // unsafe-inline only if needed
            "img-src 'self' data: https:; "+
            "font-src 'self'; "+
            "connect-src 'self'; "+
            "frame-ancestors 'none'")
        
        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions Policy
        c.Header("Permissions-Policy", 
            "geolocation=(), microphone=(), camera=()")

        c.Next()
    }
}

// Production-only: HSTS
func HSTSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Strict-Transport-Security", 
            "max-age=31536000; includeSubDomains; preload")
        c.Next()
    }
}
```

## Input Validation & Sanitization

### ✅ CORRECT: Validate Everything

```go
// Use binding tags
type CreateUserRequest struct {
    Email     string `json:"email" binding:"required,email,max=255"`
    Username  string `json:"username" binding:"required,min=3,max=50,alphanum"`
    Password  string `json:"password" binding:"required,min=8,max=100"`
    FirstName string `json:"first_name" binding:"required,max=100"`
    LastName  string `json:"last_name" binding:"required,max=100"`
}

// Additional business logic validation
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    // Validate password strength
    if err := ValidatePasswordStrength(req.Password); err != nil {
        return nil, err
    }

    // Check for SQL injection patterns (defense in depth)
    if containsSQLKeywords(req.Email) || containsSQLKeywords(req.Username) {
        return nil, errors.New("invalid characters in input")
    }

    // ... rest of create logic
}

// ❌ WRONG: Trust user input
func (s *UserService) Create(ctx context.Context, email, username, password string) error {
    // No validation!
    query := fmt.Sprintf("INSERT INTO users (email, username, password) VALUES ('%s', '%s', '%s')", 
        email, username, password)  // ❌ SQL injection + no validation!
}
```

## Error Handling (Security)

### ✅ CORRECT: Never Leak Internal Details

```go
// Public errors (safe for clients)
var (
    ErrInvalidCredentials = errors.New("invalid email or password")
    ErrAccountInactive    = errors.New("account is inactive")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrForbidden          = errors.New("insufficient permissions")
)

// Handler
func (h *AuthHandler) Login(c *gin.Context) {
    // ... validate request

    token, err := h.authService.Login(ctx, req.Email, req.Password)
    if err != nil {
        // Log detailed error server-side
        h.logger.Error("Login failed", 
            "email", req.Email, 
            "error", err,
            "ip", c.ClientIP(),
        )

        // Return generic message to client
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Error:   "Authentication failed",
            Message: "Invalid email or password",  // ✅ Generic message
            Code:    "ERR_AUTH_INVALID",
        })
        return
    }

    c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// ❌ WRONG: Expose internal errors
c.JSON(http.StatusInternalServerError, gin.H{
    "error": err.Error(),  // ❌ May leak DB schema, paths, etc.
})

// ❌ WRONG: Different messages help attackers
if userNotFound {
    c.JSON(401, gin.H{"error": "User not found"})  // ❌ Confirms user doesn't exist
} else if invalidPassword {
    c.JSON(401, gin.H{"error": "Invalid password"})  // ❌ Confirms user exists
}
```

## Logging Security

### ✅ CORRECT: Never Log Sensitive Data

```go
// ✅ CORRECT: Log without sensitive data
logger.Info("User login attempt",
    "email", email,  // OK to log email
    "ip", c.ClientIP(),
    "user_agent", c.Request.UserAgent(),
)

logger.Error("Login failed",
    "email", email,
    "error", "invalid credentials",  // Generic error
    "ip", c.ClientIP(),
)

// ❌ WRONG: Log passwords or tokens
logger.Info("Login", 
    "email", email,
    "password", password,  // ❌ NEVER log passwords!
)

logger.Info("Token generated",
    "token", token,  // ❌ NEVER log tokens!
)
```

## Security Checklist

- [ ] **Secrets**: All secrets in environment variables, never hardcoded
- [ ] **JWT**: Secret ≥32 chars, short expiry (≤24h access, ≤30d refresh)
- [ ] **Passwords**: bcrypt with cost ≥12, enforce strength policy
- [ ] **SQL**: Always use parameterized queries, never concatenate
- [ ] **CORS**: Whitelist origins, never use "*" with credentials
- [ ] **Rate Limiting**: 5/min on auth, 100/min on API
- [ ] **Headers**: X-Frame-Options, CSP, HSTS (prod), X-Content-Type-Options
- [ ] **Validation**: Validate all inputs, use binding tags
- [ ] **Errors**: Generic messages to clients, detailed logs server-side
- [ ] **Logging**: Never log passwords, tokens, or sensitive data
- [ ] **HTTPS**: Always in production, no exceptions
- [ ] **.gitignore**: Exclude .env, *.key, *.pem, secrets/

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial security best practices guide |

**Remember**: Security is not optional. Every line of code should assume malicious input. Defense in depth means multiple layers of protection.
