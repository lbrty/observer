# Observer — Agent Instructions

## Quick Ref

- **Module**: `github.com/lbrty/observer` | **Go** 1.25.\* | **Port**: 9000
- **Architecture**: DDD + Clean Architecture, manual DI (no frameworks)
- **Build**: `Justfile` (not Makefile) | **Frontend pkg manager**: `bun`
- **ADRs**: `docs/adr/` | **Testing guide**: `docs/testing-guide.md` | **Variables**: `docs/variables.md`

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
  config/         # reads env vars with defaults
  server/         # HTTP server setup
  app/            # DI container (manual wiring)
docs/adr/         # architectural decision records
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
- Full guide with examples: `docs/testing-guide.md`

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

### Components & Tooling

- Check `base-ui` and `@phosphor-icons/react` first for existing components/icons
- React compiler enabled — omit effect dependencies where possible
- Extract shared constants to `constants.ts` (root-level if cross-module)

### Tailwind `@apply` Order

When >10 rules, separate `@apply` per group on its own line:

positioning > layout > sizing > borders > background > padding/margin > text > transforms > rest

## Variables

Use values from `docs/variables.md` when implementing features. Suggest new entries when values repeat across files.
