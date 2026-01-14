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
  - Authentication/Authorization
  - Request validation
  - Response formatting
  - OpenAPI/Swagger documentation

**Key Endpoints:**
- `/api/v1/auth/*` - Authentication
- `/api/v1/users/*` - User management
- `/api/v1/requests/*` - Content requests
- `/api/v1/content/*` - Content discovery
- `/api/v1/admin/*` - Admin operations

### 2. Services Layer (`internal/services`)

Business logic organized by domain:

- **User Service** - User management, RBAC
- **Request Service** - Request lifecycle, approvals
- **Metadata Service** - Multi-provider enrichment
- **Quality Service** - Profile management
- **Notification Service** - Multi-channel notifications
- **Analytics Service** - Stats and insights

### 3. Database Layer (`internal/database`)

- **Primary Database:** PostgreSQL 16+
- **Driver:** pgx v5 (native PostgreSQL driver)
- **Connection Pool:** pgxpool (built-in pooling)
- **Migrations:** golang-migrate

**Schema Organization:**
```
- users, roles, permissions
- requests, approvals
- metadata cache
- settings, profiles
- analytics events
```

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
- **Language:** Go 1.23
- **Web Framework:** Gin
- **Database:** PostgreSQL 16
- **Cache:** Redis 7
- **Search:** Typesense 29
- **Jobs:** Asynq
- **Config:** Viper

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
