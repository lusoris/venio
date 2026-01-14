# Venio Build Script for Windows
# Usage: .\build.ps1 <command>

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Venio Development Script" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Available commands:" -ForegroundColor Yellow
    Write-Host "  dev             - Start development server with hot reload (Air)" -ForegroundColor Green
    Write-Host "  run             - Run server without hot reload" -ForegroundColor Green
    Write-Host "  test            - Run all tests" -ForegroundColor Green
    Write-Host "  lint            - Run linter" -ForegroundColor Green
    Write-Host "  format          - Format code" -ForegroundColor Green
    Write-Host "  build           - Build binaries" -ForegroundColor Green
    Write-Host "  migrate-up      - Run database migrations" -ForegroundColor Green
    Write-Host "  migrate-down    - Rollback migrations" -ForegroundColor Green
    Write-Host "  db-reset        - Reset database" -ForegroundColor Green
    Write-Host "  db-shell        - Open PostgreSQL shell" -ForegroundColor Green
    Write-Host "  docker-up       - Start PostgreSQL and Redis" -ForegroundColor Green
    Write-Host "  docker-down     - Stop Docker services" -ForegroundColor Green
    Write-Host "  docker-logs     - View Docker logs" -ForegroundColor Green
    Write-Host "  install         - Install development tools" -ForegroundColor Green
    Write-Host "  setup           - Complete project setup" -ForegroundColor Green
    Write-Host "  clean           - Clean build artifacts" -ForegroundColor Green
    Write-Host ""
}

function Invoke-Dev {
    Write-Host "🚀 Starting Venio in development mode..." -ForegroundColor Cyan
    air
}

function Invoke-Run {
    Write-Host "🚀 Starting Venio..." -ForegroundColor Cyan
    go run cmd/venio/main.go
}

function Invoke-Test {
    Write-Host "🧪 Running tests..." -ForegroundColor Cyan
    go test -v -race -coverprofile=coverage.out ./...
    Write-Host "✅ Tests complete" -ForegroundColor Green
}

function Invoke-Lint {
    Write-Host "🔍 Running linter..." -ForegroundColor Cyan
    golangci-lint run ./...
    Write-Host "✅ Linting complete" -ForegroundColor Green
}

function Invoke-Format {
    Write-Host "✨ Formatting code..." -ForegroundColor Cyan
    gofmt -s -w .
    goimports -w .
    Write-Host "✅ Formatting complete" -ForegroundColor Green
}

function Invoke-Build {
    Write-Host "🔨 Building Venio..." -ForegroundColor Cyan
    if (-not (Test-Path "bin")) { New-Item -ItemType Directory -Path "bin" | Out-Null }
    go build -v -o bin/venio.exe ./cmd/venio
    go build -v -o bin/venio-worker.exe ./cmd/worker
    Write-Host "✅ Build complete: bin/venio.exe, bin/venio-worker.exe" -ForegroundColor Green
}

function Invoke-MigrateUp {
    Write-Host "⬆️  Running migrations..." -ForegroundColor Cyan
    Get-Content migrations/001_initial_schema.up.sql | docker exec -i venio-postgres psql -U venio -d venio
    Write-Host "✅ Migrations complete" -ForegroundColor Green
}

function Invoke-MigrateDown {
    Write-Host "⬇️  Rolling back migrations..." -ForegroundColor Cyan
    Get-Content migrations/001_initial_schema.down.sql | docker exec -i venio-postgres psql -U venio -d venio
    Write-Host "✅ Rollback complete" -ForegroundColor Green
}

function Invoke-DbReset {
    Invoke-MigrateDown
    Invoke-MigrateUp
}

function Invoke-DbShell {
    Write-Host "🐘 Opening PostgreSQL shell..." -ForegroundColor Cyan
    docker exec -it venio-postgres psql -U venio -d venio
}

function Invoke-DockerUp {
    Write-Host "🐳 Starting Docker services..." -ForegroundColor Cyan
    docker compose up -d postgres redis
    Write-Host "✅ Docker services started" -ForegroundColor Green
}

function Invoke-DockerDown {
    Write-Host "🐳 Stopping Docker services..." -ForegroundColor Cyan
    docker compose down
    Write-Host "✅ Docker services stopped" -ForegroundColor Green
}

function Invoke-DockerLogs {
    docker compose logs -f
}

function Invoke-Install {
    Write-Host "🔧 Installing development tools..." -ForegroundColor Cyan
    go install github.com/air-verse/air@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install golang.org/x/tools/cmd/goimports@latest
    Write-Host "📦 Downloading dependencies..." -ForegroundColor Cyan
    go mod download
    Write-Host "✅ Installation complete" -ForegroundColor Green
}

function Invoke-Setup {
    Invoke-Install
    Invoke-DockerUp
    Start-Sleep -Seconds 5
    Invoke-MigrateUp
    Write-Host ""
    Write-Host "🎉 Setup complete!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Run '.\build.ps1 dev' to start development server" -ForegroundColor Yellow
}

function Invoke-Clean {
    Write-Host "🧹 Cleaning..." -ForegroundColor Cyan
    if (Test-Path "bin") { Remove-Item -Recurse -Force bin }
    if (Test-Path "coverage.out") { Remove-Item coverage.out }
    if (Test-Path "coverage.html") { Remove-Item coverage.html }
    go clean
    Write-Host "✅ Clean complete" -ForegroundColor Green
}

# Main command dispatcher
switch ($Command.ToLower()) {
    "help" { Show-Help }
    "dev" { Invoke-Dev }
    "run" { Invoke-Run }
    "test" { Invoke-Test }
    "lint" { Invoke-Lint }
    "format" { Invoke-Format }
    "build" { Invoke-Build }
    "migrate-up" { Invoke-MigrateUp }
    "migrate-down" { Invoke-MigrateDown }
    "db-reset" { Invoke-DbReset }
    "db-shell" { Invoke-DbShell }
    "docker-up" { Invoke-DockerUp }
    "docker-down" { Invoke-DockerDown }
    "docker-logs" { Invoke-DockerLogs }
    "install" { Invoke-Install }
    "setup" { Invoke-Setup }
    "clean" { Invoke-Clean }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}
