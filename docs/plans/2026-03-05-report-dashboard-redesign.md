# Report Dashboard Redesign

**Goal:** Transform the reports page from a flat grid of charts into an enterprise-grade dashboard with unified header, KPI scorecards, semantic sections, adaptive chart layouts, color palette, export, and polished UX.

## Layout (top to bottom)

```
+--------------------------------------------------+
| Title + [CSV] [Print]              [^ Filters]   |
| +-- Filter bar (expanded by default, collapsible) |
| | [preset pills: Month | Quarter | Year | All]   |
| | [date from] [date to] [office] [category]       |
| | [consultant] [status] [sex] [age group]         |
| | [chip: Office: X] [chip: Sex: Y]   [Clear all]  |
| +------------------------------------------------+|
+--------------------------------------------------+
| KPI Row (6 cards)                                 |
| [People] [Consultations] [Active] [IDP] [HH] [O] |
+--------------------------------------------------+
| -- Overview --                                    |
| [Status Flow Sankey -- full width]                |
+--------------------------------------------------+
| -- Services --                                    |
| [Consultations bar -- full width]                 |
| [By Sphere bar] [By Office bar]                   |
+--------------------------------------------------+
| -- Demographics --                                |
| [By Sex donut] [Family Units donut] [IDP donut]   |
| [By Age Group bar -- full width]                  |
+--------------------------------------------------+
| -- Geography & Taxonomy --                        |
| [By Region bar] [By Category bar]                 |
| [By Tag bar -- full width]                        |
+--------------------------------------------------+
```

## 1. Unified Header + Filter Bar

- Page title, export buttons (CSV + Print), and collapse toggle in one row
- Filter panel is inside the same container as the header (no separate bordered box)
- Collapse toggle hides the filter grid; active filter chips remain visible when collapsed
- Expanded by default; user can collapse manually

## 2. Date Range Presets

- Inline pill buttons above the date pickers: "This Month", "Last Quarter", "This Year", "All Time"
- Clicking a preset sets both `date_from` and `date_to` and highlights the active pill
- Manual date entry deselects any active preset

## 3. Active Filter Chips

- When filters are applied, render removable chips below the filter grid
- Each chip shows "Label: Value x" and removes that single filter on click
- "Clear all" button resets everything

## 4. KPI Summary Row

- 6 stat cards derived from existing `group.total` values:
  - Total People (`by_sex.total`)
  - Consultations (`consultations.total`)
  - Active Cases (sum of `by_idp_status` rows with active-like labels)
  - IDP Count (`by_idp_status` IDP row count)
  - Households (`family_units.total`)
  - Offices (`by_office` distinct count)
- Large number, small label underneath, no mini-charts
- Responsive: 3x2 on mobile, 6x1 on desktop

## 5. Section Headers

- Lightweight: text label + thin `border-b` divider
- Sections: Overview, Services, Demographics, Geography & Taxonomy
- No cards/boxes around sections

## 6. Chart Layout Changes

- Sankey promoted to first chart (full-width, under KPI row)
- Consultations bar chart: full-width
- Donuts grouped: `by_sex`, `family_units`, `by_idp_status` in a 3-column row
- Age Group bar: full-width
- Tag bar: full-width
- Remaining bars: 2-column grid

## 7. Adaptive Horizontal Bars

- When a bar chart has >6 data points, render as horizontal bar chart
- Labels on y-axis (read naturally), bars extend right
- Applies to: `by_category`, `by_region`, `by_tag`, `by_office`, `by_sphere`
- BarChart component gets a `direction` prop: `"vertical" | "horizontal" | "auto"` (default `"auto"`)

## 8. Category-Aware Color Palette

Semantic color assignments per data domain, consistent across bar and donut charts:

- **Sex**: distinct hue per value (female=rose, male=blue, other=purple, unknown=gray)
- **Support types**: warm palette (humanitarian=amber, legal=blue, social=green, psychological=violet, medical=red, general=slate)
- **Spheres**: cool/neutral gradient
- **Case statuses**: semantic (new=blue, active=green, closed=gray, archived=slate)
- **Age groups**: sequential warm ramp (light to dark)
- **IDP status**: traffic-light inspired
- **Categories/tags/regions**: generated from a fixed 12-color palette, assigned by index

ReportCard accepts optional `colorMap: Record<string, string>` prop. Fallback: current single-color bars.

## 9. Consistent Card Heights

- Chart containers get `min-h-[280px]` (or similar) within each grid row
- Prevents jagged grid when donut cards are shorter than bar cards

## 10. Loading Skeletons

- Replace text "Loading..." with skeleton cards matching grid layout
- Pulsing rectangles: KPI row (6 small), then chart card shapes
- Use `animate-pulse bg-bg-tertiary rounded-xl` pattern

## 11. Export

### CSV
- Client-side: iterate all report groups, build CSV with group name, label, count columns
- Trigger via download link with `Blob` URL
- Filename: `report-{projectId}-{date}.csv`

### Print
- "Print" button triggers `window.print()`
- `@media print` stylesheet:
  - Hide sidebar, header, filters, export buttons
  - Single-column layout
  - `page-break-inside: avoid` on chart cards
  - Light backgrounds for ink savings
  - Add project name + date range in print header

## Dependencies

No new packages required. All changes are frontend-only.
