# Venio Architecture

## Overview

Venio is built as a modular orchestration layer that sits between users and their media automation stack (*arr apps, media servers, etc.).

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Users (Web/Mobile)                    │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   Venio Backend                          │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐ │
│  │   API       │  │   Services   │  │   Workers      │ │
│  │  (Gin)      │  │   (Business  │  │   (Asynq)      │ │
│  │             │  │    Logic)    │  │                │ │
│  └─────────────┘  └──────────────┘  └────────────────┘ │
│         │                 │                   │          │
│  ┌─────▼─────────────────▼───────────────────▼───────┐ │
│  │          Database Layer (PostgreSQL)              │ │
│  └───────────────────────────────────────────────────┘ │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│              External Services                           │
│  ┌──────────┐ ┌─────────┐ ┌─────────┐ ┌──────────────┐│
│  │Overseerr │ │ Lidarr  │ │Whisparr │ │Media Servers ││
│  └──────────┘ └─────────┘ └─────────┘ └──────────────┘│
└──────────────────────────────────────────────────────────┘
```

## Core Components

### 1. API Layer (`internal/api`)

- **Framework:** Gin (Go)
- **Responsibilities:**
  - HTTP request handling
  - Authentication/Authorization (JWT)
  - RBAC enforcement via middleware
  - Request validation
  - Response formatting
  - OpenAPI/Swagger documentation

**Key Endpoints:**
- `/api/v1/auth/*` - Authentication (register, login, refresh)
- `/api/v1/users/*` - User management
- `/api/v1/roles/*` - Role management (admin only)
- `/api/v1/permissions/*` - Permission management (admin only)
- `/api/v1/users/:id/roles/*` - User-role assignments
- `/api/v1/requests/*` - Content requests
- `/api/v1/content/*` - Content discovery
- `/api/v1/admin/*` - Admin operations

**RBAC Middleware** (`internal/api/middleware/rbac.go`):
- `RequireRole(roleName)` - Ensures user has specific role
- `RequirePermission(permissionName)` - Ensures user has specific permission
- `RequireAnyRole(...roleNames)` - Ensures user has at least one of the roles
- `RequireAnyPermission(...permissionNames)` - Ensures user has at least one permission

Authorization flow:
1. AuthMiddleware validates JWT and extracts user ID
2. RBAC middleware checks UserRoleService for role/permission
3. Returns 401 (Unauthorized) if not authenticated
4. Returns 403 (Forbidden) if missing required role/permission
5. Proceeds to handler if authorized

### 2. Services Layer (`internal/services`)

Business logic organized by domain:

- **User Service** - User management, basic operations
- **Role Service** - Role CRUD, role-permission management
- **Permission Service** - Permission CRUD, user permission queries
- **User Role Service** - User-role assignments, authorization checks
- **Request Service** - Request lifecycle, approvals
- **Metadata Service** - Multi-provider enrichment
- **Quality Service** - Profile management
- **Notification Service** - Multi-channel notifications
- **Analytics Service** - Stats and insights

**RBAC Service Details:**

**RoleService:**
- Create, Read, Update, Delete roles
- List roles with pagination
- Get role permissions
- Assign/Remove permissions to/from roles
- Duplicate name validation

**PermissionService:**
- Create, Read, Update, Delete permissions
- List permissions with pagination
- Query user permissions
- Permission validation

**UserRoleService:**
- Assign/Remove roles to/from users
- Get all user roles
- Check if user has specific role (`HasRole`)
- Check if user has specific permission (`HasPermission`)
- Authorization checks for middleware

### 3. Database Layer (`internal/database`)

- **Primary Database:** PostgreSQL 18.1 (Alpine)
- **Driver:** pgx v5 (native PostgreSQL driver)
- **Connection Pool:** pgxpool (built-in pooling)
- **Migrations:** golang-migrate

**Schema Organization:**
```
- users, roles, permissions
- user_roles, role_permissions (junction tables)
- requests, approvals
- metadata cache
- settings, profiles
- analytics events
```

**RBAC Repositories:**

**RoleRepository** (`internal/repositories/role_repository.go`):
- `Create(role)` - Create new role
- `GetByID(id)` - Get role by ID
- `GetByName(name)` - Get role by name
- `Update(role)` - Update role details
- `Delete(id)` - Delete role
- `List(limit, offset)` - Paginated role listing
- `GetPermissions(roleID)` - Get all permissions for role

**PermissionRepository** (`internal/repositories/permission_repository.go`):
- `Create(permission)` - Create new permission
- `GetByID(id)` - Get permission by ID
- `GetByName(name)` - Get permission by name
- `Update(permission)` - Update permission details
- `Delete(id)` - Delete permission
- `List(limit, offset)` - Paginated permission listing
- `GetByUserID(userID)` - Get all permissions for user (via roles)
- `AssignToRole(roleID, permissionID)` - Assign permission to role
- `RemoveFromRole(roleID, permissionID)` - Remove permission from role

**UserRoleRepository** (`internal/repositories/user_role_repository.go`):
- `AssignRole(userID, roleID)` - Assign role to user
- `RemoveRole(userID, roleID)` - Remove role from user
- `GetUserRoles(userID)` - Get all roles for user
- `HasRole(userID, roleName)` - Check if user has role
- `HasPermission(userID, permissionName)` - Check if user has permission

**Database Pattern:**
- Raw SQL queries with pgx for flexibility
- Connection pooling with configurable limits
- Proper error handling and logging
- Repository pattern for data access

### 4. Workers (`cmd/worker`)

Background job processing with Asynq:

- **Job Types:**
  - Metadata fetching
  - Arr synchronization
  - Analytics computation
  - Retention rule execution
  - Notification delivery

### 5. Proxy Layer (`internal/proxy`)

Metadata API proxy that intercepts and enriches requests:

```
Arr/Media Server → Venio Proxy → Cache Check → Enrichment → Provider API
                                      ↓
                                   Return
```

### 6. Providers (`internal/providers`)

External API clients:
- Overseerr client
- Arr clients (Sonarr, Radarr, Lidarr, Whisparr)
- Media server clients (Jellyfin, Plex, Emby)
- Metadata providers (TMDB, TVDB, MusicBrainz, etc.)

## Data Flow

### Request Creation Flow

```
1. User creates request via UI
2. API validates request
3. Request Service:
   - Checks permissions
   - Applies auto-approval rules
   - Creates database record
4. If approved:
   - Worker sends to Overseerr/Arr
   - Tracks status
5. Notifications sent
```

### Metadata Enrichment Flow

```
1. Request arrives at Metadata Proxy
2. Check cache (Redis)
3. If cache miss:
   - Fetch from multiple providers in parallel
   - Merge results
   - Apply overrides
   - Cache result
4. Return enriched data
```

## Technology Stack

### Backend
- **Language:** Go 1.25
- **Web Framework:** Gin 1.10.0
- **Database:** PostgreSQL 18.1 (Alpine)
- **Cache:** Redis 8.4 (Alpine)
- **Search:** Typesense 29
- **Jobs:** Asynq
- **Config:** Viper

**PostgreSQL 18.1 Features:**
- Asynchronous I/O (AIO) subsystem for better performance
- Skip scan support for multicolumn B-tree indexes
- Virtual generated columns (compute during read)
- OAuth authentication support
- Improved parallel query execution
- Enhanced vacuum and analyze operations
- UUID v7 generation (timestamp-ordered)

**Redis 8.4 Benefits:**
- Enhanced performance and memory efficiency
- Improved clustering capabilities
- Better TLS/SSL support
- Extended command compatibility

### Frontend (Future)
- React + TypeScript
- shadcn/ui (Tailwind)
- Vite + React Router

### Infrastructure
- **Containers:** Docker
- **Orchestration:** Docker Compose / Kubernetes
- **Registry:** GitHub Container Registry (ghcr.io)

## Security

### Authentication
- JWT tokens
- OIDC integration (Authentik, Authelia, etc.)
- Session management (Redis)

### Authorization
- Role-Based Access Control (RBAC)
- Permission system
- Per-module access control

### Data Protection
- Encrypted API keys in database
- Secure session tokens
- Adult content isolation
- Audit logging

## Scalability

### Horizontal Scaling
- Stateless API servers
- Shared PostgreSQL/Redis
- Load balancer in front

### Caching Strategy
- Metadata cache (Redis)
- API response cache
- Static asset CDN (future)

### Database
- Connection pooling
- Read replicas (future)
- Partitioning for analytics

## Monitoring & Observability

- Structured logging (JSON)
- Health check endpoints
- Metrics (Prometheus-compatible)
- Error tracking
- Performance monitoring

## Development Principles

1. **Modularity** - Loosely coupled components
2. **Testability** - High test coverage (>80%)
3. **API-First** - OpenAPI spec drives development
4. **Configuration** - Environment-based config
5. **Documentation** - Code comments + external docs

## Future Architecture Considerations

- **Microservices:** Split into separate services if needed
- **Event Sourcing:** For audit trail
- **CQRS:** Separate read/write paths
- **GraphQL:** Alternative to REST
- **gRPC:** For service-to-service communication

---

*This is a living document and will be updated as architecture evolves.*
