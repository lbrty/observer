# ADR-003: Main Schema


| Field      | Value                      |
| ---------- | -------------------------- |
| Status     | Accepted                   |
| Date       | 2026-02-21                 |
| Supersedes | —                          |
| Components | observer, database, schema |


## Domain Overview

```text
Geography:  countries → states (conflict_zone) → places
Org:        offices
Taxonomy:   categories (multi-code via person_categories), tags (project-scoped)
Projects:   projects, project_permissions
People:     people, person_tags, person_categories, person_notes, documents
Household:  households, household_members
Movement:   migration_records (movement_reason, housing_at_destination)
Support:    support_records (referral_status, referred_to_office)
Animals:    pets
```


## Migration Files

### Updated: 000002 — Users (role fix + profile columns)

The ADR-002 users migration is updated in-place to use the correct role enum and add profile fields.
`office_id` is added separately in migration 000021 (after `offices` is created in 000010).

```sql
CREATE TABLE users (
    id          TEXT         PRIMARY KEY,
    first_name  TEXT,
    last_name   TEXT,
    email       VARCHAR(255) NOT NULL,
    phone       VARCHAR(20)  NOT NULL,
    role        VARCHAR(50)  NOT NULL CHECK (role IN ('admin', 'staff', 'consultant', 'guest')),
    is_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT uq_users_phone UNIQUE (phone)
);
```

**Platform-level role semantics:**

| Role         | Access                                          |
| ------------ | ----------------------------------------------- |
| `admin`      | Full platform access, user management           |
| `staff`      | Create/manage projects, view all cases          |
| `consultant` | Assigned to projects, works with people records |
| `guest`      | Read-only on explicitly assigned projects       |


### 000008 — States

`conflict_zone` tags an oblast/state with its conflict zone designation. IDP classification is derived at query time via `origin_place_id → places.state_id → states.conflict_zone` rather than being hardcoded on the person record. This is a free-text field — no enum constraint — so new zones can be added by seeding the states table without a schema migration.

```sql
CREATE TABLE states (
    id            TEXT        PRIMARY KEY,
    country_id    TEXT        NOT NULL REFERENCES countries (id) ON DELETE CASCADE,
    name          TEXT        NOT NULL,
    code          CITEXT,
    conflict_zone TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX        ix_states_country_id   ON states (country_id);
CREATE INDEX        ix_states_name         ON states (name);
CREATE UNIQUE INDEX uq_states_country_code ON states (country_id, code) WHERE code IS NOT NULL;
```


### 000010 — Offices

`state_id` dropped — derivable via `place_id → places.state_id`. No global name uniqueness — different cities legitimately share office names (e.g. "Main Office" in Kyiv and Kherson).

```sql
CREATE TABLE offices (
    id         TEXT        PRIMARY KEY,
    name       TEXT        NOT NULL,
    place_id   TEXT        REFERENCES places (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_offices_place_id ON offices (place_id);
```


### 000012 — Tags

Project-scoped free-form labels. Tags are unique per project, not globally — each project maintains its own tag vocabulary independently.

```sql
CREATE TABLE tags (
    id         TEXT        PRIMARY KEY,
    project_id TEXT        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_tags_project_name ON tags (project_id, name);
CREATE INDEX        ix_tags_project_id   ON tags (project_id);
```


### 000014 — Project Permissions

Per-user, per-project access. Action capabilities derive from `role`; data sensitivity is controlled by explicit boolean flags.

```sql
CREATE TYPE project_role AS ENUM ('owner', 'manager', 'consultant', 'viewer');

CREATE TABLE project_permissions (
    id                 TEXT         PRIMARY KEY,
    project_id         TEXT         NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    user_id            TEXT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role               project_role NOT NULL DEFAULT 'viewer',
    can_view_contact   BOOLEAN      NOT NULL DEFAULT FALSE,
    can_view_personal  BOOLEAN      NOT NULL DEFAULT FALSE,
    can_view_documents BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_project_permissions_user_project ON project_permissions (user_id, project_id);
CREATE INDEX        ix_project_permissions_project_id   ON project_permissions (project_id);
CREATE INDEX        ix_project_permissions_user_id      ON project_permissions (user_id);
```

**Project-role → implied action permissions (enforced in middleware):**

| Role         | read | create | update | delete | manage members     |
| ------------ | ---- | ------ | ------ | ------ | ------------------ |
| `viewer`     | ✓    |        |        |        |                    |
| `consultant` | ✓    | ✓      | ✓      |        |                    |
| `manager`    | ✓    | ✓      | ✓      | ✓      | ✓                  |
| `owner`      | ✓    | ✓      | ✓      | ✓      | ✓ + delete project |

`projects.owner_id` implicitly grants owner-level role — no `project_permissions` row required for the project owner.


### 000016 — Person Tags

Junction table replacing the legacy `TEXT[]` tags column on `people`.

```sql
CREATE TABLE person_tags (
    person_id TEXT NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    tag_id    TEXT NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (person_id, tag_id)
);

CREATE INDEX ix_person_tags_tag_id ON person_tags (tag_id);
```


### 000018 — Support Records

Tracks assistance provided to a person.

- `recorded_by` — user who entered the record (was `owner_id` in archive, which had no FK)
- `consultant_id` — consultant who delivered the service
- `office_id` — office that provided or coordinated this consultation (direct FK, not derived through the consultant's home office — see ADR-005 §office_id)
- `referred_to_office` — the office the client was referred to; non-NULL together with `referral_status` indicates this record is a referral, not direct delivery
- `referral_status` — lifecycle of an outbound referral; NULL means direct service delivery
- `provided_at` — date service was delivered (distinct from `created_at`, the DB insertion time)
- `sphere` — topic area for "by sphere of appeal" report breakdowns; constrained to `support_sphere` enum values (see ADR-005)

```sql
CREATE TYPE support_type AS ENUM (
    'humanitarian', 'legal', 'social', 'psychological', 'medical', 'general'
);

CREATE TYPE support_sphere AS ENUM (
    'housing_assistance', 'document_recovery', 'social_benefits', 'property_rights',
    'employment_rights', 'family_law', 'healthcare_access', 'education_access',
    'financial_aid', 'psychological_support', 'other'
);

CREATE TABLE support_records (
    id                  TEXT           PRIMARY KEY,
    person_id           TEXT           NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    project_id          TEXT           NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    consultant_id       TEXT           REFERENCES users (id) ON DELETE SET NULL,
    recorded_by         TEXT           REFERENCES users (id) ON DELETE SET NULL,
    office_id           TEXT           REFERENCES offices (id) ON DELETE SET NULL,
    referred_to_office  TEXT           REFERENCES offices (id) ON DELETE SET NULL,
    type                support_type   NOT NULL DEFAULT 'general',
    sphere              support_sphere,
    referral_status     TEXT           CHECK (referral_status IN (
                                           'pending', 'accepted', 'completed', 'declined', 'no_response'
                                       )),
    provided_at         DATE,
    notes               TEXT,
    created_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_support_records_person_id     ON support_records (person_id);
CREATE INDEX ix_support_records_project_id    ON support_records (project_id);
CREATE INDEX ix_support_records_consultant_id ON support_records (consultant_id);
CREATE INDEX ix_support_records_office_id     ON support_records (office_id);
CREATE INDEX ix_support_records_type          ON support_records (type);
CREATE INDEX ix_support_records_provided_at   ON support_records (provided_at);
CREATE INDEX ix_support_records_sphere        ON support_records (sphere) WHERE sphere IS NOT NULL;
```


### 000020 — Documents

Documents attached to a person record, scoped to a project.

- `path` — relative path from the storage root
- `encryption_key_ref` — KMS key identifier (not raw key material); null for unencrypted files

```sql
CREATE TABLE documents (
    id                 TEXT        PRIMARY KEY,
    person_id          TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    project_id         TEXT        NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    uploaded_by        TEXT        REFERENCES users (id) ON DELETE SET NULL,
    encryption_key_ref TEXT,
    name               TEXT        NOT NULL,
    path               TEXT        NOT NULL,
    mime_type          TEXT        NOT NULL,
    size               BIGINT      NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_documents_person_id  ON documents (person_id);
CREATE INDEX ix_documents_project_id ON documents (project_id);
```


### 000022 — Person Notes

Informal case-worker notes on a person record. Distinct from `support_records` (formal, quantified service events) — notes are internal observations and reminders.

```sql
CREATE TABLE person_notes (
    id         TEXT        PRIMARY KEY,
    person_id  TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    author_id  TEXT        REFERENCES users (id) ON DELETE SET NULL,
    body       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_person_notes_person_id ON person_notes (person_id);
CREATE INDEX ix_person_notes_author_id ON person_notes (author_id);
```


### 000024 — Person Categories

Replaces the single `people.category_id` FK with a many-to-many junction, mirroring UNHCR's PSN multi-code model. A person may hold multiple concurrent vulnerability classifications (e.g. single parent + mobility impairment).

```sql
CREATE TABLE person_categories (
    person_id   TEXT NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    category_id TEXT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (person_id, category_id)
);

CREATE INDEX ix_person_categories_category_id ON person_categories (category_id);
```


## Summary of Migration Numbers

| Number | Table / Change                                                                                                |
| ------ | ------------------------------------------------------------------------------------------------------------- |
| 000002 | users (role fix, profile columns)                                                                             |
| 000007 | countries                                                                                                     |
| 000008 | states (+ conflict_zone)                                                                                      |
| 000009 | places                                                                                                        |
| 000010 | offices (no state_id, no global name unique)                                                                  |
| 000011 | categories                                                                                                    |
| 000012 | tags (project-scoped)                                                                                         |
| 000013 | projects (+ status)                                                                                           |
| 000014 | project_permissions (project_role enum + sensitivity flags)                                                   |
| 000015 | people (split name, primary_phone, consent, age XOR, external_id unique, no idp_status/parent_id/category_id) |
| 000016 | person_tags                                                                                                   |
| 000017 | migration_records (movement_reason, housing_at_destination, immutable)                                        |
| 000018 | support_records (recorded_by, provided_at, sphere, referral_status, referred_to_office)                       |
| 000019 | pets (owner → people)                                                                                         |
| 000020 | documents (renamed from person_documents)                                                                     |
| 000021 | users office_id (alter)                                                                                       |
| 000022 | person_notes                                                                                                  |
| 000023 | households + household_members (replaces people.parent_id)                                                    |
| 000024 | person_categories junction (replaces people.category_id)                                                      |
