# Consultant QoL Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Give consultants a smoother workflow — inline support record creation from person overview, inline reference data creation from person drawer, and backend RBAC expansion so consultants can read/write reference data.

**Architecture:** Backend adds a new `adminWrite` route group for consultant+staff+admin POST/PATCH on countries/states/places/categories, and expands `adminRead` to include consultants. Frontend adds an inline collapsible form on person overview for quick support records, and `+` button dialogs next to cascading location dropdowns in the person drawer.

**Tech Stack:** Go (Gin router, middleware), React 19, TanStack Router/Query, Base UI dialogs, ky HTTP client, i18next.

---

### Task 1: Backend — Expand adminRead to include consultant

**Files:**

- Modify: `internal/server/server.go:128`

**Step 1: Update adminRead RequireRole**

Change line 128 from:

```go
adminRead := s.router.Group("/admin", authMW.Authenticate(), authMW.RequireRole(user.RoleAdmin, user.RoleStaff))
```

to:

```go
adminRead := s.router.Group("/admin", authMW.Authenticate(), authMW.RequireRole(user.RoleAdmin, user.RoleStaff, user.RoleConsultant))
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: success

**Step 3: Commit**

```bash
git add internal/server/server.go
git commit -m "allow consultants to read admin reference data"
```

---

### Task 2: Backend — Add adminWrite group for consultant reference data

**Files:**

- Modify: `internal/server/server.go:145-181`

**Step 1: Add the adminWrite group between adminRead and admin blocks**

After the `adminRead` block (after line 143), insert a new route group. Then **remove** the POST/PATCH routes for countries, states, places, categories from the `admin` block (they move to `adminWrite`). The `admin` block keeps only DELETE for these resources, plus all offices write routes.

The new `adminWrite` block:

```go
// Reference data write endpoints (admin + staff + consultant)
adminWrite := s.router.Group("/admin", authMW.Authenticate(), authMW.RequireRole(user.RoleAdmin, user.RoleStaff, user.RoleConsultant))
{
    adminWrite.POST("/countries", countryHandler.Create)
    adminWrite.PATCH("/countries/:id", countryHandler.Update)

    adminWrite.POST("/states", stateHandler.Create)
    adminWrite.PATCH("/states/:id", stateHandler.Update)

    adminWrite.POST("/places", placeHandler.Create)
    adminWrite.PATCH("/places/:id", placeHandler.Update)

    adminWrite.POST("/categories", categoryHandler.Create)
    adminWrite.PATCH("/categories/:id", categoryHandler.Update)
}
```

The `admin` block (admin-only) should keep:

- All user management (POST/PATCH/POST reset-password)
- All project management (GET/POST/PATCH projects, permissions CRUD)
- All office write routes (POST/PATCH/DELETE offices)
- DELETE for countries, states, places, categories

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: success

**Step 3: Commit**

```bash
git add internal/server/server.go
git commit -m "add adminWrite group for consultant reference data access"
```

---

### Task 3: Frontend — Add i18n keys for new UI elements

**Files:**

- Modify: `packages/observer-web/src/locales/en.json`
- Modify: `packages/observer-web/src/locales/ky.json`
- Modify: `packages/observer-web/src/locales/uk.json`
- Modify: `packages/observer-web/src/locales/de.json`
- Modify: `packages/observer-web/src/locales/tr.json`

**Step 1: Add keys to en.json**

Inside `project.people`, add:

```json
"quickSupportTitle": "Quick add support record",
"quickSupportAdd": "Add support record",
"quickSupportSaved": "Support record saved"
```

Inside `project.people`, add (for the add-reference dialogs):

```json
"addCountry": "Add country",
"addState": "Add state",
"addPlace": "Add place"
```

Inside `admin.common` or create `common` section if not present, add:

```json
"nameRequired": "Name is required"
```

**Step 2: Add equivalent keys to ky.json, uk.json, de.json, tr.json**

Use the Kyrgyz Latin transliteration rules for ky.json. Translate appropriately for the other locales.

**Step 3: Commit**

```bash
git add packages/observer-web/src/locales/*.json
git commit -m "add i18n keys for consultant QoL features"
```

---

### Task 4: Frontend — Inline support record form on person overview

**Files:**

- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId/index.tsx`

**Step 1: Add the inline form component**

Below the existing `<section>` card in `PersonOverview`, add a collapsible inline form. The component needs:

- Import `useState` (already imported via createFileRoute deps), `useTranslation`, `useCreateSupportRecord`, `UISelect`, `CheckIcon`
- State: `showForm: boolean`, `saved: boolean`, `error: string`, and form fields (`type`, `sphere`, `provided_at`, `notes`)
- Default `provided_at` to today's date (`new Date().toISOString().slice(0, 10)`)
- `+ Add support record` button that toggles `showForm`
- When expanded, show a flat form card with subtle section dividers:

```tsx
{
  /* Section divider pattern */
}
<div className="flex items-center gap-3">
  <span className="text-xs font-semibold uppercase tracking-wide text-fg-tertiary">
    {t("project.supportRecords.recordInfo")}
  </span>
  <span className="h-px flex-1 bg-border-secondary" />
</div>;
```

- Form fields in a 2-col grid: Type (select, required), Sphere (select), Provided at (date), Notes (textarea 2 rows, full width)
- Save / Cancel buttons at bottom right
- On save: call `createSupportRecord.mutateAsync({ person_id: personId, type, sphere?, provided_at?, notes? })`, invalidate `["support-records", projectId]`, show success indicator, collapse form after 1.5s
- On cancel: collapse form, reset fields

Use the same `inputClass` pattern and `UISelect` component as existing drawers. Reuse existing translation keys from `project.supportRecords.*` for type/sphere options.

**Step 2: Test manually**

Navigate to a person's overview tab, verify:

- Button appears below person details
- Form expands/collapses
- Type select works
- Save creates a record (check support records tab)
- Success indicator shows briefly

**Step 3: Commit**

```bash
git add packages/observer-web/src/routes/_app/projects/\$projectId/people/\$personId/index.tsx
git commit -m "add inline support record form on person overview"
```

---

### Task 5: Frontend — Add reference data dialog component

**Files:**

- Create: `packages/observer-web/src/components/add-reference-dialog.tsx`

**Step 1: Create a reusable dialog component**

This component renders a Base UI dialog for adding a new reference data item (country, state, or place). It accepts:

```tsx
interface AddReferenceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  children: React.ReactNode;
  onSubmit: () => void;
  isPending: boolean;
  error: string;
}
```

Use `@base-ui/react/dialog` (check existing usage in the codebase for the correct import). The dialog should be a small centered modal with:

- Title
- Error banner (if error)
- `{children}` for form fields
- Cancel / Save buttons in footer

Pattern:

```tsx
import { Dialog } from "@base-ui/react/dialog";

export function AddReferenceDialog({
  open,
  onOpenChange,
  title,
  children,
  onSubmit,
  isPending,
  error,
}: AddReferenceDialogProps) {
  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 z-50 bg-black/25 backdrop-blur-xs ..." />
        <Dialog.Popup className="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border border-border-secondary bg-bg-secondary p-6 shadow-elevated ...">
          <Dialog.Title className="font-serif text-lg font-semibold text-fg">{title}</Dialog.Title>
          {error && <error banner />}
          <form
            onSubmit={(e) => {
              e.preventDefault();
              onSubmit();
            }}
            className="mt-4 space-y-4"
          >
            {children}
            <div className="flex justify-end gap-2 pt-2">
              <Dialog.Close className="...cancel styles...">
                {t("admin.common.cancel")}
              </Dialog.Close>
              <button type="submit" disabled={isPending} className="...accent button...">
                {isPending ? t("project.people.saving") : t("project.people.save")}
              </button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
```

**Step 2: Commit**

```bash
git add packages/observer-web/src/components/add-reference-dialog.tsx
git commit -m "add reusable AddReferenceDialog component"
```

---

### Task 6: Frontend — Wire add-reference dialogs into person drawer

**Files:**

- Modify: `packages/observer-web/src/components/person-drawer.tsx`

**Step 1: Add state and imports**

Add imports for:

- `AddReferenceDialog` from `@/components/add-reference-dialog`
- `useCreateCountry` from `@/hooks/use-countries`
- `useCreateState` from `@/hooks/use-states`
- `useCreatePlace` from `@/hooks/use-places`
- `PlusIcon` from `@/components/icons`
- `Field` is already imported
- `useQueryClient` is already imported

Add state for 3 dialogs:

```tsx
const [addCountryOpen, setAddCountryOpen] = useState(false);
const [addStateOpen, setAddStateOpen] = useState<{ open: boolean; forOrigin: boolean }>({
  open: false,
  forOrigin: true,
});
const [addPlaceOpen, setAddPlaceOpen] = useState<{ open: boolean; forOrigin: boolean }>({
  open: false,
  forOrigin: true,
});

// Form state for each dialog
const [newCountryName, setNewCountryName] = useState("");
const [newCountryCode, setNewCountryCode] = useState("");
const [newStateName, setNewStateName] = useState("");
const [newStateConflictZone, setNewStateConflictZone] = useState("");
const [newPlaceName, setNewPlaceName] = useState("");

const [dialogError, setDialogError] = useState("");

const createCountry = useCreateCountry();
const createState = useCreateState();
const createPlace = useCreatePlace();
```

**Step 2: Change location grid layout**

For each of the 6 location selects (origin country/state/place, current country/state/place), change from plain `UISelect` to a flex row with the select and a `+` button:

```tsx
<div className="flex gap-1.5">
  <div className="flex-1">
    <UISelect
      value={form.origin_country}
      onValueChange={(v) => {
        set("origin_country", v);
        set("origin_state", "");
        set("origin_place_id", "");
      }}
      options={countryOptions}
      placeholder={t("project.people.selectCountry")}
      fullWidth
    />
  </div>
  <button
    type="button"
    onClick={() => {
      setDialogError("");
      setNewCountryName("");
      setNewCountryCode("");
      setAddCountryOpen(true);
    }}
    className="inline-flex size-9 shrink-0 items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary hover:border-border-primary hover:text-fg"
    title={t("project.people.addCountry")}
  >
    <PlusIcon size={16} />
  </button>
</div>
```

For state `+` button: `disabled={!form.origin_country}` (or `!form.current_country` for current section).
For place `+` button: `disabled={!form.origin_state}` (or `!form.current_state`).

**Step 3: Add dialog handlers**

```tsx
async function handleAddCountry() {
  setDialogError("");
  try {
    const created = await createCountry.mutateAsync({ name: newCountryName, code: newCountryCode });
    // Auto-select in whichever section triggered it (both use same country list)
    // Since country dialog doesn't track origin/current, just close and let user pick
    setAddCountryOpen(false);
  } catch (err) {
    if (err instanceof HTTPError) {
      const body = await err.response.json().catch(() => null);
      setDialogError(body?.error || err.message);
    }
  }
}

async function handleAddState() {
  setDialogError("");
  const countryId = addStateOpen.forOrigin ? form.origin_country : form.current_country;
  try {
    const created = await createState.mutateAsync({
      countryId,
      data: {
        name: newStateName,
        ...(newStateConflictZone && { conflict_zone: newStateConflictZone }),
      },
    });
    if (addStateOpen.forOrigin) {
      set("origin_state", created.id);
      set("origin_place_id", "");
    } else {
      set("current_state", created.id);
      set("current_place_id", "");
    }
    setAddStateOpen({ open: false, forOrigin: true });
  } catch (err) {
    if (err instanceof HTTPError) {
      const body = await err.response.json().catch(() => null);
      setDialogError(body?.error || err.message);
    }
  }
}

async function handleAddPlace() {
  setDialogError("");
  const stateId = addPlaceOpen.forOrigin ? form.origin_state : form.current_state;
  try {
    const created = await createPlace.mutateAsync({
      stateId,
      data: { name: newPlaceName },
    });
    if (addPlaceOpen.forOrigin) {
      set("origin_place_id", created.id);
    } else {
      set("current_place_id", created.id);
    }
    setAddPlaceOpen({ open: false, forOrigin: true });
  } catch (err) {
    if (err instanceof HTTPError) {
      const body = await err.response.json().catch(() => null);
      setDialogError(body?.error || err.message);
    }
  }
}
```

**Step 4: Add dialog JSX at end of form (before closing `</form>`)**

```tsx
<AddReferenceDialog
  open={addCountryOpen}
  onOpenChange={setAddCountryOpen}
  title={t("project.people.addCountry")}
  onSubmit={handleAddCountry}
  isPending={createCountry.isPending}
  error={dialogError}
>
  <Field.Root>
    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
      {t("admin.reference.countries.name")} *
    </Field.Label>
    <Field.Control
      required
      value={newCountryName}
      onChange={(e) => setNewCountryName(e.target.value)}
      className={inputClass}
    />
  </Field.Root>
  <Field.Root>
    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
      {t("admin.reference.countries.code")}
    </Field.Label>
    <Field.Control
      value={newCountryCode}
      onChange={(e) => setNewCountryCode(e.target.value)}
      className={inputClass}
      maxLength={3}
    />
  </Field.Root>
</AddReferenceDialog>

<AddReferenceDialog
  open={addStateOpen.open}
  onOpenChange={(v) => setAddStateOpen((s) => ({ ...s, open: v }))}
  title={t("project.people.addState")}
  onSubmit={handleAddState}
  isPending={createState.isPending}
  error={dialogError}
>
  <Field.Root>
    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
      {t("admin.reference.states.name")} *
    </Field.Label>
    <Field.Control
      required
      value={newStateName}
      onChange={(e) => setNewStateName(e.target.value)}
      className={inputClass}
    />
  </Field.Root>
  <Field.Root>
    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
      {t("admin.reference.states.conflictZone")}
    </Field.Label>
    <Field.Control
      value={newStateConflictZone}
      onChange={(e) => setNewStateConflictZone(e.target.value)}
      className={inputClass}
    />
  </Field.Root>
</AddReferenceDialog>

<AddReferenceDialog
  open={addPlaceOpen.open}
  onOpenChange={(v) => setAddPlaceOpen((s) => ({ ...s, open: v }))}
  title={t("project.people.addPlace")}
  onSubmit={handleAddPlace}
  isPending={createPlace.isPending}
  error={dialogError}
>
  <Field.Root>
    <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
      {t("admin.reference.places.name")} *
    </Field.Label>
    <Field.Control
      required
      value={newPlaceName}
      onChange={(e) => setNewPlaceName(e.target.value)}
      className={inputClass}
    />
  </Field.Root>
</AddReferenceDialog>
```

**Step 5: Test manually**

- Open person drawer
- Click `+` next to country → dialog opens, fill name+code, save → country appears in dropdown
- Select country → `+` next to state becomes enabled → add state → auto-selected
- Select state → `+` next to place becomes enabled → add place → auto-selected
- Repeat for current location section

**Step 6: Commit**

```bash
git add packages/observer-web/src/components/person-drawer.tsx packages/observer-web/src/components/add-reference-dialog.tsx
git commit -m "add inline reference data creation dialogs in person drawer"
```

---

### Task 7: Verify and final commit

**Step 1: Run backend tests**

Run: `just test`
Expected: all pass (except pre-existing middleware test)

**Step 2: Run frontend type check**

Run: `cd packages/observer-web && bunx tsc --noEmit`
Expected: no errors

**Step 3: Manual smoke test**

1. Log in as consultant
2. Navigate to a project → People → select a person
3. On overview tab: click "Add support record", fill type+sphere, save → success
4. Check Support Records tab → new record visible
5. Go back to people list → open person drawer
6. In location section: click `+` next to country → add country dialog → save → new country in dropdown
7. Select country → `+` next to state → add state → auto-selected
8. Verify the same works for current location section

**Step 4: Final commit if any fixups needed**

```bash
git add -A
git commit -m "consultant QoL: inline support records, reference data dialogs, RBAC expansion"
```
