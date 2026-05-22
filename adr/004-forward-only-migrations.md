# ADR-004: Forward-Only Migrations

| Field      | Value                          |
| ---------- | ------------------------------ |
| Status     | Accepted                       |
| Date       | 2026-02-21                     |
| Supersedes | —                              |
| Components | observer, database, migrations |

## Decision

All migrations are **forward-only**. There are no `.down.sql` rollback files. Every schema change is a new numbered `.up.sql` file that can only advance the schema version.

## Rationale

- **Production systems do not roll back** schema changes in practice. A failed migration is fixed by a new corrective migration, not by reverting.
- `.down.sql` files are rarely tested, frequently wrong, and destructive — dropping columns and tables deletes data.
- Forward-only migrations make the migration history the authoritative record of schema evolution. Every intentional change is permanent and auditable.
- Rollback of a bad deploy is done at the **application level** (redeploy the previous binary) while keeping the schema at the new version. Application code is written to tolerate the new schema from day one.

## Convention

Migration files live in `migrations/` and follow this naming pattern:

```
<six-digit-number>_<snake_case_description>.up.sql
```

Example: `000035_add_mfa_recovery_codes.up.sql`

There is no corresponding `000035_add_mfa_recovery_codes.down.sql`.

## Applying Migrations

```bash
just migrate-up   # applies all pending .up.sql files
```

The `migrate` CLI subcommand embedded in the observer binary also applies pending migrations on startup when invoked directly.

