# Venio Windows Development Setup Script
# Auto-installs all required tools using winget

$ErrorActionPreference = "Continue"

Write-Host "=== Venio Windows Development Setup ===" -ForegroundColor Cyan
Write-Host ""

# Check if running as admin
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")
if (-not $isAdmin) {
    Write-Host "WARNING: Some features require Administrator privileges" -ForegroundColor Yellow
}

# Color functions
function Write-Success { Write-Host $args -ForegroundColor Green }
function Write-Fail { Write-Host "ERROR: $args" -ForegroundColor Red }
function Write-Warn { Write-Host "WARNING: $args" -ForegroundColor Yellow }
function Write-Info { Write-Host $args -ForegroundColor Blue }

# 1. Check and install Go
Write-Info "1. Checking Go installation..."
$goCmd = Get-Command go -ErrorAction SilentlyContinue
if ($goCmd) {
    $goVersion = & go version 2>&1
    Write-Success "✓ Go is already installed: $goVersion"
} else {
    Write-Warn "Go not found. Installing with winget..."

    $wingetCmd = Get-Command winget -ErrorAction SilentlyContinue
    if ($wingetCmd) {
        Write-Info "Installing Go 1.23+..."
        & winget install -e --id GoLang.Go --source winget 2>&1 | Select-String -Pattern "Successfully|Failed|error" | ForEach-Object { Write-Host $_ }

        # Refresh PATH
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
        Start-Sleep -Seconds 1

        $goCmd = Get-Command go -ErrorAction SilentlyContinue
        if ($goCmd) {
            $goVersion = & go version 2>&1
            Write-Success "✓ Go installed: $goVersion"
        } else {
            Write-Fail "Go installation failed"
            exit 1
        }
    } else {
        Write-Fail "winget not available. Manual installation required."
        Write-Host "Download from: https://go.dev/dl/" -ForegroundColor Yellow
        exit 1
    }
}

# 2. Check and install Docker
Write-Info ""
Write-Info "2. Checking Docker installation..."
$dockerCmd = Get-Command docker -ErrorAction SilentlyContinue
if ($dockerCmd) {
    $dockerVersion = & docker --version 2>&1
    Write-Success "✓ Docker is already installed: $dockerVersion"
} else {
    Write-Warn "Docker not found. Installing with winget..."

    $wingetCmd = Get-Command winget -ErrorAction SilentlyContinue
    if ($wingetCmd) {
        Write-Info "Installing Docker Desktop (this may take a few minutes)..."
        & winget install -e --id Docker.DockerDesktop --source winget 2>&1 | Select-String -Pattern "Successfully|Failed|error" | ForEach-Object { Write-Host $_ }

        Write-Warn ""
        Write-Warn "⚠ Docker Desktop installed but requires system restart"
        Write-Info ""
        Write-Host "Please restart your computer, then run this script again." -ForegroundColor Yellow
        exit 0
    } else {
        Write-Fail "winget not available. Manual installation required."
        Write-Host "Download from: https://www.docker.com/products/docker-desktop" -ForegroundColor Yellow
        exit 1
    }
}

# 3. Check Docker Compose
Write-Info "3. Checking Docker Compose..."
try {
    $composeVersion = & docker compose version 2>&1
    Write-Success "✓ Docker Compose is installed"
} catch {
    Write-Fail "Docker Compose not available"
    exit 1
}

# 4. Install Go development tools
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
        Write-Fail "  ✗ $tool failed"
        $toolsFailed++
    }
}

if ($toolsFailed -gt 0) {
    Write-Warn "$toolsFailed tools had issues but continuing..."
}

# 5. Setup Lefthook
Write-Info ""
Write-Info "5. Setting up Git hooks..."
$output = & lefthook install 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Success "✓ Lefthook hooks installed"
} else {
    Write-Warn "Lefthook setup: OK (optional)"
}

# 6. Download Go modules
Write-Info ""
Write-Info "6. Downloading Go dependencies..."
Write-Info "  (This may take a minute...)"
$output = & go mod download 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Success "✓ Go modules downloaded"
} else {
    Write-Fail "Failed to download modules"
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
        Write-Warn ".env.example not found"
    }
}

# 8. Final verification
Write-Info ""
Write-Info "8. Verifying setup..."
$setupOk = $true

$verifyTools = @("Go", "Docker", "Air", "Delve", "Golangci-lint")
foreach ($tool in $verifyTools) {
    $cmd = Get-Command $tool -ErrorAction SilentlyContinue
    if ($cmd) {
        Write-Success "  ✓ $tool"
    } else {
        Write-Fail "  ✗ $tool"
        $setupOk = $false
    }
}

# Final message
Write-Host ""
if ($setupOk) {
    Write-Success "=== Setup completed successfully! ==="
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Option 1: Full Docker development stack"
    Write-Host "  docker compose -f docker-compose.yml -f docker-compose.dev.yml up"
    Write-Host ""
    Write-Host "Option 2: Local with Docker services"
    Write-Host "  docker compose up postgres redis typesense"
    Write-Host "  go run cmd/venio/main.go"
    Write-Host ""
    Write-Host "Option 3: With hot reload (air)"
    Write-Host "  air"
    Write-Host ""
} else {
    Write-Warn "Setup completed with some issues"
    Write-Host "Re-run the script after fixing the problems above" -ForegroundColor Yellow
    exit 1
}
