# Audit Log Structured Details Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a nullable `details JSONB` column to `audit_logs`, populate it with structured human-readable context at every audit call site, and render it in the frontend summary column.

**Architecture:** Forward-only migration adds the column; `audit.Entry` gains a `Details map[string]any` field; `AuditUseCase.Record()` gains a `details` parameter; each call site passes a typed map captured at write time; the frontend reads `details` when present and falls back to `summary` for older entries via a `formatAuditDetails` utility.

**Tech Stack:** Go (sqlx, encoding/json, squirrel), TypeScript/React (TanStack Router)

---

### Task 1: Database migration

**Files:**
- Create: `migrations/000037_add_audit_log_details.up.sql`

- [ ] **Step 1: Create migration file**

`migrations/000037_add_audit_log_details.up.sql`:
```sql
ALTER TABLE audit_logs ADD COLUMN details JSONB;
```

- [ ] **Step 2: Apply migration and verify**

```bash
go run ./cmd/observer migrate
```

Verify:
```bash
psql $DATABASE_URL -c "\d audit_logs" | grep details
```
Expected: `details | jsonb |`

- [ ] **Step 3: Commit**

```bash
git add migrations/000037_add_audit_log_details.up.sql
git commit -m "add details jsonb column to audit_logs"
```

---

### Task 2: Go — domain entity + repository

**Files:**
- Modify: `internal/domain/audit/entity.go`
- Modify: `internal/repository/audit/audit.go`

- [ ] **Step 1: Add Details field to audit.Entry**

Full updated `internal/domain/audit/entity.go`:
```go
package audit

import "time"

type Entry struct {
	ID         string
	ProjectID  *string
	UserID     *string // nil when the user has been deleted
	Action     string
	EntityType string
	EntityID   *string
	Summary    string
	Details    map[string]any
	IP         *string
	UserAgent  *string
	CreatedAt  time.Time

	// Populated by repository reads via LEFT JOIN with users; nil when user has been deleted.
	UserFirstName *string
	UserLastName  *string
	UserEmail     *string
}

type Filter struct {
	ProjectID  *string
	UserID     *string
	Action     *string
	EntityType *string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PerPage    int
}
```

- [ ] **Step 2: Add encoding/json import + Details to auditRow**

In `internal/repository/audit/audit.go`, add `"encoding/json"` to the import block and add `Details json.RawMessage` to `auditRow`:

```go
import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	domainaudit "github.com/lbrty/observer/internal/domain/audit"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

type auditRow struct {
	ID            string          `db:"id"`
	ProjectID     *string         `db:"project_id"`
	UserID        *string         `db:"user_id"`
	Action        string          `db:"action"`
	EntityType    string          `db:"entity_type"`
	EntityID      *string         `db:"entity_id"`
	Summary       string          `db:"summary"`
	Details       json.RawMessage `db:"details"`
	IP            *string         `db:"ip"`
	UserAgent     *string         `db:"user_agent"`
	CreatedAt     time.Time       `db:"created_at"`
	UserFirstName *string         `db:"user_first_name"`
	UserLastName  *string         `db:"user_last_name"`
	UserEmail     *string         `db:"user_email"`
}
```

- [ ] **Step 3: Update scanAuditRow to unmarshal Details**

Replace `scanAuditRow`:
```go
func scanAuditRow(r auditRow) domainaudit.Entry {
	entry := domainaudit.Entry{
		ID:            r.ID,
		ProjectID:     r.ProjectID,
		UserID:        r.UserID,
		Action:        r.Action,
		EntityType:    r.EntityType,
		EntityID:      r.EntityID,
		Summary:       r.Summary,
		IP:            r.IP,
		UserAgent:     r.UserAgent,
		CreatedAt:     r.CreatedAt,
		UserFirstName: r.UserFirstName,
		UserLastName:  r.UserLastName,
		UserEmail:     r.UserEmail,
	}
	if r.Details != nil {
		_ = json.Unmarshal(r.Details, &entry.Details)
	}
	return entry
}
```

- [ ] **Step 4: Update Log() to serialize Details**

Replace the `Log()` method body:
```go
func (r *auditLogRepo) Log(ctx context.Context, entry domainaudit.Entry) error {
	entry.ID = ulid.NewString()
	var detailsJSON []byte
	if entry.Details != nil {
		detailsJSON, _ = json.Marshal(entry.Details)
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audit_logs (id, project_id, user_id, action, entity_type, entity_id, summary, ip, user_agent, details)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		entry.ID, entry.ProjectID, entry.UserID, entry.Action, entry.EntityType, entry.EntityID,
		entry.Summary, entry.IP, entry.UserAgent, detailsJSON,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}
```

- [ ] **Step 5: Add a.details to List() SELECT**

In `List()`, change the `Select(...)` call to include `"a.details"` after `"a.summary"`:
```go
listSQL, listArgs, err := repository.PSQL.
    Select(
        "a.id", "a.project_id", "a.user_id", "a.action", "a.entity_type", "a.entity_id",
        "a.summary", "a.details", "a.ip", "a.user_agent", "a.created_at",
        "u.first_name AS user_first_name", "u.last_name AS user_last_name", "u.email AS user_email",
    ).
    From("audit_logs a").
    LeftJoin("users u ON u.id = a.user_id").
    Where(cond).
    OrderBy("a.created_at DESC").
    Limit(uint64(filter.PerPage)).
    Offset(uint64(offset)).
    ToSql()
```

- [ ] **Step 6: Build**

```bash
go build ./...
```
Expected: no errors

- [ ] **Step 7: Run tests**

```bash
just test
```
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/domain/audit/entity.go internal/repository/audit/audit.go
git commit -m "add Details to audit.Entry and repository log/scan"
```

---

### Task 3: Use case — Record() signature + DTO + all call sites pass nil

**Files:**
- Modify: `internal/usecase/audit/types.go`
- Modify: `internal/usecase/audit/audit_usecase.go`
- Modify: `internal/usecase/project/person_usecase.go`
- Modify: `internal/usecase/project/support_record_usecase.go`
- Modify: `internal/usecase/project/migration_record_usecase.go`
- Modify: `internal/usecase/project/household_usecase.go`
- Modify: `internal/usecase/project/note_usecase.go`
- Modify: `internal/usecase/project/document_usecase.go`
- Modify: `internal/usecase/project/pet_usecase.go`
- Modify: `internal/usecase/project/tag_usecase.go`
- Modify: `internal/usecase/admin/user_usecase.go`
- Modify: `internal/usecase/admin/permission_usecase.go`
- Modify: `internal/handler/report/export.go`

The goal of this task is to make the codebase compile with the new signature. All call sites pass `nil` for now; real details are added in Tasks 4 and 5.

- [ ] **Step 1: Update types.go — add Details to EntryDTO and LogInput**

Full updated `internal/usecase/audit/types.go`:
```go
package audit

type LogInput struct {
	ProjectID  *string
	UserID     *string
	Action     string
	EntityType string
	EntityID   *string
	Summary    string
	Details    map[string]any
	IP         string
	UserAgent  string
}

type ListInput struct {
	ProjectID  *string `form:"project_id"`
	UserID     *string `form:"user_id"`
	Action     *string `form:"action"`
	EntityType *string `form:"entity_type"`
	DateFrom   *string `form:"date_from"`
	DateTo     *string `form:"date_to"`
	Page       int     `form:"page"`
	PerPage    int     `form:"per_page"`
}

type EntryDTO struct {
	ID            string         `json:"id"`
	ProjectID     *string        `json:"project_id"`
	UserID        *string        `json:"user_id"`
	Action        string         `json:"action"`
	EntityType    string         `json:"entity_type"`
	EntityID      *string        `json:"entity_id"`
	Summary       string         `json:"summary"`
	Details       map[string]any `json:"details,omitempty"`
	IP            *string        `json:"ip,omitempty"`
	UserAgent     *string        `json:"user_agent,omitempty"`
	CreatedAt     string         `json:"created_at"`
	UserFirstName *string        `json:"user_first_name"`
	UserLastName  *string        `json:"user_last_name"`
	UserEmail     *string        `json:"user_email"`
}

type ListOutput struct {
	Entries []EntryDTO `json:"entries"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
}
```

- [ ] **Step 2: Update Record() signature + body in audit_usecase.go**

Replace the `Record()` method:
```go
func (uc *AuditUseCase) Record(
	ctx context.Context,
	projectID *string,
	action, entityType string,
	entityID *string,
	summary string,
	details map[string]any,
) {
	if uc == nil {
		return
	}

	var userID *string
	if uid := middleware.AuditUserID(ctx); uid != "" {
		userID = &uid
	}

	entry := domainaudit.Entry{
		ProjectID:  projectID,
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Summary:    summary,
		Details:    details,
		IP:         strPtrOrNil(middleware.AuditIP(ctx)),
		UserAgent:  strPtrOrNil(middleware.AuditUserAgent(ctx)),
	}

	if err := uc.repo.Log(ctx, entry); err != nil {
		slog.Error("audit log failed", slog.String("action", action), slog.Any("err", err))
	}
}
```

Also update `Log()` to propagate `input.Details`:
```go
func (uc *AuditUseCase) Log(ctx context.Context, input LogInput) error {
	entry := domainaudit.Entry{
		ProjectID:  input.ProjectID,
		UserID:     input.UserID,
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		Summary:    input.Summary,
		Details:    input.Details,
		IP:         strPtrOrNil(input.IP),
		UserAgent:  strPtrOrNil(input.UserAgent),
	}
	if err := uc.repo.Log(ctx, entry); err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	return nil
}
```

Update `List()` DTO mapping to include Details:
```go
dtos[i] = EntryDTO{
    ID:            e.ID,
    ProjectID:     e.ProjectID,
    UserID:        e.UserID,
    Action:        e.Action,
    EntityType:    e.EntityType,
    EntityID:      e.EntityID,
    Summary:       e.Summary,
    Details:       e.Details,
    IP:            e.IP,
    UserAgent:     e.UserAgent,
    CreatedAt:     e.CreatedAt.Format(time.RFC3339),
    UserFirstName: e.UserFirstName,
    UserLastName:  e.UserLastName,
    UserEmail:     e.UserEmail,
}
```

- [ ] **Step 3: Update all call sites — add nil as last argument**

In every file listed under Files above, find each `auditUC.Record(` call and append `, nil` before the closing `)`. There are 34 calls total across 11 files. After this step every call looks like, e.g.:

```go
// person_usecase.go — Create
uc.auditUC.Record(ctx, &projectID, "person.create", "person", &p.ID,
    fmt.Sprintf("Created person %s", p.ID), nil)

// support_record_usecase.go — Create
uc.auditUC.Record(ctx, &projectID, "support_record.create", "support_record", &r.ID,
    fmt.Sprintf("Created support record %s", r.ID), nil)

// export.go — ExportPeople
h.auditUC.Record(c.Request.Context(), &projectID, "export", "person", nil,
    fmt.Sprintf("exported %d people", len(out.People)), nil)
```

(Apply the same `, nil` append to all remaining calls.)

- [ ] **Step 4: Build**

```bash
go build ./...
```
Expected: no errors

- [ ] **Step 5: Run tests**

```bash
just test
```
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add \
  internal/usecase/audit/types.go \
  internal/usecase/audit/audit_usecase.go \
  internal/usecase/project/person_usecase.go \
  internal/usecase/project/support_record_usecase.go \
  internal/usecase/project/migration_record_usecase.go \
  internal/usecase/project/household_usecase.go \
  internal/usecase/project/note_usecase.go \
  internal/usecase/project/document_usecase.go \
  internal/usecase/project/pet_usecase.go \
  internal/usecase/project/tag_usecase.go \
  internal/usecase/admin/user_usecase.go \
  internal/usecase/admin/permission_usecase.go \
  internal/handler/report/export.go
git commit -m "add details param to Record(), thread nil through all call sites"
```

---

### Task 4: Project use case call sites — structured details

**Files:**
- Modify: `internal/usecase/project/person_usecase.go`
- Modify: `internal/usecase/project/support_record_usecase.go`
- Modify: `internal/usecase/project/migration_record_usecase.go`
- Modify: `internal/usecase/project/household_usecase.go`
- Modify: `internal/usecase/project/note_usecase.go`
- Modify: `internal/usecase/project/document_usecase.go`
- Modify: `internal/usecase/project/pet_usecase.go`
- Modify: `internal/usecase/project/tag_usecase.go`

Replace every `, nil)` added in Task 3 with a real details map.

- [ ] **Step 1: person_usecase.go — Create**

Replace the Record call in `Create()`:
```go
name := p.FirstName
if p.LastName != nil {
    name += " " + *p.LastName
}
details := map[string]any{"name": name}
if p.ExternalID != nil {
    details["external_id"] = *p.ExternalID
}
uc.auditUC.Record(ctx, &projectID, "person.create", "person", &p.ID,
    fmt.Sprintf("Created person %s", p.ID), details)
```

- [ ] **Step 2: person_usecase.go — Update**

Replace the Record call in `Update()` (p is available after the update):
```go
name := p.FirstName
if p.LastName != nil {
    name += " " + *p.LastName
}
details := map[string]any{"name": name}
if p.ExternalID != nil {
    details["external_id"] = *p.ExternalID
}
uc.auditUC.Record(ctx, &projectID, "person.update", "person", &id,
    fmt.Sprintf("Updated person %s", id), details)
```

- [ ] **Step 3: person_usecase.go — Delete**

Replace the Record call in `Delete()` (p is fetched at the top of Delete):
```go
name := p.FirstName
if p.LastName != nil {
    name += " " + *p.LastName
}
uc.auditUC.Record(ctx, &projectID, "person.delete", "person", &id,
    fmt.Sprintf("Deleted person %s", id), map[string]any{"name": name})
```

- [ ] **Step 4: support_record_usecase.go — Create**

Replace the Record call in `Create()` (p is fetched for validation):
```go
personName := p.FirstName
if p.LastName != nil {
    personName += " " + *p.LastName
}
details := map[string]any{
    "person_name": personName,
    "type":        string(r.Type),
}
if r.Sphere != nil {
    details["sphere"] = string(*r.Sphere)
}
uc.auditUC.Record(ctx, &projectID, "support_record.create", "support_record", &r.ID,
    fmt.Sprintf("Created support record %s", r.ID), details)
```

- [ ] **Step 5: support_record_usecase.go — Update**

Replace the Record call in `Update()` (r is fetched from DB):
```go
details := map[string]any{
    "person_id": r.PersonID,
    "type":      string(r.Type),
}
if r.Sphere != nil {
    details["sphere"] = string(*r.Sphere)
}
uc.auditUC.Record(ctx, &projectID, "support_record.update", "support_record", &id,
    fmt.Sprintf("Updated support record %s", id), details)
```

- [ ] **Step 6: support_record_usecase.go — Delete**

Replace the Record call in `Delete()` (r is fetched from DB):
```go
uc.auditUC.Record(ctx, &projectID, "support_record.delete", "support_record", &id,
    fmt.Sprintf("Deleted support record %s", id),
    map[string]any{"person_id": r.PersonID})
```

- [ ] **Step 7: migration_record_usecase.go — Create**

Replace the Record call in `Create()` (personID is a parameter; input holds place IDs):
```go
details := map[string]any{"person_id": personID}
if input.FromPlaceID != nil {
    details["origin_place_id"] = *input.FromPlaceID
}
if input.DestinationPlaceID != nil {
    details["destination_place_id"] = *input.DestinationPlaceID
}
uc.auditUC.Record(ctx, &projectID, "migration_record.create", "migration_record", &r.ID,
    fmt.Sprintf("Created migration record %s", r.ID), details)
```

- [ ] **Step 8: migration_record_usecase.go — Update**

Replace the Record call in `Update()` (personID is a parameter; r holds place IDs):
```go
details := map[string]any{"person_id": personID}
if r.FromPlaceID != nil {
    details["origin_place_id"] = *r.FromPlaceID
}
if r.DestinationPlaceID != nil {
    details["destination_place_id"] = *r.DestinationPlaceID
}
uc.auditUC.Record(ctx, &projectID, "migration_record.update", "migration_record", &id,
    fmt.Sprintf("Updated migration record %s", id), details)
```

- [ ] **Step 9: migration_record_usecase.go — Delete**

Replace the Record call in `Delete()` (rec is fetched from DB):
```go
uc.auditUC.Record(ctx, &projectID, "migration_record.delete", "migration_record", &id,
    fmt.Sprintf("Deleted migration record %s", id),
    map[string]any{"person_id": rec.PersonID})
```

- [ ] **Step 10: household_usecase.go — Create**

Replace the Record call in `Create()` (h is the newly-created household):
```go
details := map[string]any{}
if h.ReferenceNumber != nil {
    details["reference_number"] = *h.ReferenceNumber
}
uc.auditUC.Record(ctx, &projectID, "household.create", "household", &h.ID,
    fmt.Sprintf("Created household %s", h.ID), details)
```

- [ ] **Step 11: household_usecase.go — Update**

Replace the Record call in `Update()` (h is fetched then mutated):
```go
details := map[string]any{}
if h.ReferenceNumber != nil {
    details["reference_number"] = *h.ReferenceNumber
}
uc.auditUC.Record(ctx, &projectID, "household.update", "household", &id,
    fmt.Sprintf("Updated household %s", id), details)
```

- [ ] **Step 12: household_usecase.go — Delete**

Replace the Record call in `Delete()` (h is fetched from DB):
```go
details := map[string]any{}
if h.ReferenceNumber != nil {
    details["reference_number"] = *h.ReferenceNumber
}
uc.auditUC.Record(ctx, &projectID, "household.delete", "household", &id,
    fmt.Sprintf("Deleted household %s", id), details)
```

- [ ] **Step 13: note_usecase.go — Create**

Replace the Record call in `Create()` (personID is a parameter):
```go
uc.auditUC.Record(ctx, &projectID, "note.create", "note", &n.ID,
    fmt.Sprintf("Created note %s", n.ID),
    map[string]any{"person_id": personID})
```

- [ ] **Step 14: note_usecase.go — Update**

Replace the Record call in `Update()` (n is fetched from DB):
```go
uc.auditUC.Record(ctx, &projectID, "note.update", "note", &id,
    fmt.Sprintf("Updated note %s", id),
    map[string]any{"person_id": n.PersonID})
```

- [ ] **Step 15: note_usecase.go — Delete**

Replace the Record call in `Delete()` (n is fetched from DB):
```go
uc.auditUC.Record(ctx, &projectID, "note.delete", "note", &id,
    fmt.Sprintf("Deleted note %s", id),
    map[string]any{"person_id": n.PersonID})
```

- [ ] **Step 16: document_usecase.go — Upload**

Replace the Record call in `Upload()` (filename and personID are parameters):
```go
uc.auditUC.Record(ctx, &projectID, "document.upload", "document", &doc.ID,
    fmt.Sprintf("Uploaded document %s", doc.Name),
    map[string]any{"filename": filename, "person_id": personID})
```

- [ ] **Step 17: document_usecase.go — Download**

Replace the Record call in `Download()` (d is fetched from DB):
```go
uc.auditUC.Record(ctx, &d.ProjectID, "document.download", "document", &d.ID,
    fmt.Sprintf("Downloaded document %s", d.Name),
    map[string]any{"filename": d.Name, "person_id": d.PersonID})
```

- [ ] **Step 18: document_usecase.go — Update**

Replace the Record call in `Update()` (d is fetched from DB):
```go
uc.auditUC.Record(ctx, &projectID, "document.update", "document", &id,
    fmt.Sprintf("Updated document %s", id),
    map[string]any{"filename": d.Name, "person_id": d.PersonID})
```

- [ ] **Step 19: document_usecase.go — Delete**

Replace the Record call in `Delete()` (d is fetched from DB):
```go
uc.auditUC.Record(ctx, &d.ProjectID, "document.delete", "document", &d.ID,
    fmt.Sprintf("Deleted document %s", d.Name),
    map[string]any{"filename": d.Name, "person_id": d.PersonID})
```

- [ ] **Step 20: pet_usecase.go — Create**

Replace the Record call in `Create()` (p is the newly-created pet):
```go
details := map[string]any{"name": p.Name}
if p.OwnerID != nil {
    details["owner_id"] = *p.OwnerID
}
uc.auditUC.Record(ctx, &projectID, "pet.create", "pet", &p.ID,
    fmt.Sprintf("Created pet %s", p.ID), details)
```

- [ ] **Step 21: pet_usecase.go — Update**

Replace the Record call in `Update()` (p is fetched then mutated):
```go
details := map[string]any{"name": p.Name}
if p.OwnerID != nil {
    details["owner_id"] = *p.OwnerID
}
uc.auditUC.Record(ctx, &projectID, "pet.update", "pet", &id,
    fmt.Sprintf("Updated pet %s", id), details)
```

- [ ] **Step 22: pet_usecase.go — Delete**

Replace the Record call in `Delete()` (p is fetched from DB):
```go
details := map[string]any{"name": p.Name}
if p.OwnerID != nil {
    details["owner_id"] = *p.OwnerID
}
uc.auditUC.Record(ctx, &projectID, "pet.delete", "pet", &id,
    fmt.Sprintf("Deleted pet %s", id), details)
```

- [ ] **Step 23: tag_usecase.go — Create**

Replace the Record call in `Create()`:
```go
uc.auditUC.Record(ctx, &projectID, "tag.create", "tag", &t.ID,
    fmt.Sprintf("Created tag %s", t.Name),
    map[string]any{"name": t.Name})
```

- [ ] **Step 24: tag_usecase.go — Update**

Replace the Record call in `Update()` (t is fetched from DB):
```go
uc.auditUC.Record(ctx, &projectID, "tag.update", "tag", &id,
    fmt.Sprintf("Updated tag %s", id),
    map[string]any{"name": t.Name})
```

- [ ] **Step 25: tag_usecase.go — Delete**

Replace the Record call in `Delete()` (t is fetched from DB):
```go
uc.auditUC.Record(ctx, &projectID, "tag.delete", "tag", &id,
    fmt.Sprintf("Deleted tag %s", id),
    map[string]any{"name": t.Name})
```

- [ ] **Step 26: Build and test**

```bash
go build ./...
just test
```
Expected: PASS

- [ ] **Step 27: Commit**

```bash
git add \
  internal/usecase/project/person_usecase.go \
  internal/usecase/project/support_record_usecase.go \
  internal/usecase/project/migration_record_usecase.go \
  internal/usecase/project/household_usecase.go \
  internal/usecase/project/note_usecase.go \
  internal/usecase/project/document_usecase.go \
  internal/usecase/project/pet_usecase.go \
  internal/usecase/project/tag_usecase.go
git commit -m "populate structured details in project use case audit calls"
```

---

### Task 5: Admin + export handler call sites — structured details

**Files:**
- Modify: `internal/usecase/admin/user_usecase.go`
- Modify: `internal/usecase/admin/permission_usecase.go`
- Modify: `internal/handler/report/export.go`

- [ ] **Step 1: user_usecase.go — Create**

Replace the Record call in `Create()`:
```go
uc.auditUC.Record(ctx, nil, "admin.user.create", "user", &uid,
    fmt.Sprintf("Created user %s with role %s", newUser.Email, newUser.Role),
    map[string]any{"email": newUser.Email, "role": string(newUser.Role)})
```

- [ ] **Step 2: user_usecase.go — role_change (in Update)**

Replace the Record call inside the `if u.Role != oldRole` block:
```go
uc.auditUC.Record(ctx, nil, "user.role_change", "user", &uid,
    fmt.Sprintf("Changed role from %s to %s for user %s", oldRole, u.Role, uid),
    map[string]any{
        "email":    u.Email,
        "old_role": string(oldRole),
        "new_role": string(u.Role),
    })
```

- [ ] **Step 3: user_usecase.go — reset_password**

Replace the Record call in `ResetPassword()`:
```go
uc.auditUC.Record(ctx, nil, "admin.user.reset_password", "user", &uid,
    fmt.Sprintf("Password reset for user %s", uid),
    map[string]any{"user_id": uid})
```

- [ ] **Step 4: user_usecase.go — DeactivateUser**

Move the Record call to after the `GetByID` fetch so `u.Email` is available. Replace the current body from the Deactivate call onward:

```go
if err := uc.userRepo.Deactivate(ctx, id); err != nil {
    return nil, fmt.Errorf("deactivate user: %w", err)
}

uid := id.String()

u, err := uc.userRepo.GetByID(ctx, id)
if err != nil {
    return nil, err
}

uc.auditUC.Record(ctx, nil, "user.deactivate", "user", &uid,
    "User account deactivated",
    map[string]any{"email": u.Email, "user_id": uid})

dto := userToDTO(u)
return &dto, nil
```

- [ ] **Step 5: user_usecase.go — ReactivateUser**

Same reorder for `ReactivateUser`:

```go
if err := uc.userRepo.Reactivate(ctx, id); err != nil {
    return nil, fmt.Errorf("reactivate user: %w", err)
}

uid := id.String()

u, err := uc.userRepo.GetByID(ctx, id)
if err != nil {
    return nil, err
}

uc.auditUC.Record(ctx, nil, "user.reactivate", "user", &uid,
    "User account reactivated",
    map[string]any{"email": u.Email, "user_id": uid})

dto := userToDTO(u)
return &dto, nil
```

- [ ] **Step 6: permission_usecase.go — Assign**

The current code discards the user object: `if _, err := uc.userRepo.GetByID(ctx, uid); err != nil`. Change to capture it, then use it in details:

```go
u, err := uc.userRepo.GetByID(ctx, uid)
if err != nil {
    return nil, fmt.Errorf("verify user for assign: %w", err)
}
```

Then replace the Record call:
```go
uc.auditUC.Record(ctx, &projectID, "permission.grant", "permission", &perm.ID,
    fmt.Sprintf("Granted %s to user %s", input.Role, input.UserID),
    map[string]any{
        "subject_name":  u.FirstName + " " + u.LastName,
        "subject_email": u.Email,
        "role":          string(role),
        "project_id":    projectID,
    })
```

- [ ] **Step 7: permission_usecase.go — Update**

After `uc.permRepo.Update(ctx, perm)`, add a best-effort user fetch and replace the Record call:

```go
if err := uc.permRepo.Update(ctx, perm); err != nil {
    return nil, fmt.Errorf("update permission: %w", err)
}

details := map[string]any{"role": string(perm.Role), "project_id": projectID}
if uid, err := ulid.Parse(perm.UserID); err == nil {
    if u, err := uc.userRepo.GetByID(ctx, uid); err == nil {
        details["subject_name"] = u.FirstName + " " + u.LastName
        details["subject_email"] = u.Email
    }
}

uc.auditUC.Record(ctx, &projectID, "admin.permission.update", "permission", &id,
    fmt.Sprintf("Updated permission %s", id), details)
```

- [ ] **Step 8: permission_usecase.go — Revoke**

After `uc.permRepo.Delete(ctx, id)`, add the same user fetch pattern:

```go
if err := uc.permRepo.Delete(ctx, id); err != nil {
    return fmt.Errorf("revoke permission: %w", err)
}

details := map[string]any{"project_id": projectID}
if uid, err := ulid.Parse(perm.UserID); err == nil {
    if u, err := uc.userRepo.GetByID(ctx, uid); err == nil {
        details["subject_name"] = u.FirstName + " " + u.LastName
        details["subject_email"] = u.Email
    }
}

uc.auditUC.Record(ctx, &projectID, "permission.revoke", "permission", &id,
    fmt.Sprintf("Revoked permission %s", id), details)

return nil
```

- [ ] **Step 9: export.go — ExportPeople**

Replace the Record call:
```go
h.auditUC.Record(c.Request.Context(), &projectID, "export", "person", nil,
    fmt.Sprintf("exported %d people", len(out.People)),
    map[string]any{"entity_type": "person", "count": len(out.People)})
```

- [ ] **Step 10: export.go — ExportSupportRecords**

Replace the Record call:
```go
h.auditUC.Record(c.Request.Context(), &projectID, "export", "support_record", nil,
    fmt.Sprintf("exported %d support records", len(out.Records)),
    map[string]any{"entity_type": "support_record", "count": len(out.Records)})
```

- [ ] **Step 11: export.go — ExportPets**

Replace the Record call:
```go
h.auditUC.Record(c.Request.Context(), &projectID, "export", "pet", nil,
    fmt.Sprintf("exported %d pets", len(out.Pets)),
    map[string]any{"entity_type": "pet", "count": len(out.Pets)})
```

- [ ] **Step 12: export.go — ExportHouseholds**

Replace the Record call:
```go
h.auditUC.Record(c.Request.Context(), &projectID, "export", "household", nil,
    fmt.Sprintf("exported %d households", len(out.Households)),
    map[string]any{"entity_type": "household", "count": len(out.Households)})
```

- [ ] **Step 13: Build and test**

```bash
go build ./...
just test
```
Expected: PASS

- [ ] **Step 14: Commit**

```bash
git add \
  internal/usecase/admin/user_usecase.go \
  internal/usecase/admin/permission_usecase.go \
  internal/handler/report/export.go
git commit -m "populate structured details in admin and export audit calls"
```

---

### Task 6: TypeScript — AuditEntry type + formatAuditDetails utility

**Files:**
- Modify: `packages/observer-web/src/types/audit.ts`
- Create: `packages/observer-web/src/lib/audit-details.ts`

- [ ] **Step 1: Add details to AuditEntry**

Full updated `packages/observer-web/src/types/audit.ts`:
```typescript
export interface AuditEntry {
  id: string;
  project_id: string | null;
  user_id: string | null;
  action: string;
  entity_type: string;
  entity_id: string | null;
  summary: string;
  details: Record<string, unknown> | null;
  ip: string;
  user_agent: string;
  created_at: string;
  user_first_name: string;
  user_last_name: string;
  user_email: string;
}

export interface AuditListParams {
  project_id?: string;
  user_id?: string;
  action?: string;
  entity_type?: string;
  date_from?: string;
  date_to?: string;
  page?: number;
  per_page?: number;
}

export interface AuditListOutput {
  entries: AuditEntry[];
  total: number;
  page: number;
  per_page: number;
}
```

- [ ] **Step 2: Create src/lib/audit-details.ts**

```typescript
import type { AuditEntry } from "@/types/audit";

export function formatAuditDetails(entry: AuditEntry): string {
  const d = entry.details;
  if (!d) return entry.summary;

  const { action, entity_type } = entry;
  const ref = entry.entity_id ?? "";

  if (entity_type === "person") {
    const name = (d.name as string | undefined) ?? ref;
    if (action.endsWith(".create")) return `Created person ${name}`;
    if (action.endsWith(".update")) return `Updated person ${name}`;
    if (action.endsWith(".delete")) return `Deleted person ${name}`;
  }

  if (entity_type === "support_record") {
    const personRef =
      (d.person_name as string | undefined) ?? (d.person_id as string | undefined) ?? ref;
    const type = d.type as string | undefined;
    const sphere = d.sphere as string | undefined;
    const qualifier = [type, sphere].filter(Boolean).join("/");
    if (action.endsWith(".create"))
      return `Created ${qualifier ? qualifier + " " : ""}support record for ${personRef}`;
    if (action.endsWith(".update"))
      return `Updated ${qualifier ? qualifier + " " : ""}support record for ${personRef}`;
    if (action.endsWith(".delete")) return `Deleted support record for ${personRef}`;
  }

  if (entity_type === "migration_record") {
    const personRef = (d.person_id as string | undefined) ?? ref;
    if (action.endsWith(".create")) return `Created migration record for ${personRef}`;
    if (action.endsWith(".update")) return `Updated migration record for ${personRef}`;
    if (action.endsWith(".delete")) return `Deleted migration record for ${personRef}`;
  }

  if (entity_type === "household") {
    const householdRef = (d.reference_number as string | undefined) ?? ref;
    if (action.endsWith(".create")) return `Created household ${householdRef}`;
    if (action.endsWith(".update")) return `Updated household ${householdRef}`;
    if (action.endsWith(".delete")) return `Deleted household ${householdRef}`;
  }

  if (entity_type === "note") {
    const personRef = (d.person_id as string | undefined) ?? "";
    const forPart = personRef ? ` for ${personRef}` : "";
    if (action.endsWith(".create")) return `Created note${forPart}`;
    if (action.endsWith(".update")) return `Updated note${forPart}`;
    if (action.endsWith(".delete")) return `Deleted note${forPart}`;
  }

  if (entity_type === "document") {
    const filename = (d.filename as string | undefined) ?? ref;
    if (action === "document.upload") return `Uploaded ${filename}`;
    if (action === "document.download") return `Downloaded ${filename}`;
    if (action.endsWith(".update")) return `Updated document ${filename}`;
    if (action.endsWith(".delete")) return `Deleted ${filename}`;
  }

  if (entity_type === "pet") {
    const name = (d.name as string | undefined) ?? ref;
    const ownerId = d.owner_id as string | undefined;
    if (action.endsWith(".create"))
      return `Created pet ${name}${ownerId ? ` (owner: ${ownerId})` : ""}`;
    if (action.endsWith(".update")) return `Updated pet ${name}`;
    if (action.endsWith(".delete")) return `Deleted pet ${name}`;
  }

  if (entity_type === "tag") {
    const name = (d.name as string | undefined) ?? ref;
    if (action.endsWith(".create")) return `Created tag "${name}"`;
    if (action.endsWith(".update")) return `Updated tag "${name}"`;
    if (action.endsWith(".delete")) return `Deleted tag "${name}"`;
  }

  if (entity_type === "user") {
    if (action === "admin.user.create") {
      const email = (d.email as string | undefined) ?? "";
      const role = d.role as string | undefined;
      return `Created user ${email}${role ? ` with role ${role}` : ""}`;
    }
    if (action === "user.role_change") {
      const email = (d.email as string | undefined) ?? "";
      const oldRole = d.old_role as string | undefined;
      const newRole = d.new_role as string | undefined;
      return `Changed role ${oldRole} → ${newRole} for ${email}`;
    }
    if (action === "user.deactivate") {
      return `Deactivated user ${(d.email as string | undefined) ?? (d.user_id as string | undefined) ?? ref}`;
    }
    if (action === "user.reactivate") {
      return `Reactivated user ${(d.email as string | undefined) ?? (d.user_id as string | undefined) ?? ref}`;
    }
  }

  if (entity_type === "permission") {
    const subjectRef =
      (d.subject_name as string | undefined) ??
      (d.subject_email as string | undefined) ??
      (d.user_id as string | undefined) ??
      ref;
    const role = d.role as string | undefined;
    if (action === "permission.grant")
      return `Granted ${role ?? ""} permission to ${subjectRef}`;
    if (action === "admin.permission.update")
      return `Updated permission for ${subjectRef}${role ? ` (role: ${role})` : ""}`;
    if (action === "permission.revoke") return `Revoked permission for ${subjectRef}`;
  }

  if (action === "export") {
    const count = d.count as number | undefined;
    const entityRef = (d.entity_type as string | undefined) ?? entity_type;
    return `Exported ${count ?? ""} ${entityRef}`;
  }

  return entry.summary;
}
```

- [ ] **Step 3: Type-check**

```bash
cd packages/observer-web && bunx tsc --noEmit
```
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add \
  packages/observer-web/src/types/audit.ts \
  packages/observer-web/src/lib/audit-details.ts
git commit -m "add details field to AuditEntry and formatAuditDetails utility"
```

---

### Task 7: Frontend — render structured details in summary column

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/audit-logs.lazy.tsx`
- Modify: `packages/observer-web/src/routes/_app/admin/audit-logs.lazy.tsx`

- [ ] **Step 1: Update project audit-logs page**

In `packages/observer-web/src/routes/_app/projects/$projectId/audit-logs.lazy.tsx`:

Add import after the `AuditEntry` import:
```typescript
import { formatAuditDetails } from "@/lib/audit-details";
```

Replace the summary column definition:
```typescript
{
  key: "summary",
  header: t("audit.summary"),
  render: (e) => (
    <span className="max-w-xs truncate text-sm text-fg">{formatAuditDetails(e)}</span>
  ),
},
```

- [ ] **Step 2: Update admin audit-logs page**

In `packages/observer-web/src/routes/_app/admin/audit-logs.lazy.tsx`:

Add import after the `AuditEntry` import:
```typescript
import { formatAuditDetails } from "@/lib/audit-details";
```

Replace the summary column definition (same change as Step 1):
```typescript
{
  key: "summary",
  header: t("audit.summary"),
  render: (e) => (
    <span className="max-w-xs truncate text-sm text-fg">{formatAuditDetails(e)}</span>
  ),
},
```

- [ ] **Step 3: Type-check**

```bash
cd packages/observer-web && bunx tsc --noEmit
```
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add \
  "packages/observer-web/src/routes/_app/projects/\$projectId/audit-logs.lazy.tsx" \
  packages/observer-web/src/routes/_app/admin/audit-logs.lazy.tsx
git commit -m "render structured audit details in summary column"
```
