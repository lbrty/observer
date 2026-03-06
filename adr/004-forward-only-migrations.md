---
title: "ADR-004: Forward-Only Migrations"
weight: 4
---

| Field      | Value                          |
| ---------- | ------------------------------ |
| Status     | Accepted                       |
| Date       | 2026-02-21                     |
| Supersedes | —                              |
| Components | observer, database, migrations |

---

## Context

ADR-001 and ADR-002 described migrations with both `.up.sql` and `.down.sql` files, following the default `golang-migrate` convention. In practice, rollback migrations introduce more problems than they solve for this type of application:

- Rollback SQL is rarely tested and frequently wrong
- Data written by the up migration is destroyed or left orphaned on rollback
- Production databases are never rolled back — incidents are resolved by deploying a fix forward
- Maintaining two files per migration doubles the maintenance surface with near-zero benefit

## Decision

**All migrations are forward-only.** Only `.up.sql` files exist. Down migrations are not created, not committed, and not supported by the CLI tooling.

## Consequences

### Migration files

Each migration is a single `.up.sql` file:

```text
migrations/
├── 000001_init_extensions.up.sql
├── 000002_create_users_table.up.sql
├── ...
└── 000021_add_office_to_users.up.sql
```

No `.down.sql` files exist anywhere in the repository.

### Naming convention

Files follow the pattern:

```text
{NNNNNN}_{description}.up.sql
```

- `NNNNNN` — zero-padded 6-digit sequential integer
- `description` — snake_case description of the change
- Sequence is monotonically increasing; gaps are not permitted

### `migrate create` command

The `observer migrate create <name>` command creates a single `.up.sql` file. The sequence number auto-increments from the highest existing file in the migrations directory. An explicit number can be provided via `--seq`:

```text
observer migrate create add_index_to_people
# creates: migrations/000022_add_index_to_people.up.sql

observer migrate create --seq 25 add_column
# creates: migrations/000025_add_column.up.sql
```

### How to "undo" a migration

If a schema change needs to be reverted, write a new forward migration that undoes it:

```sql
-- 000022_drop_unused_column.up.sql
ALTER TABLE people DROP COLUMN IF EXISTS legacy_field;
```

This keeps the audit trail intact and is safe to apply in production.

### `migrate down` is not supported

The CLI exposes only `up` and `version` subcommands. There is no `down` command.

## Alternatives Considered

**Keep `.down.sql` files but never run them** — rejected; dead files create confusion about whether rollback is supported.

**Use a separate rollback script directory** — rejected; unnecessary complexity for a policy decision.
