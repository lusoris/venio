---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/**/*.go"
description: Input Validation & Sanitization Standards
---

# Input Validation & Sanitization Guidelines

## Core Principle

**Defense in Depth:** Validate all user input at multiple layers. Never trust client-side validation alone.

## Validation Layers

1. **Struct Tags** (First line of defense)
2. **Custom Validators** (Business logic)
3. **Repository Layer** (Database constraints)

## Struct Tag Validation

### ✅ Required Fields with Constraints

```go
type CreateUserRequest struct {
    // Email: Required, valid format, max length
    Email     string  `json:"email" binding:"required,email,max=255"`

    // Username: Required, 3-50 chars, alphanumeric
    Username  string  `json:"username" binding:"required,min=3,max=50,alphanum"`

    // Names: Required, max length for database column
    FirstName string  `json:"first_name" binding:"required,max=100"`
    LastName  string  `json:"last_name" binding:"required,max=100"`

    // Password: Required, 8-128 chars (bcrypt limit)
    Password  string  `json:"password" binding:"required,min=8,max=128"`

    // Optional fields
    Avatar    *string `json:"avatar,omitempty" binding:"omitempty,url,max=500"`
    Bio       *string `json:"bio,omitempty" binding:"omitempty,max=1000"`
}
```

### Common Validation Tags

| Tag | Purpose | Example |
|-----|---------|---------|
| `required` | Field must be present | `binding:"required"` |
| `email` | Valid email format | `binding:"email"` |
| `min=n` | Minimum length/value | `binding:"min=8"` |
| `max=n` | Maximum length/value | `binding:"max=255"` |
| `len=n` | Exact length | `binding:"len=10"` |
| `alphanum` | Alphanumeric only | `binding:"alphanum"` |
| `url` | Valid URL format | `binding:"url"` |
| `numeric` | Numbers only | `binding:"numeric"` |
| `uuid` | Valid UUID format | `binding:"uuid"` |
| `oneof=a b c` | One of specified values | `binding:"oneof=admin user guest"` |

## Length Limits (DoS Prevention)

### ⚠️ Why Length Limits Matter

Unlimited input can cause:
- Database column overflow
- Memory exhaustion
- JSON parsing DoS
- Slow queries on large text fields

### ✅ Always Set Maximum Lengths

```go
// ✅ GOOD: All fields have max lengths
type BlogPostRequest struct {
    Title   string `json:"title" binding:"required,min=1,max=200"`
    Slug    string `json:"slug" binding:"required,min=1,max=100,alphanum"`
    Content string `json:"content" binding:"required,min=10,max=50000"`
    Tags    string `json:"tags" binding:"max=500"`
}

// ❌ BAD: No length limits = DoS risk
type BlogPostRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content" binding:"required"`  // ❌ Could be 100MB!
}
```

### Recommended Limits by Type

| Field Type | Max Length | Reasoning |
|------------|------------|-----------|
| Email | 255 | RFC 5321 limit |
| Username | 50 | Twitter/GitHub standard |
| Name (First/Last) | 100 | Covers most names worldwide |
| Password | 128 | Bcrypt has 72-byte limit, but allow more for hashing |
| URL | 500-2048 | Most URLs < 500, max 2048 for safety |
| Short Text | 1000 | Bio, description |
| Long Text | 50000 | Article content (adjust as needed) |
| Slug/ID | 100 | URL-safe identifiers |

## Custom Validators

### Register Custom Validators

```go
// validators/validators.go
package validators

import (
    "regexp"
    "github.com/go-playground/validator/v10"
)

// UsernameValidator validates username format
func UsernameValidator(fl validator.FieldLevel) bool {
    username := fl.Field().String()

    // Must start with letter, contain only alphanumeric, dash, underscore
    matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_-]*$`, username)
    return matched
}

// NoSQLInjection prevents basic SQL injection attempts
func NoSQLInjection(fl validator.FieldLevel) bool {
    input := fl.Field().String()

    // Block common SQL keywords in username/slug fields
    dangerousPatterns := []string{
        "DROP", "DELETE", "INSERT", "UPDATE", "SELECT",
        "--", "/*", "*/", ";",
    }

    for _, pattern := range dangerousPatterns {
        if strings.Contains(strings.ToUpper(input), pattern) {
            return false
        }
    }
    return true
}

// RegisterCustomValidators registers all custom validators
func RegisterCustomValidators(v *validator.Validate) {
    v.RegisterValidation("username", UsernameValidator)
    v.RegisterValidation("nosql", NoSQLInjection)
}
```

### Usage in Structs

```go
type CreateUserRequest struct {
    // Custom validator
    Username string `json:"username" binding:"required,min=3,max=50,username"`

    // Multiple custom validators
    Slug string `json:"slug" binding:"required,alphanum,nosql,max=100"`
}
```

## Sanitization

### HTML/XSS Prevention

```go
import "html"

// ✅ GOOD: Sanitize user-generated content before storing
func (s *Service) CreatePost(ctx context.Context, req *CreatePostRequest) error {
    // Escape HTML to prevent XSS
    req.Title = html.EscapeString(req.Title)
    req.Content = html.EscapeString(req.Content)

    return s.repo.Create(ctx, req)
}
```

### Whitespace Trimming

```go
// ✅ GOOD: Trim whitespace from text inputs
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    req.Email = strings.TrimSpace(req.Email)
    req.Username = strings.TrimSpace(req.Username)
    req.FirstName = strings.TrimSpace(req.FirstName)
    req.LastName = strings.TrimSpace(req.LastName)

    // Proceed with validation
    if err := req.Validate(); err != nil {
        return err
    }
    // ...
}
```

### Email Normalization

```go
// ✅ GOOD: Normalize emails for consistency
func normalizeEmail(email string) string {
    email = strings.TrimSpace(email)
    email = strings.ToLower(email)
    return email
}

func (s *Service) Register(ctx context.Context, req *CreateUserRequest) error {
    req.Email = normalizeEmail(req.Email)
    // ...
}
```

## Enum Validation

### Using `oneof` Tag

```go
type UpdateUserRequest struct {
    Status string `json:"status" binding:"oneof=active inactive suspended"`
    Role   string `json:"role" binding:"oneof=admin moderator user guest"`
}
```

### Using Constants

```go
const (
    StatusActive    = "active"
    StatusInactive  = "inactive"
    StatusSuspended = "suspended"
)

func (r *UpdateUserRequest) Validate() error {
    validStatuses := map[string]bool{
        StatusActive:    true,
        StatusInactive:  true,
        StatusSuspended: true,
    }

    if !validStatuses[r.Status] {
        return errors.New("invalid status")
    }
    return nil
}
```

## Array/Slice Validation

### Limit Array Size

```go
type BulkCreateRequest struct {
    // Limit to prevent DoS
    Users []CreateUserRequest `json:"users" binding:"required,min=1,max=100,dive"`
}
```

The `dive` tag validates each element in the array.

## SQL Injection Prevention

### ✅ ALWAYS Use Parameterized Queries

```go
// ✅ GOOD: Parameterized query (pgx does this automatically)
row := r.pool.QueryRow(ctx,
    `SELECT id, email FROM users WHERE email = $1`,
    email,  // ✅ Safe: Parameter binding
)

// ❌ BAD: String concatenation
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)  // ❌ SQL Injection!
row := r.pool.QueryRow(ctx, query)
```

### Dynamic Query Building (Advanced)

If you MUST build dynamic queries (rare cases):

```go
// ✅ GOOD: Whitelist allowed fields
func buildOrderByClause(sortField string) (string, error) {
    allowedFields := map[string]bool{
        "created_at": true,
        "username":   true,
        "email":      true,
    }

    if !allowedFields[sortField] {
        return "", errors.New("invalid sort field")
    }

    // Safe to use in query now
    return fmt.Sprintf("ORDER BY %s", sortField), nil
}

// Usage
orderBy, err := buildOrderByClause(req.SortBy)
if err != nil {
    return err
}

query := fmt.Sprintf(`
    SELECT id, email, username
    FROM users
    %s
    LIMIT $1 OFFSET $2
`, orderBy)

rows, err := r.pool.Query(ctx, query, limit, offset)
```

## File Upload Validation

```go
type UploadRequest struct {
    // Validate file type
    ContentType string `json:"content_type" binding:"required,oneof=image/jpeg image/png image/gif"`

    // Validate file size (in bytes)
    Size int64 `json:"size" binding:"required,min=1,max=5242880"` // Max 5MB

    // Validate filename
    Filename string `json:"filename" binding:"required,max=255"`
}

func validateFile(file *multipart.FileHeader) error {
    // Check file size
    if file.Size > 5*1024*1024 {  // 5MB
        return errors.New("file too large")
    }

    // Check file extension
    allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
    ext := strings.ToLower(filepath.Ext(file.Filename))

    allowed := false
    for _, allowedExt := range allowedExtensions {
        if ext == allowedExt {
            allowed = true
            break
        }
    }

    if !allowed {
        return errors.New("invalid file type")
    }

    return nil
}
```

## URL Validation

```go
import "net/url"

// ✅ GOOD: Validate and sanitize URLs
func validateURL(rawURL string) error {
    u, err := url.Parse(rawURL)
    if err != nil {
        return errors.New("invalid URL format")
    }

    // Only allow HTTPS (for security)
    if u.Scheme != "https" {
        return errors.New("only HTTPS URLs allowed")
    }

    // Block private/internal IPs (SSRF prevention)
    host := u.Hostname()
    if isPrivateIP(host) {
        return errors.New("private IP addresses not allowed")
    }

    return nil
}

func isPrivateIP(host string) bool {
    // Check for localhost
    if host == "localhost" || host == "127.0.0.1" {
        return true
    }

    // Check for private IP ranges
    privateRanges := []string{
        "10.", "172.16.", "172.17.", "172.18.", "172.19.",
        "172.20.", "172.21.", "172.22.", "172.23.", "172.24.",
        "172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
        "172.30.", "172.31.", "192.168.",
    }

    for _, prefix := range privateRanges {
        if strings.HasPrefix(host, prefix) {
            return true
        }
    }

    return false
}
```

## Checklist for Every Input Struct

- [ ] All string fields have `max` length
- [ ] Email fields use `email` validator
- [ ] URLs use `url` validator
- [ ] Enums use `oneof` or custom validator
- [ ] Passwords have `min=8,max=128`
- [ ] Optional fields use `omitempty`
- [ ] Arrays have `max` size with `dive`
- [ ] Custom business logic has validator function
- [ ] No SQL injection risk (parameterized queries)
- [ ] XSS prevention (HTML escape if needed)

## Common Validation Patterns

### Registration Form
```go
type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email,max=255"`
    Username  string `json:"username" binding:"required,min=3,max=50,alphanum"`
    Password  string `json:"password" binding:"required,min=8,max=128"`
    FirstName string `json:"first_name" binding:"required,max=100"`
    LastName  string `json:"last_name" binding:"required,max=100"`
}
```

### Search Query
```go
type SearchRequest struct {
    Query  string `json:"query" binding:"required,min=1,max=200"`
    Limit  int    `json:"limit" binding:"min=1,max=100"`
    Offset int    `json:"offset" binding:"min=0"`
}
```

### Update Request (Partial)
```go
type UpdateUserRequest struct {
    Email     *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
    Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
    FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=100"`
}
```

---

**Remember:** Validate early, validate often. Defense in depth prevents attacks at multiple layers.
