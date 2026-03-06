---
title: Landscape
weight: 2
---

A comparison of Observer against established IDP and humanitarian case management systems.

## Systems

| System | Owner | Use |
| --- | --- | --- |
| **proGres v4 / PRIMES** | UNHCR | Refugee/IDP registration, multi-module case management |
| **Primero** | UNICEF / UNHCR (open source) | Child protection, GBV case management |
| **ActivityInfo** | ActivityInfo (SaaS) | Humanitarian monitoring, protection IM |
| **KoBoToolbox** | KoBoToolbox Foundation (open source) | Field data collection |
| **MiMOSA / DTM** | IOM | Migrant assistance, displacement tracking |
| **OSCaR** | Children in Families (open source) | Vulnerable children (SE Asia) |

## How Observer compares

### Entity model

| System | Individual | Family | Approach |
| --- | --- | --- | --- |
| Observer | `people` (individual-centric) | `households` + `household_members` with typed relationships | First-class household entity, 9 relationship types |
| proGres v4 | Individual Record | Registration Group | Group-first with individuals nested within |
| Primero | Case (= individual) | Family record overlay | 1:1 case/individual, family as subform |
| ActivityInfo | Anonymous serial | Subform | No persistent entity |

### Access control

| System | Platform roles | Project scoping |
| --- | --- | --- |
| Observer | admin / staff / consultant / guest | project_role + 3 sensitivity flags |
| Primero | Admin / Manager / Caseworker | Agency + module scoped |
| proGres v4 | Module-specific roles | Operation (country) scoped |
| ActivityInfo | Owner / Admin / Editor / Viewer | Per-resource grants |

### Support / aid taxonomy

| System | Approach |
| --- | --- |
| Observer | 6 support types + 11 support spheres (schema-constrained enums) |
| proGres v4 | Assistance module: cash, in-kind, documentation, legal, medical, etc. |
| Primero | Configurable service types |
| MiMOSA | By stage × type matrix |

### Movement tracking

| System | Approach |
| --- | --- |
| Observer | `migration_records`: from/to places, date, movement reason, housing at destination |
| proGres v4 | Process status + Voluntary Repatriation module |
| MiMOSA / DTM | Corridor analysis, transport mode, departure/arrival dates |
| Primero | No movement tracking |

### Reporting

| System | Approach |
| --- | --- |
| Observer | 39 per-project SQL report types matching Ukrainian NGO donor obligations |
| proGres v4 | UNHCR Dataport, global population statistics |
| Primero | Caseload dashboards, age/sex disaggregation |
| ActivityInfo | Pivot tables, maps, 5W/3W cluster reporting |

## What Observer gets right

- **Individual-centric model** with household overlay — consistent with Primero, OSCaR, and humanitarian field consensus
- **Dual-level RBAC** — more expressive than most tools without configuration ceremony
- **`origin_place_id` vs `current_place_id`** — correct semantic separation between biography and current location
- **`registered_at` separate from `created_at`** — prevents batch import from corrupting registration date reports
- **`support_sphere` enum** — 11 constrained values enabling GROUP BY-safe consultation topic breakdowns
- **Forward-only migrations** — operational discipline rare in this sector
- **Project-scoped tags** — avoids cross-project vocabulary pollution
- **Households with typed relationships** — head, spouse, child, parent, sibling, grandchild, grandparent, other_relative, non_relative
- **Multi-code vulnerability classification** — `person_categories` junction table, aligned with UNHCR PSN codes
- **Consent tracking** — `consent_given` + `consent_date` for GDPR compliance
- **Movement causality** — `movement_reason` and `housing_at_destination` on migration records
- **Referral tracking** — `referral_status` + `referred_to_office` on support records

## Known limitations

- **Single project per person** — a person enrolled in multiple projects requires separate records
- **No inter-agency referral workflows** — referral tracking is lightweight (status field), not a full Primero-style referral module
- **No biometric deduplication** — relies on `external_id` (РНОКПП) partial unique index per project
- **Single case status per person** — cannot express simultaneous case states across service streams
