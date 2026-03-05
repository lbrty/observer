# DatePicker Component Design

**Goal:** Replace all native `<input type="date">` with reusable DatePicker components using react-day-picker v9, styled to match the existing dark theme.

## Components

### DatePicker (single date)

```tsx
<DatePicker
  value="2026-03-05"           // YYYY-MM-DD string
  onChange={(date) => ...}      // receives YYYY-MM-DD string or ""
  placeholder="dd.mm.yyyy"     // optional, defaults to "dd.mm.yyyy"
  disabled={false}             // optional
  className="..."              // optional, applied to trigger wrapper
/>
```

- Trigger: read-only text input + calendar icon, using existing inputClass pattern
- Popover: Base UI Popover wrapping react-day-picker `<DayPicker mode="single" />`
- Display format: DD.MM.YYYY universal, internal value stays YYYY-MM-DD
- Clicking a day selects it, closes the popover, calls onChange
- Clicking the selected day again deselects it (for optional fields)

### DateRangePicker (reports page)

```tsx
<DateRangePicker
  from="2026-03-01"
  to="2026-03-05"
  onChange={({ from, to }) => ...}
/>
```

- Two triggers side-by-side ("From" / "To") opening a single shared calendar popover
- Uses react-day-picker `mode="range"` selection
- Selecting start and end date closes popover and calls onChange with both values

## Architecture

- `packages/observer-web/src/components/date-picker.tsx` — both components
- Calendar styled via CSS overrides using existing CSS variables (--bg, --bg-secondary, --accent, --fg, etc.)
- No changes to API payloads — values stay YYYY-MM-DD strings

## Replacement Targets (7 inputs, 5 files)

1. `reports/index.tsx` — `date_from`, `date_to` (replace with DateRangePicker)
2. `migration-records.tsx` — `migration_date`
3. `support-record-drawer.tsx` — `provided_at`
4. `people/$personId/index.tsx` — `providedAt`
5. `person-drawer.tsx` — `birth_date`, `consent_date`

## Dependencies

- `react-day-picker` v9
