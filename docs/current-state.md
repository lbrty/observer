# Observer — Current State

This document summarises the current architecture, all key decisions made, and the rationale behind simplifications applied during the initial design phase.

---

## What Observer Is

Observer is an IDP (Internally Displaced Persons) case management platform. It tracks displaced people, their geographic movement, case workers, support records, and associated assets (documents, pets) within project contexts.

---

## Project Structure

```text
observer/
├── cmd/observer/              # CLI entry point (serve, migrate, keygen)
├── internal/
│   ├── app/                   # DI container
│   ├── application/auth/      # Use cases: register, login, refresh, logout
│   ├── config/                # Environment-based configuration
│   ├── database/              # DB interface + sqlx implementation
│   ├── domain/
│   │   ├── auth/              # Session entity + repository interface
│   │   └── user/              # User, Credentials, MFAConfig, VerificationToken
│   ├── health/                # Health check handler
│   ├── infrastructure/
│   │   ├── auth/              # RSA JWT, Argon2id password hashing
│   │   └── persistence/postgres/  # Repository implementations
│   ├── interfaces/http/
│   │   ├── handlers/          # Auth HTTP handlers
│   │   └── middleware/        # JWT authentication middleware
│   ├── logger/                # slog JSON logger + Gin middleware
│   ├── server/                # Gin HTTP server setup
│   ├── testutil/              # Testcontainers helpers (Postgres, Redis)
│   └── ulid/                  # Thread-safe ULID generator
├── migrations/                # Forward-only .up.sql files (000001–000021)
├── docs/adr/                  # Architecture Decision Records
├── Justfile                   # Developer task runner
└── docker-compose.yml         # Postgres + Redis for local dev
```

---

## ADR Index

| ADR | Title | Status |
|-----|-------|--------|
| [ADR-001](adr/001-bootstrapping.md) | Basic project structure, CLI, database, health check, server | Accepted |
| [ADR-002](adr/002-users-and-auth.md) | User domain, JWT auth with RSA signing, Argon2id passwords | Accepted |
| [ADR-003](adr/003-main-schema.md) | Main database schema (geography, projects, people, support) | Accepted |
| [ADR-004](adr/004-forward-only-migrations.md) | Forward-only migrations, no down files | Accepted |
| [ADR-005](adr/005-reports.md) | Reports specification (39 report requirements from IDP archive) | Accepted |

---

## Migration Sequence

All 22 migrations are forward-only `.up.sql` files. No `.down.sql` files exist.

| Number | Migration |
|--------|-----------|
| 000001 | Init PostgreSQL extensions (citext, pg_trgm, unaccent, pgcrypto) |
| 000002 | users |
| 000003 | credentials |
| 000004 | mfa_configs |
| 000005 | sessions |
| 000006 | verification_tokens |
| 000007 | countries |
| 000008 | states |
| 000009 | places |
| 000010 | offices |
| 000011 | categories |
| 000012 | tags (project-scoped) |
| 000013 | projects |
| 000014 | project_permissions |
| 000015 | people |
| 000016 | person_tags |
| 000017 | migration_records |
| 000018 | support_records |
| 000019 | pets |
| 000020 | documents |
| 000021 | add office_id column to users |
| 000022 | person_notes |

---

## Key Simplifications and Why

### 1. Forward-Only Migrations (ADR-004)

**What changed:** Removed all `.down.sql` files. The `migrate create` command now generates only a single `.up.sql`.

**Why:** Rollback SQL is almost never tested, frequently incorrect, and destroys data written by the up migration. Production databases are never rolled back — issues are fixed by applying a new forward migration. Maintaining two files per change doubled the surface area with near-zero practical benefit.

---

### 2. Eliminated `schema_health` Table

**What changed:** Deleted `migrations/000002_init_schema.up.sql` which created a `schema_health` table with a ULID primary key and `checked_at` timestamp.

**Why:** The table was never referenced anywhere in the codebase. The health handler used `db.Ping()` instead. A dead table with no consumer is noise.

---

### 3. Simplified Health Endpoint

**What changed:** `GET /health` now returns only `{"status":"ok"}` (200) or `{"status":"not ok"}` (503).

**Before:**

```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-02-01T12:00:00Z"
}
```

**After:**

```json
{ "status": "ok" }
```

**Why:** Health check endpoints are consumed by load balancers and orchestrators (Kubernetes, ECS) that only check the HTTP status code. The extra fields added no operational value and the timestamp was redundant with HTTP response headers.

---

### 4. Fixed User Role Enum

**What changed:** `users.role` changed from `('user', 'seller', 'admin')` to `('admin', 'staff', 'consultant', 'guest')`.

**Why:** The original values were carried over from a generic template and had no meaning in an IDP platform. The new roles reflect the actual actors: administrators, staff who manage projects, consultants who work cases, and read-only guests.

---

### 5. Role-Based Project Permissions

**What changed:** `project_permissions` replaced 7 boolean columns (`can_read`, `can_create`, `can_update`, `can_delete`, `can_view_contact`, `can_view_personal`, `can_view_documents`) with a `project_role` enum + 3 sensitivity flags.

**Before:** 7 independent booleans → 128 possible combinations, most meaningless.

**After:**

```sql
role               project_role NOT NULL DEFAULT 'viewer',  -- owner|manager|consultant|viewer
can_view_contact   BOOLEAN      NOT NULL DEFAULT FALSE,
can_view_personal  BOOLEAN      NOT NULL DEFAULT FALSE,
can_view_documents BOOLEAN      NOT NULL DEFAULT FALSE,
```

**Why:** Action permissions (read/create/update/delete) naturally follow from a role. Only data sensitivity access is genuinely per-user and kept as explicit flags. This halved the table width, eliminated invalid flag combinations, and moved authorization logic into a single role check in middleware.

---

### 6. Dropped `offices.state_id`

**What changed:** Removed the `state_id` column from `offices`, keeping only `place_id`.

**Why:** `place_id` already points to a row in `places` which has its own `state_id`. Storing `state_id` directly on `offices` created a redundant column that could drift out of sync. The state is derivable via a single join: `offices → places → states`.

---

### 7. Renamed `from_place_id` → `origin_place_id` on `people`

**What changed:** The column tracking a person's home location on the `people` table was renamed.

**Why:** `from_place_id` was ambiguous — it could mean "where they moved from last" (a movement event) or "their permanent origin" (a biographical fact). The `migration_records` table already tracks movement events with its own `from_place_id`. On `people`, `origin_place_id` unambiguously means hometown/origin.

---

### 8. Tags Normalized to Junction Table

**What changed:** Replaced `TEXT[]` tags on the `people` table with a `tags` table and `person_tags` junction table.

**Why:** Arrays prevent:

- Efficient reverse lookup ("find all people with tag X" requires a scan)
- Atomic rename of a tag across all records
- Referential integrity (any string is valid)

The normalized model handles all three and has no meaningful query overhead with a proper index.

---

### 9. Fixed `support_records.owner_id` (archive bug)

**What changed:** `owner_id` changed from a bare `TEXT` column to `TEXT REFERENCES users (id) ON DELETE SET NULL`.

**Why:** In the original Python migration, `owner_id` was declared as `Column(Text)` with no foreign key. This allowed arbitrary strings to be stored, breaking referential integrity. The corrected version enforces that `owner_id` must reference a real user.

---

### 10. `pets.owner_id` References `people`, Not `users`

**What changed:** The FK target for `pets.owner_id` changed from `users` to `people`.

**Why:** In an IDP context the owner of a pet is a displaced person being tracked in the system, not a platform operator. The archive's `users` FK was semantically wrong.

---

### 11. Added `first_name`, `last_name`, `office_id` to `users`

**What changed:** Platform users now have a display name and office affiliation.

**Why:** Staff and consultants need to be identified by name in case files and filtered by office. The original `users` table only had `email` + `phone` + `role`, which is sufficient for authentication but not for display or organizational queries. `office_id` is added in a separate migration (000021) to avoid a circular dependency with the `offices` table (000010).

---

## Technology Choices

| Concern | Choice | Reason |
|---|---|---|
| Language | Go 1.22+ | Strong concurrency, fast binaries, excellent stdlib |
| HTTP | Gin | Low overhead, familiar middleware pattern |
| Database | PostgreSQL | JSONB, CITEXT, pg_trgm, strong FK enforcement |
| Migrations | golang-migrate (file source) | Simple, no magic, forward-only per ADR-004 |
| Auth tokens | RS256 JWT | Asymmetric — public key can be shared with other services without exposing signing key |
| Password hashing | Argon2id | Winner of Password Hashing Competition, memory-hard |
| IDs | ULID (TEXT) | Sortable by time, no UUID extension dependency, human-readable |
| CLI | Cobra + godotenv | Standard Go CLI pattern |
| Mocks | go.uber.org/mock (mockgen) | Interface-based, generated, type-safe |
| Tests | testify + testcontainers | Real database in integration tests, no manual setup |
| Logging | slog (stdlib) | Structured JSON, no external dependency |

---

## Authentication Flow

```text
POST /auth/register  → validate → check uniqueness → hash password → create user + credentials
POST /auth/login     → verify credentials → check MFA → create session → return token pair
POST /auth/refresh   → find session → delete old → create new session → return new token pair
POST /auth/logout    → find session by refresh token → delete (no auth header required)
```

Token types:

- **Access token** (RS256 JWT, 15 min) — carries `uid`, `role`, `type=access`
- **Refresh token** (ULID string, 7 days) — stored in `sessions` table, rotated on each refresh
- **MFA pending token** (RS256 JWT, 5 min) — carries `type=mfa_pending`, used during MFA verification flow

---

## Running Locally

```bash
# Start dependencies
just docker-up

# Generate RSA keys
just generate-keys

# Copy and edit environment
cp .env.example .env

# Apply migrations
just migrate-up

# Start server
just run
```
