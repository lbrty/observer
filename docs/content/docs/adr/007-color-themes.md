---
title: "ADR-007: Color Themes"
weight: 7
---

| Field      | Value                   |
| ---------- | ----------------------- |
| Status     | Accepted                |
| Date       | 2026-02-23              |
| Components | observer-web (main.css) |

---

## Context

Observer is case management infrastructure for humanitarian case workers operating in NGO offices and field locations — low-light environments, direct sunlight, degraded displays. The color system must:

1. Communicate case states (verified, pending, attention-required) with zero ambiguity.
2. Work across four lighting conditions without losing semantic clarity.
3. Avoid cool tones that feel clinical — the palette should feel warm and deliberate, not sterile.
4. Meet WCAG AAA (7.0:1) for body text and AA18 (3.0:1) for interactive elements across all themes.

---

## Decision

### Warm-shifted palette

All colors are warm-shifted at every luminance level. No cool grays. Backgrounds use warm charcoal (dark) and warm parchment (light), never slate or pure white/black.

### Four theme variants

| Theme    | Selector                                             | Purpose                     | Activation               |
| -------- | ---------------------------------------------------- | --------------------------- | ------------------------ |
| Light    | `:root` (default CSS)                                | Office environments         | Explicit user choice     |
| Dark     | `[data-theme="dark"]` / `prefers-color-scheme: dark` | Field conditions, low-light | Default for new installs |
| Light HC | `[data-theme="light-hc"]`                            | Direct sunlight             | Explicit opt-in          |
| Dark HC  | `[data-theme="dark-hc"]`                             | Hostile/degraded displays   | Explicit opt-in          |

Dark is the default. Platform `prefers-color-scheme: dark` maps to Dark, not Dark HC. High Contrast requires explicit opt-in. Never auto-switch based on time of day.

### Semantic color groups

Three accent groups carry fixed meanings across the entire UI:

| Group   | Meaning                      | Usage                                                                                   |
| ------- | ---------------------------- | --------------------------------------------------------------------------------------- |
| Sage    | Verified, complete           | Verified registrations, completed support records. Never before verification completes. |
| Amber   | Active, interactive, pending | Active cases, pending referrals, focus rings, CTAs, interactive elements.               |
| Sienna  | Danger, immediate attention  | Conflict zones, data sensitivity warnings, destructive actions, urgent alerts.          |
| Neutral | Structural UI                | Surfaces, borders, text hierarchy. No semantic meaning.                                 |

These assignments are non-negotiable. Sage is never decorative. Sienna is never a mild warning. Amber is never an error state.

### Token architecture

Tokens are defined as CSS custom properties in `packages/observer-web/src/main.css` and mapped to Tailwind via `@theme inline`.

**Structural tokens** (per-theme values):

```
--bg, --bg-secondary, --bg-tertiary     Surfaces (ground, raised, sunken)
--fg, --fg-secondary, --fg-tertiary     Text hierarchy
--border, --border-secondary            Structural borders
--accent, --accent-fg                   Primary interactive color + text on accent
--ring                                  Focus ring (semi-transparent accent)
--fg-error                              Error text
```

**Semantic extras** (per-theme values):

```
--color-foam    Sage group — verified/complete states
--color-gold    Amber group — active/pending states
--color-rose    Sienna group — danger/attention
```

**Shadows** (per-theme values):

```
--shadow-card, --shadow-elevated, --shadow-overlay
```

### Accent scale shifting

Accent values shift across themes to preserve visual weight:

| Token      | Light     | Dark      | Light HC  | Dark HC   |
| ---------- | --------- | --------- | --------- | --------- |
| `--accent` | `#b07825` | `#c9953a` | `#885010` | `#e0a840` |

One stop lighter on Dark, two stops lighter on Dark HC, one stop darker on Light HC.

### Current CSS snapshot

Source: `packages/observer-web/src/main.css`

#### Light (`:root` default)

| Token                | Value                      |
| -------------------- | -------------------------- |
| `--bg`               | `#f2f0eb`                  |
| `--bg-secondary`     | `#fdfcfa`                  |
| `--bg-tertiary`      | `#eae8e2`                  |
| `--fg`               | `#1a1916`                  |
| `--fg-secondary`     | `#3d3a32`                  |
| `--fg-tertiary`      | `#7a7568`                  |
| `--border`           | `#b8b4a9`                  |
| `--border-secondary` | `#d6d2c8`                  |
| `--accent`           | `#b07825`                  |
| `--accent-fg`        | `#fdfcfa`                  |
| `--ring`             | `rgba(176, 120, 37, 0.12)` |
| `--fg-error`         | `#a33b2c`                  |
| `--color-foam`       | `#668c57`                  |
| `--color-gold`       | `#c9953a`                  |
| `--color-rose`       | `#c05040`                  |
| `--shadow-card`      | `rgba(26, 25, 22, 0.07)`   |
| `--shadow-elevated`  | `rgba(26, 25, 22, 0.15)`   |
| `--shadow-overlay`   | `rgba(26, 25, 22, 0.04)`   |

#### Dark (`[data-theme="dark"]` / `prefers-color-scheme: dark`)

| Token                | Value                      |
| -------------------- | -------------------------- |
| `--bg`               | `#191714`                  |
| `--bg-secondary`     | `#201e1a`                  |
| `--bg-tertiary`      | `#141210`                  |
| `--fg`               | `#f0ede6`                  |
| `--fg-secondary`     | `#a09890`                  |
| `--fg-tertiary`      | `#6b6560`                  |
| `--border`           | `#3d3a35`                  |
| `--border-secondary` | `#2e2b27`                  |
| `--accent`           | `#c9953a`                  |
| `--accent-fg`        | `#191714`                  |
| `--ring`             | `rgba(201, 149, 58, 0.12)` |
| `--fg-error`         | `#c05040`                  |
| `--color-foam`       | `#668c57`                  |
| `--color-gold`       | `#deb96a`                  |
| `--color-rose`       | `#d6816a`                  |
| `--shadow-card`      | `rgba(0, 0, 0, 0.3)`       |
| `--shadow-elevated`  | `rgba(0, 0, 0, 0.52)`      |
| `--shadow-overlay`   | `rgba(0, 0, 0, 0.08)`      |

#### Dark HC (`[data-theme="dark-hc"]`)

| Token                | Value                      |
| -------------------- | -------------------------- |
| `--bg`               | `#0d0c0a`                  |
| `--bg-secondary`     | `#141210`                  |
| `--bg-tertiary`      | `#080706`                  |
| `--fg`               | `#ffffff`                  |
| `--fg-secondary`     | `#d6d0c6`                  |
| `--fg-tertiary`      | `#9a9288`                  |
| `--border`           | `#4f4a43`                  |
| `--border-secondary` | `#302c27`                  |
| `--accent`           | `#e0a840`                  |
| `--accent-fg`        | `#0d0c0a`                  |
| `--ring`             | `rgba(224, 168, 64, 0.18)` |
| `--fg-error`         | `#d8604e`                  |
| `--color-foam`       | `#7aad68`                  |
| `--color-gold`       | `#f0c870`                  |
| `--color-rose`       | `#ee8e7c`                  |
| `--shadow-card`      | `rgba(0, 0, 0, 0.5)`       |
| `--shadow-elevated`  | `rgba(0, 0, 0, 0.76)`      |
| `--shadow-overlay`   | `rgba(0, 0, 0, 0.04)`      |

#### Light HC (`[data-theme="light-hc"]`)

| Token                | Value                     |
| -------------------- | ------------------------- |
| `--bg`               | `#ede9e2`                 |
| `--bg-secondary`     | `#f5f3ed`                 |
| `--bg-tertiary`      | `#d9d5cc`                 |
| `--fg`               | `#0c0b09`                 |
| `--fg-secondary`     | `#2a2820`                 |
| `--fg-tertiary`      | `#5c5849`                 |
| `--border`           | `#9e9a90`                 |
| `--border-secondary` | `#c0bcb2`                 |
| `--accent`           | `#885010`                 |
| `--accent-fg`        | `#ffffff`                 |
| `--ring`             | `rgba(136, 80, 16, 0.25)` |
| `--fg-error`         | `#7a2015`                 |
| `--color-foam`       | `#36682a`                 |
| `--color-gold`       | `#c9953a`                 |
| `--color-rose`       | `#7a2015`                 |
| `--shadow-card`      | `rgba(26, 25, 22, 0.12)`  |
| `--shadow-elevated`  | `rgba(26, 25, 22, 0.3)`   |
| `--shadow-overlay`   | `rgba(26, 25, 22, 0.08)`  |

### Backend representation

Theme preference is stored in the browser via `localStorage` under the key `observer-theme`. The frontend reads this value on load and sets `data-theme` on the root `<html>` element. Valid values: `"dark"`, `"light"`, `"dark-hc"`, `"light-hc"`, `"system"`. `"system"` defers to `prefers-color-scheme`.

No backend persistence is needed — theme is a per-device preference, not a per-account setting.

---

## Consequences

### Positive

- **Unambiguous state communication**: sage/amber/sienna carry fixed meanings — case workers never have to guess what a color means.
- **Field-ready**: four themes cover the full range of working conditions from NGO offices to field locations in direct sunlight.
- **Warm palette**: avoids the clinical feel of cool grays, appropriate for sensitive humanitarian case work with vulnerable populations.
- **Token indirection**: components reference tokens, not hex values — theme changes propagate automatically.

### Negative

- **Four theme variants to maintain**: every new color token requires four hex values. Mitigated by keeping the token count small.
- **localStorage-only persistence**: theme preference does not roam across devices. Acceptable trade-off — field workers typically use a single assigned device.

---

## Alternatives Considered

### A. Cool-neutral palette

Standard cool grays (slate, zinc) used by most design systems.

**Rejected because**: feels clinical and impersonal. Humanitarian case work benefits from a warm, grounded aesthetic that signals care and permanence when handling sensitive beneficiary data.

### B. Single theme with system preference only

Ship one light and one dark variant, no HC.

**Rejected because**: field workers operate in extreme conditions (direct sunlight, degraded displays) where standard contrast ratios are insufficient. HC variants are a field requirement, not a nice-to-have.

### C. Semantic colors derived from a single hue

Generate sage/sienna from the amber accent via hue rotation.

**Rejected because**: semantic colors must be independently tunable per theme to hit contrast targets. Derived colors would compromise either the light or dark HC variant.
