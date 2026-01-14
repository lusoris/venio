# Repositories Package

This package contains the repository pattern implementations for data access layer.

## Structure

- `user_repository.go` - User data access layer with CRUD operations

## Architecture

The repository pattern provides a clean abstraction between the service layer and database layer:

- **Interface**: Defines the contract for user data access operations
- **Implementation**: PostgreSQL implementation using pgx and pgxpool
- **Queries**: Direct SQL queries with prepared statement patterns for security

## Usage

```go
// Initialize repository
userRepo := repositories.NewPostgresUserRepository(pool)

// Get user by ID
user, err := userRepo.GetByID(ctx, 1)

// Create new user
id, err := userRepo.Create(ctx, &models.User{
    Email: "user@example.com",
    Username: "john_doe",
    // ... other fields
})

// Update user
err := userRepo.Update(ctx, user)

// List users with pagination
users, err := userRepo.List(ctx, 10, 0)
```

## Security Considerations

- All queries use parameterized statements (`$1`, `$2`, etc.) to prevent SQL injection
- Password field is only selected in methods that need it (GetByEmail, GetByUsername)
- GetByID specifically excludes password for security
