# ADR-001: Bootstrapping

| Field      | Value                   |
| ---------- | ----------------------- |
| Status     | Accepted                |
| Date       | 2026-02-21              |
| Supersedes | —                       |
| Components | observer, bootstrapping |

## Decision

Observer is a Go service using DDD + Clean Architecture with manual dependency injection (no DI frameworks). Gin handles HTTP routing; sqlx provides database access; configuration is read from environment variables with typed defaults.

## Project Layout

```
cmd/              # entrypoints (serve, migrate, keygen)
internal/
  domain/         # entities, errors, value objects — no framework imports
  usecase/        # business logic; one package per aggregate or feature area
  handler/        # thin HTTP adapters — bind request, call use case, return response
  middleware/     # HTTP middleware (auth, RBAC, rate limiting, CSRF, security headers)
  repository/     # all repository interfaces + implementations; mock/ subdir
  crypto/         # RSA keys, Argon2id hasher, JWT token generator
  storage/        # FileStorage interface + LocalStorage implementation
  config/         # environment-variable config with typed defaults
  server/         # HTTP server setup, route wiring, middleware chain
  app/            # DI container — manual wiring of all dependencies
  database/       # sqlx DB wrapper and interface
  logger/         # slog JSON logger + Gin request logging middleware
  ulid/           # ULID helpers (New, NewString)
  spa/            # embedded SPA filesystem (production build tag)
adr/              # architecture decision records
migrations/       # forward-only .up.sql migration files
api/swagger/      # generated Swagger docs
packages/         # frontend monorepo (observer-web)
```

## Architecture Rules

- Business logic lives exclusively in `internal/usecase/`. Handlers and repositories must not contain business rules.
- Handlers are thin: bind the request, call a use case, return the response.
- Domain packages define repository interfaces. `internal/repository/interfaces.go` collects all interfaces in a single file with one `go:generate` directive for mocks.
- Manual DI: `internal/app/container.go` wires every dependency explicitly. No reflect-based injection.
- `ulid.ULID` is used for all entity IDs; DTOs expose them as `string` (via `.String()`).

## Configuration

All config is read by `internal/config/config.Load()` from environment variables with hard-coded defaults so the binary starts with sane values for local development.

| Variable                    | Default                    | Purpose                               |
| --------------------------- | -------------------------- | ------------------------------------- |
| `DEV_MODE`                  | `false`                    | Disables CORS, CSRF, security headers |
| `SERVER_HOST`               | `localhost`                | Bind address                          |
| `SERVER_PORT`               | `9000`                     | Bind port                             |
| `SERVER_READ_TIMEOUT`       | `30s`                      | HTTP read timeout                     |
| `SERVER_WRITE_TIMEOUT`      | `30s`                      | HTTP write timeout                    |
| `DATABASE_DSN`              | `""`                       | PostgreSQL connection string          |
| `REDIS_URL`                 | `redis://localhost:6379/0` | Redis connection URL                  |
| `LOG_LEVEL`                 | `info`                     | slog log level                        |
| `JWT_PRIVATE_KEY_PATH`      | `keys/jwt_rsa`             | RSA private key path                  |
| `JWT_PUBLIC_KEY_PATH`       | `keys/jwt_rsa.pub`         | RSA public key path                   |
| `JWT_ACCESS_TTL`            | `15m`                      | Access token lifetime                 |
| `JWT_REFRESH_TTL`           | `168h`                     | Refresh token lifetime (7 days)       |
| `JWT_MFA_TEMP_TTL`          | `5m`                       | MFA temporary token lifetime          |
| `JWT_ISSUER`                | `observer`                 | JWT issuer claim                      |
| `SWAGGER_ENABLED`           | `false`                    | Enable `/swagger/*` UI                |
| `CORS_ORIGINS`              | `http://localhost:5173`    | Comma-separated allowed origins       |
| `COOKIE_DOMAIN`             | `""`                       | Cookie domain                         |
| `COOKIE_SECURE`             | `true`                     | Cookie Secure flag                    |
| `COOKIE_SAME_SITE`          | `lax`                      | Cookie SameSite policy                |
| `COOKIE_MAX_AGE`            | `2h`                       | Cookie max age                        |
| `RATE_LIMIT_LOGIN`          | `10`                       | Login attempts per minute             |
| `RATE_LIMIT_REGISTER`       | `5`                        | Register attempts per minute          |
| `STORAGE_PATH`              | `data/uploads`             | File storage root                     |
| `SENTRY_DSN`                | `""`                       | Sentry DSN (empty = disabled)         |
| `SENTRY_TRACES_SAMPLE_RATE` | `0.1`                      | Sentry trace sample rate              |

## Tech Stack

| Layer       | Library                                       |
| ----------- | --------------------------------------------- |
| HTTP        | `github.com/gin-gonic/gin`                    |
| Database    | `github.com/jmoiron/sqlx` (PostgreSQL)        |
| Migrations  | `github.com/golang-migrate/migrate`           |
| IDs         | `github.com/oklog/ulid/v2`                    |
| Hashing     | Argon2id (`golang.org/x/crypto`)              |
| JWT         | RSA-signed (`github.com/golang-jwt/jwt/v5`)   |
| Mocks       | `go.uber.org/mock/gomock`                     |
| Test assert | `github.com/stretchr/testify`                 |
| Integration | `github.com/testcontainers/testcontainers-go` |
| Logging     | `log/slog` (stdlib JSON handler)              |
| Error track | `github.com/getsentry/sentry-go`              |

## Build

The project uses a `Justfile` (not a Makefile):

```just
just test            # unit tests (no Docker)
just test-all        # unit + integration tests (requires Docker)
just generate-mocks  # regenerate all gomock mocks via go:generate
just generate-types  # regenerate TypeScript types from Go DTOs (tygo)
just migrate-up      # run forward migrations
just build-prod      # build frontend + Go binary with -tags production
```

The `production` build tag causes `internal/spa/embed.go` to embed the compiled frontend into the binary via `//go:embed all:dist`. Without that tag, `internal/spa/noembed.go` provides a no-op stub.

## Route Structure

All API routes are mounted under `/api`. The frontend SPA owns everything else via a `NoRoute` catch-all handler (production builds only).

```
GET  /health                       # load-balancer health check (no auth)
/api/auth/*                        # registration, login, refresh, MFA, profile
/api/my/*                          # authenticated user's own project list
/api/search?q=<query>              # global search (authenticated)
/api/admin/*                       # admin/staff/consultant management endpoints
/api/projects/:project_id/*        # project-scoped endpoints (role-gated)
/swagger/*                         # Swagger UI (when SWAGGER_ENABLED=true)
```
