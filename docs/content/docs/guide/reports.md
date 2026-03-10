---
title: Reports
weight: 7
---

Observer generates 39 structured report types designed to match actual Ukrainian NGO donor reporting obligations. All reports are scoped to a single project and filtered by a date range.

## Running a Report

All reports are available under:

```http
GET /projects/:project_id/reports/:report_type?start=2024-01-01&end=2024-12-31
```

Replace `:report_type` with one of the identifiers in the tables below.

## Count Types

Every report uses one of three counting methods:

| Type       | Meaning                                                                        |
| ---------- | ------------------------------------------------------------------------------ |
| **Events** | Rows in `support_records` — one person with three consultations counts 3       |
| **People** | Distinct individuals — same person counted once regardless of how many records |
| **Units**  | Distinct family units (one household = one unit)                               |

Date filtering uses `support_records.provided_at` for consultations and `people.registered_at` for registration reports — not `created_at`.

## Report Groups

### Group 1 — General Consultation Counts

| #   | Report                           | Type   |
| --- | -------------------------------- | ------ |
| 1   | Total consultations of all types | Events |
| 2   | Total legal consultations        | Events |
| 3   | Total social consultations       | Events |

### Group 2 — Sex Breakdown

| #   | Report                                  | Type   |
| --- | --------------------------------------- | ------ |
| 12  | Men registered in period                | People |
| 13  | Women registered in period              | People |
| 14  | Women who received legal consultations  | People |
| 15  | Women who received social consultations | People |
| 16  | Men who received legal consultations    | People |
| 17  | Men who received social consultations   | People |

Reports 14–17 count distinct people, not consultation events.

### Group 3 — Geographic / IDP Status

| #    | Report                                                            | Type   |
| ---- | ----------------------------------------------------------------- | ------ |
| 4    | Total people registered                                           | People |
| 5–6  | People from conflict zone registered                              | People |
| 7–10 | People from conflict zone who received legal/social consultations | People |
| 11   | Non-IDPs registered                                               | People |

IDP status is derived automatically: `origin_place_id → places → states.conflict_zone`. Reports 5–6 and 7–10 are parameterised by conflict zone label.

### Group 4 — Vulnerability Category Breakdown

| #   | Report                                                 | Type   |
| --- | ------------------------------------------------------ | ------ |
| 18  | People registered — by vulnerability category          | People |
| 19  | People who received social consultations — by category | People |
| 20  | People who received legal consultations — by category  | People |

People with no assigned category appear in an "uncategorised" bucket.

### Group 5 — Current Region of Stay

| #   | Report                                               | Type   |
| --- | ---------------------------------------------------- | ------ |
| 21  | People registered — by current region                | People |
| 22  | People who received legal consultations — by region  | People |
| 23  | People who received social consultations — by region | People |

### Group 6 — Support Sphere Breakdown

| #   | Report                                               | Type   |
| --- | ---------------------------------------------------- | ------ |
| 24  | Legal consultation count — by sphere                 | Events |
| 25  | People who received legal consultations — by sphere  | People |
| 29  | Social consultation count — by sphere                | Events |
| 30  | People who received social consultations — by sphere | People |

**Support spheres:**

| Value                   | Description                                       |
| ----------------------- | ------------------------------------------------- |
| `housing_assistance`    | Housing rights, eviction, social housing          |
| `document_recovery`     | Passports, birth certificates, property documents |
| `social_benefits`       | IDP registration, social payments                 |
| `property_rights`       | Property left in occupied territories             |
| `employment_rights`     | Labour law, dismissal, job placement              |
| `family_law`            | Divorce, custody, alimony                         |
| `healthcare_access`     | Medical coverage, disability documentation        |
| `education_access`      | School enrolment, educational rights              |
| `financial_aid`         | Emergency financial assistance                    |
| `psychological_support` | Mental health referrals, counselling              |
| `other`                 | Unlisted or cross-cutting topics                  |

### Group 7 — Office Breakdown

| #   | Report                                | Type   |
| --- | ------------------------------------- | ------ |
| 28  | Legal consultation count — by office  | Events |
| 32  | Social consultation count — by office | Events |
| 33  | Total consultation count — by office  | Events |

### Group 8 — Age Group Breakdown

| #   | Report                                                  | Type   |
| --- | ------------------------------------------------------- | ------ |
| 26  | Legal consultation count — by age group                 | Events |
| 27  | People who received legal consultations — by age group  | People |
| 31a | Social consultation count — by age group                | Events |
| 31b | People who received social consultations — by age group | People |
| 34  | Total consultation count — by age group                 | Events |

**Age groups:**

| Value               | Age range |
| ------------------- | --------- |
| `infant`            | 0–1       |
| `toddler`           | 1–3       |
| `pre_school`        | 3–6       |
| `middle_childhood`  | 6–12      |
| `young_teen`        | 12–14     |
| `teenager`          | 14–18     |
| `young_adult`       | 18–25     |
| `early_adult`       | 25–35     |
| `middle_aged_adult` | 35–55     |
| `old_adult`         | 55+       |

When `birth_date` is set and `age_group` is null, the application computes the bucket automatically.

### Group 9 — Tag Search

| #   | Report                                        | Type   |
| --- | --------------------------------------------- | ------ |
| 35  | Support records for people with specific tags | Events |
| 36  | People registered with specific tags          | People |

Pass one or more tag IDs as parameters. Useful for ad-hoc donor queries.

### Group 10 — Family Units

| #   | Report                                                      | Type           |
| --- | ----------------------------------------------------------- | -------------- |
| 37  | People and family members who received legal consultations  | People + Units |
| 38  | People and family members who received social consultations | People + Units |
| 39  | People and family members registered in period              | People + Units |

A family unit is one household record. Counts are returned as both total individuals and distinct household units.

## Pet Reports

Animal-related reports are available separately under:

```http
GET /projects/:project_id/pet-reports/:report_type
```

Pet reports cover species breakdown, vaccination status, sterilisation status, and pet-to-person ratios.
