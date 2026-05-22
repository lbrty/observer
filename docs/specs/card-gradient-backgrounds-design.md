# Card Gradient Backgrounds

**Date:** 2026-05-14

## Goal

Replace the existing SVG-mask card background decorations (`card-bg-topo`, `card-bg-dots`, `card-bg-waves`) with a set of 8 gradient variants (opacities 0.035–0.07; see Opacity Calibration) that are randomly (but deterministically) assigned to cards. Gradients adapt between light and dark mode — dark mode uses `mix-blend-mode: screen` to keep colours from reading muddy on dark surfaces.

## Gradient Variants

8 named variants, each a single CSS class:

| Class              | Name   | Shape                         |
| ------------------ | ------ | ----------------------------- |
| `card-grad-indigo` | Indigo | radial, top-right corner      |
| `card-grad-foam`   | Foam   | diagonal linear, bottom-right |
| `card-grad-gold`   | Gold   | radial, bottom-right corner   |
| `card-grad-rose`   | Rose   | radial, top-left corner       |
| `card-grad-sky`    | Sky    | linear, top edge              |
| `card-grad-violet` | Violet | radial, top-right corner      |
| `card-grad-amber`  | Amber  | diagonal linear, bottom-right |
| `card-grad-teal`   | Teal   | radial, bottom-left corner    |

## Opacity Calibration

| Mode  | Blend                    | Opacity range |
| ----- | ------------------------ | ------------- |
| Light | normal                   | 0.05 – 0.07   |
| Dark  | `mix-blend-mode: screen` | 0.035 – 0.04  |

## Implementation

### 1. CSS — `src/main.css`

- Remove `.card-bg-topo`, `.card-bg-dots`, `.card-bg-waves`, and `.auth-backdrop` blocks.
- Add 8 `.card-grad-*` classes using `::after` pseudo-element (same pattern as existing classes: `position: absolute; inset: 0; pointer-events: none; border-radius: inherit`).
- Add `[data-theme="dark"] .card-grad-*::after { mix-blend-mode: screen; }` and `@media (prefers-color-scheme: dark) { :root:not([data-theme]) .card-grad-*::after { mix-blend-mode: screen; } }` overrides.

### 2. Utility — `src/lib/card-gradient.ts`

```ts
const CARD_GRADIENTS = [
  "card-grad-indigo",
  "card-grad-foam",
  "card-grad-gold",
  "card-grad-rose",
  "card-grad-sky",
  "card-grad-violet",
  "card-grad-amber",
  "card-grad-teal",
] as const;

export type CardGradient = (typeof CARD_GRADIENTS)[number];

export function cardGradient(index: number): CardGradient {
  return CARD_GRADIENTS[index % CARD_GRADIENTS.length];
}
```

Using `index % 8` keeps assignment deterministic and stable across renders.

### 3. Apply to card components

Replace existing `card-bg-*` class with `cardGradient(index)` call in each location:

| File                                         | Change                                           |
| -------------------------------------------- | ------------------------------------------------ |
| `routes/_app/index.tsx`                      | Stat action cards (topo) + project cards (waves) |
| `components/reports/shared/kpi-card.tsx`     | Pass `index` prop, apply gradient class          |
| `components/my-stats/my-stats-kpi-cards.tsx` | Pass index from map                              |
| `components/reports/people-kpi-cards.tsx`    | Pass index from map                              |
| `components/reports/pets-kpi-cards.tsx`      | Pass index from map                              |

### 4. Remove orphaned CSS

After all usages are updated, delete the `.card-bg-dots` class (currently defined but unused according to the audit).

## What is NOT changing

- Card layout, padding, border, shadow — untouched.
- `auth-backdrop` — kept as-is (full-page backdrop, different purpose).
- Any non-card components.
