
# UX/UI Polish Design

**Goal:** Make the Observer web app feel sleek, responsive, and effortless to use by fixing inconsistencies, adding missing primitives, and improving mobile responsiveness.

**Approach:** Bottom-up — build missing component primitives first, then apply them across pages. No new features, only polish.

**Scope:** Frontend only (`packages/observer-web/src/`).


## 2. Collapsible Sidebar

**Problem:** The w-52 sidebar is fixed width with no responsive behavior. On tablets and small laptops it eats too much horizontal space. On mobile it's unusable.

**Solution:**
- **Desktop (>=1024px):** Full sidebar, current behavior
- **Tablet (768-1023px):** Icon-only collapsed sidebar (w-14), tooltip on hover for labels
- **Mobile (<768px):** Hidden sidebar, hamburger menu button in header opens it as a slide-over overlay

Add a `SidebarContext` to manage open/collapsed state. Store preference in localStorage.


## 4. Toast Notifications

**Problem:** Success/error feedback uses inline banners (`SuccessBanner`, `ErrorBanner`) that are positioned inside forms and auto-clear. They're invisible if the user scrolls, and can't show feedback for actions taken outside forms (e.g. delete confirmation).

**Solution:** Add a global toast system using Base-UI's `Toast` component (or a lightweight custom one):
- Position: bottom-right
- Auto-dismiss: 4s for success, sticky for errors
- Variants: success (foam), error (rose), info (accent)
- Accessible: `role="status"` with `aria-live="polite"`
- Zustand store: `useToast()` with `toast.success(msg)`, `toast.error(msg)`

Keep inline banners for form-level validation errors, use toasts for action confirmations.


## 6. Better Empty States

**Problem:** Empty tables show a plain text string ("No data") with no visual weight. Empty states for person tabs (no notes, no documents) are similarly bare.

**Solution:** Create `components/empty-state.tsx`:
- Centered layout with a muted icon (from Phosphor), heading, description, optional action button
- Reuse across DataTable and tab content pages
- Each entity gets a contextual empty message (e.g. "No notes yet — add the first note")


## 8. Tooltip Component

**Problem:** Icon-only buttons (edit, delete, close) have no labels. Users must guess what each icon does. Some have `title` attributes but no styled tooltip.

**Solution:** Create `components/tooltip.tsx` wrapping Base-UI's `Tooltip`:
- Simple API: `<Tooltip label="Edit"><button>...</button></Tooltip>`
- Styled: dark bg, small text, subtle shadow, appears on hover/focus with 200ms delay
- Apply to all icon-only buttons: edit pencils, delete trash cans, close X, pagination arrows


## 10. Micro-interactions & Transitions

**Problem:** Most UI state changes are instant. Drawer backdrop and popup animate, but tab switches, page loads, and button presses lack tactile feedback.

**Solution:**
- Buttons: `active:scale-[0.98]` for a subtle press effect
- Tab switching: content fade-in (opacity 0→1, 150ms)
- Sidebar link active state: smooth background color transition (already has `transition-colors`)
- Loading → content: fade in with `animate-in fade-in` (Tailwind animate plugin or custom keyframe)
- Skeleton → real content: crossfade rather than hard swap

