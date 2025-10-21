# Memcached Management - Makefile

.PHONY: help build run test test-unit test-integration clean docker-up docker-down

# Default target
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run all tests"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-integration- Run integration tests only"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-up       - Start Memcached with Docker"
	@echo "  docker-down     - Stop Memcached Docker container"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/memcached-app cmd/main.go

# Run the application
run:
	@echo "Starting application..."
	go run cmd/main.go

# Run all tests
test:
	@echo "Running all tests..."
	go test -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v ./services ./models

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Start Memcached with Docker
docker-up:
	@echo "Starting Memcached..."
	docker-compose up -d

# Stop Memcached Docker container
docker-down:
	@echo "Stopping Memcached..."
	docker-compose down

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"