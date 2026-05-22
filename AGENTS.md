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
  usecase/        # business logic — subdirs mirror domain groups (admin/, auth/, project/, …)
  handler/        # thin HTTP adapters — subdirs mirror usecase groups (admin/, auth/, project/, …)
    errors.go     # shared: HandleError, BindJSON, MapDomainError
    handlertest/  # shared test helpers for handler subdirectory tests
  repository/     # repository implementations — subdirs mirror domain groups (user/, project/, …)
    interfaces.go # all repository interfaces
    mock/         # gomock mocks
  middleware/     # HTTP middleware (auth, RBAC)
  crypto/         # RSA keys, Argon hasher, token generator
  storage/        # file storage interface + local filesystem + S3 impl (ADR-010)
  config/         # reads env vars via github.com/sultaniman/env
  server/         # HTTP server setup + route wiring
  app/            # DI container (manual wiring)
adr/              # architectural decision records
migrations/       # forward-only SQL migrations

packages/observer-web/src/
  components/     # React components, organized by domain
  hooks/          # React hooks, mirroring component domain layout
  routes/         # TanStack Router file-based routes
  stores/         # client state (Zustand)
  constants/      # shared constants by domain (i18n.ts, person.ts, support.ts, …)
  types/          # TypeScript types generated from Go (see ADR-008) + manual
  lib/            # utilities (form-error, export-csv, tag-color, …)
  locales/        # i18n JSON files (en, ky, ru, uk, de, tr)
```

## Code Conventions

### Naming & Language

- Short, clear names: `ix` (index), `uq` (unique)
- Default UI/content language: Kyrgyz Latin (ky), not Russian

### Comments

- Simple docstrings only
- No decorative separators (`//-----`, `//=====`, `/* ── ... ── */`), no ASCII art
- Complex logic: mermaid diagrams + module README instead of lengthy inline comments

### Architecture Rules

- Business logic in `internal/usecase/`, never in handlers or SQL
- Handlers are thin — bind request, call use case, return response
- Manual DI wired in `internal/app/container.go`
- Domain entities define repository interfaces; `internal/repository/<group>/` implements them
- `ulid.ULID` for entity IDs, `string` in DTOs (via `.String()`)
- Prefer libs in this set: Gin, testify, gomock, testcontainers-go, sqlx
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
- Constants live in `constants/<domain>.ts`; i18n key maps go in `constants/i18n.ts`

## Do's and Don'ts

### Do

**Before implementing**

- State assumptions explicitly. If uncertain, ask — don't guess silently.
- If multiple interpretations exist, present them rather than picking one silently.
- Share a brief plan and wait for confirmation. For small, well-scoped changes, stating the plan is enough.
- Check `adr/` before exploring the codebase broadly.
- Read the entire plan file before starting, when one exists.

**While implementing**

- Match existing style, even if you'd do it differently.
- Keep new tests consistent with existing test design.
- Remove imports, variables, and functions your changes leave unused.
- Use `bun`/`bunx` for all frontend package management — never `npm`, `npx`, or `node`.
- Use `@/` imports for all frontend modules except colocated siblings (`./`). Follow the import group order.
- Check `base-ui` and `@phosphor-icons/react` before adding new components or icons.
- Put business logic in `internal/usecase/` — never in handlers or SQL.
- Keep handlers thin: bind request, call use case, return response.
- Use `ulid.ULID` for entity IDs and `string` in DTOs via `.String()`.
- Use `just` commands for building and testing, not `make`.

**After implementing**

- Run pre-commit hooks and fix any errors.
- Update documentation when the components it describes change.

### Don't

**Scope**

- Don't add features beyond what was asked.
- Don't add abstractions for single-use code.
- Don't add error handling for impossible scenarios.
- Don't improve adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Don't remove pre-existing dead code — mention it instead.
- Don't skip or work around failing tests — fix them before proceeding.

**Git**

- Don't commit unless explicitly asked or after manual review. Commit messages: lowercase, short.
- Don't run `git add -A`. Stage files individually.

**Code style**

- Don't write decorative comments (`// ─────`, `// =====`) or ASCII art.
- Don't ignore linting errors without a strong, stated reason.
- Don't write lengthy inline comments for complex logic — use mermaid diagrams or a module README instead.

**Files**

- Don't edit `src/apiTypes.ts` manually.
- Don't wire dependencies outside `internal/app/container.go`.
