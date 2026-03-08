# Observer CLI Reference

## Overview

Observer provides a CLI for managing the server, database migrations, key generation, user administration, and development utilities.

## Installation

```bash
# Build from source
just build

# Or install directly
go install github.com/lbrty/observer/cmd/observer@latest
```

## Commands

### serve

Start the HTTP server.

```bash
observer serve [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--host` | string | `localhost` | Server host (overrides `SERVER_HOST` env) |
| `--port` | int | `9000` | Server port (overrides `SERVER_PORT` env) |

**Examples:**

```bash
# Start with defaults (localhost:9000)
observer serve

# Custom host and port
observer serve --host 0.0.0.0 --port 8080

# With environment configuration
DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve
```

In production builds, embedded migrations are applied automatically on startup. The server shuts down gracefully on SIGINT/SIGTERM with a 30-second timeout.

---

### migrate

Database migration management.

#### migrate up

Apply all pending migrations.

```bash
observer migrate up [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--path` | string | `migrations` | Path to migrations directory |

```bash
# Apply all pending migrations
observer migrate up

# Use a custom migrations directory
observer migrate up --path ./db/migrations
```

#### migrate create

Create a new forward-only migration file.

```bash
observer migrate create [name] [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--path` | string | `migrations` | Path to migrations directory |
| `--seq` | uint | auto | Explicit sequence number |

```bash
# Create a migration (auto-numbered)
observer migrate create add_audit_log

# Create with explicit sequence number
observer migrate create add_audit_log --seq 25
```

#### migrate version

Show the current migration version.

```bash
observer migrate version [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--path` | string | `migrations` | Path to migrations directory |

---

### keygen

Generate an RSA key pair for JWT signing.

```bash
observer keygen [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--bits` | int | `4096` | RSA key size (minimum 4096) |
| `--output` | string | `.` | Output directory for key files |

**Examples:**

```bash
# Generate keys in the current directory
observer keygen

# Generate 8192-bit keys in the keys/ directory
observer keygen --bits 8192 --output keys
```

Output files: `private_key.pem` (0600) and `public_key.pem` (0644).

---

### create-admin

Create a platform administrator account.

```bash
observer create-admin [flags]
```

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--email` | string | yes | Admin email |
| `--password` | string | yes | Admin password (min 8 chars) |
| `--first-name` | string | no | First name |
| `--last-name` | string | no | Last name |
| `--phone` | string | no | Phone number |

**Examples:**

```bash
# Create an admin with required fields
observer create-admin --email admin@example.com --password "s3cure-p4ss"

# With optional profile fields
observer create-admin \
  --email admin@example.com \
  --password "s3cure-p4ss" \
  --first-name Admin \
  --last-name User \
  --phone "+1234567890"
```

Connects to the database, hashes the password with Argon2id, and inserts the user with admin role (verified + active). Rejects duplicate emails and phone numbers.

---

### seed

Seed the database with realistic mock data for development.

```bash
observer seed [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--people` | int | `50` | Number of people per project |
| `--projects` | int | `2` | Number of projects |
| `--seed` | int64 | `0` | Random seed (0 = random) |

**Examples:**

```bash
# Seed with defaults (2 projects, 50 people each)
observer seed

# Custom counts
observer seed --projects 5 --people 200

# Reproducible seed
observer seed --seed 42
```

**WARNING:** This command truncates ALL tables before inserting data. Do not run against a production database.

Creates reference data (countries, states, places, offices, categories), users with known passwords (`password`), projects with permissions, and populates people with support records, migration records, notes, pets, and households.

---

### setup

Run first-time project setup interactively.

```bash
observer setup
```

This command:

1. Creates a `.env` file with sensible defaults (prompts before overwriting)
2. Creates required directories (`keys/`, `data/uploads/`)
3. Generates a 4096-bit RSA key pair for JWT signing
4. Prints next-steps instructions

**Example output:**

```
Created .env with default configuration.
Created directory: keys
Created directory: data/uploads
Generating 4096-bit RSA key pair...
Private key written to: keys/private_key.pem
Public key written to: keys/public_key.pem

Setup complete!

Next steps:
  1. Start Postgres and Redis:
     docker compose up -d

  2. Run database migrations:
     observer migrate up

  3. Create an admin user:
     observer create-admin --email admin@example.com --password "your-password"

  4. Start the server:
     observer serve
```

---

## Environment Variables

All configuration is read from environment variables. A `.env` file in the working directory is loaded automatically via [godotenv](https://github.com/joho/godotenv).

### Server

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `localhost` | Bind address |
| `SERVER_PORT` | `9000` | Listen port |
| `SERVER_READ_TIMEOUT` | `30s` | HTTP read timeout |
| `SERVER_WRITE_TIMEOUT` | `30s` | HTTP write timeout |

### Database

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_DSN` | _(none)_ | PostgreSQL connection string |

### Redis

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_URL` | `redis://localhost:6379/0` | Redis connection URL |

### JWT

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa` | Path to RSA private key |
| `JWT_PUBLIC_KEY_PATH` | `keys/jwt_rsa.pub` | Path to RSA public key |
| `JWT_ACCESS_TTL` | `15m` | Access token lifetime |
| `JWT_REFRESH_TTL` | `168h` | Refresh token lifetime (7 days) |
| `JWT_MFA_TEMP_TTL` | `5m` | MFA temporary token lifetime |
| `JWT_ISSUER` | `observer` | JWT issuer claim |

### CORS

| Variable | Default | Description |
|----------|---------|-------------|
| `CORS_ORIGINS` | `http://localhost:5173` | Comma-separated allowed origins |

### Cookies

| Variable | Default | Description |
|----------|---------|-------------|
| `COOKIE_DOMAIN` | _(empty)_ | Cookie domain |
| `COOKIE_SECURE` | `true` | Set Secure flag on cookies |
| `COOKIE_SAME_SITE` | `lax` | SameSite policy (`lax`, `strict`, `none`) |
| `COOKIE_MAX_AGE` | `2h` | Cookie max age |

### Rate Limiting

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_LOGIN` | `10` | Login attempts per window |
| `RATE_LIMIT_REGISTER` | `5` | Registration attempts per window |

### Storage

| Variable | Default | Description |
|----------|---------|-------------|
| `STORAGE_PATH` | `data/uploads` | Root directory for file uploads |

### Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |

### Swagger

| Variable | Default | Description |
|----------|---------|-------------|
| `SWAGGER_ENABLED` | `false` | Enable Swagger UI at `/swagger/` |

### Sentry

| Variable | Default | Description |
|----------|---------|-------------|
| `SENTRY_DSN` | _(empty)_ | Sentry DSN (empty = disabled) |
| `SENTRY_TRACES_SAMPLE_RATE` | `0.1` | Sentry traces sample rate |

---

## Common Workflows

### First-time setup

```bash
# 1. Run interactive setup (creates .env, keys, directories)
observer setup

# 2. Start Postgres and Redis
docker compose up -d

# 3. Run migrations
observer migrate up

# 4. Create an admin user
observer create-admin --email admin@example.com --password "your-password"

# 5. Start the server
observer serve
```

### Adding a new migration

```bash
# Create the migration file
observer migrate create add_audit_log

# Edit the generated SQL file
$EDITOR migrations/000025_add_audit_log.up.sql

# Apply it
observer migrate up
```

### Seeding development data

```bash
# Make sure migrations are applied first
observer migrate up

# Seed with defaults
observer seed

# Login with: admin@example.com / password
```

### Generating new JWT keys

```bash
# Generate new keys
observer keygen --output keys

# Update .env if paths differ from defaults
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem

# Restart the server — existing tokens will be invalidated
observer serve
```
