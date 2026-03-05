# Report Dashboard Redesign — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Transform the reports page into an enterprise-grade dashboard with unified header, KPI scorecards, semantic sections, adaptive charts, color palette, export, and polished UX.

**Architecture:** All changes are frontend-only within the existing `packages/observer-web/` package. We modify the BarChart component to support horizontal mode and color maps, add new i18n keys across 6 locales, create a color palette module, and restructure the reports page layout. No backend or database changes.

**Tech Stack:** React, d3, Tailwind CSS v4, @phosphor-icons/react, @base-ui, react-i18next, existing design tokens from `main.css`.

**Design doc:** `docs/plans/2026-03-05-report-dashboard-redesign.md`

---

### Task 1: Add New Icons to Icon Barrel

**Files:**
- Modify: `packages/observer-web/src/components/icons.ts`

**Step 1: Add icon exports**

Add these exports to the existing barrel file (alphabetical order within the list):

```ts
// Add to the existing export block:
CaretDownIcon,
CaretUpIcon,
DownloadSimpleIcon,
FunnelIcon,
PrinterIcon,
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bunx --bun vite build --mode development 2>&1 | head -5`
Expected: no errors about missing exports

**Step 3: Commit**

```bash
git add packages/observer-web/src/components/icons.ts
git commit -m "add icons for report dashboard redesign"
```

---

### Task 2: Add i18n Keys for All 6 Locales

**Files:**
- Modify: `packages/observer-web/src/locales/en.json`
- Modify: `packages/observer-web/src/locales/ky.json`
- Modify: `packages/observer-web/src/locales/ru.json`
- Modify: `packages/observer-web/src/locales/uk.json`
- Modify: `packages/observer-web/src/locales/de.json`
- Modify: `packages/observer-web/src/locales/tr.json`

**Step 1: Add keys to the `project.reports` section in each locale**

New keys to add (shown in English — translate for each locale):

```json
{
  "project": {
    "reports": {
      "exportCsv": "Export CSV",
      "print": "Print",
      "toggleFilters": "Filters",
      "presetThisMonth": "This month",
      "presetLastQuarter": "Last quarter",
      "presetThisYear": "This year",
      "presetAllTime": "All time",
      "clearAll": "Clear all",
      "sectionOverview": "Overview",
      "sectionServices": "Services",
      "sectionDemographics": "Demographics",
      "sectionGeography": "Geography & Taxonomy",
      "kpiPeople": "People",
      "kpiConsultations": "Consultations",
      "kpiActiveCases": "Active cases",
      "kpiIdp": "IDPs",
      "kpiHouseholds": "Households",
      "kpiOffices": "Offices"
    }
  }
}
```

For `ky.json`, use Kyrgyz Latin transliteration per project convention.

**Step 2: Verify JSON is valid**

Run: `for f in packages/observer-web/src/locales/*.json; do echo "$f"; python3 -c "import json; json.load(open('$f'))"; done`
Expected: no errors

**Step 3: Commit**

```bash
git add packages/observer-web/src/locales/*.json
git commit -m "add i18n keys for report dashboard redesign"
```

---

### Task 3: Create Color Palette Module

**Files:**
- Create: `packages/observer-web/src/components/charts/colors.ts`

**Step 1: Create the semantic color palette**

This module exports color maps for each data domain, plus a fallback palette and a helper to get a color for a given label.

```ts
export const SEX_COLORS: Record<string, string> = {
  female: "#e05a8a",
  male: "#5a8ae0",
  other: "#8b5cf6",
  unknown: "#94a3b8",
};

export const SUPPORT_TYPE_COLORS: Record<string, string> = {
  humanitarian: "#d97706",
  legal: "#3b82f6",
  social: "#10b981",
  psychological: "#8b5cf6",
  medical: "#ef4444",
  general: "#64748b",
};

export const SPHERE_COLORS: Record<string, string> = {
  housing_assistance: "#0ea5e9",
  document_recovery: "#6366f1",
  social_benefits: "#10b981",
  property_rights: "#f59e0b",
  employment_rights: "#ec4899",
  family_law: "#8b5cf6",
  healthcare_access: "#ef4444",
  education_access: "#14b8a6",
  financial_aid: "#d97706",
  psychological_support: "#a78bfa",
  other: "#94a3b8",
  unspecified: "#94a3b8",
};

export const CASE_STATUS_COLORS: Record<string, string> = {
  new: "#6366f1",
  active: "#10b981",
  closed: "#94a3b8",
  archived: "#64748b",
};

export const IDP_STATUS_COLORS: Record<string, string> = {
  idp: "#ef4444",
  non_idp: "#10b981",
  unknown: "#94a3b8",
};

export const AGE_GROUP_COLORS: Record<string, string> = {
  infant: "#fef3c7",
  toddler: "#fde68a",
  pre_school: "#fcd34d",
  middle_childhood: "#fbbf24",
  young_teen: "#f59e0b",
  teenager: "#d97706",
  young_adult: "#b45309",
  early_adult: "#92400e",
  middle_aged: "#78350f",
  older_adult: "#451a03",
};

export const FALLBACK_PALETTE = [
  "#6366f1", "#f59e0b", "#10b981", "#ef4444",
  "#8b5cf6", "#ec4899", "#14b8a6", "#f97316",
  "#3b82f6", "#84cc16", "#e879f9", "#06b6d4",
];

export function getColor(
  label: string,
  colorMap?: Record<string, string>,
  index?: number,
): string {
  if (colorMap?.[label]) return colorMap[label];
  return FALLBACK_PALETTE[(index ?? 0) % FALLBACK_PALETTE.length];
}
```

**Step 2: Commit**

```bash
git add packages/observer-web/src/components/charts/colors.ts
git commit -m "add semantic color palette for report charts"
```

---

### Task 4: Add Horizontal Mode to BarChart

**Files:**
- Modify: `packages/observer-web/src/components/charts/bar-chart.tsx`

**Step 1: Extend the BarChart props**

Add to `BarChartProps`:

```ts
interface BarChartProps {
  data: CountResult[];
  width?: number;
  height?: number;
  yAxisLabel?: string;
  legend?: BarLegendItem[];
  direction?: "vertical" | "horizontal" | "auto";
  colorMap?: Record<string, string>;
}
```

**Step 2: Implement horizontal rendering**

Inside the `BarChart` component, determine orientation:

```ts
const resolvedDirection =
  direction === "auto"
    ? data.length > 6 ? "horizontal" : "vertical"
    : (direction ?? "vertical");
```

When `resolvedDirection === "horizontal"`:
- Swap axes: labels on Y (scaleBand), values on X (scaleLinear)
- Increase left margin to ~120px to accommodate label text
- Auto-scale height: `Math.max(300, data.length * 32)` to give each bar breathing room
- Bars extend right from Y axis
- Labels are not rotated (they read naturally on Y axis)
- Value labels appear at the end of each bar

When `resolvedDirection === "vertical"`: existing logic unchanged.

For `colorMap`: use `getColor(d.label, colorMap, i)` from the colors module instead of the hardcoded `var(--color-accent)`:

```ts
import { getColor } from "./colors";
// ...
.attr("fill", (d, i) => getColor(d.label, colorMap, i))
```

**Step 3: Verify build**

Run: `cd packages/observer-web && bunx --bun vite build --mode development 2>&1 | head -5`

**Step 4: Commit**

```bash
git add packages/observer-web/src/components/charts/bar-chart.tsx
git commit -m "add horizontal mode and color map support to BarChart"
```

---

### Task 5: Update PieChart to Use Color Maps

**Files:**
- Modify: `packages/observer-web/src/components/charts/pie-chart.tsx`

**Step 1: Accept colorMap prop**

Add `colorMap?: Record<string, string>` to `PieChart` props.

**Step 2: Replace hardcoded COLORS**

Instead of the static `COLORS` array with `d3.scaleOrdinal`, use:

```ts
import { getColor, FALLBACK_PALETTE } from "./colors";
// ...
const color = (label: string, i: number) => getColor(label, colorMap, i);
```

Update the `.attr("fill", ...)` in slices and the legend dot `backgroundColor` to use this function.

Keep the existing `COLORS` array as `FALLBACK_PALETTE` is now the source of truth — delete the local `COLORS` constant.

**Step 3: Verify build**

Run: `cd packages/observer-web && bunx --bun vite build --mode development 2>&1 | head -5`

**Step 4: Commit**

```bash
git add packages/observer-web/src/components/charts/pie-chart.tsx
git commit -m "use semantic color maps in PieChart"
```

---

### Task 6: Create CSV Export Utility

**Files:**
- Create: `packages/observer-web/src/lib/export-csv.ts`

**Step 1: Implement CSV export**

```ts
import type { FullReport } from "@/types/report";

function escapeCSV(value: string): string {
  if (value.includes(",") || value.includes('"') || value.includes("\n")) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}

export function exportReportCSV(data: FullReport, projectId: string) {
  const rows: string[] = ["Group,Label,Count"];

  const groups: [string, { rows: { label: string; count: number }[]; total: number }][] = [
    ["Consultations", data.consultations],
    ["By Sex", data.by_sex],
    ["By IDP Status", data.by_idp_status],
    ["By Category", data.by_category],
    ["By Region", data.by_region],
    ["By Sphere", data.by_sphere],
    ["By Office", data.by_office],
    ["By Age Group", data.by_age_group],
    ["By Tag", data.by_tag],
    ["Family Units", data.family_units],
  ];

  for (const [name, group] of groups) {
    for (const row of group.rows) {
      rows.push(`${escapeCSV(name)},${escapeCSV(row.label)},${row.count}`);
    }
    rows.push(`${escapeCSV(name)},TOTAL,${group.total}`);
  }

  const blob = new Blob([rows.join("\n")], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const date = new Date().toISOString().slice(0, 10);
  const a = document.createElement("a");
  a.href = url;
  a.download = `report-${projectId}-${date}.csv`;
  a.click();
  URL.revokeObjectURL(url);
}
```

**Step 2: Commit**

```bash
git add packages/observer-web/src/lib/export-csv.ts
git commit -m "add client-side CSV export for reports"
```

---

### Task 7: Add Print Stylesheet

**Files:**
- Modify: `packages/observer-web/src/main.css`

**Step 1: Add @media print rules at the end of main.css**

```css
@media print {
  /* Hide non-essential UI */
  nav,
  aside,
  [data-print-hide] {
    display: none !important;
  }

  /* Reset backgrounds for ink savings */
  body,
  :root,
  [data-theme] {
    background: #fff !important;
    color: #000 !important;
  }

  /* Single-column chart layout */
  .report-grid {
    display: block !important;
  }
  .report-grid > * {
    break-inside: avoid;
    margin-bottom: 1rem;
  }

  /* Ensure charts are visible */
  svg {
    max-width: 100% !important;
    height: auto !important;
  }

  /* Print header */
  .print-header {
    display: block !important;
  }
}
```

**Step 2: Commit**

```bash
git add packages/observer-web/src/main.css
git commit -m "add print stylesheet for report dashboard"
```

---

### Task 8: Rewrite Reports Page — Unified Header + Collapsible Filters

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/reports/index.tsx`

This is the largest task. We rewrite the `ReportsPage` component with the new layout. All the building blocks from Tasks 1–7 are now available.

**Step 1: Add imports for new dependencies**

Add to the top of the file:

```ts
import {
  CaretDownIcon,
  CaretUpIcon,
  DownloadSimpleIcon,
  FunnelIcon,
  PrinterIcon,
  XIcon,
} from "@/components/icons";
import {
  SEX_COLORS,
  SUPPORT_TYPE_COLORS,
  SPHERE_COLORS,
  CASE_STATUS_COLORS,
  IDP_STATUS_COLORS,
  AGE_GROUP_COLORS,
  FALLBACK_PALETTE,
} from "@/components/charts/colors";
import { exportReportCSV } from "@/lib/export-csv";
```

**Step 2: Add date preset helper**

Above `ReportsPage`, add:

```ts
type DatePreset = "month" | "quarter" | "year" | "all";

function getPresetDates(preset: DatePreset): { date_from?: string; date_to?: string } {
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
```

**Step 3: Add SectionHeader component**

```tsx
function SectionHeader({ title }: { title: string }) {
  return (
    <div className="col-span-full border-b border-border-secondary pb-1 pt-4">
      <h2 className="text-xs font-semibold uppercase tracking-wider text-fg-tertiary">{title}</h2>
    </div>
  );
}
```

**Step 4: Add KpiCard component**

```tsx
function KpiCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-4">
      <p className="text-2xl font-bold tabular-nums text-fg">{value.toLocaleString()}</p>
      <p className="mt-0.5 text-xs font-medium text-fg-tertiary">{label}</p>
    </div>
  );
}
```

**Step 5: Add FilterChip component**

```tsx
function FilterChip({
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

**Step 6: Add ReportSkeleton component**

```tsx
function ReportSkeleton() {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-3 gap-4 lg:grid-cols-6">
        {Array.from({ length: 6 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
      <div className="h-64 animate-pulse rounded-xl bg-bg-tertiary" />
      <div className="grid gap-6 lg:grid-cols-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-72 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
    </div>
  );
}
```

**Step 7: Rewrite the ReportsPage component**

Update the `ReportsPage` function. Keep all existing state/hooks, add:

```ts
const [filtersOpen, setFiltersOpen] = useState(true);
const [activePreset, setActivePreset] = useState<DatePreset | null>(null);
```

The JSX structure becomes:

```tsx
<div>
  {/* Print-only header */}
  <div className="print-header hidden">
    <h1 className="text-lg font-bold">{t("project.reports.title")}</h1>
    {params.date_from && <p>From: {params.date_from} To: {params.date_to ?? "now"}</p>}
  </div>

  {/* Unified header + filter panel */}
  <div data-print-hide className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary">
    {/* Top bar: title + actions + filter toggle */}
    <div className="flex items-center justify-between px-5 py-3">
      <h1 className="font-serif text-xl font-bold tracking-tight text-fg">
        {t("project.reports.title")}
      </h1>
      <div className="flex items-center gap-2">
        {data && (
          <>
            <button
              type="button"
              onClick={() => exportReportCSV(data, projectId)}
              className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
            >
              <DownloadSimpleIcon size={14} />
              {t("project.reports.exportCsv")}
            </button>
            <button
              type="button"
              onClick={() => window.print()}
              className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
            >
              <PrinterIcon size={14} />
              {t("project.reports.print")}
            </button>
          </>
        )}
        <button
          type="button"
          onClick={() => setFiltersOpen((o) => !o)}
          className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
        >
          <FunnelIcon size={14} />
          {t("project.reports.toggleFilters")}
          {filtersOpen ? <CaretUpIcon size={12} /> : <CaretDownIcon size={12} />}
        </button>
      </div>
    </div>

    {/* Collapsible filter panel */}
    {filtersOpen && (
      <div className="border-t border-border-secondary px-5 pb-4 pt-3">
        {/* Date presets */}
        <div className="mb-3 flex flex-wrap gap-1.5">
          {(["month", "quarter", "year", "all"] as const).map((preset) => (
            <button
              key={preset}
              type="button"
              onClick={() => {
                const dates = getPresetDates(preset);
                setParams((p) => ({ ...p, ...dates }));
                setActivePreset(preset);
              }}
              className={`rounded-md px-2.5 py-1 text-xs font-medium transition-colors ${
                activePreset === preset
                  ? "bg-accent text-accent-fg"
                  : "bg-bg-tertiary text-fg-secondary hover:text-fg"
              }`}
            >
              {t(`project.reports.preset${preset[0].toUpperCase()}${preset.slice(1)}`)}
            </button>
          ))}
        </div>

        {/* Filter grid (same as existing, reusing FilterField + UISelect) */}
        <div className="grid grid-cols-2 gap-x-4 gap-y-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-8">
          {/* ... existing filter fields, keep them identical ... */}
          {/* But update DatePicker onChange to clear activePreset: */}
          {/* onChange={(v) => { setParams(p => ({ ...p, date_from: v || undefined })); setActivePreset(null); }} */}
        </div>
      </div>
    )}

    {/* Active filter chips (always visible even when collapsed) */}
    {hasFilters && (
      <div className="flex flex-wrap items-center gap-1.5 border-t border-border-secondary px-5 py-2.5">
        {params.office_id && (
          <FilterChip
            label={t("project.reports.filterOffice")}
            value={officeOptions.find((o) => o.value === params.office_id)?.label ?? params.office_id}
            onRemove={() => setParams((p) => ({ ...p, office_id: undefined }))}
          />
        )}
        {params.category_id && (
          <FilterChip
            label={t("project.reports.filterCategory")}
            value={categoryOptions.find((c) => c.value === params.category_id)?.label ?? params.category_id}
            onRemove={() => setParams((p) => ({ ...p, category_id: undefined }))}
          />
        )}
        {params.consultant_id && (
          <FilterChip
            label={t("project.reports.filterConsultant")}
            value={consultantOptions.find((c) => c.value === params.consultant_id)?.label ?? params.consultant_id}
            onRemove={() => setParams((p) => ({ ...p, consultant_id: undefined }))}
          />
        )}
        {params.case_status && (
          <FilterChip
            label={t("project.reports.filterCaseStatus")}
            value={caseStatusOptions.find((s) => s.value === params.case_status)?.label ?? params.case_status}
            onRemove={() => setParams((p) => ({ ...p, case_status: undefined }))}
          />
        )}
        {params.sex && (
          <FilterChip
            label={t("project.reports.filterSex")}
            value={sexOptions.find((s) => s.value === params.sex)?.label ?? params.sex}
            onRemove={() => setParams((p) => ({ ...p, sex: undefined }))}
          />
        )}
        {params.age_group && (
          <FilterChip
            label={t("project.reports.filterAgeGroup")}
            value={ageGroupOptions.find((g) => g.value === params.age_group)?.label ?? params.age_group}
            onRemove={() => setParams((p) => ({ ...p, age_group: undefined }))}
          />
        )}
        {params.date_from && (
          <FilterChip
            label={t("project.reports.dateFrom")}
            value={params.date_from}
            onRemove={() => { setParams((p) => ({ ...p, date_from: undefined })); setActivePreset(null); }}
          />
        )}
        {params.date_to && (
          <FilterChip
            label={t("project.reports.dateTo")}
            value={params.date_to}
            onRemove={() => { setParams((p) => ({ ...p, date_to: undefined })); setActivePreset(null); }}
          />
        )}
        <button
          type="button"
          onClick={() => { setParams({}); setActivePreset(null); }}
          className="ml-1 text-xs font-medium text-fg-tertiary underline transition-colors hover:text-fg"
        >
          {t("project.reports.clearAll")}
        </button>
      </div>
    )}
  </div>

  {/* Loading skeleton */}
  {isLoading && <ReportSkeleton />}

  {/* Dashboard content */}
  {data && (
    <div className="report-grid grid gap-6 lg:grid-cols-2">
      {/* KPI Row */}
      <div className="col-span-full grid grid-cols-3 gap-4 lg:grid-cols-6">
        <KpiCard label={t("project.reports.kpiPeople")} value={data.by_sex.total} />
        <KpiCard label={t("project.reports.kpiConsultations")} value={data.consultations.total} />
        <KpiCard
          label={t("project.reports.kpiActiveCases")}
          value={data.by_idp_status.rows.find((r) => r.label === "active")?.count ?? 0}
        />
        <KpiCard
          label={t("project.reports.kpiIdp")}
          value={data.by_idp_status.rows.find((r) => r.label === "idp")?.count ?? data.by_idp_status.total}
        />
        <KpiCard label={t("project.reports.kpiHouseholds")} value={data.family_units.total} />
        <KpiCard label={t("project.reports.kpiOffices")} value={data.by_office.rows.length} />
      </div>

      {/* Overview section */}
      <SectionHeader title={t("project.reports.sectionOverview")} />
      {data.status_flow && data.status_flow.length > 0 && (
        <div className="col-span-full rounded-xl border border-border-secondary bg-bg-secondary p-5">
          <h3 className="mb-3 text-sm font-semibold text-fg">
            {t("project.reports.statusFlow")}
          </h3>
          <SankeyChart
            data={data.status_flow}
            translateLabel={(l) => {
              const key = labelKeyMap[l];
              return key ? t(key) : t(`project.people.${l}`, l);
            }}
          />
        </div>
      )}

      {/* Services section */}
      <SectionHeader title={t("project.reports.sectionServices")} />
      <div className="col-span-full">
        <ReportCard
          group={data.consultations}
          title={t("project.reports.consultations")}
          chart="bar"
          yAxisLabel={axisLabel}
          colorMap={SUPPORT_TYPE_COLORS}
        />
      </div>
      <ReportCard
        group={data.by_sphere}
        title={t("project.reports.bySphere")}
        chart="bar"
        yAxisLabel={axisLabel}
        colorMap={SPHERE_COLORS}
        direction="auto"
      />
      <ReportCard
        group={data.by_office}
        title={t("project.reports.byOffice")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />

      {/* Demographics section */}
      <SectionHeader title={t("project.reports.sectionDemographics")} />
      <div className="col-span-full grid grid-cols-1 gap-6 md:grid-cols-3">
        <ReportCard
          group={data.by_sex}
          title={t("project.reports.bySex")}
          chart="pie"
          colorMap={SEX_COLORS}
        />
        <ReportCard
          group={data.family_units}
          title={t("project.reports.familyUnits")}
          chart="pie"
        />
        <ReportCard
          group={data.by_idp_status}
          title={t("project.reports.byIdpStatus")}
          chart="pie"
          colorMap={IDP_STATUS_COLORS}
        />
      </div>
      <div className="col-span-full">
        <ReportCard
          group={data.by_age_group}
          title={t("project.reports.byAgeGroup")}
          chart="bar"
          yAxisLabel={axisLabel}
          skipTranslation
          mapLabel={(l) => AGE_RANGE_MAP[l] ?? l}
          legend={ageGroupLegend}
          colorMap={AGE_GROUP_COLORS}
        />
      </div>

      {/* Geography & Taxonomy section */}
      <SectionHeader title={t("project.reports.sectionGeography")} />
      <ReportCard
        group={data.by_region}
        title={t("project.reports.byRegion")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />
      <ReportCard
        group={data.by_category}
        title={t("project.reports.byCategory")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />
      <div className="col-span-full">
        <ReportCard
          group={data.by_tag}
          title={t("project.reports.byTag")}
          chart="bar"
          yAxisLabel={axisLabel}
          direction="auto"
        />
      </div>
    </div>
  )}
</div>
```

**Step 8: Update ReportCard to pass through new props**

The `ReportCard` component needs to accept and forward `colorMap` and `direction`:

```tsx
function ReportCard({
  group,
  title,
  chart,
  yAxisLabel,
  legend,
  mapLabel,
  skipTranslation,
  colorMap,
  direction,
}: {
  group: ReportGroup;
  title: string;
  chart: "bar" | "pie";
  yAxisLabel?: string;
  legend?: BarLegendItem[];
  mapLabel?: (label: string) => string;
  skipTranslation?: boolean;
  colorMap?: Record<string, string>;
  direction?: "vertical" | "horizontal" | "auto";
}) {
  const translated = useTranslatedRows(group.rows);
  const source = skipTranslation ? group.rows : translated;
  const rows = mapLabel ? source.map((r) => ({ ...r, label: mapLabel(r.label) })) : source;
  return (
    <div className="min-h-[280px] rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-3 flex items-baseline justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <span className="tabular-nums text-xs font-medium text-fg-tertiary">
          {group.total.toLocaleString()}
        </span>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart data={rows} yAxisLabel={yAxisLabel} legend={legend} colorMap={colorMap} direction={direction} />
        ) : (
          <PieChart data={rows} colorMap={colorMap} />
        )
      ) : (
        <p className="py-8 text-center text-sm text-fg-tertiary">—</p>
      )}
    </div>
  );
}
```

Note `min-h-[280px]` added for consistent card heights.

**Step 9: Verify the app builds and renders**

Run: `cd packages/observer-web && bunx --bun vite build --mode development 2>&1 | tail -3`
Expected: build succeeds

**Step 10: Commit**

```bash
git add packages/observer-web/src/routes/_app/projects/\$projectId/reports/index.tsx
git commit -m "redesign report dashboard with unified header, KPI row, sections, and adaptive charts"
```

---

### Task 9: Translate i18n Keys for preset labels

**Files:**
- Modify: `packages/observer-web/src/locales/ky.json`
- Modify: `packages/observer-web/src/locales/ru.json`
- Modify: `packages/observer-web/src/locales/uk.json`
- Modify: `packages/observer-web/src/locales/de.json`
- Modify: `packages/observer-web/src/locales/tr.json`

The preset key names in code are:
- `presetMonth` → maps to `presetThisMonth`
- `presetQuarter` → maps to `presetLastQuarter`
- `presetYear` → maps to `presetThisYear`
- `presetAll` → maps to `presetAllTime`

Wait — check the preset button code. The key construction is:
```ts
t(`project.reports.preset${preset[0].toUpperCase()}${preset.slice(1)}`)
```
With presets `"month" | "quarter" | "year" | "all"` this produces:
- `project.reports.presetMonth`
- `project.reports.presetQuarter`
- `project.reports.presetYear`
- `project.reports.presetAll`

So the i18n keys need to match exactly. Update Task 2's keys to use:
- `"presetMonth"` (not `"presetThisMonth"`)
- `"presetQuarter"` (not `"presetLastQuarter"`)
- `"presetYear"` (not `"presetThisYear"`)
- `"presetAll"` (not `"presetAllTime"`)

With display text:
- en: "This month", "Last quarter", "This year", "All time"
- ru: "Этот месяц", "Прошлый квартал", "Этот год", "Всё время"
- uk: "Цей місяць", "Минулий квартал", "Цей рік", "Весь час"
- de: "Dieser Monat", "Letztes Quartal", "Dieses Jahr", "Gesamter Zeitraum"
- tr: "Bu ay", "Geçen çeyrek", "Bu yıl", "Tüm zamanlar"
- ky: Use Kyrgyz Latin transliteration

This task is just a correction pass to ensure key names match the code. Should be handled as part of Task 2.

**Step 1: Verify all keys match code usage**

Search the reports page for all `t("project.reports.` calls and cross-check with locale files.

**Step 2: Commit if any fixes needed**

```bash
git add packages/observer-web/src/locales/*.json
git commit -m "fix i18n key names to match preset code"
```

---

### Task 10: Visual QA and Polish

**Files:**
- Possibly touch: any of the above files for minor adjustments

**Step 1: Run the dev server**

Run: `cd packages/observer-web && bun run dev`

**Step 2: Visual check in browser**

Open `http://localhost:5173` (or whatever port), navigate to a project's Reports page. Verify:

- [ ] Unified header with title, export buttons, filter toggle
- [ ] Filter collapse/expand works
- [ ] Date presets set dates and highlight active pill
- [ ] Active filter chips appear and are individually removable
- [ ] KPI cards show correct totals in a 6-column row
- [ ] Sankey appears first under "Overview" section
- [ ] Consultations is full-width under "Services"
- [ ] 3 donuts side by side under "Demographics"
- [ ] Charts with >6 items render horizontal bars
- [ ] Colors are semantic per domain
- [ ] Loading shows skeleton cards
- [ ] CSV export downloads a valid file
- [ ] Print button opens print dialog with clean layout
- [ ] All cards have consistent min-height

**Step 3: Fix any visual issues found**

Adjust spacing, colors, or layout as needed.

**Step 4: Final commit**

```bash
git add -A
git commit -m "polish report dashboard visual QA fixes"
```

---

## Task Dependency Graph

```
Task 1 (icons) ──┐
Task 2 (i18n) ───┤
Task 3 (colors) ─┼── Task 8 (page rewrite) ── Task 9 (key fix) ── Task 10 (QA)
Task 4 (bar-h) ──┤
Task 5 (pie) ────┤
Task 6 (csv) ────┤
Task 7 (print) ──┘
```

Tasks 1–7 are independent and can be parallelized. Task 8 depends on all of them. Tasks 9–10 are sequential after 8.
