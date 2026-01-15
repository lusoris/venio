# API Documentation

Venio provides comprehensive API documentation using Swagger/OpenAPI 3.0.

## Accessing the Documentation

### Swagger UI

Interactive API documentation is available at:

**Development:**
```
http://localhost:3690/swagger/index.html
```

**Production:**
```
https://api.venio.dev/swagger/index.html
```

### Features

- **Interactive Testing:** Try out API endpoints directly from the browser
- **Request/Response Examples:** See example payloads for all endpoints
- **Schema Documentation:** Detailed information about all data models
- **Authentication:** Test authenticated endpoints with Bearer tokens

## Regenerating Documentation

When you add or modify API endpoints, regenerate the Swagger documentation:

```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init --parseDependency --parseInternal -g cmd/venio/main.go -o docs/swagger
```

## Adding Swagger Annotations

### Handler Functions

Add Swagger annotations as comments above handler functions:

```go
// GetUser godoc
// @Summary Get user by ID
// @Description Get detailed information about a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // Implementation
}
```

### Main Package

Add API metadata in `cmd/venio/main.go`:

```go
// @title Venio API
// @version 1.0
// @description Enterprise-grade authentication and user management API
// @termsOfService https://venio.dev/terms

// @contact.name API Support
// @contact.url https://venio.dev/support
// @contact.email support@venio.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3690
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
```

### Model Structs

Add example tags to struct fields:

```go
type User struct {
    ID        int64     `json:"id" example:"1"`
    Email     string    `json:"email" example:"user@example.com"`
    Username  string    `json:"username" example:"johndoe"`
    FirstName string    `json:"first_name" example:"John"`
    LastName  string    `json:"last_name" example:"Doe"`
}
```

## Swagger Annotation Reference

### Common Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `@Summary` | Brief endpoint description | `@Summary Get user by ID` |
| `@Description` | Detailed description | `@Description Retrieves user information...` |
| `@Tags` | Group endpoints | `@Tags users` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Param` | Request parameter | `@Param id path int true "User ID"` |
| `@Success` | Success response | `@Success 200 {object} models.User` |
| `@Failure` | Error response | `@Failure 404 {object} ErrorResponse` |
| `@Security` | Authentication requirement | `@Security BearerAuth` |
| `@Router` | Endpoint path and method | `@Router /api/v1/users/{id} [get]` |

### Parameter Locations

- `path` - URL path parameter (e.g., `/users/{id}`)
- `query` - Query string parameter (e.g., `?page=1`)
- `header` - HTTP header
- `body` - Request body
- `formData` - Form data

### Parameter Types

- `string` - String value
- `integer` or `int` - Integer number
- `number` - Floating point number
- `boolean` or `bool` - Boolean value
- `object` - Complex object
- `array` - Array of values

## API Endpoints

### Authentication

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/auth/register` | POST | Register new user |
| `/api/v1/auth/login` | POST | Login and get tokens |
| `/api/v1/auth/refresh` | POST | Refresh access token |

### Users

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/v1/users` | GET | ✓ | List all users |
| `/api/v1/users/:id` | GET | ✓ | Get user by ID |
| `/api/v1/users/:id` | PUT | ✓ | Update user |
| `/api/v1/users/:id` | DELETE | ✓ | Delete user |
| `/api/v1/users/:id/roles` | GET | ✓ | Get user roles |
| `/api/v1/users/:id/roles` | POST | ✓ (Admin) | Assign role |
| `/api/v1/users/:id/roles/:roleId` | DELETE | ✓ (Admin) | Remove role |

### Roles

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/v1/roles` | GET | ✓ (Admin) | List all roles |
| `/api/v1/roles/:id` | GET | ✓ (Admin) | Get role by ID |
| `/api/v1/roles` | POST | ✓ (Admin) | Create role |
| `/api/v1/roles/:id` | PUT | ✓ (Admin) | Update role |
| `/api/v1/roles/:id` | DELETE | ✓ (Admin) | Delete role |
| `/api/v1/roles/:id/permissions` | GET | ✓ (Admin) | Get role permissions |
| `/api/v1/roles/:id/permissions` | POST | ✓ (Admin) | Assign permission |
| `/api/v1/roles/:id/permissions/:permId` | DELETE | ✓ (Admin) | Remove permission |

### Permissions

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/v1/permissions` | GET | ✓ (Admin) | List all permissions |
| `/api/v1/permissions/:id` | GET | ✓ (Admin) | Get permission by ID |
| `/api/v1/permissions` | POST | ✓ (Admin) | Create permission |
| `/api/v1/permissions/:id` | PUT | ✓ (Admin) | Update permission |
| `/api/v1/permissions/:id` | DELETE | ✓ (Admin) | Delete permission |

### Health Checks

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health/live` | GET | Liveness probe |
| `/health/ready` | GET | Readiness probe |

### Metrics & Docs

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/metrics` | GET | Prometheus metrics |
| `/swagger/*` | GET | Swagger UI |

## Authentication

Most endpoints require authentication using JWT Bearer tokens.

### Getting a Token

1. **Register** a new account:
   ```bash
   curl -X POST http://localhost:3690/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "username": "johndoe",
       "first_name": "John",
       "last_name": "Doe",
       "password": "SecurePass123!"
     }'
   ```

2. **Login** to get tokens:
   ```bash
   curl -X POST http://localhost:3690/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "SecurePass123!"
     }'
   ```

   Response:
   ```json
   {
     "access_token": "eyJhbGciOiJIUzI1NiIs...",
     "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
     "user": { ... }
   }
   ```

3. **Use the access token** in subsequent requests:
   ```bash
   curl -X GET http://localhost:3690/api/v1/users \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
   ```

### Refreshing Tokens

When the access token expires, use the refresh token:

```bash
curl -X POST http://localhost:3690/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }'
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error Type",
  "message": "Human-readable error message"
}
```

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request succeeded |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Service temporarily unavailable |

## Rate Limiting

API endpoints are rate limited to prevent abuse:

- **Authentication endpoints:** 5 requests per minute per IP
- **General API endpoints:** 100 requests per minute per user

Rate limit information is included in response headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642090800
```

When rate limited, you'll receive a `429 Too Many Requests` response.

## Testing the API

### Using Swagger UI

1. Navigate to `http://localhost:3690/swagger/index.html`
2. Click "Authorize" button
3. Enter: `Bearer <your-access-token>`
4. Try out any endpoint

### Using curl

See examples in the [Authentication](#authentication) section above.

### Using Postman

Import the OpenAPI spec into Postman:

1. In Postman, click "Import"
2. Select "Link"
3. Enter: `http://localhost:3690/swagger/doc.json`
4. Click "Import"

## Best Practices

1. **Always use HTTPS in production**
2. **Store tokens securely** (HttpOnly cookies or secure storage)
3. **Refresh tokens before they expire**
4. **Handle rate limits gracefully** (exponential backoff)
5. **Validate input** on client-side before sending
6. **Handle errors appropriately** - don't expose internal details to users

## Further Reading

- [Swagger Specification](https://swagger.io/specification/)
- [Swaggo Documentation](https://github.com/swaggo/swag)
- [REST API Best Practices](https://restfulapi.net/)
