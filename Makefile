.PHONY: help dev run watch test test-coverage test-integration lint format build docker-build docker-push migrate-up migrate-down clean

# Default target
help:
	@echo "Venio Development Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  dev               - Start development environment (Docker Compose)"
	@echo "  run               - Run locally without Docker"
	@echo "  watch             - Run with hot reload"
	@echo "  test              - Run tests"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-integration  - Run integration tests"
	@echo "  lint              - Run linters"
	@echo "  format            - Format code"
	@echo "  build             - Build binary"
	@echo "  docker-build      - Build Docker image"
	@echo "  docker-push       - Push Docker image to registry"
	@echo "  migrate-up        - Run database migrations"
	@echo "  migrate-down      - Rollback database migrations"
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
	@echo "Running migrations..."
	# Add migration tool command here

migrate-down:
	@echo "Rolling back migrations..."
	# Add migration tool command here

# Cleanup
clean:
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean
