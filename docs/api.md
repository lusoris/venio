# Venio API Documentation

## Base URL
```
http://localhost:3690/api/v1
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <access_token>
```

## Endpoints

### Authentication

#### Register
Creates a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "password": "SecurePassword123!"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "is_active": true,
  "created_at": "2026-01-14T12:00:00Z",
  "updated_at": "2026-01-14T12:00:00Z"
}
```

#### Login
Authenticates a user and returns JWT tokens.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  }
}
```

#### Refresh Token
Gets a new access token using a refresh token.

**Endpoint:** `POST /auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Users (Protected)

All user endpoints require authentication.

#### List Users
Get a paginated list of users.

**Endpoint:** `GET /users?limit=10&offset=0`

**Query Parameters:**
- `limit` (optional): Number of users to return (default: 10, max: 100)
- `offset` (optional): Number of users to skip (default: 0)

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  }
]
```

#### Get User
Get a specific user by ID.

**Endpoint:** `GET /users/:id`

**Response:** `200 OK`
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "is_active": true,
  "created_at": "2026-01-14T12:00:00Z",
  "updated_at": "2026-01-14T12:00:00Z"
}
```

#### Update User
Update user information.

**Endpoint:** `PUT /users/:id`

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "first_name": "Jane",
  "is_active": false
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "email": "newemail@example.com",
  "username": "johndoe",
  "first_name": "Jane",
  "last_name": "Doe",
  "is_active": false,
  "created_at": "2026-01-14T12:00:00Z",
  "updated_at": "2026-01-14T13:00:00Z"
}
```

#### Delete User
Delete a user by ID.

**Endpoint:** `DELETE /users/:id`

**Response:** `204 No Content`

## Error Responses

All errors return a consistent structure:

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

**Common Status Codes:**
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Example Usage

### Using cURL

```bash
# Register a new user
curl -X POST http://localhost:3690/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "password": "SecurePassword123!"
  }'

# Login
curl -X POST http://localhost:3690/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123!"
  }'

# Get users (with token)
curl -X GET http://localhost:3690/api/v1/users \
  -H "Authorization: Bearer <your_access_token>"
```

### Using PowerShell

```powershell
# Register a new user
$body = @{
    email = "user@example.com"
    username = "johndoe"
    first_name = "John"
    last_name = "Doe"
    password = "SecurePassword123!"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:3690/api/v1/auth/register" `
    -Method POST -Body $body -ContentType "application/json"

# Login
$loginBody = @{
    email = "user@example.com"
    password = "SecurePassword123!"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri "http://localhost:3690/api/v1/auth/login" `
    -Method POST -Body $loginBody -ContentType "application/json"

$token = ($response.Content | ConvertFrom-Json).access_token

# Get users (with token)
$headers = @{ Authorization = "Bearer $token" }
Invoke-WebRequest -Uri "http://localhost:3690/api/v1/users" `
    -Method GET -Headers $headers
```

## JWT Token Claims

Access tokens contain the following claims:

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "roles": [],
  "iss": "venio",
  "exp": 1768482270,
  "iat": 1768395870
}
```

- `user_id`: User's database ID
- `email`: User's email address
- `username`: User's username
- `roles`: Array of user's roles (TODO: populated from database)
- `iss`: Token issuer (always "venio")
- `exp`: Token expiration timestamp
- `iat`: Token issued at timestamp

## Security Notes

1. Passwords are hashed using bcrypt before storage
2. Passwords are never returned in API responses
3. All protected routes require valid JWT authentication
4. Access tokens expire after 24 hours (configurable)
5. Refresh tokens expire after 7 days (configurable)
6. SQL injection protection via parameterized queries
