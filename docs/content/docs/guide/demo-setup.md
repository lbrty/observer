---
title: Demo Setup
weight: 3
---

Try Observer with a fully populated database. The `seed` command fills the database with realistic test data — people, projects, support records, migration history, households, and more — so you can explore every feature without entering data by hand.

## Quick start

```bash
git clone https://github.com/lbrty/observer.git
cd observer
just build
./observer setup
docker compose up -d
./observer migrate up
./observer seed
./observer serve
```

Open `http://localhost:9000/health` — you should see `{"status":"healthy"}`.

Start the web interface:

```bash
cd packages/observer-web
bun install
bun run dev
```

Open `http://localhost:5173` and log in with `admin@example.com` / `password`.

## What the seed command creates

| What                  | Details                                              |
| --------------------- | ---------------------------------------------------- |
| **Reference data**    | Countries, states, places, offices, categories       |
| **Users**             | Admin + staff accounts, all with password `password` |
| **Projects**          | 2 projects (configurable with `--projects`)          |
| **People**            | 50 per project (configurable with `--people`)        |
| **Support records**   | Consultation records linked to people                |
| **Migration records** | Movement history with origin/destination places      |
| **Households**        | Family groupings with members                        |
| **Notes**             | Case notes attached to people                        |
| **Pets**              | Pet records with tags                                |
| **Tags**              | Project-scoped labels for categorization             |

### Default credentials

| Email               | Password   | Role  |
| ------------------- | ---------- | ----- |
| `admin@example.com` | `password` | Admin |

### Custom seed options

```bash
# More projects and people
./observer seed --projects 5 --people 200

# Reproducible data (same seed = same output)
./observer seed --seed 42
```

{{% callout type="warning" %}}
The seed command **truncates ALL tables** before inserting data. Never run it against a production database.
{{% /callout %}}

## Step-by-step breakdown

### 1. Build

```bash
just build
```

Compiles the `observer` binary.

### 2. Setup

```bash
./observer setup
```

Creates `.env` with sensible defaults, `keys/` and `data/uploads/` directories, and generates a 4096-bit RSA key pair for JWT signing. If `.env` already exists, it will ask before overwriting.

### 3. Start services

```bash
docker compose up -d
```

Starts PostgreSQL and Redis in the background.

### 4. Migrate

```bash
./observer migrate up
```

Applies all database migrations.

### 5. Seed

```bash
./observer seed
```

Fills the database with demo data.

### 6. Run

```bash
./observer serve
```

Starts the API server on `http://localhost:9000`.

### 7. Frontend (optional)

```bash
cd packages/observer-web
bun install
bun run dev
```

Starts the web interface on `http://localhost:5173`.

## Reset

To wipe everything and start fresh:

```bash
docker compose down -v
```

Then re-run from step 3.
