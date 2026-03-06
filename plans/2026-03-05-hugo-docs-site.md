# Hugo Documentation Site Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace the flat `docs/` markdown files with a Hugo site using the Hextra theme, producing a proper user-facing documentation website.

**Architecture:** Hugo site lives at `docs/` root. Hextra theme added as a Hugo module. Existing markdown content migrated into Hugo content structure with front matter. Plans directory moved to project root since it's internal-only.

**Tech Stack:** Hugo (v0.157+), Hextra theme (Hugo module), Go modules for Hugo


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

title: Observer
layout: hextra-home

### Task 4: Create the docs section — Getting Started

**Files:**
- Create: `docs/content/docs/_index.md`
- Create: `docs/content/docs/getting-started.md`

**Step 1: Create docs section index**

```markdown
```

**Step 2: Create getting started page (from local-development.md)**

Migrate content from the old `docs/local-development.md`, adding Hugo front matter:

```markdown

(content of local-development.md with front matter added)
```

Copy the full content of `local-development.md` under the front matter.

**Step 3: Commit**

```bash
git add docs/content/docs/
git commit -m "add getting-started docs page"
```

title: Architecture
weight: 2

### Task 6: Create Frontend page

**Files:**
- Create: `docs/content/docs/frontend.md`

**Step 1: Create frontend page**

```markdown

(full content of the old frontend.md)
```

**Step 2: Commit**

```bash
git add docs/content/docs/frontend.md
git commit -m "add frontend docs page"
```

title: Testing
weight: 4

### Task 8: Create Reference section

**Files:**
- Create: `docs/content/docs/reference/_index.md`
- Create: `docs/content/docs/reference/variables.md`
- Create: `docs/content/docs/reference/kyrgyz-latin.md`

**Step 1: Create reference section index**

```markdown
```

**Step 2: Create variables page**

```markdown

(full content of the old variables.md)
```

**Step 3: Create kyrgyz latin page**

```markdown

(full content of the old qyrgyz-latin.md)
```

**Step 4: Commit**

```bash
git add docs/content/docs/reference/
git commit -m "add reference docs section (variables, kyrgyz latin)"
```

title: Architecture Decision Records
weight: 7
sidebar:
  open: false
title: "ADR-001: Bootstrapping"
weight: 1

### Task 10: Create About section

**Files:**
- Create: `docs/content/docs/about/_index.md`
- Create: `docs/content/docs/about/value-offering.md`
- Create: `docs/content/docs/about/existing-solutions.md`
- Create: `docs/content/docs/about/current-state.md`

**Step 1: Create about section index**

```markdown
```

**Step 2: Migrate value-offering.md**

```markdown

(full content of value-offering.md)
```

**Step 3: Migrate existing-solutions.md**

```markdown

(full content of existing-solutions.md)
```

**Step 4: Migrate current-state.md**

```markdown

(full content of current-state.md)
```

Replace any ```mermaid blocks with `{{< mermaid >}}` shortcodes.

**Step 5: Commit**

```bash
git add docs/content/docs/about/
git commit -m "add about docs section"
```


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
