# Component Extraction — Design Spec

**Date:** 2026-04-20

## Goal

Extract reusable sub-components from all oversized frontend files so that no component file exceeds 150 lines (170 max where formatting forces it). Route files are not split but their extractable UI pieces are moved to `components/`.

## Constraints

- **Line limit:** 150 lines hard, ~170 soft ceiling for formatting-heavy files.
- **No directory reorganisation:** `components/` root stays flat.
- **No route file splitting:** Routes remain single files; sub-components are extracted out.
- **Extraction pattern:** Sibling files for component families (e.g. `search-palette/index.tsx` + `search-palette/result.tsx`). Sub-components sourced from route files land in `components/`, not alongside the route.
- **Imports:** Follow existing `@/` alias convention; colocated siblings use `./`.

---

## Extraction Targets

### A — Component files (sibling extraction)

Each file below is converted to a folder. The index re-exports so existing imports don't change.

| Current file                               | New structure                         | Sub-components to extract                      |
| ------------------------------------------ | ------------------------------------- | ---------------------------------------------- |
| `components/search-palette.tsx`            | `search-palette/index.tsx`            | `result.tsx` — SearchResult, SearchResultGroup |
| `components/charts/bar-chart.tsx`          | `charts/bar-chart/index.tsx`          | `axes.tsx`, `bars.tsx`, `tooltip.tsx`          |
| `components/charts/sankey-chart.tsx`       | `charts/sankey-chart/index.tsx`       | `nodes.tsx`, `links.tsx`, `labels.tsx`         |
| `components/date-picker.tsx`               | `date-picker/index.tsx`               | `trigger.tsx`, `nav.tsx`                       |
| `components/profile/mfa-settings.tsx`      | `profile/mfa-settings/index.tsx`      | `setup.tsx`, `backup-codes.tsx`, `disable.tsx` |
| `components/permissions/assign-dialog.tsx` | `permissions/assign-dialog/index.tsx` | `permission-row.tsx`                           |

Drawer indexes that still exceed 150 lines after existing sub-sections get additional extraction:

| File                                      | Sub-components to extract                     |
| ----------------------------------------- | --------------------------------------------- |
| `migration-record-drawer/index.tsx` (252) | `movement-section.tsx`, `housing-section.tsx` |
| `person-drawer/index.tsx` (238)           | `identity-section.tsx`, `contact-section.tsx` |
| `support-record-drawer/index.tsx` (211)   | `referral-section.tsx`                        |
| `household-drawer/index.tsx` (195)        | `members-section.tsx`                         |

### B — Route files (extraction to `components/`)

Sub-components are new files in `components/`. Route files shrink but are not split.

| Route                                   | Extracted to `components/`                                                                       |
| --------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `reports/people.lazy.tsx` (518)         | `reports/people-kpi-cards.tsx`, `reports/people-chart-section.tsx`, `reports/report-filters.tsx` |
| `reports/pets.lazy.tsx` (465)           | `reports/pets-kpi-cards.tsx`, `reports/pets-chart-section.tsx`                                   |
| `documents.lazy.tsx` (397)              | `document-row.tsx`, `upload-zone.tsx`, `document-preview.tsx`                                    |
| `people/index.tsx` (334)                | `people-columns.tsx`, `people-filters.tsx`                                                       |
| `admin/users/index.lazy.tsx` (312)      | `user-row.tsx`, `role-select.tsx`                                                                |
| `support-records/-page.tsx` (309)       | `support-record-row.tsx`, `support-record-filters.tsx`                                           |
| `my-stats/index.lazy.tsx` (307)         | `stat-card.tsx`, `activity-section.tsx`                                                          |
| `people/$personId/index.lazy.tsx` (306) | `person-header.tsx`, `person-details.tsx`                                                        |
| `reports/custom.lazy.tsx` (300)         | `reports/custom-report-form.tsx`, `reports/report-result.tsx`                                    |
| `pets/-page.tsx` (286)                  | `pet-row.tsx`, `pet-filters.tsx`                                                                 |

---

## Approach

**B — Sibling files.** Each oversized component becomes a folder. The `index.tsx` orchestrates and re-exports the public API. Sub-files are co-located siblings and unexported outside the folder unless they have standalone use.

This follows the existing drawer pattern (`person-drawer/index.tsx` + section files) which is already proven in the codebase.

---

## Import Compatibility

Converting `search-palette.tsx` → `search-palette/index.tsx` is transparent to all importers because bundlers and TypeScript resolve `@/components/search-palette` to either the file or the folder's `index.tsx`. No import sites need to change.

---

## New Directories

`components/reports/` will be created to house report-specific sub-components extracted from route files. This is a new directory, not a reorganisation of existing files — existing `components/` files stay where they are.

---

## What Is Not Changing

- Hook files (all under 90 lines — within limit).
- Store files.
- Route file boundaries.
- Directory structure of `components/` root.
- Any file currently under 150 lines.

---

## Success Criteria

1. No `.tsx` file in `components/` or `profile/` or `permissions/` exceeds 170 lines.
2. All existing imports resolve without change.
3. `bun run build` passes with no type errors.
4. Extracted sub-components are not exported beyond their folder unless they have a standalone use case.
