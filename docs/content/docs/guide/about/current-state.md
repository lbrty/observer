---
title: Current State
weight: 3
---

## What Observer is

A self-hosted IDP case management platform. It tracks displaced people, their geographic movement, support records, households, documents, and pets — scoped per project with role-based access control and built-in reporting.

## Stack

| Concern | Choice |
| --- | --- |
| Backend | Go 1.25, Gin, sqlx, PostgreSQL |
| Frontend | React 19, Vite 6, TanStack Router/Query, Tailwind v4 |
| Auth | RS256 JWT (HttpOnly cookies), Argon2id passwords |
| IDs | ULID (sortable, human-readable) |
| Migrations | golang-migrate, forward-only |
| Testing | testify + gomock + testcontainers-go |
| CLI | Cobra |
| Build | Justfile, Bun, Docker multistage |

## Project structure

```
cmd/observer/              # CLI (serve, migrate, keygen)
internal/
  app/                     # Manual DI container
  config/                  # Env-based configuration
  crypto/                  # RSA JWT, Argon2id hashing
  database/                # sqlx wrapper
  domain/                  # Entities + repository interfaces
    auth/ user/ project/ reference/
    person/ tag/ support/ migration/
    household/ note/ document/ pet/
  handler/                 # Thin HTTP adapters
  middleware/              # JWT auth + project RBAC
  repository/              # Repository implementations
  usecase/                 # Business logic
  server/                  # Gin setup, routes, embedded SPA
  spa/                     # go:embed frontend for production
packages/
  observer-web/            # React frontend
migrations/                # Forward-only .up.sql files
```

## Database schema

24 migrations, forward-only. Reference data first, then auth, then projects, then core domain.

| Table | Purpose |
| --- | --- |
| countries, states, places | Geographic hierarchy |
| offices | Service delivery locations |
| categories | Vulnerability categories |
| users, credentials, mfa_configs | Platform users and auth |
| sessions, verification_tokens | Auth lifecycle |
| projects, project_permissions | Multi-project scoping |
| tags, person_tags | Project-scoped tagging |
| people | Individual records |
| person_categories | Multi-code vulnerability (junction) |
| households, household_members | Family units with typed relationships |
| migration_records | Movement history with causality |
| support_records | Service delivery + referral tracking |
| person_notes | Case notes |
| documents | Document metadata |
| pets | Companion animals |

## Authentication

```
POST /auth/register  → hash password → create user + credentials
POST /auth/login     → verify password → check MFA → create session → set cookies
POST /auth/refresh   → rotate refresh token → set new cookies
POST /auth/logout    → delete session → clear cookies
```

- Access token: RS256 JWT, 15 min
- Refresh token: ULID, 7 days, rotated on each refresh
- Transport: HttpOnly cookies (`access_token` on `/`, `refresh_token` on `/auth`)

## Authorization

Two levels:

1. **Platform role** (admin / staff / consultant / guest) — set on the user record
2. **Project role** (owner / manager / consultant / viewer) — set per user per project, plus three sensitivity flags: `can_view_contact`, `can_view_personal`, `can_view_documents`

Action permissions (read/create/update/delete) derive from the project role rank. Sensitivity flags are independent boolean grants.

## Deployment

Production builds embed the React frontend into the Go binary via `go:embed`. The final Docker image is built from `scratch` — a single static binary + SSL certs + migrations.

```bash
docker compose up    # postgres + redis + observer on :9000
```
