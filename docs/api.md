# API Documentation

## Overview

Venio provides a RESTful API for all operations. Full OpenAPI/Swagger documentation is available at:

**Swagger UI:** `http://localhost:3690/swagger/`

## Authentication

### JWT Authentication

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2026-01-11T00:00:00Z"
}
```

Use token in subsequent requests:
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## API Endpoints

### Users

- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Requests

- `POST /api/v1/requests` - Create request
- `GET /api/v1/requests` - List requests
- `GET /api/v1/requests/:id` - Get request
- `PUT /api/v1/requests/:id/approve` - Approve request
- `DELETE /api/v1/requests/:id` - Cancel request

### Content

- `GET /api/v1/content/search` - Search content
- `GET /api/v1/content/:type/:id` - Get details

*Full endpoint documentation available in Swagger UI*

## Rate Limiting

API is rate limited to:
- **Authenticated:** 1000 requests/hour
- **Unauthenticated:** 100 requests/hour

Rate limit headers:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1609459200
```

## Error Responses

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "field": "email",
      "issue": "invalid format"
    }
  }
}
```

Common error codes:
- `UNAUTHORIZED` - 401
- `FORBIDDEN` - 403
- `NOT_FOUND` - 404
- `VALIDATION_ERROR` - 400
- `INTERNAL_ERROR` - 500

## Webhooks

Venio can send webhooks for events. Configure in settings.

Events:
- `request.created`
- `request.approved`
- `request.completed`
- `content.added`

---

*Detailed API reference: Coming soon*
