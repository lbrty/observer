---
title: Why Observer
weight: 1
---

## The problem

Small and mid-size NGOs working with displaced persons in Ukraine face a gap: paper ledgers don't scale, but the systems that do — proGres, Primero, ActivityInfo — require UN partnerships, technical partners, or SaaS subscriptions these organizations can't access.

Observer fills that gap.

## What it is

A self-hosted case management platform. One Go binary, one PostgreSQL database. A single sysadmin can deploy it on a VPS.

It tracks people, their movement history, support records, households, documents, and pets — scoped per project, with role-based access control and built-in reporting.

## Who it's for

Organizations with:

- 5–50 staff, one or zero IT people
- One or more donor-funded aid projects
- Reporting obligations to EU, USAID, or bilateral donors
- Data sensitivity requirements that prevent using third-party SaaS
- No existing UN agency partnership granting access to Primero or proGres

## Why not the alternatives

| System | Barrier |
| --- | --- |
| **proGres v4** | Requires UNHCR partnership; not self-hostable |
| **Primero** | Requires UNICEF/IRC technical partner for deployment |
| **ActivityInfo** | SaaS with per-user pricing; designed for aggregate monitoring, not case management |
| **KoBoToolbox** | Data collection only — no persistent case records |

## What makes Observer different

**Deployable without permission.** No partner onboarding, no SaaS agreement. `docker compose up` and you're running.

**Schema-enforced access control.** Platform roles (admin, staff, consultant, guest) combined with project roles (owner, manager, consultant, viewer) and three data sensitivity flags. Not configurable per-grant — declarative and immediate.

**Reports that match donor requirements.** 39 report types covering consultation counts, IDP origin breakdowns, sex/age disaggregation, sphere of appeal breakdowns, and family unit counts. These are the specific reports Ukrainian legal aid organizations submit to donors.

**Support sphere taxonomy.** 11 constrained consultation topics (housing, documents, social benefits, property rights, etc.) enabling GROUP BY-safe breakdowns. No incumbent system offers this at the schema level.

**Households as first-class entities.** Head, spouse, child, parent, sibling — typed relationships that support family unit reports and group-level document attribution.

**IDP classification derived from geography.** No hardcoded political categories. IDP origin is computed from `origin_place → state → conflict_zone`. New conflict zones are added by updating reference data, not the schema.

**Forward-only migrations.** The database schema evolves through versioned, append-only SQL migrations. No rollback files, no manual patches.

## What Observer does not do

- Scale to millions of registrations (proGres territory)
- Inter-agency referral workflows (Primero territory)
- Biometric deduplication (proGres + BIMS)
- Cluster-level 5W/3W aggregate reporting (ActivityInfo territory)

Observer produces the raw data. If you need cluster dashboards, export to ActivityInfo.
