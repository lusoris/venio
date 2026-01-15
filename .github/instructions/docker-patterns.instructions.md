---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "Dockerfile,Dockerfile.*,docker-compose*.yml,.dockerignore"
description: Docker & Containerization Best Practices
---

# Docker & Containerization Best Practices

## Core Principle

**Consistency Across Environments**: Development, staging, and production should use identical Go versions. Multi-stage builds for minimal, secure images.

## Go Version Consistency (CRITICAL)

### ✅ CORRECT: Same Version Everywhere

```dockerfile
# Dockerfile (production)
FROM golang:1.25-alpine AS builder

# Dockerfile.dev (development)
FROM golang:1.25-alpine
```

### ❌ WRONG: Version Mismatch

```dockerfile
# Dockerfile (production)
FROM golang:1.25-alpine  # ❌ Latest Go

# Dockerfile.dev (development)
FROM golang:1.23-alpine  # ❌ Different version!
# Risk: Environment-specific bugs, incompatible dependencies
```

### Version Update Strategy

When updating Go version:
1. Update BOTH Dockerfile and Dockerfile.dev simultaneously
2. Update go.mod if needed: `go mod edit -go=1.25`
3. Test in development first
4. Deploy to staging
5. Deploy to production

## Multi-Stage Builds (Production)

### ✅ CORRECT: Minimal Production Image

```dockerfile
# Dockerfile
# ===================================================================
# Stage 1: Build
# ===================================================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/venio \
    ./cmd/venio/main.go

# ===================================================================
# Stage 2: Runtime
# ===================================================================
FROM alpine:3.21

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 venio && \
    adduser -D -u 1000 -G venio venio

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/venio /app/venio

# Copy migrations (if needed)
COPY --from=builder /app/migrations /app/migrations

# Change ownership
RUN chown -R venio:venio /app

# Switch to non-root user
USER venio

# Expose port
EXPOSE 3690

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/app/venio", "health"]

# Run application
CMD ["/app/venio"]
```

## Development Dockerfile

### ✅ CORRECT: Hot-Reload Friendly

```dockerfile
# Dockerfile.dev
FROM golang:1.25-alpine

# Install development tools
RUN apk add --no-cache git make gcc musl-dev

# Install air for hot-reload
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (or use volumes)
COPY . .

# Expose port
EXPOSE 3690

# Use air for hot-reload
CMD ["air", "-c", ".air.toml"]
```

## Docker Compose (Development)

### ✅ CORRECT: Development Stack

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  # Backend API
  api:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "3690:3690"
    volumes:
      # Mount source code for hot-reload
      - .:/app
      - go-modules:/go/pkg/mod  # Cache Go modules
    environment:
      - VENIO_APP_ENV=development
      - VENIO_DATABASE_HOST=postgres
      - VENIO_DATABASE_PORT=5432
      - VENIO_DATABASE_USER=venio
      - VENIO_DATABASE_PASSWORD=venio-dev-password
      - VENIO_DATABASE_NAME=venio
      - VENIO_REDIS_HOST=redis
      - VENIO_REDIS_PORT=6379
      - VENIO_JWT_SECRET=dev-secret-at-least-32-characters-long
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - venio-network

  # PostgreSQL
  postgres:
    image: postgres:18.1-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=venio
      - POSTGRES_PASSWORD=venio-dev-password
      - POSTGRES_DB=venio
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U venio"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - venio-network

  # Redis
  redis:
    image: redis:8.4-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - venio-network

  # Frontend (Next.js)
  web:
    build:
      context: ./web
      dockerfile: Dockerfile.dev
    ports:
      - "3000:3000"
    volumes:
      - ./web:/app
      - /app/node_modules
      - /app/.next
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:3690
    depends_on:
      - api
    networks:
      - venio-network

volumes:
  postgres-data:
  redis-data:
  go-modules:

networks:
  venio-network:
    driver: bridge
```

## Docker Compose (Production)

### ✅ CORRECT: Production Stack

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3690:3690"
    environment:
      - VENIO_APP_ENV=production
      - VENIO_DATABASE_HOST=postgres
      - VENIO_REDIS_HOST=redis
    env_file:
      - .env.production  # Secrets from .env file
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - venio-network
    healthcheck:
      test: ["CMD", "/app/venio", "health"]
      interval: 30s
      timeout: 5s
      retries: 3

  postgres:
    image: postgres:18.1-alpine
    environment:
      - POSTGRES_USER=${DATABASE_USER}
      - POSTGRES_PASSWORD=${DATABASE_PASSWORD}
      - POSTGRES_DB=${DATABASE_NAME}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - venio-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DATABASE_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:8.4-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data
    restart: unless-stopped
    networks:
      - venio-network
    healthcheck:
      test: ["CMD", "redis-cli", "--pass", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres-data:
  redis-data:

networks:
  venio-network:
    driver: bridge
```

## .dockerignore

### ✅ CORRECT: Exclude Unnecessary Files

```dockerignore
# Git
.git
.gitignore

# Environment
.env
.env.*
!.env.example

# Documentation
*.md
docs/

# Tests
*_test.go
**/*_test.go

# IDE
.vscode/
.idea/
*.swp
*.swo

# Build artifacts
bin/
build/
dist/
*.exe

# Dependencies (will be downloaded)
vendor/

# Logs
*.log

# OS
.DS_Store
Thumbs.db

# Docker
docker-compose*.yml
Dockerfile*
.dockerignore
```

## Build Optimization

### Layer Caching Strategy

```dockerfile
# ✅ CORRECT: Leverage layer caching

# 1. Base dependencies (rarely change)
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git make

# 2. Go modules (change occasionally)
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download  # ✅ Cached if go.mod/go.sum unchanged

# 3. Source code (changes frequently)
COPY . .
RUN go build -o venio ./cmd/venio/main.go

# ❌ WRONG: Copy everything first
COPY . .
RUN go mod download  # ❌ Re-downloads on any file change!
RUN go build -o venio ./cmd/venio/main.go
```

### Build Arguments

```dockerfile
# Dockerfile
ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder

ARG BUILD_VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN go build \
    -ldflags="-w -s \
              -X main.Version=${BUILD_VERSION} \
              -X main.BuildTime=${BUILD_TIME} \
              -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/venio \
    ./cmd/venio/main.go

# Build command:
# docker build \
#   --build-arg BUILD_VERSION=1.0.0 \
#   --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
#   --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
#   -t venio:1.0.0 .
```

## Security Best Practices

### ✅ CORRECT: Non-Root User

```dockerfile
# Create non-root user
RUN addgroup -g 1000 venio && \
    adduser -D -u 1000 -G venio venio

# Change ownership
RUN chown -R venio:venio /app

# Switch to non-root
USER venio

# ❌ WRONG: Run as root
# Default if you don't specify USER
```

### Minimal Base Image

```dockerfile
# ✅ CORRECT: Alpine Linux (small, secure)
FROM alpine:3.21

# Install only what's needed
RUN apk add --no-cache ca-certificates tzdata

# ❌ WRONG: Large base images
FROM ubuntu:latest  # ❌ ~70MB vs alpine's ~5MB
FROM golang:1.25    # ❌ ~300MB (includes build tools)
```

### Distroless Alternative

```dockerfile
# Alternative: Distroless (even smaller, more secure)
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/venio /venio
ENTRYPOINT ["/venio"]

# Benefits:
# - No shell (reduces attack surface)
# - No package manager
# - Minimal dependencies
# - ~2MB image size
```

## Health Checks

### Application Health Check

```dockerfile
# In Dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/app/venio", "health"]
```

### Health Check Handler

```go
// cmd/venio/main.go
func main() {
    if len(os.Args) > 1 && os.Args[1] == "health" {
        if err := healthCheck(); err != nil {
            os.Exit(1)
        }
        os.Exit(0)
    }

    // Normal startup
    // ...
}

func healthCheck() error {
    resp, err := http.Get("http://localhost:3690/health")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unhealthy status: %d", resp.StatusCode)
    }

    return nil
}
```

## Common Docker Commands

### Development

```bash
# Build development image
docker compose -f docker-compose.dev.yml build

# Start development stack
docker compose -f docker-compose.dev.yml up

# Stop stack
docker compose -f docker-compose.dev.yml down

# View logs
docker compose -f docker-compose.dev.yml logs -f api

# Rebuild single service
docker compose -f docker-compose.dev.yml build --no-cache api

# Execute command in running container
docker compose -f docker-compose.dev.yml exec api sh

# Run tests
docker compose -f docker-compose.dev.yml exec api go test ./...
```

### Production

```bash
# Build production image
docker build -t venio:latest .

# Tag for registry
docker tag venio:latest registry.example.com/venio:1.0.0

# Push to registry
docker push registry.example.com/venio:1.0.0

# Run production stack
docker compose up -d

# Scale service
docker compose up -d --scale api=3

# View resource usage
docker stats

# Prune unused resources
docker system prune -a
```

## Debugging Containers

### Inspect Container

```bash
# View container logs
docker logs venio-api-1

# Execute shell in running container
docker exec -it venio-api-1 sh

# Inspect container details
docker inspect venio-api-1

# View environment variables
docker exec venio-api-1 env

# Check running processes
docker exec venio-api-1 ps aux
```

### Common Issues

```bash
# Issue: Container exits immediately
# Check logs
docker logs venio-api-1
# Common causes: Missing env vars, connection failures

# Issue: Database connection refused
# Check network
docker network inspect venio-network
# Verify service names match (postgres, not localhost)

# Issue: Port already in use
# Find process using port
netstat -ano | findstr :3690  # Windows
lsof -i :3690                 # Linux/macOS

# Issue: Build cache issues
# Clear cache and rebuild
docker builder prune
docker compose build --no-cache
```

## Container Checklist

- [ ] **Same Go version in all Dockerfiles**
- [ ] **Multi-stage build for production**
- [ ] **Non-root user in production image**
- [ ] **Minimal base image (alpine or distroless)**
- [ ] **.dockerignore to exclude unnecessary files**
- [ ] **Layer caching optimized (go.mod before COPY . .)**
- [ ] **Health checks defined**
- [ ] **Secrets from environment, not baked in**
- [ ] **Volume mounts for development hot-reload**
- [ ] **Named volumes for persistent data**
- [ ] **Health checks in docker-compose**
- [ ] **Restart policy for production**

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial Docker best practices guide |

**Remember**: Development and production MUST use the same Go version. Multi-stage builds keep production images small and secure. Never bake secrets into images.
