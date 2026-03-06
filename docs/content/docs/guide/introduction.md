---
title: Introduction
weight: 1
---

Observer is a self-hosted IDP case management platform for NGOs.

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

**Schema-enforced access control.** Platform roles (admin, staff, consultant, guest) combined with project roles (owner, manager, consultant, viewer) and three data sensitivity flags. Declarative and immediate.

**Reports that match donor requirements.** 39 report types covering consultation counts, IDP origin breakdowns, sex/age disaggregation, sphere of appeal breakdowns, and family unit counts.

**Support sphere taxonomy.** 11 constrained consultation topics enabling GROUP BY-safe breakdowns.

**Households as first-class entities.** Typed relationships (head, spouse, child, parent, sibling) that support family unit reports.

**IDP classification derived from geography.** Origin computed from `origin_place → state → conflict_zone` — no hardcoded political categories.

**Forward-only migrations.** Versioned, append-only SQL migrations. No rollback files.

## What Observer does not do

- Scale to millions of registrations (proGres territory)
- Inter-agency referral workflows (Primero territory)
- Biometric deduplication (proGres + BIMS)
- Cluster-level 5W/3W aggregate reporting (ActivityInfo territory)

Observer produces the raw data. If you need cluster dashboards, export to ActivityInfo.
