# Frontend Component Extraction Design

Date: 2026-03-08

## Goal

Reduce cognitive load for newcomers by breaking large files into small, self-documenting pieces. No file should require scrolling to understand.

## Principles

- Readability and discoverability over clever abstractions
- Small files that are self-documenting
- A newcomer reads the index to understand the flow, dives into a section for details

## 1. FormSection Compound

Generic layout component replacing the repeated `SectionHeading + grid` pattern across all 5 drawers.

```tsx
<FormSection title="Identyfikatsiya" columns={2}>
  <FormField label="Aty" .../>
  <FormField label="Jöny" .../>
</FormSection>
```

- Renders SectionHeading + responsive grid (grid-cols-1 sm:grid-cols-{columns}, gap-4)
- Optional `columns` prop (default 2)

## 2. Domain Constants

Extract repeated option/label maps into per-domain constant files:

```
web/src/constants/
  person.ts      # sex, age-group, case-status options + labels
  support.ts     # support-type, support-sphere, referral-status
  migration.ts   # movement-reason, housing-at-destination
  pet.ts         # pet-status
  household.ts   # relationship types
  user.ts        # roles, user statuses
```

Each exports typed `{ value, label }` arrays consumed by selects, filters, badges, and reports.

## 3. DataTablePage Compound

Single compound for the full list-page pattern (15+ pages currently repeat this):

```tsx
<DataTablePage
  title="Adamdar"
  columns={columns}
  data={people}
  isLoading={isLoading}
  pagination={{ page, totalPages, onPageChange }}
  emptyState={{ icon: Users, title: "Adamdar jok" }}
  createAction={{ label: "Jany adam", onClick: openDrawer }}
  filters={[
    { type: "search", placeholder: "Izdöö...", value, onChange },
    { type: "select", label: "Status", options: caseStatusOptions, value, onChange },
  ]}
/>
```

Renders: PageHeader + FilterBar + DataTable + Pagination + EmptyState. Pages become ~50-80 lines.

## 4. Drawer Refactoring

Each drawer split into folder with index + section files:

```
components/
  person-drawer/
    index.tsx              # ~80 lines — hooks, submit, composes sections
    identity-section.tsx   # ~50 lines — name, sex, birth, age
    location-section.tsx   # ~40 lines — origin/current place pickers
    case-section.tsx       # ~50 lines — status, external ID, office, consent
  support-record-drawer/
    index.tsx              # ~70 lines
    info-section.tsx       # ~50 lines — type, sphere, date, person
    referral-section.tsx   # ~40 lines — referral status, office
  household-drawer/
    index.tsx              # ~60 lines
    head-section.tsx       # ~40 lines — reference number, head picker
    members-section.tsx    # ~80 lines — member list + add form
  migration-record-drawer/
    index.tsx              # ~60 lines
    origin-section.tsx     # ~50 lines — origin place selectors
    destination-section.tsx # ~50 lines — destination place selectors
    details-section.tsx    # ~30 lines — reason, housing
  pet-drawer/
    index.tsx              # ~60 lines
    info-section.tsx       # ~40 lines — name, species, breed
    status-section.tsx     # ~30 lines — status, notes
```

Sections receive form state + handlers as props. All use `<FormSection>` for layout.

## 5. Report Page Extraction

Shared report components:

```
components/
  report/
    report-filters.tsx     # ~60 lines — date range, category, office selects
    kpi-grid.tsx           # ~40 lines — grid of KPI stat cards
    chart-card.tsx         # ~30 lines — card wrapper for a chart with title
    report-skeleton.tsx    # ~30 lines — loading skeleton for reports
```

Page-specific extractions:

```
routes/projects/$projectId/
  reports/
    people.tsx             # ~100 lines — composes filters + KPIs + charts
    people/
      label-maps.ts        # ~80 lines — report type labels + chart configs
      people-kpis.tsx      # ~40 lines — KPI definitions
      people-charts.tsx    # ~60 lines — chart layout
    pets.tsx               # ~80 lines — composes shared report components
    pets/
      label-maps.ts
      pets-charts.tsx
  my-stats/
    index.tsx              # ~80 lines — composes filters + cards + charts
    stats-cards.tsx        # ~50 lines
    stats-charts.tsx       # ~60 lines
```

## 6. Permissions Page Extraction

```
routes/admin/projects/$projectId/
  permissions.tsx          # ~80 lines — data fetching, composes matrix
  permissions/
    permission-matrix.tsx  # ~120 lines — table layout, header row
    permission-row.tsx     # ~60 lines — single user row with toggle cells
    role-select.tsx        # ~40 lines — project role dropdown per user
```

## Expected Impact

- Largest file: 740 → ~120 lines
- Drawers: 350-420 → ~60-80 lines (index) + small sections
- List pages: 200-330 → ~50-80 lines
- No new abstractions beyond splitting and composing
