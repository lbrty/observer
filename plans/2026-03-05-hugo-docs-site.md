# Hugo Documentation Site Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace the flat `docs/` markdown files with a Hugo site using the Hextra theme, producing a proper user-facing documentation website.

**Architecture:** Hugo site lives at `docs/` root. Hextra theme added as a Hugo module. Existing markdown content migrated into Hugo content structure with front matter. Plans directory moved to project root since it's internal-only.

**Tech Stack:** Hugo (v0.157+), Hextra theme (Hugo module), Go modules for Hugo

---

## Pre-work: Move internal files out of docs/

Before initializing Hugo, relocate files that should NOT be user-facing.

### Task 1: Move plans and internal docs out of docs/

**Files:**
- Move: `docs/plans/` -> `plans/` (project root)
- Move: `docs/old-project-summary.md` -> delete or move to `plans/`
- Modify: `.gitignore` if needed

**Step 1: Move plans directory to project root**

```bash
mv docs/plans plans
```

**Step 2: Move old-project-summary.md**

```bash
mv docs/old-project-summary.md plans/
```

**Step 3: Commit**

```bash
git add -A plans/ docs/plans/ docs/old-project-summary.md
git commit -m "move internal plans out of docs/ before Hugo migration"
```

---

## Hugo Site Initialization

### Task 2: Initialize Hugo site in docs/

**Files:**
- Create: `docs/hugo.toml`
- Create: `docs/go.mod`
- Create: `docs/go.sum`
- Remove: all old flat markdown files from `docs/` root

**Step 1: Clear old flat docs (content will be re-added as Hugo pages)**

```bash
# Back up ADR folder content and all .md files
# They'll be recreated as Hugo content pages
rm docs/architecture.md docs/current-state.md docs/existing-solutions.md \
   docs/frontend.md docs/local-development.md docs/testing-guide.md \
   docs/value-offering.md docs/variables.md docs/qyrgyz-latin.md
rm -rf docs/adr/
```

**Step 2: Initialize Hugo module in docs/**

```bash
cd docs
hugo mod init github.com/lbrty/observer/docs
```

**Step 3: Create hugo.toml**

```toml
baseURL = "https://lbrty.github.io/observer/"
title = "Observer"
languageCode = "en"

[module]
  [[module.imports]]
    path = "github.com/imfing/hextra"

[markup.goldmark.renderer]
  unsafe = true

[markup.highlight]
  noClasses = false

[[menu.main]]
  name = "GitHub"
  url = "https://github.com/lbrty/observer"
  weight = 100
  [menu.main.params]
    icon = "github"

[params]
  description = "Self-hosted IDP case management for NGOs"

[params.navbar]
  displayTitle = true
  displayLogo = false

[params.footer]
  displayPoweredBy = false

[params.page]
  width = "wide"
```

**Step 4: Download the Hextra module**

```bash
cd docs && hugo mod get -u
```

**Step 5: Verify Hugo builds**

```bash
cd docs && hugo
```

Expected: successful build with no errors.

**Step 6: Commit**

```bash
git add docs/hugo.toml docs/go.mod docs/go.sum
git commit -m "initialize Hugo site with Hextra theme in docs/"
```

---

## Content Migration

### Task 3: Create the home page

**Files:**
- Create: `docs/content/_index.md`

**Step 1: Create home page**

The home page uses Hextra's landing page layout with hero and feature cards.

```markdown
---
title: Observer
layout: hextra-home
---

<div class="hx-mt-6 hx-mb-6">
{{< hextra/hero-headline >}}
  Self-hosted case management for NGOs
{{< /hextra/hero-headline >}}
</div>

<div class="hx-mb-12">
{{< hextra/hero-subtitle >}}
  A secure, organized way to manage displaced persons cases, track humanitarian support, and store sensitive personal information — replacing unprotected spreadsheets with a proper system.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx-mb-6">
{{< hextra/hero-button text="Get Started" link="docs/getting-started/" >}}
</div>

<div class="hx-mt-6"></div>

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Actually deployable"
    subtitle="A Go binary, a PostgreSQL database, and a Justfile. Deploy on a single VPS — no UN agency sponsorship required."
  >}}
  {{< hextra/feature-card
    title="Dual-level RBAC"
    subtitle="Platform roles combined with project roles and data sensitivity flags. Declarative, schema-enforced access control."
  >}}
  {{< hextra/feature-card
    title="39 report types"
    subtitle="Matching actual Ukrainian NGO donor reporting obligations — EU, USAID, and bilateral donor formats built in."
  >}}
  {{< hextra/feature-card
    title="Forward-only migrations"
    subtitle="Disciplined schema evolution with no rollback footguns. Built for systems that evolve over years."
  >}}
  {{< hextra/feature-card
    title="Consultation taxonomy"
    subtitle="11-value support_sphere enum enables GROUP BY-safe breakdown by consultation topic — unique in the field."
  >}}
  {{< hextra/feature-card
    title="Privacy-first"
    subtitle="Self-hosted, no third-party SaaS, GDPR consent tracking. Sensitive data stays where you control it."
  >}}
{{< /hextra/feature-grid >}}
```

**Step 2: Verify**

```bash
cd docs && hugo server -D
```

Open http://localhost:1313 — should see the landing page.

**Step 3: Commit**

```bash
git add docs/content/_index.md
git commit -m "add Hugo home page with Hextra hero and feature cards"
```

---

### Task 4: Create the docs section — Getting Started

**Files:**
- Create: `docs/content/docs/_index.md`
- Create: `docs/content/docs/getting-started.md`

**Step 1: Create docs section index**

```markdown
---
title: Documentation
---
```

**Step 2: Create getting started page (from local-development.md)**

Migrate content from the old `docs/local-development.md`, adding Hugo front matter:

```markdown
---
title: Getting Started
weight: 1
---

(content of local-development.md with front matter added)
```

Copy the full content of `local-development.md` under the front matter.

**Step 3: Commit**

```bash
git add docs/content/docs/
git commit -m "add getting-started docs page"
```

---

### Task 5: Create Architecture page

**Files:**
- Create: `docs/content/docs/architecture.md`

**Step 1: Create architecture page**

```markdown
---
title: Architecture
weight: 2
---

(full content of the old architecture.md)
```

Note: Mermaid diagrams are supported by Hextra natively. Wrap each mermaid block in the shortcode:

````
{{< mermaid >}}
graph TD
  ...
{{< /mermaid >}}
````

Replace all ```mermaid fenced blocks with `{{< mermaid >}}` shortcodes.

**Step 2: Commit**

```bash
git add docs/content/docs/architecture.md
git commit -m "add architecture docs page with mermaid diagrams"
```

---

### Task 6: Create Frontend page

**Files:**
- Create: `docs/content/docs/frontend.md`

**Step 1: Create frontend page**

```markdown
---
title: Frontend
weight: 3
---

(full content of the old frontend.md)
```

**Step 2: Commit**

```bash
git add docs/content/docs/frontend.md
git commit -m "add frontend docs page"
```

---

### Task 7: Create Testing page

**Files:**
- Create: `docs/content/docs/testing.md`

**Step 1: Create testing page**

```markdown
---
title: Testing
weight: 4
---

(full content of the old testing-guide.md)
```

Replace all ```mermaid blocks with `{{< mermaid >}}` shortcodes.

**Step 2: Commit**

```bash
git add docs/content/docs/testing.md
git commit -m "add testing docs page"
```

---

### Task 8: Create Reference section

**Files:**
- Create: `docs/content/docs/reference/_index.md`
- Create: `docs/content/docs/reference/variables.md`
- Create: `docs/content/docs/reference/kyrgyz-latin.md`

**Step 1: Create reference section index**

```markdown
---
title: Reference
weight: 5
---
```

**Step 2: Create variables page**

```markdown
---
title: Environment Variables
weight: 1
---

(full content of the old variables.md)
```

**Step 3: Create kyrgyz latin page**

```markdown
---
title: Kyrgyz Latin Transliteration
weight: 2
---

(full content of the old qyrgyz-latin.md)
```

**Step 4: Commit**

```bash
git add docs/content/docs/reference/
git commit -m "add reference docs section (variables, kyrgyz latin)"
```

---

### Task 9: Create ADR section

**Files:**
- Create: `docs/content/docs/adr/_index.md`
- Create: `docs/content/docs/adr/001-bootstrapping.md`
- Create: `docs/content/docs/adr/002-users-and-auth.md`
- Create: `docs/content/docs/adr/003-main-schema.md`
- Create: `docs/content/docs/adr/004-forward-only-migrations.md`
- Create: `docs/content/docs/adr/005-reports.md`
- Create: `docs/content/docs/adr/006-authorization-middleware.md`
- Create: `docs/content/docs/adr/007-color-themes.md`
- Create: `docs/content/docs/adr/008-tygo-type-generation.md`
- Create: `docs/content/docs/adr/009-deferred-features.md`

**Step 1: Create ADR section index**

```markdown
---
title: Architecture Decision Records
weight: 7
sidebar:
  open: false
---

These records document significant architectural decisions made during development.
```

**Step 2: Migrate each ADR file**

For each ADR file, add front matter with title and weight matching the ADR number:

```markdown
---
title: "ADR-001: Bootstrapping"
weight: 1
---

(original ADR content)
```

Replace any ```mermaid blocks with `{{< mermaid >}}` shortcodes.

**Step 3: Commit**

```bash
git add docs/content/docs/adr/
git commit -m "add ADR docs section"
```

---

### Task 10: Create About section

**Files:**
- Create: `docs/content/docs/about/_index.md`
- Create: `docs/content/docs/about/value-offering.md`
- Create: `docs/content/docs/about/existing-solutions.md`
- Create: `docs/content/docs/about/current-state.md`

**Step 1: Create about section index**

```markdown
---
title: About
weight: 6
---
```

**Step 2: Migrate value-offering.md**

```markdown
---
title: Value Offering
weight: 1
---

(full content of value-offering.md)
```

**Step 3: Migrate existing-solutions.md**

```markdown
---
title: Existing Solutions
weight: 2
---

(full content of existing-solutions.md)
```

**Step 4: Migrate current-state.md**

```markdown
---
title: Current State
weight: 3
---

(full content of current-state.md)
```

Replace any ```mermaid blocks with `{{< mermaid >}}` shortcodes.

**Step 5: Commit**

```bash
git add docs/content/docs/about/
git commit -m "add about docs section"
```

---

## Finishing Touches

### Task 11: Add Justfile commands for docs

**Files:**
- Modify: `Justfile`

**Step 1: Add docs commands to Justfile**

```just
# Start docs dev server
docs-dev:
    cd docs && hugo server -D

# Build docs site
docs-build:
    cd docs && hugo --minify

# Clean docs build
docs-clean:
    rm -rf docs/public/
```

**Step 2: Commit**

```bash
git add Justfile
git commit -m "add docs-dev, docs-build, docs-clean to Justfile"
```

---

### Task 12: Add .gitignore entries

**Files:**
- Modify: `.gitignore`

**Step 1: Add Hugo build output to .gitignore**

```
# Hugo
docs/public/
docs/resources/_gen/
docs/.hugo_build.lock
```

**Step 2: Commit**

```bash
git add .gitignore
git commit -m "add Hugo build artifacts to .gitignore"
```

---

### Task 13: Verify full site build

**Step 1: Build the site**

```bash
cd docs && hugo --minify
```

Expected: clean build, no errors or warnings.

**Step 2: Run dev server and check all pages**

```bash
cd docs && hugo server
```

Verify:
- Home page renders with hero + feature cards
- Sidebar navigation shows all sections
- All pages render content correctly
- Mermaid diagrams render
- Search works
- Dark mode toggle works

**Step 3: Final commit if any fixes needed**

---

## Content Structure Summary

```
docs/
  hugo.toml
  go.mod
  go.sum
  content/
    _index.md                          # landing page
    docs/
      _index.md                        # docs section root
      getting-started.md               # local development guide
      architecture.md                  # architecture overview + mermaid
      frontend.md                      # frontend guide
      testing.md                       # testing guide
      reference/
        _index.md
        variables.md                   # env vars reference
        kyrgyz-latin.md                # transliteration rules
      about/
        _index.md
        value-offering.md              # value proposition
        existing-solutions.md          # comparison matrix
        current-state.md               # current architecture state
      adr/
        _index.md
        001-bootstrapping.md
        002-users-and-auth.md
        003-main-schema.md
        004-forward-only-migrations.md
        005-reports.md
        006-authorization-middleware.md
        007-color-themes.md
        008-tygo-type-generation.md
        009-deferred-features.md
```
