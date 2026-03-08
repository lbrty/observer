# Reports, Export, Audit Logs & Profile Settings — Design

| Field  | Value      |
| ------ | ---------- |
| Date   | 2026-03-08 |
| Status | Approved   |


## 1. Reports — Consolidated Presets + Custom Builder

### Approach

Hybrid: curated preset reports for stable donor deliverables + a form-based custom report builder for ad-hoc slicing.

ADR-005's 39 reports collapse into 14 presets. The former Group 9 (tag search) moves entirely into the custom builder. A single `ReportUseCase.GenerateReport(params)` method backs both presets (named parameter sets) and custom queries.

### Preset Reports (14)

| #  | Preset                                  | Metric         | Grouped by                | Filters    |
| -- | --------------------------------------- | -------------- | ------------------------- | ---------- |
| 1  | Consultation totals                     | Events         | type                      | date range |
| 2  | Registrations by sex                    | People         | sex                       | date range |
| 3  | Consultations by sex                    | People         | sex x type                | date range |
| 4  | Registrations by IDP status             | People         | conflict_zone             | date range |
| 5  | Consultations by IDP status             | People         | conflict_zone x type      | date range |
| 6  | Registrations by region                 | People         | current region            | date range |
| 7  | Consultations by region                 | People         | current region x type     | date range |
| 8  | Consultations by sphere                 | Events + People| sphere x type             | date range |
| 9  | Consultations by office                 | Events         | office x type             | date range |
| 10 | Consultations by age group              | Events + People| age_group x type          | date range |
| 11 | Consultations by vulnerability category | People         | category x type           | date range |
| 12 | Family unit summary                     | People + Units | type                      | date range |
| 13 | People with pets                        | People         | pet status                | date range |
| 14 | Pet counts by status                    | Pets           | pet status                | date range |

### Custom Report Builder

Form-based UI:

- **Metric**: events / people / family units / pets
- **Group by** (1-2 dimensions): sex, age group, region, conflict zone, office, sphere, category, person tags, pet tags, pet status
- **Filters**: date range, consultation type, specific tag IDs, has pets (boolean)
- **Output**: table + CSV export (gated by `ActionExport`)

### Backend

Single `ReportUseCase` with `GenerateReport(params ReportParams)`. Presets are named parameter constants, not separate query paths. Dynamic SQL generation with parameterized filters (no string concatenation — safe from injection).


## 2. Data Table Filters + Export

### Filters

Filter bar above every data table. Filter state stored in URL search params (shareable/bookmarkable). Collapses to summary chips when not editing.

**Common to all tables:** date range, tags (multi-select), free text search

**Per-table filters:**

| Table             | Additional filters                                              |
| ----------------- | --------------------------------------------------------------- |
| People            | sex, age group, case status, conflict zone, region, category, has pets |
| Support records   | type (legal/social), sphere, office, referral status            |
| Migration records | movement reason, housing at destination                         |
| Pets              | pet status, pet tags                                            |
| Households        | member count range                                              |

### Filter mechanics

- Filters applied as query params to existing list endpoints (extend `PersonListFilter`, `RecordListFilter` etc.)
- Reuse existing `tag-filter.tsx` component pattern for multi-select filters

### Export

- **Format**: CSV (single format for MVP)
- **Scope**: current filtered view — what the user sees is what gets exported
- **Flow**: export button click -> backend streams CSV matching current filters -> audit log entry created
- **Permission**: `ActionExport` (see below)
- Export button hidden when user lacks permission


## 3. Export Permission

New project-scoped action: `ActionExport`

- Added to the `Action` enum alongside read/create/update/delete/manage_members
- Minimum role: `ProjectRoleConsultant` (rank 2)
- Added to `MinRoleForAction` map
- Enforced in export endpoints via `RequireProjectRole(ActionExport)`
- Viewers cannot export; consultants, managers, owners can


## 4. Theme & Language in Profile

### Change

- Create profile/settings page at route `/settings/profile`
- Move theme picker (system/light/dark/light-hc/dark-hc) and language picker (ky/en/ru/uk/de/tr) from avatar menu + footer into profile page
- Avatar menu links to profile page instead
- Storage remains `localStorage` — no backend changes

### Profile page contents (MVP)

- User info (name, email) — read-only display
- Theme selection (radio group or segmented control)
- Language selection (dropdown)


## 5. Audit Logs

### Schema

```sql
CREATE TABLE audit_logs (
    id          TEXT PRIMARY KEY,   -- ULID
    project_id  TEXT,               -- NULL for platform-level actions
    user_id     TEXT NOT NULL,
    action      TEXT NOT NULL,      -- e.g. 'export.people', 'person.create'
    entity_type TEXT NOT NULL,      -- e.g. 'person', 'pet', 'document'
    entity_id   TEXT,               -- ULID of affected entity, NULL for bulk ops
    summary     TEXT NOT NULL,      -- human-readable context
    ip          TEXT,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_audit_logs_project_time ON audit_logs (project_id, created_at);
CREATE INDEX ix_audit_logs_user_time    ON audit_logs (user_id, created_at);
CREATE INDEX ix_audit_logs_action_time  ON audit_logs (action, created_at);
```

### Audited Actions

| Bucket           | Actions                                                                                    |
| ---------------- | ------------------------------------------------------------------------------------------ |
| Data export      | export.people, export.support_records, export.migration_records, export.households, export.pets |
| Documents        | document.upload, document.download, document.delete                                        |
| Record lifecycle | person.create, person.delete, pet.create, pet.delete, household.create, household.delete, support_record.create, support_record.delete, migration_record.create, migration_record.delete |
| Admin            | user.role_change, project.create, project.delete, permission.grant, permission.revoke      |

### Implementation

- `AuditLogger` interface in domain layer: `Log(ctx context.Context, entry AuditEntry) error`
- Postgres implementation writes to `audit_logs` table
- Called from use cases (not handlers or middleware) — business logic layer
- Injected via DI in `container.go`

### Admin UI

- `/admin/audit-log` — full audit table, admin role required, filters by user/action/entity type/date/project
- `/projects/:id/audit-log` — project-scoped view, project manager+ role required
- Read-only tables
