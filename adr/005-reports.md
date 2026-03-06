# ADR-005: Reports


| Field      | Value                     |
| ---------- | ------------------------- |
| Status     | Proposed                  |
| Date       | 2026-02-21                |
| Supersedes | —                         |
| Components | observer, reports, schema |


## Conventions

**Count type** is explicit for every report:

- **Events** — rows in `support_records`; a person receiving three consultations contributes 3
- **People** — distinct `person_id` values; the same person counted once regardless of consultation count
- **Units** — distinct family units (heads of household)

**Date filtering:**

- Consultation date range → `support_records.provided_at` (not `created_at`)
- Registration date range → `people.registered_at` (not `people.created_at`; see §registered_at below)

**Project scoping:** All reports operate within a single project context unless stated otherwise.


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


### Group 4 — Vulnerability Category Breakdown

| #   | Report                                                                       | Count type |
| --- | ---------------------------------------------------------------------------- | ---------- |
| 18  | People registered in a given period — by **vulnerability category**          | People     |
| 19  | People who received **social** consultations — by **vulnerability category** | People     |
| 20  | People who received **legal** consultations — by **vulnerability category**  | People     |

**Schema:** `people.category_id → categories.name`. Joined to `support_records` for consultation-specific counts.

> People with `category_id IS NULL` appear as an "uncategorised" bucket or are excluded depending on the report context.


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


### Group 10 — Family Units

| #   | Report                                                                | Count type     |
| --- | --------------------------------------------------------------------- | -------------- |
| 37  | People and their family members who received **legal** consultations  | People + Units |
| 38  | People and their family members who received **social** consultations | People + Units |
| 39  | People and their family members registered in a given period          | People + Units |

**Schema:** `households` + `household_members`. A family unit is one `households` row; members are all rows in `household_members` with that `household_id`. Use `household_members.relationship = 'head'` to identify the nominated head. The `people.parent_id` self-reference has been replaced by this model (see ADR-003 migration 000023).


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


## Legacy System Mapping

| Archive (Django)       | Observer                           |
| ---------------------- | ---------------------------------- |
| `Settler` model        | `people` table                     |
| `aid_list` relation    | `support_records` table            |
| `aid_type = 'l'`       | `support_records.type = 'legal'`   |
| `aid_type = 's'`       | `support_records.type = 'social'`  |
| `gender = 'f'` / `'m'` | `people.sex = 'female'` / `'male'` |
