# ADR-001: Observer Basic Structure Implementation Plan

| Field      | Value                   |
| ---------- | ----------------------- |
| Status     | Accepted                |
| Date       | 2026-02-21              |
| Supersedes | —                       |
| Components | observer, bootstrapping |

---

Important: use variables defined in: `../variables.md`

## Phase 1: Project Initialization

### 1.1 Create Root Structure

```text
[project_name]/
├── cmd/
│   └── [project_name]/
│       ├── main.go
│       ├── main_test.go
│       └── cmd/
│           ├── serve.go
│           ├── serve_test.go
│           ├── migrate.go
│           ├── migrate_test.go
│           ├── keygen.go
│           └── keygen_test.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── database/
│   │   ├── database.go
│   │   ├── database_test.go
│   │   └── mock/
│   ├── server/
│   │   ├── server.go
│   │   └── server_test.go
│   ├── health/
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   └── mock/
│   ├── logger/
│   │   └── logger.go
│   ├── ulid/
│   │   ├── ulid.go
│   │   └── ulid_test.go
│   └── testutil/
│       ├── postgres.go
│       └── redis.go
├── migrations/
│   ├── 000001_init_extensions.up.sql
│   └── README.md
├── .env.example
├── .gitignore
├── Justfile
├── docker-compose.yml
└── README.md
```

### 1.2 Install Dependencies

```bash
# Core
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/gin-gonic/gin@latest
go get github.com/jmoiron/sqlx@latest
go get github.com/lib/pq@latest
go get github.com/golang-migrate/migrate/v4@latest
go get github.com/golang-migrate/migrate/v4/database/postgres@latest
go get github.com/golang-migrate/migrate/v4/source/file@latest
go get github.com/joho/godotenv@latest
go get github.com/oklog/ulid/v2@latest

# Testing
go get go.uber.org/mock@latest
go install go.uber.org/mock/mockgen@latest
go get github.com/stretchr/testify@latest
go get github.com/testcontainers/testcontainers-go@latest
go get github.com/testcontainers/testcontainers-go/modules/postgres@latest
go get github.com/testcontainers/testcontainers-go/modules/redis@latest
```

## Phase 2: Configuration

### 2.1 `internal/config/config.go`

```go
package config

import (
    "time"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Log      LogConfig
}

type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

type DatabaseConfig struct {
    DSN string
}

type LogConfig struct {
    Level string
}

func Load() (*Config, error)
```

**Configuration loading priority:**

1. Default constant values
2. Environment variables
3. CLI flags (override all)
4. .env file support in development (using godotenv)

**Default values:**

```go
const (
    DefaultServerHost         = "localhost"
    DefaultServerPort         = 9000
    DefaultServerReadTimeout  = 30 * time.Second
    DefaultServerWriteTimeout = 30 * time.Second
    DefaultLogLevel           = "info"
)
```

### 2.2 `.env.example`

```env
SERVER_HOST=localhost
SERVER_PORT=9000
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

DATABASE_DSN=postgresql://postgres:postgres@localhost:5432/[project_name]?sslmode=disable

LOG_LEVEL=info
```

### 2.3 `.gitignore`

```text
bin/
*.env
!.env.example
*.pem
coverage.out
coverage.html
vendor/
```

## Phase 3: Logging

### 3.1 `internal/logger/logger.go`

```go
package logger

import (
    "log/slog"
    "os"
)

func New(level string) *slog.Logger

func GinMiddleware(logger *slog.Logger) gin.HandlerFunc
```

**Features:**

- JSON format output to stdout
- Configurable log level (debug, info, warn, error)
- Request logging middleware for Gin
- Structured logging with context

## Phase 4: ULID

### 4.1 `internal/ulid/ulid.go`

```go
package ulid

import (
    "crypto/rand"
    "sync"

    "github.com/oklog/ulid/v2"
)

var (
    entropy     = ulid.Monotonic(rand.Reader, 0)
    entropyLock sync.Mutex
)

func New() string {
    entropyLock.Lock()
    defer entropyLock.Unlock()
    return ulid.MustNew(ulid.Now(), entropy).String()
}

func MustNew() string {
    return New()
}

func Parse(s string) (ulid.ULID, error) {
    return ulid.Parse(s)
}

func IsValid(s string) bool {
    _, err := ulid.Parse(s)
    return err == nil
}
```

### 4.2 `internal/ulid/ulid_test.go`

**Test cases:**

- ULID generation uniqueness
- ULID validation
- ULID parsing
- Thread safety

## Phase 5: Database

### 5.1 `internal/database/database.go`

```go
package database

import (
    "context"

    "github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=mock/database.go -package=mock [package_name]/internal/database DB

type DB interface {
    Ping(ctx context.Context) error
    Close() error
    GetDB() *sqlx.DB
}

type database struct {
    db *sqlx.DB
}

func New(dsn string) (DB, error)

func (d *database) Ping(ctx context.Context) error

func (d *database) Close() error

func (d *database) GetDB() *sqlx.DB
```

### 5.2 `migrations/000001_init_extensions.up.sql`

```sql
-- Enable useful PostgreSQL extensions

-- Case-insensitive text type (useful for emails, usernames)
CREATE EXTENSION IF NOT EXISTS citext;

-- Trigram matching for fuzzy text search (animal names, full names etc.)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Remove accents from text (useful for Cyrillic/Latin search)
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Cryptographic functions
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

**Index Naming Convention:**

- Regular indexes: `ix_<table_name>_<column_name(s)>`
- Unique constraints: `uq_<table_name>_<column_name(s)>`

**Examples:**

```sql
-- Regular index
CREATE INDEX ix_users_created_at ON users(created_at);

-- Composite index
CREATE INDEX ix_anima,_id_date ON users(id, date);

-- Unique constraint
CREATE UNIQUE INDEX uq_users_email ON users(email);

-- Unique composite
CREATE UNIQUE INDEX uq_animal_name ON time_slots(user_id, name);
```

### 5.4 `internal/database/database_test.go`

**Integration tests with testcontainers:**

- Database connection
- Ping functionality
- Connection pool

### 5.5 `internal/testutil/postgres.go`

```go
package testutil

import (
    "context"
    "testing"

    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupPostgres(t *testing.T) (string, func())
```

### 5.6 `internal/testutil/redis.go`

```go
package testutil

import (
    "context"
    "testing"

    "github.com/testcontainers/testcontainers-go/modules/redis"
)

func SetupRedis(t *testing.T) (string, func())
```

## Phase 6: Health Check

### 6.1 `internal/health/handler.go`

```go
package health

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "[package_name]/internal/database"
)

//go:generate mockgen -destination=mock/handler.go -package=mock [package_name]/internal/health Handler

type Handler interface {
    Health(c *gin.Context)
}

type handler struct {
    db database.DB
}

func NewHandler(db database.DB) Handler

func (h *handler) Health(c *gin.Context)
```

**Response format:**

```json
{ "status": "ok" }
```

**Error response (503):**

```json
{ "status": "not ok" }
```

### 6.2 `internal/health/handler_test.go`

**Test cases:**

- Successful health check
- Database connection failure
- Response format validation

## Phase 7: HTTP Server

### 7.1 `internal/server/server.go`

```go
package server

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "[package_name]/internal/config"
    "[package_name]/internal/database"
    "[package_name]/internal/health"
    "[package_name]/internal/logger"
    "[package_name]/internal/ulid"
)

type Server struct {
    router *gin.Engine
    server *http.Server
    config *config.ServerConfig
}

func New(cfg *config.Config, db database.DB, log *slog.Logger) *Server

func (s *Server) Start() error

func (s *Server) Shutdown(ctx context.Context) error

func (s *Server) setupMiddleware(log *slog.Logger)

func (s *Server) setupRoutes(healthHandler health.Handler)
```

**Middleware stack:**

1. Request ID (ULID)
2. Logger (slog)
3. Recovery
4. CORS

**Routes:**

- GET /health

### 7.2 `internal/server/server_test.go`

**Test cases:**

- Server initialization
- Health endpoint
- Graceful shutdown

## Phase 8: CLI Commands

### 8.1 `cmd/[project_name]/main.go`

```go
package main

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "Observer",
    Short: "Observer - IDP plarform",
}

func init() {
    cobra.OnInitialize(initConfig)
}

func initConfig() {
    // Load .env file in development
    _ = godotenv.Load()

    // Set up Viper
    viper.SetEnvPrefix("[project_name_uppercased]")
    viper.AutomaticEnv()

    // Bind environment variables
    // Bind flags
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

### 8.2 `cmd/[project_name]/main_test.go`

**Test cases:**

- Root command initialization
- Help text
- Version flag

### 8.3 `cmd/[project_name]/cmd/serve.go`

```go
package cmd

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the HTTP server",
    RunE:  runServe,
}

func init() {
    serveCmd.Flags().String("host", "", "Server host")
    serveCmd.Flags().Int("port", 0, "Server port")
}

func runServe(cmd *cobra.Command, args []string) error {
    // Load config
    // Initialize logger
    // Initialize database
    // Initialize server
    // Start server
    // Wait for shutdown signal
    // Graceful shutdown
}
```

**Flags:**

- --host (override config)
- --port (override config)

**Graceful shutdown:**

- Listen for SIGINT, SIGTERM
- Shutdown timeout: 30 seconds

### 8.4 `cmd/[project_name]/cmd/serve_test.go`

**Test cases:**

- Serve command initialization
- Flag parsing (host, port)
- Config override from flags
- Integration test with test server

### 8.5 `cmd/[project_name]/cmd/migrate.go`

```go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Database migration management",
}

var migrateUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Apply all pending migrations",
    RunE:  runMigrateUp,
}

var migrateCreateCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new migration file",
    Args:  cobra.ExactArgs(1),
    RunE:  runMigrateCreate,
}

var migrateVersionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show current migration version",
    RunE:  runMigrateVersion,
}

func init() {
    migrateCmd.AddCommand(migrateUpCmd)
    migrateCmd.AddCommand(migrateCreateCmd)
    migrateCmd.AddCommand(migrateVersionCmd)
}
```

**Subcommands:**

**migrate up:**

- Apply all pending .up.sql files
- Forward only, no rollback
- Show migration progress

**migrate create <name>:**

- Generate timestamped migration file
- Format: `{timestamp}_{name}.up.sql`
- Timestamp format: Unix timestamp or sequential

**migrate version:**

- Show current schema version
- Show dirty state if any

### 8.6 `cmd/[project_name]/cmd/migrate_test.go`

**Test cases:**

- Migrate command initialization
- Subcommand registration
- Create command validation (requires name argument)
- Migration file name format
- Integration test with test database (up, version)

### 8.7 `cmd/[project_name]/cmd/keygen.go`

```go
package cmd

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "os"

    "github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
    Use:   "keygen",
    Short: "Generate RSA key pair",
    RunE:  runKeygen,
}

func init() {
    keygenCmd.Flags().Int("bits", 4096, "RSA key size (minimum 4096)")
    keygenCmd.Flags().String("output", ".", "Output directory")
}

func runKeygen(cmd *cobra.Command, args []string) error {
    // Validate bits >= 4096
    // Generate RSA key pair
    // Save private_key.pem
    // Save public_key.pem
}
```

**Flags:**

- --bits (default: 4096, minimum: 4096)
- --output (default: current directory)

**Output files:**

- private_key.pem
- public_key.pem

### 8.8 `cmd/[project_name]/cmd/keygen_test.go`

**Test cases:**

- Keygen command initialization
- Flag parsing (bits, output)
- Validation: bits >= 4096
- Key generation in temp directory
- Verify generated files exist
- Verify PEM format
- Verify key size

## Phase 9: Testing

### 9.1 Mock Generation

**Generate mocks for interfaces:**

```bash
go generate ./...
```

**Files with go:generate directives:**

- `internal/database/database.go`
- `internal/health/handler.go`

### 9.2 Unit Tests

**internal/config/config_test.go:**

- Configuration loading
- Default values
- Environment variable override
- Flag override

**internal/ulid/ulid_test.go:**

- ULID generation
- Uniqueness
- Validation
- Parsing

**internal/health/handler_test.go:**

```go
package health_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"

    "[package_name]/internal/health"
    mock_database "[package_name]/internal/database/mock"
)

func TestHealthHandler_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDB := mock_database.NewMockDB(ctrl)
    mockDB.EXPECT().Ping(gomock.Any()).Return(nil)

    handler := health.NewHandler(mockDB)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    handler.Health(c)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealthHandler_DatabaseError(t *testing.T) {
    // Test database connection failure
}
```

### 9.3 Integration Tests

**internal/database/database_test.go:**

```go
package database_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "[package_name]/internal/database"
    "[package_name]/internal/testutil"
)

func TestDatabase_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    dsn, cleanup := testutil.SetupPostgres(t)
    defer cleanup()

    db, err := database.New(dsn)
    require.NoError(t, err)
    defer db.Close()

    ctx := context.Background()
    err = db.Ping(ctx)
    assert.NoError(t, err)
}
```

### 9.4 CLI Tests

**cmd/[project_name]/cmd/serve_test.go:**

```go
package cmd

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestServeCommand(t *testing.T) {
    // Test command initialization
    assert.NotNil(t, serveCmd)
    assert.Equal(t, "serve", serveCmd.Use)
}

func TestServeCommand_Flags(t *testing.T) {
    // Test host flag
    // Test port flag
}
```

**cmd/[project_name]/cmd/migrate_test.go:**

```go
package cmd

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMigrateCommand(t *testing.T) {
    assert.NotNil(t, migrateCmd)
}

func TestMigrateCreateCommand_RequiresName(t *testing.T) {
    // Test that create requires exactly 1 argument
}

func TestMigrationFileFormat(t *testing.T) {
    // Test migration file naming format
}
```

**cmd/[project_name]/cmd/keygen_test.go:**

```go
package cmd

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestKeygenCommand(t *testing.T) {
    assert.NotNil(t, keygenCmd)
}

func TestKeygenCommand_MinimumBits(t *testing.T) {
    // Test validation: bits < 4096 should fail
}

func TestKeygenCommand_GeneratesKeys(t *testing.T) {
    tmpDir := t.TempDir()

    // Run keygen in temp directory
    // Verify private_key.pem exists
    // Verify public_key.pem exists
    // Verify PEM format
}
```

## Phase 10: Build Tools

### 10.1 `Justfile`

```just
# Default command
default:
    @just --list

# Build the application
build:
    go build -o bin/[project_name] ./cmd/[project_name]

# Run the server
run:
    go run ./cmd/[project_name] serve

# Create new migration
migrate-create name:
    go run ./cmd/[project_name] migrate create {{name}}

# Apply migrations (forward only)
migrate-up:
    go run ./cmd/[project_name] migrate up

# Show migration version
migrate-version:
    go run ./cmd/[project_name] migrate version

# Generate RSA keys
keygen:
    go run ./cmd/[project_name] keygen

# Run unit tests only (fast)
test:
    go test -v -short ./...

# Run all tests including integration tests
test-all:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

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
```

### 10.2 `docker-compose.yml`

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: [project_name]
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:8-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

### 10.3 `README.md`

````markdown
# [Project_name]

Flower marked platform for Central Asia for buyes and animals.

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Just (optional, for task automation)

## Quick Start

1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Start services: `just docker-up`
4. Run migrations: `just migrate-up`
5. Start server: `just run`

Server will start on <http://localhost:9000>

## Development

### Commands

- `just build` - Build the application
- `just run` - Run the server
- `just test` - Run unit tests
- `just test-all` - Run all tests including integration
- `just generate-mocks` - Generate mocks

### Database Migrations

- `just migrate-up` - Apply all pending migrations
- `just migrate-create <name>` - Create new migration
- `just migrate-version` - Show current version

### Testing

Run unit tests:

```bash
just test
```
````

Run integration tests:

```bash
just test-all
```

Generate coverage report:

```bash
just test-coverage
```

## Project Structure

```text
[project_name]/
├── cmd/[project_name]/  # CLI entrypoint
├── internal/            # Internal packages
│   ├── config/          # Configuration
│   ├── database/        # Database layer
│   ├── server/          # HTTP server
│   ├── health/          # Health check
│   ├── logger/          # Logging
│   └── ulid/            # ULID utilities
└── migrations/          # Database migrations
```

## License

Private

## Implementation Order

### Step 1: Foundation

1. Create directory structure
2. Install dependencies
3. Create `.env.example` and `.gitignore`

### Step 2: Configuration

4. Implement `internal/config/config.go` (default port: 9000)
5. Implement `internal/config/config_test.go`
6. Viper setup with defaults → env → flags priority

### Step 3: Logging

7. Implement `internal/logger/logger.go`
8. JSON output with slog

### Step 4: ULID

9. Implement `internal/ulid/ulid.go`
10. Implement `internal/ulid/ulid_test.go`

### Step 5: CLI Root

11. Create `cmd/[project_name]/main.go`
12. Create `cmd/[project_name]/main_test.go`
13. Viper + godotenv integration

### Step 6: Database

14. Implement `internal/database/database.go`
15. Implement `internal/database/database_test.go`
16. Add gomock generation directive
17. Create `migrations/000001_init_extensions.up.sql`

### Step 7: Migrate Command

19. Implement `cmd/[project_name]/cmd/migrate.go`
20. Implement `cmd/[project_name]/cmd/migrate_test.go`
21. Subcommands: up, create, version

### Step 8: Keygen Command

22. Implement `cmd/[project_name]/cmd/keygen.go`
23. Implement `cmd/[project_name]/cmd/keygen_test.go`
24. RSA key generation (4096+ bits)

### Step 9: Health Handler

25. Implement `internal/health/handler.go`
26. Implement `internal/health/handler_test.go`
27. Add gomock generation directive

### Step 10: HTTP Server

28. Implement `internal/server/server.go`
29. Implement `internal/server/server_test.go`
30. Register health endpoint
31. Add middleware (request ID with ULID, logger, recovery, CORS)

### Step 11: Serve Command

32. Implement `cmd/[project_name]/cmd/serve.go`
33. Implement `cmd/[project_name]/cmd/serve_test.go`
34. Graceful shutdown with signals

### Step 12: Testing Infrastructure

35. Implement `internal/testutil/postgres.go`
36. Implement `internal/testutil/redis.go`

### Step 13: Build Tools

37. Create `Justfile`
38. Create `docker-compose.yml`
39. Create `README.md`
