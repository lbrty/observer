---
title: Getting Started
weight: 2
---

## Try it in 5 minutes

You don't need a server, a hosting provider, or an IT department. If you have a laptop with Docker installed, you can see Observer running right now.

```bash
git clone https://github.com/lbrty/observer.git
cd observer
cp .env.example .env
just generate-keys
just docker-up
just run
```

Open `http://localhost:9000/health` — if you see `"status":"healthy"`, the backend is running.

Then start the web interface:

```bash
just web-dev
```

Open `http://localhost:5173` — you're looking at Observer.

## What you just started

- A **Go backend** serving the API on port 9000
- A **PostgreSQL database** with 24 tables for people, households, support records, migration history, documents, and pets
- A **React frontend** with project management, role-based access, and 39 built-in report types
- **JWT authentication** with automatic token rotation

All of this runs on a single machine. In production, it compiles down to one binary.

## What you'll need for a real deployment

| Requirement | Why |
| --- | --- |
| A VPS or on-premise server | Observer is self-hosted — your data stays on your infrastructure |
| PostgreSQL database | The only external dependency |
| 30 minutes of sysadmin time | `docker compose up` on a server with your domain pointed at it |

No SaaS subscription. No per-user fees. No vendor lock-in. You own the data and the deployment.

See [Deployment](/docs/guide/deployment/) for the full production setup guide.

## Prerequisites for development

| Tool | Version | Install |
| --- | --- | --- |
| Go | 1.25.* | https://go.dev/dl/ |
| Bun | latest | https://bun.sh/ |
| Docker + Compose | latest | https://docs.docker.com/get-docker/ |
| Just | latest | https://github.com/casey/just#installation |

## Step by step

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

The defaults work out of the box with the provided `docker-compose.yml`.

### 3. Generate RSA keys

```bash
just generate-keys
```

Creates `keys/jwt_rsa` and `keys/jwt_rsa.pub` for JWT signing.

### 4. Start services and run

```bash
just docker-up    # starts PostgreSQL and Redis
just run          # starts the backend on :9000 (runs migrations automatically)
just web-dev      # starts the frontend on :5173
```

## Troubleshooting

**Port 5432 already in use** — A local Postgres instance may be running. Stop it or change the port mapping in `docker-compose.yml`.

**"no such file or directory" for key paths** — Run `just generate-keys` first. The `keys/` directory is gitignored.

**Migration fails with "connection refused"** — Docker services may not be ready yet. Wait a few seconds after `just docker-up`.
