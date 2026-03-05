---
title: Existing Solutions
weight: 2
---

This document surveys established IDP and humanitarian case management systems, compares their design decisions against Observer's current architecture, identifies gaps and misalignments, and proposes concrete improvements.

---

## Systems Surveyed

| System                   | Owner                                | Primary Use                                                |
| ------------------------ | ------------------------------------ | ---------------------------------------------------------- |
| **Primero**              | UNICEF / UNHCR (open source)         | Child protection, GBV case management                      |
| **proGres v4 / PRIMES**  | UNHCR                                | Refugee and IDP registration, multi-module case management |
| **OSCaR**                | Children in Families (open source)   | Vulnerable children and family strengthening (SE Asia)     |
| **ActivityInfo**         | ActivityInfo (SaaS)                  | Humanitarian monitoring, reporting, protection IM          |
| **KoBoToolbox**          | KoBoToolbox Foundation (open source) | Field data collection, needs assessments                   |
| **MiMOSA / DTM**         | IOM                                  | Migrant assistance management, displacement tracking       |
| **Ukraine IDP Register** | Ukraine Ministry of Social Policy    | Official IDP status registration, benefit eligibility      |

---

## Comparison Matrix

### 1. Primary Entity Model

| System           | Individual entity            | Case entity                      | Distinction                                         |
| ---------------- | ---------------------------- | -------------------------------- | --------------------------------------------------- |
| Observer         | `people`                     | —                                | Case IS the person record; no separation            |
| Primero          | Case (= individual)          | Case                             | 1:1 — Case and individual are the same object       |
| proGres v4       | Individual Record            | Protection/RSD/resettlement Case | 1:many — one person, multiple active case processes |
| OSCaR            | Client                       | Case linked to Client            | 1:many — clients can hold multiple open cases       |
| ActivityInfo     | Serial number (anonymous)    | Record                           | Flexible, implementation-dependent                  |
| KoBoToolbox      | Submission                   | None                             | No persistent entity between submissions            |
| MiMOSA           | Individual / Beneficiary     | Assistance Case                  | 1:many                                              |
| Ukraine Register | Individual (tax ID / РНОКПП) | IDP status record                | Effectively 1:1 per status period                   |

### 2. Family / Household Model

| System           | Household entity                     | Approach                                                                                                  |
| ---------------- | ------------------------------------ | --------------------------------------------------------------------------------------------------------- |
| Observer         | Implicit via `parent_id`             | Self-referencing FK on `people`; one level only                                                           |
| proGres v4       | **Registration Group** (first-class) | Group has its own ID and reference number; individuals nested within                                      |
| Primero          | Family record (v2+)                  | Overlay linking individual Case records; family members also storable as a non-case subform within a Case |
| OSCaR            | Implicit                             | Client-centric; family grouping through program enrollment                                                |
| ActivityInfo     | Subform / repeat group               | No persistent household entity                                                                            |
| KoBoToolbox      | Repeat group                         | Within-submission roster only                                                                             |
| MiMOSA           | Household                            | Collected alongside the individual record                                                                 |
| Ukraine Register | Application bundle                   | Family bundled at application time; individual records per person                                         |

### 3. IDP / Displacement Origin Tracking

| System           | Approach                                                                      |
| ---------------- | ----------------------------------------------------------------------------- |
| Observer         | Hardcoded enum: `crimea / east_ukraine / non_idp`                             |
| proGres v4       | Population type field + last habitual residence location (place/settlement)   |
| Ukraine Register | Prior registered address (administrative address, not a category)             |
| DTM Ukraine      | Oblast of origin captured as a survey field; no hardcoded categories          |
| Primero          | No native field; requires custom configuration (lookup list of oblasts/zones) |
| ActivityInfo     | Custom reference fields per deployment                                        |

### 4. Vulnerability Classification

| System           | Approach                                                                                    |
| ---------------- | ------------------------------------------------------------------------------------------- |
| Observer         | `categories` table (free, single FK per person)                                             |
| proGres v4       | UNHCR Specific Needs Codes (PSN): 11 categories, 71 subcodes; multiple codes per individual |
| Primero          | "Protection Concerns" configurable lookup; multiple per case                                |
| ActivityInfo     | No built-in taxonomy; implementers use UNHCR PSN codes as custom lookup values              |
| Ukraine Register | No structured vulnerability taxonomy; linked to benefit eligibility rules                   |

### 5. Support / Aid Type Taxonomy

| System           | Approach                                                                                                                           |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Observer         | `support_type` enum (6 values) + `support_sphere` enum (11 values)                                                                 |
| proGres v4       | Assistance module: cash, in-kind, documentation, protection interventions, legal, medical, psychosocial, shelter                   |
| Primero          | Configurable service types; Primero–proGres interoperability defines 4 cross-agency types                                          |
| MiMOSA           | By stage (pre-departure, in-transit, post-arrival) × type (transportation, reception, medical, legal, psychosocial, reintegration) |
| Ukraine Register | Linked to social benefit codes, not service types                                                                                  |

### 6. Movement / Displacement Tracking

| System           | Approach                                                                                                                                             |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Observer         | `migration_records` table: `from_place_id`, `destination_place_id`, `migration_date`, `notes`                                                        |
| proGres v4       | Process status (Active/Hold/Inactive/Closed) + Voluntary Repatriation module + RApp for border/transit point monitoring                              |
| MiMOSA / DTM     | Most sophisticated: corridor analysis, mode of transport, reason for movement, departure/arrival dates; DTM publishes aggregate flow monitoring data |
| Primero          | No movement tracking module                                                                                                                          |
| Ukraine Register | Not tracked; only current place of residence vs. prior registered address                                                                            |

### 7. Roles and Permissions

| System       | Platform roles                                        | Project/program-level scoping                                                         |
| ------------ | ----------------------------------------------------- | ------------------------------------------------------------------------------------- |
| Observer     | `admin / staff / consultant / guest`                  | `project_role` (owner/manager/consultant/viewer) + 3 sensitivity flags                |
| Primero      | Admin / Manager / Caseworker                          | Agency-scoped; module-scoped; record-level "created by" scoping                       |
| proGres v4   | Module-specific roles (Registrar, Case Manager, etc.) | Operation (country) scoped; module-level restricted access                            |
| ActivityInfo | Owner / Admin / Editor / Viewer per resource          | Grant-based with parameter binding to restrict by reference value (e.g., partner org) |
| OSCaR        | Admin / Manager / Social Worker                       | Organization-scoped (multi-tenant)                                                    |

### 8. Reporting Capabilities

| System       | Strengths                                                                                             |
| ------------ | ----------------------------------------------------------------------------------------------------- |
| Observer     | Per-project SQL reports; see ADR-005                                                                  |
| proGres v4   | UNHCR Dataport, population statistics, mandate reporting; feeds global statistics portal              |
| Primero      | Caseload dashboards; age/sex disaggregation; CSV/PDF/Excel export                                     |
| ActivityInfo | Pivot tables, maps, charts; Power BI integration; cluster-level 5W/3W reporting; used as Ukraine RPMP |
| KoBoToolbox  | Basic form-level charts; data export for external analysis                                            |
| DTM          | Open data on HDX; authoritative IDP presence estimates                                                |

### 9. Document Management

| System           | Approach                                                                        |
| ---------------- | ------------------------------------------------------------------------------- |
| Observer         | `documents` table: path (relative), encryption_key_ref, mime_type, size         |
| Primero          | Photo and document attachments per case                                         |
| proGres v4       | Documentation module: registration cards, attestation letters, issuance history |
| ActivityInfo     | File attachment fields within records                                           |
| Ukraine Register | Certificate issuance tracked in register; no document content storage           |

### 10. Deduplication / Identity Resolution

| System           | Approach                                                                           |
| ---------------- | ---------------------------------------------------------------------------------- |
| Observer         | `external_id` field (no uniqueness constraint); no deduplication logic             |
| proGres v4       | BIMS (biometrics) integration; UNHCR individual ID; duplicate detection algorithms |
| Ukraine Register | РНОКПП (tax registration number) as universal unique identifier                    |
| Primero          | No biometrics; relies on case worker diligence; duplicate cases possible           |
| KoBoToolbox      | None                                                                               |

---

## What Observer Gets Right

### Individual-centric model

Consistent with Primero and OSCaR. The humanitarian field has largely converged on individual-centric records with household as an overlay, not the other way around (the exception being proGres, which is household-first for entitlement attribution). For a protection case management context — which Observer serves — individual-centric is the correct primary model.

### Dual-level RBAC (platform + project)

The `admin/staff/consultant/guest` platform roles combined with `project_role` (owner/manager/consultant/viewer) + data sensitivity flags is more expressive than most off-the-shelf tools. ActivityInfo achieves similar granularity but requires more manual configuration per grant. Observer's approach is declarative and schema-enforced.

### `origin_place_id` vs `current_place_id`

The semantic separation between where a person is from (biography) and where they are now (operational) is correct and mirrors proGres's `last place of habitual residence` vs. `current location`. Many simpler systems collapse this.

### `migration_records` as a separate table

Tracking movement as discrete historical events rather than overwriting a current location is the right approach and consistent with MiMOSA's movement tracking model. The immutability constraint (no `updated_at`) is sound.

### `registered_at` separate from `created_at`

This distinction — which almost all systems get wrong in their initial implementation — is explicitly present. Batch imports create false registration date clusters when `created_at` is used. Observer's `registered_at DATE` field mirrors proGres's mandatory "registration date" biographical field.

### `support_sphere` enum

Most humanitarian systems operate with coarse aid type categories and no topic-level breakdown. Observer's `support_sphere` with 11 constrained values is more expressive than proGres's assistance module and is directly required by the operational reports (ADR-005). The constrained enum prevents the GROUP BY fragmentation that plagues free-text sphere fields in other systems.

### Forward-only migrations

Not universal in this space. Many humanitarian systems use database dumps and manual schema patches in production (especially legacy systems like MiMOSA's 32-application ecosystem). Observer's disciplined migration policy is a genuine operational advantage.

### Project-scoped tags

Global tag namespaces are a well-known problem in Primero deployments where different country configurations share a single installation. Observer's project-scoped unique constraint on tags avoids cross-project vocabulary pollution.

---

## What Is Misaligned with Reality

### 1. The family model is structurally inadequate

Observer's `people.parent_id` self-reference encodes family relationships as a single-level parent/child hierarchy. In practice:

- A **spouse/partner** is a peer, not a "child" of the head. The current model cannot represent a married couple without making one a "child" of the other.
- **Relationship type** is not captured. Whether a person is a spouse, child, parent, sibling, or other dependent matters for entitlement calculations (proGres tracks this) and for family-unit reports.
- The **Registration Group** concept in proGres — a first-class entity with its own reference number — reflects a real operational need. Family units receive group-level documents (ration cards, evacuation orders, housing allocations). Without a group entity, all of these require application-level workarounds.
- The `people.parent_id IS NULL` assumption in the family unit query (ADR-005 Pattern G) is incorrect: a person who is genuinely the family head and a person whose family relationship has simply never been entered both have `parent_id IS NULL`. These are indistinguishable in the current schema.

### 2. `case_status` on `people` conflates person identity with case lifecycle

A person record in Observer is permanently tied to one project and has one `case_status`. In real operations:

- A person may be discharged from one service stream (housing resolved) while remaining active in another (employment support ongoing). The current model cannot express two simultaneous case states for the same person.
- A person may leave and return to a project area. Re-opening requires either creating a duplicate person record or resurrecting a closed record — neither is clean.
- proGres solves this by separating Individual Records (biographical, persistent) from Cases (module-specific, can be opened and closed independently). Primero addresses it partially through inter-agency referrals and service records.

### 3. `idp_status` is a hardcoded geopolitical taxonomy

The enum `('crimea', 'east_ukraine', 'non_idp')` encodes the 2014-era conflict geography. By 2022, the conflict expanded significantly: Zaporizhzhia, Mykolaiv, Kherson, Kharkiv, and parts of Odesa oblast became displacement origin zones. The enum cannot accommodate this without a schema migration.

The official Ukraine IDP register does not maintain this three-way classification at the data layer. It records the prior registered address (the actual place/settlement) and derives the category analytically. DTM surveys capture origin at the oblast level. Hardcoding a categorical taxonomy in the schema creates brittleness; the classification logic belongs in application code or a lookup table, not an enum.

### 4. No deduplication or identity resolution mechanism

Observer's `external_id` field is nullable with no uniqueness constraint. In the Ukrainian context, the РНОКПП (individual tax number) is the authoritative identity anchor — it links to the national IDP register, pension fund, and social benefit systems. Without enforcing uniqueness of `external_id` within a project (at minimum), a person can be registered multiple times — a known, prevalent problem in all humanitarian systems that lack biometrics or a national ID anchor.

### 5. No referral workflow

Primero's most operationally significant feature is inter-agency referrals: a caseworker can refer a client to a specific service provider (another NGO, a government office, a clinic), and the referral is tracked from sent → accepted → completed/rejected. This closes the loop on service delivery.

Observer's `support_records` records delivered services but has no referral chain. If an NGO refers a client to a government agency for housing support and the housing agency has its own system, there is no mechanism to track that referral's outcome. For organizations operating multi-partner programs, this is a significant gap.

### 6. One `category_id` per person is too narrow

Observer uses a single `category_id` FK on `people`. UNHCR's PSN code system supports multiple concurrent specific-needs codes per individual (e.g., a person can simultaneously be DS-MM — person with mobility impairment — and SM-MC — serious chronic medical condition — and SP-PT — single parent). Organizations working to UNHCR standards expect multi-code vulnerability classification.

The current model forces a choice of one vulnerability category. This means either:

- Creating composite categories ("single parent + disabled") which proliferate the `categories` table with combinatorial entries
- Recording only the "primary" vulnerability and losing secondary ones — common but analytically lossy

### 7. `people.project_id` prevents multi-project enrollment

A person is permanently tied to one project. This is inconsistent with how organizations operate: a person may receive legal services from one project run by the same NGO while also participating in a housing project, and both case managers need to see the person's full history.

proGres sidesteps this by being a single centralized system across all operations. Primero handles it through inter-agency referrals (a record effectively moves between organizational silos). Observer's current model requires a person to be duplicated if enrolled in multiple projects — a deduplication nightmare.

### 8. Consent tracking is absent

GDPR applies to humanitarian organizations processing personal data of EU residents and citizens. The Ukrainian context involves data subjects who are EU residents (those who fled to EU countries) and data processed by EU-based and EU-funded organizations. Primero explicitly supports consent tracking. proGres has explicit consent documentation requirements per UNHCR data protection standards. Observer has no `consent_given`, `consent_date`, or data sharing consent field.

### 9. Phone numbers as JSONB limits operational queries

`phone_numbers JSONB` handles multiple numbers elegantly for display but makes operational use difficult:

- "Find all people with phone number X" requires a GIN scan with JSON operators — still possible but complex
- Integration with SMS gateways (RapidPro, Twilio) typically requires a canonical primary phone number
- Deduplication by phone number is expensive
- No way to enforce format (E.164 standard) at the database level

Real systems that need to send verification SMS or call for follow-up contact typically store at least a `primary_phone VARCHAR(20)` separately, with `phone_numbers JSONB` as supplementary.

### 10. `migration_records` lacks movement causality

Observer tracks `from_place_id`, `destination_place_id`, `migration_date`, and `notes`. Real displacement tracking (especially DTM-aligned) also captures:

- **Reason for movement**: fled conflict, security concerns, access to services, return, economic
- **Mode of transport**: relevant for flow monitoring
- **Housing situation at destination**: IDP collective site, private, with relatives, renting

These dimensions are required for the "intentions and conditions" reporting that donors increasingly require.

---

## Proposed Improvements and Pivots

### Improvement 1 — Households as a first-class entity (High impact, moderate complexity)

Add a `households` table and a `household_members` junction table, replacing `people.parent_id`.

```sql
CREATE TABLE households (
    id               TEXT        PRIMARY KEY,
    project_id       TEXT        NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    reference_number TEXT,          -- human-readable case reference (e.g., KYV-2024-00142)
    head_person_id   TEXT        REFERENCES people (id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE household_members (
    household_id TEXT        NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    person_id    TEXT        NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    relationship TEXT        NOT NULL CHECK (relationship IN (
        'head', 'spouse', 'child', 'parent', 'sibling',
        'grandchild', 'grandparent', 'other_relative', 'non_relative'
    )),
    PRIMARY KEY (household_id, person_id)
);
```

`people.parent_id` would be removed. Group-level queries use `household_members` with a `WHERE relationship = 'head'` filter for identifying the registration head. The `reference_number` supports human-readable case file citations on paper documents.

**Trade-off:** Breaking change in the people schema; all existing family queries must be rewritten. The gain is substantial: correct peer relationships (spouse ≠ child), relationship semantics, and a group-level identity that can receive group-level documents.

---

### Improvement 2 — Enforce `external_id` uniqueness within a project (Low complexity, high value)

```sql
CREATE UNIQUE INDEX uq_people_project_external_id
    ON people (project_id, external_id)
    WHERE external_id IS NOT NULL;
```

Document that `external_id` should hold the РНОКПП (Ukrainian tax number) when available. This is the only practical deduplication anchor available without biometrics. The partial index (`WHERE external_id IS NOT NULL`) avoids unique constraint conflicts for records where the ID is not yet known.

---

### Improvement 3 — Replace `idp_status` enum with origin-place-derived classification (Moderate complexity)

Remove the hardcoded `idp_status CHECK (IN 'crimea', 'east_ukraine', 'non_idp')` column. Instead:

- Add a `conflict_zone TEXT` column to the `states` table to tag oblasts with their conflict zone designation
- The IDP classification is then computed: `origin_place_id → places.state_id → states.conflict_zone`

```sql
-- states table gains:
ALTER TABLE states ADD COLUMN conflict_zone TEXT
    CHECK (conflict_zone IN ('crimea', 'east_ukraine', 'south_east', NULL));
```

Advantages:

- Classification updates (new conflict zones) require a data update, not a schema migration
- Origin place remains queryable at the settlement level for geographic reports
- Consistent with how the official Ukraine register derives the classification (from prior address)

**Trade-off:** Requires seeding the `states` table with conflict zone data. Application queries become a join. Acceptable for this scale.

---

### Improvement 4 — Multi-code vulnerability classification (Low complexity)

Replace the single `people.category_id` FK with a junction table, mirroring the UNHCR PSN multi-code model.

```sql
CREATE TABLE person_categories (
    person_id   TEXT NOT NULL REFERENCES people (id) ON DELETE CASCADE,
    category_id TEXT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (person_id, category_id)
);
```

Remove `people.category_id`. This aligns with proGres PSN codes and allows recording all applicable vulnerability dimensions per person without a combinatorial explosion in the `categories` table.

---

### Improvement 5 — Lightweight referral tracking on support records (Moderate complexity)

Add referral lifecycle fields to `support_records` for outbound referrals to other organizations or offices.

```sql
ALTER TABLE support_records
    ADD COLUMN referral_status  TEXT CHECK (referral_status IN
        ('pending', 'accepted', 'completed', 'declined', 'no_response')),
    ADD COLUMN referred_to_office TEXT REFERENCES offices (id) ON DELETE SET NULL;
```

A support record with `referral_status IS NOT NULL` is a referral; one without is direct service delivery. This enables reporting on referral closure rates without a separate referral entity. A full inter-agency referral module (like Primero's) would require a separate `referrals` table — treat this as an interim solution.

---

### Improvement 6 — Add `primary_phone` alongside JSONB (Low complexity)

```sql
ALTER TABLE people
    ADD COLUMN primary_phone VARCHAR(20);  -- E.164 format recommended
```

`phone_numbers JSONB` retains additional numbers. `primary_phone` is the canonical contact number used for SMS integrations, deduplication queries, and display. This mirrors how every humanitarian system handles the "multiple phones, one primary" pattern.

---

### Improvement 7 — Consent tracking (Low complexity, compliance-critical)

```sql
ALTER TABLE people
    ADD COLUMN consent_given BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN consent_date  DATE;
```

Consent governs whether a person's data can be shared with other organizations via referrals or data exports. Default `FALSE` enforces opt-in. Required for GDPR compliance when processing data of persons in EU member states.

---

### Improvement 8 — Movement causality on migration records (Low complexity, domain alignment)

```sql
ALTER TABLE migration_records
    ADD COLUMN movement_reason TEXT CHECK (movement_reason IN (
        'conflict', 'security', 'service_access', 'return', 'relocation_program', 'economic', 'other'
    )),
    ADD COLUMN housing_at_destination TEXT CHECK (housing_at_destination IN (
        'own_property', 'renting', 'with_relatives', 'collective_site', 'hotel', 'other', 'unknown'
    ));
```

These two dimensions are the minimum required by DTM-aligned reporting and UNHCR return/displacement monitoring. `movement_reason` distinguishes voluntary relocation from forced displacement — a legally and operationally significant distinction.

---

## Summary Priority Table

| #   | Improvement                                                      | Complexity | Impact            | Schema change                                       |
| --- | ---------------------------------------------------------------- | ---------- | ----------------- | --------------------------------------------------- |
| 1   | Households as first-class entity + relationship types            | High       | High              | `households`, `household_members`, drop `parent_id` |
| 2   | Enforce `external_id` uniqueness per project                     | Low        | High              | Partial unique index                                |
| 3   | Replace `idp_status` enum with `states.conflict_zone` derivation | Moderate   | Medium            | `states.conflict_zone`, drop `people.idp_status`    |
| 4   | Multi-code vulnerability (`person_categories` junction)          | Low        | Medium            | `person_categories`, drop `people.category_id`      |
| 5   | Lightweight referral tracking on support records                 | Low        | Medium            | 2 columns on `support_records`                      |
| 6   | `primary_phone` canonical field                                  | Low        | Low               | 1 column on `people`                                |
| 7   | Consent tracking                                                 | Low        | High (compliance) | 2 columns on `people`                               |
| 8   | Movement causality fields                                        | Low        | Medium            | 2 columns on `migration_records`                    |

---

## Architectural Pivot: Person/Case Separation

The most significant structural divergence from professional systems (proGres, OSCaR) is that Observer conflates a person's identity record with their case record. This is also the most disruptive change to consider.

**Current:** `people` row = biographical identity + case lifecycle + project enrollment (all in one)

**Alternative:** separate `individuals` (biographical, cross-project) from `enrollments` or `cases` (project-specific, can be opened and closed independently)

This pivot would support:

- The same person receiving services across multiple projects without record duplication
- Multiple concurrent case states (housing case closed, legal case active)
- Full cross-project history for returning beneficiaries

**Recommendation:** Defer this pivot. The current user base (single-NGO, project-scoped operations) does not require cross-project enrollment. The complexity of a case/person split is significant and the benefit is only realized when multiple projects share beneficiaries. Monitor for this requirement as the platform matures; it is easier to add a `cases` table and migrate `people` to `individuals` as a forward migration when the need is demonstrated than to build the complexity upfront.

---

## References

- Primero: [github.com/primeroIMS/primero](https://github.com/primeroIMS/primero), [support.primero.org](https://support.primero.org)
- proGres / PRIMES: UNHCR Registration Guidance Chapter 7; OIOS Audit 2024/056
- UNHCR PSN Codes: UNHCR Guidance on Standardized Specific Needs Codes
- ActivityInfo: [activityinfo.org/support/docs](https://www.activityinfo.org/support/docs)
- KoBoToolbox: [support.kobotoolbox.org](https://support.kobotoolbox.org)
- IOM DTM Ukraine: [dtm.iom.int/ukraine](https://dtm.iom.int/ukraine)
- MiMOSA NextGen: IOM RFP December 2020
- Ukraine IDP Register: IDMC Analysis "IDP Registration in Ukraine: Who's in, who's out, and who's counting"; UNDP Ukraine press releases
- EU4Recovery / Case Manager Online Cabinet: EEAS delegation Ukraine
