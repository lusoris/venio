# Services Package

This package contains the business logic layer implementations.

## Structure

- `user_service.go` - User business logic and validation

## Architecture

The service layer sits between the HTTP handlers and repositories:

- **Repositories**: Handle data access and persistence
- **Services**: Implement business logic, validation, and orchestration
- **Handlers**: Handle HTTP requests and responses

## Features

### User Service

- User registration with validation and password hashing
- User retrieval by ID, email, or username
- User updates with field-level validation
- User deletion
- Paginated user listing
- Email validation
- Username length validation (3-50 characters)
- Automatic timestamp management

## Password Security

Passwords are hashed using bcrypt with the default cost factor before storage:

```go
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

## Validation

Services perform validation on input requests:

- Email format validation using regex
- Username length validation (3-50 chars)
- Password strength validation (≥8 chars)
- Duplicate email detection before user creation

## Error Handling

All service methods return structured errors with context:

```go
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

This allows callers to use `errors.Is()` and `errors.As()` for specific error handling.
