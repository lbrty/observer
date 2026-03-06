---
title: "UX/UI Polish Design"
weight: 1
---

# UX/UI Polish Design

**Goal:** Make the Observer web app feel sleek, responsive, and effortless to use by fixing inconsistencies, adding missing primitives, and improving mobile responsiveness.

**Approach:** Bottom-up — build missing component primitives first, then apply them across pages. No new features, only polish.

**Scope:** Frontend only (`packages/observer-web/src/`).

---

## 1. Button Component

**Problem:** Button styles are copy-pasted as inline className strings across 20+ locations. Three variants exist (primary/secondary/destructive) but with no shared component — leading to inconsistent padding, hover states, and disabled styling.

**Solution:** Create `components/button.tsx` with variants:
- `primary` — `bg-accent text-accent-fg` (create, save, submit)
- `secondary` — `border border-border-secondary` (cancel, back)
- `ghost` — transparent, `hover:bg-bg-tertiary` (icon buttons, toolbar)
- `danger` — `bg-rose text-white` (delete, destructive)

Props: `variant`, `size` (sm/md), `loading`, `disabled`, `icon` (leading icon slot), `asChild` (for Link composition). All buttons get `focus-visible:ring-2 focus-visible:ring-ring` for keyboard users.

---

## 2. Collapsible Sidebar

**Problem:** The w-52 sidebar is fixed width with no responsive behavior. On tablets and small laptops it eats too much horizontal space. On mobile it's unusable.

**Solution:**
- **Desktop (>=1024px):** Full sidebar, current behavior
- **Tablet (768-1023px):** Icon-only collapsed sidebar (w-14), tooltip on hover for labels
- **Mobile (<768px):** Hidden sidebar, hamburger menu button in header opens it as a slide-over overlay

Add a `SidebarContext` to manage open/collapsed state. Store preference in localStorage.

---

## 3. Focus-Visible Rings

**Problem:** Inputs use only `focus:border-accent` with no outline. Buttons have no focus indicator at all. Keyboard users can't see where focus is.

**Solution:** Add global focus-visible styles in `main.css`:
```css
:focus-visible {
  outline: 2px solid var(--ring);
  outline-offset: 2px;
}
```

Remove `outline-none` from inputs, replace with `outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1`. This preserves the clean look for mouse users while giving keyboard users clear focus indicators.

---

## 4. Toast Notifications

**Problem:** Success/error feedback uses inline banners (`SuccessBanner`, `ErrorBanner`) that are positioned inside forms and auto-clear. They're invisible if the user scrolls, and can't show feedback for actions taken outside forms (e.g. delete confirmation).

**Solution:** Add a global toast system using Base-UI's `Toast` component (or a lightweight custom one):
- Position: bottom-right
- Auto-dismiss: 4s for success, sticky for errors
- Variants: success (foam), error (rose), info (accent)
- Accessible: `role="status"` with `aria-live="polite"`
- Zustand store: `useToast()` with `toast.success(msg)`, `toast.error(msg)`

Keep inline banners for form-level validation errors, use toasts for action confirmations.

---

## 5. Responsive Drawer

**Problem:** `DrawerShell` uses `max-w-[840px]` fixed width. On mobile it covers the whole screen but doesn't adapt — no full-screen mode, no safe-area padding.

**Solution:**
- Mobile (<768px): Full-screen drawer with no border-radius, proper safe-area insets
- Desktop: Current right-side slide-in behavior
- Add `size` prop: `md` (560px, simple forms) and `lg` (840px, complex forms like person edit)

---

## 6. Better Empty States

**Problem:** Empty tables show a plain text string ("No data") with no visual weight. Empty states for person tabs (no notes, no documents) are similarly bare.

**Solution:** Create `components/empty-state.tsx`:
- Centered layout with a muted icon (from Phosphor), heading, description, optional action button
- Reuse across DataTable and tab content pages
- Each entity gets a contextual empty message (e.g. "No notes yet — add the first note")

---

## 7. Inline Form Validation

**Problem:** Form validation is top-level only (a single ErrorBanner). Individual fields show no per-field error state. Required fields show `*` but no color/styling on missing input.

**Solution:** Enhance `FormField` to accept an `error` prop:
- Show error text below the input in rose/red
- Apply `border-rose` to the input when error is present
- Base-UI `Field.Error` already supports this pattern
- Keep form-level ErrorBanner for server errors

---

## 8. Tooltip Component

**Problem:** Icon-only buttons (edit, delete, close) have no labels. Users must guess what each icon does. Some have `title` attributes but no styled tooltip.

**Solution:** Create `components/tooltip.tsx` wrapping Base-UI's `Tooltip`:
- Simple API: `<Tooltip label="Edit"><button>...</button></Tooltip>`
- Styled: dark bg, small text, subtle shadow, appears on hover/focus with 200ms delay
- Apply to all icon-only buttons: edit pencils, delete trash cans, close X, pagination arrows

---

## 9. Improved Table UX

**Problem:** DataTable works but lacks sorting indicators and has no responsive strategy for narrow screens.

**Solution:**
- Add optional `sortable` flag to columns, render sort arrows in header
- Accept `sortKey` + `sortDir` + `onSort` props
- On narrow screens (<768px): switch to a card/list layout instead of a table (using a `renderCard` prop or CSS-only responsive approach)
- Better skeleton that matches actual content width patterns

---

## 10. Micro-interactions & Transitions

**Problem:** Most UI state changes are instant. Drawer backdrop and popup animate, but tab switches, page loads, and button presses lack tactile feedback.

**Solution:**
- Buttons: `active:scale-[0.98]` for a subtle press effect
- Tab switching: content fade-in (opacity 0→1, 150ms)
- Sidebar link active state: smooth background color transition (already has `transition-colors`)
- Loading → content: fade in with `animate-in fade-in` (Tailwind animate plugin or custom keyframe)
- Skeleton → real content: crossfade rather than hard swap

---

## Priority Order

| # | Change | Impact | Effort |
|---|--------|--------|--------|
| 1 | Button component | High — fixes 20+ inconsistencies | Small |
| 2 | Focus-visible rings | High — accessibility compliance | Tiny |
| 3 | Collapsible sidebar | High — mobile usability | Medium |
| 4 | Toast notifications | Medium — better action feedback | Medium |
| 5 | Responsive drawer | Medium — mobile form UX | Small |
| 6 | Tooltip component | Medium — discoverability | Small |
| 7 | Better empty states | Medium — perceived polish | Small |
| 8 | Inline form validation | Medium — form usability | Small |
| 9 | Improved table UX | Medium — data-heavy pages | Medium |
| 10 | Micro-interactions | Low — perceived quality | Tiny |
