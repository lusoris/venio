# Project Guidelines & AI Assistant Instructions

This document outlines project standards, coding guidelines, and instructions for AI assistants working on Venio.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Development Standards](#development-standards)
3. [AI Assistant Instructions](#ai-assistant-instructions)
4. [Code Organization](#code-organization)
5. [Security Guidelines](#security-guidelines)

---

## Project Overview

### Project Name
**Venio** - Unified Media Management System

### Vision
A comprehensive orchestration layer for Movies, TV Shows, Music, and Adult content that unifies Overseerr, Lidarr, Whisparr, and media servers into a single, Netflix-like interface with intelligent content lifecycle management and community-driven features.

### Tech Stack
- **Backend:** Go 1.23+ with Gin web framework
- **Frontend:** Next.js 16+ with React and TypeScript
- **Database:** PostgreSQL 16 with pgx driver
- **Cache:** Redis 7
- **Authentication:** JWT (access + refresh tokens)
- **Deployment:** Docker & Docker Compose

### Core Repositories
- **Backend:** `cmd/`, `internal/`
- **Frontend:** `web/`
- **Infrastructure:** `deployments/`, `docker-compose.yml`
- **Database:** `migrations/`
- **Documentation:** `docs/`

---

## Development Standards

### Code Style

#### Go
- **Standard:** [Effective Go](https://golang.org/doc/effective_go)
- **Formatter:** `gofmt` (automated by editor)
- **Linter:** golangci-lint (enforced in CI)
- **Conventions:**
  - Package names: lowercase, single word
  - Exported functions/types: PascalCase
  - Private functions/types: camelCase
  - Interfaces: usually named with `-er` suffix (Reader, Writer, Handler)
  - Error handling: explicit, propagate up call stack
  - Comments: on exported types/functions (godoc format)

Example:
```go
package services

// UserService defines user-related business logic
type UserService interface {
    GetUser(ctx context.Context, id int64) (*models.User, error)
    CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
}

// NewUserService creates a new user service
func NewUserService(repo repositories.UserRepository) UserService {
    return &userService{repo: repo}
}
```

#### TypeScript/React
- **Standard:** [TypeScript Best Practices](https://www.typescriptlang.org/docs/handbook/)
- **Formatter:** Prettier (enforced via ESLint)
- **Linter:** ESLint (enforced in CI)
- **Conventions:**
  - Components: PascalCase in separate files
  - Hooks: camelCase with `use` prefix
  - Types/Interfaces: PascalCase
  - Constants: SCREAMING_SNAKE_CASE
  - Files: kebab-case for utilities, PascalCase for components

Example:
```typescript
// contexts/AuthContext.tsx
interface AuthContextType {
  user: User | null;
  login(credentials: LoginRequest): Promise<void>;
  logout(): void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  // ...
  return <AuthContext.Provider value={{user, login, logout}}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
}
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting)
- `refactor:` Code refactoring without feature changes
- `perf:` Performance improvements
- `test:` Test additions/changes
- `chore:` Build, dependencies, tooling
- `ci:` CI/CD pipeline changes
- `build:` Build script or tool changes

**Examples:**
```
feat(auth): add JWT refresh token support
fix(api): handle null user roles in list endpoint
docs(windows): add setup guide for Windows developers
refactor(database): consolidate connection pooling logic
```

### Testing

- **Go:** `testing` package + [testify](https://github.com/stretchr/testify)
- **TypeScript:** Jest + React Testing Library
- **Standards:**
  - Unit tests for business logic (services, utilities)
  - Integration tests for database operations
  - E2E tests for critical user flows
  - Minimum 70% code coverage for new code
  - Test files: `*_test.go` or `*.test.ts` adjacent to source

### Documentation

- **README.md:** Project overview, quick start, basic setup
- **docs/:** Detailed guides (architecture, API, configuration, deployment)
- **Code Comments:** Godoc-style for Go, JSDoc for TypeScript
- **Commits:** Detailed messages with context
- **PRs:** Clear description of changes, related issues, testing info

### File Organization

```
venio/
├── cmd/                      # Application entry points
│   ├── venio/               # Main server
│   └── worker/              # Background jobs
├── internal/                # Private application code
│   ├── api/                 # HTTP handlers, middleware, routes
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── middleware/      # HTTP middleware
│   │   └── routes.go        # Route configuration
│   ├── services/            # Business logic layer
│   ├── repositories/        # Data access layer
│   ├── models/              # Data structures
│   ├── config/              # Configuration management
│   └── database/            # Database connection
├── web/                     # Next.js frontend
│   ├── src/
│   │   ├── app/             # Next.js App Router pages
│   │   ├── components/      # Reusable React components
│   │   ├── contexts/        # React contexts (AuthContext, etc.)
│   │   ├── hooks/           # Custom React hooks
│   │   ├── lib/             # Utilities, API clients, helpers
│   │   └── styles/          # Global styles
├── migrations/              # Database migrations
├── docs/                    # Project documentation
├── scripts/                 # Development scripts
└── deployments/             # Deployment configurations
```

---

## AI Assistant Instructions

### General Principles

1. **Follow Existing Patterns:** Study existing code before making changes. Maintain consistency with current patterns and structures.

2. **Use Typed Structures:** Always use explicit types/interfaces. Avoid `any` type in TypeScript and bare interfaces in Go.

3. **Error Handling:**
   - Go: Always check and propagate errors with context
   - TypeScript: Use specific error types, provide meaningful error messages

4. **Security First:**
   - Validate all inputs
   - Use parameterized queries (prevent SQL injection)
   - Sanitize user data before logging
   - Follow principle of least privilege
   - Never log sensitive data (passwords, tokens, keys)

5. **Documentation:**
   - Add comments for non-obvious logic
   - Document exported functions/types
   - Update relevant docs when changing APIs

### Code Generation

#### When Creating Go Code
```go
// ALWAYS include:
// 1. Package declaration
// 2. Imports (sorted, organized)
// 3. Exported type documentation (godoc style)
// 4. Explicit error handling
// 5. Logging for debugging

package services

import (
    "context"
    "fmt"
    "log"

    "github.com/lusoris/venio/internal/models"
)

// UserService handles user-related operations
type UserService struct {
    repo repositories.UserRepository
    // ... other dependencies
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        log.Printf("failed to get user %d: %w", id, err)
        return nil, fmt.Errorf("get user: %w", err)
    }
    return user, nil
}
```

#### When Creating TypeScript Code
```typescript
// ALWAYS include:
// 1. Explicit type definitions
// 2. JSDoc comments for exports
// 3. Error handling with meaningful messages
// 4. Proper hook dependencies
// 5. Client-side validation

import { useState, useCallback } from 'react';

/** Request payload for user login */
export interface LoginRequest {
  email: string;
  password: string;
}

/** Response containing JWT tokens */
export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

/**
 * Hook for user authentication
 * @throws Error if login fails
 */
export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const login = useCallback(async (credentials: LoginRequest) => {
    setLoading(true);
    setError(null);
    try {
      // ... login logic
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return { user, login, error, loading };
}
```

### Testing Requirements

When creating new code, ensure:

1. **Unit Tests:** Business logic has test coverage
2. **Integration Tests:** Database operations are tested with real DB
3. **Error Cases:** Test both success and failure paths
4. **Edge Cases:** Empty strings, null values, boundary conditions

Example Go test:
```go
func TestUserService_GetUser(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&models.User{ID: 1}, nil)
    service := services.NewUserService(mockRepo)

    // Act
    user, err := service.GetUser(context.Background(), 1)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, int64(1), user.ID)
    mockRepo.AssertExpectations(t)
}
```

### Review Checklist

Before suggesting code changes, verify:

- [ ] Code follows project style guides
- [ ] Error handling is explicit and meaningful
- [ ] Sensitive data is not exposed
- [ ] Database queries use parameterized statements
- [ ] Types are explicit (no `any` in TS, no bare interfaces in Go)
- [ ] Comments explain "why", not "what"
- [ ] All imports are used
- [ ] Tests cover new functionality
- [ ] Commit message follows Conventional Commits
- [ ] Changes don't break existing functionality

### When You're Unsure

1. Check existing similar code in the codebase
2. Review relevant documentation in `docs/`
3. Look at git history for similar changes (`git log -p`)
4. Ask clarifying questions rather than making assumptions
5. Suggest multiple approaches with tradeoffs

---

## Code Organization

### Layered Architecture

Venio follows a three-layer architecture:

```
HTTP Layer (Handlers, Middleware)
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Database / External Services
```

**Benefits:**
- Clear separation of concerns
- Easy to test each layer independently
- Business logic not coupled to HTTP framework
- Easy to swap implementations (e.g., different database)

### Example: User Management

```
api/handlers/user_handler.go    # HTTP handlers
    ↓
services/user_service.go        # Business logic, validation
    ↓
repositories/user_repository.go # Database queries
    ↓
models/user.go                  # Data structures
```

---

## Security Guidelines

### Authentication & Authorization
- JWT tokens stored in HTTP-only cookies (frontend)
- Always validate tokens before accessing protected resources
- Implement token refresh with expiration
- Use bcrypt for password hashing (not MD5, SHA1)
- Implement rate limiting on auth endpoints

### Data Protection
- Validate all user inputs (length, type, format)
- Sanitize before logging (remove passwords, keys, tokens)
- Use parameterized queries to prevent SQL injection
- Encrypt sensitive data at rest (if applicable)
- Never commit secrets to git (use `.env` files, mark as private)

### API Security
- CORS: Configure properly for your frontend domain
- Rate limiting: Prevent brute force attacks
- Input validation: Check request size, format, types
- Output encoding: Prevent XSS attacks
- HTTPS: Required in production

### Code Review Checklist for Security
- [ ] No hardcoded secrets or credentials
- [ ] Input validation on all endpoints
- [ ] Proper error messages (don't leak system details)
- [ ] Authentication/Authorization checks present
- [ ] Sensitive data not logged
- [ ] SQL injection prevention (parameterized queries)
- [ ] CORS properly configured
- [ ] Rate limiting considered

---

## Frequently Referenced Patterns

### Error Handling (Go)
```go
// DO: Wrap errors with context
if err := someFunc(); err != nil {
    return fmt.Errorf("operation: %w", err)
}

// DON'T: Just propagate without context
return err
```

### API Responses (Go)
```go
// Consistent response structure
type Response struct {
    Data  interface{} `json:"data,omitempty"`
    Error string      `json:"error,omitempty"`
}

// Send success response
c.JSON(200, gin.H{"data": user})

// Send error response
c.JSON(400, gin.H{"error": "validation failed", "message": err.Error()})
```

### React Hooks (TypeScript)
```typescript
// Always include dependencies array
useEffect(() => {
  fetchData();
}, [dependency1, dependency2]); // <- required

// Extract API calls to custom hooks
const { data, loading, error } = useFetchUsers();
```

---

## Document Version History

| Date | Author | Changes |
|------|--------|---------|
| 2026-01-14 | AI Assistant | Initial creation with all project guidelines and standards |

---

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
- [TypeScript Best Practices](https://www.typescriptlang.org/docs/handbook/2/narrowing.html)
- [React Hooks Rules](https://react.dev/reference/rules/rules-of-hooks)
- [OWASP Security Guidelines](https://owasp.org/)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/sql-syntax.html)
