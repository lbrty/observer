# Filter Search Button Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add an optional "Search" button to `FilterBar` that complements instant onChange filtering with a manual trigger.

**Architecture:** Add `onSearch?: () => void` to `FilterBar`; render a secondary `Button` at the end of the flex row when the prop is present. Thread the prop through `DataTablePage`. Each table/report page destructures `refetch` from its query hook and passes it as `onSearch`.

**Tech Stack:** React, TypeScript, react-i18next, `@tanstack/react-query` (`refetch`), `Button` component (`variant="secondary" size="md"`)

---

### Task 1: Add `common.search` i18n key to all locales

**Files:**
- Modify: `packages/observer-web/src/locales/en.json`
- Modify: `packages/observer-web/src/locales/ru.json`
- Modify: `packages/observer-web/src/locales/uk.json`
- Modify: `packages/observer-web/src/locales/de.json`
- Modify: `packages/observer-web/src/locales/tr.json`
- Modify: `packages/observer-web/src/locales/ky.json`

**Step 1: Add key to each locale's `common` object**

In each file, add `"search": "<translation>"` inside the `"common": { ... }` object (after `"disabled"`):

- `en.json`: `"search": "Search"`
- `ru.json`: `"search": "Поиск"`
- `uk.json`: `"search": "Пошук"`
- `de.json`: `"search": "Suchen"`
- `tr.json`: `"search": "Ara"`
- `ky.json`: `"search": "Izdöö"` — use the `kyrgyz-latin-transliteration` skill if unsure

**Step 2: Verify JSON is valid**

```bash
cd packages/observer-web
for f in src/locales/*.json; do python3 -c "import json,sys; json.load(open('$f'))" && echo "$f OK"; done
```

Expected: all files print `OK`

**Step 3: Commit**

```bash
git add packages/observer-web/src/locales/
git commit -m "Add common.search i18n key to all locales"
```

---

### Task 2: Add `onSearch` prop to `FilterBar`

**Files:**
- Modify: `packages/observer-web/src/components/filter-bar.tsx`

**Step 1: Update `FilterBarProps` and add the button**

Add `onSearch?: () => void` to the interface and render a `Button` after `{trailing}`:

```tsx
import { useTranslation } from "react-i18next";
import { Button } from "@/components/button";

// ...

interface FilterBarProps {
  filters: FilterDef[];
  trailing?: ReactNode;
  onSearch?: () => void;
}

export function FilterBar({ filters, trailing, onSearch }: FilterBarProps) {
  const { t } = useTranslation();
  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      {filters.map((f, i) => {
        // ... existing map unchanged ...
      })}
      {trailing}
      {onSearch && (
        <Button variant="secondary" onClick={onSearch}>
          {t("common.search")}
        </Button>
      )}
    </div>
  );
}
```

**Step 2: Verify TypeScript compiles**

```bash
cd packages/observer-web && bunx tsc --noEmit
```

Expected: no errors

**Step 3: Commit**

```bash
git add packages/observer-web/src/components/filter-bar.tsx
git commit -m "Add onSearch button to FilterBar"
```

---

### Task 3: Thread `onSearch` through `DataTablePage`

**Files:**
- Modify: `packages/observer-web/src/components/data-table-page.tsx`

**Step 1: Add prop and pass it to `FilterBar`**

```tsx
interface DataTablePageProps<T> {
  // ... existing props ...
  onSearch?: () => void;
}

export function DataTablePage<T>({
  // ... existing destructuring ...
  onSearch,
}: DataTablePageProps<T>) {
  return (
    <div>
      <PageHeader title={title} action={createAction} />

      {filters && filters.length > 0 && (
        <FilterBar filters={filters} trailing={filterTrailing} onSearch={onSearch} />
      )}

      {/* ... rest unchanged ... */}
    </div>
  );
}
```

**Step 2: Verify TypeScript compiles**

```bash
cd packages/observer-web && bunx tsc --noEmit
```

Expected: no errors

**Step 3: Commit**

```bash
git add packages/observer-web/src/components/data-table-page.tsx
git commit -m "Thread onSearch through DataTablePage to FilterBar"
```

---

### Task 4: Wire `onSearch` in people table page

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/people/index.tsx`

**Step 1: Destructure `refetch` from `usePeople`**

Find (line ~78):
```tsx
const { data, isLoading } = usePeople(projectId, params);
```

Replace with:
```tsx
const { data, isLoading, refetch } = usePeople(projectId, params);
```

**Step 2: Pass `refetch` as `onSearch` to `DataTablePage`**

Find the `<DataTablePage` JSX and add:
```tsx
onSearch={refetch}
```

**Step 3: Verify TypeScript compiles**

```bash
cd packages/observer-web && bunx tsc --noEmit
```

**Step 4: Commit**

```bash
git add packages/observer-web/src/routes/_app/projects/\$projectId/people/index.tsx
git commit -m "Add search button to people table"
```

---

### Task 5: Wire `onSearch` in pets table page

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/pets/-pets-page.tsx`

**Step 1: Destructure `refetch` from `usePets`**

Find (line ~78):
```tsx
const { data, isLoading } = usePets(projectId, params);
```

Replace with:
```tsx
const { data, isLoading, refetch } = usePets(projectId, params);
```

**Step 2: Pass to `DataTablePage`**

```tsx
onSearch={refetch}
```

**Step 3: Compile check + commit**

```bash
cd packages/observer-web && bunx tsc --noEmit
git add packages/observer-web/src/routes/_app/projects/\$projectId/pets/-pets-page.tsx
git commit -m "Add search button to pets table"
```

---

### Task 6: Wire `onSearch` in households table page

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/households/index.tsx`

**Step 1: Destructure `refetch` from `useHouseholds`**

Find (line ~53):
```tsx
const { data, isLoading } = useHouseholds(projectId, params);
```

Replace with:
```tsx
const { data, isLoading, refetch } = useHouseholds(projectId, params);
```

**Step 2: Pass to `DataTablePage`**

```tsx
onSearch={refetch}
```

**Step 3: Compile check + commit**

```bash
cd packages/observer-web && bunx tsc --noEmit
git add packages/observer-web/src/routes/_app/projects/\$projectId/households/index.tsx
git commit -m "Add search button to households table"
```

---

### Task 7: Wire `onSearch` in audit logs page

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/audit-logs.tsx`

**Step 1: Destructure `refetch` from `useProjectAuditLogs`**

Find (line ~49):
```tsx
const { data, isLoading } = useProjectAuditLogs(projectId, params);
```

Replace with:
```tsx
const { data, isLoading, refetch } = useProjectAuditLogs(projectId, params);
```

**Step 2: Pass to `DataTablePage`**

```tsx
onSearch={refetch}
```

**Step 3: Compile check + commit**

```bash
cd packages/observer-web && bunx tsc --noEmit
git add packages/observer-web/src/routes/_app/projects/\$projectId/audit-logs.tsx
git commit -m "Add search button to audit logs table"
```

---

### Task 8: Add search button inline to people report page

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/reports/people.tsx`

**Step 1: Destructure `refetch` from `useReport`**

Find (line ~69):
```tsx
const { data, isLoading } = useReport(projectId, params);
```

Replace with:
```tsx
const { data, isLoading, refetch } = useReport(projectId, params);
```

**Step 2: Add Button import at top of file**

```tsx
import { Button } from "@/components/button";
```

**Step 3: Add button at end of filter row**

The filter row div ends at line ~263:
```tsx
            </div>  {/* closes flex flex-wrap items-start gap-4 */}
```

Before that closing `</div>`, add:
```tsx
              <div className="self-end">
                <Button variant="secondary" onClick={() => refetch()}>
                  {t("common.search")}
                </Button>
              </div>
```

**Step 4: Compile check + commit**

```bash
cd packages/observer-web && bunx tsc --noEmit
git add packages/observer-web/src/routes/_app/projects/\$projectId/reports/people.tsx
git commit -m "Add search button to people report"
```

---

### Task 9: Add search button inline to pets report page

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/reports/pets.tsx`

**Step 1: Destructure `refetch` from `usePetReport`**

Find (line ~221):
```tsx
const { data, isLoading } = usePetReport(projectId, params);
```

Replace with:
```tsx
const { data, isLoading, refetch } = usePetReport(projectId, params);
```

**Step 2: Add Button import**

```tsx
import { Button } from "@/components/button";
```

**Step 3: Add button at end of filter row (line ~317)**

Before the closing `</div>` of `flex flex-wrap items-start gap-4`:
```tsx
              <div className="self-end">
                <Button variant="secondary" onClick={() => refetch()}>
                  {t("common.search")}
                </Button>
              </div>
```

**Step 4: Compile check + commit**

```bash
cd packages/observer-web && bunx tsc --noEmit
git add packages/observer-web/src/routes/_app/projects/\$projectId/reports/pets.tsx
git commit -m "Add search button to pets report"
```

---

### Task 10: Final verification

**Step 1: Full TypeScript check**

```bash
cd packages/observer-web && bunx tsc --noEmit
```

Expected: zero errors

**Step 2: Run unit tests**

```bash
cd /Users/sultan/Projects/observer && just test
```

Expected: all pass

**Step 3: Build check**

```bash
cd packages/observer-web && bun run build
```

Expected: successful build, no type errors
