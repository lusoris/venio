# Windows Development Setup Instructions

## Status: ✗ Missing Prerequisites

The setup script has been created and tested. It correctly identifies that the following tools are **NOT installed**:

### Missing Tools
- **Go 1.23+** ❌
- **Docker Desktop** ❌

## Installation Instructions

### 1. Install Go 1.23+

1. Visit: https://go.dev/dl/
2. Download: `go1.23.x.windows-amd64.msi` (latest 1.23.x version)
3. Run the installer (double-click the .msi file)
4. Follow the installation wizard
5. **Important**: Restart your terminal/PowerShell after installation
6. Verify installation by running: `go version`

### 2. Install Docker Desktop

1. Visit: https://www.docker.com/products/docker-desktop
2. Download: Docker Desktop for Windows
3. Run the installer
4. Follow the installation wizard
5. **Important**: Allow WSL 2 integration (recommended)
6. **Important**: Restart your computer after installation
7. Verify installation by running: `docker --version`

## Running the Setup Script

Once both Go and Docker are installed:

```powershell
# Navigate to the project directory
cd c:\Users\ms\dev\venio

# Run the setup script
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force
.\scripts\setup-windows.ps1
```

### What the Script Does

1. ✓ Verifies Go installation
2. ✓ Verifies Docker & Docker Compose installation
3. ✓ Installs Go development tools:
   - Air (hot reload)
   - Delve (debugger)
   - goimports (import formatter)
   - golangci-lint (linter)
   - Lefthook (git hooks)
4. ✓ Sets up Git hooks
5. ✓ Downloads Go dependencies
6. ✓ Creates .env file from .env.example
7. ✓ Verifies all installations

## After Successful Setup

You'll have three options to run the development environment:

### Option 1: Full Development Stack (Docker)
```powershell
docker compose -f docker-compose.yml -f docker-compose.dev.yml up
```
Access at: http://localhost:3690

### Option 2: Run Locally (Services in Docker)
```powershell
# Terminal 1: Start services
docker compose up postgres redis typesense

# Terminal 2: Run Venio
go run cmd/venio/main.go
```

### Option 3: Run with Hot Reload
```powershell
air
```

## Next Steps

1. Install Go from: https://go.dev/dl/
2. Install Docker Desktop from: https://www.docker.com/products/docker-desktop
3. Restart your terminal
4. Run the setup script: `.\scripts\setup-windows.ps1`
5. Choose one of the three options above to start development
