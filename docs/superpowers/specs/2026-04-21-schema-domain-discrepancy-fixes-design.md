---
title: Schema ↔ Domain Discrepancy Fixes
date: 2026-04-21
status: approved
---

# Schema ↔ Domain Discrepancy Fixes

## Background

A cross-reference of the 35 SQL migrations against all domain entities, repository
interfaces, and use-case layers revealed 7 discrepancies. This document captures the
agreed design for fixing all of them in a single sweep (Approach A).

## Items

### 1. `SearchHits` — move to domain layer

**Problem:** `SearchHits` and its nested types are defined in `internal/repository/search_types.go`,
putting a domain result type below the domain boundary.

**Fix:**
- Create `internal/domain/search/types.go` with `SearchHits` and all nested types.
- Update `internal/repository/interfaces.go` to import and reference `search.SearchHits`.
- Update `internal/repository/search_repository.go` implementation to use the domain type.
- Update all call sites (use cases, handlers) to import from `internal/domain/search`.
- Delete `internal/repository/search_types.go`.

---

### 2. `PetRepository.List` — introduce `PetListFilter`

**Problem:** `PetRepository.List` takes raw parameters instead of a filter struct,
making it a breaking change to add any new filter field.

**Fix:**
- Add `PetListFilter` struct to `internal/domain/pet/entity.go`:
  ```go
  type PetListFilter struct {
      ProjectID string
      Status    *PetStatus
      TagIDs    []string
      Page      int
      PerPage   int
  }
  ```
- Change `PetRepository.List` signature to `List(ctx context.Context, filter pet.PetListFilter) ([]*pet.Pet, int, error)`.
- Update mock, postgres implementation, use case, and handler call sites.

---

### 3. `MigrationRecordRepository` — add `Delete`

**Problem:** `MigrationRecordRepository` is the only mutable entity repository without
a `Delete` method. Migration records are editable (`UpdatedAt` exists), so omission
is an oversight, not intentional immutability.

**Fix:**
- Add `Delete(ctx context.Context, id string) error` to `MigrationRecordRepository` interface.
- Implement in postgres repository.
- Regenerate mock.
- Wire into the migration record use case and handler (DELETE endpoint).

---

### 4. `MFARecoveryCode.ID` — change to `ulid.ULID`

**Problem:** `MFARecoveryCode.ID` is `string` while all other user/auth domain entities
use `ulid.ULID` for their IDs (`User.ID`, `Session.ID`, `VerificationToken.ID`).

**Fix:**
- Change `MFARecoveryCode.ID` from `string` to `ulid.ULID`.
- Update `MFARecoveryCodeRepository.MarkUsed(ctx, id string)` to `MarkUsed(ctx, id ulid.ULID)`.
- Update postgres implementation and mock.
- Update all call sites that construct or read `MFARecoveryCode.ID`.

---

### 5. `audit.Entry` — fix IP/UserAgent nullability

**Problem:** `audit_logs.ip` and `audit_logs.user_agent` have no `NOT NULL` constraint in
the DB, but `audit.Entry.IP` and `audit.Entry.UserAgent` are non-nullable `string` in Go.
A system-generated entry with no IP/user agent would scan incorrectly.

**Fix:**
- Change `Entry.IP` from `string` to `*string`.
- Change `Entry.UserAgent` from `string` to `*string`.
- Update audit log repository implementation (scan + insert).
- Update all call sites that construct `audit.Entry` or read these fields.
- No migration needed — DB columns are already nullable.

---

### 6. Enrichment fields — document as read-only projections

**Problem:** Three entity structs carry JOIN-populated fields that don't exist as columns.
Create/Update callers must know to leave them zero-valued, but nothing in the code signals this.

Affected:
- `household.Household`: `HeadPersonName *string`, `MemberCount int`
- `support.Record`: `PersonFirstName *string`, `PersonLastName *string`
- `audit.Entry`: `UserFirstName string`, `UserLastName string`, `UserEmail string`

**Fix:**
- Add a `// Populated by repository reads; zero-valued on writes.` comment block above
  each group of enrichment fields in the struct definition.
- No logic changes required.

---

### 7. `PersonCategoryRepository` — add `ListBulk`

**Problem:** `PersonTagRepository` and `PetTagRepository` both expose `ListBulk` for
efficient list-view hydration. `PersonCategoryRepository` only has `List(personID string)`,
forcing N+1 queries on any list view that shows categories per person.

**Fix:**
- Add `ListBulk(ctx context.Context, personIDs []string) (map[string][]string, error)` to
  `PersonCategoryRepository` interface.
- Implement in postgres repository using `WHERE person_id = ANY($1)`.
- Regenerate mock.
- Update any list use case that currently N+1s category hydration.

---

### 8. `//go:generate` — replace giant reflect-mode list with source mode

**Problem:** Line 24 of `internal/repository/interfaces.go` is a 280-character `//go:generate`
directive listing 29 comma-separated interface names in reflect mode. It is hard to read,
hard to diff, and requires manual maintenance every time an interface is added.
A bonus bug: `SearchRepository` is defined in `interfaces.go` but absent from the list,
meaning it has never been mocked.

**Fix:**
- Create `internal/repository/generate.go`:
  ```go
  package repository

  //go:generate mockgen -source=interfaces.go -destination=mock/repository.go -package=mock
  ```
- Remove the `//go:generate` line from `interfaces.go`.
- Source mode auto-discovers all interfaces in the file, including `SearchRepository`.
- Run `just generate-mocks` to verify the mock regenerates correctly and includes `SearchRepository`.

---

## Execution order

Dependencies run top to bottom; items with no dependencies on each other can be done
in any order within the same pass.

1. `PetListFilter` struct (new type, no deps)
2. `SearchHits` move to `internal/domain/search/` (new package, then update imports)
3. `PersonCategoryRepository.ListBulk` (additive — interface + mock + impl)
4. `MigrationRecordRepository.Delete` (additive — interface + mock + impl + handler)
5. `MFARecoveryCode.ID` → `ulid.ULID` (type change — propagate through callers)
6. `audit.Entry` IP/UserAgent → `*string` (type change — propagate through callers)
7. Enrichment field comments (doc-only, no logic)
8. `generate.go` + source mode (do last, after all interface changes are in place — one `just generate-mocks` regenerates everything)

## Out of scope

- `person_status_history` domain entity — trigger-managed writes + `StatusFlowReport`
  aggregation already provide sufficient abstraction. No Go entity needed.
- Migrations for audit nullability — DB is already correct; only Go types change.
- Enrichment field extraction to separate DTOs — acceptable as-is for the current read patterns.
