# Audit Log Structured Details

**Date:** 2026-04-28  
**Status:** Approved

## Problem

Audit log summaries currently display raw ULID entity IDs (e.g., `Updated permission 01K6Y9T7XBGD2KHBNWJ6ZPA`). These are opaque to users and provide no actionable context.

## Goal

Capture human-readable context (names, roles, relationships) at write time so audit entries remain meaningful even if the referenced entity is later deleted.

## Approach

Add a nullable `details JSONB` column to `audit_logs`. Populate it alongside the existing `summary` text at every `Record()` call site. The frontend renders from `details` when present and falls back to `summary` for older entries.

## Database

**Migration** `000037_add_audit_log_details.up.sql`:

```sql
ALTER TABLE audit_logs ADD COLUMN details JSONB;
```

No index on `details` — existing indices on `project_id`, `action`, `entity_type`, and `created_at` cover all query patterns.

## Go Layer

### Domain (`internal/domain/audit/entity.go`)

Add `Details map[string]any` field to `Entry`:

```go
type Entry struct {
    // ... existing fields ...
    Details map[string]any // nil for legacy entries
}
```

### Use case (`internal/usecase/audit/audit_usecase.go`)

Extend `Record()` signature:

```go
func (uc *AuditUseCase) Record(
    ctx context.Context,
    projectID *string,
    action, entityType string,
    entityID *string,
    summary string,
    details map[string]any,
) {
```

All existing call sites pass `nil` as `details` initially, then are updated to pass structured maps.

### Repository (`internal/repository/audit/audit.go`)

Add `details` as the 10th parameter in the `INSERT`:

```go
INSERT INTO audit_logs
  (id, project_id, user_id, action, entity_type, entity_id, summary, ip, user_agent, details)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
```

Pass `entry.Details` serialized with `json.Marshal` → `[]byte` (consistent with how `PhoneNumbers` JSONB is handled in `person_usecase.go`). Pass `nil` when `Details` is nil — the column is nullable.

### DTO (`internal/usecase/audit/types.go`)

Add `Details map[string]any` to `EntryDTO`:

```go
Details map[string]any `json:"details,omitempty"`
```

## Details Payload Per Entity Type

Each call site captures its natural human label plus contextual fields:

| Entity type | Fields captured |
|---|---|
| `permission` | `subject_name`, `subject_email`, `project_name`, `role` |
| `person` | `name`, `external_id` (if set) |
| `support_record` | `person_name`, `type`, `sphere` |
| `migration_record` | `person_name`, `origin`, `destination` |
| `document` | `person_name`, `filename` |
| `note` | `person_name` |
| `pet` | `name`, `person_name` |
| `user` | `email`, `old_role`, `new_role` (for role changes); `email`, `role` (for create) |
| `tag` | `name` |
| `export` | `entity_type`, `count`, `format` |

All fields are strings or numbers — no nested objects.

## Frontend

### Type (`packages/observer-web/src/types/audit.ts`)

```typescript
export interface AuditEntry {
  // ... existing fields ...
  details: Record<string, unknown> | null;
}
```

### Utility (`packages/observer-web/src/lib/audit-details.ts`)

`formatAuditDetails(entry: AuditEntry): string` — returns a human-readable string built from `details` fields, keyed by `entity_type` and `action`. Falls back to `entry.summary` when `details` is null.

Examples:
- permission update → `"Updated permission for Jane Doe (consultant) on Bishkek IDP Project"`
- person create → `"Created person Asel Mamytova"`
- support_record create → `"Created legal/housing support record for Asel Mamytova"`
- migration_record create → `"Created migration record for Asel Mamytova (Donetsk → Bishkek)"`
- document upload → `"Uploaded passport.pdf for Asel Mamytova"`
- user role change → `"Changed role admin → staff for jane@example.com"`
- tag create → `"Created tag \"vulnerable\""`
- export → `"Exported 42 people as CSV"`

### Column renderer

The summary column in both audit log pages (`_app/projects/$projectId/audit-logs.lazy.tsx` and the admin equivalent) calls `formatAuditDetails(entry)` instead of `e.summary`.

## Backward Compatibility

- Existing audit entries have `details = NULL` — `formatAuditDetails` returns `entry.summary` unchanged.
- No data migration required.

## Out of Scope

- Diffing old/new field values for update operations (Phase 2).
- Full-text search on details content.
- i18n of the formatted string on the backend — rendering is client-side only.
