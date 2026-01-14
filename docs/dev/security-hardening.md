# Security Hardening Guide

Comprehensive security hardening guidelines based on OWASP Top 10, industry best practices, and Venio-specific configurations.

## Table of Contents

1. [OWASP Top 10 2023 Coverage](#owasp-top-10-2023-coverage)
2. [Application Security](#application-security)
3. [Database Security](#database-security)
4. [API Security](#api-security)
5. [Container Security](#container-security)
6. [Secrets Management](#secrets-management)
7. [Network Security](#network-security)
8. [Monitoring & Auditing](#monitoring--auditing)

---

## OWASP Top 10 2023 Coverage

### A01:2023 – Broken Access Control

**Risk:** Unauthorized users accessing protected resources.

**Mitigations:**

```go
// ✅ DO: Implement middleware-based authorization
func (m *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt64("user_id")

        // Check permission in database
        hasPermission, err := m.rbacService.UserHasPermission(c.Request.Context(), userID, permission)
        if err != nil || !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// ✅ DO: Validate resource ownership
func (h *UserHandler) UpdateUser(c *gin.Context) {
    requestUserID := c.GetInt64("user_id")
    targetUserID := c.Param("id")

    // Users can only update themselves unless they're admin
    if !h.isAdmin(requestUserID) && requestUserID != targetUserID {
        c.JSON(http.StatusForbidden, gin.H{"error": "cannot update other users"})
        return
    }

    // Continue with update
}
```

**Checklist:**
- [x] RBAC implemented with roles and permissions
- [x] Middleware validates permissions on every request
- [x] Resource ownership verified before modifications
- [x] Default deny policy (explicit grants only)
- [x] Regular permission audits

---

### A02:2023 – Cryptographic Failures

**Risk:** Sensitive data exposed due to weak encryption.

**Mitigations:**

```go
// ✅ DO: Use bcrypt for password hashing (cost 12+)
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12) // Cost factor 12
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

// ✅ DO: Use TLS 1.3 minimum
// In production deployment (nginx/traefik):
// ssl_protocols TLSv1.3;
// ssl_ciphers HIGH:!aNULL:!MD5;

// ✅ DO: Use HTTPS everywhere
server := &http.Server{
    Addr:         ":443",
    TLSConfig:    tlsConfig,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
}
server.ListenAndServeTLS("cert.pem", "key.pem")
```

**Environment Variables Protection:**

```bash
# ✅ DO: Use strong, random secrets
JWT_SECRET=$(openssl rand -base64 48)  # 64+ characters
API_KEY=$(openssl rand -base64 32)     # 44+ characters

# ❌ DON'T: Use predictable secrets
JWT_SECRET="supersecret"  # ❌ Too weak!
```

**Checklist:**
- [x] Passwords hashed with bcrypt (cost 12)
- [x] JWT secrets minimum 32 bytes (44 base64 chars)
- [x] TLS 1.3 enforced in production
- [x] No sensitive data logged
- [x] Database credentials encrypted at rest

---

### A03:2023 – Injection

**Risk:** SQL/Command injection attacks.

**Mitigations:**

```go
// ✅ DO: Always use parameterized queries
const query = `
    SELECT id, email, username
    FROM users
    WHERE email = $1 AND deleted_at IS NULL
`
err := pool.QueryRow(ctx, query, userEmail).Scan(&user.ID, &user.Email, &user.Username)

// ❌ DON'T: Concatenate user input into queries
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", userEmail) // ❌ SQL Injection!
```

**Input Validation:**

```go
// ✅ DO: Validate and sanitize all inputs
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
    Password string `json:"password" validate:"required,min=8,max=128"`
}

func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Additional validation
    validate := validator.New()
    if err := validate.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Safe to use req now
}
```

**Checklist:**
- [x] All queries use parameterized statements
- [x] Input validation on all API endpoints
- [x] Length limits enforced
- [x] Special characters sanitized
- [x] No shell commands with user input

---

### A04:2023 – Insecure Design

**Risk:** Architectural flaws in security controls.

**Mitigations:**

**Rate Limiting:**

```go
// ✅ DO: Implement rate limiting on sensitive endpoints
import "github.com/gin-contrib/limiter"

func setupRateLimiting(router *gin.Engine) {
    // Global rate limit: 100 req/min per IP
    store := memory.NewStore()
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
    }
    router.Use(limiter.New(store, rate))

    // Stricter limits for auth endpoints
    authGroup := router.Group("/api/v1/auth")
    authRate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  5, // Only 5 login attempts per minute
    }
    authGroup.Use(limiter.New(store, authRate))
}
```

**Account Lockout:**

```go
// ✅ DO: Implement progressive delays after failed login attempts
type LoginAttemptTracker struct {
    attempts map[string]int
    mu       sync.RWMutex
}

func (t *LoginAttemptTracker) CheckAndDelay(email string) error {
    t.mu.Lock()
    defer t.mu.Unlock()

    attempts := t.attempts[email]
    if attempts >= 5 {
        return errors.New("account temporarily locked, try again later")
    }

    // Progressive delay: 2^attempts seconds
    if attempts > 0 {
        delay := time.Duration(1<<attempts) * time.Second
        time.Sleep(delay)
    }

    return nil
}

func (t *LoginAttemptTracker) RecordFailure(email string) {
    t.mu.Lock()
    t.attempts[email]++
    t.mu.Unlock()

    // Clear after 15 minutes
    time.AfterFunc(15*time.Minute, func() {
        t.mu.Lock()
        delete(t.attempts, email)
        t.mu.Unlock()
    })
}
```

**Checklist:**
- [x] Rate limiting on all API endpoints
- [x] Account lockout after failed logins
- [x] CAPTCHA on registration/login forms
- [x] Secure password reset flow
- [x] Multi-factor authentication (MFA) support

---

### A05:2023 – Security Misconfiguration

**Risk:** Insecure default configurations.

**Mitigations:**

**HTTP Security Headers:**

```go
// ✅ DO: Set comprehensive security headers
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Prevent clickjacking
        c.Header("X-Frame-Options", "DENY")

        // Prevent MIME sniffing
        c.Header("X-Content-Type-Options", "nosniff")

        // Enable XSS protection
        c.Header("X-XSS-Protection", "1; mode=block")

        // HSTS (force HTTPS for 2 years)
        c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

        // Content Security Policy
        c.Header("Content-Security-Policy",
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data: https:; "+
            "font-src 'self' data:; "+
            "connect-src 'self'; "+
            "frame-ancestors 'none'")

        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

        // Permissions Policy
        c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

        c.Next()
    }
}
```

**CORS Configuration:**

```go
// ✅ DO: Configure CORS properly
import "github.com/gin-contrib/cors"

func setupCORS(router *gin.Engine) {
    config := cors.Config{
        AllowOrigins:     []string{os.Getenv("FRONTEND_URL")}, // Specific origins only
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    router.Use(cors.New(config))
}

// ❌ DON'T: Allow all origins
AllowOrigins: []string{"*"}  // ❌ Security risk!
```

**Checklist:**
- [x] All security headers configured
- [x] CORS restricted to known origins
- [x] Error messages don't leak system details
- [x] Debug mode disabled in production
- [x] Unnecessary services disabled

---

### A06:2023 – Vulnerable and Outdated Components

**Risk:** Known vulnerabilities in dependencies.

**Mitigations:**

```bash
# ✅ DO: Run Snyk scans before each commit
snyk test

# ✅ DO: Audit Go dependencies
go list -m all | nancy sleuth

# ✅ DO: Keep dependencies updated
go get -u ./...
go mod tidy
```

**Automated Dependency Checks:**

```yaml
# .github/workflows/security.yml
name: Security Scan
on: [push, pull_request]
jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Snyk
        uses: snyk/actions/golang@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
```

**Checklist:**
- [x] Snyk integrated into CI/CD
- [x] Weekly dependency updates
- [x] Vulnerability alerts enabled
- [x] No dependencies with known CVEs
- [x] Using latest stable versions

---

### A07:2023 – Identification and Authentication Failures

**Risk:** Weak authentication mechanisms.

**Mitigations:**

**Password Requirements:**

```go
// ✅ DO: Enforce strong password policy
func ValidatePassword(password string) error {
    if len(password) < 12 {
        return errors.New("password must be at least 12 characters")
    }

    var (
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )

    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        return errors.New("password must contain uppercase, lowercase, number, and special character")
    }

    return nil
}
```

**Session Management:**

```go
// ✅ DO: Implement secure session management
type SessionManager struct {
    store    *redis.Client
    duration time.Duration
}

func (sm *SessionManager) CreateSession(ctx context.Context, userID int64) (string, error) {
    sessionID := generateSecureToken() // Cryptographically secure random

    sessionData := map[string]interface{}{
        "user_id":    userID,
        "created_at": time.Now().Unix(),
        "expires_at": time.Now().Add(sm.duration).Unix(),
    }

    data, _ := json.Marshal(sessionData)
    err := sm.store.Set(ctx, sessionID, data, sm.duration).Err()
    return sessionID, err
}

func (sm *SessionManager) InvalidateSession(ctx context.Context, sessionID string) error {
    return sm.store.Del(ctx, sessionID).Err()
}
```

**Checklist:**
- [x] Minimum password length 12 characters
- [x] Password complexity enforced
- [x] Account lockout after failed attempts
- [x] Session timeout configured (24 hours)
- [x] Secure session token generation

---

### A08:2023 – Software and Data Integrity Failures

**Risk:** Unsigned code or data tampering.

**Mitigations:**

```go
// ✅ DO: Sign critical payloads
import "crypto/hmac"
import "crypto/sha256"

func SignPayload(payload []byte, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write(payload)
    return hex.EncodeToString(h.Sum(nil))
}

func VerifyPayload(payload []byte, signature string, secret string) bool {
    expectedSignature := SignPayload(payload, secret)
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

**Checklist:**
- [x] Docker images signed
- [x] Dependencies verified with checksums
- [x] CI/CD pipeline secured
- [x] Code signing for releases
- [x] Audit logs for deployments

---

### A09:2023 – Security Logging and Monitoring Failures

**Risk:** Attacks go undetected.

**Mitigations:**

```go
// ✅ DO: Log security-relevant events
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)

    user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        // Log failed login attempt
        h.logger.Warn("failed login attempt",
            zap.String("email", req.Email),
            zap.String("ip", c.ClientIP()),
            zap.Time("timestamp", time.Now()),
        )
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // Log successful login
    h.logger.Info("successful login",
        zap.Int64("user_id", user.ID),
        zap.String("ip", c.ClientIP()),
    )

    c.JSON(http.StatusOK, gin.H{"data": user})
}
```

**Checklist:**
- [x] Failed login attempts logged
- [x] Permission changes logged
- [x] Data access logged (audit trail)
- [x] Alerts configured for anomalies
- [x] Logs retained for 90+ days

---

### A10:2023 – Server-Side Request Forgery (SSRF)

**Risk:** Server makes unauthorized requests.

**Mitigations:**

```go
// ✅ DO: Validate URLs before fetching
import "net/url"

func isSafeURL(rawURL string) error {
    u, err := url.Parse(rawURL)
    if err != nil {
        return err
    }

    // Block private IP ranges
    if u.Hostname() == "localhost" ||
       strings.HasPrefix(u.Hostname(), "127.") ||
       strings.HasPrefix(u.Hostname(), "192.168.") ||
       strings.HasPrefix(u.Hostname(), "10.") {
        return errors.New("private IP addresses not allowed")
    }

    // Only allow HTTPS
    if u.Scheme != "https" {
        return errors.New("only HTTPS allowed")
    }

    return nil
}
```

**Checklist:**
- [x] URL validation before external requests
- [x] Whitelist of allowed domains
- [x] Timeout on HTTP clients (5 seconds)
- [x] No redirects followed
- [x] Response size limits enforced

---

## Application Security

### Secure Coding Guidelines

```go
// ✅ DO: Validate file uploads
func (h *FileHandler) UploadFile(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
        return
    }
    defer file.Close()

    // Validate file size (max 10MB)
    if header.Size > 10*1024*1024 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
        return
    }

    // Validate file type by content (not just extension)
    buffer := make([]byte, 512)
    _, err = file.Read(buffer)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
        return
    }

    contentType := http.DetectContentType(buffer)
    if !isAllowedContentType(contentType) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
        return
    }

    // Continue with safe file processing
}
```

---

## Database Security

### PostgreSQL Hardening

```sql
-- ✅ DO: Use Row-Level Security (RLS)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

CREATE POLICY user_access_policy ON users
    FOR ALL
    USING (id = current_setting('app.current_user_id')::bigint);

-- ✅ DO: Restrict user permissions
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM venio_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON users TO venio_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO venio_app;

-- ✅ DO: Use SCRAM-SHA-256 authentication
-- In postgresql.conf:
password_encryption = 'scram-sha-256'
```

**Connection Security:**

```bash
# ✅ DO: Enable SSL for database connections
# postgresql.conf
ssl = on
ssl_cert_file = '/path/to/server.crt'
ssl_key_file = '/path/to/server.key'
ssl_ca_file = '/path/to/ca.crt'

# In DATABASE_URL
DATABASE_URL="postgres://user:pass@host/db?sslmode=require"
```

---

## Container Security

### Dockerfile Hardening

```dockerfile
# ✅ DO: Use minimal base images
FROM golang:1.25-alpine AS builder

# ✅ DO: Run as non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# ✅ DO: Multi-stage build to minimize attack surface
FROM alpine:3.21
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

COPY --from=builder --chown=appuser:appuser /app/venio /app/venio

USER appuser

# ✅ DO: Set read-only filesystem
VOLUME ["/tmp"]
RUN mkdir -p /home/appuser/.cache && chown appuser:appuser /home/appuser/.cache
```

**docker-compose.yml Security:**

```yaml
# ✅ DO: Limit container capabilities
services:
  venio:
    image: venio:latest
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
```

---

## Secrets Management

### Environment Variables

```bash
# ✅ DO: Use secrets management tools
# Docker Swarm secrets
docker secret create jwt_secret /run/secrets/jwt_secret

# Kubernetes secrets
kubectl create secret generic venio-secrets \
  --from-literal=jwt-secret="$(openssl rand -base64 48)" \
  --from-literal=database-password="$(openssl rand -base64 32)"
```

**Never Commit Secrets:**

```gitignore
# .gitignore
.env
.env.local
.env.*.local
*.key
*.pem
secrets/
```

---

## Document Revision

| Version | Date | Author | Changes | Source Version |
|---------|------|--------|---------|----------------|
| 1.0.0 | 2026-01-14 | AI Assistant | Initial security hardening guide | - |

## Referenced Documentation

- **OWASP Top 10 2023:** [OWASP Top Ten](https://owasp.org/www-project-top-ten/) (Released: 2023-09-24)
- **OWASP ASVS:** [Application Security Verification Standard](https://owasp.org/www-project-application-security-verification-standard/) (v4.0.3)
- **NIST Cybersecurity Framework:** [NIST CSF](https://www.nist.gov/cyberframework)
- **CIS Docker Benchmark:** [Docker Security Benchmarks](https://www.cisecurity.org/benchmark/docker) (v1.6.0)
- **PostgreSQL Security:** [PostgreSQL Security Best Practices](https://www.postgresql.org/docs/18/security.html)
- **Go Security:** [Go Security Policy](https://go.dev/security/)
- **JWT Best Practices:** [RFC 8725 - JWT Best Current Practices](https://datatracker.ietf.org/doc/html/rfc8725)
