#!/bin/bash

################################################################################
# Venio Development Environment Setup for macOS
#
# This script sets up a complete development environment for Venio on macOS
# using Homebrew (brew) as the package manager.
#
# Usage:
#   bash scripts/setup-macos-dev.sh
#
# Requirements:
#   - macOS 11.0 or later
#   - Homebrew installed (https://brew.sh)
#   - Internet connection
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
INSTALLED=0
SKIPPED=0
FAILED=0

################################################################################
# Helper Functions
################################################################################

print_header() {
    echo -e "\n${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}${1}${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
    ((INSTALLED++))
}

print_skip() {
    echo -e "${YELLOW}⊘${NC} $1"
    ((SKIPPED++))
}

print_error() {
    echo -e "${RED}✗${NC} $1"
    ((FAILED++))
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

check_command() {
    if command -v "$1" &> /dev/null; then
        return 0
    else
        return 1
    fi
}

install_brew_package() {
    local package=$1
    local name=${2:-$package}

    if check_command "$package"; then
        print_skip "$name already installed"
    else
        echo "Installing $name..."
        if brew install "$package" &> /dev/null; then
            print_success "$name installed"
        else
            print_error "Failed to install $name"
        fi
    fi
}

install_brew_cask() {
    local package=$1
    local name=${2:-$package}

    if brew list --cask "$package" &> /dev/null 2>&1; then
        print_skip "$name already installed"
    else
        echo "Installing $name..."
        if brew install --cask "$package" &> /dev/null; then
            print_success "$name installed"
        else
            print_error "Failed to install $name"
        fi
    fi
}

################################################################################
# System Checks
################################################################################

print_header "System Requirements Check"

# Check macOS version
OS_VERSION=$(sw_vers -productVersion)
echo "macOS version: $OS_VERSION"

# Check if Homebrew is installed
if ! check_command "brew"; then
    print_error "Homebrew is not installed"
    echo -e "\nPlease install Homebrew first:"
    echo "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    exit 1
fi

print_success "Homebrew is installed"
HOMEBREW_VERSION=$(brew --version | head -n 1)
echo "  $HOMEBREW_VERSION"

# Update Homebrew
print_info "Updating Homebrew..."
if brew update &> /dev/null; then
    print_success "Homebrew updated"
else
    print_error "Failed to update Homebrew"
fi

################################################################################
# Core Development Tools
################################################################################

print_header "Installing Core Development Tools"

# Git
install_brew_package "git" "Git"

# Make
install_brew_package "make" "Make"

# Curl & Wget
install_brew_package "curl" "Curl"
install_brew_package "wget" "Wget"

# jq (JSON processor)
install_brew_package "jq" "jq (JSON processor)"

# SQLite
install_brew_package "sqlite" "SQLite"

################################################################################
# Go Development
################################################################################

print_header "Installing Go Development Stack"

# Go 1.25
echo "Installing Go 1.25..."
if check_command "go"; then
    CURRENT_VERSION=$(go version | awk '{print $3}')
    echo "Current Go version: $CURRENT_VERSION"

    if [[ "$CURRENT_VERSION" == "go1.25"* ]]; then
        print_skip "Go 1.25 already installed"
    else
        print_info "Upgrading to Go 1.25..."
        brew upgrade go &> /dev/null
        print_success "Go upgraded to $(go version | awk '{print $3}')"
    fi
else
    if brew install go &> /dev/null; then
        print_success "Go 1.25 installed"
        echo "  $(go version)"
    else
        print_error "Failed to install Go"
    fi
fi

# Go development tools
print_info "\nInstalling Go development tools..."

GO_TOOLS=(
    "github.com/cosmtrek/air@latest:Air (hot reload)"
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest:golangci-lint (linter)"
    "golang.org/x/tools/cmd/goimports@latest:goimports (imports formatter)"
    "github.com/go-delve/delve/cmd/dlv@latest:Delve (debugger)"
    "github.com/evilmartians/lefthook@latest:Lefthook (git hooks)"
)

for tool in "${GO_TOOLS[@]}"; do
    IFS=':' read -r pkg desc <<< "$tool"
    echo "Installing $desc..."
    if go install "$pkg" &> /dev/null; then
        print_success "$desc installed"
    else
        print_error "Failed to install $desc"
    fi
done

################################################################################
# Node.js Development
################################################################################

print_header "Installing Node.js Development Stack"

# Node.js LTS via Homebrew
echo "Installing Node.js LTS..."
if check_command "node"; then
    NODE_VERSION=$(node --version)
    print_skip "Node.js already installed ($NODE_VERSION)"
else
    if brew install node &> /dev/null; then
        print_success "Node.js installed ($(node --version))"
    else
        print_error "Failed to install Node.js"
    fi
fi

# npm packages
if check_command "npm"; then
    echo "Installing npm global packages..."

    # snyk
    if npm list -g snyk &> /dev/null; then
        print_skip "snyk already installed globally"
    else
        echo "Installing snyk..."
        if npm install -g snyk &> /dev/null; then
            print_success "snyk installed globally"
        else
            print_error "Failed to install snyk"
        fi
    fi
fi

################################################################################
# Database Tools
################################################################################

print_header "Installing Database Tools"

# PostgreSQL client
install_brew_package "postgresql@18" "PostgreSQL 18 client"

# Redis CLI
install_brew_package "redis" "Redis"

################################################################################
# Docker & Compose
################################################################################

print_header "Installing Docker & Docker Compose"

# Note: Docker Desktop for Mac is the recommended approach
if check_command "docker"; then
    print_skip "Docker is already installed"
    echo "  $(docker --version)"
else
    print_info "Installing Docker Desktop for Mac..."
    install_brew_cask "docker" "Docker Desktop"
fi

# Verify Docker Compose is available
if check_command "docker-compose"; then
    print_skip "Docker Compose is already installed"
    echo "  $(docker-compose --version)"
else
    print_info "Docker Compose should be included with Docker Desktop"
    print_info "If missing, run: brew install docker-compose"
fi

################################################################################
# Additional Development Tools
################################################################################

print_header "Installing Additional Development Tools"

# Colima (Docker runtime alternative - optional)
if check_command "colima"; then
    print_skip "Colima already installed"
else
    print_info "Colima (Docker runtime for macOS) is optional"
    print_info "To install: brew install colima"
fi

# OrbStack (faster Docker alternative - optional)
print_info "For better performance, consider: brew install orbstack"

################################################################################
# Shell Configuration
################################################################################

print_header "Shell Configuration"

# Detect shell
SHELL_NAME=$(basename "$SHELL")
echo "Current shell: $SHELL_NAME"

if [[ "$SHELL_NAME" == "bash" ]]; then
    SHELL_RC="$HOME/.bashrc"
elif [[ "$SHELL_NAME" == "zsh" ]]; then
    SHELL_RC="$HOME/.zshrc"
elif [[ "$SHELL_NAME" == "fish" ]]; then
    SHELL_RC="$HOME/.config/fish/config.fish"
else
    SHELL_RC="unknown"
fi

if [[ "$SHELL_RC" != "unknown" ]]; then
    print_info "Configuration file: $SHELL_RC"

    # Ensure Go bin directory is in PATH
    if grep -q "GOPATH/bin" "$SHELL_RC" 2>/dev/null || grep -q "go/bin" "$SHELL_RC" 2>/dev/null; then
        print_skip "GOPATH/bin already in PATH"
    else
        print_info "Adding Go bin directory to PATH..."
        echo "" >> "$SHELL_RC"
        echo "# Venio Development Environment" >> "$SHELL_RC"
        echo "export GOPATH=\$(go env GOPATH)" >> "$SHELL_RC"
        echo "export PATH=\$GOPATH/bin:\$PATH" >> "$SHELL_RC"
        print_success "PATH updated"
    fi
fi

################################################################################
# Verification
################################################################################

print_header "Verifying Installation"

TOOLS=(
    "git:Git"
    "make:Make"
    "go:Go"
    "node:Node.js"
    "npm:npm"
    "docker:Docker"
    "psql:PostgreSQL Client"
    "redis-cli:Redis CLI"
)

for tool_check in "${TOOLS[@]}"; do
    IFS=':' read -r cmd name <<< "$tool_check"
    if check_command "$cmd"; then
        echo -e "${GREEN}✓${NC} $name - $(${cmd} --version 2>&1 | head -n 1)"
    else
        echo -e "${RED}✗${NC} $name - NOT FOUND"
    fi
done

################################################################################
# Post-Installation Summary
################################################################################

print_header "Installation Summary"

echo "Installed: $INSTALLED"
echo "Skipped:   $SKIPPED"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed:    $FAILED${NC}"
fi

################################################################################
# Next Steps
################################################################################

print_header "Next Steps"

echo -e "${YELLOW}1. Start Docker Desktop${NC}"
echo "   - Open Applications/Docker.app"
echo "   - Wait for Docker daemon to start"

echo -e "\n${YELLOW}2. Reload shell configuration${NC}"
if [[ "$SHELL_NAME" == "zsh" ]]; then
    echo "   source ~/.zshrc"
elif [[ "$SHELL_NAME" == "bash" ]]; then
    echo "   source ~/.bashrc"
elif [[ "$SHELL_NAME" == "fish" ]]; then
    echo "   source ~/.config/fish/config.fish"
fi

echo -e "\n${YELLOW}3. Verify Go tools are in PATH${NC}"
echo "   go version"
echo "   which air"
echo "   which golangci-lint"

echo -e "\n${YELLOW}4. Start development${NC}"
echo "   cd c:\\Users\\ms\\dev\\venio"
echo "   make dev          # Start Docker services"
echo "   make watch        # Run with hot reload"
echo "   make test-api     # Run API tests"

echo -e "\n${YELLOW}5. Optional: Install snyk authentication${NC}"
echo "   snyk auth"

echo -e "\n${BLUE}Happy coding!${NC}\n"
