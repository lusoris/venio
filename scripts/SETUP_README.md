# Complete Development Environment Setup

Automated setup scripts for all major platforms. Choose your OS below to get started.

## Quick Start by Platform

### Windows
```powershell
# 1. Open PowerShell as Administrator
# 2. Run the setup script:
.\scripts\setup-windows-dev.ps1

# 3. Restart your terminal
# 4. Clone and run:
git clone https://github.com/lusoris/venio.git
cd venio
docker compose up postgres redis -d
go run cmd/venio/main.go
```

### macOS
```bash
# 1. Install Homebrew (if not already installed):
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Run the setup script:
bash scripts/setup-macos-dev.sh

# 3. Start Docker Desktop from Applications/Docker.app

# 4. Reload your shell:
source ~/.zshrc  # or ~/.bashrc / ~/.config/fish/config.fish

# 5. Clone and run:
git clone https://github.com/lusoris/venio.git
cd venio
docker compose up postgres redis -d
go run cmd/venio/main.go
```

### Linux - Debian/Ubuntu
```bash
# 1. Run with sudo:
sudo bash scripts/setup-linux-debian.sh

# 2. Reload your shell:
source /etc/profile.d/venio.sh

# 3. Clone and run:
git clone https://github.com/lusoris/venio.git
cd venio
docker compose up postgres redis -d
go run cmd/venio/main.go
```

### Linux - Arch/Manjaro
```bash
# 1. Run with sudo:
sudo bash scripts/setup-linux-arch.sh

# 2. Reload your shell and add to docker group:
newgrp docker
source ~/.bashrc  # or ~/.zshrc / ~/.config/fish/config.fish

# 3. Clone and run:
git clone https://github.com/lusoris/venio.git
cd venio
docker compose up postgres redis -d
go run cmd/venio/main.go
```

### Linux - Fedora/RHEL/CentOS
```bash
# 1. Run with sudo:
sudo bash scripts/setup-fedora-dev.sh

# 2. Apply group changes:
newgrp docker

# 3. Reload your shell:
source /etc/profile.d/venio.sh

# 4. Clone and run:
git clone https://github.com/lusoris/venio.git
cd venio
docker compose up postgres redis -d
go run cmd/venio/main.go
```

## What Gets Installed

All setup scripts install:

### Core Tools
- **Go 1.25** - Backend language
- **Node.js LTS + npm** - Frontend toolchain
- **Git** - Version control
- **Docker + Docker Compose** - Containerization
- **PostgreSQL client tools** - Database connectivity
- **Build tools** - Compilers and build utilities

### Go Development Tools
- **Air** - Hot reload for development
- **golangci-lint** - Code linting
- **goimports** - Code formatting
- **Delve** - Go debugger
- **Lefthook** - Git hooks manager

### Node.js Global Tools
- **npm** (latest) - Package manager
- **snyk** - Security scanning

### Database Tools
- PostgreSQL client (psql)
- Redis CLI tools
- SQLite 3

## System Requirements

### Windows
- Windows 10 or 11
- 8GB RAM minimum (16GB recommended)
- 20GB free disk space
- Administrator access
- PowerShell 5.0+

### Linux - Debian/Ubuntu
- Ubuntu 22.04 LTS or later, or Debian 12+
- 8GB RAM minimum (16GB recommended)
- 20GB free disk space
- sudo access
- Supported architectures: x86_64, ARM64

### Linux - Arch
- Arch Linux or Manjaro
- 8GB RAM minimum (16GB recommended)
- 20GB free disk space
- sudo access

### Linux - Fedora/RHEL/CentOS
- Fedora 39+ (recommended) or RHEL 9+ or CentOS Stream 9+
- 8GB RAM minimum (16GB recommended)
- 20GB free disk space
- sudo access
- Supported architectures: x86_64, ARM64

### macOS
- macOS 11.0 or later (12.0+ recommended for M1/M2 support)
- 8GB RAM minimum (16GB recommended)
- 20GB free disk space
- Homebrew installed
- Intel or Apple Silicon (M1/M2/M3)
- Supported architectures: x86_64, ARM64

## Manual Installation

If scripts fail, you can install components manually:

### Install Go 1.25
```bash
# Download latest from https://go.dev/dl/
# Or use your package manager:

# Debian/Ubuntu
sudo apt install golang-go

# Arch
sudo pacman -S go
```

### Install Node.js
```bash
# Debian/Ubuntu
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# Arch
sudo pacman -S nodejs npm
```

### Install Docker
```bash
# Debian/Ubuntu
curl -fsSL https://get.docker.com -o get-docker.sh
sudo bash get-docker.sh

# Arch
sudo pacman -S docker docker-compose
```

### Install Go Tools
```bash
go install github.com/cosmtrek/air@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/evilmartians/lefthook@latest
```

## Post-Installation

### 1. Verify Installation
```bash
go version
node --version
docker --version
docker compose version
git --version
```

### 2. Clone Repository
```bash
git clone https://github.com/lusoris/venio.git
cd venio
```

### 3. Create Environment File
```bash
cp .env.example .env
# Edit .env with your settings
```

### 4. Start Database Services
```bash
docker compose up postgres redis -d
```

### 5. Run Migrations
```bash
make migrate-up
```

### 6. Seed Test Data
```bash
make seed-data
```

### 7. Run Backend
```bash
# Terminal 1: Backend
go run cmd/venio/main.go

# Terminal 2: Frontend
cd web
npm install
npm run dev
```

### 8. Access Application
- **Backend API:** http://localhost:3690
- **Frontend:** http://localhost:3000
- **API Docs:** http://localhost:3690/docs (if available)

## Troubleshooting

### "Command not found" after installation

**Solution:** Reload your shell to update PATH environment:

```bash
# Windows (PowerShell)
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

# Linux (Bash)
source ~/.bashrc

# Linux (Fish)
source ~/.config/fish/config.fish

# Or simply open a new terminal
```

### Docker daemon not running

**Solution:** Start Docker service:

```bash
# Windows: Docker Desktop should auto-start. If not:
docker run hello-world

# Linux
sudo systemctl start docker
sudo systemctl enable docker
```

### Permission denied on scripts

**Solution:** Make scripts executable:

```bash
chmod +x scripts/setup-*.sh
chmod +x scripts/test-*.sh
```

### Go tools installation fails

**Solution:** Ensure `$HOME/go/bin` is in PATH:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:$HOME/go/bin

# Then install again:
go install github.com/cosmtrek/air@latest
```

### Docker permission errors

**Solution:** Add your user to docker group:

```bash
# Linux
sudo usermod -aG docker $USER
newgrp docker

# Or use sudo
sudo docker compose up
```

### Port already in use

**Solution:** Change ports in `docker-compose.yml` or stop conflicting services:

```bash
# Find process on port 5432 (PostgreSQL)
lsof -i :5432
sudo kill -9 <PID>

# Or change ports in docker-compose.yml:
services:
  postgres:
    ports:
      - "5433:5432"  # Changed from 5432
```

## Development Workflow

### Start Development
```bash
# Terminal 1: Backend (with hot reload)
make watch

# Terminal 2: Frontend
cd web && npm run dev

# Terminal 3: Database management
make db-shell  # Connect to PostgreSQL
```

### Run Tests
```bash
# Unit tests
make test

# With coverage
make test-coverage

# API integration tests
make test-api

# Integration tests
make test-integration
```

### Seed Data
```bash
# Create test users with different roles
make seed-data

# Test credentials:
# admin@test.local / AdminPassword123!
# user@test.local / UserPassword123!
# moderator@test.local / ModeratorPassword123!
# guest@test.local / GuestPassword123!
```

### Code Quality
```bash
# Format code
make format

# Lint code
make lint

# Security scan
snyk test
```

## Cross-Platform Development

When syncing code between Windows and Linux:

### Configure Git Line Endings

```bash
# Windows
git config --global core.autocrlf true

# Linux/macOS
git config --global core.autocrlf input
```

### Use Platform-Specific Scripts

- **Windows:** Use PowerShell scripts (`.ps1`)
- **Linux/macOS:** Use Bash scripts (`.sh`) or Fish (`.fish`)

### Keep Scripts Executable

In `.gitattributes`:
```
scripts/*.ps1 text eol=crlf
scripts/*.sh text eol=lf
scripts/*.fish text eol=lf
```

## Continuous Integration

### GitHub Actions Setup
```yaml
name: CI
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      - uses: actions/setup-node@v3
        with:
          node-version: 'lts/*'
      - run: make test
      - run: make lint
```

## Environment Variables

### Required (.env)
```env
# Server
SERVER_PORT=3690
SERVER_ENV=development

# Database
DATABASE_URL=postgres://venio:venio@localhost:5432/venio?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h

# Frontend (in web/.env.local)
VITE_API_URL=http://localhost:3690
```

## Getting Help

- **Documentation:** See `docs/` directory
- **Issues:** GitHub Issues
- **Discussions:** GitHub Discussions
- **Email:** Support contact in README

## Updating Dependencies

### Go Dependencies
```bash
go get -u ./...
go mod tidy
```

### Node.js Dependencies
```bash
cd web
npm update
npm audit fix
```

### Docker Images
```bash
docker compose pull
docker compose up -d
```

## Platform-Specific Notes

### Windows
- Use `\` or `/` in paths (PowerShell handles both)
- CRLF line endings (set by `core.autocrlf=true`)
- Docker Desktop required (WSL2 backend recommended)
- PowerShell execution policy may need adjustment:
  ```powershell
  Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
  ```

### Linux - Debian/Ubuntu
- LF line endings
- systemd for service management
- PostgreSQL 18.1+ recommended
- Ensure `universe` repo is enabled for some packages

### Linux - Arch
- LF line endings
- systemd for service management
- Rolling release (always latest stable versions)
- May need to rebuild some AUR packages

## Next Steps

1. Run the setup script for your platform
2. Verify installation: `go version && node --version && docker --version`
3. Clone the repository
4. Follow "Post-Installation" steps above
5. Read [CONTRIBUTING.md](../CONTRIBUTING.md) for development guidelines

---

**Last Updated:** January 2026
**Supported Platforms:** Windows 10+, Ubuntu 22.04+, Debian 12+, Arch Linux
**Go Version:** 1.25 (latest stable)
**Node.js Version:** LTS (currently 22.x)
