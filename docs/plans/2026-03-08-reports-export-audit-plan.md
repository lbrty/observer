# Reports, Export, Audit Logs & Profile Settings — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Consolidate reports into 14 presets + custom builder, add data table filters with CSV export gated by a new ActionExport permission, move theme/language into a profile page, and implement simple audit logging.

**Architecture:** Extends the existing DDD + Clean Architecture layers. New domain types for audit logs, new repository interfaces + Postgres implementations, use case orchestration, thin handlers. Frontend adds filter bars to data tables, a custom report builder form, a profile settings page, and admin audit log views.

**Tech Stack:** Go (Gin, sqlx, gomock, testify), React (TanStack Router + Query, base-ui), PostgreSQL, CSV streaming

---

## Phase 1: Export Permission (ActionExport)

Foundation for export and audit features. No external dependencies.

### Task 1: Add ActionExport to domain

**Files:**
- Modify: `internal/domain/project/entity.go`
- Modify: `internal/middleware/project_auth.go` (no code change needed — it reads from MinRoleForAction dynamically)

**Step 1: Add ActionExport constant and map entry**

In `internal/domain/project/entity.go`, add to the Action constants:

```go
const (
	ActionRead          Action = "read"
	ActionCreate        Action = "create"
	ActionUpdate        Action = "update"
	ActionDelete        Action = "delete"
	ActionManageMembers Action = "manage_members"
	ActionExport        Action = "export"
)
```

And to `MinRoleForAction`:

```go
var MinRoleForAction = map[Action]ProjectRole{
	ActionRead:          ProjectRoleViewer,
	ActionCreate:        ProjectRoleConsultant,
	ActionUpdate:        ProjectRoleConsultant,
	ActionDelete:        ProjectRoleManager,
	ActionManageMembers: ProjectRoleManager,
	ActionExport:        ProjectRoleConsultant,
}
```

**Step 2: Run existing tests to verify nothing breaks**

Run: `just test`
Expected: All pass — ActionExport is additive, middleware reads from map dynamically.

**Step 3: Commit**

```bash
git add internal/domain/project/entity.go
git commit -m "add ActionExport permission for project-scoped data export"
```

---

## Phase 2: Audit Logs (Backend)

New feature end-to-end: migration → domain → repository → use case → handler → routes.

### Task 2: Create audit_logs migration

**Files:**
- Create: `migrations/000028_create_audit_logs.up.sql`

**Step 1: Write the migration**

```sql
CREATE TABLE audit_logs (
    id          TEXT        PRIMARY KEY,
    project_id  TEXT        REFERENCES projects(id) ON DELETE SET NULL,
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action      TEXT        NOT NULL,
    entity_type TEXT        NOT NULL,
    entity_id   TEXT,
    summary     TEXT        NOT NULL,
    ip          TEXT,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_audit_logs_project_time ON audit_logs (project_id, created_at DESC);
CREATE INDEX ix_audit_logs_user_time    ON audit_logs (user_id, created_at DESC);
CREATE INDEX ix_audit_logs_action_time  ON audit_logs (action, created_at DESC);
```

**Step 2: Run migration locally**

Run: `just migrate`
Expected: Migration 000028 applied successfully.

**Step 3: Commit**

```bash
git add migrations/000028_create_audit_logs.up.sql
git commit -m "add audit_logs table migration"
```

### Task 3: Audit log domain types

**Files:**
- Create: `internal/domain/audit/entity.go`
- Create: `internal/domain/audit/errors.go`

**Step 1: Define audit domain types**

`internal/domain/audit/entity.go`:

```go
package audit

import "time"

type Entry struct {
	ID         string
	ProjectID  *string
	UserID     string
	Action     string
	EntityType string
	EntityID   *string
	Summary    string
	IP         string
	UserAgent  string
	CreatedAt  time.Time
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

`internal/domain/audit/errors.go`:

```go
package audit

import "errors"

var ErrNotFound = errors.New("audit log entry not found")
```

**Step 2: Commit**

```bash
git add internal/domain/audit/
git commit -m "add audit log domain types"
```

### Task 4: Audit log repository interface + implementation

**Files:**
- Modify: `internal/repository/interfaces.go` — add AuditLogRepository interface
- Create: `internal/repository/audit_repository.go` — Postgres implementation

**Step 1: Add interface to `internal/repository/interfaces.go`**

```go
type AuditLogRepository interface {
	Log(ctx context.Context, entry audit.Entry) error
	List(ctx context.Context, filter audit.Filter) ([]audit.Entry, int, error)
}
```

Add `audit.Entry` import: `"github.com/lbrty/observer/internal/domain/audit"`

Update the `go:generate` line to include `AuditLogRepository`.

**Step 2: Implement in `internal/repository/audit_repository.go`**

```go
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lbrty/observer/internal/domain/audit"
	"github.com/lbrty/observer/internal/ulid"
)

type auditLogRepository struct {
	db *sqlx.DB
}

func NewAuditLogRepository(db *sqlx.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Log(ctx context.Context, entry audit.Entry) error {
	entry.ID = ulid.NewString()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audit_logs (id, project_id, user_id, action, entity_type, entity_id, summary, ip, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		entry.ID, entry.ProjectID, entry.UserID, entry.Action, entry.EntityType, entry.EntityID,
		entry.Summary, entry.IP, entry.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func (r *auditLogRepository) List(ctx context.Context, filter audit.Filter) ([]audit.Entry, int, error) {
	q := `SELECT id, project_id, user_id, action, entity_type, entity_id, summary, ip, user_agent, created_at
	      FROM audit_logs WHERE 1=1`
	countQ := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`
	args := []any{}
	ix := 0

	if filter.ProjectID != nil {
		ix++
		clause := fmt.Sprintf(" AND project_id = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.ProjectID)
	}
	if filter.UserID != nil {
		ix++
		clause := fmt.Sprintf(" AND user_id = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.UserID)
	}
	if filter.Action != nil {
		ix++
		clause := fmt.Sprintf(" AND action = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.Action)
	}
	if filter.EntityType != nil {
		ix++
		clause := fmt.Sprintf(" AND entity_type = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.EntityType)
	}
	if filter.DateFrom != nil {
		ix++
		clause := fmt.Sprintf(" AND created_at >= $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.DateFrom)
	}
	if filter.DateTo != nil {
		ix++
		clause := fmt.Sprintf(" AND created_at <= $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.DateTo)
	}

	var total int
	if err := r.db.GetContext(ctx, &total, countQ, args...); err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	q += " ORDER BY created_at DESC"
	offset := (filter.Page - 1) * filter.PerPage
	ix++
	q += fmt.Sprintf(" LIMIT $%d", ix)
	args = append(args, filter.PerPage)
	ix++
	q += fmt.Sprintf(" OFFSET $%d", ix)
	args = append(args, offset)

	var entries []audit.Entry
	if err := r.db.SelectContext(ctx, &entries, q, args...); err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return entries, total, nil
}
```

**Step 3: Regenerate mocks**

Run: `just generate-mocks`

**Step 4: Run tests**

Run: `just test`
Expected: All pass.

**Step 5: Commit**

```bash
git add internal/repository/interfaces.go internal/repository/audit_repository.go internal/repository/mock/
git commit -m "add audit log repository interface and Postgres implementation"
```

### Task 5: Audit logger use case

**Files:**
- Create: `internal/usecase/audit/audit_usecase.go`
- Create: `internal/usecase/audit/types.go`

**Step 1: Define types in `internal/usecase/audit/types.go`**

```go
package audit

type LogInput struct {
	ProjectID  *string
	UserID     string
	Action     string
	EntityType string
	EntityID   *string
	Summary    string
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
	ID         string  `json:"id"`
	ProjectID  *string `json:"project_id"`
	UserID     string  `json:"user_id"`
	Action     string  `json:"action"`
	EntityType string  `json:"entity_type"`
	EntityID   *string `json:"entity_id"`
	Summary    string  `json:"summary"`
	IP         string  `json:"ip"`
	UserAgent  string  `json:"user_agent"`
	CreatedAt  string  `json:"created_at"`
}

type ListOutput struct {
	Entries []EntryDTO `json:"entries"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
}
```

**Step 2: Implement use case in `internal/usecase/audit/audit_usecase.go`**

```go
package audit

import (
	"context"
	"fmt"
	"time"

	domainaudit "github.com/lbrty/observer/internal/domain/audit"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/usecase"
)

type AuditUseCase struct {
	repo repository.AuditLogRepository
}

func NewAuditUseCase(repo repository.AuditLogRepository) *AuditUseCase {
	return &AuditUseCase{repo: repo}
}

// Log records an audit event. Called from other use cases.
func (uc *AuditUseCase) Log(ctx context.Context, input LogInput) error {
	entry := domainaudit.Entry{
		ProjectID:  input.ProjectID,
		UserID:     input.UserID,
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		Summary:    input.Summary,
		IP:         input.IP,
		UserAgent:  input.UserAgent,
	}
	if err := uc.repo.Log(ctx, entry); err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	return nil
}

// List retrieves paginated audit log entries with filters.
func (uc *AuditUseCase) List(ctx context.Context, input ListInput) (*ListOutput, error) {
	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)
	filter := domainaudit.Filter{
		ProjectID:  input.ProjectID,
		UserID:     input.UserID,
		Action:     input.Action,
		EntityType: input.EntityType,
		Page:       page,
		PerPage:    perPage,
	}

	if input.DateFrom != nil {
		t, err := time.Parse(time.DateOnly, *input.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("parse date_from: %w", err)
		}
		filter.DateFrom = &t
	}
	if input.DateTo != nil {
		t, err := time.Parse(time.DateOnly, *input.DateTo)
		if err != nil {
			return nil, fmt.Errorf("parse date_to: %w", err)
		}
		filter.DateTo = &t
	}

	entries, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	dtos := make([]EntryDTO, len(entries))
	for i, e := range entries {
		dtos[i] = EntryDTO{
			ID:         e.ID,
			ProjectID:  e.ProjectID,
			UserID:     e.UserID,
			Action:     e.Action,
			EntityType: e.EntityType,
			EntityID:   e.EntityID,
			Summary:    e.Summary,
			IP:         e.IP,
			UserAgent:  e.UserAgent,
			CreatedAt:  e.CreatedAt.Format(time.RFC3339),
		}
	}
	return &ListOutput{Entries: dtos, Total: total, Page: page, PerPage: perPage}, nil
}
```

**Step 3: Commit**

```bash
git add internal/usecase/audit/
git commit -m "add audit log use case with Log and List methods"
```

### Task 6: Audit log handler + routes

**Files:**
- Create: `internal/handler/audit_handler.go`
- Modify: `internal/server/server.go` — register audit routes
- Modify: `internal/app/container.go` — wire AuditUseCase

**Step 1: Create handler in `internal/handler/audit_handler.go`**

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

type AuditHandler struct {
	auditUC *ucaudit.AuditUseCase
}

func NewAuditHandler(auditUC *ucaudit.AuditUseCase) *AuditHandler {
	return &AuditHandler{auditUC: auditUC}
}

// ListAll returns audit logs for admins (all projects).
func (h *AuditHandler) ListAll(c *gin.Context) {
	var input ucaudit.ListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.auditUC.List(c.Request.Context(), input)
	if err != nil {
		internalError(c, "list audit logs", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// ListByProject returns audit logs scoped to a project.
func (h *AuditHandler) ListByProject(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucaudit.ListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	input.ProjectID = &projectID
	out, err := h.auditUC.List(c.Request.Context(), input)
	if err != nil {
		internalError(c, "list project audit logs", err)
		return
	}
	c.JSON(http.StatusOK, out)
}
```

**Step 2: Wire in `internal/app/container.go`**

Add `AuditUC *ucaudit.AuditUseCase` field to Container struct. Wire:

```go
auditRepo := repository.NewAuditLogRepository(db.DB())
container.AuditUC = ucaudit.NewAuditUseCase(auditRepo)
```

**Step 3: Register routes in `internal/server/server.go`**

Admin route group:
```go
auditHandler := handler.NewAuditHandler(container.AuditUC)
admin.GET("/audit-logs", auditHandler.ListAll)
```

Project route group (under the manager-level or read-level group — managers+ only):
```go
// Inside the project routes, under a group requiring ActionRead at minimum
// but the handler scopes by project_id automatically
delete.GET("/audit-logs", auditHandler.ListByProject)
```

Note: Use the `delete` group (ActionDelete, manager+) for project audit logs since only managers should see them.

**Step 4: Run tests**

Run: `just test`
Expected: All pass.

**Step 5: Commit**

```bash
git add internal/handler/audit_handler.go internal/app/container.go internal/server/server.go
git commit -m "add audit log handler and routes for admin and project scope"
```

### Task 7: Integrate audit logging into existing use cases

**Files:**
- Modify: `internal/usecase/project/person_usecase.go` — add audit on create/delete
- Modify: `internal/usecase/project/pet_usecase.go` — add audit on create/delete
- Modify: `internal/usecase/project/household_usecase.go` — add audit on create/delete
- Modify: `internal/usecase/project/support_record_usecase.go` — add audit on create/delete
- Modify: `internal/usecase/project/migration_record_usecase.go` — add audit on create/delete
- Modify: `internal/usecase/project/document_usecase.go` — add audit on upload/download/delete
- Modify: `internal/usecase/admin/project_usecase.go` — add audit on create/delete
- Modify: `internal/usecase/admin/permission_usecase.go` — add audit on grant/revoke
- Modify: `internal/usecase/admin/users_usecase.go` — add audit on role change
- Modify: `internal/app/container.go` — inject AuditUseCase into all affected use cases

**Step 1: Add `AuditUseCase` dependency to each use case constructor**

Pattern for each use case — add `auditUC *ucaudit.AuditUseCase` as a constructor parameter and struct field. Then in create/delete/upload/download methods, call:

```go
_ = uc.auditUC.Log(ctx, ucaudit.LogInput{
	ProjectID:  &projectID,
	UserID:     userID,       // extract from context
	Action:     "person.create",
	EntityType: "person",
	EntityID:   &created.ID,
	Summary:    fmt.Sprintf("Created person %s", created.ID),
	IP:         ip,           // extract from context
	UserAgent:  userAgent,    // extract from context
})
```

Audit log failures should not fail the primary operation — use fire-and-forget (log the error but don't return it).

**Step 2: Update DI wiring in `container.go`**

Pass `container.AuditUC` to all modified use case constructors.

**Step 3: Update existing tests**

Add mock AuditUseCase to test constructors. Use `gomock.Any()` for audit calls.

**Step 4: Run tests**

Run: `just test`
Expected: All pass.

**Step 5: Commit**

```bash
git add internal/usecase/ internal/app/container.go
git commit -m "integrate audit logging into record lifecycle, document, admin, and permission use cases"
```

---

## Phase 3: Data Table Filters + CSV Export (Backend)

### Task 8: Extend list filter types

**Files:**
- Modify: `internal/domain/person/entity.go` — extend PersonListFilter
- Modify: `internal/domain/support/entity.go` — extend RecordListFilter
- Modify: `internal/domain/migration/entity.go` — extend RecordListFilter (if exists)
- Modify: `internal/domain/pet/entity.go` — add PetListFilter if not present

**Step 1: Extend filter structs**

Add missing filter fields that the design requires. Example for PersonListFilter:

```go
type PersonListFilter struct {
	// existing fields...
	Sex          *string
	AgeGroup     *string
	CaseStatus   *string
	CategoryID   *string
	ConflictZone *string
	RegionID     *string
	HasPets      *bool
	TagIDs       []string
	Search       *string
}
```

Similarly extend support record, migration record, and pet filters per the design table.

**Step 2: Update repository implementations**

Modify `internal/repository/person_repository.go` (and others) to handle the new filter fields in their `List()` SQL queries using the existing `applyPeopleFilters` / conditional clause pattern.

**Step 3: Run tests**

Run: `just test`
Expected: All pass.

**Step 4: Commit**

```bash
git add internal/domain/ internal/repository/
git commit -m "extend list filters for people, support records, migration records, and pets"
```

### Task 9: CSV export endpoints

**Files:**
- Create: `internal/handler/export_handler.go`
- Modify: `internal/server/server.go` — register export routes under ActionExport
- Modify: `internal/app/container.go` — wire if needed

**Step 1: Create export handler**

The export handler reuses existing list use cases but streams CSV instead of JSON. Pattern:

```go
package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lbrty/observer/internal/middleware"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

type ExportHandler struct {
	personUC *ucproject.PersonUseCase
	petUC    *ucproject.PetUseCase
	// ... other use cases
	auditUC  *ucaudit.AuditUseCase
}

func NewExportHandler(
	personUC *ucproject.PersonUseCase,
	petUC *ucproject.PetUseCase,
	auditUC *ucaudit.AuditUseCase,
) *ExportHandler {
	return &ExportHandler{personUC: personUC, petUC: petUC, auditUC: auditUC}
}

func (h *ExportHandler) ExportPeople(c *gin.Context) {
	projectID := c.Param("project_id")
	// bind filters from query params (same as List)
	// call personUC.List with high perPage to get all matching
	// stream as CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=people-%s.csv", projectID))
	w := csv.NewWriter(c.Writer)
	// write header row
	// write data rows
	w.Flush()

	// audit log
	userID := middleware.UserIDFrom(c)
	_ = h.auditUC.Log(c.Request.Context(), ucaudit.LogInput{
		ProjectID:  &projectID,
		UserID:     userID,
		Action:     "export.people",
		EntityType: "person",
		Summary:    fmt.Sprintf("Exported people from project %s", projectID),
	})
}
```

Repeat pattern for support records, migration records, households, pets.

**Step 2: Register routes**

In `server.go`, add an export route group under ActionExport:

```go
export := projectGroup.Group("")
export.Use(authMw.RequireProjectRole(project.ActionExport))
{
	exportHandler := handler.NewExportHandler(container.PersonUC, container.PetUC, container.AuditUC)
	export.GET("/export/people", exportHandler.ExportPeople)
	export.GET("/export/support-records", exportHandler.ExportSupportRecords)
	export.GET("/export/migration-records", exportHandler.ExportMigrationRecords)
	export.GET("/export/households", exportHandler.ExportHouseholds)
	export.GET("/export/pets", exportHandler.ExportPets)
}
```

**Step 3: Run tests**

Run: `just test`
Expected: All pass.

**Step 4: Commit**

```bash
git add internal/handler/export_handler.go internal/server/server.go internal/app/container.go
git commit -m "add CSV export endpoints gated by ActionExport permission"
```

---

## Phase 4: Custom Report Builder (Backend)

### Task 10: Custom report builder endpoint

**Files:**
- Modify: `internal/usecase/report/types.go` — add CustomReportInput/Output
- Modify: `internal/usecase/report/report_usecase.go` — add GenerateCustom method
- Modify: `internal/repository/interfaces.go` — add CustomQuery method to ReportRepository
- Modify: `internal/repository/report_repository.go` — implement dynamic SQL builder
- Modify: `internal/handler/report_handler.go` — add CustomReport handler
- Modify: `internal/server/server.go` — register route

**Step 1: Define custom report types**

In `internal/usecase/report/types.go`:

```go
type CustomReportInput struct {
	Metric    string   `form:"metric" binding:"required,oneof=events people units pets"`
	GroupBy   []string `form:"group_by" binding:"required,min=1,max=2"`
	DateFrom  *string  `form:"date_from"`
	DateTo    *string  `form:"date_to"`
	Type      *string  `form:"support_type"`
	TagIDs    []string `form:"tag_ids"`
	HasPets   *bool    `form:"has_pets"`
}

type CustomReportOutput struct {
	Metric  string         `json:"metric"`
	GroupBy []string       `json:"group_by"`
	Rows    []CustomRow    `json:"rows"`
	Total   int            `json:"total"`
}

type CustomRow struct {
	Dimensions map[string]string `json:"dimensions"`
	Count      int               `json:"count"`
}
```

**Step 2: Add repository method**

In `ReportRepository` interface:

```go
CustomQuery(ctx context.Context, projectID string, metric string, groupBy []string, filter domainreport.ReportFilter) ([]domainreport.CustomResult, error)
```

Implementation builds SQL dynamically using safe column mapping (whitelist of allowed dimension → SQL expression):

```go
var dimensionSQL = map[string]string{
	"sex":           "p.sex",
	"age_group":     "COALESCE(p.age_group, ...)",
	"region":        "COALESCE(st.name, 'unknown')",
	"conflict_zone": "COALESCE(st.conflict_zone, 'none')",
	"office":        "COALESCE(o.name, 'unknown')",
	"sphere":        "sr.sphere",
	"category":      "c.name",
	"person_tag":    "t.name",
	"pet_tag":       "pt_tag.name",
	"pet_status":    "pets.status",
}
```

No string concatenation from user input — dimension names are validated against the whitelist.

**Step 3: Implement use case method**

`GenerateCustom()` validates input, maps to domain filter, calls repository, returns DTO.

**Step 4: Add handler + route**

```go
read.GET("/reports/custom", reportHandler.GenerateCustom)
```

**Step 5: Run tests**

Run: `just test`
Expected: All pass.

**Step 6: Commit**

```bash
git add internal/usecase/report/ internal/repository/ internal/handler/report_handler.go internal/server/server.go
git commit -m "add custom report builder with dynamic dimension grouping"
```

---

## Phase 5: Frontend — Profile Page

### Task 11: Create profile settings page

**Files:**
- Create: `packages/observer-web/src/routes/_app/settings.tsx` (layout if needed)
- Create: `packages/observer-web/src/routes/_app/settings/profile.tsx`
- Modify: `packages/observer-web/src/routes/_app.tsx` — remove theme/language from AvatarMenu, add profile link
- Modify: `packages/observer-web/src/components/app-footer.tsx` — remove theme/language controls

**Step 1: Create profile settings page**

`packages/observer-web/src/routes/_app/settings/profile.tsx`:

- Display user info (name, email) read-only from auth store
- Theme selection: segmented control with system/light/dark/light-hc/dark-hc options
- Language selection: dropdown using LANGUAGES from constants
- Both read/write localStorage, same keys (THEME_KEY, LANG_KEY)
- Apply theme via `document.documentElement.dataset.theme`
- Apply language via `i18n.changeLanguage()`

**Step 2: Update AvatarMenu in `_app.tsx`**

- Remove inline theme/language dropdowns
- Add "Settings" link pointing to `/settings/profile`
- Keep logout button

**Step 3: Update app-footer.tsx**

- Remove theme/language controls from footer
- Keep any other footer content

**Step 4: Add translations**

Add keys to all 6 locale files: `settings.profile`, `settings.theme`, `settings.language`, etc.

**Step 5: Verify in browser**

Run: `cd packages/observer-web && bun dev`
Navigate to `/settings/profile`, verify theme and language controls work.

**Step 6: Commit**

```bash
git add packages/observer-web/src/
git commit -m "move theme and language selection to profile settings page"
```

---

## Phase 6: Frontend — Data Table Filters

### Task 12: Build generic filter bar component

**Files:**
- Create: `packages/observer-web/src/components/filter-bar.tsx`

**Step 1: Create reusable filter bar**

A composable filter bar that accepts filter field definitions and syncs state with URL search params:

```tsx
interface FilterField {
  key: string
  label: string
  type: "select" | "multi-select" | "date-range" | "search" | "boolean"
  options?: { value: string; label: string }[]
}

interface FilterBarProps {
  fields: FilterField[]
  values: Record<string, any>
  onChange: (values: Record<string, any>) => void
}
```

- Uses URL search params via TanStack Router's `useSearch` / `useNavigate`
- Renders appropriate input for each field type
- Collapse/expand toggle
- "Clear all" button
- Reuses `tag-filter.tsx` pattern for multi-select

**Step 2: Commit**

```bash
git add packages/observer-web/src/components/filter-bar.tsx
git commit -m "add reusable filter bar component with URL param sync"
```

### Task 13: Add filters to people table

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects.$projectId/people/index.tsx`
- Modify: `packages/observer-web/src/hooks/use-people.ts` — pass filter params to API

**Step 1: Add FilterBar to people list page**

Import FilterBar, define fields (sex, age_group, case_status, search, date_range, tags, has_pets), wire to URL search params, pass to `usePeople` hook.

**Step 2: Update hook to pass filters**

Modify `usePeople` to accept filter params and pass them as query params to the API call.

**Step 3: Add export button**

Add CSV export button that calls `/export/people` with current filters. Only visible when user has export permission (check via project role from context).

**Step 4: Verify in browser**

**Step 5: Commit**

```bash
git add packages/observer-web/src/routes/_app/projects.\$projectId/people/ packages/observer-web/src/hooks/use-people.ts
git commit -m "add filter bar and export button to people table"
```

### Task 14: Add filters to remaining tables

**Files:**
- Modify: support records table + hook
- Modify: migration records table + hook
- Modify: pets table + hook
- Modify: households table + hook

**Step 1: Repeat Task 13 pattern for each table**

Each table gets its specific filter fields per the design doc. Each gets an export button.

**Step 2: Verify in browser**

**Step 3: Commit**

```bash
git add packages/observer-web/src/
git commit -m "add filter bars and export buttons to all data tables"
```

---

## Phase 7: Frontend — Custom Report Builder + Audit Log UI

### Task 15: Custom report builder page

**Files:**
- Create: `packages/observer-web/src/routes/_app/projects.$projectId/reports/custom.tsx`
- Modify: `packages/observer-web/src/hooks/use-reports.ts` — add useCustomReport hook
- Modify: `packages/observer-web/src/types/report.ts` — add CustomReport types

**Step 1: Add types and hook**

```typescript
interface CustomReportParams {
  metric: "events" | "people" | "units" | "pets"
  group_by: string[]
  date_from?: string
  date_to?: string
  support_type?: string
  tag_ids?: string[]
  has_pets?: boolean
}

interface CustomRow {
  dimensions: Record<string, string>
  count: number
}

interface CustomReportOutput {
  metric: string
  group_by: string[]
  rows: CustomRow[]
  total: number
}
```

**Step 2: Build the form-based builder page**

- Metric selector (radio group: events/people/units/pets)
- Dimension picker (multi-select, max 2: sex, age_group, region, conflict_zone, office, sphere, category, person_tag, pet_tag, pet_status)
- Filter section (date range, support type, tag multi-select, has_pets toggle)
- Results table rendering CustomRow[] with dynamic columns based on selected dimensions
- Export button (CSV of current results, gated by ActionExport)

**Step 3: Add translations**

**Step 4: Verify in browser**

**Step 5: Commit**

```bash
git add packages/observer-web/src/
git commit -m "add custom report builder page with form-based dimension selection"
```

### Task 16: Admin audit log page

**Files:**
- Create: `packages/observer-web/src/routes/_app/admin/audit-logs.tsx`
- Create: `packages/observer-web/src/hooks/use-audit-logs.ts`
- Create: `packages/observer-web/src/types/audit.ts`

**Step 1: Add types**

```typescript
interface AuditEntry {
  id: string
  project_id: string | null
  user_id: string
  action: string
  entity_type: string
  entity_id: string | null
  summary: string
  ip: string
  user_agent: string
  created_at: string
}

interface AuditListOutput {
  entries: AuditEntry[]
  total: number
  page: number
  per_page: number
}
```

**Step 2: Add hook**

```typescript
export function useAuditLogs(params: AuditListParams) {
  return useQuery({
    queryKey: ["audit-logs", params],
    queryFn: () => api.get("admin/audit-logs", { searchParams: params }).json<AuditListOutput>(),
  })
}
```

**Step 3: Build admin audit log page**

- DataTable with columns: timestamp, user, action, entity type, summary, IP
- FilterBar with: user, action type, entity type, date range, project
- Pagination

**Step 4: Add translations**

**Step 5: Verify in browser**

**Step 6: Commit**

```bash
git add packages/observer-web/src/
git commit -m "add admin audit log page with filters"
```

### Task 17: Project-scoped audit log page

**Files:**
- Create: `packages/observer-web/src/routes/_app/projects.$projectId/audit-logs.tsx`

**Step 1: Build project audit log page**

Reuses `useAuditLogs` hook with `project_id` param. Same DataTable + FilterBar pattern but without project filter (already scoped). Only visible to manager+ role.

**Step 2: Add navigation link**

Add "Audit Log" link to project sidebar/nav for managers+.

**Step 3: Add translations**

**Step 4: Commit**

```bash
git add packages/observer-web/src/
git commit -m "add project-scoped audit log page for managers"
```

---

## Phase 8: Integration Testing

### Task 18: Integration tests for audit logs

**Files:**
- Create: `internal/repository/audit_repository_test.go`

**Step 1: Write integration tests**

Using testcontainers-go (Postgres), test:
- `Log()` inserts entry and returns no error
- `List()` with no filters returns all entries paginated
- `List()` with project_id filter returns scoped entries
- `List()` with date range filter
- `List()` with action filter

Guard with `testing.Short()`.

**Step 2: Run integration tests**

Run: `just test-all`
Expected: All pass.

**Step 3: Commit**

```bash
git add internal/repository/audit_repository_test.go
git commit -m "add integration tests for audit log repository"
```

### Task 19: Integration tests for export endpoints

**Files:**
- Create: `internal/handler/export_handler_test.go`

**Step 1: Write handler tests**

Test CSV export endpoints:
- Returns 200 with CSV content-type and valid CSV body
- Returns 403 when user lacks ActionExport permission
- Respects filter params in exported data
- Creates audit log entry on export

**Step 2: Run tests**

Run: `just test-all`
Expected: All pass.

**Step 3: Commit**

```bash
git add internal/handler/export_handler_test.go
git commit -m "add integration tests for CSV export endpoints"
```

### Task 20: Unit tests for custom report builder

**Files:**
- Create: `internal/usecase/report/report_usecase_test.go` (or extend existing)

**Step 1: Write unit tests**

Using gomock for ReportRepository:
- `GenerateCustom()` with valid single dimension
- `GenerateCustom()` with two dimensions
- `GenerateCustom()` rejects invalid dimension names
- `GenerateCustom()` rejects more than 2 dimensions

**Step 2: Run tests**

Run: `just test`
Expected: All pass.

**Step 3: Commit**

```bash
git add internal/usecase/report/
git commit -m "add unit tests for custom report builder use case"
```
