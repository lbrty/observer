# Frontend Component Extraction Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Break large frontend files into small, self-documenting pieces that a newcomer can understand at a glance.

**Architecture:** Extract shared layout compounds (FormSection, DataTablePage), domain constants, and report sub-components. Split 5 drawers into folder-based section components. Refactor report pages to share duplicated utilities. Extract permission dialogs into separate files.

**Tech Stack:** React, TypeScript, TanStack React Table, TanStack Router, react-i18next, base-ui, Tailwind CSS

---

## Phase 1: Foundation Components

### Task 1: Create FormSection compound component

This is the shared layout used by all drawer sections. It replaces the repeated `SectionHeading + grid` pattern.

**Files:**
- Create: `packages/observer-web/src/components/form-section.tsx`

**Step 1: Create FormSection**

```tsx
import type { ReactNode } from "react";
import { SectionHeading } from "@/components/section-heading";

interface FormSectionProps {
  title: string;
  columns?: 1 | 2 | 3;
  children: ReactNode;
  className?: string;
}

const gridCols = {
  1: "grid-cols-1",
  2: "grid-cols-1 sm:grid-cols-2",
  3: "grid-cols-1 sm:grid-cols-3",
};

export function FormSection({ title, columns = 2, children, className }: FormSectionProps) {
  return (
    <>
      <SectionHeading>{title}</SectionHeading>
      <div className={`grid gap-4 ${gridCols[columns]} ${className ?? ""}`}>
        {children}
      </div>
    </>
  );
}
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`
Expected: no errors

**Step 3: Commit**

```bash
git add packages/observer-web/src/components/form-section.tsx
git commit -m "add FormSection compound component"
```

---

### Task 2: Extract domain constants — person

Extract the sex, age-group, and case-status option arrays that are duplicated in person-drawer, reports, and my-stats.

**Files:**
- Create: `packages/observer-web/src/constants/person.ts`

**Step 1: Create person constants**

These are the i18n key maps and value arrays used across person-drawer.tsx:194-218, reports/people.tsx:42-56,218-231, and my-stats/index.tsx:33-52.

```tsx
export const sexKeys: Record<string, string> = {
  male: "project.people.sexMale",
  female: "project.people.sexFemale",
  other: "project.people.sexOther",
  unknown: "project.people.sexUnknown",
};

export const SEX_VALUES = ["male", "female", "other", "unknown"] as const;

export const ageGroupKeys: Record<string, string> = {
  infant: "project.people.ageInfant",
  toddler: "project.people.ageToddler",
  pre_school: "project.people.agePreSchool",
  middle_childhood: "project.people.ageMiddleChildhood",
  young_teen: "project.people.ageYoungTeen",
  teenager: "project.people.ageTeenager",
  young_adult: "project.people.ageYoungAdult",
  early_adult: "project.people.ageEarlyAdult",
  middle_aged_adult: "project.people.ageMiddleAgedAdult",
  old_adult: "project.people.ageOldAdult",
};

export const AGE_GROUP_VALUES = [
  "infant",
  "toddler",
  "pre_school",
  "middle_childhood",
  "young_teen",
  "teenager",
  "young_adult",
  "early_adult",
  "middle_aged_adult",
  "old_adult",
] as const;

export const AGE_RANGE_MAP: Record<string, string> = {
  infant: "0-1",
  toddler: "1-3",
  pre_school: "3-6",
  middle_childhood: "6-12",
  young_teen: "12-14",
  teenager: "14-18",
  young_adult: "18-25",
  early_adult: "25-35",
  middle_aged_adult: "35-55",
  old_adult: "55+",
};

export const caseStatusKeys: Record<string, string> = {
  new: "project.people.new",
  active: "project.people.active",
  closed: "project.people.closed",
  archived: "project.people.archived",
};

export const CASE_STATUS_VALUES = ["new", "active", "closed", "archived"] as const;
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Commit**

```bash
git add packages/observer-web/src/constants/person.ts
git commit -m "extract person domain constants"
```

---

### Task 3: Extract domain constants — migration, household, pet, user

**Files:**
- Create: `packages/observer-web/src/constants/migration.ts`
- Create: `packages/observer-web/src/constants/household.ts`
- Create: `packages/observer-web/src/constants/pet.ts`
- Create: `packages/observer-web/src/constants/user.ts`

**Step 1: Create migration constants**

From migration-record-drawer.tsx:192-210:

```tsx
export const reasonKeys: Record<string, string> = {
  conflict: "project.migrationRecords.reasonConflict",
  security: "project.migrationRecords.reasonSecurity",
  service_access: "project.migrationRecords.reasonService",
  return: "project.migrationRecords.reasonReturn",
  relocation_program: "project.migrationRecords.reasonRelocation",
  economic: "project.migrationRecords.reasonEconomic",
  other: "project.migrationRecords.reasonOther",
};

export const housingKeys: Record<string, string> = {
  own_property: "project.migrationRecords.housingOwn",
  renting: "project.migrationRecords.housingRenting",
  with_relatives: "project.migrationRecords.housingRelatives",
  collective_site: "project.migrationRecords.housingCollective",
  hotel: "project.migrationRecords.housingHotel",
  other: "project.migrationRecords.housingOther",
  unknown: "project.migrationRecords.housingUnknown",
};
```

**Step 2: Create household constants**

From household-drawer.tsx:156-178:

```tsx
export const relationshipKeys: Record<string, string> = {
  head: "project.households.relationshipHead",
  spouse: "project.households.relationshipSpouse",
  child: "project.households.relationshipChild",
  parent: "project.households.relationshipParent",
  sibling: "project.households.relationshipSibling",
  grandchild: "project.households.relationshipGrandchild",
  grandparent: "project.households.relationshipGrandparent",
  other_relative: "project.households.relationshipOtherRelative",
  non_relative: "project.households.relationshipNonRelative",
};
```

**Step 3: Create pet constants**

From pet-drawer.tsx:117-123:

```tsx
export const petStatusKeys: Record<string, string> = {
  registered: "project.pets.statusRegistered",
  adopted: "project.pets.statusAdopted",
  owner_found: "project.pets.statusOwnerFound",
  needs_shelter: "project.pets.statusNeedsShelter",
  unknown: "project.pets.statusUnknown",
};
```

**Step 4: Create user constants**

From admin/users/index.tsx:46-58 and permissions.tsx:210-218:

```tsx
export const roleKeys: Record<string, string> = {
  admin: "admin.users.roleAdmin",
  staff: "admin.users.roleStaff",
  consultant: "admin.users.roleConsultant",
  guest: "admin.users.roleGuest",
};

export const projectRoleKeys: Record<string, string> = {
  owner: "admin.permissions.roleOwner",
  manager: "admin.permissions.roleManager",
  consultant: "admin.permissions.roleConsultant",
  viewer: "admin.permissions.roleViewer",
};
```

**Step 5: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 6: Commit**

```bash
git add packages/observer-web/src/constants/migration.ts packages/observer-web/src/constants/household.ts packages/observer-web/src/constants/pet.ts packages/observer-web/src/constants/user.ts
git commit -m "extract migration, household, pet, user domain constants"
```

---

## Phase 2: Report Shared Components

The reports/people.tsx (741 lines) and my-stats/index.tsx (503 lines) share heavily duplicated code: `ReportCard`/`StatsCard`, `KpiCard`, `FilterChip`, `FilterField`, `getPresetDates`, `labelKeyMap`, `useTranslatedRows`, `AGE_RANGE_MAP`, skeleton components. Extract these to shared report components.

### Task 4: Extract shared report utilities and label maps

**Files:**
- Create: `packages/observer-web/src/components/report/label-maps.ts`

**Step 1: Create shared label maps and utilities**

This consolidates the duplicated `labelKeyMap`, `useTranslatedRows`, `AGE_RANGE_MAP`, and date preset logic from both reports/people.tsx and my-stats/index.tsx.

```tsx
import { useTranslation } from "react-i18next";
import { typeKeys, sphereKeys, referralKeys } from "@/constants/support";
import { sexKeys, ageGroupKeys, AGE_RANGE_MAP as AGE_RANGES } from "@/constants/person";
import type { CountResult } from "@/types/report";

export { AGE_RANGES as AGE_RANGE_MAP };

export const labelKeyMap: Record<string, string> = {
  ...typeKeys,
  ...sphereKeys,
  unspecified: "project.supportRecords.sphereOther",
  ...referralKeys,
  ...sexKeys,
  ...ageGroupKeys,
};

export function useTranslatedRows(rows: CountResult[]): CountResult[] {
  const { t } = useTranslation();
  const translated = rows.map((r) => {
    const key = labelKeyMap[r.label];
    return key ? { ...r, label: t(key) } : r;
  });
  const merged = new Map<string, number>();
  for (const r of translated) {
    merged.set(r.label, (merged.get(r.label) ?? 0) + r.count);
  }
  return Array.from(merged, ([label, count]) => ({ label, count }));
}

export type DatePreset = "month" | "quarter" | "year" | "all";

export function getPresetDates(preset: DatePreset): { date_from?: string; date_to?: string } {
  const now = new Date();
  const fmt = (d: Date) => d.toISOString().slice(0, 10);
  const today = fmt(now);

  switch (preset) {
    case "month": {
      const from = new Date(now.getFullYear(), now.getMonth(), 1);
      return { date_from: fmt(from), date_to: today };
    }
    case "quarter": {
      const qMonth = Math.floor(now.getMonth() / 3) * 3 - 3;
      const from = new Date(now.getFullYear(), qMonth, 1);
      const to = new Date(now.getFullYear(), qMonth + 3, 0);
      return { date_from: fmt(from), date_to: fmt(to) };
    }
    case "year": {
      const from = new Date(now.getFullYear(), 0, 1);
      return { date_from: fmt(from), date_to: today };
    }
    case "all":
      return { date_from: undefined, date_to: undefined };
  }
}

export const PRESET_KEYS: { key: DatePreset; i18n: string }[] = [
  { key: "month", i18n: "project.reports.presetMonth" },
  { key: "quarter", i18n: "project.reports.presetQuarter" },
  { key: "year", i18n: "project.reports.presetYear" },
  { key: "all", i18n: "project.reports.presetAll" },
];
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Commit**

```bash
git add packages/observer-web/src/components/report/label-maps.ts
git commit -m "extract shared report label maps and date utilities"
```

---

### Task 5: Extract ReportCard, KpiCard, FilterChip, FilterField, ReportSkeleton

**Files:**
- Create: `packages/observer-web/src/components/report/report-card.tsx`
- Create: `packages/observer-web/src/components/report/kpi-card.tsx`
- Create: `packages/observer-web/src/components/report/filter-chip.tsx`
- Create: `packages/observer-web/src/components/report/filter-field.tsx`
- Create: `packages/observer-web/src/components/report/report-skeleton.tsx`
- Create: `packages/observer-web/src/components/report/index.ts`

**Step 1: Create ReportCard**

This replaces both `ReportCard` from reports/people.tsx:71-130 and the identical `StatsCard` from my-stats/index.tsx:67-126.

```tsx
import { BarChart, type BarLegendItem } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import { DownloadSimpleIcon } from "@/components/icons";
import { exportGroupCSV } from "@/lib/export-csv";
import { useTranslatedRows } from "@/components/report/label-maps";
import type { ReportGroup } from "@/types/report";

interface ReportCardProps {
  group: ReportGroup;
  title: string;
  chart: "bar" | "pie";
  yAxisLabel?: string;
  legend?: BarLegendItem[];
  mapLabel?: (label: string) => string;
  skipTranslation?: boolean;
  colorMap?: Record<string, string>;
  direction?: "vertical" | "horizontal" | "auto";
}

export function ReportCard({
  group,
  title,
  chart,
  yAxisLabel,
  legend,
  mapLabel,
  skipTranslation,
  colorMap,
  direction,
}: ReportCardProps) {
  const translated = useTranslatedRows(group.rows);
  const source = skipTranslation ? group.rows : translated;
  const rows = mapLabel ? source.map((r) => ({ ...r, label: mapLabel(r.label) })) : source;
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-3 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => exportGroupCSV(title, rows)}
            className="text-fg-tertiary transition-colors hover:text-fg"
            title="Download CSV"
          >
            <DownloadSimpleIcon size={14} />
          </button>
          <span className="tabular-nums text-xs font-medium text-fg-tertiary">
            {group.total.toLocaleString()}
          </span>
        </div>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart
            data={rows}
            yAxisLabel={yAxisLabel}
            legend={legend}
            colorMap={colorMap}
            direction={direction}
          />
        ) : (
          <PieChart data={rows} colorMap={colorMap} />
        )
      ) : (
        <p className="py-8 text-center text-sm text-fg-tertiary">&mdash;</p>
      )}
    </div>
  );
}
```

**Step 2: Create KpiCard**

From reports/people.tsx:149-156, identical to my-stats/index.tsx:128-135:

```tsx
export function KpiCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-4">
      <p className="text-2xl font-bold tabular-nums text-fg">{value.toLocaleString()}</p>
      <p className="mt-0.5 text-xs font-medium text-fg-tertiary">{label}</p>
    </div>
  );
}
```

**Step 3: Create FilterChip**

From reports/people.tsx:158-177, identical to my-stats/index.tsx:137-156:

```tsx
import { XIcon } from "@/components/icons";

export function FilterChip({
  label,
  value,
  onRemove,
}: {
  label: string;
  value: string;
  onRemove: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onRemove}
      className="inline-flex items-center gap-1 rounded-md bg-bg-tertiary px-2 py-0.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
    >
      <span className="text-fg-tertiary">{label}:</span> {value}
      <XIcon size={10} />
    </button>
  );
}
```

**Step 4: Create FilterField**

From reports/people.tsx:132-139, identical to my-stats/index.tsx:197-204:

```tsx
import type { ReactNode } from "react";

export function FilterField({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="space-y-1.5">
      <span className="block text-xs font-medium text-fg-secondary">{label}</span>
      {children}
    </div>
  );
}
```

**Step 5: Create ReportSkeleton**

A configurable version combining reports/people.tsx:179-195 and my-stats/index.tsx:158-173:

```tsx
export function ReportSkeleton({ kpiCount = 6, chartCount = 4 }: { kpiCount?: number; chartCount?: number }) {
  return (
    <div className="space-y-6">
      <div className={`grid gap-4 ${kpiCount > 4 ? "grid-cols-3 lg:grid-cols-6" : "grid-cols-2 sm:grid-cols-4"}`}>
        {Array.from({ length: kpiCount }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
      <div className="grid gap-6 lg:grid-cols-2">
        {Array.from({ length: chartCount }).map((_, i) => (
          <div key={i} className="h-72 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
    </div>
  );
}
```

**Step 6: Create barrel export**

```tsx
export { ReportCard } from "./report-card";
export { KpiCard } from "./kpi-card";
export { FilterChip } from "./filter-chip";
export { FilterField } from "./filter-field";
export { ReportSkeleton } from "./report-skeleton";
export { labelKeyMap, useTranslatedRows, getPresetDates, PRESET_KEYS, AGE_RANGE_MAP } from "./label-maps";
export type { DatePreset } from "./label-maps";
```

**Step 7: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 8: Commit**

```bash
git add packages/observer-web/src/components/report/
git commit -m "extract shared report components"
```

---

### Task 6: Refactor reports/people.tsx to use shared report components

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/reports/people.tsx`

**Step 1: Replace inline components with imports**

Update imports at the top of the file. Remove the duplicated `ReportCard`, `KpiCard`, `FilterChip`, `FilterField`, `ReportSkeleton`, `SectionHeader`, `labelKeyMap`, `useTranslatedRows`, `AGE_RANGE_MAP`, `getPresetDates`, `PRESET_KEYS`, `DatePreset` definitions. Import them from `@/components/report` instead.

The `SectionHeader` component (reports/people.tsx:141-147) stays inline since it's report-specific layout (border-bottom separator), or can be extracted to `@/components/report/section-header.tsx`.

After refactoring, the page file should contain only:
1. Route definition
2. `ReportsPage` function with state, hooks, option arrays, and JSX composition
3. `SectionHeader` (if kept inline)

The file should drop from ~741 to ~350 lines.

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Verify the page renders**

Run: `cd packages/observer-web && bun run dev`
Navigate to a project's reports page and verify charts, filters, KPIs all render correctly.

**Step 4: Commit**

```bash
git add packages/observer-web/src/routes/_app/projects/\$projectId/reports/people.tsx
git commit -m "refactor reports/people to use shared report components"
```

---

### Task 7: Refactor my-stats to use shared report components

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/my-stats/index.tsx`

**Step 1: Replace inline components with imports**

Same approach as Task 6. Remove duplicated `StatsCard` (use `ReportCard` instead), `KpiCard`, `FilterChip`, `FilterField`, `StatsSkeleton` (use `ReportSkeleton` instead), `labelKeyMap`, `useTranslatedRows`, `AGE_RANGE_MAP`, `getPresetDates`, `PRESET_KEYS`, `DatePreset`, `SUPPORT_TYPE_OPTIONS`, `FilterField`.

The file should drop from ~503 to ~180 lines.

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Verify the page renders**

Run: `cd packages/observer-web && bun run dev`
Navigate to a project's my-stats page and verify charts, filters, KPIs all render correctly.

**Step 4: Commit**

```bash
git add packages/observer-web/src/routes/_app/projects/\$projectId/my-stats/index.tsx
git commit -m "refactor my-stats to use shared report components"
```

---

## Phase 3: DataTablePage Compound

### Task 8: Create FilterBar component

**Files:**
- Create: `packages/observer-web/src/components/filter-bar.tsx`

**Step 1: Create FilterBar**

Extract the repeated search-input + filter-selects pattern from admin/users/index.tsx:115-149 and similar pages.

```tsx
import type { ReactNode } from "react";

import { MagnifyingGlassIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";

export interface SearchFilter {
  type: "search";
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
}

export interface SelectFilter {
  type: "select";
  value: string;
  onValueChange: (value: string) => void;
  options: { label: string; value: string }[];
  placeholder?: string;
}

export type FilterDef = SearchFilter | SelectFilter;

interface FilterBarProps {
  filters: FilterDef[];
  trailing?: ReactNode;
}

export function FilterBar({ filters, trailing }: FilterBarProps) {
  return (
    <div className="mb-4 flex gap-3">
      {filters.map((f, i) => {
        if (f.type === "search") {
          return (
            <div key={i} className="relative">
              <MagnifyingGlassIcon
                size={14}
                className="absolute top-1/2 left-3 -translate-y-1/2 text-fg-tertiary"
              />
              <input
                placeholder={f.placeholder}
                value={f.value}
                onChange={(e) => f.onChange(e.target.value)}
                className="rounded-lg border border-border-secondary bg-bg-secondary py-1.5 pr-3 pl-8 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
              />
            </div>
          );
        }
        return (
          <UISelect
            key={i}
            value={f.value}
            onValueChange={f.onValueChange}
            options={f.options}
            placeholder={f.placeholder}
          />
        );
      })}
      {trailing}
    </div>
  );
}
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Commit**

```bash
git add packages/observer-web/src/components/filter-bar.tsx
git commit -m "add FilterBar component"
```

---

### Task 9: Create DataTablePage compound component

**Files:**
- Create: `packages/observer-web/src/components/data-table-page.tsx`

**Step 1: Create DataTablePage**

This wraps PageHeader + FilterBar + DataTable + Pagination + EmptyState into a single compound.

```tsx
import type { ReactNode } from "react";

import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { type FilterDef, FilterBar } from "@/components/filter-bar";
import type { Icon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";

interface PaginationConfig {
  page: number;
  perPage: number;
  total: number;
  onChange: (page: number) => void;
}

interface DataTablePageProps<T> {
  title: string;
  columns: Column<T>[];
  data: T[];
  keyExtractor: (item: T) => string;
  isLoading?: boolean;
  onRowClick?: (item: T) => void;
  pagination?: PaginationConfig;
  filters?: FilterDef[];
  filterTrailing?: ReactNode;
  emptyIcon?: Icon;
  emptyTitle?: string;
  emptyDescription?: string;
  emptyAction?: ReactNode;
  createAction?: ReactNode;
  children?: ReactNode;
}

export function DataTablePage<T>({
  title,
  columns,
  data,
  keyExtractor,
  isLoading,
  onRowClick,
  pagination,
  filters,
  filterTrailing,
  emptyIcon,
  emptyTitle,
  emptyDescription,
  emptyAction,
  createAction,
  children,
}: DataTablePageProps<T>) {
  return (
    <div>
      <PageHeader title={title} action={createAction} />

      {filters && filters.length > 0 && (
        <FilterBar filters={filters} trailing={filterTrailing} />
      )}

      <DataTable
        columns={columns}
        data={data}
        keyExtractor={keyExtractor}
        onRowClick={onRowClick}
        isLoading={isLoading}
        emptyState={
          emptyIcon && emptyTitle ? (
            <EmptyState
              icon={emptyIcon}
              title={emptyTitle}
              description={emptyDescription}
              action={emptyAction}
            />
          ) : undefined
        }
      />

      {pagination && (
        <Pagination
          page={pagination.page}
          perPage={pagination.perPage}
          total={pagination.total}
          onChange={pagination.onChange}
        />
      )}

      {children}
    </div>
  );
}
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Commit**

```bash
git add packages/observer-web/src/components/data-table-page.tsx
git commit -m "add DataTablePage compound component"
```

---

### Task 10: Refactor admin/users page to use DataTablePage

**Files:**
- Modify: `packages/observer-web/src/routes/_app/admin/users/index.tsx`

**Step 1: Replace manual layout with DataTablePage**

Replace the manual `PageHeader` + search input + UISelect filters + `DataTable` + `Pagination` block (lines 104-176) with a single `<DataTablePage>` call. Keep the `CreateUserDialog` as-is (it's already a separate function).

The page should drop from ~327 to ~200 lines. The `UsersPage` function body becomes:
1. State (page, search, role, isActive, createOpen)
2. Hook calls (useUsers, useNavigate)
3. Option arrays
4. Column definitions
5. Single `<DataTablePage>` JSX with `filters` and `children` for the dialog

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Verify the page renders**

Run: `cd packages/observer-web && bun run dev`
Navigate to /admin/users and verify search, filters, table, pagination, row click, and create dialog all work.

**Step 4: Commit**

```bash
git add packages/observer-web/src/routes/_app/admin/users/index.tsx
git commit -m "refactor admin/users to use DataTablePage"
```

---

### Task 11: Refactor remaining list pages to use DataTablePage

Apply the same `DataTablePage` refactoring to all other list pages that follow the pattern. These are independent and can be done in one batch or one at a time.

**Files to modify:**
- `packages/observer-web/src/routes/_app/admin/projects/index.tsx`
- `packages/observer-web/src/routes/_app/admin/reference/countries/index.tsx`
- `packages/observer-web/src/routes/_app/admin/reference/countries/$countryId/index.tsx` (states list)
- `packages/observer-web/src/routes/_app/admin/reference/countries/$countryId/states/$stateId.tsx` (places list)
- `packages/observer-web/src/routes/_app/admin/reference/offices.tsx`
- `packages/observer-web/src/routes/_app/admin/reference/categories.tsx`
- `packages/observer-web/src/routes/_app/projects/$projectId/tags/index.tsx`
- `packages/observer-web/src/routes/_app/projects/$projectId/people/index.tsx`
- `packages/observer-web/src/routes/_app/projects/$projectId/support-records/-support-records-page.tsx`
- `packages/observer-web/src/routes/_app/projects/$projectId/pets/-pets-page.tsx`

**Step 1: For each page, apply the same pattern**

Read the page, identify the PageHeader + filters + DataTable + Pagination block, and replace with `<DataTablePage>`. Keep page-specific dialogs, drawers, and tab logic as-is.

Some pages may have additional layout (tabs, nested routes) that doesn't fit the compound — use `children` prop or keep those parts outside. Don't force pages into the compound if they have unusual layout.

**Step 2: Verify build after each page**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 3: Spot-check a few pages in dev**

Run: `cd packages/observer-web && bun run dev`

**Step 4: Commit per batch of 2-3 pages**

```bash
git commit -m "refactor admin list pages to use DataTablePage"
git commit -m "refactor project list pages to use DataTablePage"
```

---

## Phase 4: Drawer Extraction

### Task 12: Extract person-drawer into folder with sections

**Files:**
- Create: `packages/observer-web/src/components/person-drawer/identity-section.tsx`
- Create: `packages/observer-web/src/components/person-drawer/location-section.tsx`
- Create: `packages/observer-web/src/components/person-drawer/case-section.tsx`
- Modify: `packages/observer-web/src/components/person-drawer.tsx` → move to `packages/observer-web/src/components/person-drawer/index.tsx`

**Step 1: Create identity-section.tsx**

Extracts lines 232-299 from person-drawer.tsx. Receives form state and set function as props.

```tsx
import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { FormField } from "@/components/form-field";
import { FormSection } from "@/components/form-section";
import { UISelect } from "@/components/ui-select";
import { sexKeys, ageGroupKeys } from "@/constants/person";

interface IdentitySectionProps {
  form: {
    first_name: string;
    last_name: string;
    patronymic: string;
    sex: string;
    birth_date: string;
    age_group: string;
    primary_phone: string;
    email: string;
  };
  set: (key: string, value: string) => void;
}

export function IdentitySection({ form, set }: IdentitySectionProps) {
  const { t } = useTranslation();

  const sexOptions = Object.entries(sexKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  const ageGroupOptions = Object.entries(ageGroupKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <FormSection title={t("project.people.identity")}>
      <FormField
        label={t("project.people.firstName")}
        value={form.first_name}
        onChange={(v) => set("first_name", v)}
        required
      />
      <FormField
        label={t("project.people.lastName")}
        value={form.last_name}
        onChange={(v) => set("last_name", v)}
      />
      <FormField
        label={t("project.people.patronymic")}
        value={form.patronymic}
        onChange={(v) => set("patronymic", v)}
      />

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.sexLabel")}
        </Field.Label>
        <UISelect
          value={form.sex}
          onValueChange={(v) => set("sex", v)}
          options={sexOptions}
          fullWidth
        />
      </Field.Root>

      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.birthDate")}
        </span>
        <DatePicker
          value={form.birth_date}
          onChange={(v) => set("birth_date", v)}
          captionLayout="dropdown"
          clearable
        />
      </div>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.ageGroup")}
        </Field.Label>
        <UISelect
          value={form.age_group}
          onValueChange={(v) => set("age_group", v)}
          options={ageGroupOptions}
          fullWidth
          clearable
        />
      </Field.Root>

      <FormField
        label={t("project.people.phone")}
        value={form.primary_phone}
        onChange={(v) => set("primary_phone", v)}
      />
      <FormField
        label={t("project.people.email")}
        value={form.email}
        onChange={(v) => set("email", v)}
        type="email"
      />
    </FormSection>
  );
}
```

**Step 2: Create location-section.tsx**

Extracts lines 301-362. Receives form, set, place labels, and place resolution as props.

```tsx
import { useTranslation } from "react-i18next";

import { FormSection } from "@/components/form-section";
import { PlaceCombobox } from "@/components/place-combobox";

interface LocationSectionProps {
  originPlaceId: string;
  currentPlaceId: string;
  originLabel: string;
  currentLabel: string;
  onOriginSelect: (placeId: string, label: string) => void;
  onCurrentSelect: (placeId: string, label: string) => void;
  onOriginClear: () => void;
  onCurrentClear: () => void;
}

export function LocationSection({
  originPlaceId,
  currentPlaceId,
  originLabel,
  currentLabel,
  onOriginSelect,
  onCurrentSelect,
  onOriginClear,
  onCurrentClear,
}: LocationSectionProps) {
  const { t } = useTranslation();

  return (
    <FormSection title={t("project.people.location")}>
      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.originPlace")}
        </span>
        {originPlaceId ? (
          <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
            <span className="flex-1 truncate text-sm text-fg">
              {originLabel || originPlaceId}
            </span>
            <button
              type="button"
              onClick={onOriginClear}
              className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
            >
              ×
            </button>
          </div>
        ) : (
          <PlaceCombobox
            onSelect={(place, state, country) => {
              onOriginSelect(place.id, `${place.name}, ${state.name}, ${country.name}`);
            }}
          />
        )}
      </div>

      <div>
        <span className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.currentPlace")}
        </span>
        {currentPlaceId ? (
          <div className="flex h-9 items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3">
            <span className="flex-1 truncate text-sm text-fg">
              {currentLabel || currentPlaceId}
            </span>
            <button
              type="button"
              onClick={onCurrentClear}
              className="shrink-0 cursor-pointer text-fg-tertiary hover:text-fg"
            >
              ×
            </button>
          </div>
        ) : (
          <PlaceCombobox
            onSelect={(place, state, country) => {
              onCurrentSelect(place.id, `${place.name}, ${state.name}, ${country.name}`);
            }}
          />
        )}
      </div>
    </FormSection>
  );
}
```

**Step 3: Create case-section.tsx**

Extracts lines 364-414. Receives form, set, officeOptions as props.

```tsx
import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { FormField } from "@/components/form-field";
import { FormSection } from "@/components/form-section";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { caseStatusKeys } from "@/constants/person";

interface CaseSectionProps {
  form: {
    case_status: string;
    external_id: string;
    office_id: string;
    consent_given: boolean;
    consent_date: string;
  };
  set: (key: string, value: string | boolean) => void;
  officeOptions: { label: string; value: string }[];
}

export function CaseSection({ form, set, officeOptions }: CaseSectionProps) {
  const { t } = useTranslation();

  const caseStatusOptions = Object.entries(caseStatusKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));

  return (
    <FormSection title={t("project.people.case")}>
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("project.people.caseStatusLabel")}
        </Field.Label>
        <UISelect
          value={form.case_status}
          onValueChange={(v) => set("case_status", v)}
          options={caseStatusOptions}
          fullWidth
        />
      </Field.Root>

      <FormField
        label={t("project.people.externalId")}
        value={form.external_id}
        onChange={(v) => set("external_id", v)}
      />

      {officeOptions.length > 0 && (
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.people.office")}
          </Field.Label>
          <UISelect
            value={form.office_id}
            onValueChange={(v) => set("office_id", v)}
            options={officeOptions}
            fullWidth
          />
        </Field.Root>
      )}

      <div className="col-span-full space-y-4">
        <UISwitch
          checked={form.consent_given}
          onCheckedChange={(v) => set("consent_given", v)}
          label={t("project.people.consentGiven")}
        />

        {form.consent_given && (
          <div>
            <span className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("project.people.consentDate")}
            </span>
            <DatePicker value={form.consent_date} onChange={(v) => set("consent_date", v)} />
          </div>
        )}
      </div>
    </FormSection>
  );
}
```

**Step 4: Move person-drawer.tsx to person-drawer/index.tsx and refactor**

1. `mkdir -p packages/observer-web/src/components/person-drawer`
2. `git mv packages/observer-web/src/components/person-drawer.tsx packages/observer-web/src/components/person-drawer/index.tsx`
3. Replace the JSX body (lines 230-418) with:

```tsx
<DrawerShell ...>
  <ErrorBanner message={error} />
  <IdentitySection form={form} set={set} />
  <LocationSection
    originPlaceId={form.origin_place_id}
    currentPlaceId={form.current_place_id}
    originLabel={resolvedOriginLabel}
    currentLabel={resolvedCurrentLabel}
    onOriginSelect={(id, label) => { set("origin_place_id", id); setOriginPlaceLabel(label); }}
    onCurrentSelect={(id, label) => { set("current_place_id", id); setCurrentPlaceLabel(label); }}
    onOriginClear={() => { set("origin_place_id", ""); setOriginPlaceLabel(""); }}
    onCurrentClear={() => { set("current_place_id", ""); setCurrentPlaceLabel(""); }}
  />
  <CaseSection form={form} set={set} officeOptions={officeOptions} />
  <SectionHeading>{t("project.tags.title")}</SectionHeading>
  <TagPicker projectId={projectId} selectedIds={tagIds} onChange={setTagIds} />
</DrawerShell>
```

Remove the inline sexOptions, caseStatusOptions, ageGroupOptions arrays (they're now in the section components via constants).

**Step 5: Verify build**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 6: Verify the drawer renders**

Run: `cd packages/observer-web && bun run dev`
Navigate to a project's people page, open the create/edit drawer, verify all sections work.

**Step 7: Commit**

```bash
git add packages/observer-web/src/components/person-drawer/
git commit -m "extract person-drawer into folder with section components"
```

---

### Task 13: Extract support-record-drawer into folder with sections

Same pattern as Task 12. Extract from the 393-line file.

**Files:**
- Create: `packages/observer-web/src/components/support-record-drawer/info-section.tsx` (lines 270-350: type, sphere, date, person picker)
- Create: `packages/observer-web/src/components/support-record-drawer/referral-section.tsx` (lines 352-381: referral status, office)
- Move: `packages/observer-web/src/components/support-record-drawer.tsx` → `packages/observer-web/src/components/support-record-drawer/index.tsx`

**Step 1: Create info-section.tsx**

Use `FormSection` for layout. Receive form, set, and person combobox props. Use `typeKeys` and `sphereKeys` from `@/constants/support` for options.

**Step 2: Create referral-section.tsx**

Use `FormSection`. Receive form, set, officeOptions. Use `referralKeys` from `@/constants/support`.

**Step 3: Move and refactor index.tsx**

1. `mkdir -p packages/observer-web/src/components/support-record-drawer`
2. `git mv packages/observer-web/src/components/support-record-drawer.tsx packages/observer-web/src/components/support-record-drawer/index.tsx`
3. Replace JSX body with `<InfoSection>`, `<ReferralSection>`, and the notes `FormTextarea`.
4. Remove inline option arrays (typeOptions, sphereOptions, referralStatusOptions).

**Step 4: Verify build and test**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 5: Commit**

```bash
git add packages/observer-web/src/components/support-record-drawer/
git commit -m "extract support-record-drawer into folder with section components"
```

---

### Task 14: Extract household-drawer into folder with sections

**Files:**
- Create: `packages/observer-web/src/components/household-drawer/head-section.tsx` (lines 192-229: reference number, head person picker)
- Create: `packages/observer-web/src/components/household-drawer/members-section.tsx` (lines 231-352: member table, add member form)
- Move: `packages/observer-web/src/components/household-drawer.tsx` → `packages/observer-web/src/components/household-drawer/index.tsx`

**Step 1: Create head-section.tsx**

Use `FormSection`. Receive form, set, headPersonLabel, person combobox callbacks.

**Step 2: Create members-section.tsx**

This is the most complex section — it has the member table + add-member form. Receive editingId, household, projectId, onAddMember, onRemoveMember, memberForm, setMemberForm, memberPersonName callbacks. Use `relationshipKeys` from `@/constants/household`.

**Step 3: Move and refactor index.tsx**

1. `mkdir -p packages/observer-web/src/components/household-drawer`
2. `git mv packages/observer-web/src/components/household-drawer.tsx packages/observer-web/src/components/household-drawer/index.tsx`
3. Replace JSX with `<HeadSection>` + `<MembersSection>`.

**Step 4: Verify build and test**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 5: Commit**

```bash
git add packages/observer-web/src/components/household-drawer/
git commit -m "extract household-drawer into folder with section components"
```

---

### Task 15: Extract migration-record-drawer into folder with sections

**Files:**
- Create: `packages/observer-web/src/components/migration-record-drawer/place-section.tsx` (shared component for both origin and destination — lines 228-294 follow the same 3-column country→state→place pattern)
- Create: `packages/observer-web/src/components/migration-record-drawer/details-section.tsx` (lines 296-328: date, reason, housing)
- Move: `packages/observer-web/src/components/migration-record-drawer.tsx` → `packages/observer-web/src/components/migration-record-drawer/index.tsx`

**Step 1: Create place-section.tsx**

A reusable 3-column cascading select for country→state→place. Used twice (origin, destination). Props: title, countryValue, stateValue, placeValue, onCountryChange, onStateChange, onPlaceChange, countryOptions, stateOptions, placeOptions.

```tsx
import { useTranslation } from "react-i18next";

import { SectionHeading } from "@/components/section-heading";
import { UISelect } from "@/components/ui-select";

interface PlaceSectionProps {
  title: string;
  country: string;
  state: string;
  place: string;
  onCountryChange: (v: string) => void;
  onStateChange: (v: string) => void;
  onPlaceChange: (v: string) => void;
  countryOptions: { label: string; value: string }[];
  stateOptions: { label: string; value: string }[];
  placeOptions: { label: string; value: string }[];
}

export function PlaceSection({
  title,
  country,
  state,
  place,
  onCountryChange,
  onStateChange,
  onPlaceChange,
  countryOptions,
  stateOptions,
  placeOptions,
}: PlaceSectionProps) {
  const { t } = useTranslation();

  return (
    <>
      <SectionHeading>{title}</SectionHeading>
      <div className="grid grid-cols-3 gap-3">
        <UISelect
          value={country}
          onValueChange={onCountryChange}
          options={countryOptions}
          placeholder={t("project.people.selectCountry")}
          fullWidth
        />
        <UISelect
          value={state}
          onValueChange={onStateChange}
          options={stateOptions}
          placeholder={t("project.people.selectState")}
          disabled={!country}
          fullWidth
        />
        <UISelect
          value={place}
          onValueChange={onPlaceChange}
          options={placeOptions}
          placeholder={t("project.people.selectPlace")}
          disabled={!state}
          fullWidth
        />
      </div>
    </>
  );
}
```

**Step 2: Create details-section.tsx**

Use `FormSection`. Receive form, set. Use `reasonKeys` and `housingKeys` from `@/constants/migration`.

**Step 3: Move and refactor index.tsx**

1. `mkdir -p packages/observer-web/src/components/migration-record-drawer`
2. `git mv packages/observer-web/src/components/migration-record-drawer.tsx packages/observer-web/src/components/migration-record-drawer/index.tsx`
3. Replace JSX with two `<PlaceSection>` instances + `<DetailsSection>` + notes textarea.

**Step 4: Verify build and test**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 5: Commit**

```bash
git add packages/observer-web/src/components/migration-record-drawer/
git commit -m "extract migration-record-drawer into folder with section components"
```

---

### Task 16: Extract pet-drawer into folder with sections

**Files:**
- Create: `packages/observer-web/src/components/pet-drawer/info-section.tsx` (lines 137-206: name, status, owner, registration ID, notes)
- Move: `packages/observer-web/src/components/pet-drawer.tsx` → `packages/observer-web/src/components/pet-drawer/index.tsx`

**Step 1: Create info-section.tsx**

Use `FormSection`. The pet drawer is simpler — just one section for info (name, status, owner, registrationId, notes). Use `petStatusKeys` from `@/constants/pet`.

**Step 2: Move and refactor index.tsx**

1. `mkdir -p packages/observer-web/src/components/pet-drawer`
2. `git mv packages/observer-web/src/components/pet-drawer.tsx packages/observer-web/src/components/pet-drawer/index.tsx`
3. Replace JSX with `<InfoSection>` + TagPicker.

**Step 3: Verify build and test**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 4: Commit**

```bash
git add packages/observer-web/src/components/pet-drawer/
git commit -m "extract pet-drawer into folder with section components"
```

---

## Phase 5: Permissions Page Extraction

### Task 17: Extract permissions dialogs into separate files

The permissions page (486 lines) has two large dialog functions (AssignDialog: 140 lines, EditDialog: 118 lines) that can be moved to their own files.

**Files:**
- Create: `packages/observer-web/src/routes/_app/admin/projects/$projectId/permissions/assign-dialog.tsx`
- Create: `packages/observer-web/src/routes/_app/admin/projects/$projectId/permissions/edit-dialog.tsx`
- Create: `packages/observer-web/src/routes/_app/admin/projects/$projectId/permissions/role-select.tsx`
- Modify: `packages/observer-web/src/routes/_app/admin/projects/$projectId/permissions.tsx`

**Step 1: Create role-select.tsx**

Extract `useRoleOptions` hook and `RoleDescription` component (lines 210-224) shared by both dialogs. Use `projectRoleKeys` from `@/constants/user`.

```tsx
import { useTranslation } from "react-i18next";
import { projectRoleKeys } from "@/constants/user";

export function useRoleOptions() {
  const { t } = useTranslation();
  return Object.entries(projectRoleKeys).map(([value, key]) => ({
    label: t(key),
    value,
  }));
}

export function RoleDescription({ role }: { role: string }) {
  const { t } = useTranslation();
  const key = `admin.permissions.role${role.charAt(0).toUpperCase() + role.slice(1)}Desc`;
  return <p className="text-xs text-fg-tertiary">{t(key)}</p>;
}
```

**Step 2: Create assign-dialog.tsx**

Move AssignDialog function (lines 226-366) to its own file. Import `useRoleOptions` and `RoleDescription` from `./role-select`.

**Step 3: Create edit-dialog.tsx**

Move EditDialog function (lines 368-485) to its own file. Import from `./role-select`.

**Step 4: Update permissions.tsx**

The permissions route file should now contain only the route definition and `PermissionsPage` function (~130 lines). It must create a `permissions/` directory structure:

Option A: If TanStack Router allows colocated files alongside the route, create the directory alongside the route file.
Option B: Put the dialogs in `@/components/permissions/` instead if the router structure doesn't support colocated files.

Check TanStack Router conventions — if `permissions.tsx` already defines the route, adding a `permissions/` folder could conflict. In that case, put the dialogs in `packages/observer-web/src/components/permissions/` instead.

**Step 5: Verify build and test**

Run: `cd packages/observer-web && bunx tsc --noEmit`

**Step 6: Commit**

```bash
git add packages/observer-web/src/routes/_app/admin/projects/\$projectId/permissions* packages/observer-web/src/components/permissions/
git commit -m "extract permissions dialogs into separate files"
```

---

## Phase 6: Verification

### Task 18: Full build and type check

**Step 1: Type check**

Run: `cd packages/observer-web && bunx tsc --noEmit`
Expected: no errors

**Step 2: Dev server**

Run: `cd packages/observer-web && bun run dev`
Check these pages:
- /admin/users — table, filters, create dialog
- /admin/projects — table
- Any project's people page — drawer sections
- Any project's support records — drawer sections
- Any project's households — drawer with members
- Any project's reports — charts, KPIs, filters
- Any project's my-stats — charts, KPIs
- Any project's permissions — assign/edit/revoke dialogs

**Step 3: Commit any fixes**

```bash
git commit -m "fix post-extraction issues"
```

---

## Dependency Graph

```
Phase 1 (Tasks 1-3): Foundation — no dependencies
Phase 2 (Tasks 4-7): Reports — depends on Task 2 (person constants) + existing support constants
Phase 3 (Tasks 8-11): DataTablePage — no dependencies on other phases
Phase 4 (Tasks 12-16): Drawers — depends on Task 1 (FormSection) + Tasks 2-3 (constants)
Phase 5 (Task 17): Permissions — depends on Task 3 (user constants)
Phase 6 (Task 18): Verification — after all phases
```

Phases 2, 3, and 4 can run in parallel after Phase 1 completes. Phase 5 can run in parallel with Phases 2-4.
