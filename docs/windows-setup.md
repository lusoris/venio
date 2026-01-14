# Windows Development Setup Guide

This guide provides step-by-step instructions for setting up a complete Venio development environment on Windows.

## Prerequisites

- **Windows 10/11** (21H2 or later recommended)
- **Administrator privileges** (for winget and package installation)
- **Internet connection**
- **PowerShell 5.1+** (or Windows Terminal)

## Automated Setup

### Option 1: Complete Automated Setup (Recommended)

We provide an automated PowerShell script that installs everything:

```powershell
# Run this from an Administrator PowerShell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
.\scripts\setup-windows-dev.ps1
```

This script will:
- ✅ Install winget (Windows Package Manager) if needed
- ✅ Install Go 1.25+
- ✅ Install Docker Desktop
- ✅ Install Git
- ✅ Install GNU Make
- ✅ Install development tools (Air, Delve, golangci-lint, etc.)
- ✅ Setup Lefthook pre-commit hooks
- ✅ Configure environment variables
- ✅ Download and run initial database migrations

### Option 2: Manual Installation

If you prefer manual installation or have issues with the script, follow these steps:

#### 1. Install Go

Using winget:
```powershell
winget install GoLang.Go -e
```

Or download from [go.dev](https://go.dev/dl/) and run the installer.

Verify installation:
```powershell
go version
```

#### 2. Install Docker Desktop

Using winget:
```powershell
winget install Docker.DockerDesktop -e
```

Or download from [Docker Desktop](https://www.docker.com/products/docker-desktop).

Start Docker Desktop and ensure the WSL 2 backend is configured.

Verify installation:
```powershell
docker --version
docker compose version
```

#### 3. Install Git

Using winget:
```powershell
winget install Git.Git -e
```

Verify:
```powershell
git --version
```

#### 4. Install GNU Make

Using winget:
```powershell
winget install GnuWin32.Make -e
```

Or using Chocolatey:
```powershell
choco install make
```

Verify:
```powershell
make --version
```

**Note:** After installing Make, restart PowerShell/Terminal for PATH updates to take effect.

#### 5. Install Development Tools

```powershell
# Hot reload (Air)
go install github.com/cosmtrek/air@latest

# Debugger (Delve)
go install github.com/go-delve/delve/cmd/dlv@latest

# Import formatter
go install golang.org/x/tools/cmd/goimports@latest

# Linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Pre-commit hooks
go install github.com/evilmartians/lefthook@latest
```

#### 6. Configure Environment

```powershell
cd venio
cp .env.example .env
# Edit .env with your settings (see Configuration section)
```

#### 7. Setup Lefthook

```powershell
lefthook install
```

## Configuration

Create a `.env` file from the template:

```powershell
cp .env.example .env
```

Edit `.env` with your settings:

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
POSTGRES_PASSWORD=YourSecurePassword123!
POSTGRES_DB=venio

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=YourRedisPassword123!

# JWT
JWT_SECRET=YourVerySecureJWTSecretAt32CharsMinimum12345
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRY_DAYS=7
```

## First Run

### 1. Start Docker Services

```powershell
cd c:\path\to\venio
docker compose up postgres redis -d
```

Verify containers are running:
```powershell
docker compose ps
```

Expected output:
```
NAME             IMAGE              STATUS
venio-postgres   postgres:16-alpine Up
venio-redis      redis:7-alpine     Up
```

### 2. Run Database Migrations

```powershell
# Using the build script
.\build.ps1 migrate-up

# Or manually with Go
go run cmd/venio/main.go  # Auto-runs migrations on startup
```

### 3. Start Backend Server

```powershell
# Option 1: Direct Go execution
go run cmd/venio/main.go

# Option 2: With build script (after Make is installed)
make dev

# Option 3: With hot reload (Air)
air
```

Expected output:
```
2026/01/14 14:45:56 🚀 Starting Venio Server...
2026/01/14 14:45:56 ✓ Database connection established
2026/01/14 14:45:56 ✅ Venio Server running on http://localhost:3690
[GIN-debug] Listening and serving HTTP on 0.0.0.0:3690
```

### 4. Start Frontend (in new PowerShell window)

```powershell
cd c:\path\to\venio\web
npm install  # First time only
npm run dev
```

Expected output:
```
▲ Next.js 16.1.1
- Local:         http://localhost:3000
✓ Ready in 562ms
```

### 5. Health Check

```powershell
# Test backend
Invoke-WebRequest -Uri "http://localhost:3690/health" | ConvertFrom-Json

# Test frontend
Start-Process "http://localhost:3000"
```

## Build Tools

### Using build.ps1 (PowerShell)

The `build.ps1` script provides all common development commands:

```powershell
# Show all available commands
.\build.ps1 help

# Common commands
.\build.ps1 dev              # Start all services
.\build.ps1 run              # Run backend only
.\build.ps1 watch            # Hot reload backend
.\build.ps1 test             # Run tests
.\build.ps1 lint             # Run linter
.\build.ps1 format           # Format code
.\build.ps1 build            # Build binary
.\build.ps1 migrate-up       # Run migrations
.\build.ps1 migrate-down     # Rollback migrations
.\build.ps1 db-reset         # Reset database
.\build.ps1 db-shell         # Connect to PostgreSQL
.\build.ps1 docker-up        # Start Docker services
.\build.ps1 docker-down      # Stop Docker services
```

### Using Makefile (After Make is installed)

Once GNU Make is installed, use the standard Makefile:

```powershell
make help
make dev
make test
make lint
make build
```

## Troubleshooting

### Issue: "make: command not found"

**Solution:** Install GNU Make and restart PowerShell:
```powershell
winget install GnuWin32.Make -e
# Restart PowerShell
```

### Issue: "docker: command not found"

**Solution:** Ensure Docker Desktop is installed and running:
1. Install Docker Desktop: `winget install Docker.DockerDesktop -e`
2. Open Docker Desktop application
3. Wait for initialization (may take 30 seconds)
4. Restart PowerShell

### Issue: "Port 5432 already in use"

**Solution:** Kill existing PostgreSQL:
```powershell
# Find and stop the container
docker compose down
docker ps -a  # Verify it's stopped
docker compose up postgres -d  # Restart
```

### Issue: "Database connection failed"

**Solution:** Check credentials in `.env`:
```powershell
# Verify container is running
docker compose ps

# Check logs
docker compose logs postgres

# Test connection
docker exec venio-postgres psql -U venio -d venio -c "SELECT 1"
```

### Issue: Frontend not connecting to backend

**Solution:** Ensure backend is running and check `.env.local`:

```powershell
# Backend should be running on 3690
curl http://localhost:3690/health

# Check web/.env.local
cat web\.env.local
# Should contain: NEXT_PUBLIC_API_URL=http://localhost:3690/api/v1
```

### Issue: Go modules not found

**Solution:** Download Go modules:
```powershell
go mod download
go mod tidy
```

### Issue: "lefthook: command not found"

**Solution:** Ensure `$HOME\go\bin` is in PATH:
```powershell
# Check if Go bin is in PATH
$env:PATH -split ";" | grep -i "go\\bin"

# If not found, add it manually
$env:PATH += ";$HOME\go\bin"
[System.Environment]::SetEnvironmentVariable("PATH", $env:PATH, [System.EnvironmentVariableTarget]::User)
```

## VSCode Setup

### Install Extensions

Option 1: Manual installation via the setup script:
```powershell
.\scripts\install-vscode-extensions.ps1
```

Option 2: Manual installation from VSCode command palette:
- Press `Ctrl+Shift+X`
- Search for and install:
  - `golang.go` - Go support
  - `ms-azuretools.vscode-docker` - Docker support
  - `eamodio.gitlens` - Git integration
  - `mtxr.sqltools` - Database tools
  - `redhat.vscode-yaml` - YAML/Config support

### Debug Configuration

Create `.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to Delve",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "remotePath": "",
      "port": 2345,
      "host": "127.0.0.1",
      "showLog": true
    }
  ]
}
```

### Launch Backend with Debugger

```powershell
dlv debug ./cmd/venio --listen=:2345 --headless --api-version=2
```

Then click "Run and Debug" in VSCode.

## Environment Verification

Run this script to verify your setup:

```powershell
Write-Host "=== Venio Development Environment Check ===" -ForegroundColor Cyan

$checks = @{
  "Go" = "go version"
  "Docker" = "docker --version"
  "Git" = "git --version"
  "Make" = "make --version"
  "Air" = "air --version"
  "Delve" = "dlv version"
  "golangci-lint" = "golangci-lint --version"
}

foreach ($check in $checks.GetEnumerator()) {
  Write-Host "Checking $($check.Name)..." -NoNewline
  
  try {
    $output = Invoke-Expression $check.Value 2>&1
    if ($LASTEXITCODE -eq 0) {
      Write-Host " ✓" -ForegroundColor Green
    } else {
      Write-Host " ✗" -ForegroundColor Red
    }
  } catch {
    Write-Host " ✗" -ForegroundColor Red
  }
}

Write-Host "`nChecking Docker services..." -ForegroundColor Cyan
docker compose ps

Write-Host "`nEnvironment check complete!" -ForegroundColor Green
```

## Next Steps

1. Read the [Development Guide](development.md) for detailed workflow
2. Check [Architecture Overview](architecture.md) to understand the codebase
3. Review [API Documentation](api.md) for endpoint details
4. See [Contributing Guide](../CONTRIBUTING.md) for coding standards

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Gin Web Framework](https://gin-gonic.com/)

## Support

If you encounter issues:

1. Check the [Troubleshooting](#troubleshooting) section above
2. Review GitHub Issues: https://github.com/lusoris/venio/issues
3. Check Development Guide: [development.md](development.md)
4. Ask in discussions: https://github.com/lusoris/venio/discussions
