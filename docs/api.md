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

## Roles & Permissions (RBAC)

### Role Endpoints

All role endpoints require admin authentication.

#### Create Role
Create a new role.

**Endpoint:** `POST /roles`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "name": "moderator",
  "description": "Content moderator with limited admin rights"
}
```

**Response:** `201 Created`
```json
{
  "id": 4,
  "name": "moderator",
  "description": "Content moderator with limited admin rights",
  "created_at": "2026-01-14T15:00:00Z",
  "updated_at": "2026-01-14T15:00:00Z"
}
```

#### List Roles
Get a paginated list of all roles.

**Endpoint:** `GET /roles?page=1&limit=10`

**Authorization:** Requires `admin` role

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "admin",
    "description": "Full system administrator",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  },
  {
    "id": 2,
    "name": "user",
    "description": "Standard user with basic permissions",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  }
]
```

#### Get Role
Get a specific role by ID.

**Endpoint:** `GET /roles/:id`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "admin",
  "description": "Full system administrator",
  "created_at": "2026-01-14T12:00:00Z",
  "updated_at": "2026-01-14T12:00:00Z"
}
```

#### Update Role
Update role information.

**Endpoint:** `PUT /roles/:id`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "name": "super_admin",
  "description": "Super administrator with all permissions"
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "super_admin",
  "description": "Super administrator with all permissions",
  "created_at": "2026-01-14T12:00:00Z",
  "updated_at": "2026-01-14T15:30:00Z"
}
```

#### Delete Role
Delete a role by ID.

**Endpoint:** `DELETE /roles/:id`

**Authorization:** Requires `admin` role

**Response:** `204 No Content`

#### Get Role Permissions
Get all permissions assigned to a role.

**Endpoint:** `GET /roles/:id/permissions`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "user:read",
    "description": "Read user information",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  },
  {
    "id": 2,
    "name": "user:write",
    "description": "Create and modify users",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  }
]
```

#### Assign Permission to Role
Add a permission to a role.

**Endpoint:** `POST /roles/:id/permissions`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "permission_id": 5
}
```

**Response:** `200 OK`
```json
{
  "message": "Permission assigned to role successfully"
}
```

#### Remove Permission from Role
Remove a permission from a role.

**Endpoint:** `DELETE /roles/:roleId/permissions/:permissionId`

**Authorization:** Requires `admin` role

**Response:** `204 No Content`

### Permission Endpoints

All permission endpoints require admin authentication.

#### Create Permission
Create a new permission.

**Endpoint:** `POST /permissions`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "name": "media:delete",
  "description": "Delete media items"
}
```

**Response:** `201 Created`
```json
{
  "id": 9,
  "name": "media:delete",
  "description": "Delete media items",
  "created_at": "2026-01-14T15:00:00Z",
  "updated_at": "2026-01-14T15:00:00Z"
}
```

#### List Permissions
Get a paginated list of all permissions.

**Endpoint:** `GET /permissions?page=1&limit=10`

**Authorization:** Requires `admin` role

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "user:read",
    "description": "Read user information",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  },
  {
    "id": 2,
    "name": "user:write",
    "description": "Create and modify users",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  }
]
```

#### Get Permission
Get a specific permission by ID.

**Endpoint:** `GET /permissions/:id`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "user:read",
  "description": "Read user information",
  "created_at": "2026-01-14T12:00:00Z",
  "updated_at": "2026-01-14T12:00:00Z"
}
```

#### Update Permission
Update permission information.

**Endpoint:** `PUT /permissions/:id`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "name": "user:read:all",
  "description": "Read all user information including sensitive data"
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "user:read:all",
  "description": "Read all user information including sensitive data",
  "created_at": "2026-01-14T12:00:00Z",
  "updated_at": "2026-01-14T15:45:00Z"
}
```

#### Delete Permission
Delete a permission by ID.

**Endpoint:** `DELETE /permissions/:id`

**Authorization:** Requires `admin` role

**Response:** `204 No Content`

### User-Role Endpoints

#### Get User Roles
Get all roles assigned to a user.

**Endpoint:** `GET /users/:userId/roles`

**Authorization:** Authenticated users can view their own roles, admin can view any user's roles

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "admin",
    "description": "Full system administrator",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  },
  {
    "id": 2,
    "name": "user",
    "description": "Standard user with basic permissions",
    "created_at": "2026-01-14T12:00:00Z",
    "updated_at": "2026-01-14T12:00:00Z"
  }
]
```

#### Assign Role to User
Add a role to a user.

**Endpoint:** `POST /users/:userId/roles`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "role_id": 2
}
```

**Response:** `200 OK`
```json
{
  "message": "Role assigned to user successfully"
}
```

#### Remove Role from User
Remove a role from a user.

**Endpoint:** `DELETE /users/:userId/roles/:roleId`

**Authorization:** Requires `admin` role

**Response:** `204 No Content`

### Admin Dashboard

The admin dashboard provides dedicated endpoints for managing users, roles, and assignments. These endpoints are designed for administrative UIs and return data optimized for management views.

#### List All Users (Admin)
Get all users with pagination and extended details.

**Endpoint:** `GET /admin/users`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "is_active": true,
      "created_at": "2026-01-14T12:00:00Z"
    }
  ]
}
```

#### Create User (Admin)
Create a new user with role assignments.

**Endpoint:** `POST /admin/users`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "username": "newuser",
  "first_name": "Jane",
  "last_name": "Smith",
  "password": "SecurePassword123!",
  "roles": [1, 2]
}
```

**Response:** `201 Created`
```json
{
  "id": 5,
  "email": "newuser@example.com",
  "username": "newuser",
  "firstName": "Jane",
  "lastName": "Smith"
}
```

#### Delete User (Admin)
Delete a user from the system.

**Endpoint:** `DELETE /admin/users/:id`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "message": "User deleted successfully"
}
```

#### List All Roles (Admin)
Get all roles with user counts.

**Endpoint:** `GET /admin/roles`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "roles": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrator role",
      "user_count": 2,
      "created_at": "2026-01-14T12:00:00Z"
    }
  ]
}
```

#### Create Role (Admin)
Create a new role with permission assignments.

**Endpoint:** `POST /admin/roles`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "name": "moderator",
  "description": "Content moderator",
  "permissions": [3, 5, 7]
}
```

**Response:** `201 Created`
```json
{
  "id": 4,
  "name": "moderator"
}
```

#### Delete Role (Admin)
Delete a role from the system.

**Endpoint:** `DELETE /admin/roles/:id`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "message": "Role deleted successfully"
}
```

#### List All Permissions (Admin)
Get all available permissions.

**Endpoint:** `GET /admin/permissions`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "permissions": [
    {
      "id": 1,
      "name": "users.read",
      "description": "Read user data",
      "created_at": "2026-01-14T12:00:00Z"
    }
  ]
}
```

#### List User-Role Assignments (Admin)
Get all user-role assignments for management.

**Endpoint:** `GET /admin/user-roles`

**Authorization:** Requires `admin` role

**Response:** `200 OK`
```json
{
  "assignments": [
    {
      "user_id": 1,
      "user_email": "user@example.com",
      "role_name": "admin",
      "assigned_at": "2026-01-14T12:00:00Z"
    }
  ]
}
```

#### Remove User-Role Assignment (Admin)
Remove a specific role assignment from a user.

**Endpoint:** `DELETE /admin/user-roles/:id`

**Authorization:** Requires `admin` role

**Request Body:**
```json
{
  "user_id": 1,
  "role_id": 2
}
```

**Response:** `200 OK`
```json
{
  "message": "Role assignment removed successfully"
}
```

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
- `403 Forbidden`: Insufficient permissions (e.g., not admin)
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (e.g., duplicate role name)
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
  "roles": ["admin", "user"],
  "iss": "venio",
  "exp": 1768482270,
  "iat": 1768395870
}
```

- `user_id`: User's database ID
- `email`: User's email address
- `username`: User's username
- `roles`: Array of user's role names
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
