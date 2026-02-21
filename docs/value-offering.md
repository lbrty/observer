# Observer — Value Offering

---

## Positioning

Observer fills the gap between a paper ledger and a full UNHCR-scale system: a **self-hosted, project-scoped case management tool for organizations too small or too politically constrained to use proGres v4 or Primero**, with enough structure to meet the specific reporting obligations Ukrainian NGOs face toward donors.

The incumbents either cannot be deployed by these organizations or do not produce the right reports. That is the gap Observer occupies.

---

## Target User

A Ukrainian or Ukraine-supporting NGO with:

- 5–50 staff, one or zero dedicated IT people
- One or more projects receiving legal, social, or humanitarian aid funding
- Donor reporting requirements expressed as the 39 report types in ADR-005
- Sensitivity requirements that prevent uploading beneficiary data to third-party SaaS platforms
- No existing relationship with UNICEF, UNHCR, or IOM that would grant access to Primero or proGres

---

## Why the Incumbents Do Not Serve This Market

| System | Why it is unavailable to a small NGO |
| --- | --- |
| **proGres v4 / PRIMES** | Requires UNHCR operational partnership; not self-hostable; global system with country-level access governance |
| **Primero** | Requires a UNICEF/IRC technical partner for deployment and configuration; server infrastructure managed by the partner, not the NGO |
| **ActivityInfo** | SaaS with per-user pricing; designed for aggregate monitoring (5W/3W), not individual case management |
| **KoBoToolbox** | Survey collection only; no persistent case entity between submissions |
| **MiMOSA / DTM** | IOM internal system; not available to external organizations |
| **Ukraine IDP Register** | Government system; NGOs can query but not contribute structured case data |

---

## Differentiators

### 1. Actually deployable

A Go binary, a PostgreSQL database, and a Justfile. A single sysadmin can run it on a €20/month VPS. No partner onboarding, no SaaS agreement, no UN agency sponsorship required.

### 2. `support_sphere` — constrained consultation topic taxonomy

Observer's 11-value `support_sphere` enum (`housing_assistance`, `document_recovery`, `social_benefits`, `property_rights`, `employment_rights`, `family_law`, `healthcare_access`, `education_access`, `financial_aid`, `psychological_support`, `other`) is a first-class constrained field enabling GROUP BY-safe breakdown by consultation topic.

No incumbent system offers this at the schema level. Primero has configurable service types but no GROUP-BY-safe topic taxonomy. proGres's assistance module is coarser. This directly produces the "by sphere of appeal" breakdowns that Ukrainian NGOs report to legal aid donors (EU4Justice, USAID, donor-funded legal clinics).

### 3. Dual-level RBAC — declarative and schema-enforced

`platform_role` (`admin / staff / consultant / guest`) combined with `project_role` (`owner / manager / consultant / viewer`) and three data-sensitivity flags (`can_view_contact`, `can_view_personal`, `can_view_documents`) gives more expressive access control than most off-the-shelf tools without configuration ceremony.

ActivityInfo achieves similar granularity but requires manual per-grant configuration. Observer's model is declarative, schema-enforced, and immediately queryable.

### 4. `registered_at` separate from `created_at`

Batch imports from field visits (paper intake forms digitised after a visit) give all records the same `created_at`. `registered_at DATE` records the actual date a person was formally registered. Almost every small system gets this wrong, corrupting all registration-window reports.

### 5. Forward-only migrations

Most humanitarian systems in the Ukraine context run on whatever schema the local IT person deployed years ago. Observer's ADR-004 migration discipline is operationally unusual and valuable for a system expected to evolve over years of changing donor requirements.

### 6. 39 report types matching Ukrainian NGO donor obligations

ADR-005 encodes all 39 report queries previously implemented in the legacy system, covering: consultation counts by type, IDP geographic origin breakdowns, sex breakdowns, vulnerability category breakdowns, regional breakdowns, sphere of appeal breakdowns, office breakdowns, age group breakdowns, tag-based searches, and family unit counts.

These are the specific reports Ukrainian legal aid and social service organizations submit to EU, USAID, and bilateral donors. The schema is shaped from day one to make these queries efficient.

### 7. IDP classification derived from geography, not hardcoded enum

IDP origin classification (`crimea`, `east_ukraine`, or any future zone) is derived at query time via `origin_place_id → places.state_id → states.conflict_zone`. New conflict zones are added by seeding the `states` table — no schema migration required. This mirrors how the official Ukraine IDP register derives the classification (from prior registered address) rather than storing a hardcoded political category.

### 8. Households as first-class entities with typed relationships

The `households` + `household_members` model (relationship types: head, spouse, child, parent, sibling, grandchild, grandparent, other_relative, non_relative) supports the family-unit reports required by donors and enables group-level document attribution (housing allocations, ration cards, evacuation orders). Most small systems either omit family structure entirely or use a fragile self-referencing FK that cannot represent peer relationships.

---

## What Observer Does Not Compete On

| Capability | Who owns it | Observer's stance |
| --- | --- | --- |
| Scale (10M+ registrations) | proGres v4 | Not a target; single-project deployments |
| Inter-agency referral workflows | Primero | `referral_status` column is a lightweight tracker, not a full referral module |
| Biometric deduplication | proGres + BIMS | `external_id` unique index is the deduplication anchor; national ID (РНОКПП) required |
| Cluster / 5W / 3W aggregate reporting | ActivityInfo | Observer produces the raw query results; ActivityInfo is still the right tool for cluster-level dashboards |
| Cross-country mandate reporting | proGres / UNHCR Dataport | Out of scope |

---

## Summary

Observer is purpose-built for the organization that is too small for proGres, too politically independent for Primero, and too data-sensitive for SaaS. The reporting schema matches what Ukrainian NGOs actually submit. The deployment model matches what they can actually operate.
