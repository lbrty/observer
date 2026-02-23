# Default command
default:
    @just --list

# Build the application
build:
    go build -o bin/observer ./cmd/observer

# Run the server
run:
    go run ./cmd/observer serve

# Create new migration
migrate-create name:
    go run ./cmd/observer migrate create {{name}}

# Apply migrations (forward only)
migrate-up:
    go run ./cmd/observer migrate up

# Show migration version
migrate-version:
    go run ./cmd/observer migrate version

# Generate RSA keys using the built-in keygen command
keygen:
    go run ./cmd/observer keygen

# Generate RSA keys using openssl (alternative)
generate-keys:
    #!/usr/bin/env bash
    mkdir -p keys
    echo "Generating RSA private key (4096 bits)..."
    openssl genrsa -out keys/jwt_rsa 4096
    echo "Generating RSA public key..."
    openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub
    echo "Setting permissions..."
    chmod 600 keys/jwt_rsa
    chmod 644 keys/jwt_rsa.pub
    echo "RSA keys generated successfully in keys/ directory"

# Run unit tests only (fast, no Docker)
test:
    go test -v -short ./...

# Run all tests including integration tests
test-all:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

# Generate OpenAPI spec from annotations
openapi:
    swag init -g cmd/observer/main.go -o api/swagger --parseDependency --parseInternal

# Generate mocks
generate-mocks:
    go generate ./...

# Run tests with race detector
test-race:
    go test -v -race ./...

# Benchmark tests
bench:
    go test -bench=. -benchmem ./...

# Start docker compose
docker-up:
    docker-compose up -d

# Stop docker compose
docker-down:
    docker-compose down

# Clean build artifacts
clean:
    rm -rf bin/
    rm -f *.pem
    rm -f coverage.out coverage.html

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    golangci-lint run

# Tidy dependencies
tidy:
    go mod tidy
