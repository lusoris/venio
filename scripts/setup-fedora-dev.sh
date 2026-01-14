#!/bin/bash

################################################################################
# Venio Development Environment Setup for Fedora/RHEL/CentOS
#
# This script sets up a complete development environment for Venio on Fedora,
# RHEL (Red Hat Enterprise Linux), and CentOS using DNF (Dandified Yum) or
# YUM package manager.
#
# Supported distributions:
#   - Fedora 39+ (recommended)
#   - RHEL 9+ (with EPEL repository)
#   - CentOS Stream 9+
#
# Usage:
#   sudo bash scripts/setup-fedora-dev.sh
#
# Requirements:
#   - sudo access
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

check_sudo() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}Error: This script must be run with sudo${NC}"
        echo "Usage: sudo bash scripts/setup-fedora-dev.sh"
        exit 1
    fi
}

install_package() {
    local package=$1
    local name=${2:-$package}
    
    echo "Installing $name..."
    if dnf install -y "$package" &> /dev/null; then
        print_success "$name installed"
    else
        print_error "Failed to install $name"
    fi
}

check_and_install() {
    local command=$1
    local package=$2
    local name=${3:-$package}
    
    if check_command "$command"; then
        print_skip "$name already installed"
    else
        install_package "$package" "$name"
    fi
}

################################################################################
# System Checks
################################################################################

print_header "System Requirements Check"

check_sudo

# Detect distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO=$ID
    VERSION=$VERSION_ID
    echo "Distribution: $PRETTY_NAME"
else
    print_error "Cannot detect Linux distribution"
    exit 1
fi

# Check if it's RHEL/Fedora-based
if [[ ! "$DISTRO" =~ ^(fedora|rhel|centos)$ ]]; then
    echo -e "${YELLOW}Warning: This script is optimized for Fedora/RHEL/CentOS${NC}"
    echo "Detected: $PRETTY_NAME"
fi

# Check for DNF/YUM
if check_command "dnf"; then
    PACKAGE_MANAGER="dnf"
    print_success "DNF package manager found"
elif check_command "yum"; then
    PACKAGE_MANAGER="yum"
    print_success "YUM package manager found"
else
    print_error "Neither DNF nor YUM found"
    exit 1
fi

# Update package manager
print_info "Updating package cache..."
if $PACKAGE_MANAGER update -y &> /dev/null; then
    print_success "Package cache updated"
else
    print_error "Failed to update package cache"
fi

# Enable EPEL for RHEL/CentOS (if needed)
if [[ "$DISTRO" =~ ^(rhel|centos)$ ]]; then
    print_info "Setting up EPEL repository for RHEL/CentOS..."
    if $PACKAGE_MANAGER install -y epel-release &> /dev/null; then
        print_success "EPEL repository enabled"
    fi
fi

################################################################################
# Core Development Tools
################################################################################

print_header "Installing Core Development Tools"

# Git
check_and_install "git" "git" "Git"

# Development tools group
echo "Installing development tools group..."
if $PACKAGE_MANAGER groupinstall -y "Development Tools" &> /dev/null; then
    print_success "Development tools installed (gcc, make, etc.)"
else
    print_error "Failed to install development tools"
fi

# Additional build tools
install_package "cmake" "CMake"
install_package "automake" "Automake"

# Curl & Wget
check_and_install "curl" "curl" "Curl"
check_and_install "wget" "wget" "Wget"

# jq (JSON processor)
install_package "jq" "jq (JSON processor)"

# SQLite
install_package "sqlite" "SQLite"

# Useful utilities
install_package "htop" "htop (process monitor)"
install_package "tmux" "tmux (terminal multiplexer)"

################################################################################
# Go Development
################################################################################

print_header "Installing Go Development Stack"

# Check current Go installation
if check_command "go"; then
    CURRENT_VERSION=$(go version | awk '{print $3}')
    echo "Current Go version: $CURRENT_VERSION"
    
    if [[ "$CURRENT_VERSION" == "go1.25"* ]]; then
        print_skip "Go 1.25 already installed"
    else
        print_info "Downloading Go 1.25..."
    fi
else
    print_info "Downloading Go 1.25..."
fi

# Download and install Go 1.25 if not already installed or if older version
if [[ ! $(go version 2>/dev/null) == *"go1.25"* ]]; then
    GO_VERSION="1.25"
    
    # Detect architecture
    ARCH=$(uname -m)
    if [[ "$ARCH" == "x86_64" ]]; then
        GO_ARCH="amd64"
    elif [[ "$ARCH" == "aarch64" ]]; then
        GO_ARCH="arm64"
    else
        print_error "Unsupported architecture: $ARCH"
        GO_ARCH="amd64"  # Default fallback
    fi
    
    GO_TARBALL="go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    GO_URL="https://golang.org/dl/${GO_TARBALL}"
    
    echo "Downloading Go from $GO_URL..."
    if curl -fsSL "$GO_URL" -o "/tmp/$GO_TARBALL"; then
        echo "Extracting Go..."
        
        # Remove old Go installation if it exists
        rm -rf /usr/local/go
        
        tar -C /usr/local -xzf "/tmp/$GO_TARBALL"
        rm "/tmp/$GO_TARBALL"
        
        # Update PATH
        if ! grep -q "/usr/local/go/bin" /etc/profile.d/venio.sh 2>/dev/null; then
            echo "export PATH=/usr/local/go/bin:\$PATH" | tee /etc/profile.d/venio.sh > /dev/null
        fi
        
        # Source the new PATH
        export PATH=/usr/local/go/bin:$PATH
        
        print_success "Go 1.25 installed"
        echo "  $(go version)"
    else
        print_error "Failed to download Go"
    fi
fi

# Go development tools
print_info "\nInstalling Go development tools..."

export PATH=/usr/local/go/bin:$PATH
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH

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

# NodeSource repository for latest Node.js LTS
if [[ "$DISTRO" == "fedora" ]]; then
    print_info "Installing Node.js from distribution repositories..."
    install_package "nodejs" "Node.js"
    install_package "npm" "npm"
else
    # For RHEL/CentOS, use NodeSource repository
    print_info "Setting up NodeSource repository..."
    
    # Get Node major version (currently LTS is 22.x)
    NODE_MAJOR=22
    
    if ! rpm --import https://rpm.nodesource.com/pubkey.gpg &> /dev/null; then
        print_error "Failed to import NodeSource GPG key"
    else
        print_success "NodeSource GPG key imported"
    fi
    
    # Install NodeSource repository
    if curl -sL https://rpm.nodesource.com/setup_${NODE_MAJOR}.x | bash - &> /dev/null; then
        install_package "nodejs" "Node.js"
    else
        print_error "Failed to setup NodeSource repository, using distro packages"
        install_package "nodejs" "Node.js"
    fi
fi

# npm global packages
if check_command "npm"; then
    echo "Installing npm global packages..."
    
    # snyk
    echo "Installing snyk..."
    if npm install -g snyk &> /dev/null; then
        print_success "snyk installed globally"
    else
        print_error "Failed to install snyk"
    fi
fi

################################################################################
# Database Tools
################################################################################

print_header "Installing Database Tools"

# PostgreSQL client
install_package "postgresql" "PostgreSQL 18 client"

# Redis CLI
install_package "redis" "Redis"

################################################################################
# Docker & Container Tools
################################################################################

print_header "Installing Docker & Container Tools"

# Docker
echo "Installing Docker..."
if dnf config-manager --add-repo=https://download.docker.com/linux/fedora/docker-ce.repo &> /dev/null && \
   $PACKAGE_MANAGER install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin &> /dev/null; then
    print_success "Docker installed"
    
    # Enable and start Docker daemon
    echo "Enabling Docker daemon..."
    systemctl enable docker
    systemctl start docker
    print_success "Docker daemon configured to start on boot"
    
    # Add docker group
    if ! getent group docker > /dev/null; then
        echo "Creating docker group..."
        groupadd docker
        print_success "docker group created"
    else
        print_skip "docker group already exists"
    fi
    
    # Add current user to docker group
    CURRENT_USER=${SUDO_USER:-$USER}
    if ! groups "$CURRENT_USER" | grep -q "\bdocker\b"; then
        echo "Adding user to docker group..."
        usermod -aG docker "$CURRENT_USER"
        print_success "User added to docker group"
        print_info "Note: You may need to log out and log back in for this to take effect"
        print_info "Or run: newgrp docker"
    fi
else
    print_error "Failed to install Docker"
    print_info "Alternative: Follow Docker documentation at https://docs.docker.com/engine/install/fedora/"
fi

# Docker Compose (if not already installed via docker-compose-plugin)
if check_command "docker-compose"; then
    print_skip "Docker Compose already installed"
else
    print_info "Docker Compose should be available via 'docker compose' command"
fi

################################################################################
# Container Security & Scanning
################################################################################

print_header "Installing Container Security Tools"

# Podman (alternative to Docker)
# Uncomment if you prefer Podman
# install_package "podman" "Podman"
# install_package "podman-compose" "Podman Compose"

# Skopeo (inspect container images)
install_package "skopeo" "Skopeo (image inspector)"

################################################################################
# Shell Configuration
################################################################################

print_header "Shell Configuration"

# Setup PATH for all users via /etc/profile.d
if [ ! -f /etc/profile.d/venio.sh ]; then
    print_info "Creating system-wide Venio configuration..."
    cat > /etc/profile.d/venio.sh << 'EOF'
# Venio Development Environment
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
EOF
    chmod 644 /etc/profile.d/venio.sh
    print_success "System configuration created"
fi

# For current shell session
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

print_info "Shell configuration file: /etc/profile.d/venio.sh"

################################################################################
# Verification
################################################################################

print_header "Verifying Installation"

TOOLS=(
    "git:Git"
    "gcc:GCC"
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

CURRENT_USER=${SUDO_USER:-$USER}

echo -e "${YELLOW}1. Apply group changes${NC}"
echo "   Run one of:"
echo "   - newgrp docker"
echo "   - Log out and log back in"

echo -e "\n${YELLOW}2. Start Docker daemon${NC}"
echo "   sudo systemctl start docker"

echo -e "\n${YELLOW}3. Reload environment variables${NC}"
echo "   source /etc/profile.d/venio.sh"
echo "   # Or open a new terminal"

echo -e "\n${YELLOW}4. Verify Go tools are in PATH${NC}"
echo "   which air"
echo "   which golangci-lint"

echo -e "\n${YELLOW}5. Start development${NC}"
echo "   cd /path/to/venio"
echo "   make dev          # Start Docker services"
echo "   make watch        # Run with hot reload"
echo "   make test-api     # Run API tests"

echo -e "\n${YELLOW}6. Optional: Set up Snyk${NC}"
echo "   snyk auth"
echo "   snyk test"

echo -e "\n${BLUE}Happy coding!${NC}\n"
