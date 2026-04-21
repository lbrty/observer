# Default command
default:
    @just --list

# ── Dev ──────────────────────────────────────────────────────────────────────

# Run backend and frontend dev servers concurrently
dev:
    DEV_MODE=true SWAGGER_ENABLED=true COOKIE_SECURE=false go run ./cmd/observer serve & cd packages/observer-web && bun run dev

# Run the server (backend only)
run:
    SWAGGER_ENABLED=true go run ./cmd/observer serve

# ── Build ─────────────────────────────────────────────────────────────────────

# Build the Go binary (dev)
build:
    go build -o bin/observer ./cmd/observer

# Build frontend + Go binary with embedded SPA (production)
build-prod:
    cd packages/observer-web && bun run build
    go build -tags production -o bin/observer ./cmd/observer

# Clean build artifacts
clean:
    rm -rf bin/
    rm -f *.pem

# ── Database ──────────────────────────────────────────────────────────────────

# Apply migrations (forward only)
migrate-up:
    go run ./cmd/observer migrate up

# Create a new migration file
migrate-create name:
    go run ./cmd/observer migrate create {{name}}

# Show current migration version
migrate-version:
    go run ./cmd/observer migrate version

# Seed the database with realistic demo data
seed *args='':
    go run ./cmd/observer seed {{args}}

# Start docker compose
docker-up:
    docker-compose up -d

# Stop docker compose
docker-down:
    docker-compose down

# ── Auth & Admin ──────────────────────────────────────────────────────────────

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

# Create an admin user
create-admin email password *args='':
    go run ./cmd/observer create-admin --email {{email}} --password {{password}} {{args}}

# ── Test ──────────────────────────────────────────────────────────────────────

# Run unit tests only (fast, no Docker)
test:
    go test -v -short ./...

# Run all tests including integration tests
test-all:
    go test -v ./...

# ── Code quality ──────────────────────────────────────────────────────────────

# Format code (Go + frontend)
fmt:
    go fmt ./...
    cd packages/observer-web && bun run fmt

# Lint Go code
lint:
    golangci-lint run

# Tidy Go dependencies
tidy:
    go mod tidy

# Generate mocks
generate-mocks:
    go generate ./...

# Generate OpenAPI spec from annotations
openapi:
    swag init -g cmd/observer/main.go -o api/swagger --parseDependency --parseInternal
