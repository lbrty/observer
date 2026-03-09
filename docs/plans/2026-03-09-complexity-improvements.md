# Complexity Improvements Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reduce cyclomatic complexity in three identified hotspots without changing behavior.

**Architecture:** Three independent refactors — a generic SQL clause helper in the repository layer, a pointer-copy helper in the usecase layer, and URL search params replacing local filter state in report pages.

**Tech Stack:** Go generics (1.21+), TanStack Router `validateSearch`

---

### Task 1: Generic SQL Filter Helper — `appendIf`

**Files:**
- Modify: `internal/repository/helpers.go`
- Modify: `internal/repository/report_repository.go`

**Step 1: Add helper to `helpers.go`**

Append at the bottom of `internal/repository/helpers.go`:

```go
// appendIf appends a SQL clause with a positional arg if v is non-nil.
// clause must contain a single %d verb for the parameter index, e.g. " AND p.sex = $%d".
func appendIf[T any](q string, args []any, ix int, clause string, v *T) (string, []any, int) {
	if v == nil {
		return q, args, ix
	}
	return q + fmt.Sprintf(clause, ix), append(args, *v), ix + 1
}
```

**Step 2: Run tests to confirm nothing is broken yet**

```bash
just test
```

Expected: all pass.

**Step 3: Refactor `applyPeopleFilters` in `report_repository.go`**

Replace the entire `applyPeopleFilters` function body with:

```go
func applyPeopleFilters(q string, f report.ReportFilter, dateCol string, args []any, ix int) (string, []any, int) {
	q, args, ix = appendIf(q, args, ix, " AND "+dateCol+" >= $%d", f.DateFrom)
	q, args, ix = appendIf(q, args, ix, " AND "+dateCol+" <= $%d", f.DateTo)
	q, args, ix = appendIf(q, args, ix, " AND p.case_status = $%d", f.CaseStatus)
	q, args, ix = appendIf(q, args, ix, " AND p.sex = $%d", f.Sex)
	q, args, ix = appendIf(q, args, ix, " AND p.id IN (SELECT person_id FROM person_categories WHERE category_id = $%d)", f.CategoryID)
	q, args, ix = appendIf(q, args, ix, " AND p.office_id = $%d", f.OfficeID)
	q, args, ix = appendIf(q, args, ix, " AND p.consultant_id = $%d", f.ConsultantID)
	q, args, ix = appendIf(q, args, ix, " AND p.id IN (SELECT sr2.person_id FROM support_records sr2 WHERE sr2.project_id = $1 AND sr2.type = $%d)", f.SupportType)
	return q, args, ix
}
```

**Step 4: Refactor `applySupportFilters` — simple clauses only**

Replace the first block of `applySupportFilters` (before the `personClauses` section) with `appendIf` calls:

```go
func applySupportFilters(q string, f report.ReportFilter, dateCol string, args []any, ix int) (string, []any, int) {
	q, args, ix = appendIf(q, args, ix, " AND "+dateCol+" >= $%d", f.DateFrom)
	q, args, ix = appendIf(q, args, ix, " AND "+dateCol+" <= $%d", f.DateTo)
	q, args, ix = appendIf(q, args, ix, " AND sr.office_id = $%d", f.OfficeID)
	q, args, ix = appendIf(q, args, ix, " AND sr.consultant_id = $%d", f.ConsultantID)
	q, args, ix = appendIf(q, args, ix, " AND sr.type = $%d", f.SupportType)

	// Person-level filters via subquery — multi-clause, kept explicit.
	var personClauses []string
	var personArgs []any
	if f.CategoryID != nil {
		personClauses = append(personClauses, fmt.Sprintf("id IN (SELECT person_id FROM person_categories WHERE category_id = $%d)", ix))
		personArgs = append(personArgs, *f.CategoryID)
		ix++
	}
	if f.Sex != nil {
		personClauses = append(personClauses, fmt.Sprintf("sex = $%d", ix))
		personArgs = append(personArgs, *f.Sex)
		ix++
	}
	if f.CaseStatus != nil {
		personClauses = append(personClauses, fmt.Sprintf("case_status = $%d", ix))
		personArgs = append(personArgs, *f.CaseStatus)
		ix++
	}
	if len(personClauses) > 0 {
		sub := " AND sr.person_id IN (SELECT id FROM people WHERE "
		for i, clause := range personClauses {
			if i > 0 {
				sub += " AND "
			}
			sub += clause
		}
		q += sub + ")"
		args = append(args, personArgs...)
	}

	return q, args, ix
}
```

**Step 5: Run tests**

```bash
just test
```

Expected: all pass.

**Step 6: Commit**

```bash
git add internal/repository/helpers.go internal/repository/report_repository.go
git commit -m "Add generic appendIf helper, reduce filter builder branches"
```

---

### Task 2: Pointer-Copy Helper — `setPtr` and `applyOpt`

**Files:**
- Create: `internal/usecase/project/patch.go`
- Modify: `internal/usecase/project/person_usecase.go`

**Step 1: Create `patch.go`**

```go
package project

// setPtr sets *dst to src when src is non-nil. Use for pointer-typed entity fields.
func setPtr[T any](dst **T, src *T) {
	if src != nil {
		*dst = src
	}
}

// applyOpt dereferences src into dst when src is non-nil. Use for value-typed entity fields.
func applyOpt[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}
```

**Step 2: Run tests**

```bash
just test
```

Expected: all pass (no behavior change yet).

**Step 3: Refactor `Update` in `person_usecase.go`**

Replace the nil-check block at the top of `Update` (pointer and simple value fields) with:

```go
setPtr(&p.ConsultantID, input.ConsultantID)
setPtr(&p.OfficeID, input.OfficeID)
setPtr(&p.CurrentPlaceID, input.CurrentPlaceID)
setPtr(&p.OriginPlaceID, input.OriginPlaceID)
setPtr(&p.ExternalID, input.ExternalID)
setPtr(&p.LastName, input.LastName)
setPtr(&p.Patronymic, input.Patronymic)
setPtr(&p.Email, input.Email)
setPtr(&p.PrimaryPhone, input.PrimaryPhone)
applyOpt(&p.FirstName, input.FirstName)
applyOpt(&p.ConsentGiven, input.ConsentGiven)
```

Keep the enum conversions and complex fields explicit (they cannot be generalized without losing type safety):

```go
if input.Sex != nil {
    p.Sex = person.Sex(*input.Sex)
}
if input.CaseStatus != nil {
    p.CaseStatus = person.CaseStatus(*input.CaseStatus)
}
if input.AgeGroup != nil {
    ag := person.AgeGroup(*input.AgeGroup)
    p.AgeGroup = &ag
}
if input.PhoneNumbers != nil {
    b, _ := json.Marshal(input.PhoneNumbers)
    p.PhoneNumbers = b
}
```

Apply the same `setPtr`/`applyOpt` pattern in `Create` where identical assignments appear.

**Step 4: Run tests**

```bash
just test
```

Expected: all pass.

**Step 5: Commit**

```bash
git add internal/usecase/project/patch.go internal/usecase/project/person_usecase.go
git commit -m "Extract setPtr/applyOpt helpers, reduce person update nil-check branches"
```

---

### Task 3: URL Search Params for Report Filter State

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/reports/people.tsx`
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/reports/pets.tsx`

**Step 1: Add `validateSearch` to the people report route**

Replace the `createFileRoute` call:

```ts
export const Route = createFileRoute("/_app/projects/$projectId/reports/people")({
  validateSearch: (search: Record<string, unknown>): ReportParams => ({
    date_from: search.date_from as string | undefined,
    date_to: search.date_to as string | undefined,
    office_id: search.office_id as string | undefined,
    category_id: search.category_id as string | undefined,
    consultant_id: search.consultant_id as string | undefined,
    case_status: search.case_status as string | undefined,
    sex: search.sex as string | undefined,
    age_group: search.age_group as string | undefined,
    support_type: search.support_type as string | undefined,
  }),
  component: ReportsPage,
});
```

**Step 2: Replace local filter state with URL search in `ReportsPage`**

Remove:
```ts
const [params, setParams] = useState<ReportParams>({});
```

Add:
```ts
const params = Route.useSearch();
const navigate = Route.useNavigate();

function setParams(update: ReportParams | ((prev: ReportParams) => ReportParams)) {
  navigate({
    search: typeof update === "function" ? update : () => update,
    replace: true,
  });
}
```

Keep `filtersOpen` and `activePreset` as local `useState` — they are UI-only and should not be in the URL.

**Step 3: Build and verify**

```bash
cd packages/observer-web && bun run build
```

Expected: no type errors. Navigate to the report page, apply filters, copy the URL — filters should survive a page reload.

**Step 4: Repeat for `pets.tsx`**

Apply the same `validateSearch` + `useSearch`/`useNavigate` pattern to `pets.tsx` using `PetReportParams` instead of `ReportParams`.

**Step 5: Commit**

```bash
git add packages/observer-web/src/routes/_app/projects/\$projectId/reports/
git commit -m "Move report filter state to URL search params"
```

---

## Execution Options

**1. Subagent-Driven (this session)** — dispatch a fresh subagent per task, review between tasks

**2. Parallel Session (separate)** — open new session in worktree, use `superpowers:executing-plans`

Which approach?
