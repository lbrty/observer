---
title: "ADR-005: Reports"
weight: 5
---

| Field      | Value                     |
| ---------- | ------------------------- |
| Status     | Proposed                  |
| Date       | 2026-02-21                |
| Supersedes | —                         |
| Components | observer, reports, schema |

---

## Context

This ADR captures the reporting requirements derived from the legacy `idp-archive/reports.txt` file. That document lists 39 report queries previously implemented against the Django/Python archive system. The schema defined in ADR-003 is shaped specifically to make these queries efficient.

All requirements are translated from Ukrainian and reorganised by theme.

---

## Conventions

**Count type** is explicit for every report:

- **Events** — rows in `support_records`; a person receiving three consultations contributes 3
- **People** — distinct `person_id` values; the same person counted once regardless of consultation count
- **Units** — distinct family units (heads of household)

**Date filtering:**

- Consultation date range → `support_records.provided_at` (not `created_at`)
- Registration date range → `people.registered_at` (not `people.created_at`; see §registered_at below)

**Project scoping:** All reports operate within a single project context unless stated otherwise.

---

## Report Requirements

### Group 1 — General Consultation Counts

| #   | Report                                             | Count type |
| --- | -------------------------------------------------- | ---------- |
| 1   | Total consultations of all types in a given period | Events     |
| 2   | Total **legal** consultations in a given period    | Events     |
| 3   | Total **social** consultations in a given period   | Events     |

**Schema:** `support_records` filtered by `provided_at` and `type`.

---

### Group 2 — Sex Breakdown

| #   | Report                                                        | Count type |
| --- | ------------------------------------------------------------- | ---------- |
| 12  | **Men** registered in a given period                          | People     |
| 13  | **Women** registered in a given period                        | People     |
| 14  | **Women** who received legal consultations in a given period  | People     |
| 15  | **Women** who received social consultations in a given period | People     |
| 16  | **Men** who received legal consultations in a given period    | People     |
| 17  | **Men** who received social consultations in a given period   | People     |

**Schema:** `people.sex` joined to `support_records` via `person_id`. Registration window uses `people.registered_at`.

> Reports 14–17 count distinct people, not consultation events. A woman who had three legal consultations in the period counts as 1.

---

### Group 3 — Geographic / IDP Status Breakdown

| #   | Report                                                             | Count type |
| --- | ------------------------------------------------------------------ | ---------- |
| 4   | Total people registered in a given period                          | People     |
| 5   | People from **Crimea** registered in a given period                | People     |
| 6   | People from **Eastern Ukraine (ATO)** registered in a given period | People     |
| 7   | People from Crimea who received **legal** consultations            | People     |
| 8   | People from Crimea who received **social** consultations           | People     |
| 9   | People from Eastern Ukraine who received **legal** consultations   | People     |
| 10  | People from Eastern Ukraine who received **social** consultations  | People     |
| 11  | **Non-IDPs** registered in a given period                          | People     |

**Schema:** IDP classification is derived via `people.origin_place_id → places.state_id → states.conflict_zone`. The `conflict_zone` column on `states` is a free-text label; values such as `'crimea'` and `'east_ukraine'` are seeded into the states table rather than stored on the person record.
Registration window uses `people.registered_at`.
Consultation window uses `support_records.provided_at` joined via `person_id`.

> Records where `origin_place_id IS NULL` or the resolved `conflict_zone IS NULL` are excluded from reports 5–10. Report 11 selects people whose `conflict_zone IS NULL` (no conflict-origin designation). Report 4 includes all records regardless of origin.

**Pattern B (updated) — People count by conflict zone:**

```sql
SELECT st.conflict_zone, COUNT(DISTINCT p.id) AS total
FROM   people p
JOIN   places pl ON pl.id = p.origin_place_id
JOIN   states st ON st.id = pl.state_id
WHERE  p.project_id    = :project_id
  AND  p.registered_at BETWEEN :start AND :end
  AND  st.conflict_zone IS NOT NULL
GROUP  BY st.conflict_zone;
```

---

### Group 4 — Vulnerability Category Breakdown

| #   | Report                                                                       | Count type |
| --- | ---------------------------------------------------------------------------- | ---------- |
| 18  | People registered in a given period — by **vulnerability category**          | People     |
| 19  | People who received **social** consultations — by **vulnerability category** | People     |
| 20  | People who received **legal** consultations — by **vulnerability category**  | People     |

**Schema:** `people.category_id → categories.name`. Joined to `support_records` for consultation-specific counts.

> People with `category_id IS NULL` appear as an "uncategorised" bucket or are excluded depending on the report context.

---

### Group 5 — Current Region of Stay Breakdown

| #   | Report                                                               | Count type |
| --- | -------------------------------------------------------------------- | ---------- |
| 21  | People registered — by **current region of stay**                    | People     |
| 22  | People who received **legal** consultations — by **current region**  | People     |
| 23  | People who received **social** consultations — by **current region** | People     |

**Schema:** `people.current_place_id → places.state_id → states.name`.

> Records with `current_place_id IS NULL` are excluded from regional breakdowns. This is expected for newly registered people whose location has not yet been entered. Applications should surface this as an "unknown region" bucket so the omission is visible.

---

### Group 6 — Support Sphere Breakdown

| #   | Report                                                   | Count type |
| --- | -------------------------------------------------------- | ---------- |
| 24  | **Legal** consultation count — by sphere of appeal       | Events     |
| 25  | People who received **legal** consultations — by sphere  | People     |
| 29  | **Social** consultation count — by sphere of appeal      | Events     |
| 30  | People who received **social** consultations — by sphere | People     |

**Schema:** `support_records.sphere` (`support_sphere` enum) grouped alongside `type`.

**Defined sphere values:**

| Value                   | Description                                                 |
| ----------------------- | ----------------------------------------------------------- |
| `housing_assistance`    | Housing rights, eviction, social housing, temporary shelter |
| `document_recovery`     | Passports, birth certificates, property documents           |
| `social_benefits`       | IDP registration, social payments, benefit eligibility      |
| `property_rights`       | Property left in occupied territories, restitution          |
| `employment_rights`     | Labour law, dismissal, job placement                        |
| `family_law`            | Divorce, custody, alimony, guardianship                     |
| `healthcare_access`     | Medical coverage, disability documentation                  |
| `education_access`      | School enrolment, tutoring, educational rights              |
| `financial_aid`         | Emergency financial assistance, humanitarian grants         |
| `psychological_support` | Mental health referrals, counselling coordination           |
| `other`                 | Unlisted or cross-cutting topics                            |

> Records with `sphere IS NULL` are excluded from sphere breakdowns. The sphere field is optional but encouraged for all legal and social records.

---

### Group 7 — Office Breakdown

| #   | Report                                                 | Count type |
| --- | ------------------------------------------------------ | ---------- |
| 28  | **Legal** consultation count — by **providing office** | Events     |
| 32  | **Social** consultation count — by providing office    | Events     |
| 33  | Total consultation count — by providing office         | Events     |

**Schema:** `support_records.office_id → offices.name`.

> `office_id` on `support_records` is the direct field for this dimension. It records which office provided or coordinated the consultation, independent of where the consultant's home office is. Records with `office_id IS NULL` are excluded from office breakdowns.

---

### Group 8 — Age Group Breakdown

| #   | Report                                                       | Count type |
| --- | ------------------------------------------------------------ | ---------- |
| 26  | **Legal** consultation count — by **age group** of recipient | Events     |
| 27  | People who received **legal** consultations — by age group   | People     |
| 31a | **Social** consultation count — by age group                 | Events     |
| 31b | People who received **social** consultations — by age group  | People     |
| 34  | Total consultation count — by age group                      | Events     |

**Schema:** `people.age_group` (`person_age_group` enum).

**Age group definitions:**

| Value               | Age range     |
| ------------------- | ------------- |
| `infant`            | 0–1           |
| `toddler`           | 1–3           |
| `pre_school`        | 3–6           |
| `middle_childhood`  | 6–12          |
| `young_teen`        | 12–14         |
| `teenager`          | 14–18         |
| `young_adult`       | 18–25         |
| `early_adult`       | 25–35         |
| `middle_aged_adult` | 35–55         |
| `old_adult`         | 55+           |
| `unknown`           | Not specified |

When `birth_date` is set and `age_group` is NULL, the application layer computes the bucket using `age = EXTRACT(YEAR FROM AGE(birth_date))` and maps to the ranges above. The `chk_people_age_xor` constraint ensures both are never set simultaneously.

---

### Group 9 — Tag Search

| #   | Report                                                                  | Count type |
| --- | ----------------------------------------------------------------------- | ---------- |
| 35  | Support records in a period associated with people having specific tags | Events     |
| 36  | People registered in a period whose records include specific tags       | People     |

**Schema:** `person_tags` junction (`person_id`, `tag_id`) → `tags.name`. Tags are project-scoped (`tags.project_id`).

> Report 35 filters on `support_records.provided_at`. Report 36 filters on `people.registered_at`. Both require a tag name match within the same project.

---

### Group 10 — Family Units

| #   | Report                                                                | Count type     |
| --- | --------------------------------------------------------------------- | -------------- |
| 37  | People and their family members who received **legal** consultations  | People + Units |
| 38  | People and their family members who received **social** consultations | People + Units |
| 39  | People and their family members registered in a given period          | People + Units |

**Schema:** `households` + `household_members`. A family unit is one `households` row; members are all rows in `household_members` with that `household_id`. Use `household_members.relationship = 'head'` to identify the nominated head. The `people.parent_id` self-reference has been replaced by this model (see ADR-003 migration 000023).

---

## Schema Implications

| Schema element                                           | Reports enabled                                     |
| -------------------------------------------------------- | --------------------------------------------------- |
| `people.registered_at`                                   | 4–6, 11–13, 18, 21, 36, 39                          |
| `people.origin_place_id → places → states.conflict_zone` | 4–11                                                |
| `people.sex`                                             | 12–17                                               |
| `person_categories → categories`                         | 18–20                                               |
| `people.current_place_id → places → states`              | 21–23                                               |
| `support_records.type`                                   | 1–3, 7–10, 14–17, 19–20, 22–25, 28–30, 32–34, 37–38 |
| `support_records.provided_at`                            | All consultation date-range filtering               |
| `support_records.sphere` (`support_sphere` enum)         | 24, 25, 29, 30                                      |
| `support_records.office_id → offices`                    | 28, 32, 33                                          |
| `people.age_group` / `people.birth_date`                 | 26, 27, 31, 34                                      |
| `person_tags` + `tags` (project-scoped)                  | 35, 36                                              |
| `households` + `household_members`                       | 37, 38, 39                                          |

---

## Reference Query Patterns

### Pattern A — Count by type and period

```sql
SELECT type, COUNT(*) AS total
FROM   support_records
WHERE  project_id  = :project_id
  AND  provided_at BETWEEN :start AND :end
GROUP  BY type;
```

### Pattern B — People count by IDP status and registration window

```sql
SELECT idp_status, COUNT(*) AS total
FROM   people
WHERE  project_id    = :project_id
  AND  registered_at BETWEEN :start AND :end
GROUP  BY idp_status;
```

### Pattern C — Distinct people by sex who received a consultation type

```sql
SELECT p.sex, COUNT(DISTINCT p.id) AS people
FROM   support_records sr
JOIN   people p ON p.id = sr.person_id
WHERE  sr.project_id  = :project_id
  AND  sr.type        = 'legal'
  AND  sr.provided_at BETWEEN :start AND :end
GROUP  BY p.sex;
```

### Pattern D — People by current region (NULL-aware)

```sql
SELECT
    COALESCE(st.name, 'unknown') AS region,
    COUNT(DISTINCT p.id)         AS people
FROM   people p
LEFT   JOIN places pl ON pl.id = p.current_place_id
LEFT   JOIN states st ON st.id = pl.state_id
WHERE  p.project_id    = :project_id
  AND  p.registered_at BETWEEN :start AND :end
GROUP  BY st.name;
```

> LEFT JOIN surfaces the "unknown region" bucket for records with `current_place_id IS NULL`.

### Pattern E — Consultations by providing office

```sql
SELECT COALESCE(o.name, 'unknown') AS office, COUNT(sr.id) AS consultations
FROM   support_records sr
LEFT   JOIN offices o ON o.id = sr.office_id
WHERE  sr.project_id  = :project_id
  AND  sr.type        = 'legal'
  AND  sr.provided_at BETWEEN :start AND :end
GROUP  BY o.name;
```

### Pattern F — Age group from birth_date (application-layer bucketing)

When `age_group IS NULL` and `birth_date IS NOT NULL`, compute the bucket in the query:

```sql
SELECT
    CASE
        WHEN age <  1  THEN 'infant'
        WHEN age <  3  THEN 'toddler'
        WHEN age <  6  THEN 'pre_school'
        WHEN age < 12  THEN 'middle_childhood'
        WHEN age < 14  THEN 'young_teen'
        WHEN age < 18  THEN 'teenager'
        WHEN age < 25  THEN 'young_adult'
        WHEN age < 35  THEN 'early_adult'
        WHEN age < 55  THEN 'middle_aged_adult'
        ELSE                'old_adult'
    END                       AS age_bucket,
    COUNT(DISTINCT p.id)      AS people
FROM (
    SELECT id,
           EXTRACT(YEAR FROM AGE(birth_date))::INT AS age
    FROM   people
    WHERE  project_id = :project_id
      AND  birth_date IS NOT NULL
      AND  age_group  IS NULL
) p
JOIN support_records sr ON sr.person_id = p.id
WHERE sr.type        = 'legal'
  AND sr.provided_at BETWEEN :start AND :end
GROUP BY age_bucket;
-- Union with rows where age_group IS NOT NULL to cover both input paths.
```

### Pattern G — Family units

```sql
SELECT
    COUNT(DISTINCT hm.household_id) AS family_units,
    COUNT(DISTINCT hm.person_id)    AS total_individuals
FROM   household_members hm
JOIN   households h        ON h.id  = hm.household_id
JOIN   support_records sr  ON sr.person_id = hm.person_id
WHERE  h.project_id   = :project_id
  AND  sr.project_id  = :project_id
  AND  sr.type        = 'legal'
  AND  sr.provided_at BETWEEN :start AND :end;
```

> Uses `households` + `household_members`. No recursive CTE needed — the household entity is explicit and flat.

---

## Design Notes

### `registered_at` vs `created_at`

`people.created_at` is the database insertion timestamp. Batch imports (e.g. digitising paper intake forms after a field visit) give all rows the same `created_at`. `registered_at DATE` is the actual date a person was registered with the office. Reports filtering by registration window use `registered_at`; it is nullable (not set for legacy imports where the original date is unknown).

### `support_records.office_id`

Office attribution on a consultation is a direct field, not derived through `consultant_id → users.office_id`. A consultant from the Kyiv office may provide a consultation coordinated by the Kherson office during a field visit — routing through the consultant's home office would misattribute it. The providing office is recorded explicitly at the time of entry.

### `support_sphere` enum

Sphere values use snake_case. New spheres are added via a forward migration (`ALTER TYPE support_sphere ADD VALUE`). Free-text entry is not permitted — consistent values are required for GROUP BY in reports 24, 25, 29, 30.

### Null handling in breakdowns

Geographic (Group 5), sphere (Group 6), and office (Group 7) breakdowns silently drop records with NULL in the grouping column unless explicitly handled with `LEFT JOIN` + `COALESCE`. All reference query patterns above use `LEFT JOIN` to surface the null bucket as `'unknown'`.

---

## Legacy System Mapping

| Archive (Django)       | Observer                           |
| ---------------------- | ---------------------------------- |
| `Settler` model        | `people` table                     |
| `aid_list` relation    | `support_records` table            |
| `aid_type = 'l'`       | `support_records.type = 'legal'`   |
| `aid_type = 's'`       | `support_records.type = 'social'`  |
| `gender = 'f'` / `'m'` | `people.sex = 'female'` / `'male'` |
