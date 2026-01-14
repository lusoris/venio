.PHONY: help dev run watch test test-coverage test-integration lint format build docker-build docker-push migrate-up migrate-down db-reset db-shell docker-up docker-down docker-logs docker-ps install-tools install setup clean

# Default target
help:
	@echo "Venio Development Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  dev               - Start development environment (Docker Compose)"
	@echo "  run               - Run locally without Docker"
	@echo "  watch             - Run with hot reload (Air)"
	@echo "  test              - Run tests"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-integration  - Run integration tests"
	@echo "  lint              - Run linters"
	@echo "  format            - Format code"
	@echo "  build             - Build binary"
	@echo "  docker-build      - Build Docker image"
	@echo "  docker-push       - Push Docker image to registry"
	@echo "  docker-up         - Start PostgreSQL and Redis"
	@echo "  docker-down       - Stop all Docker services"
	@echo "  docker-logs       - View Docker logs"
	@echo "  migrate-up        - Run database migrations"
	@echo "  migrate-down      - Rollback database migrations"
	@echo "  db-reset          - Reset database (down then up)"
	@echo "  db-shell          - Open PostgreSQL shell"
	@echo "  install-tools     - Install development tools"
	@echo "  install           - Install tools and dependencies"
	@echo "  setup             - Complete project setup"
	@echo "  clean             - Remove build artifacts"

# Development
dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up

run:
	go run cmd/venio/main.go

watch:
	air

# Testing
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration:
	docker compose -f docker-compose.test.yml up --abort-on-container-exit
	docker compose -f docker-compose.test.yml down

# Code Quality
lint:
	golangci-lint run

format:
	goimports -w .
	gofmt -s -w .

# Building
build:
	mkdir -p bin
	go build -o bin/venio cmd/venio/main.go
	go build -o bin/worker cmd/worker/main.go

docker-build:
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/USERNAME/venio:dev .

docker-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/USERNAME/venio:dev --push .

# Database
migrate-up:
	@echo "⬆️  Running migrations..."
	@powershell -Command "Get-Content migrations/001_initial_schema.up.sql | docker exec -i venio-postgres psql -U venio -d venio"
	@echo "✅ Migrations complete"

migrate-down:
	@echo "⬇️  Rolling back migrations..."
	@powershell -Command "Get-Content migrations/001_initial_schema.down.sql | docker exec -i venio-postgres psql -U venio -d venio"
	@echo "✅ Rollback complete"

db-reset: migrate-down migrate-up

db-shell:
	@echo "🐘 Opening PostgreSQL shell..."
	docker exec -it venio-postgres psql -U venio -d venio

# Docker helpers
docker-up:
	@echo "🐳 Starting Docker services..."
	docker compose up -d postgres redis
	@echo "✅ Docker services started"

docker-down:
	@echo "🐳 Stopping Docker services..."
	docker compose down

docker-logs:
	docker compose logs -f

docker-ps:
	docker compose ps

# Installation
install-tools:
	@echo "🔧 Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ Tools installed"

install: install-tools
	@echo "📦 Downloading dependencies..."
	go mod download
	@echo "✅ Dependencies downloaded"

# Complete setup
setup: install docker-up migrate-up
	@echo ""
	@echo "🎉 Setup complete!"
	@echo ""
	@echo "Run 'make watch' to start development server with hot reload"
	@echo "Run 'make dev' to start full Docker environment"

# Cleanup
clean:
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean
