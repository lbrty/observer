---
title: Introduction
weight: 1
---

Observer is a self-hosted case management platform for NGOs working with internally displaced persons.

## The problem

Small and mid-size NGOs working with displaced persons in Ukraine face a gap: paper ledgers don't scale, but the systems that do — proGres, Primero, ActivityInfo — require UN partnerships, technical partners, or SaaS subscriptions these organizations can't access.

Observer fills that gap.

## What it does

Observer gives your organization a private, secure system to track the people you serve. It runs on your own server — no cloud service, no subscription, no third party ever sees your data.

With Observer, your team can:

- **Register people and families** — record personal details, household relationships, and documents
- **Track support** — log consultations, referrals, and the type of assistance provided
- **Follow movement** — where people came from, where they moved, and why
- **Control access** — decide who on your team can see what, down to contact details and documents
- **Generate reports** — built-in breakdowns by sex, age, region, support type, vulnerability category, and more — filterable to match EU, USAID, and bilateral donor requirements

One person with basic server skills can set it up in under an hour.

## Who it's for

Organizations with:

- 5 to 50 staff, one or zero IT people
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

**You don't need anyone's permission.** No partner onboarding, no SaaS agreement. Install it on your server and start working.

**Your data stays yours.** Everything runs on infrastructure you control. No data leaves your server unless you export it.

**Access control is built in.** Platform roles (admin, staff, consultant, guest) combine with project roles (owner, manager, consultant, viewer) and sensitivity flags that control who sees contact info, personal details, and documents.

**Reports match what donors actually ask for.** 12 report dimensions — consultation counts, IDP origin, sex/age disaggregation, support sphere, vulnerability category, region, office, tags, family units, and case status — each filterable by date range, support type, and demographic criteria.

**Families are tracked as units.** Typed household relationships (head, spouse, child, parent, sibling) power family-level reporting.

**IDP status is computed, not guessed.** A person's displacement status is derived from their origin location and whether that area is a conflict zone — no manual classification needed.

## Supported languages

The UI ships with six languages: English, Ukrainian, Russian, German, Turkish, and Kyrgyz (Latin script). Kyrgyz uses a custom Latin transliteration because the official Kyrgyz Latin alphabet was adopted in 2023 and standard translation tools don't support it yet — we maintain our own transliteration rules to provide accurate, native-feeling text for Central Asian deployments.

## What Observer does not do

Observer is not designed for:

- Inter-agency referral workflows (Primero territory)
- Biometric deduplication
- Cluster-level 5W/3W aggregate reporting (ActivityInfo territory)

If you need cluster dashboards, export your data from Observer to ActivityInfo.
