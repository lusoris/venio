# Development Guide

This guide covers setting up a development environment for Venio.

## Prerequisites

### Required

- **Go 1.23+** - [Download](https://go.dev/dl/)
- **Docker 20+** - [Install](https://docs.docker.com/get-docker/)
- **Docker Compose** - [Install](https://docs.docker.com/compose/install/)
- **Git** - [Install](https://git-scm.com/downloads)

### Recommended

- **VSCode** with Go extension
- **Make** (usually pre-installed on Linux/macOS)
- **golangci-lint** - [Install](https://golangci-lint.run/usage/install/)
- **air** (hot reload) - `go install github.com/cosmtrek/air@latest`

## Initial Setup

### 1. Clone Repository

```bash
git clone https://github.com/lusoris/venio.git
cd venio
```

### 2. Install Development Tools

```bash
# Go tools
go install github.com/cosmtrek/air@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lefthook (pre-commit hooks)
go install github.com/evilmartians/lefthook@latest
lefthook install
```

### 3. Environment Configuration

```bash
cp .env.example .env
# Edit .env with your settings
```

Key settings to configure:
```env
POSTGRES_PASSWORD=your_secure_password
REDIS_PASSWORD=your_redis_password
JWT_SECRET=your_jwt_secret_min_32_chars
```

### 4. VSCode Setup

If using VSCode, recommended extensions will be suggested automatically.

Install them for the best experience:
- Go (golang.go)
- Docker (ms-azuretools.vscode-docker)
- GitLens (eamodio.gitlens)

## Running Locally

### Option 1: Services Only (Recommended for Development)

```bash
# Terminal 1: Start Docker services
docker compose up postgres redis

# Terminal 2: Run Venio locally
go run cmd/venio/main.go
```

This is best for:
- Faster iteration (no Docker rebuild)
- Easier debugging with IDE
- Hot reload with `air`

Access:
- Venio API: http://localhost:3690
- PostgreSQL: localhost:5432 (in Docker)
- Redis: localhost:6379 (in Docker)

**Health Check:**
```bash
curl http://localhost:3690/health
```

### Option 2: With Hot Reload (air)

```bash
# Terminal 1: Start services
docker compose up postgres redis

# Terminal 2: Run with hot reload
air
```

Benefits:
- Auto-rebuild on file changes
- Same environment as production
- Faster development cycle

### Option 3: Full Development Stack

```bash
make dev
```

This starts everything including Venio in Docker (less common for active development).

## Configuration

Configuration is loaded from `.env` file using Viper:

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=3690
APP_ENV=development
DEBUG=true

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=venio
POSTGRES_PASSWORD=<secure_password>
POSTGRES_DB=venio

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=<secure_password>

# JWT
JWT_SECRET=<at-least-32-characters>
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRY_DAYS=7
```

**Important:** JWT_SECRET must be at least 32 characters long for production.

## Architecture

### Package Structure

```
internal/
├── api/          # HTTP handlers and routes
├── services/     # Business logic
├── database/     # Database connection and utilities
├── models/       # Data structures
├── config/       # Configuration management
└── repositories/ # Data access layer (TBD)

cmd/
├── venio/        # Main application
└── worker/       # Background worker
```

### Database Access

We use **pgx v5** with a connection pool for PostgreSQL:

```go
// Typical database operation
db, err := database.Connect(ctx, &cfg.Database)
rows, err := db.Query(ctx, "SELECT * FROM users WHERE id = $1", userID)
```

Pattern:
- No ORM, direct SQL queries for control
- Connection pooling with configurable limits
- Proper error handling and context usage
- Repository pattern for data access

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
```

### 2. Make Changes

Edit code, following our [coding standards](../CONTRIBUTING.md#coding-standards).

### 3. Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./internal/services/...
```

### 4. Lint Code

```bash
make lint

# Auto-fix issues
make format
```

### 5. Commit Changes

```bash
git add .
git commit -m "feat: your feature description"
```

Pre-commit hooks will automatically:
- Format code
- Run linter
- Run tests

### 6. Push & Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Project Structure

```
venio/
├── cmd/
│   ├── venio/          # Main backend server
│   └── worker/         # Background worker
├── internal/           # Private application code
│   ├── api/           # HTTP handlers
│   ├── services/      # Business logic
│   ├── database/      # Database layer
│   ├── models/        # Data models
│   ├── config/        # Configuration
│   └── providers/     # External API clients
├── api/               # OpenAPI specs
├── migrations/        # Database migrations
├── configs/           # Config files
├── deployments/       # Docker & K8s
├── scripts/           # Build/dev scripts
└── docs/             # Documentation
```

## Common Tasks

### Database Migrations

```bash
# Create new migration
# TODO: Add migration tool

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

### Generate Mocks

```bash
# TODO: Add mockgen commands
```

### Generate OpenAPI Docs

```bash
# TODO: Add swag commands
```

### Build Binary

```bash
make build

# Binaries will be in ./bin/
./bin/venio
./bin/worker
```

### Build Docker Image

```bash
make docker-build
```

## Debugging

### VSCode Debugging

1. Set breakpoints in code
2. Press F5 or use Debug panel
3. Choose "Debug Venio" configuration

### Delve (CLI)

```bash
dlv debug cmd/venio/main.go
```

### Docker Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f venio
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Integration Tests

```bash
make test-integration
```

This spins up test databases in Docker.

### Test Coverage

```bash
make test-coverage
# Opens coverage.html in browser
```

## Code Quality

### Linting

```bash
# Run linter
golangci-lint run

# Auto-fix
golangci-lint run --fix
```

### Formatting

```bash
# Format all code
goimports -w .
gofmt -s -w .

# Or use Make
make format
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 3690
lsof -i :3690

# Kill it
kill -9 <PID>
```

### Docker Issues

```bash
# Clean up everything
docker compose down -v

# Rebuild from scratch
docker compose build --no-cache
```

### Go Module Issues

```bash
# Clear module cache
go clean -modcache

# Re-download modules
go mod download
```

## Performance Profiling

### CPU Profiling

```bash
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Memory Profiling

```bash
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

## Tips & Tricks

1. **Use Make targets** - They handle common tasks
2. **Enable pre-commit hooks** - Catches issues before commit
3. **Write tests first** - TDD helps design better APIs
4. **Use table-driven tests** - Easier to add test cases
5. **Check CI before pushing** - Save time

## Getting Help

- [GitHub Discussions](https://github.com/lusoris/venio/discussions)
- [Discord](#) (coming soon)
- [Contributing Guide](../CONTRIBUTING.md)

---

Happy coding! 🚀
