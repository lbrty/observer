---
title: Export
weight: 8
---

Observer can export case data and reports as CSV files. Export requires the `consultant` project role **and** the `can_export` permission flag set on your project membership. Platform admins, staff, and project owners always have export access. See [Permissions](/docs/guide/permissions/) for details.

## Export Case Data

Export all people and their associated records for a project:

```http
GET /projects/:project_id/export/people?format=csv
```

Optional query parameters:

| Parameter     | Description                              |
| ------------- | ---------------------------------------- |
| `start`       | Filter by registration date (YYYY-MM-DD) |
| `end`         | Filter by registration date (YYYY-MM-DD) |
| `category_id` | Filter by vulnerability category         |
| `tag_id`      | Filter by tag                            |

The response is streamed — large datasets are sent progressively so the connection does not time out.

## Export Support Records

Export consultation records for a project:

```http
GET /projects/:project_id/export/support-records?format=csv&start=2024-01-01&end=2024-12-31
```

| Parameter | Description                                |
| --------- | ------------------------------------------ |
| `start`   | Filter by `provided_at` date (YYYY-MM-DD)  |
| `end`     | Filter by `provided_at` date (YYYY-MM-DD)  |
| `type`    | `legal` or `social`                        |
| `sphere`  | Support sphere (e.g. `housing_assistance`) |

## Export Migration Records

Export movement/displacement history:

```http
GET /projects/:project_id/export/migration-records?format=csv
```

## CSV Format

All exports use UTF-8 CSV with a header row. Fields follow the same structure as the API responses. Sensitive fields (`contact`, `personal`, `documents`) are included or omitted based on your project permission flags — the same rules that apply to API reads apply to exports.

## Downloading in the Web Interface

1. Open a project
2. Click **Export** in the navigation
3. Choose the record type and date range
4. Click **Download CSV**

The file downloads directly to your browser.

{{% callout type="note" %}}
Exports are scoped to your project. You cannot export data from a project you are not assigned to, and your sensitivity flags apply — if `can_view_personal` is off, the exported CSV will not contain national IDs or birth dates.
{{% /callout %}}
