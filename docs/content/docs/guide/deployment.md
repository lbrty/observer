---
title: Deployment
weight: 4
---

This guide walks you through putting Observer on a server where your team can use it. You don't need deep technical knowledge — if you can SSH into a server and run a few commands, you can do this.

## Docker (recommended)

This is the simplest path. Observer ships as a single Docker image with the web interface already baked in — there's nothing extra to install or configure on the frontend side.

### What you need

- A server with Docker and Docker Compose installed
- A domain name pointed at your server (for HTTPS)

### Step 1: Generate signing keys

Observer uses RSA keys to sign login tokens. Run these commands on your server to create them:

```bash
mkdir -p keys
openssl genrsa -out keys/jwt_rsa 4096
openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub
```

Keep these keys safe. If you lose them, everyone will need to log in again.

### Step 2: Configure your environment

Copy the example environment file and edit it for your setup:

```bash
cp .env.example .env
```

The most important variables:

| Variable               | What it does                                      | Default                    |
| ---------------------- | ------------------------------------------------- | -------------------------- |
| `DATABASE_DSN`         | How Observer connects to PostgreSQL               | _(must be set)_            |
| `REDIS_URL`            | How Observer connects to Redis                    | `redis://localhost:6379/0` |
| `JWT_PRIVATE_KEY_PATH` | Where you put the private key from Step 1         | `keys/jwt_rsa`             |
| `JWT_PUBLIC_KEY_PATH`  | Where you put the public key from Step 1          | `keys/jwt_rsa.pub`         |
| `CORS_ORIGINS`         | Your domain (e.g. `https://observer.yourorg.org`) | `http://localhost:5173`    |
| `COOKIE_SECURE`        | Set to `true` when using HTTPS (you should)       | `true`                     |
| `SERVER_HOST`          | Which address to listen on                        | `localhost`                |
| `SERVER_PORT`          | Which port to listen on                           | `9000`                     |

See [Environment Variables](/docs/developers/reference/variables/) for the full list.

### Step 3: Start Observer

```bash
docker compose up -d
```

This starts PostgreSQL, Redis, and Observer. The database schema is created automatically on first launch — no manual migration step needed.

### Step 4: Verify it's running

```bash
curl http://localhost:9000/health
```

You should see:

```json
{ "status": "healthy", "database": "connected", "timestamp": "..." }
```

If you see this, Observer is ready. Open your domain in a browser to access the web interface.

## Without Docker (VPS / bare metal)

If you prefer to run Observer directly, build the binary:

```bash
CGO_ENABLED=0 go build -tags production -ldflags="-s -w" -o observer ./cmd/observer
```

The `-tags production` flag embeds the web interface into the binary. You get a single file you can copy anywhere.

Run it:

```bash
./observer serve --host 0.0.0.0
```

You'll need PostgreSQL and Redis running separately. Point `DATABASE_DSN` and `REDIS_URL` to them.

## Setting up HTTPS

You should always run Observer behind a reverse proxy that handles HTTPS. This keeps login credentials and personal data encrypted in transit.

[Caddy](https://caddyserver.com/) is the easiest option — it handles certificates automatically:

```
observer.yourorg.org {
    reverse_proxy localhost:9000
}
```

If you use Nginx or another proxy, make sure to set:

- `COOKIE_SECURE=true` in your environment
- `CORS_ORIGINS` to your actual domain (e.g. `https://observer.yourorg.org`)
