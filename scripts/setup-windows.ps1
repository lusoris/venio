# Venio Windows Development Setup Script
# This script installs all required tools and dependencies for Windows development

$ErrorActionPreference = "Continue"

Write-Host "=== Venio Windows Development Setup ===" -ForegroundColor Cyan
Write-Host ""

# Check if running as admin
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")
if (-not $isAdmin) {
    Write-Host "WARNING: Recommend running as Administrator" -ForegroundColor Yellow
}

# Color functions
function Write-Success { Write-Host $args -ForegroundColor Green }
function Write-Fail { Write-Host "ERROR: $args" -ForegroundColor Red }
function Write-Warn { Write-Host "WARNING: $args" -ForegroundColor Yellow }
function Write-Info { Write-Host $args -ForegroundColor Blue }

# 1. Check Go
Write-Info "1. Checking Go installation..."
$goCmd = Get-Command go -ErrorAction SilentlyContinue
if ($goCmd) {
    $goVersion = & go version 2>&1
    Write-Success "✓ Go is installed: $goVersion"
} else {
    Write-Fail "✗ Go is NOT installed"
    Write-Info ""
    Write-Host "To install Go 1.23+:" -ForegroundColor Yellow
    Write-Host "  1. Download from: https://go.dev/dl/"
    Write-Host "  2. Choose 'go1.23.x.windows-amd64.msi' for your system"
    Write-Host "  3. Run the installer"
    Write-Host "  4. Restart your terminal/PowerShell"
    Write-Host "  5. Run this script again"
    Write-Host ""
    exit 1
}

# 2. Check Docker
Write-Info ""
Write-Info "2. Checking Docker..."
$dockerCmd = Get-Command docker -ErrorAction SilentlyContinue
if ($dockerCmd) {
    $dockerVersion = & docker --version 2>&1
    Write-Success "✓ Docker is installed: $dockerVersion"
} else {
    Write-Fail "✗ Docker is NOT installed"
    Write-Info ""
    Write-Host "To install Docker Desktop:" -ForegroundColor Yellow
    Write-Host "  1. Download from: https://www.docker.com/products/docker-desktop"
    Write-Host "  2. Run the installer"
    Write-Host "  3. Complete Docker Desktop setup"
    Write-Host "  4. Restart your terminal/PowerShell"
    Write-Host "  5. Run this script again"
    Write-Host ""
    exit 1
}

# 3. Check Docker Compose
Write-Info "3. Checking Docker Compose..."
$composeCmd = Get-Command "docker" -ErrorAction SilentlyContinue
if ($composeCmd) {
    try {
        $composeVersion = & docker compose version 2>&1
        Write-Success "✓ Docker Compose is installed: $composeVersion"
    } catch {
        Write-Fail "✗ Docker Compose is NOT available"
        Write-Info "Docker Compose is included with Docker Desktop"
        Write-Info "Please reinstall Docker Desktop from: https://www.docker.com/products/docker-desktop"
        exit 1
    }
} else {
    Write-Fail "✗ Docker is not available"
    exit 1
}

# 4. Install Go Development Tools
Write-Info ""
Write-Info "4. Installing Go development tools..."

$tools = @(
    "github.com/cosmtrek/air@latest",
    "github.com/go-delve/delve/cmd/dlv@latest",
    "golang.org/x/tools/cmd/goimports@latest",
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
    "github.com/evilmartians/lefthook@latest"
)

$toolsFailed = 0
foreach ($tool in $tools) {
    Write-Host "  Installing $tool..."
    $output = & go install $tool 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "  ✓ $tool"
    } else {
        Write-Fail "  ✗ $tool"
        Write-Host "    Error: $output" -ForegroundColor Red
        $toolsFailed++
    }
}

if ($toolsFailed -gt 0) {
    Write-Fail "$toolsFailed tools failed to install"
    exit 1
}

# 5. Setup Lefthook
Write-Info ""
Write-Info "5. Setting up Git hooks with Lefthook..."
$output = & lefthook install 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Success "✓ Lefthook hooks installed"
} else {
    Write-Warn "Lefthook setup: $output"
}

# 6. Download Go modules
Write-Info ""
Write-Info "6. Downloading Go dependencies..."
Write-Info "  (This may take a minute...)"
$output = & go mod download 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Success "✓ Go modules downloaded"
} else {
    Write-Fail "Failed to download Go modules"
    Write-Host "Error: $output" -ForegroundColor Red
    exit 1
}

# 7. Create .env file
Write-Info ""
Write-Info "7. Setting up environment file..."
$envFile = ".env"
$envExample = ".env.example"

if (Test-Path $envFile) {
    Write-Success "✓ .env file already exists"
} else {
    if (Test-Path $envExample) {
        Copy-Item $envExample $envFile
        Write-Success "✓ Created .env from .env.example"
        Write-Host ""
        Write-Host "⚠ Please update .env with your settings:" -ForegroundColor Yellow
        Write-Host "  - POSTGRES_PASSWORD"
        Write-Host "  - REDIS_PASSWORD"
        Write-Host "  - JWT_SECRET (min 32 chars)"
    } else {
        Write-Warn ".env.example not found, skipping .env creation"
    }
}

# 8. Verify setup
Write-Info ""
Write-Info "8. Verifying installation..."
$setupOk = $true

@(
    "Go",
    "Docker",
    "Air",
    "Delve",
    "Golangci-lint"
) | ForEach-Object {
    $cmd = Get-Command $_ -ErrorAction SilentlyContinue
    if ($cmd) {
        Write-Success "  ✓ $_"
    } else {
        Write-Fail "  ✗ $_"
        $setupOk = $false
    }
}

# Final message
Write-Host ""
if ($setupOk) {
    Write-Success "=== Setup completed successfully! ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Option 1: Start full development environment (Docker)"
    Write-Host "  docker compose -f docker-compose.yml -f docker-compose.dev.yml up"
    Write-Host ""
    Write-Host "Option 2: Run locally without Docker"
    Write-Host "  docker compose up postgres redis typesense"
    Write-Host "  go run cmd/venio/main.go"
    Write-Host ""
    Write-Host "Option 3: Run with hot reload"
    Write-Host "  air"
    Write-Host ""
} else {
    Write-Fail "Setup completed with errors. Please fix the issues above."
    exit 1
}
