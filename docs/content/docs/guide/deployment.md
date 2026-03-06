---
title: Deployment
weight: 2
---

## Docker (recommended)

Observer ships as a single Docker image. The frontend is embedded into the Go binary at build time.

### Prerequisites

- Docker + Docker Compose
- RSA key pair for JWT signing

### Generate keys

```bash
mkdir -p keys
openssl genrsa -out keys/jwt_rsa 4096
openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub
```

### Start

```bash
docker compose up -d
```

This starts PostgreSQL, Redis, and Observer on port 9000. The web UI is served from the binary — no separate frontend deployment needed.

### Environment

Key variables (set in `docker-compose.yml` or via env):

| Variable | Default | Description |
| --- | --- | --- |
| `DATABASE_DSN` | — | PostgreSQL connection string |
| `REDIS_URL` | `redis://localhost:6379/0` | Redis connection string |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa` | Path to RSA private key |
| `JWT_PUBLIC_KEY_PATH` | `keys/jwt_rsa.pub` | Path to RSA public key |
| `CORS_ORIGINS` | `http://localhost:5173` | Allowed CORS origins |
| `COOKIE_SECURE` | `true` | Set `false` for non-HTTPS |
| `SERVER_HOST` | `localhost` | Bind address |
| `SERVER_PORT` | `9000` | Listen port |

See [Environment Variables](/docs/developers/reference/variables/) for the full list.

### Migrations

Migrations run automatically on startup. The binary includes all migration files.

### Health check

```bash
curl http://localhost:9000/health
# {"status":"healthy","database":"connected","timestamp":"..."}
```

## VPS / bare metal

Build the binary:

```bash
CGO_ENABLED=0 go build -tags production -ldflags="-s -w" -o observer ./cmd/observer
```

The `-tags production` flag embeds the frontend. Run:

```bash
./observer serve --host 0.0.0.0
```

You'll need PostgreSQL and Redis running separately. Set `DATABASE_DSN` and `REDIS_URL` accordingly.

## Reverse proxy

Put Nginx, Caddy, or similar in front for TLS termination. Set `COOKIE_SECURE=true` and `CORS_ORIGINS` to your domain.

Caddy example:

```
observer.example.org {
    reverse_proxy localhost:9000
}
```
