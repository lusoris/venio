# install-vscode-extensions.ps1
# VSCode Extension Installer for Venio Project

Write-Host "Installing VSCode Extensions for Venio..." -ForegroundColor Cyan

# Check if code command is available
try {
    $null = Get-Command code -ErrorAction Stop
} catch {
    Write-Host "ERROR: VSCode 'code' command not found in PATH!" -ForegroundColor Red
    Write-Host "Please add VSCode to your PATH or run from VSCode terminal." -ForegroundColor Yellow
    exit 1
}

# Extension list
$extensions = @(
    # Go Development
    "golang.go",
    
    # Docker & Containers
    "ms-azuretools.vscode-docker",
    
    # Database
    "mtxr.sqltools",
    "mtxr.sqltools-driver-pg",
    
    # API Development & Testing
    "humao.rest-client",
    "rangav.vscode-thunder-client",
    
    # YAML & Config
    "redhat.vscode-yaml",
    "editorconfig.editorconfig",
    
    # Git
    "eamodio.gitlens",
    "github.vscode-github-actions",
    
    # Testing & Quality
    "hbenl.vscode-test-explorer",
    "ryanluker.vscode-coverage-gutters",
    "usernamehw.errorlens",
    
    # Documentation
    "davidanson.vscode-markdownlint",
    "yzhang.markdown-all-in-one",
    
    # Productivity
    "aaron-bond.better-comments",
    "gruntfuggly.todo-tree",
    "streetsidesoftware.code-spell-checker"
)

# Optional extensions (ask user)
$optionalExtensions = @(
    "github.copilot"
)

Write-Host "`nInstalling $($extensions.Count) required extensions..." -ForegroundColor Green

$installed = 0
$failed = 0

foreach ($ext in $extensions) {
    Write-Host "Installing $ext..." -NoNewline
    
    try {
        $output = code --install-extension $ext --force 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host " ✓" -ForegroundColor Green
            $installed++
        } else {
            Write-Host " ✗" -ForegroundColor Red
            Write-Host "  Error: $output" -ForegroundColor Red
            $failed++
        }
    } catch {
        Write-Host " ✗" -ForegroundColor Red
        Write-Host "  Error: $_" -ForegroundColor Red
        $failed++
    }
}

# Ask about optional extensions
Write-Host "`n--- Optional Extensions ---" -ForegroundColor Yellow
foreach ($ext in $optionalExtensions) {
    $response = Read-Host "Install $ext? (y/n)"
    if ($response -eq 'y' -or $response -eq 'Y') {
        Write-Host "Installing $ext..." -NoNewline
        
        try {
            $output = code --install-extension $ext --force 2>&1
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host " ✓" -ForegroundColor Green
                $installed++
            } else {
                Write-Host " ✗" -ForegroundColor Red
                $failed++
            }
        } catch {
            Write-Host " ✗" -ForegroundColor Red
            $failed++
        }
    }
}

# Summary
Write-Host "`n=== Installation Summary ===" -ForegroundColor Cyan
Write-Host "Successfully installed: $installed" -ForegroundColor Green
if ($failed -gt 0) {
    Write-Host "Failed: $failed" -ForegroundColor Red
}

Write-Host "`nRestart VSCode to activate all extensions." -ForegroundColor Yellow
Write-Host "Done!" -ForegroundColor Green
