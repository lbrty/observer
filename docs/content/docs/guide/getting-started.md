---
title: Getting Started
weight: 2
---

## See it running in 5 minutes

You don't need a server, a hosting provider, or an IT department. If you have a laptop with Docker installed, you can see Observer running right now.

```bash
git clone https://github.com/lbrty/observer.git
cd observer
cp .env.example .env
just generate-keys
just docker-up
just run
```

Open `http://localhost:9000/health` in your browser. If you see `"status":"healthy"`, the backend is running.

Then start the web interface:

```bash
just web-dev
```

Open `http://localhost:5173` — you're looking at Observer.

## What you just started

- A **backend** serving the API — handles authentication, data storage, and reports
- A **database** with tables for people, households, support records, migration history, documents, and pets
- A **web interface** with project management, role-based access, and built-in reporting
- **Automatic login security** — tokens rotate on every session refresh

All of this runs on a single machine. In production, it compiles down to one file you can copy to any server.

## Ready to deploy for real?

To move from "trying it out" to "my team uses this every day," you need:

| What | Why |
| --- | --- |
| A server (VPS or on-premise) | Observer is self-hosted — your data never leaves your infrastructure |
| PostgreSQL | The only external service Observer needs |
| About 30 minutes | Run `docker compose up` on a server with your domain pointed at it |

No subscription. No per-user fees. No vendor lock-in. You own the data and the deployment.

See [Deployment](/docs/guide/deployment/) for the step-by-step production setup.

## For developers: local setup

If you want to work on Observer itself, you'll need these tools installed:

| Tool | Version | Install |
| --- | --- | --- |
| Go | 1.25.* | https://go.dev/dl/ |
| Bun | latest | https://bun.sh/ |
| Docker + Compose | latest | https://docs.docker.com/get-docker/ |
| Just | latest | https://github.com/casey/just#installation |

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

The defaults work out of the box with the provided `docker-compose.yml`. No editing needed.

### 3. Generate signing keys

```bash
just generate-keys
```

This creates a key pair that Observer uses to sign login tokens. The `keys/` directory is gitignored — each developer generates their own.

### 4. Start everything

```bash
just docker-up    # starts PostgreSQL and Redis
just run          # starts the backend on :9000 (runs migrations automatically)
just web-dev      # starts the frontend on :5173
```

## Something not working?

**Port 5432 already in use** — You probably have a local PostgreSQL running. Stop it, or change the port in `docker-compose.yml`.

**"no such file or directory" for key paths** — You need to run `just generate-keys` first.

**Migration fails with "connection refused"** — The database container might not be ready yet. Wait a few seconds after `just docker-up` and try again.
