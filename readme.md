# Observer

Case management platform for NGOs working with forcibly displaced people and animals. Replaces spreadsheets with a structured, auditable system for registering individuals, tracking support history, managing documents, and reporting.

**Stack:** Go 1.26 · PostgreSQL · Redis · React 19 · TanStack Router · Tailwind

## Development setup

**Prerequisites:** Go 1.26+, PostgreSQL, Redis, [Bun](https://bun.sh), [Just](https://just.systems)

```sh
# 1. Start dependencies
just docker-up

# 2. Generate RSA keys for JWT
just keygen

# 3. Apply migrations
just migrate-up

# 4. Install frontend dependencies
cd packages/observer-web && bun install && cd ../..

# 5. Run backend + frontend dev servers concurrently (port 9000 + Vite on 5173)
just dev
```

The backend exposes Swagger at `http://localhost:9000/swagger/index.html` in dev mode.

To create the first admin account:

```sh
just create-admin admin@example.com yourpassword
```

### Environment variables

All config is read from environment variables with sensible defaults. Key ones:

| Variable               | Default                    | Description                     |
| ---------------------- | -------------------------- | ------------------------------- |
| `DATABASE_DSN`         | —                          | Postgres connection string      |
| `REDIS_URL`            | `redis://localhost:6379/0` | Redis URL                       |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa`             | RSA private key                 |
| `JWT_PUBLIC_KEY_PATH`  | `keys/jwt_rsa.pub`         | RSA public key                  |
| `STORAGE_BACKEND`      | `local`                    | `local` or `s3`                 |
| `STORAGE_PATH`         | `data/uploads`             | Local upload directory          |
| `S3_BUCKET`            | —                          | S3 bucket (when backend=s3)     |
| `S3_REGION`            | `us-east-1`                | S3 region                       |
| `CORS_ORIGINS`         | `http://localhost:5173`    | Comma-separated allowed origins |
| `COOKIE_SECURE`        | `true`                     | Set `false` for local HTTP dev  |
| `SWAGGER_ENABLED`      | `false`                    | Enable Swagger UI               |
| `SENTRY_DSN`           | —                          | Sentry error tracking           |

## Project structure

```
cmd/observer/          CLI entrypoints (serve, migrate, keygen, create-admin, seed)
internal/
  domain/              Entities, repository interfaces, domain errors
    user/ auth/ project/ person/ support/ migration/
    household/ document/ pet/ note/ tag/ reference/
  usecase/             Business logic — all non-trivial logic lives here
    auth/ admin/ project/
  handler/             Thin HTTP adapters (bind → usecase → respond)
  middleware/          Auth, RBAC, project role checks
  postgres/            Repository implementations (sqlx)
  storage/             FileStorage interface + local + S3 backends
  crypto/              RSA keys, Argon2 hasher, token generator
  config/              Env-based config
  app/                 Manual DI container
  server/              HTTP server (Gin)
  spa/                 Embedded SPA for production builds
migrations/            Forward-only .up.sql files
packages/observer-web/ React SPA
  src/
    routes/            TanStack Router file-based routes
    components/        UI components grouped by domain
      ui/              Atoms (button, badge, icons, toast…)
      layout/          App shell (page-header, sidebar-link…)
      forms/           Form primitives (form-field, filter-bar, comboboxes)
      table/           Data table building blocks
      dialogs/         Modal dialogs
      drawer/          Drawer shell
      people/ pets/ support/ households/ migration/
      documents/ tags/ users/ permissions/ profile/
      reports/ charts/ search-palette/ auth/ date-picker/
    hooks/             React Query data hooks grouped by domain
      reference/       Geography reference data (countries, states, places…)
      people/ pets/ support/ households/ migration/
      documents/ tags/ notes/ users/ projects/ reports/
    stores/            Zustand stores (auth, toast)
    types/             TypeScript types matching API responses
    lib/               Utilities (api client, export, params, i18n)
    constants/         Enum → i18n key maps
adr/                   Architecture decision records
```

### Architecture

- **DDD + Clean Architecture** — domain layer has no dependencies; usecases depend on domain interfaces; handlers depend on usecases
- **Manual DI** — wired in `internal/app/container.go`, no DI framework
- **ULID IDs** — `ulid.ULID` in entities, `string` in DTOs
- **Forward-only migrations** — no rollbacks; schema only moves forward
- **Dual-level RBAC** — platform role (`admin/staff/consultant/guest`) × project role × 3 sensitivity flags

## Common tasks

```sh
just test              # Unit tests (fast, no Docker)
just test-all          # All tests including integration (requires Docker)
just generate-mocks    # Regenerate gomock mocks
just migrate-create add_foo        # New migration
just seed              # Seed demo data
just build-prod        # Build frontend + Go binary with embedded SPA
```

## Supporting the project

Contributions and [donations](https://paypal.me/SultanIman) are welcome.
