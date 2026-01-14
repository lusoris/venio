# Venio Scripts

Utility scripts for development, testing, deployment, and environment setup.

## Script Directory

### Environment Setup

These scripts automate the installation of all required development tools.

| Script | Platform | Purpose |
|--------|----------|---------|
| [setup-windows-dev.ps1](setup-windows-dev.ps1) | Windows | Complete dev environment setup with winget package manager |
| [setup-macos-dev.sh](setup-macos-dev.sh) | macOS | Complete dev environment setup with Homebrew |
| [setup-linux-debian.sh](setup-linux-debian.sh) | Debian/Ubuntu | Complete dev environment setup with apt package manager |
| [setup-linux-arch.sh](setup-linux-arch.sh) | Arch/Manjaro | Complete dev environment setup with pacman package manager |
| [setup-fedora-dev.sh](setup-fedora-dev.sh) | Fedora/RHEL/CentOS | Complete dev environment setup with dnf/yum package manager |
| [SETUP_README.md](SETUP_README.md) | All Platforms | Comprehensive setup guide with troubleshooting |

**Quick Start:**
```bash
# Windows (PowerShell as Admin)
.\scripts\setup-windows-dev.ps1

# macOS (with brew)
bash scripts/setup-macos-dev.sh

# Debian/Ubuntu (with sudo)
sudo bash scripts/setup-linux-debian.sh

# Arch Linux (with sudo)
sudo bash scripts/setup-linux-arch.sh

# Fedora/RHEL/CentOS (with sudo)
sudo bash scripts/setup-fedora-dev.sh
```

### Testing

These scripts automate API testing with comprehensive RBAC endpoint coverage.

| Script | Shell | Purpose |
|--------|-------|---------|
| [test-api.ps1](test-api.ps1) | PowerShell | RBAC API test suite for Windows |
| [test-api.sh](test-api.sh) | Bash | RBAC API test suite for Linux/macOS/Git Bash |
| [test-api.fish](test-api.fish) | Fish | RBAC API test suite for Fish shell users |
| [TEST_API_README.md](TEST_API_README.md) | All Shells | Detailed test documentation and examples |

**Quick Start:**
```bash
# PowerShell
.\scripts\test-api.ps1 -BaseURL http://localhost:8080 -Verbose 1

# Bash
./scripts/test-api.sh http://localhost:8080 1

# Fish
./scripts/test-api.fish http://localhost:8080 1

# Or via Makefile
make test-api
```

### VSCode Extensions

| Script | Purpose |
|--------|---------|
| [install-vscode-extensions.ps1](install-vscode-extensions.ps1) | Auto-install recommended VSCode extensions |

**Usage:**
```powershell
.\scripts\install-vscode-extensions.ps1
```

## What Each Setup Script Installs

### Go Development Tools
- **Go 1.25** - Latest stable Go compiler
- **Air** - Hot reload for development
- **golangci-lint** - Advanced Go linter
- **goimports** - Automatic Go import formatting
- **Delve** - Go debugger
- **Lefthook** - Git hooks manager

### Node.js Development
- **Node.js LTS** - JavaScript runtime (latest LTS)
- **npm** - Package manager (auto-updated to latest)
- **snyk** - Security vulnerability scanning

### Infrastructure Tools
- **Docker** - Container runtime
- **Docker Compose** - Multi-container orchestration
- **PostgreSQL Client** - psql CLI for database management
- **Redis Tools** - redis-cli for Redis management
- **Git** - Version control
- **Make** - Build automation (Windows only)

### Build Tools
- **C/C++ Compiler** - For native module compilation
- **Build Essentials** - Linux build toolchain
- **curl/wget** - HTTP clients for downloads

## Test Scripts Features

### Test Phases
All test scripts execute the same 7 test phases:

1. **Authentication** - Login as all 4 test users, generate JWT tokens
2. **Role Management** - CRUD operations on roles (admin only)
3. **Permission Management** - CRUD operations on permissions (admin only)
4. **Role-Permission Assignment** - Link permissions to roles
5. **User Role Management** - Assign/revoke roles from users
6. **Permission-Based Access Control** - Verify authorization enforcement
7. **Cleanup** - Remove test data after testing

### Test Users
All three test scripts use identical credentials:

```
admin@test.local / AdminPassword123! (admin role)
moderator@test.local / ModeratorPassword123! (moderator role)
user@test.local / UserPassword123! (user role)
guest@test.local / GuestPassword123! (guest role)
```

### Test Coverage
- **45+ API endpoint tests**
- **Authorization checks** for each role
- **CRUD operations** (Create, Read, Update, Delete)
- **Permission enforcement**
- **Cross-role permission validation**
- **Automatic cleanup** of test data

## Setup Script Compatibility

| Feature | Windows | Debian/Ubuntu | Arch Linux |
|---------|---------|---------------|-----------|
| Go 1.25 | ✓ | ✓ | ✓ |
| Node.js LTS | ✓ | ✓ | ✓ |
| Docker | ✓ | ✓ | ✓ |
| PostgreSQL Tools | ✓ | ✓ | ✓ |
| Go Tools | ✓ | ✓ | ✓ |
| npm Tools | ✓ | ✓ | ✓ |
| Path Configuration | ✓ | ✓ | ✓ |
| Auto Group Setup | ✓ | ✓ | ✓ |
| Package Detection | ✓ | ✓ | ✓ |

## Manual Tool Installation

If any tool fails to install automatically, install manually:

### Go
```bash
# https://golang.org/dl/ - Download Go 1.25
# Windows: Use installer or winget install golang.Go
# Linux: pacman -S go (Arch) or apt install golang-go (Debian)
```

### Node.js
```bash
# Windows: winget install OpenJS.NodeJS
# Debian: curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo bash
# Arch: pacman -S nodejs npm
```

### Docker
```bash
# Windows: winget install docker.docker
# Debian: curl -fsSL https://get.docker.com | sh
# Arch: pacman -S docker
```

### Go Tools
```bash
# After Go is installed:
go install github.com/cosmtrek/air@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/evilmartians/lefthook@latest
```

### npm Global Packages
```bash
npm install -g snyk
```

## Troubleshooting

### "Permission denied" on Linux scripts
```bash
chmod +x scripts/setup-linux-*.sh
chmod +x scripts/test-*.sh
```

### "Command not found" after installation
```bash
# Reload your shell's configuration:

# Bash
source ~/.bashrc

# Zsh
source ~/.zshrc

# Fish
source ~/.config/fish/config.fish

# Or open a new terminal window
```

### Docker permission errors
```bash
# Add your user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Then logout and login, or:
sudo docker compose up
```

### Port conflicts
```bash
# Check what's using port 5432 (PostgreSQL)
lsof -i :5432

# Change ports in docker-compose.yml if needed
# Or kill the process:
sudo kill -9 <PID>
```

## Development Commands via Makefile

Once setup is complete, use these commands:

```bash
# Setup
make install          # Install tools and dependencies
make setup            # Full setup: install + docker + migrate

# Development
make dev              # Start full Docker environment
make run              # Run app locally without Docker
make watch            # Run with hot reload (Air)

# Testing
make test             # Run unit tests
make test-api         # Run API integration tests
make test-coverage    # Run tests with coverage report

# Database
make db-shell         # Connect to PostgreSQL CLI
make migrate-up       # Run database migrations
make seed-data        # Seed test users and roles

# Code Quality
make format           # Format Go code
make lint             # Lint Go code
make build            # Build binaries
```

## Cross-Platform Development Tips

When syncing code between Windows and Linux:

1. **Configure Git line endings:**
   ```bash
   git config --global core.autocrlf true   # Windows
   git config --global core.autocrlf input  # Linux
   ```

2. **Use platform-specific scripts:**
   - Windows → PowerShell (`.ps1`)
   - Linux → Bash (`.sh`) or Fish (`.fish`)

3. **Keep scripts executable:**
   - Add to `.gitattributes`:
     ```
     scripts/*.ps1 text eol=crlf
     scripts/*.sh text eol=lf
     scripts/*.fish text eol=lf
     ```

## Contributing

When adding new scripts:

1. Create platform-specific versions (Windows, Linux)
2. Add comprehensive error handling
3. Include color-coded output
4. Document in this README
5. Test on target platform
6. Add to appropriate section above

## Related Documentation

- [SETUP_README.md](SETUP_README.md) - Detailed setup guide
- [TEST_API_README.md](TEST_API_README.md) - API testing documentation
- [../docs/development.md](../docs/development.md) - Development guidelines
- [../CONTRIBUTING.md](../CONTRIBUTING.md) - Contribution guidelines

---

**Last Updated:** January 2026
**Platforms Supported:** Windows 10+, Ubuntu 22.04+, Debian 12+, Arch Linux
**Go Version:** 1.25 (latest stable)
**Node.js:** LTS (currently 22.x)
