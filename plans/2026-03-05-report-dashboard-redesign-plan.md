# Report Dashboard Redesign — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Transform the reports page into an enterprise-grade dashboard with unified header, KPI scorecards, semantic sections, adaptive charts, color palette, export, and polished UX.

**Architecture:** All changes are frontend-only within the existing `packages/observer-web/` package. We modify the BarChart component to support horizontal mode and color maps, add new i18n keys across 6 locales, create a color palette module, and restructure the reports page layout. No backend or database changes.

**Tech Stack:** React, d3, Tailwind CSS v4, @phosphor-icons/react, @base-ui, react-i18next, existing design tokens from `main.css`.

**Design doc:** `docs/plans/2026-03-05-report-dashboard-redesign.md`


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

