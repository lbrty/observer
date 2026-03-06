
# UX/UI Polish Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make the Observer web app feel sleek, responsive, and effortless by adding missing component primitives and applying them consistently across all pages.

**Architecture:** Bottom-up — build shared components first (Button, Tooltip, EmptyState, Toast), then apply them across the app. Collapsible sidebar and responsive drawer improve mobile UX. Focus-visible rings and micro-interactions add polish. All work is in `packages/observer-web/src/`.

**Tech Stack:** React 19, Tailwind CSS v4, Base-UI (headless), Phosphor Icons, TanStack Router, i18next

**Design doc:** `docs/plans/2026-03-05-ux-polish-design.md`


### Task 2: Replace Primary Buttons with Button Component

Replace all inline primary button patterns across the codebase. There are ~26 occurrences. The pattern to find:

```
bg-accent px-[45] py-2 text-sm font-medium text-accent-fg
```

**Files to modify (all in `packages/observer-web/src/`):**
- `components/drawer-shell.tsx` — submit button
- `components/form-dialog.tsx` — submit button
- `components/add-reference-dialog.tsx` — submit button
- `components/confirm-dialog.tsx` — cancel (secondary) + delete (danger)
- `routes/_app/projects/$projectId/people/index.tsx` — register button
- `routes/_app/projects/$projectId/support-records/index.tsx` — add button
- `routes/_app/projects/$projectId/households/index.tsx` — add button
- `routes/_app/projects/$projectId/tags/index.tsx` — add button
- `routes/_app/projects/$projectId/pets/index.tsx` — add button
- `routes/_app/projects/$projectId/people/$personId/index.tsx` — edit button
- `routes/_app/projects/$projectId/people/$personId/notes.tsx` — add button
- `routes/_app/projects/$projectId/people/$personId/support-records.tsx` — add button
- `routes/_app/admin/users/index.tsx` — add button
- `routes/_app/admin/users/$userId.tsx` — save button
- `routes/_app/admin/projects/index.tsx` — create button
- `routes/_app/admin/projects/$projectId/index.tsx` — save button
- `routes/_app/admin/projects/$projectId/permissions.tsx` — save/add buttons
- `routes/_app/admin/reference/countries/index.tsx` — add button
- `routes/_app/admin/reference/countries/$countryId/index.tsx` — save button
- `routes/_app/admin/reference/countries/$countryId/states/$stateId.tsx` — save button
- `routes/_app/admin/reference/categories.tsx` — add button
- `routes/_app/admin/reference/offices.tsx` — add button
- `routes/_app/profile.tsx` — save buttons
- `routes/_auth/login.tsx` — submit button
- `routes/_auth/register.tsx` — submit button
- `components/household-drawer.tsx` — add member button

**Step 1: Update each file**

For each file, add `import { Button } from "@/components/button"` and replace inline buttons.

Example — **PageHeader action buttons** (people/index.tsx):

Before:
```tsx
<button
  type="button"
  onClick={openCreate}
  className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
>
  + {t("project.people.register")}
</button>
```

After:
```tsx
<Button onClick={openCreate} icon={<PlusIcon size={16} />}>
  {t("project.people.register")}
</Button>
```

Example — **DrawerShell** footer buttons:

Before:
```tsx
<Drawer.Close className="cursor-pointer rounded-lg border border-border-secondary px-4 py-2 text-sm font-medium text-fg-secondary shadow-card hover:bg-bg-tertiary">
  {t("admin.common.cancel")}
</Drawer.Close>
<button
  type="submit"
  disabled={isPending}
  className="cursor-pointer rounded-lg bg-accent px-5 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
>
  {isPending ? savingText : saveText}
</button>
```

After:
```tsx
<Button variant="secondary" asChild>
  <Drawer.Close>{t("admin.common.cancel")}</Drawer.Close>
</Button>
<Button type="submit" loading={isPending}>
  {isPending ? savingText : saveText}
</Button>
```

Example — **ConfirmDialog** (secondary cancel + danger delete):

```tsx
<Button variant="secondary" asChild>
  <Dialog.Close>{t("admin.common.cancel")}</Dialog.Close>
</Button>
<Button variant="danger" onClick={onConfirm} loading={loading}>
  {t("admin.common.delete")}
</Button>
```

Example — **RowActions** (ghost icon buttons):

```tsx
<Button variant="ghost" onClick={(e) => { e.stopPropagation(); onEdit(); }}>
  <PencilSimpleIcon size={16} />
</Button>
<Button variant="ghost" className="hover:text-rose" onClick={(e) => { e.stopPropagation(); onDelete(); }}>
  <TrashIcon size={16} />
</Button>
```

Example — **Login/Register** (full-width submit):

```tsx
<Button type="submit" loading={submitting} className="w-full">
  {submitting ? t("auth.loggingIn") : t("auth.login")}
</Button>
```

**Step 2: Verify the app compiles**

Run: `cd packages/observer-web && bunx tsc --noEmit --pretty 2>&1 | head -30`

**Step 3: Verify tests pass**

Run: `cd packages/observer-web && bun test`

**Step 4: Commit**

```bash
git add -u packages/observer-web/src/
git commit -m "refactor(web): replace inline button styles with Button component"
```


### Task 4: Tooltip Component

**Files:**
- Create: `packages/observer-web/src/components/tooltip.tsx`
- Modify: `packages/observer-web/src/components/row-actions.tsx`

**Step 1: Create the Tooltip component**

```tsx
import { Tooltip as BaseTooltip } from "@base-ui/react/tooltip";
import type { ReactNode } from "react";

interface TooltipProps {
  label: string;
  children: ReactNode;
  side?: "top" | "bottom" | "left" | "right";
}

export function Tooltip({ label, children, side = "top" }: TooltipProps) {
  return (
    <BaseTooltip.Provider>
      <BaseTooltip.Root delay={200}>
        <BaseTooltip.Trigger render={children as React.ReactElement} />
        <BaseTooltip.Portal>
          <BaseTooltip.Positioner side={side} sideOffset={6}>
            <BaseTooltip.Popup className="rounded-md bg-fg px-2 py-1 text-xs text-bg shadow-elevated">
              {label}
            </BaseTooltip.Popup>
          </BaseTooltip.Positioner>
        </BaseTooltip.Portal>
      </BaseTooltip.Root>
    </BaseTooltip.Provider>
  );
}
```

**Step 2: Apply tooltips to RowActions**

```tsx
import { PencilSimpleIcon, TrashIcon } from "@/components/icons";
import { Button } from "@/components/button";
import { Tooltip } from "@/components/tooltip";
import { useTranslation } from "react-i18next";

interface RowActionsProps {
  onEdit: () => void;
  onDelete: () => void;
}

export function RowActions({ onEdit, onDelete }: RowActionsProps) {
  const { t } = useTranslation();

  return (
    <div className="flex gap-1">
      <Tooltip label={t("admin.common.edit")}>
        <Button variant="ghost" onClick={(e) => { e.stopPropagation(); onEdit(); }}>
          <PencilSimpleIcon size={16} />
        </Button>
      </Tooltip>
      <Tooltip label={t("admin.common.delete")}>
        <Button variant="ghost" className="hover:text-rose" onClick={(e) => { e.stopPropagation(); onDelete(); }}>
          <TrashIcon size={16} />
        </Button>
      </Tooltip>
    </div>
  );
}
```

**Step 3: Apply tooltips to other icon-only buttons**

Search for icon-only button patterns (buttons that render only an Icon with no text). Key locations:
- `drawer-shell.tsx` — close X button
- `pagination.tsx` — prev/next caret buttons
- Person detail edit pencil buttons in various tab pages

Add `Tooltip` wrapping these buttons with appropriate i18n labels.

**Step 4: Add i18n keys if missing**

Check that `admin.common.edit`, `admin.common.delete`, `admin.common.close`, `admin.common.previousPage`, `admin.common.nextPage` exist in all 6 locale files. Add any missing ones.

**Step 5: Verify it compiles**

Run: `cd packages/observer-web && bunx tsc --noEmit --pretty 2>&1 | head -20`

**Step 6: Commit**

```bash
git add packages/observer-web/src/components/tooltip.tsx
git add -u packages/observer-web/src/
git commit -m "feat(web): add Tooltip component and apply to icon buttons"
```


### Task 6: Toast Notification System

**Files:**
- Create: `packages/observer-web/src/stores/toast.tsx`
- Create: `packages/observer-web/src/components/toast.tsx`
- Modify: `packages/observer-web/src/routes/__root.tsx`

**Step 1: Create toast store**

Use React context + useState (same pattern as auth store):

```tsx
import { createContext, useCallback, useContext, useState } from "react";
import type { ReactNode } from "react";

type ToastVariant = "success" | "error" | "info";

interface ToastItem {
  id: string;
  message: string;
  variant: ToastVariant;
}

interface ToastActions {
  success: (message: string) => void;
  error: (message: string) => void;
  info: (message: string) => void;
  dismiss: (id: string) => void;
}

interface ToastContextValue {
  toasts: ToastItem[];
  toast: ToastActions;
}

const ToastContext = createContext<ToastContextValue | null>(null);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const addToast = useCallback((message: string, variant: ToastVariant) => {
    const id = crypto.randomUUID();
    setToasts((prev) => [...prev, { id, message, variant }]);

    if (variant !== "error") {
      setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
      }, 4000);
    }
  }, []);

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const toast: ToastActions = {
    success: (msg) => addToast(msg, "success"),
    error: (msg) => addToast(msg, "error"),
    info: (msg) => addToast(msg, "info"),
    dismiss,
  };

  return (
    <ToastContext value={{ toasts, toast }}>
      {children}
    </ToastContext>
  );
}

export function useToast(): ToastActions {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used within ToastProvider");
  return ctx.toast;
}

export function useToasts(): ToastItem[] {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToasts must be used within ToastProvider");
  return ctx.toasts;
}
```

**Step 2: Create Toast UI component**

```tsx
import { CheckIcon, WarningIcon, XIcon } from "@/components/icons";
import { useToasts, useToast } from "@/stores/toast";

const variantStyles = {
  success: "border-foam/20 bg-foam/10 text-foam",
  error: "border-rose/20 bg-rose/10 text-rose",
  info: "border-accent/20 bg-accent/10 text-accent",
};

const variantIcons = {
  success: CheckIcon,
  error: WarningIcon,
  info: CheckIcon,
};

export function ToastContainer() {
  const toasts = useToasts();
  const { dismiss } = useToast();

  if (toasts.length === 0) return null;

  return (
    <div className="fixed right-4 bottom-4 z-[200] flex flex-col gap-2">
      {toasts.map((t) => {
        const Icon = variantIcons[t.variant];
        return (
          <div
            key={t.id}
            role="status"
            aria-live="polite"
            className={`flex items-center gap-2 rounded-lg border px-4 py-3 text-sm font-medium shadow-elevated backdrop-blur-sm animate-in slide-in-from-right ${variantStyles[t.variant]}`}
          >
            <Icon size={16} weight="bold" className="shrink-0" />
            <span className="flex-1">{t.message}</span>
            <button
              type="button"
              onClick={() => dismiss(t.id)}
              className="shrink-0 cursor-pointer rounded p-0.5 opacity-60 hover:opacity-100"
            >
              <XIcon size={14} />
            </button>
          </div>
        );
      })}
    </div>
  );
}
```

**Step 3: Add slide-in animation to main.css**

```css
@keyframes slide-in-from-right {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.animate-in {
  animation-duration: 200ms;
  animation-timing-function: ease-out;
  animation-fill-mode: both;
}

.slide-in-from-right {
  animation-name: slide-in-from-right;
}
```

**Step 4: Wire into __root.tsx**

Wrap the app tree with `ToastProvider` and render `ToastContainer`:

```tsx
import { ToastProvider } from "@/stores/toast";
import { ToastContainer } from "@/components/toast";

// Inside the root component:
<ToastProvider>
  <AuthProvider>
    <Outlet />
  </AuthProvider>
  <ToastContainer />
</ToastProvider>
```

**Step 5: Verify it compiles**

Run: `cd packages/observer-web && bunx tsc --noEmit --pretty 2>&1 | head -20`

**Step 6: Commit**

```bash
git add packages/observer-web/src/stores/toast.tsx packages/observer-web/src/components/toast.tsx
git add -u packages/observer-web/src/
git commit -m "feat(web): add toast notification system with success/error/info variants"
```


### Task 8: Responsive Drawer

**Files:**
- Modify: `packages/observer-web/src/components/drawer-shell.tsx`

**Step 1: Add size prop and responsive classes**

Add a `size` prop (`md` | `lg`, default `lg`) and responsive breakpoints:

```tsx
interface DrawerShellProps {
  // ... existing props
  size?: "md" | "lg";
}

const sizeClasses = {
  md: "sm:max-w-[560px]",
  lg: "sm:max-w-[840px]",
};
```

Update the Popup className:

```tsx
<Drawer.Popup className={`fixed top-0 right-0 flex h-dvh w-full flex-col border-l border-border-secondary bg-bg-secondary shadow-elevated transition-transform duration-200 ease-out data-ending-style:translate-x-full data-starting-style:translate-x-full ${sizeClasses[size ?? "lg"]}`}>
```

This makes the drawer full-width on mobile (<640px) and constrained on larger screens.

**Step 2: Add safe-area padding for mobile**

Add `pb-[env(safe-area-inset-bottom)]` to the footer div to account for iOS safe areas.

**Step 3: Verify it compiles**

Run: `cd packages/observer-web && bunx tsc --noEmit --pretty 2>&1 | head -20`

**Step 4: Commit**

```bash
git add -u packages/observer-web/src/components/drawer-shell.tsx
git commit -m "feat(web): make drawer responsive with size prop and mobile full-screen"
```


### Task 10: Inline Form Validation

**Files:**
- Modify: `packages/observer-web/src/components/form-field.tsx`

**Step 1: Add error prop to FormField**

```tsx
interface FormFieldProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  required?: boolean;
  disabled?: boolean;
  type?: string;
  maxLength?: number;
  className?: string;
  error?: string;
}

export function FormField({
  label,
  value,
  onChange,
  required,
  disabled,
  type,
  maxLength,
  className,
  error,
}: FormFieldProps) {
  return (
    <Field.Root className={className} invalid={!!error}>
      <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
        {label}
        {required && " *"}
      </Field.Label>
      <Field.Control
        required={required}
        disabled={disabled}
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        maxLength={maxLength}
        className={`${inputClass} ${error ? "!border-rose" : ""}`}
      />
      {error && (
        <Field.Error className="mt-1 text-xs text-rose" forceShow>
          {error}
        </Field.Error>
      )}
    </Field.Root>
  );
}
```

Do the same for `FormTextarea`.

**Step 2: Verify it compiles**

Run: `cd packages/observer-web && bunx tsc --noEmit --pretty 2>&1 | head -20`

**Step 3: Commit**

```bash
git add -u packages/observer-web/src/components/form-field.tsx
git commit -m "feat(web): add inline error display to FormField and FormTextarea"
```


## Summary

| Task | What | Key Files |
|------|------|-----------|
| 1 | Button component | `components/button.tsx` |
| 2 | Replace all inline buttons | ~25 files |
| 3 | Focus-visible rings | `main.css`, `form-field.tsx`, auth pages |
| 4 | Tooltip component | `components/tooltip.tsx`, `row-actions.tsx` |
| 5 | EmptyState component | `components/empty-state.tsx`, `data-table.tsx` |
| 6 | Toast system | `stores/toast.tsx`, `components/toast.tsx`, `__root.tsx` |
| 7 | Apply toasts to mutations | ~8 drawer/page files |
| 8 | Responsive drawer | `drawer-shell.tsx` |
| 9 | Collapsible sidebar | `stores/sidebar.tsx`, layouts, sidebar-link |
| 10 | Inline form validation | `form-field.tsx` |
| 11 | Micro-interactions | `main.css`, layout files |
