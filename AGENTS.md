# Observer — Agent Instructions

## Quick Ref

- **Module**: `github.com/lbrty/observer` | **Go** 1.26.\* | **Port**: 9000
- **Architecture**: DDD + Clean Architecture, manual DI (no frameworks)
- **Build**: `Justfile` (not Makefile) | **Frontend pkg manager**: `bun` | **Always use `bun`/`bunx`** (never `npm`/`npx`/`node`)
- **ADRs**: `adr/`

## Project Layout

```
cmd/              # entrypoints
internal/
  domain/         # entities + repository interfaces
  usecase/        # business logic (use cases live here, NOT in handlers/DB)
  handler/        # thin HTTP adapters — delegate to use cases
  middleware/     # HTTP middleware (auth, RBAC)
  postgres/       # repository implementations
  crypto/         # RSA keys, Argon hasher, token generator
  storage/        # file storage interface + local filesystem + S3 impl (ADR-010)
  config/         # reads env vars with defaults
  server/         # HTTP server setup
  app/            # DI container (manual wiring)
adr/              # architectural decision records
migrations/       # forward-only SQL migrations
```

## Code Conventions

### Naming & Language

- Short, clear names: `ix` (index), `uq` (unique)
- Default UI/content language: Kyrgyz Latin (ky), not Russian

### Comments

- Simple docstrings only
- No decorative separators (`//-----`, `//=====`, `/* ── ... ── */`), no ASCII art
- Complex logic: mermaid diagrams + module README instead of lengthy text

### Architecture Rules

- Business logic in `internal/usecase/`, never in handlers or SQL
- Handlers are thin — bind request, call use case, return response
- Manual DI wired in `internal/app/container.go`
- Domain entities define repository interfaces; `internal/postgres/` implements them
- `ulid.ULID` for entity IDs, `string` in DTOs (via `.String()`)
- Prefer well-maintained, widely-known libs (Gin, testify, gomock, testcontainers-go, sqlx)
- Pragmatic MVP: core functionality first, defer advanced features (MEK/DEK, detailed audit logs) to Phase 2

## Testing

```bash
just test            # unit tests only (fast, no Docker)
just test-all        # all tests including integration (Docker required)
just generate-mocks  # regenerate gomock mocks
```

- Unit: testify `assert`/`require` + gomock
- Integration: testcontainers-go (Postgres, Redis), guarded by `testing.Short()`
- Never skip verification — fix failing tests before proceeding

## Frontend

### Imports

`@/` alias for all imports. Exception: colocated siblings use `./`.

Order (blank line between groups):

1. `react`, `react-dom`
2. External libs (`@tanstack/*`, `@zxcvbn-ts/*`)
3. Workspace packages (`@observer/*`)
4. App aliases (`@/components/*`, `@/stores/*`, `@/hooks/*`)
5. Colocated (`./constants`, `./types`)
6. Styles (`.module.css`) — always last

### Component & Hook Structure

Both `src/components/` and `src/hooks/` are organized by domain:

- **`ui/`** — atoms: button, status-badge, icons, toast, tooltip, ui-select, ui-switch, alert-banner, empty-state
- **`layout/`** — app shell: page-header, section-heading, sidebar-link, app-footer
- **`table/`** — data table: data-table, data-table-page, pagination, row-actions
- **`forms/`** — form primitives: form-field, form-section, filter-bar, base-combobox, place-combobox
- **`dialogs/`** — modal dialogs: confirm-dialog, form-dialog
- **`drawer/`** — drawer shell
- Domain folders: `people/`, `pets/`, `support/`, `households/`, `migration/`, `documents/`, `tags/`, `users/`, `permissions/`, `profile/`, `reports/`, `charts/`, `auth/`, `date-picker/`, `search-palette/`, `my-stats/`

Hooks mirror the same domain layout under `src/hooks/`: `reference/`, `people/`, `pets/`, `support/`, `households/`, `migration/`, `documents/`, `tags/`, `notes/`, `users/`, `projects/`, `reports/`. General-purpose hooks (`use-drawer-form`, `use-export-csv`, `use-search`, `use-audit-logs`, `use-schema-status`) stay at `src/hooks/`.

### Components & Tooling

- Check `base-ui` and `@phosphor-icons/react` first for existing components/icons
- React compiler enabled — omit effect dependencies where possible
- Extract shared constants to `constants.ts` (root-level if cross-module)

### Tailwind `@apply` Order

When >10 rules, separate `@apply` per group on its own line:

positioning > layout > sizing > borders > background > padding/margin > text > transforms > rest

