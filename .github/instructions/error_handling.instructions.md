---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/**/*.go"
description: Error Handling & Logging Guidelines
---

# Error Handling & Logging Guidelines

## Core Principle

**Security-First Error Handling:** Never expose internal system details to API clients. Log detailed errors server-side, return generic messages client-side.

## HTTP Error Responses

### ✅ DO: Separate Public and Internal Errors

```go
// ✅ GOOD: Generic message for client, detailed log for server
user, err := h.userService.Register(c.Request.Context(), &req)
if err != nil {
    // Server-side: Detailed logging
    log.Printf("Registration failed for %s: %v", req.Email, err)

    // Client-side: Generic response
    c.JSON(http.StatusBadRequest, gin.H{
        "error": "Registration failed",
        "message": "Unable to create account. Email may already be registered.",
    })
    return
}
```

### ❌ DON'T: Leak Internal Details

```go
// ❌ BAD: Exposes database structure, table names, file paths
c.JSON(http.StatusInternalServerError, gin.H{
    "error": err.Error(),  // ❌ May contain: "pq: duplicate key value violates unique constraint users_email_key"
})

// ❌ BAD: Validation errors expose field names
c.JSON(http.StatusBadRequest, gin.H{
    "message": err.Error(),  // ❌ May contain: "Key: 'User.Password' Error:Field validation for 'Password' failed on the 'min' tag"
})
```

## Error Code System

Use consistent error codes for frontend handling:

```go
type ErrorResponse struct {
    Error   string `json:"error"`   // Generic category
    Message string `json:"message"` // User-friendly message
    Code    string `json:"code"`    // Machine-readable code (optional)
}

// Error code constants
const (
    ErrCodeAuthInvalid      = "AUTH_INVALID_CREDENTIALS"
    ErrCodeAuthExpired      = "AUTH_TOKEN_EXPIRED"
    ErrCodeUserExists       = "USER_EMAIL_EXISTS"
    ErrCodeUserNotFound     = "USER_NOT_FOUND"
    ErrCodeValidationFailed = "VALIDATION_FAILED"
    ErrCodeServerError      = "SERVER_ERROR"
)

// Usage
c.JSON(http.StatusUnauthorized, ErrorResponse{
    Error:   "Authentication failed",
    Message: "Invalid email or password",
    Code:    ErrCodeAuthInvalid,
})
```

## Logging Levels

### ERROR: System failures requiring attention
```go
log.Printf("ERROR: Database connection failed: %v", err)
log.Printf("ERROR: Failed to generate JWT: %v", err)
```

### WARN: Unexpected but handled situations
```go
log.Printf("WARN: User roles not found for user %d, using empty roles", userID)
log.Printf("WARN: Rate limit exceeded for IP %s", ip)
```

### INFO: Normal operations
```go
log.Printf("INFO: User registered: %s (ID: %d)", user.Email, user.ID)
log.Printf("INFO: Successful login: %s", user.Email)
```

### DEBUG: Development-only details (never in production)
```go
if cfg.App.Debug {
    log.Printf("DEBUG: JWT token generated with claims: %+v", claims)
}
```

## Sensitive Data Protection

### ❌ NEVER Log:
- Passwords (plain or hashed)
- JWT tokens
- API keys
- Credit card numbers
- Personal identifiable information (PII) without explicit need

```go
// ❌ BAD
log.Printf("User login attempt: email=%s, password=%s", email, password)

// ✅ GOOD
log.Printf("User login attempt: email=%s", email)
```

### Sanitize Before Logging

```go
// ✅ DO: Redact sensitive fields
func sanitizeUserForLogging(user *models.User) string {
    return fmt.Sprintf("User{ID:%d, Email:%s, IsActive:%v}",
        user.ID, user.Email, user.IsActive)
}

log.Printf("User updated: %s", sanitizeUserForLogging(user))
```

## Common Patterns

### Database Errors
```go
user, err := s.repo.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        // Specific handling for not found
        return nil, fmt.Errorf("user not found")
    }
    // Generic database error
    log.Printf("Database error fetching user %d: %v", id, err)
    return nil, fmt.Errorf("database operation failed")
}
```

### Validation Errors
```go
if err := req.Validate(); err != nil {
    // Don't expose validation details
    c.JSON(http.StatusBadRequest, ErrorResponse{
        Error:   "Invalid input",
        Message: "Please check your input and try again",
        Code:    ErrCodeValidationFailed,
    })
    return
}
```

### Authentication Errors
```go
// ✅ DO: Same message for all auth failures (prevents user enumeration)
c.JSON(http.StatusUnauthorized, ErrorResponse{
    Error:   "Authentication failed",
    Message: "Invalid email or password",
    Code:    ErrCodeAuthInvalid,
})

// ❌ DON'T: Different messages reveal account existence
// "User not found" vs "Invalid password" ← Helps attackers
```

## Structured Logging (Future Enhancement)

When upgrading to structured logging (e.g., `log/slog`):

```go
logger.Error("User registration failed",
    "error", err,
    "email", req.Email,
    "username", req.Username,
    "timestamp", time.Now(),
)
```

## Error Wrapping

Use error wrapping for context:

```go
// ✅ DO: Wrap errors to preserve context
if err := s.repo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("failed to create user in repository: %w", err)
}

// ❌ DON'T: Lose context
if err := s.repo.Create(ctx, user); err != nil {
    return nil, err  // ❌ No context about where error occurred
}
```

## HTTP Status Code Mapping

| Status | Use Case | Example |
|--------|----------|---------|
| 400 Bad Request | Invalid input | "Email format invalid" |
| 401 Unauthorized | Authentication failed | "Invalid credentials" |
| 403 Forbidden | Insufficient permissions | "Admin role required" |
| 404 Not Found | Resource doesn't exist | "User not found" |
| 409 Conflict | Duplicate resource | "Email already registered" |
| 422 Unprocessable Entity | Valid format, business logic fail | "Username already taken" |
| 429 Too Many Requests | Rate limit exceeded | "Too many attempts" |
| 500 Internal Server Error | Unexpected server error | "Unable to process request" |

---

## Checklist for Every API Handler

- [ ] Validation errors return generic messages
- [ ] Database errors don't leak schema details
- [ ] Authentication failures use consistent messages
- [ ] Sensitive data is never logged
- [ ] Detailed errors logged server-side
- [ ] Appropriate HTTP status codes used
- [ ] Error codes defined for frontend
- [ ] Context preserved through error wrapping

---

**Remember:** User-friendly messages for clients, detailed logs for debugging. Never expose internal system structure or sensitive data.
