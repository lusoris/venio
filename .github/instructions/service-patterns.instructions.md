---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/services/**/*.go"
description: Service Layer Business Logic Patterns
---

# Service Layer Patterns

## Core Principle

**Services contain ALL business logic.** They validate input, orchestrate repository calls, enforce business rules, and transform data between layers.

## Service Method Patterns

### 1. Standard Method Call Flow

Every service method follows this pattern:

```go
func (s *service) MethodName(ctx context.Context, params...) (*Model, error) {
    // 1. Validate input parameters
    if invalid {
        return nil, errors.New("validation error: ...")
    }

    // 2. Check business rules (existence, duplicates, permissions)
    existing, err := s.repository.CheckExistence(ctx, ...)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("already exists")
    }

    // 3. Call repository operation(s)
    result, err := s.repository.Operation(ctx, ...)
    if err != nil {
        return nil, err
    }

    // 4. Check result validity
    if result == nil {
        return nil, errors.New("not found")
    }

    // 5. Return result
    return result, nil
}
```

### 2. Create Operations Pattern

```go
// ✅ CORRECT: Complete Create pattern
func (s *userService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    // 1. Validate request (often done by request struct tags, but double-check)
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("validation error: %w", err)
    }

    // 2. Check for duplicates
    exists, err := s.userRepository.Exists(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("check email existence: %w", err)
    }
    if exists {
        return nil, errors.New("email already registered")
    }

    // 3. Hash password (business logic)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("hash password: %w", err)
    }
    req.Password = string(hashedPassword)

    // 4. Create in repository
    user, err := s.userRepository.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}
```

**Important**: Repository returns the created entity, so no additional `GetByID` call needed.

### 3. Read Operations Pattern

#### Get Single Entity

```go
func (s *roleService) GetByID(ctx context.Context, id int64) (*models.Role, error) {
    // 1. Validate ID
    if id <= 0 {
        return nil, errors.New("invalid role ID")
    }

    // 2. Get from repository
    role, err := s.roleRepository.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. Check if found
    if role == nil {
        return nil, errors.New("role not found")
    }

    return role, nil
}
```

**Key Point**: Repository returns `(nil, nil)` for not found, service converts to error.

#### List Operations

```go
func (s *permissionService) List(ctx context.Context, limit, offset int) ([]*models.Permission, int64, error) {
    // 1. Validate and set defaults
    if limit <= 0 {
        limit = 10
    }
    if limit > 100 {
        limit = 100
    }
    if offset < 0 {
        offset = 0
    }

    // 2. Get from repository (returns []models.Permission)
    permissions, total, err := s.permissionRepository.List(ctx, limit, offset)
    if err != nil {
        return nil, 0, err
    }

    // 3. Convert []Model to []*Model (service layer returns pointers)
    permPointers := make([]*models.Permission, len(permissions))
    for i := range permissions {
        permPointers[i] = &permissions[i]
    }

    return permPointers, total, nil
}
```

**Key Point**: Repository returns `[]models.Model` (values), service converts to `[]*models.Model` (pointers) for handlers.

### 4. Update Operations Pattern

```go
func (s *roleService) Update(ctx context.Context, id int64, req models.UpdateRoleRequest) (*models.Role, error) {
    // 1. Validate ID
    if id <= 0 {
        return nil, errors.New("invalid role ID")
    }

    // 2. Check if entity exists
    existing, err := s.roleRepository.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        return nil, errors.New("role not found")
    }

    // 3. Check business rules (e.g., duplicate name)
    if req.Name != nil && *req.Name != existing.Name {
        duplicate, dupErr := s.roleRepository.GetByName(ctx, *req.Name)
        if dupErr != nil {
            return nil, dupErr
        }
        if duplicate != nil {
            return nil, errors.New("role with this name already exists")
        }
    }

    // 4. Update in repository
    role, err := s.roleRepository.Update(ctx, id, &req)
    if err != nil {
        return nil, err
    }

    return role, nil
}
```

**Key Point**: Repository returns updated entity, avoiding extra round-trip.

### 5. Delete Operations Pattern

```go
func (s *userService) Delete(ctx context.Context, id int64) error {
    // 1. Validate ID
    if id <= 0 {
        return errors.New("invalid user ID")
    }

    // 2. Check if entity exists (business rule: can't delete non-existent)
    user, err := s.userRepository.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if user == nil {
        return errors.New("user not found")
    }

    // 3. Business logic: Check if safe to delete
    // Example: Can't delete users with active sessions, orders, etc.
    // (Add checks as needed)

    // 4. Delete from repository
    return s.userRepository.Delete(ctx, id)
}
```

**Key Point**: Always check existence first - deletion should be explicit, not silent.

## Service Interface Design

### 1. Standard Interface Structure

```go
// internal/services/user_service.go

package services

import (
    "context"
    "github.com/lusoris/venio/internal/models"
    "github.com/lusoris/venio/internal/repositories"
)

// UserService handles business logic for users
type UserService interface {
    // Registration & Authentication
    Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
    
    // CRUD operations
    GetByID(ctx context.Context, id int64) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByUsername(ctx context.Context, username string) (*models.User, error)
    Update(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error)
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, limit, offset int) ([]*models.User, int64, error)
}

// Private implementation
type userService struct {
    userRepository repositories.UserRepository
}

// Constructor returns interface
func NewUserService(userRepository repositories.UserRepository) UserService {
    return &userService{
        userRepository: userRepository,
    }
}
```

### 2. Service Return Types

**CRITICAL**: Services return pointers for collections (handlers need pointers).

| Operation | Service Return | Repository Return |
|-----------|---------------|-------------------|
| GetByID | `(*Model, error)` | `(*Model, error)` |
| Create | `(*Model, error)` | `(*Model, error)` |
| Update | `(*Model, error)` | `(*Model, error)` |
| Delete | `error` | `error` |
| List | `([]*Model, int64, error)` | `([]Model, int64, error)` |
| GetByUser | `([]*Model, error)` | `([]Model, error)` |

**Why different?**
- Repository returns values for efficiency
- Service converts to pointers for JSON marshaling in handlers
- Conversion happens once in service layer

### 3. Service Composition

When services depend on multiple repositories:

```go
type roleService struct {
    roleRepository       repositories.RoleRepository
    permissionRepository repositories.PermissionRepository
}

func NewRoleService(
    roleRepository repositories.RoleRepository,
    permissionRepository repositories.PermissionRepository,
) RoleService {
    return &roleService{
        roleRepository:       roleRepository,
        permissionRepository: permissionRepository,
    }
}

func (s *roleService) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
    // 1. Validate IDs
    if roleID <= 0 || permissionID <= 0 {
        return errors.New("invalid IDs")
    }

    // 2. Check role exists
    role, err := s.roleRepository.GetByID(ctx, roleID)
    if err != nil {
        return err
    }
    if role == nil {
        return errors.New("role not found")
    }

    // 3. Check permission exists
    perm, err := s.permissionRepository.GetByID(ctx, permissionID)
    if err != nil {
        return err
    }
    if perm == nil {
        return errors.New("permission not found")
    }

    // 4. Assign via repository
    return s.permissionRepository.AssignToRole(ctx, roleID, permissionID)
}
```

## Validation Patterns

### 1. Input Validation Order

```go
func (s *service) Create(ctx context.Context, req *Request) (*Model, error) {
    // 1. Struct validation (tags: required, email, min, max)
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("validation error: %w", err)
    }

    // 2. Business logic validation (duplicates, constraints)
    exists, err := s.repository.CheckDuplicate(ctx, req.UniqueField)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("already exists")
    }

    // 3. Proceed with operation
    // ...
}
```

### 2. Validation Error Messages

```go
// ✅ CORRECT: Descriptive validation errors
if id <= 0 {
    return nil, errors.New("invalid role ID: must be positive")
}

if email == "" {
    return nil, errors.New("validation error: email is required")
}

if len(password) < 8 {
    return nil, errors.New("validation error: password must be at least 8 characters")
}

// ❌ WRONG: Generic errors
if id <= 0 {
    return nil, errors.New("invalid ID")
}
```

## Error Handling

### 1. Error Wrapping

```go
// ✅ CORRECT: Wrap errors with context
user, err := s.userRepository.Create(ctx, req)
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// ❌ WRONG: Return raw error
user, err := s.userRepository.Create(ctx, req)
if err != nil {
    return nil, err
}
```

### 2. Not Found Handling

```go
// ✅ CORRECT: Convert repository (nil, nil) to service error
user, err := s.userRepository.GetByID(ctx, id)
if err != nil {
    return nil, err  // Database error
}
if user == nil {
    return nil, errors.New("user not found")  // Not found error
}
return user, nil

// ❌ WRONG: Don't check for nil
user, err := s.userRepository.GetByID(ctx, id)
if err != nil {
    return nil, err
}
return user, nil  // Could return nil user!
```

## Business Logic Examples

### 1. Password Hashing

```go
func (s *authService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    // Validate
    if err := req.Validate(); err != nil {
        return nil, err
    }

    // Hash password (business logic)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("hash password: %w", err)
    }
    req.Password = string(hashedPassword)

    // Create user
    return s.userRepository.Create(ctx, req)
}
```

### 2. Token Generation

```go
func (s *authService) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
    // 1. Get user
    user, err := s.userRepository.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid credentials")
    }

    // 2. Verify password
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // 3. Check active status
    if !user.IsActive {
        return nil, errors.New("account is inactive")
    }

    // 4. Generate JWT token (business logic)
    token, err := jwt.GenerateToken(user.ID, user.Email)
    if err != nil {
        return nil, fmt.Errorf("generate token: %w", err)
    }

    return &models.LoginResponse{
        User:  user,
        Token: token,
    }, nil
}
```

### 3. Permission Checks

```go
func (s *roleService) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
    // 1. Validate
    if roleID <= 0 || permissionID <= 0 {
        return errors.New("invalid IDs")
    }

    // 2. Check role exists
    role, err := s.roleRepository.GetByID(ctx, roleID)
    if err != nil {
        return err
    }
    if role == nil {
        return errors.New("role not found")
    }

    // 3. Check permission exists
    perm, err := s.permissionRepository.GetByID(ctx, permissionID)
    if err != nil {
        return err
    }
    if perm == nil {
        return errors.New("permission not found")
    }

    // 4. Business rule: Don't assign duplicate permissions
    existing, err := s.roleRepository.GetPermissions(ctx, roleID)
    if err != nil {
        return err
    }
    for _, p := range existing {
        if p.ID == permissionID {
            return nil  // Already assigned, silently succeed
        }
    }

    // 5. Assign
    return s.permissionRepository.AssignToRole(ctx, roleID, permissionID)
}
```

## Type Conversion Patterns

### 1. Repository Values → Service Pointers

```go
func (s *roleService) List(ctx context.Context, limit, offset int) ([]*models.Role, int64, error) {
    // Repository returns []models.Role
    roles, total, err := s.roleRepository.List(ctx, limit, offset)
    if err != nil {
        return nil, 0, err
    }

    // Convert to []*models.Role for handler
    rolePointers := make([]*models.Role, len(roles))
    for i := range roles {
        rolePointers[i] = &roles[i]
    }

    return rolePointers, total, nil
}
```

### 2. Request DTO → Repository DTO

```go
// Handler passes models.CreateUserRequest
// Service validates and passes to repository as-is
func (s *userService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    // Validate
    if err := req.Validate(); err != nil {
        return nil, err
    }

    // Transform if needed (e.g., hash password)
    req.Password = hashPassword(req.Password)

    // Pass to repository (same DTO type)
    return s.userRepository.Create(ctx, req)
}
```

## Service Factory Pattern

```go
// internal/services/factory.go

type Factory struct {
    userRepo       repositories.UserRepository
    roleRepo       repositories.RoleRepository
    permissionRepo repositories.PermissionRepository
    userRoleRepo   repositories.UserRoleRepository
}

func NewFactory(
    userRepo repositories.UserRepository,
    roleRepo repositories.RoleRepository,
    permissionRepo repositories.PermissionRepository,
    userRoleRepo repositories.UserRoleRepository,
) *Factory {
    return &Factory{
        userRepo:       userRepo,
        roleRepo:       roleRepo,
        permissionRepo: permissionRepo,
        userRoleRepo:   userRoleRepo,
    }
}

func (f *Factory) All() (
    UserService,
    AuthService,
    RoleService,
    PermissionService,
    UserRoleService,
) {
    userService := NewUserService(f.userRepo)
    authService := NewAuthService(f.userRepo)
    roleService := NewRoleService(f.roleRepo, f.permissionRepo)
    permissionService := NewPermissionService(f.permissionRepo)
    userRoleService := NewUserRoleService(f.userRoleRepo)

    return userService, authService, roleService, permissionService, userRoleService
}
```

## Checklist for New Services

- [ ] Interface defined with all methods
- [ ] Constructor returns interface, not concrete type
- [ ] Input validation at start of each method
- [ ] Business rule checks before repository calls
- [ ] Error wrapping with context
- [ ] Convert `(nil, nil)` from repository to `errors.New("not found")`
- [ ] Convert `[]Model` to `[]*Model` for list operations
- [ ] Check for duplicates in Create operations
- [ ] Check existence in Update/Delete operations
- [ ] Descriptive error messages
- [ ] No direct database access (use repositories)

## Common Mistakes

❌ **Don't:**
- Skip validation (services must validate)
- Return repository errors directly (wrap them)
- Skip existence checks in Update/Delete
- Return `nil` without error for not found
- Access database directly (use repositories)
- Put presentation logic in services (belongs in handlers)

✅ **Do:**
- Validate all inputs
- Wrap errors with context
- Check existence before Update/Delete
- Return clear error for not found
- Use repositories for all data access
- Keep services focused on business logic

---

**Remember**: Services orchestrate business logic. They validate, check rules, coordinate repositories, and transform data. Keep them thin but complete.
