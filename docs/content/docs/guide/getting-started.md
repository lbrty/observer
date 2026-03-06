---
title: Getting Started
weight: 2
---

## Prerequisites

| Tool             | Version | Install                                                                 |
| ---------------- | ------- | ----------------------------------------------------------------------- |
| Go               | 1.25.\* | https://go.dev/dl/                                                      |
| Bun              | latest  | https://bun.sh/                                                         |
| Docker + Compose | latest  | https://docs.docker.com/get-docker/                                     |
| Just             | latest  | https://github.com/casey/just#installation                              |
| golangci-lint    | latest  | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |

Optional for OpenAPI generation:

| Tool | Install                                             |
| ---- | --------------------------------------------------- |
| swag | `go install github.com/swaggo/swag/cmd/swag@latest` |

## Setup

### 1. Clone and install dependencies

```bash
git clone https://github.com/lbrty/observer.git
cd observer
go mod download
bun install
```

### 2. Configure environment

```bash
cp .env.example .env
```

The defaults work out of the box with the provided `docker-compose.yml`. Edit `.env` only if you need non-default ports or credentials.

### 3. Generate RSA keys

The server needs an RSA key pair for JWT signing. Use the OpenSSL method — it matches the default config paths:

```bash
just generate-keys
```

This creates `keys/jwt_rsa` (private) and `keys/jwt_rsa.pub` (public) with correct permissions.

Alternatively, use the built-in Go command (requires adjusting env vars):

```bash
mkdir -p keys
go run ./cmd/observer keygen --output keys
```

This outputs `keys/private_key.pem` and `keys/public_key.pem`, so you'd update `.env`:

```
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem
```

### 4. Start Postgres and Redis

```bash
just docker-up
```

Verify both services are healthy:

```bash
docker-compose ps
```

You should see `postgres` and `redis` with status `healthy`.

### 5. Run database migrations

```bash
just migrate-up
```

Check current version:

```bash
just migrate-version
```

### 6. Start the server

```bash
just run
```

The server starts on `http://localhost:9000` with Swagger UI enabled at `http://localhost:9000/swagger/index.html`.

### 7. Verify the backend

```bash
curl http://localhost:9000/health
```

Expected response:

```json
{ "status": "healthy", "database": "connected", "timestamp": "..." }
```

### 8. Start the frontend

```bash
just web-dev
```

Opens at `http://localhost:5173`. The frontend proxies API requests to the backend at `:9000` via cookies (CORS is pre-configured for `localhost:5173`).

See [docs/frontend.md](frontend.md) for full frontend documentation.

## Running Tests

```bash
just test          # unit tests only — fast, no Docker needed
just test-all      # all tests including integration — Docker must be running
just test-coverage # generate HTML coverage report (opens coverage.html)
just test-race     # run with Go race detector
```

## Common Tasks

| Task                  | Command                      |
| --------------------- | ---------------------------- |
| Build binary          | `just build`                 |
| Format code           | `just fmt`                   |
| Lint                  | `just lint`                  |
| Tidy modules          | `just tidy`                  |
| Regenerate mocks      | `just generate-mocks`        |
| Generate OpenAPI spec | `just openapi`               |
| Create new migration  | `just migrate-create <name>` |
| Stop Docker services  | `just docker-down`           |
| Frontend dev server   | `just web-dev`               |
| Frontend build        | `just web-build`             |
| Frontend dependencies | `just web-install`           |
| List all commands     | `just`                       |

## Troubleshooting

**Port 5432 already in use** — A local Postgres instance may be running. Stop it or change the port mapping in `docker-compose.yml` and update `DATABASE_DSN` in `.env`.

**"no such file or directory" for key paths** — Run `just generate-keys` first. The `keys/` directory is gitignored and must be created locally.

**Migration fails with "connection refused"** — Docker services may not be ready yet. Wait a few seconds after `just docker-up` or check `docker-compose ps` for health status.

**Tests fail with "short mode"** — Integration tests are skipped with `just test`. Use `just test-all` (requires Docker) to run the full suite.
