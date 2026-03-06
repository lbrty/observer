# Remaining Features Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement all missing features: frontend project-scoped pages, auth flows (admin approval, profile, password change/reset), and reports (ADR-005) with D3 visualizations.

**Architecture:** Follow existing DDD + Clean Architecture patterns. Frontend pages follow the people list page pattern (DataTable + drawer/dialog). Backend reports add a new `report` domain + repository + usecase + handler. Auth flows extend existing auth usecase. Registration is open but users start inactive until admin approval.

**Tech Stack:** Go 1.25, Gin, sqlx, PostgreSQL | React 19, TanStack Router/Query, ky, D3, Tailwind, Base UI, i18next


### Task 2: Frontend types for remaining project-scoped resources

**Files:**

- Create: `packages/observer-web/src/types/support-record.ts`
- Create: `packages/observer-web/src/types/migration-record.ts`
- Create: `packages/observer-web/src/types/household.ts`
- Create: `packages/observer-web/src/types/note.ts`
- Create: `packages/observer-web/src/types/document.ts`
- Create: `packages/observer-web/src/types/pet.ts`
- Create: `packages/observer-web/src/types/tag.ts`

**Step 1: Create type files**

Each type file matches the backend DTOs. Pattern from `person.ts`:

```ts
// types/support-record.ts
export type SupportType =
  | "humanitarian"
  | "legal"
  | "social"
  | "psychological"
  | "medical"
  | "general";
export type SupportSphere =
  | "housing_assistance"
  | "document_recovery"
  | "social_benefits"
  | "property_rights"
  | "employment_rights"
  | "family_law"
  | "healthcare_access"
  | "education_access"
  | "financial_aid"
  | "psychological_support"
  | "other";
export type ReferralStatus = "pending" | "accepted" | "completed" | "declined" | "no_response";

export interface SupportRecord {
  id: string;
  person_id: string;
  project_id: string;
  consultant_id?: string;
  recorded_by?: string;
  office_id?: string;
  referred_to_office?: string;
  type: SupportType;
  sphere?: SupportSphere;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateSupportRecordInput {
  person_id: string;
  type: SupportType;
  sphere?: SupportSphere;
  consultant_id?: string;
  office_id?: string;
  referred_to_office?: string;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
}

export interface UpdateSupportRecordInput {
  type?: SupportType;
  sphere?: SupportSphere;
  consultant_id?: string;
  office_id?: string;
  referred_to_office?: string;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
}

export interface ListSupportRecordsParams {
  person_id?: string;
  consultant_id?: string;
  office_id?: string;
  type?: SupportType;
  sphere?: SupportSphere;
  page?: number;
  per_page?: number;
}

export interface ListSupportRecordsOutput {
  records: SupportRecord[];
  total: number;
  page: number;
  per_page: number;
}
```

```ts
// types/migration-record.ts
export type MovementReason =
  | "conflict"
  | "security"
  | "service_access"
  | "return"
  | "relocation_program"
  | "economic"
  | "other";
export type HousingAtDestination =
  | "own_property"
  | "renting"
  | "with_relatives"
  | "collective_site"
  | "hotel"
  | "other"
  | "unknown";

export interface MigrationRecord {
  id: string;
  person_id: string;
  from_place_id?: string;
  destination_place_id?: string;
  migration_date?: string;
  movement_reason?: MovementReason;
  housing_at_destination?: HousingAtDestination;
  notes?: string;
  created_at: string;
}

export interface CreateMigrationRecordInput {
  from_place_id?: string;
  destination_place_id?: string;
  migration_date?: string;
  movement_reason?: MovementReason;
  housing_at_destination?: HousingAtDestination;
  notes?: string;
}
```

```ts
// types/household.ts
export type Relationship =
  | "head"
  | "spouse"
  | "child"
  | "parent"
  | "sibling"
  | "grandchild"
  | "grandparent"
  | "other_relative"
  | "non_relative";

export interface Household {
  id: string;
  project_id: string;
  reference_number?: string;
  head_person_id?: string;
  members?: HouseholdMember[];
  created_at: string;
  updated_at: string;
}

export interface HouseholdMember {
  person_id: string;
  relationship: string;
}

export interface CreateHouseholdInput {
  reference_number?: string;
  head_person_id?: string;
}

export interface UpdateHouseholdInput {
  reference_number?: string;
  head_person_id?: string;
}

export interface AddMemberInput {
  person_id: string;
  relationship: Relationship;
}

export interface ListHouseholdsParams {
  page?: number;
  per_page?: number;
}

export interface ListHouseholdsOutput {
  households: Household[];
  total: number;
  page: number;
  per_page: number;
}
```

```ts
// types/note.ts
export interface Note {
  id: string;
  person_id: string;
  author_id?: string;
  body: string;
  created_at: string;
}

export interface CreateNoteInput {
  body: string;
}
```

```ts
// types/document.ts
export interface Document {
  id: string;
  person_id: string;
  project_id: string;
  uploaded_by?: string;
  name: string;
  path: string;
  mime_type: string;
  size: number;
  created_at: string;
}

export interface CreateDocumentInput {
  person_id: string;
  name: string;
  path: string;
  mime_type: string;
  size: number;
}
```

```ts
// types/pet.ts
export type PetStatus = "registered" | "adopted" | "owner_found" | "needs_shelter" | "unknown";

export interface Pet {
  id: string;
  project_id: string;
  owner_id?: string;
  name: string;
  status: PetStatus;
  registration_id?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePetInput {
  owner_id?: string;
  name: string;
  status?: PetStatus;
  registration_id?: string;
  notes?: string;
}

export interface UpdatePetInput {
  owner_id?: string;
  name?: string;
  status?: PetStatus;
  registration_id?: string;
  notes?: string;
}

export interface ListPetsParams {
  page?: number;
  per_page?: number;
}

export interface ListPetsOutput {
  pets: Pet[];
  total: number;
  page: number;
  per_page: number;
}
```

```ts
// types/tag.ts
export interface Tag {
  id: string;
  project_id: string;
  name: string;
  created_at: string;
}

export interface CreateTagInput {
  name: string;
}
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bun run build`


### Task 4: Tags management page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/tags/index.tsx`

Simple list + create dialog + delete confirmation. No edit (tags are name-only, just delete and recreate).

**Step 1: Create tags page**

Follow the countries reference page pattern — DataTable + FormDialog for create + ConfirmDialog for delete.

Columns: name, created_at, actions (delete only).

**Step 2: Add i18n keys**

Add `project.tags.*` keys to all 4 locale files:

```json
{
  "project": {
    "tags": {
      "title": "Tegder",
      "add": "Teg koshuu",
      "name": "Atalyşy",
      "deleteConfirm": "Tegdi jokko çygarasyzby?",
      "deleted": "Teg öçürüldü"
    }
  }
}
```

**Step 3: Verify and commit**

Run: `cd packages/observer-web && bun run build`


### Task 6: Households page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/households/index.tsx`
- Create: `packages/observer-web/src/components/household-drawer.tsx`

CRUD with members management. Drawer shows household details + member list with add/remove.

**Step 1: Create household drawer**

Form sections:

- Household info: reference_number (text), head_person_id (person search/select)
- Members list: table of current members with relationship + remove button
- Add member form: person_id (search/select), relationship (select)

**Step 2: Create households list page**

DataTable columns: reference_number, head person name, member count, created_at, actions.
Pagination.

**Step 3: Add i18n keys for households**

Add `project.households.*` keys including relationship type labels.

**Step 4: Verify and commit**

Run: `cd packages/observer-web && bun run build`


### Task 8: Pets page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/pets/index.tsx`
- Create: `packages/observer-web/src/components/pet-drawer.tsx`

Full CRUD with drawer. Simple entity.

**Step 1: Create pet drawer**

Form fields: name, status (select), owner_id (person search/select), registration_id, notes.

**Step 2: Create pets list page**

DataTable columns: name, status (badge), owner, registration_id, actions.
Tab filtering by status (all, registered, adopted, etc.).

**Step 3: Add i18n keys**

Add `project.pets.*` keys including status labels.

**Step 4: Verify and commit**

Run: `cd packages/observer-web && bun run build`


## Phase 2: Backend + Frontend — Auth Flows

### Task 10: Registration approval flow (users start inactive)

Registration is open, but users start as inactive (`is_active = false`). Admins approve users by activating them via `PATCH /admin/users/:id`. No email verification needed — admins vouch for users directly.

**Files:**

- Modify: `internal/usecase/auth/register_usecase.go` — set `IsActive: false`
- Modify: `packages/observer-web/src/stores/auth.tsx` — handle inactive user error on login
- Modify: `packages/observer-web/src/routes/login.tsx` — show "pending approval" message

**Step 1: Update register usecase**

Change `IsActive: true` to `IsActive: false` and update the success message:

```go
// In register_usecase.go, change:
newUser := &user.User{
    // ...
    IsActive:   false,  // was: true — admin must approve
    // ...
}

// Update response message:
return &RegisterOutput{
    Message: "Registration successful. Your account is pending admin approval.",
}, nil
```

**Step 2: Update frontend login error handling**

When login returns `errors.auth.userNotActive`, show a clear message like "Your account is pending admin approval" instead of a generic error.

**Step 3: Add admin user activation UI**

The admin user edit page already has an `is_active` toggle. Add a visual indicator on the admin users list for pending (inactive) users — e.g., a "Pending" badge using the `gold` variant.

**Step 4: Add i18n keys**

```json
{
  "auth": {
    "pendingApproval": "Sizdin akkauntunuz admin tarabyndan tastyktaluunu küttö"
  }
}
```

**Step 5: Verify and commit**

Run: `just test && cd packages/observer-web && bun run build`


### Task 12: Password change endpoint (authenticated)

**Files:**

- Create: `internal/usecase/auth/change_password_usecase.go`
- Modify: `internal/handler/auth_handler.go` — add ChangePassword method
- Modify: `internal/server/server.go` — add `POST /auth/change-password` route
- Modify: `internal/app/container.go` — wire ChangePasswordUseCase

**Step 1: Create change password usecase**

```go
// internal/usecase/auth/change_password_usecase.go
package auth

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/crypto"
	domainuser "github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type ChangePasswordUseCase struct {
	userRepo repository.UserRepository
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

func NewChangePasswordUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
		credRepo: credRepo,
		hasher:   hasher,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, userID string, input ChangePasswordInput) error {
	cred, err := uc.credRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get credentials: %w", err)
	}

	if !uc.hasher.Verify(input.CurrentPassword, cred.PasswordHash) {
		return domainuser.ErrInvalidCredentials
	}

	hash, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	cred.PasswordHash = hash
	if err := uc.credRepo.Update(ctx, cred); err != nil {
		return fmt.Errorf("update credentials: %w", err)
	}

	return nil
}
```

**Step 2: Add handler method and route**

`POST /auth/change-password` — authenticated, extracts userID, returns 200 on success.

**Step 3: Add frontend change password form**

Add a "Change Password" section to the profile page with current_password, new_password, confirm_password fields.

**Step 4: Add i18n keys and verify**

Run: `just test`


## Phase 3: Reports (ADR-005) with D3 Visualizations

### Task 14: Report domain types

**Files:**

- Create: `internal/domain/report/types.go`

**Step 1: Define report request/response types**

```go
package report

import "time"

// ReportFilter contains common filter parameters for all reports.
type ReportFilter struct {
	ProjectID string
	DateFrom  *time.Time
	DateTo    *time.Time
}

// CountResult represents a single count in a report breakdown.
type CountResult struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// GroupedReport contains a title and breakdown rows.
type GroupedReport struct {
	Title   string        `json:"title"`
	Total   int           `json:"total"`
	Rows    []CountResult `json:"rows"`
}
```

**Step 2: Verify and commit**

Run: `just test`


### Task 16: Report repository SQL implementation

**Files:**

- Create: `internal/repository/report_repository.go`

**Step 1: Implement report SQL queries**

Each method builds a query following ADR-005 patterns. Example for Group 1:

```go
func (r *reportRepository) CountConsultations(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	query := `
		SELECT
			sr.type AS label,
			COUNT(*) AS count
		FROM support_records sr
		WHERE sr.project_id = $1
	`
	args := []any{f.ProjectID}
	ix := 2

	if f.DateFrom != nil {
		query += fmt.Sprintf(" AND sr.provided_at >= $%d", ix)
		args = append(args, *f.DateFrom)
		ix++
	}
	if f.DateTo != nil {
		query += fmt.Sprintf(" AND sr.provided_at <= $%d", ix)
		args = append(args, *f.DateTo)
		ix++
	}

	query += " GROUP BY sr.type ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, fmt.Errorf("count consultations: %w", err)
	}
	return results, nil
}
```

Group 3 (IDP status) joins through places → states:

```sql
SELECT s.conflict_zone AS label, COUNT(DISTINCT p.id) AS count
FROM people p
JOIN places pl ON p.origin_place_id = pl.id
JOIN states s ON pl.state_id = s.id
WHERE p.project_id = $1
GROUP BY s.conflict_zone
```

Group 8 (age group) uses CASE for birth_date bucketing:

```sql
SELECT
  CASE
    WHEN age < 1 THEN 'infant'
    WHEN age < 3 THEN 'toddler'
    WHEN age < 6 THEN 'pre_school'
    WHEN age < 12 THEN 'middle_childhood'
    WHEN age < 14 THEN 'young_teen'
    WHEN age < 18 THEN 'teenager'
    WHEN age < 25 THEN 'young_adult'
    WHEN age < 35 THEN 'early_adult'
    WHEN age < 55 THEN 'middle_aged_adult'
    ELSE 'old_adult'
  END AS label,
  COUNT(*) AS count
FROM (
  SELECT EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) AS age
  FROM people p WHERE p.project_id = $1 AND p.birth_date IS NOT NULL
) sub
GROUP BY label ORDER BY label
```

Group 10 (family units) counts distinct households:

```sql
SELECT 'households' AS label, COUNT(DISTINCT h.id) AS count
FROM households h
JOIN household_members hm ON h.id = hm.household_id
JOIN people p ON hm.person_id = p.id
JOIN support_records sr ON sr.person_id = p.id
WHERE h.project_id = $1
```

**Step 2: Verify and commit**

Run: `just test`


### Task 18: Report handler and routes

**Files:**

- Create: `internal/handler/report_handler.go`
- Modify: `internal/server/server.go` — add report routes
- Modify: `internal/app/container.go` — wire ReportUseCase + ReportRepository

**Step 1: Create report handler**

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucreport "github.com/lbrty/observer/internal/usecase/report"
)

type ReportHandler struct {
	uc *ucreport.ReportUseCase
}

func NewReportHandler(uc *ucreport.ReportUseCase) *ReportHandler {
	return &ReportHandler{uc: uc}
}

// Generate handles GET /projects/:project_id/reports.
func (h *ReportHandler) Generate(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucreport.ReportInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Generate(c.Request.Context(), projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, out)
}
```

**Step 2: Add route**

In `server.go`, in the project read group:

```go
readGroup.GET("/reports", reportHandler.Generate)
```

**Step 3: Wire in container**

Add `ReportRepository` and `ReportUseCase` to container.

**Step 4: Verify and commit**

Run: `just test`


## Phase 4: Deferred Features (Documented for Future)

### Task 20: Document future work

**Files:**

- Create: `docs/adr/009-deferred-features.md`

Document the following as future work:

1. **Email service integration** — SMTP client for sending notifications (optional, not required for core flow)
2. **MFA/TOTP verification** — TOTP library integration, enrollment flow, verification endpoint
3. **Document file upload/download** — Storage abstraction (local FS / S3), multipart upload, binary download, MEK/DEK encryption
4. **Data export/import** — CSV/JSON batch export of people + support records, batch import with deduplication by external_id
5. **Advanced search/filtering** — Date range filters on list endpoints, multi-field sorting, saved filters
6. **Go → TypeScript type generation** — tygo integration per ADR-008

