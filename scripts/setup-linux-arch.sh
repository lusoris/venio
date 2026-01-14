#!/bin/bash
# Setup Linux Development Environment for Venio (Arch Linux)
# This script automates the installation of all required tools

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

function print_success() {
    echo -e "${GREEN}✓ $*${NC}"
}

function print_error() {
    echo -e "${RED}✗ $*${NC}"
}

function print_warning() {
    echo -e "${YELLOW}⚠ $*${NC}"
}

function print_info() {
    echo -e "${CYAN}ℹ $*${NC}"
}

function print_header() {
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}$*${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Check if running with sudo (for pacman)
if [[ $EUID -ne 0 ]]; then
    print_error "This script must be run as root (use: sudo bash scripts/setup-linux-arch.sh)"
    exit 1
fi

echo ""
print_info "╔═══════════════════════════════════════════════════════════════╗"
print_info "║    Venio Linux Development Environment Setup (Arch Linux)     ║"
print_info "║    This script will install all required development tools    ║"
print_info "╚═══════════════════════════════════════════════════════════════╝"
echo ""

# Update package lists
print_header "Updating Package Lists"
pacman -Syy --noconfirm &>/dev/null
print_success "Package lists updated"
echo ""

# Define packages to install with pacman
print_header "Installing System Packages"

packages=(
    "base-devel"            # Compiler, make, git, etc.
    "curl"                  # HTTP client
    "wget"                  # Download utility
    "gnu-netcat"            # netcat implementation
    "sqlite"                # SQLite database
    "jq"                    # JSON processor
    "postgresql-libs"       # PostgreSQL client libraries
    "redis"                 # Redis (can also use container)
)

installed=0
skipped=0
failed=0

for package in "${packages[@]}"; do
    if pacman -Q "$package" &>/dev/null 2>&1; then
        print_info "  $package (already installed)"
        ((skipped++))
    else
        echo -n "  Installing $package... "
        if pacman -S "$package" --noconfirm &>/dev/null 2>&1; then
            print_success ""
            ((installed++))
        else
            print_error ""
            ((failed++))
        fi
    fi
done

echo ""
print_success "System packages: $installed installed, $skipped already present"
[[ $failed -gt 0 ]] && print_warning "  $failed packages failed to install"
echo ""

# Install Go 1.25 from AUR or official repo
print_header "Installing Go 1.25"

if command -v go &> /dev/null; then
    CURRENT_GO=$(go version | awk '{print $3}')
    print_success "Go already installed: $CURRENT_GO"
else
    echo -n "Installing Go from Arch repository... "

    if pacman -S go --noconfirm &>/dev/null 2>&1; then
        print_success ""
        print_info "Go installed"
        export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    else
        print_error ""
        print_warning "Please install Go manually or from AUR"
        ((failed++))
    fi
fi

echo ""

# Install Node.js (LTS)
print_header "Installing Node.js (LTS)"

if command -v node &> /dev/null; then
    NODE_VERSION=$(node -v)
    print_success "Node.js already installed: $NODE_VERSION"
else
    echo -n "Installing Node.js from Arch repository... "

    if pacman -S nodejs npm --noconfirm &>/dev/null 2>&1; then
        print_success ""
        NODE_VERSION=$(node -v)
        print_info "Node.js $NODE_VERSION installed"

        # Update npm
        npm install -g npm@latest &>/dev/null
    else
        print_error ""
        ((failed++))
    fi
fi

echo ""

# Install Docker (or docker with systemd)
print_header "Installing Docker"

if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version)
    print_success "Docker already installed: $DOCKER_VERSION"
else
    echo -n "Installing Docker from Arch repository... "

    if pacman -S docker docker-compose --noconfirm &>/dev/null 2>&1; then
        print_success ""

        # Enable docker daemon
        systemctl enable docker
        systemctl start docker

        # Add current user to docker group
        if [[ -n "$SUDO_USER" ]]; then
            usermod -aG docker "$SUDO_USER"
            print_info "Added $SUDO_USER to docker group (logout/login required)"
        fi
    else
        print_error ""
        print_warning "Please install Docker manually: pacman -S docker"
        ((failed++))
    fi
fi

echo ""

# Install PostgreSQL client and tools
print_header "Installing Database Tools"

if command -v psql &> /dev/null; then
    PG_VERSION=$(psql --version)
    print_success "PostgreSQL client installed: $PG_VERSION"
else
    echo -n "Installing PostgreSQL client... "
    if pacman -S postgresql-libs --noconfirm &>/dev/null 2>&1; then
        print_success ""
    else
        print_error ""
        ((failed++))
    fi
fi

echo ""

# Install Go development tools
print_header "Installing Go Development Tools"

export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

goTools=(
    "Air (hot reload)|github.com/cosmtrek/air@latest"
    "Delve (debugger)|github.com/go-delve/delve/cmd/dlv@latest"
    "goimports (formatter)|golang.org/x/tools/cmd/goimports@latest"
    "golangci-lint (linter)|github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    "Lefthook (git hooks)|github.com/evilmartians/lefthook@latest"
)

toolsInstalled=0
toolsFailed=0

for tool in "${goTools[@]}"; do
    IFS='|' read -r name package <<< "$tool"
    echo -n "Installing $name... "

    if go install "$package" &>/dev/null; then
        print_success ""
        ((toolsInstalled++))
    else
        print_error ""
        print_warning "  Try manually: go install $package"
        ((toolsFailed++))
    fi
done

echo ""
print_success "Go tools: $toolsInstalled installed"
[[ $toolsFailed -gt 0 ]] && print_warning "  $toolsFailed Go tools failed"
echo ""

# Install npm global tools
print_header "Installing Node.js Development Tools"

npmTools=(
    "snyk"
)

npmInstalled=0
npmFailed=0

for tool in "${npmTools[@]}"; do
    echo -n "Installing $tool... "

    if npm install -g "$tool" &>/dev/null; then
        print_success ""
        ((npmInstalled++))
    else
        print_error ""
        print_warning "  Try manually: npm install -g $tool"
        ((npmFailed++))
    fi
done

echo ""
print_success "npm tools: $npmInstalled installed"
[[ $npmFailed -gt 0 ]] && print_warning "  $npmFailed npm tools failed"
echo ""

# Final checks
print_header "Environment Verification"

checks=(
    "Go|go version"
    "Docker|docker --version"
    "Docker Compose|docker compose version"
    "Git|git --version"
    "Node.js|node --version"
    "npm|npm --version"
)

verified=0
unverified=0

for check in "${checks[@]}"; do
    IFS='|' read -r name command <<< "$check"
    echo -n "Checking $name... "

    if output=$($command 2>/dev/null); then
        print_success ""
        print_info "  $output"
        ((verified++))
    else
        print_error ""
        ((unverified++))
    fi
done

echo ""
print_success "Verification: $verified working, $unverified need attention"
echo ""

# Setup instructions
print_header "Next Steps"
cat <<EOF
1. RELOAD YOUR SHELL (if PATH was updated):
   source ~/.bashrc  # or ~/.zshrc / ~/.config/fish/config.fish
   Or open a new terminal

2. Clone the repository:
   git clone https://github.com/lusoris/venio.git
   cd venio

3. Copy environment template:
   cp .env.example .env

4. Edit .env with your settings (passwords, secrets, etc.)

5. Start Docker services:
   docker compose up postgres redis -d

6. Run migrations and seed data:
   make migrate-up
   make seed-data

7. Run the application:
   go run cmd/venio/main.go

8. In a new terminal, start the frontend:
   cd web
   npm install
   npm run dev

9. Access the application:
   Backend: http://localhost:3690
   Frontend: http://localhost:3000

For more details, see: scripts/SETUP_README.md
EOF

echo ""
print_success "╔═══════════════════════════════════════════════════════════════╗"
print_success "║       Setup Complete! Reload your shell to continue.         ║"
print_success "╚═══════════════════════════════════════════════════════════════╝"
echo ""
