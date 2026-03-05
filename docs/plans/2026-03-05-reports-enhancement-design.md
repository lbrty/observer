# Reports Enhancement Design

**Goal:** Extend the reports page with global filtering, chart interactivity, translated labels, and a Sankey diagram for case status flow.

## 1. Global Filter Bar

Six dropdowns added below date range pickers using existing `UISelect`:
- **Office** — from `useOffices` hook
- **Category** — from `useCategories` hook (via admin reference)
- **Consultant** — from project members or `useUsers`
- **Case status** — static: new, active, closed, archived
- **Sex** — static: male, female, other, unknown
- **Age group** — static: 10 groups

All sent as query params to backend. Backend `ReportFilter` struct extended with 6 optional fields. Report queries get additional WHERE clauses. Some queries need conditional JOINs (e.g. filtering by category requires `JOIN person_categories`).

"Clear filters" button resets all filters including date range.

## 2. Axis Titles

- Y-axis: localized "Count" label (`project.reports.axisCount`)
- X-axis: no title (chart heading names the dimension)

## 3. Tooltips

- Bar chart: hover shows tooltip with label + count, themed with `bg-bg-secondary border-border-secondary shadow-elevated`
- Pie/donut chart: same tooltip near cursor
- Implemented with absolute-positioned div, no external library

## 4. Click-to-Highlight

- Click a bar/slice: full opacity on selected, others fade to 30%
- Click again or click elsewhere to deselect
- Frontend-only, no data re-fetch

## 5. Sankey Diagram — Case Status Flow

### Database

New migration creates:

```sql
CREATE TABLE person_status_history (
  id         TEXT PRIMARY KEY DEFAULT gen_random_ulid(),
  person_id  TEXT NOT NULL REFERENCES people(id) ON DELETE CASCADE,
  from_status TEXT NOT NULL,
  to_status   TEXT NOT NULL,
  changed_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_psh_person ON person_status_history(person_id);
CREATE INDEX idx_psh_changed ON person_status_history(changed_at);
```

Postgres trigger on `people` table — fires on UPDATE of `case_status`, inserts history row automatically:

```sql
CREATE FUNCTION track_person_status_change() RETURNS TRIGGER AS $$
BEGIN
  IF OLD.case_status IS DISTINCT FROM NEW.case_status THEN
    INSERT INTO person_status_history(person_id, from_status, to_status)
    VALUES (NEW.id, OLD.case_status, NEW.case_status);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_person_status_change
  AFTER UPDATE ON people
  FOR EACH ROW EXECUTE FUNCTION track_person_status_change();
```

### Seed Command

Updated to generate realistic transition histories:
- new → active: 1-14 days after registration
- active → closed: 7-90 days after activation
- Some closed → archived: 30-180 days later
- Uses person's `registered_at` as baseline

### Backend

New report query: aggregate transitions + mean days between each status pair.

```sql
SELECT from_status, to_status,
       COUNT(*) AS count,
       AVG(EXTRACT(EPOCH FROM (changed_at - lag_at)) / 86400)::numeric(10,1) AS avg_days
FROM (
  SELECT from_status, to_status, changed_at,
         LAG(changed_at) OVER (PARTITION BY person_id ORDER BY changed_at) AS lag_at
  FROM person_status_history
  WHERE person_id IN (SELECT id FROM people WHERE project_id = $1)
) sub
WHERE lag_at IS NOT NULL OR from_status = 'new'
GROUP BY from_status, to_status
```

Returns `[]StatusFlow{FromStatus, ToStatus, Count, AvgDays}`.

### Frontend

- D3 Sankey layout via `d3-sankey`
- Nodes: new, active, closed, archived (4 statuses)
- Links: thickness = transition count, label = "avg X days"
- Full-width card at bottom of reports grid
- Theme-aware colors per status node

## 6. Donut Chart Fix

Current issue: fixed `w-48` (192px) SVG with 300x300 viewBox wastes card space.

Fix:
- Remove fixed `w-48`, use responsive width
- Center donut + legend using flexbox
- Increase donut size to fill available space

## 7. i18n

New keys added to all 6 locales (en, ky, ru, uk, de, tr):
- `project.reports.axisCount`
- `project.reports.filterOffice`
- `project.reports.filterCategory`
- `project.reports.filterConsultant`
- `project.reports.filterCaseStatus`
- `project.reports.filterSex`
- `project.reports.filterAgeGroup`
- `project.reports.statusFlow`
- `project.reports.avgDays`

## Dependencies

- `d3-sankey` + `@types/d3-sankey` for Sankey diagram
