# Setup Windows Development Environment for Venio
# This script automates the installation of all required tools

param(
    [switch]$SkipConfirm = $false
)

# Colors for output
$ErrorColor = 'Red'
$SuccessColor = 'Green'
$WarningColor = 'Yellow'
$InfoColor = 'Cyan'

function Write-Success {
    param([string]$Message)
    Write-Host $Message -ForegroundColor $SuccessColor
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host $Message -ForegroundColor $ErrorColor
}

function Write-Warning-Custom {
    param([string]$Message)
    Write-Host $Message -ForegroundColor $WarningColor
}

function Write-Info {
    param([string]$Message)
    Write-Host $Message -ForegroundColor $InfoColor
}

# Check if running as Administrator
$isAdmin = [Security.Principal.WindowsWindowsPrincipal]::new([Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Error-Custom "ERROR: This script must be run as Administrator!"
    Write-Warning-Custom "Please right-click PowerShell and select 'Run as Administrator'"
    exit 1
}

Write-Info "╔═══════════════════════════════════════════════════════════════╗"
Write-Info "║    Venio Windows Development Environment Setup              ║"
Write-Info "║    This script will install all required development tools   ║"
Write-Info "╚═══════════════════════════════════════════════════════════════╝"
Write-Host ""

# Check winget
Write-Info "Checking for winget (Windows Package Manager)..."
$wingetPath = Get-Command winget -ErrorAction SilentlyContinue

if (-not $wingetPath) {
    Write-Warning-Custom "⚠️  winget not found. Installing App Installer..."
    
    try {
        Add-AppxPackage -RegisterByFamilyName -MainPackage Microsoft.DesktopAppInstaller_8wekyb3d8bbwe
        Write-Success "✓ App Installer installed"
    } catch {
        Write-Error-Custom "✗ Failed to install App Installer"
        Write-Info "Please download from Microsoft Store or visit: https://github.com/microsoft/winget-cli"
        exit 1
    }
}

Write-Success "✓ winget is available"
Write-Host ""

# Define packages to install
$packages = @(
    @{ name = "Go"; id = "GoLang.Go"; command = "go" },
    @{ name = "Docker Desktop"; id = "Docker.DockerDesktop"; command = "docker" },
    @{ name = "Git"; id = "Git.Git"; command = "git" },
    @{ name = "GNU Make"; id = "GnuWin32.Make"; command = "make" },
    @{ name = "Node.js"; id = "OpenJS.NodeJS"; command = "node" }
)

Write-Info "═══════════════════════════════════════════════════════════════"
Write-Info "Installing Required Packages"
Write-Info "═══════════════════════════════════════════════════════════════"
Write-Host ""

$installed = 0
$skipped = 0
$failed = 0

foreach ($package in $packages) {
    Write-Host "Checking $($package.name)..." -NoNewline
    
    # Check if already installed
    $existingCommand = Get-Command $package.command -ErrorAction SilentlyContinue
    
    if ($existingCommand) {
        Write-Success " ✓ (already installed)"
        $skipped++
        continue
    }
    
    Write-Host ""
    Write-Host "  Installing $($package.name)..." -NoNewline
    
    try {
        $output = & winget install $package.id -e --accept-source-agreements --accept-package-agreements 2>&1
        
        if ($LASTEXITCODE -eq 0 -or $output -match "Successfully installed") {
            Write-Success " ✓"
            $installed++
        } else {
            Write-Error-Custom " ✗"
            Write-Info "  Command: winget install $($package.id) -e"
            $failed++
        }
    } catch {
        Write-Error-Custom " ✗"
        Write-Info "  Error: $_"
        $failed++
    }
}

Write-Host ""
Write-Success "Package Installation Summary: $installed installed, $skipped already present"

if ($failed -gt 0) {
    Write-Error-Custom "⚠️  $failed packages failed. Please install manually and re-run this script."
}

# Install Go development tools
Write-Host ""
Write-Info "═══════════════════════════════════════════════════════════════"
Write-Info "Installing Go Development Tools"
Write-Info "═══════════════════════════════════════════════════════════════"
Write-Host ""

$goTools = @(
    @{ name = "Air (hot reload)"; package = "github.com/cosmtrek/air@latest" },
    @{ name = "Delve (debugger)"; package = "github.com/go-delve/delve/cmd/dlv@latest" },
    @{ name = "goimports (formatter)"; package = "golang.org/x/tools/cmd/goimports@latest" },
    @{ name = "golangci-lint (linter)"; package = "github.com/golangci/golangci-lint/cmd/golangci-lint@latest" },
    @{ name = "Lefthook (git hooks)"; package = "github.com/evilmartians/lefthook@latest" }
)

$toolsInstalled = 0
$toolsFailed = 0

foreach ($tool in $goTools) {
    Write-Host "Installing $($tool.name)..." -NoNewline
    
    try {
        $output = & go install $tool.package 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success " ✓"
            $toolsInstalled++
        } else {
            Write-Error-Custom " ✗"
            Write-Info "  Command: go install $($tool.package)"
            Write-Info "  Output: $output"
            $toolsFailed++
        }
    } catch {
        Write-Error-Custom " ✗"
        Write-Info "  Error: $_"
        $toolsFailed++
    }
}

Write-Host ""
Write-Success "Go Tools Installation Summary: $toolsInstalled installed"

if ($toolsFailed -gt 0) {
    Write-Warning-Custom "⚠️  $toolsFailed Go tools failed. You can install them manually:"
    foreach ($tool in $goTools) {
        Write-Info "  go install $($tool.package)"
    }
}

# Restart explorer to update PATH
Write-Host ""
Write-Info "Restarting Windows Explorer to update PATH..."
Stop-Process -Name explorer -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 2
Start-Process explorer

# Final checks
Write-Host ""
Write-Info "═══════════════════════════════════════════════════════════════"
Write-Info "Environment Verification"
Write-Info "═══════════════════════════════════════════════════════════════"
Write-Host ""

$checks = @(
    @{ name = "Go"; command = "go version" },
    @{ name = "Docker"; command = "docker --version" },
    @{ name = "Git"; command = "git --version" },
    @{ name = "Make"; command = "make --version" },
    @{ name = "Node.js"; command = "node --version" }
)

$verified = 0
$unverified = 0

foreach ($check in $checks) {
    Write-Host "Checking $($check.name)..." -NoNewline
    
    # Need to reload PATH in current session for newly installed tools
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
    
    try {
        $output = Invoke-Expression $check.command 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success " ✓"
            Write-Info "  $output"
            $verified++
        } else {
            Write-Error-Custom " ✗ (failed to execute)"
            $unverified++
        }
    } catch {
        Write-Error-Custom " ✗ (not found - may need terminal restart)"
        $unverified++
    }
}

Write-Host ""
Write-Success "Verification Summary: $verified working, $unverified need attention"

# Setup instructions
Write-Host ""
Write-Info "═══════════════════════════════════════════════════════════════"
Write-Info "Next Steps"
Write-Info "═══════════════════════════════════════════════════════════════"
Write-Host ""
Write-Info "1. RESTART YOUR TERMINAL/POWERSHELL (important for PATH updates)"
Write-Host ""
Write-Info "2. Clone the repository:"
Write-Host "   git clone https://github.com/lusoris/venio.git"
Write-Host "   cd venio"
Write-Host ""
Write-Info "3. Copy environment template:"
Write-Host "   cp .env.example .env"
Write-Host ""
Write-Info "4. Edit .env with your settings (passwords, secrets, etc.)"
Write-Host ""
Write-Info "5. Start Docker services:"
Write-Host "   docker compose up postgres redis -d"
Write-Host ""
Write-Info "6. Run the application:"
Write-Host "   go run cmd/venio/main.go"
Write-Host ""
Write-Info "7. In a new terminal, start the frontend:"
Write-Host "   cd web"
Write-Host "   npm install"
Write-Host "   npm run dev"
Write-Host ""
Write-Info "8. Access the application:"
Write-Host "   Backend: http://localhost:3690"
Write-Host "   Frontend: http://localhost:3000"
Write-Host ""
Write-Info "For more details, see: docs/windows-setup.md"
Write-Host ""
Write-Success "╔═══════════════════════════════════════════════════════════════╗"
Write-Success "║         Setup Complete! Restart your terminal.              ║"
Write-Success "╚═══════════════════════════════════════════════════════════════╝"
