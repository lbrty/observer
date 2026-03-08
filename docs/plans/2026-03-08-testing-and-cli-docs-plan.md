# Testing & CLI Documentation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Comprehensive test coverage for backend handlers, missing use cases, repository integration tests, frontend hooks/components/routes, plus CLI documentation.

**Architecture:** Backend tests use httptest + gomock for handlers, gomock + testify for use cases, testcontainers-go for repo integration. Frontend tests use bun:test + @testing-library/react + mock.module(). CLI docs combine improved cobra help text with a standalone markdown reference.

**Tech Stack:** Go (testify, gomock, httptest, testcontainers-go), TypeScript (bun:test, @testing-library/react, Happy DOM), cobra (CLI)

---

## Phase 1: Backend Handler Test Infrastructure

### Task 1: Create handler test helper

**Files:**
- Create: `internal/handler/testhelpers_test.go`

**Step 1: Write the test helper**

Create a shared test helper for all handler tests. This sets up Gin in test mode with a mock context, common test fixtures, and helper functions.

```go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func newTestContextWithParams(method, path string, body any, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := newTestContext(method, path, body)
	c.Params = params
	return c, w
}

func setAuthContext(c *gin.Context, userID ulid.ULID) {
	c.Set(middleware.UserIDKey, userID)
}

func testID() ulid.ULID {
	return ulid.Make()
}

func parseResponse[T any](w *httptest.ResponseRecorder) T {
	var result T
	_ = json.Unmarshal(w.Body.Bytes(), &result)
	return result
}
```

**Step 2: Verify it compiles**

Run: `cd /Users/sultan/Projects/observer && go build ./internal/handler/...`
Expected: no errors

**Step 3: Commit**

```bash
git add internal/handler/testhelpers_test.go
git commit -m "add handler test helpers for httptest + gin test mode"
```

---

## Phase 2: Backend Handler Tests

Each handler follows the same pattern: mock the use case, call the handler method, assert status code + response body. Error paths first, then happy paths.

### Task 2: Auth handler tests

**Files:**
- Create: `internal/handler/auth_handler_test.go`
- Reference: `internal/handler/auth_handler.go`

**Step 1: Write tests**

Test cases (error-path heavy):
- `TestAuthHandler_Register_ValidationError` — missing/invalid fields → 400
- `TestAuthHandler_Register_DuplicateEmail` — use case returns ErrEmailTaken → 409
- `TestAuthHandler_Register_Success` — valid input → 201
- `TestAuthHandler_Login_ValidationError` — missing email/password → 400
- `TestAuthHandler_Login_InvalidCredentials` — use case returns ErrInvalidCredentials → 401
- `TestAuthHandler_Login_InactiveUser` — use case returns ErrUserNotActive → 403
- `TestAuthHandler_Login_Success` — valid creds → 200 + cookies set
- `TestAuthHandler_Me_NoAuth` — no user in context → 401
- `TestAuthHandler_Me_Success` — user in context → 200
- `TestAuthHandler_RefreshToken_NoCookie` — no refresh cookie → 401
- `TestAuthHandler_RefreshToken_InvalidToken` — use case returns error → 401
- `TestAuthHandler_RefreshToken_Success` — valid token → 200 + new cookies
- `TestAuthHandler_Logout_NoCookie` — no refresh cookie → still 204 (idempotent)
- `TestAuthHandler_Logout_Success` — valid cookie → 204
- `TestAuthHandler_UpdateProfile_ValidationError` — invalid input → 400
- `TestAuthHandler_UpdateProfile_Success` — valid input → 200
- `TestAuthHandler_ChangePassword_ValidationError` — mismatched passwords → 400
- `TestAuthHandler_ChangePassword_WrongCurrent` — use case returns error → 400
- `TestAuthHandler_ChangePassword_Success` — valid → 204

The handler depends on `*ucauth.AuthUseCase` which is a concrete struct, not an interface. We need to mock the underlying repositories it uses. Check how `AuthUseCase` is constructed and mock its dependencies (UserRepository, CredentialsRepository, SessionRepository, MFARepository, PasswordHasher, TokenGenerator).

**Step 2: Run tests**

Run: `go test ./internal/handler/ -run TestAuthHandler -v`
Expected: all PASS

**Step 3: Commit**

```bash
git add internal/handler/auth_handler_test.go
git commit -m "add auth handler tests with error path coverage"
```

### Task 3: Admin handler tests

**Files:**
- Create: `internal/handler/admin_handler_test.go`
- Reference: `internal/handler/admin_handler.go`

**Step 1: Write tests**

Test cases:
- `TestAdminHandler_ListUsers_Success` — returns user list → 200
- `TestAdminHandler_ListUsers_InternalError` — use case fails → 500
- `TestAdminHandler_GetUser_NotFound` — unknown ID → 404
- `TestAdminHandler_GetUser_InvalidID` — bad ULID → 400
- `TestAdminHandler_GetUser_Success` → 200
- `TestAdminHandler_CreateUser_ValidationError` — invalid input → 400
- `TestAdminHandler_CreateUser_DuplicateEmail` — email taken → 409
- `TestAdminHandler_CreateUser_Success` → 201
- `TestAdminHandler_UpdateUser_NotFound` → 404
- `TestAdminHandler_UpdateUser_Success` → 200
- `TestAdminHandler_ResetPassword_ValidationError` → 400
- `TestAdminHandler_ResetPassword_Success` → 204

**Step 2: Run and verify**

Run: `go test ./internal/handler/ -run TestAdminHandler -v`

**Step 3: Commit**

```bash
git add internal/handler/admin_handler_test.go
git commit -m "add admin handler tests"
```

### Task 4: Reference handler tests (country, state, place, office, category)

**Files:**
- Create: `internal/handler/country_handler_test.go`
- Create: `internal/handler/state_handler_test.go`
- Create: `internal/handler/place_handler_test.go`
- Create: `internal/handler/office_handler_test.go`
- Create: `internal/handler/category_handler_test.go`

**Step 1: Write tests**

All reference handlers follow the same CRUD pattern. For each handler, test:
- `List_Success` → 200
- `List_InternalError` → 500
- `Get_NotFound` → 404
- `Get_InvalidID` → 400
- `Get_Success` → 200
- `Create_ValidationError` → 400
- `Create_DuplicateName` → 409 (unique constraint)
- `Create_Success` → 201
- `Update_NotFound` → 404
- `Update_ValidationError` → 400
- `Update_Success` → 200
- `Delete_NotFound` → 404
- `Delete_Success` → 204

State and Place have extra scoped list methods (ListByCountry, ListByState). Test those too.

**Step 2: Run**

Run: `go test ./internal/handler/ -run "TestCountryHandler|TestStateHandler|TestPlaceHandler|TestOfficeHandler|TestCategoryHandler" -v`

**Step 3: Commit**

```bash
git add internal/handler/{country,state,place,office,category}_handler_test.go
git commit -m "add reference data handler tests (country, state, place, office, category)"
```

### Task 5: Project handler tests

**Files:**
- Create: `internal/handler/project_handler_test.go`

**Step 1: Write tests**

- `TestProjectHandler_List_Success` → 200
- `TestProjectHandler_Get_NotFound` → 404
- `TestProjectHandler_Get_Success` → 200
- `TestProjectHandler_Create_ValidationError` → 400
- `TestProjectHandler_Create_Success` → 201
- `TestProjectHandler_Update_NotFound` → 404
- `TestProjectHandler_Update_Success` → 200

**Step 2: Run and commit**

### Task 6: Permission handler tests

**Files:**
- Create: `internal/handler/permission_handler_test.go`

**Step 1: Write tests**

- `TestPermissionHandler_List_Success` → 200
- `TestPermissionHandler_Assign_ValidationError` → 400
- `TestPermissionHandler_Assign_UserNotFound` → 404
- `TestPermissionHandler_Assign_Success` → 201
- `TestPermissionHandler_Update_NotFound` → 404
- `TestPermissionHandler_Update_Success` → 200
- `TestPermissionHandler_Revoke_NotFound` → 404
- `TestPermissionHandler_Revoke_Success` → 204

**Step 2: Run and commit**

### Task 7: Tag handler tests

**Files:**
- Create: `internal/handler/tag_handler_test.go`

**Step 1: Write tests**

- `TestTagHandler_List_Success` → 200
- `TestTagHandler_Create_ValidationError` → 400
- `TestTagHandler_Create_DuplicateName` → 409
- `TestTagHandler_Create_Success` → 201
- `TestTagHandler_Update_NotFound` → 404
- `TestTagHandler_Update_Success` → 200
- `TestTagHandler_Delete_NotFound` → 404
- `TestTagHandler_Delete_Success` → 204

**Step 2: Run and commit**

### Task 8: Person handler tests

**Files:**
- Create: `internal/handler/person_handler_test.go`

**Step 1: Write tests**

PersonHandler has 3 use cases (personUC, categoryUC, tagUC). Tests:
- `TestPersonHandler_List_Success` → 200
- `TestPersonHandler_List_WithFilters` → query params → 200
- `TestPersonHandler_Get_NotFound` → 404
- `TestPersonHandler_Get_Success` → 200
- `TestPersonHandler_Create_ValidationError` → 400
- `TestPersonHandler_Create_Success` → 201
- `TestPersonHandler_Update_NotFound` → 404
- `TestPersonHandler_Update_Success` → 200
- `TestPersonHandler_Delete_NotFound` → 404
- `TestPersonHandler_Delete_Success` → 204
- `TestPersonHandler_ListCategories_Success` → 200
- `TestPersonHandler_ReplaceCategories_Success` → 200
- `TestPersonHandler_ListTags_Success` → 200
- `TestPersonHandler_ReplaceTags_Success` → 200

**Step 2: Run and commit**

### Task 9: Support record handler tests

**Files:**
- Create: `internal/handler/support_record_handler_test.go`

**Step 1: Write tests**

- `TestSupportRecordHandler_List_Success` → 200
- `TestSupportRecordHandler_Get_NotFound` → 404
- `TestSupportRecordHandler_Get_Success` → 200
- `TestSupportRecordHandler_Create_ValidationError` → 400
- `TestSupportRecordHandler_Create_Success` → 201
- `TestSupportRecordHandler_Update_NotFound` → 404
- `TestSupportRecordHandler_Update_Success` → 200
- `TestSupportRecordHandler_Delete_NotFound` → 404
- `TestSupportRecordHandler_Delete_Success` → 204

**Step 2: Run and commit**

### Task 10: Migration record handler tests

**Files:**
- Create: `internal/handler/migration_record_handler_test.go`

Same CRUD pattern as support records.

### Task 11: Household handler tests

**Files:**
- Create: `internal/handler/household_handler_test.go`

CRUD + AddMember/RemoveMember:
- Standard CRUD tests (list, get, create, update, delete)
- `TestHouseholdHandler_AddMember_ValidationError` → 400
- `TestHouseholdHandler_AddMember_Success` → 201
- `TestHouseholdHandler_RemoveMember_NotFound` → 404
- `TestHouseholdHandler_RemoveMember_Success` → 204

### Task 12: Note handler tests

**Files:**
- Create: `internal/handler/note_handler_test.go`

Standard CRUD pattern.

### Task 13: Document handler tests

**Files:**
- Create: `internal/handler/document_handler_test.go`

CRUD + Upload/Download:
- Standard CRUD tests
- `TestDocumentHandler_Upload_NoFile` → 400
- `TestDocumentHandler_Upload_Success` → 201
- `TestDocumentHandler_Download_NotFound` → 404
- `TestDocumentHandler_Download_Success` → 200 + file body
- `TestDocumentHandler_Thumbnail_NotFound` → 404
- `TestDocumentHandler_Thumbnail_Success` → 200

### Task 14: Pet handler tests

**Files:**
- Create: `internal/handler/pet_handler_test.go`

CRUD + tag management (similar to person handler with tags).

### Task 15: Report handler tests

**Files:**
- Create: `internal/handler/report_handler_test.go`
- Create: `internal/handler/pet_report_handler_test.go`

- `TestReportHandler_Generate_ValidationError` → 400
- `TestReportHandler_Generate_Success` → 200
- `TestPetReportHandler_Generate_Success` → 200

### Task 16: My handler tests

**Files:**
- Create: `internal/handler/my_handler_test.go`

- `TestMyHandler_ListProjects_Success` → 200
- `TestMyHandler_ListProjects_InternalError` → 500

**Step: Run all handler tests**

Run: `go test ./internal/handler/ -v`
Expected: all PASS

**Commit:**
```bash
git add internal/handler/*_test.go
git commit -m "add handler tests for all remaining handlers"
```

---

## Phase 3: Backend Use Case Tests (fill gaps)

### Task 17: Tag use case tests

**Files:**
- Create: `internal/usecase/project/tag_usecase_test.go`

**Step 1: Write tests**

```go
// Pattern: mock TagRepository, test each method
// Error paths:
// - List: repo error → propagated
// - Create: duplicate name → ErrTagAlreadyExists
// - Update: not found → ErrTagNotFound
// - Delete: not found → ErrTagNotFound
// Happy paths:
// - List: returns tags
// - Create: returns new tag
// - Update: returns updated tag
// - Delete: success
```

**Step 2: Run**

Run: `go test ./internal/usecase/project/ -run TestTagUseCase -v`

**Step 3: Commit**

### Task 18: Person category use case tests

**Files:**
- Create: `internal/usecase/project/person_category_usecase_test.go`

Tests:
- `TestPersonCategoryUseCase_List_Success`
- `TestPersonCategoryUseCase_List_RepoError`
- `TestPersonCategoryUseCase_Replace_Success`
- `TestPersonCategoryUseCase_Replace_RepoError`

### Task 19: Person tag use case tests

**Files:**
- Create: `internal/usecase/project/person_tag_usecase_test.go`

Same pattern as person category.

### Task 20: Support record use case tests

**Files:**
- Create: `internal/usecase/project/support_record_usecase_test.go`

Tests:
- `TestSupportRecordUseCase_List_Success`
- `TestSupportRecordUseCase_List_RepoError`
- `TestSupportRecordUseCase_Get_Success`
- `TestSupportRecordUseCase_Get_NotFound`
- `TestSupportRecordUseCase_Create_Success`
- `TestSupportRecordUseCase_Create_RepoError`
- `TestSupportRecordUseCase_Update_Success`
- `TestSupportRecordUseCase_Update_NotFound`
- `TestSupportRecordUseCase_Delete_Success`
- `TestSupportRecordUseCase_Delete_NotFound`

### Task 21: Household use case tests

**Files:**
- Create: `internal/usecase/project/household_usecase_test.go`

CRUD + member management:
- Standard CRUD tests
- `TestHouseholdUseCase_AddMember_Success`
- `TestHouseholdUseCase_AddMember_RepoError`
- `TestHouseholdUseCase_RemoveMember_Success`
- `TestHouseholdUseCase_RemoveMember_NotFound`

### Task 22: Pet use case tests

**Files:**
- Create: `internal/usecase/project/pet_usecase_test.go`

Standard CRUD pattern + tag interactions.

### Task 23: Pet tag use case tests

**Files:**
- Create: `internal/usecase/project/pet_tag_usecase_test.go`

Same pattern as person tag.

### Task 24: Migration record use case tests

**Files:**
- Create: `internal/usecase/project/migration_record_usecase_test.go`

CRUD pattern — ListByPerson, Get, Create, Update.

### Task 25: Report use case tests

**Files:**
- Create: `internal/usecase/report/report_usecase_test.go`
- Create: `internal/usecase/report/pet_report_usecase_test.go`

Test Generate method with various inputs and repo errors.

### Task 26: My projects use case tests

**Files:**
- Create: `internal/usecase/my/projects_usecase_test.go`

- `TestMyProjectsUseCase_Execute_Admin` — returns all projects
- `TestMyProjectsUseCase_Execute_User` — returns only permitted projects
- `TestMyProjectsUseCase_Execute_RepoError` — propagates error

### Task 27: Remaining admin use case tests

Check what's already tested. Fill gaps for:
- `state_usecase_test.go` (if missing)
- `place_usecase_test.go` (if missing)
- `office_usecase_test.go` (if missing)
- `category_usecase_test.go` (if missing — check, country exists)
- `permission_usecase_test.go` (if missing — assign exists, check full coverage)
- `user_usecase_test.go` (list/update exist, check Create/ResetPassword)

**Run all use case tests:**

Run: `go test ./internal/usecase/... -v`
Expected: all PASS

**Commit:**
```bash
git add internal/usecase/**/*_test.go
git commit -m "add missing use case tests for tag, person category/tag, support, household, pet, migration, reports, my projects"
```

---

## Phase 4: Backend Repository Integration Tests

### Task 28: Create repository integration test infrastructure

**Files:**
- Create: `internal/repository/integration_test.go`

**Step 1: Write test infrastructure**

```go
//go:build !short

package repository_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("observer_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}

	t.Cleanup(func() { pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	// Run migrations
	runMigrations(t, db)

	return db
}
```

Include a `runMigrations` helper that applies all migration files from `../../migrations/`.

**Step 2: Verify it starts**

Run: `go test ./internal/repository/ -run TestSetup -v` (without -short)

**Step 3: Commit**

### Task 29: User repository integration tests

**Files:**
- Create: `internal/repository/user_repository_integration_test.go`

Tests:
- Create user → read back
- Create duplicate email → unique violation
- GetByEmail → found / not found
- Update user → read back updated fields
- List with filter

### Task 30: Auth repository integration tests

**Files:**
- Create: `internal/repository/auth_repository_integration_test.go`

Tests:
- Create session → get by refresh token
- Delete session
- Create credentials → get by user ID
- Update credentials

### Task 31: Reference repository integration tests

**Files:**
- Create: `internal/repository/reference_repository_integration_test.go`

Tests for country, state, place, office, category repos:
- CRUD cycle
- Cascade delete (delete country → states deleted)
- Unique name constraints

### Task 32: Project & permission repository integration tests

**Files:**
- Create: `internal/repository/project_repository_integration_test.go`

Tests:
- Project CRUD
- Permission assignment / revocation
- List projects filtered by permission

### Task 33: Person, tag, support, migration, household, note, document, pet repository integration tests

**Files:**
- Create: `internal/repository/person_repository_integration_test.go`
- Create: `internal/repository/records_repository_integration_test.go`

Test full CRUD cycles with real Postgres. Focus on:
- Foreign key constraints
- Pagination
- Filter combinations
- Cascade deletes

**Run all integration tests:**

Run: `go test ./internal/repository/ -v` (no -short flag, Docker required)

**Commit:**
```bash
git add internal/repository/*_integration_test.go internal/repository/integration_test.go
git commit -m "add repository integration tests with testcontainers"
```

---

## Phase 5: Frontend Hook Tests

### Task 34: Create hook test factory

Many hooks follow the same pattern (CRUD via ky). Create a helper to reduce boilerplate.

**Files:**
- Create: `packages/observer-web/src/test/hook-helpers.ts`

```typescript
import { mock } from "bun:test";

export function mockApi(overrides: Record<string, () => any> = {}) {
  const defaults = {
    get: () => ({ json: () => Promise.resolve({}) }),
    post: () => ({ json: () => Promise.resolve({}) }),
    patch: () => ({ json: () => Promise.resolve({}) }),
    put: () => ({ json: () => Promise.resolve({}) }),
    delete: () => Promise.resolve(),
  };

  mock.module("@/lib/api", () => ({
    api: { ...defaults, ...overrides },
  }));
}

export function mockApiError(method: string, status: number) {
  mock.module("@/lib/api", () => ({
    api: {
      get: () => ({ json: () => Promise.resolve({}) }),
      post: () => ({ json: () => Promise.resolve({}) }),
      patch: () => ({ json: () => Promise.resolve({}) }),
      delete: () => Promise.resolve(),
      [method]: () => {
        throw { response: { status } };
      },
    },
  }));
}
```

**Commit after creating**

### Task 35: Test use-people hook

**Files:**
- Create: `packages/observer-web/src/hooks/use-people.test.ts`

Test pattern (matching existing use-notes.test.ts):
- `usePeople` — fetches people list, handles empty, handles enabled=false
- `useCreatePerson` — returns mutation
- `useUpdatePerson` — returns mutation
- `useDeletePerson` — returns mutation

### Task 36: Test use-tags hook

**Files:**
- Create: `packages/observer-web/src/hooks/use-tags.test.ts`

### Task 37: Test use-pets hook

**Files:**
- Create: `packages/observer-web/src/hooks/use-pets.test.ts`

### Task 38: Test use-households hook

**Files:**
- Create: `packages/observer-web/src/hooks/use-households.test.ts`

### Task 39: Test use-support-records hook

**Files:**
- Create: `packages/observer-web/src/hooks/use-support-records.test.ts`

### Task 40: Test use-projects hook

**Files:**
- Create: `packages/observer-web/src/hooks/use-projects.test.ts`

### Task 41: Test use-users hook

**Files:**
- Create: `packages/observer-web/src/hooks/use-users.test.ts`

### Task 42: Test remaining hooks

**Files:**
- Create: `packages/observer-web/src/hooks/use-categories.test.ts`
- Create: `packages/observer-web/src/hooks/use-countries.test.ts`
- Create: `packages/observer-web/src/hooks/use-states.test.ts`
- Create: `packages/observer-web/src/hooks/use-places.test.ts`
- Create: `packages/observer-web/src/hooks/use-offices.test.ts`
- Create: `packages/observer-web/src/hooks/use-permissions.test.ts`
- Create: `packages/observer-web/src/hooks/use-my-projects.test.ts`
- Create: `packages/observer-web/src/hooks/use-reports.test.ts`
- Create: `packages/observer-web/src/hooks/use-pet-reports.test.ts`

Each follows the same mock.module + renderHook + waitFor pattern.

**Run all hook tests:**

Run: `cd packages/observer-web && bun test src/hooks/`
Expected: all PASS

**Commit:**
```bash
git add packages/observer-web/src/hooks/*.test.ts packages/observer-web/src/test/hook-helpers.ts
git commit -m "add tests for all frontend hooks"
```

---

## Phase 6: Frontend Component Tests

### Task 43: Test auth components (login/register forms)

**Files:**
- Create: `packages/observer-web/src/routes/_auth/login.test.tsx`
- Create: `packages/observer-web/src/routes/_auth/register.test.tsx`

Tests:
- Renders form fields
- Shows validation errors on empty submit
- Shows error message on API failure (401, 500)
- Calls mutate on valid submit

Use `@testing-library/react` render + `@testing-library/user-event` for interactions.

### Task 44: Test data-table component

**Files:**
- Create: `packages/observer-web/src/components/data-table.test.tsx`

Tests:
- Renders with data
- Renders empty state when no data
- Renders loading state
- Pagination controls work

### Task 45: Test form components

**Files:**
- Create: `packages/observer-web/src/components/confirm-dialog.test.tsx`
- Create: `packages/observer-web/src/components/form-field.test.tsx`
- Create: `packages/observer-web/src/components/drawer-shell.test.tsx`

Tests:
- ConfirmDialog: renders, calls onConfirm/onCancel
- FormField: renders label, shows error
- DrawerShell: renders children, close button works

### Task 46: Test chart components

**Files:**
- Create: `packages/observer-web/src/components/charts/bar-chart.test.tsx`
- Create: `packages/observer-web/src/components/charts/pie-chart.test.tsx`

Tests:
- Renders without crashing with valid data
- Renders without crashing with empty data
- Renders without crashing with single data point

### Task 47: Test utility components

**Files:**
- Create: `packages/observer-web/src/components/empty-state.test.tsx`
- Create: `packages/observer-web/src/components/status-badge.test.tsx`
- Create: `packages/observer-web/src/components/pagination.test.tsx`
- Create: `packages/observer-web/src/components/alert-banner.test.tsx`

Tests:
- EmptyState: renders message
- StatusBadge: renders correct variant for each status
- Pagination: renders page numbers, calls onChange
- AlertBanner: renders message, close button works

**Run all component tests:**

Run: `cd packages/observer-web && bun test src/components/ src/routes/_auth/`
Expected: all PASS

**Commit:**
```bash
git add packages/observer-web/src/components/*.test.tsx packages/observer-web/src/components/**/*.test.tsx packages/observer-web/src/routes/_auth/*.test.tsx
git commit -m "add frontend component tests for auth forms, data table, charts, utility components"
```

---

## Phase 7: Frontend Route Smoke Tests

### Task 48: Create route smoke test infrastructure

**Files:**
- Create: `packages/observer-web/src/test/route-smoke.test.tsx`

Test that each route component can be imported and rendered in isolation (with mocked data/hooks). This catches import errors, missing providers, and render crashes.

For app routes that need auth context, mock the auth store. For routes that fetch data, mock the hooks.

Test cases:
- Login route renders
- Register route renders
- App index (dashboard) renders with mocked data
- Project list renders with mocked data
- People list renders with mocked data
- Person detail renders with mocked data
- Tags page renders
- Households page renders
- Documents page renders
- Pets page renders
- Reports page renders
- Admin pages render (users, projects, reference)

**Run:**

Run: `cd packages/observer-web && bun test src/test/route-smoke`
Expected: all PASS

**Commit:**
```bash
git add packages/observer-web/src/test/route-smoke.test.tsx
git commit -m "add route smoke tests for all frontend routes"
```

---

## Phase 8: CLI --help Improvements

### Task 49: Enhance cobra command help text

**Files:**
- Modify: `cmd/observer/cmd/serve.go`
- Modify: `cmd/observer/cmd/migrate.go`
- Modify: `cmd/observer/cmd/keygen.go`
- Modify: `cmd/observer/cmd/create_admin.go`
- Modify: `cmd/observer/cmd/seed.go`
- Modify: `cmd/observer/main.go`

For each command, add/improve:
- `Long` description with clear explanation
- `Example` field with realistic usage examples

Example for serve:
```go
Long: `Start the Observer HTTP server.

Reads configuration from environment variables (DATABASE_DSN, REDIS_URL, etc.).
In production builds, embedded migrations are applied automatically on startup.
Graceful shutdown on SIGINT/SIGTERM with a 30-second timeout.`,
Example: `  # Start with defaults (localhost:9000)
  observer serve

  # Custom host and port
  observer serve --host 0.0.0.0 --port 8080

  # With environment configuration
  DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve`,
```

Apply similar improvements to all commands.

**Run: verify help output**

Run: `go run ./cmd/observer serve --help`
Expected: shows Long description and examples

**Commit:**
```bash
git add cmd/observer/cmd/*.go cmd/observer/main.go
git commit -m "improve CLI --help with descriptions and usage examples"
```

---

## Phase 9: CLI Documentation

### Task 50: Write docs/guides/cli.md

**Files:**
- Create: `docs/guides/cli.md`

**Content structure:**

```markdown
# Observer CLI Reference

## Overview
Observer provides a CLI for managing the server, database migrations, and development utilities.

## Installation
go install github.com/lbrty/observer/cmd/observer@latest
# or build from source
just build

## Commands

### serve
Start the HTTP server.
[flags table, env vars table, examples]

### migrate
Database migration management.
#### migrate up
[flags, examples]
#### migrate create
[flags, examples]
#### migrate version
[flags, examples]

### keygen
Generate RSA key pair for JWT signing.
[flags, examples]

### create-admin
Create a platform administrator account.
[flags, examples]

### seed
Seed database with development data (destructive).
[flags, examples]

## Common Workflows

### First-time setup
### Adding a new migration
### Seeding development data
### Generating new JWT keys

## Environment Variables
[complete table of all env vars from config.go]
```

**Commit:**
```bash
git add docs/guides/cli.md
git commit -m "add CLI reference documentation"
```

---

## Phase 10: `observer setup` Command

### Task 51: Create the setup command

**Files:**
- Create: `cmd/observer/cmd/setup.go`
- Modify: `cmd/observer/main.go` (register command)

**What it does:**
1. Generates `.env` with sensible defaults for all config vars
2. Auto-generates RSA keys into `keys/` directory (calls keygen logic)
3. Prompts for admin email + password (stdin), creates admin user
4. Prints next-steps summary with useful links

**Step 1: Write the setup command**

```go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "First-time project setup",
	Long: `Run first-time setup for Observer.

Generates a .env file with sensible defaults, creates RSA keys for JWT signing,
and prompts for an admin account. After setup, start Postgres + Redis and run
the server.`,
	Example: `  # Run interactive setup
  observer setup

  # Then start the server
  just serve`,
	RunE: runSetup,
}
```

**Step 2: Implement `runSetup`**

The function should:

1. **Check if `.env` already exists** — if so, ask to overwrite or skip
2. **Write `.env`** with these defaults:
```env
# Server
SERVER_HOST=localhost
SERVER_PORT=9000

# Database
DATABASE_DSN=postgres://observer:observer@localhost:5432/observer?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# JWT
JWT_PRIVATE_KEY_PATH=keys/jwt_rsa
JWT_PUBLIC_KEY_PATH=keys/jwt_rsa.pub
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
JWT_ISSUER=observer

# CORS
CORS_ORIGINS=http://localhost:5173

# Cookies
COOKIE_DOMAIN=
COOKIE_SECURE=false
COOKIE_SAME_SITE=lax

# Storage
STORAGE_PATH=data/uploads

# Logging
LOG_LEVEL=info
```

3. **Generate RSA keys** — reuse keygen logic (`generateRSAKeys(4096, "keys")`)
4. **Create required directories** — `keys/`, `data/uploads/`
5. **Prompt for admin** — read email from stdin, read password with `term.ReadPassword` (hidden input), confirm password matches, then create admin (reuse create-admin logic: connect DB, hash password, insert user)
6. **Print next-steps:**
```
Setup complete!

Next steps:
  1. Start Postgres and Redis:
     docker compose up -d

  2. Run database migrations:
     observer migrate up
     # or: just migrate

  3. Start the server:
     observer serve
     # or: just serve

Useful links:
  Releases: https://github.com/sultaniman/observer/releases
  Justfile commands: just --list
```

**Step 3: Write tests**

**Files:**
- Create: `cmd/observer/cmd/setup_test.go`

Test cases:
- `TestSetup_WritesEnvFile` — verify `.env` content written correctly (use temp dir)
- `TestSetup_SkipsExistingEnv` — when `.env` exists and user declines overwrite
- `TestSetup_GeneratesKeys` — verify key files created in temp dir
- `TestSetup_ValidatesAdminPassword` — rejects password < 8 chars
- `TestSetup_ValidatesAdminEmail` — rejects empty email

**Step 4: Run tests**

Run: `go test ./cmd/observer/cmd/ -run TestSetup -v`
Expected: all PASS

**Step 5: Register in main.go**

Add `rootCmd.AddCommand(setupCmd)` in the init or main function.

**Step 6: Verify end-to-end**

Run: `go run ./cmd/observer setup --help`
Expected: shows Long description and examples

**Step 7: Commit**

```bash
git add cmd/observer/cmd/setup.go cmd/observer/cmd/setup_test.go cmd/observer/main.go
git commit -m "add observer setup command for first-time project initialization"
```

### Task 52: Update CLI docs with setup command

**Files:**
- Modify: `docs/guides/cli.md` (add setup section)

Add `setup` command documentation:
- Description
- What it generates (.env, keys, admin user)
- Example output
- Move "First-time setup" workflow to reference `observer setup`

**Commit:**
```bash
git add docs/guides/cli.md
git commit -m "add setup command to CLI documentation"
```

---

## Execution Order Summary

| Phase | Tasks | Description |
|-------|-------|-------------|
| 1 | 1 | Handler test infrastructure |
| 2 | 2–16 | Handler tests (all 17 handlers) |
| 3 | 17–27 | Use case tests (fill gaps) |
| 4 | 28–33 | Repository integration tests |
| 5 | 34–42 | Frontend hook tests |
| 6 | 43–47 | Frontend component tests |
| 7 | 48 | Frontend route smoke tests |
| 8 | 49 | CLI --help improvements |
| 9 | 50 | CLI docs/guides/cli.md |
| 10 | 51–52 | `observer setup` command + docs |
