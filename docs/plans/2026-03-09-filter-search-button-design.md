# Filter Search Button

## Problem

Date range and other filter changes apply instantly via `onChange`. Users sometimes want to set multiple filters before triggering a fetch, or simply re-apply the current filters. There is no manual trigger.

## Solution

Add an optional `onSearch` prop to `FilterBar`. When provided, a "Search" button renders at the end of the filter bar (right side).

## Component Changes

### `FilterBar` (`filter-bar.tsx`)

- Add `onSearch?: () => void` to `FilterBarProps`
- When present, render a button after `{trailing}` at the end of the flex row

### `DataTablePage` (`data-table-page.tsx`)

- Add `onSearch?: () => void` to `DataTablePageProps`
- Pass it to `FilterBar`

## Pages to Update

All pages that use `DataTablePage` with filters:

- `people/index.tsx` — pass `refetch` from `usePeople`
- `pets/-pets-page.tsx` — pass `refetch` from query
- `households/index.tsx` — pass `refetch` from query
- `audit-logs.tsx` — pass `refetch` from query

Report pages with custom filter layouts:

- `reports/people.tsx` — add button inline in filter section
- `reports/pets.tsx` — add button inline in filter section

## Button Spec

- Label: "Search" (use i18n key `common.search` if it exists, else add it)
- Style: secondary/outline button, same height as filter inputs (`h-9`)
- Position: end of the filter bar flex row, after all filters and `trailing`
- Behavior: calls `onSearch()` on click — no state changes, just triggers the callback
