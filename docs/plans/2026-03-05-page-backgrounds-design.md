# Bold Page Backgrounds — Design

**Goal:** Give every page a distinctive visual identity using large-scale decorative SVG patterns, extending the existing topographic design language.

**Architecture:** Pure CSS using mask-image SVG data URIs (same technique as existing `card-bg-*` classes). No JS, no extra assets.

---

## Two layers

### 1. Base topo layer (every page)

Large-scale topographic contour pattern applied via `::before` on the content wrapper in `_app.tsx`. Positioned in the top area of the page, fading downward. Opacity ~5%. Uses `var(--fg)` as background-color so it adapts to all 4 themes automatically.

### 2. Accent illustration layer (per page zone)

Thematic SVG watermark positioned top area, offset ~15% from the right edge. Opacity 6-10%. Applied as a CSS class on each page's outer `<div>`.

| Zone | Motif | CSS class |
|------|-------|-----------|
| Dashboard | Compass rose / cardinal directions | `.page-bg-dashboard` |
| People | Abstract overlapping circles (community) | `.page-bg-people` |
| Reports | Rising bar/line chart silhouette | `.page-bg-reports` |
| Admin | Hexagonal grid (structure) | `.page-bg-admin` |
| Support records | Interlocking chain links | `.page-bg-support` |
| Profile | Single circle with radiating lines | `.page-bg-profile` |
| Reference data | Globe with latitude lines | `.page-bg-reference` |
| Households/Tags/Pets/Docs | Base topo only (no accent) | — |

## Technical details

- All patterns use inline SVG in CSS `mask-image` (data URI encoded)
- `background-color: var(--fg)` with low opacity = automatic theme support
- `pointer-events: none` on all pseudo-elements
- `position: absolute` within `overflow: hidden` containers
- Base layer lives in `_app.tsx` content wrapper (`::before`)
- Accent classes applied per-page on the outermost content `<div>`

## Scope exclusions

- No animation or parallax
- No per-theme pattern variations
- No changes to existing card-level patterns (`card-bg-topo`, `card-bg-dots`, `card-bg-waves`)
