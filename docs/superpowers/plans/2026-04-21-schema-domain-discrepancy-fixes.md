# Schema-Domain Discrepancy Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Resolve 8 schema-domain discrepancies: introduce `PetListFilter`, move `SearchHits` to the domain layer, add `PersonCategoryRepository.ListBulk`, add `MigrationRecordRepository.Delete` with handler and route, fix `MFARecoveryCode.ID` type, fix `audit.Entry` nullability, annotate enrichment fields, and replace the giant `//go:generate` directive with source mode.

**Architecture:** All changes are surgical and independent — each touches a narrow vertical slice (domain entity → repository interface → implementation → use case → handler). Mocks are regenerated once at the end after all interface changes land. Tasks are ordered to minimize broken intermediate states: additive changes first, type propagation changes second, doc-only last, mock regeneration final.

**Tech Stack:** Go 1.26, Gin, sqlx, gomock (`go.uber.org/mock`), testify, `just` task runner.

---

## File Map

| File                                                   | Changes                                                                                                                                                                              |
| ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `internal/domain/pet/entity.go`                        | Add `PetListFilter` struct                                                                                                                                                           |
| `internal/domain/search/types.go`                      | **Create** — move `SearchHits`, `PersonHit`, `PetHit`, `ProjectHit` here                                                                                                             |
| `internal/repository/interfaces.go`                    | Update `PetRepository.List`, `SearchRepository.Search`, add `PersonCategoryRepository.ListBulk`, add `MigrationRecordRepository.Delete`, update `MFARecoveryCodeRepository.MarkUsed` |
| `internal/repository/search_types.go`                  | **Delete**                                                                                                                                                                           |
| `internal/repository/pet_repository.go`                | Update `petRepo.List` to accept `pet.PetListFilter`                                                                                                                                  |
| `internal/repository/person_repository.go`             | Add `ListBulk` to `personCategoryRepo`                                                                                                                                               |
| `internal/repository/migration_record_repository.go`   | Add `Delete` method                                                                                                                                                                  |
| `internal/repository/mfa_recovery_code_repository.go`  | Update `FindUnused` scan, `MarkUsed` param, `CreateBatch` ID serialisation                                                                                                           |
| `internal/repository/audit_repository.go`              | Update insert and scan for `*string` IP/UserAgent                                                                                                                                    |
| `internal/repository/search_repository.go`             | Update to use `domainsearch.SearchHits` etc.                                                                                                                                         |
| `internal/repository/generate.go`                      | **Create** — single source-mode `//go:generate`                                                                                                                                      |
| `internal/repository/mock/repository.go`               | Regenerated (do not edit manually)                                                                                                                                                   |
| `internal/domain/user/entity.go`                       | Change `MFARecoveryCode.ID` to `ulid.ULID`                                                                                                                                           |
| `internal/domain/audit/entity.go`                      | Change `IP`, `UserAgent` to `*string`; annotate enrichment fields                                                                                                                    |
| `internal/domain/household/entity.go`                  | Annotate enrichment fields                                                                                                                                                           |
| `internal/domain/support/entity.go`                    | Annotate enrichment fields                                                                                                                                                           |
| `internal/usecase/project/pet_usecase.go`              | Build `pet.PetListFilter` before calling repo                                                                                                                                        |
| `internal/usecase/project/migration_record_usecase.go` | Add `Delete` method                                                                                                                                                                  |
| `internal/usecase/auth/mfa_usecase.go`                 | Update `generateRecoveryCodes` to use `ulid.New()`                                                                                                                                   |
| `internal/usecase/audit/audit_usecase.go`              | Update `Record` and `Log` to pass `*string` IP/UserAgent                                                                                                                             |
| `internal/usecase/audit/types.go`                      | Change `EntryDTO.IP` and `EntryDTO.UserAgent` to `*string`                                                                                                                           |
| `internal/handler/migration_record_handler.go`         | Add `Delete` handler method                                                                                                                                                          |
| `internal/server/routes.go`                            | Register `DELETE /people/:person_id/migration-records/:id`                                                                                                                           |

---

## Task 1: Add `PetListFilter` and update `PetRepository.List`

**Files:**

- Modify: `internal/domain/pet/entity.go`
- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/pet_repository.go`
- Modify: `internal/usecase/project/pet_usecase.go`
- Modify: `internal/usecase/project/pet_usecase_test.go`

- [ ] **Step 1: Add `PetListFilter` to the pet domain**

In `internal/domain/pet/entity.go`, append after the `Pet` struct:

```go
// PetListFilter holds filter and pagination parameters for pet list queries.
type PetListFilter struct {
	ProjectID string
	Status    *PetStatus
	TagIDs    []string
	Page      int
	PerPage   int
}
```

- [ ] **Step 2: Update `PetRepository.List` in the interface**

In `internal/repository/interfaces.go`, replace:

```go
List(ctx context.Context, projectID string, status string, tagIDs []string, page, perPage int) ([]*pet.Pet, int, error)
```

with:

```go
List(ctx context.Context, filter pet.PetListFilter) ([]*pet.Pet, int, error)
```

- [ ] **Step 3: Update the repository implementation**

In `internal/repository/pet_repository.go`, replace the `List` signature and its body's parameter extraction:

```go
func (r *petRepo) List(ctx context.Context, filter pet.PetListFilter) ([]*pet.Pet, int, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 20
	}

	where := []string{"project_id = $1"}
	args := []any{filter.ProjectID}
	ix := 1

	if filter.Status != nil {
		ix++
		where = append(where, "status = $"+strconv.Itoa(ix))
		args = append(args, *filter.Status)
	}

	var tagJoin string
	if len(filter.TagIDs) > 0 {
		placeholders := make([]string, len(filter.TagIDs))
		for i, tagID := range filter.TagIDs {
			ix++
			placeholders[i] = "$" + strconv.Itoa(ix)
			args = append(args, tagID)
		}
		tagJoin = " JOIN pet_tags pt ON pt.pet_id = pets.id AND pt.tag_id IN (" + strings.Join(placeholders, ",") + ")"
		where = append(where, "1=1 GROUP BY pets.id HAVING COUNT(DISTINCT pt.tag_id) = "+strconv.Itoa(len(filter.TagIDs)))
	}
	// remainder of the method body is unchanged
```

The rest of the method body (building `whereClause`, count query, pagination, SELECT query, row scanning) is unchanged — only the parameter extraction lines at the top change.

- [ ] **Step 4: Update the use case call site**

In `internal/usecase/project/pet_usecase.go`, replace the `List` method body's repo call:

```go
func (uc *PetUseCase) List(ctx context.Context, projectID string, input ListPetsInput) (*ListPetsOutput, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 {
		perPage = 20
	}

	filter := pet.PetListFilter{
		ProjectID: projectID,
		TagIDs:    input.TagIDs,
		Page:      page,
		PerPage:   perPage,
	}
	if input.Status != "" {
		s := pet.PetStatus(input.Status)
		filter.Status = &s
	}

	pets, total, err := uc.repo.List(ctx, filter)
```

Add import `"github.com/lbrty/observer/internal/domain/pet"` if not already present.

- [ ] **Step 5: Regenerate mocks**

```bash
just generate-mocks
```

Expected: `internal/repository/mock/repository.go` regenerated with the new `List` signature. No other output.

- [ ] **Step 6: Update the existing pet use case test**

In `internal/usecase/project/pet_usecase_test.go`, the mock expectation on line 28 currently uses the old raw-parameter signature. Replace:

```go
mockRepo.EXPECT().List(gomock.Any(), "proj1", "", []string(nil), gomock.Any(), gomock.Any()).Return(...)
```

with:

```go
mockRepo.EXPECT().List(gomock.Any(), pet.PetListFilter{
    ProjectID: "proj1",
    Page:      1,
    PerPage:   20,
}).Return(...)
```

Also update the `TestPetUseCase_List_RepoError` expectation on line 54:

```go
mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, repoErr)
```

- [ ] **Step 7: Run tests**

```bash
just test
```

Expected: all pet use case tests pass.

- [ ] **Step 8: Commit**

```bash
git add internal/domain/pet/entity.go \
        internal/repository/interfaces.go \
        internal/repository/pet_repository.go \
        internal/repository/mock/repository.go \
        internal/usecase/project/pet_usecase.go \
        internal/usecase/project/pet_usecase_test.go
git commit -m "Introduce PetListFilter and update PetRepository.List"
```

---

## Task 2: Move `SearchHits` to the domain layer

**Files:**

- Create: `internal/domain/search/types.go`
- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/search_repository.go`
- Delete: `internal/repository/search_types.go`

- [ ] **Step 1: Create the domain search package**

Create `internal/domain/search/types.go`:

```go
package search

// PersonHit is a search result for a person entity.
type PersonHit struct {
	ID        string
	FirstName string
	LastName  string
	ProjectID string
}

// PetHit is a search result for a pet entity.
type PetHit struct {
	ID        string
	Name      string
	ProjectID string
}

// ProjectHit is a search result for a project entity.
type ProjectHit struct {
	ID   string
	Name string
}

// SearchHits holds search results across all entity types.
type SearchHits struct {
	People   []PersonHit
	Pets     []PetHit
	Projects []ProjectHit
}
```

- [ ] **Step 2: Update the repository interface**

In `internal/repository/interfaces.go`, add the import:

```go
"github.com/lbrty/observer/internal/domain/search"
```

Then update `SearchRepository.Search` return type:

```go
Search(ctx context.Context, projectIDs []string, query string, limit int) (*search.SearchHits, error)
```

- [ ] **Step 3: Update the search repository implementation**

In `internal/repository/search_repository.go`, add the import:

```go
domainsearch "github.com/lbrty/observer/internal/domain/search"
```

Replace all unqualified type references:

- `&SearchHits{}` → `&domainsearch.SearchHits{}`
- `SearchHits` (as a variable type) → `domainsearch.SearchHits`
- `PersonHit` → `domainsearch.PersonHit`
- `PetHit` → `domainsearch.PetHit`
- `ProjectHit` → `domainsearch.ProjectHit`

The four affected functions are `Search`, `searchPeople`, `searchPets`, `searchProjects`.

- [ ] **Step 4: Delete `internal/repository/search_types.go`**

```bash
git rm internal/repository/search_types.go
```

- [ ] **Step 5: Regenerate mocks**

```bash
just generate-mocks
```

Expected: mock regenerated cleanly. The `SearchRepository` mock now uses `*search.SearchHits`.

- [ ] **Step 6: Build and test**

```bash
just test
```

Expected: all tests pass. The search use case (`internal/usecase/search/search_usecase.go`) uses `hits.People`, `hits.Pets`, `hits.Projects` without naming the type directly — no changes needed there.

- [ ] **Step 7: Commit**

```bash
git add internal/domain/search/types.go \
        internal/repository/interfaces.go \
        internal/repository/search_repository.go \
        internal/repository/mock/repository.go
git rm internal/repository/search_types.go
git commit -m "Move SearchHits to internal/domain/search"
```

---

## Task 3: Add `PersonCategoryRepository.ListBulk`

**Files:**

- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/person_repository.go`

- [ ] **Step 1: Add the method to the interface**

In `internal/repository/interfaces.go`, inside `PersonCategoryRepository`, add after the existing `List` method:

```go
// ListBulk fetches category IDs for multiple people in one query.
// Returns a map of person ID → []category ID.
ListBulk(ctx context.Context, personIDs []string) (map[string][]string, error)
```

- [ ] **Step 2: Implement in the repository**

In `internal/repository/person_repository.go`, add after `personCategoryRepo.ReplaceAll`:

```go
func (r *personCategoryRepo) ListBulk(ctx context.Context, personIDs []string) (map[string][]string, error) {
	if len(personIDs) == 0 {
		return map[string][]string{}, nil
	}
	params := make([]string, len(personIDs))
	args := make([]any, len(personIDs))
	for i, id := range personIDs {
		args[i] = id
		params[i] = fmt.Sprintf("$%d", i+1)
	}
	q := fmt.Sprintf(
		"SELECT person_id, category_id FROM person_categories WHERE person_id IN (%s) ORDER BY person_id, category_id",
		joinStrings(params, ", "),
	)
	return queryBulkTags(ctx, r.db, q, args)
}
```

`joinStrings` and `queryBulkTags` are already defined in `internal/repository/helpers.go` and are package-internal.

- [ ] **Step 3: Regenerate mocks**

```bash
just generate-mocks
```

- [ ] **Step 4: Write a unit test**

Add to `internal/usecase/project/person_usecase_test.go` (or whichever test file covers the person use case). If the person list use case already calls `PersonCategoryRepository.List` in a loop, verify the mock expectation compiles with the new interface. If there is no existing test covering bulk category loading, add:

```go
func TestPersonCategoryRepo_ListBulk_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	mockCatRepo.EXPECT().ListBulk(gomock.Any(), []string{}).Return(map[string][]string{}, nil)

	result, err := mockCatRepo.ListBulk(context.Background(), []string{})
	require.NoError(t, err)
	assert.Empty(t, result)
}
```

- [ ] **Step 5: Check for N+1 call sites**

Search for any use case that calls `PersonCategoryRepository.List` inside a loop over people:

```bash
grep -rn "\.List(ctx\|\.List(c\.Request" internal/usecase/ --include="*.go" | grep -i "categor"
```

If found, update that use case to call `ListBulk` once outside the loop (same pattern as how `PetTagUseCase` and `PersonTagUseCase` use `ListBulk` to hydrate tag IDs on list views).

- [ ] **Step 6: Run tests**

```bash
just test
```

Expected: all tests pass.

- [ ] **Step 7: Commit**

```bash
git add internal/repository/interfaces.go \
        internal/repository/person_repository.go \
        internal/repository/mock/repository.go
git commit -m "Add PersonCategoryRepository.ListBulk"
```

---

## Task 4: Add `MigrationRecordRepository.Delete` with use case, handler, and route

**Files:**

- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/migration_record_repository.go`
- Modify: `internal/usecase/project/migration_record_usecase.go`
- Modify: `internal/usecase/project/migration_record_usecase_test.go`
- Modify: `internal/handler/migration_record_handler.go`
- Modify: `internal/server/routes.go`

- [ ] **Step 1: Add `Delete` to the interface**

In `internal/repository/interfaces.go`, inside `MigrationRecordRepository`, add after `Update`:

```go
Delete(ctx context.Context, id string) error
```

- [ ] **Step 2: Implement in the repository**

In `internal/repository/migration_record_repository.go`, add after `Update`:

```go
func (r *migrationRecordRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM migration_records WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete migration record: %w", err)
	}
	return CheckRowsAffected(res, migration.ErrRecordNotFound)
}
```

- [ ] **Step 3: Regenerate mocks**

```bash
just generate-mocks
```

- [ ] **Step 4: Write a failing use case test**

In `internal/usecase/project/migration_record_usecase_test.go`, add:

```go
func TestMigrationRecordUseCase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, auditUC)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(&migration.Record{
		ID:       "mr1",
		PersonID: "p1",
	}, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), "mr1").Return(nil)

	err := uc.Delete(context.Background(), "proj1", "p1", "mr1")
	require.NoError(t, err)
}

func TestMigrationRecordUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr99").Return(nil, migration.ErrRecordNotFound)

	err := uc.Delete(context.Background(), "proj1", "p1", "mr99")
	assert.ErrorIs(t, err, migration.ErrRecordNotFound)
}
```

- [ ] **Step 5: Run the test to confirm it fails**

```bash
just test
```

Expected: compile error — `Delete` method does not exist on `MigrationRecordUseCase`.

- [ ] **Step 6: Add `Delete` to the use case**

In `internal/usecase/project/migration_record_usecase.go`, add after `Update`:

```go
// Delete removes a migration record after verifying it belongs to the given person.
func (uc *MigrationRecordUseCase) Delete(ctx context.Context, projectID, personID, id string) error {
	rec, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rec.PersonID != personID {
		return migration.ErrRecordNotFound
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	uc.auditUC.Record(ctx, &projectID, "migration_record.delete", "migration_record", &id,
		fmt.Sprintf("Deleted migration record %s", id))
	return nil
}
```

- [ ] **Step 7: Run tests to confirm they pass**

```bash
just test
```

Expected: all migration record use case tests pass.

- [ ] **Step 8: Add the handler method**

In `internal/handler/migration_record_handler.go`, add after `Update`:

```go
// Delete handles DELETE /projects/:project_id/people/:person_id/migration-records/:id.
func (h *MigrationRecordHandler) Delete(c *gin.Context) {
	err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
```

- [ ] **Step 9: Register the route**

In `internal/server/routes.go`, inside the `del` group (which already contains `del.DELETE("/pets/:id", petHandler.Delete)`), add:

```go
del.DELETE("/people/:person_id/migration-records/:id", migrationHandler.Delete)
```

- [ ] **Step 10: Build**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 11: Commit**

```bash
git add internal/repository/interfaces.go \
        internal/repository/migration_record_repository.go \
        internal/repository/mock/repository.go \
        internal/usecase/project/migration_record_usecase.go \
        internal/usecase/project/migration_record_usecase_test.go \
        internal/handler/migration_record_handler.go \
        internal/server/routes.go
git commit -m "Add Delete to MigrationRecordRepository, use case, handler, and route"
```

---

## Task 5: Change `MFARecoveryCode.ID` from `string` to `ulid.ULID`

**Files:**

- Modify: `internal/domain/user/entity.go`
- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/mfa_recovery_code_repository.go`
- Modify: `internal/usecase/auth/mfa_usecase.go`

- [ ] **Step 1: Change the entity field**

In `internal/domain/user/entity.go`, change `MFARecoveryCode`:

```go
// MFARecoveryCode is a single-use backup code stored hashed.
type MFARecoveryCode struct {
	ID        ulid.ULID
	UserID    ulid.ULID
	CodeHash  string
	UsedAt    *time.Time
	CreatedAt time.Time
}
```

`ulid` is already imported (`github.com/oklog/ulid/v2`).

- [ ] **Step 2: Update `MarkUsed` in the interface**

In `internal/repository/interfaces.go`, inside `MFARecoveryCodeRepository`, change:

```go
MarkUsed(ctx context.Context, id string) error
```

to:

```go
MarkUsed(ctx context.Context, id ulid.ULID) error
```

- [ ] **Step 3: Update the repository implementation**

In `internal/repository/mfa_recovery_code_repository.go`:

a) In `CreateBatch`, change the insert to serialize `ID` to string:

```go
if _, err := r.db.ExecContext(ctx, q, c.ID.String(), c.UserID.String(), c.CodeHash, c.CreatedAt); err != nil {
```

b) In `FindUnused`, parse the scanned `ID` string into `ulid.ULID`:

```go
parsedID, err := ulid.Parse(row.ID)
if err != nil {
    return nil, err
}
return &domainuser.MFARecoveryCode{
    ID:        parsedID,
    UserID:    uid,
    CodeHash:  row.CodeHash,
    UsedAt:    row.UsedAt,
    CreatedAt: row.CreatedAt,
}, nil
```

c) In `MarkUsed`, change the parameter type and use `.String()` in the query:

```go
func (r *mfaRecoveryCodeRepo) MarkUsed(ctx context.Context, id ulid.ULID) error {
	const q = `UPDATE mfa_recovery_codes SET used_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id.String())
	return err
}
```

- [ ] **Step 4: Update `generateRecoveryCodes` in the MFA use case**

In `internal/usecase/auth/mfa_usecase.go`, `generateRecoveryCodes` currently sets:

```go
ID: iulid.NewString(),
```

Change to:

```go
ID: iulid.New(),
```

(`iulid` is `github.com/lbrty/observer/internal/ulid`, which exposes both `New() ulid.ULID` and `NewString() string`.)

Also on line 57, `uc.recoveryRepo.MarkUsed(ctx, rc.ID)` — `rc.ID` is now `ulid.ULID`, which matches the updated interface. No change needed.

- [ ] **Step 5: Regenerate mocks**

```bash
just generate-mocks
```

- [ ] **Step 6: Run tests**

```bash
just test
```

Expected: all tests pass. The existing MFA use case tests that set up `MockMFARecoveryCodeRepository.MarkUsed` expectations will need their `id` argument updated to `ulid.ULID`. Check `internal/usecase/auth/mfa_usecase_test.go` if it exists and update accordingly.

- [ ] **Step 7: Commit**

```bash
git add internal/domain/user/entity.go \
        internal/repository/interfaces.go \
        internal/repository/mfa_recovery_code_repository.go \
        internal/repository/mock/repository.go \
        internal/usecase/auth/mfa_usecase.go
git commit -m "Change MFARecoveryCode.ID to ulid.ULID"
```

---

## Task 6: Fix `audit.Entry` IP/UserAgent nullability

**Files:**

- Modify: `internal/domain/audit/entity.go`
- Modify: `internal/repository/audit_repository.go`
- Modify: `internal/usecase/audit/audit_usecase.go`
- Modify: `internal/usecase/audit/types.go`

- [ ] **Step 1: Change the entity fields**

In `internal/domain/audit/entity.go`, change:

```go
IP            string    `db:"ip"`
UserAgent     string    `db:"user_agent"`
```

to:

```go
IP            *string   `db:"ip"`
UserAgent     *string   `db:"user_agent"`
```

- [ ] **Step 2: Update the audit repository**

In `internal/repository/audit_repository.go`, the `Log` method passes `entry.IP` and `entry.UserAgent` directly to `ExecContext`. Since they are now `*string`, the `database/sql` driver will write NULL for nil — no code change required in the `Log` method body itself.

The `List` method uses `r.db.SelectContext(ctx, &entries, ...)` with sqlx struct scanning. Because `IP` and `UserAgent` are now `*string` with `db` tags, sqlx handles NULL → nil correctly — no change required in `List` either.

Verify by reading the method once — if it uses `sqlx.SelectContext` into `[]audit.Entry`, the struct tags do all the work.

- [ ] **Step 3: Add a `strPtrOrNil` helper to the audit use case**

In `internal/usecase/audit/audit_usecase.go`, add at package level (after the imports):

```go
func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
```

- [ ] **Step 4: Update `Record` in the use case**

In `internal/usecase/audit/audit_usecase.go`, the `Record` method currently sets:

```go
IP:         middleware.AuditIP(ctx),
UserAgent:  middleware.AuditUserAgent(ctx),
```

Change to:

```go
IP:         strPtrOrNil(middleware.AuditIP(ctx)),
UserAgent:  strPtrOrNil(middleware.AuditUserAgent(ctx)),
```

- [ ] **Step 5: Update `Log` in the use case**

In the same file, the `Log` method sets:

```go
IP:         input.IP,
UserAgent:  input.UserAgent,
```

Change to:

```go
IP:         strPtrOrNil(input.IP),
UserAgent:  strPtrOrNil(input.UserAgent),
```

(`LogInput.IP` and `LogInput.UserAgent` remain `string` — the conversion happens here.)

- [ ] **Step 6: Update the `List` method**

In the same file, the `List` method builds `EntryDTO` from `audit.Entry`. Currently:

```go
IP:            e.IP,
UserAgent:     e.UserAgent,
```

These become direct assignments since both the source (`audit.Entry.IP`) and destination (`EntryDTO.IP`) will be `*string` after the next step. No change in this assignment line itself.

- [ ] **Step 7: Update `EntryDTO`**

In `internal/usecase/audit/types.go`, change:

```go
IP            string  `json:"ip"`
UserAgent     string  `json:"user_agent"`
```

to:

```go
IP            *string `json:"ip,omitempty"`
UserAgent     *string `json:"user_agent,omitempty"`
```

- [ ] **Step 8: Run tests**

```bash
just test
```

Expected: all tests pass. If any existing test constructs `audit.Entry{IP: "...", UserAgent: "..."}` directly, update those to `IP: strPtr("...")` using a `func strPtr(s string) *string { return &s }` helper in the test file.

- [ ] **Step 9: Commit**

```bash
git add internal/domain/audit/entity.go \
        internal/repository/audit_repository.go \
        internal/usecase/audit/audit_usecase.go \
        internal/usecase/audit/types.go
git commit -m "Fix audit.Entry IP and UserAgent to *string matching nullable DB columns"
```

---

## Task 7: Annotate enrichment fields as read-only projections

**Files:**

- Modify: `internal/domain/household/entity.go`
- Modify: `internal/domain/support/entity.go`
- Modify: `internal/domain/audit/entity.go`

- [ ] **Step 1: Annotate `household.Household`**

In `internal/domain/household/entity.go`, add a comment block above the enrichment fields:

```go
// Household represents a family unit within a project.
type Household struct {
	ID              string
	ProjectID       string
	ReferenceNumber *string
	HeadPersonID    *string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Populated by repository reads via JOIN; zero-valued on Create/Update.
	HeadPersonName *string
	MemberCount    int
}
```

- [ ] **Step 2: Annotate `support.Record`**

In `internal/domain/support/entity.go`, add a comment block above the enrichment fields:

```go
	// Populated by repository reads via JOIN with people; zero-valued on Create/Update.
	PersonFirstName *string
	PersonLastName  *string
```

- [ ] **Step 3: Annotate `audit.Entry`**

In `internal/domain/audit/entity.go`, add a comment block above the enrichment fields:

```go
	// Populated by repository reads via LEFT JOIN with users; zero-valued on Log.
	UserFirstName string `db:"user_first_name"`
	UserLastName  string `db:"user_last_name"`
	UserEmail     string `db:"user_email"`
```

- [ ] **Step 4: Build**

```bash
go build ./...
```

Expected: clean build. No logic changed.

- [ ] **Step 5: Commit**

```bash
git add internal/domain/household/entity.go \
        internal/domain/support/entity.go \
        internal/domain/audit/entity.go
git commit -m "Annotate enrichment fields as read-only projections on entity structs"
```

---

## Task 8: Replace giant `//go:generate` with source mode

**Files:**

- Create: `internal/repository/generate.go`
- Modify: `internal/repository/interfaces.go`

- [ ] **Step 1: Create `generate.go`**

Create `internal/repository/generate.go`:

```go
package repository

//go:generate mockgen -source=interfaces.go -destination=mock/repository.go -package=mock
```

- [ ] **Step 2: Remove the directive from `interfaces.go`**

In `internal/repository/interfaces.go`, delete line 24 (the `//go:generate mockgen ...` line).

- [ ] **Step 3: Regenerate mocks and verify**

```bash
just generate-mocks
```

Expected: `internal/repository/mock/repository.go` regenerated. The file header now reads:

```
// Source: interfaces.go
```

and includes `MockSearchRepository` which was previously missing from the reflect-mode list.

- [ ] **Step 4: Run full test suite**

```bash
just test
```

Expected: all tests pass.

- [ ] **Step 5: Commit**

```bash
git add internal/repository/generate.go \
        internal/repository/interfaces.go \
        internal/repository/mock/repository.go
git commit -m "Replace reflect-mode go:generate list with source-mode in generate.go"
```

---

## Verification

After all 8 tasks:

```bash
just test
go build ./...
```

Both must succeed with zero errors.
