# Remaining Features Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement all missing features: frontend project-scoped pages, auth flows (admin approval, profile, password change/reset), and reports (ADR-005) with D3 visualizations.

**Architecture:** Follow existing DDD + Clean Architecture patterns. Frontend pages follow the people list page pattern (DataTable + drawer/dialog). Backend reports add a new `report` domain + repository + usecase + handler. Auth flows extend existing auth usecase. Registration is open but users start inactive until admin approval.

**Tech Stack:** Go 1.25, Gin, sqlx, PostgreSQL | React 19, TanStack Router/Query, ky, D3, Tailwind, Base UI, i18next

---

## Phase 1: Frontend — Project-Scoped Pages

Backend APIs already exist for all these resources. We need frontend pages under `/projects/$projectId/`.

### Task 1: Project layout with sidebar navigation

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId.tsx` (layout)
- Create: `packages/observer-web/src/components/project-sidebar.tsx`

The project section needs a layout with sidebar navigation to switch between people, support records, migration records, households, tags, documents, pets.

**Step 1: Create project sidebar component**

```tsx
// packages/observer-web/src/components/project-sidebar.tsx
import { useTranslation } from "react-i18next";
import { SidebarLink } from "./sidebar-link";
import { Users, HandHeart, Path, HouseSimple, Tag, Files, PawPrint } from "@phosphor-icons/react";

interface ProjectSidebarProps {
  projectId: string;
}

export function ProjectSidebar({ projectId }: ProjectSidebarProps) {
  const { t } = useTranslation();
  const base = `/projects/${projectId}`;

  const links = [
    { to: `${base}/people`, label: t("project.nav.people"), icon: Users },
    { to: `${base}/support-records`, label: t("project.nav.supportRecords"), icon: HandHeart },
    { to: `${base}/migration-records`, label: t("project.nav.migrationRecords"), icon: Path },
    { to: `${base}/households`, label: t("project.nav.households"), icon: HouseSimple },
    { to: `${base}/tags`, label: t("project.nav.tags"), icon: Tag },
    { to: `${base}/documents`, label: t("project.nav.documents"), icon: Files },
    { to: `${base}/pets`, label: t("project.nav.pets"), icon: PawPrint },
  ];

  return (
    <nav className="project-sidebar">
      {links.map((link) => (
        <SidebarLink key={link.to} to={link.to} icon={link.icon}>
          {link.label}
        </SidebarLink>
      ))}
    </nav>
  );
}
```

**Step 2: Create project layout route**

```tsx
// packages/observer-web/src/routes/_app/projects/$projectId.tsx
import { Outlet, createFileRoute } from "@tanstack/react-router";
import { ProjectSidebar } from "@/components/project-sidebar";

export const Route = createFileRoute("/_app/projects/$projectId")({
  component: ProjectLayout,
});

function ProjectLayout() {
  const { projectId } = Route.useParams();

  return (
    <div className="project-layout">
      <ProjectSidebar projectId={projectId} />
      <main className="project-content">
        <Outlet />
      </main>
    </div>
  );
}
```

**Step 3: Add i18n keys for project navigation**

Add to all 4 locale files under `project.nav`:

```json
{
  "project": {
    "nav": {
      "people": "Adamdar",
      "supportRecords": "Koldo körsötüü",
      "migrationRecords": "Köchüü",
      "households": "Üy bülöö",
      "tags": "Tegder",
      "documents": "Dokumentter",
      "pets": "Üy janybarlar"
    }
  }
}
```

Equivalent translations for en, uk, de.

**Step 4: Add styles for project layout**

Add to `packages/observer-web/src/main.css`:

```css
.project-layout {
  @apply flex gap-6;
}

.project-sidebar {
  @apply flex flex-col gap-1 w-56 shrink-0;
}

.project-content {
  @apply flex-1 min-w-0;
}
```

**Step 5: Verify and commit**

Run: `cd packages/observer-web && bunx tsr generate`
Run: `cd packages/observer-web && bun run build`

---

### Task 2: Frontend types for remaining project-scoped resources

**Files:**

- Create: `packages/observer-web/src/types/support-record.ts`
- Create: `packages/observer-web/src/types/migration-record.ts`
- Create: `packages/observer-web/src/types/household.ts`
- Create: `packages/observer-web/src/types/note.ts`
- Create: `packages/observer-web/src/types/document.ts`
- Create: `packages/observer-web/src/types/pet.ts`
- Create: `packages/observer-web/src/types/tag.ts`

**Step 1: Create type files**

Each type file matches the backend DTOs. Pattern from `person.ts`:

```ts
// types/support-record.ts
export type SupportType = "humanitarian" | "legal" | "social" | "psychological" | "medical" | "general";
export type SupportSphere =
  | "housing_assistance"
  | "document_recovery"
  | "social_benefits"
  | "property_rights"
  | "employment_rights"
  | "family_law"
  | "healthcare_access"
  | "education_access"
  | "financial_aid"
  | "psychological_support"
  | "other";
export type ReferralStatus = "pending" | "accepted" | "completed" | "declined" | "no_response";

export interface SupportRecord {
  id: string;
  person_id: string;
  project_id: string;
  consultant_id?: string;
  recorded_by?: string;
  office_id?: string;
  referred_to_office?: string;
  type: SupportType;
  sphere?: SupportSphere;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateSupportRecordInput {
  person_id: string;
  type: SupportType;
  sphere?: SupportSphere;
  consultant_id?: string;
  office_id?: string;
  referred_to_office?: string;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
}

export interface UpdateSupportRecordInput {
  type?: SupportType;
  sphere?: SupportSphere;
  consultant_id?: string;
  office_id?: string;
  referred_to_office?: string;
  referral_status?: ReferralStatus;
  provided_at?: string;
  notes?: string;
}

export interface ListSupportRecordsParams {
  person_id?: string;
  consultant_id?: string;
  office_id?: string;
  type?: SupportType;
  sphere?: SupportSphere;
  page?: number;
  per_page?: number;
}

export interface ListSupportRecordsOutput {
  records: SupportRecord[];
  total: number;
  page: number;
  per_page: number;
}
```

```ts
// types/migration-record.ts
export type MovementReason =
  | "conflict"
  | "security"
  | "service_access"
  | "return"
  | "relocation_program"
  | "economic"
  | "other";
export type HousingAtDestination =
  | "own_property"
  | "renting"
  | "with_relatives"
  | "collective_site"
  | "hotel"
  | "other"
  | "unknown";

export interface MigrationRecord {
  id: string;
  person_id: string;
  from_place_id?: string;
  destination_place_id?: string;
  migration_date?: string;
  movement_reason?: MovementReason;
  housing_at_destination?: HousingAtDestination;
  notes?: string;
  created_at: string;
}

export interface CreateMigrationRecordInput {
  from_place_id?: string;
  destination_place_id?: string;
  migration_date?: string;
  movement_reason?: MovementReason;
  housing_at_destination?: HousingAtDestination;
  notes?: string;
}
```

```ts
// types/household.ts
export type Relationship =
  | "head"
  | "spouse"
  | "child"
  | "parent"
  | "sibling"
  | "grandchild"
  | "grandparent"
  | "other_relative"
  | "non_relative";

export interface Household {
  id: string;
  project_id: string;
  reference_number?: string;
  head_person_id?: string;
  members?: HouseholdMember[];
  created_at: string;
  updated_at: string;
}

export interface HouseholdMember {
  person_id: string;
  relationship: string;
}

export interface CreateHouseholdInput {
  reference_number?: string;
  head_person_id?: string;
}

export interface UpdateHouseholdInput {
  reference_number?: string;
  head_person_id?: string;
}

export interface AddMemberInput {
  person_id: string;
  relationship: Relationship;
}

export interface ListHouseholdsParams {
  page?: number;
  per_page?: number;
}

export interface ListHouseholdsOutput {
  households: Household[];
  total: number;
  page: number;
  per_page: number;
}
```

```ts
// types/note.ts
export interface Note {
  id: string;
  person_id: string;
  author_id?: string;
  body: string;
  created_at: string;
}

export interface CreateNoteInput {
  body: string;
}
```

```ts
// types/document.ts
export interface Document {
  id: string;
  person_id: string;
  project_id: string;
  uploaded_by?: string;
  name: string;
  path: string;
  mime_type: string;
  size: number;
  created_at: string;
}

export interface CreateDocumentInput {
  person_id: string;
  name: string;
  path: string;
  mime_type: string;
  size: number;
}
```

```ts
// types/pet.ts
export type PetStatus = "registered" | "adopted" | "owner_found" | "needs_shelter" | "unknown";

export interface Pet {
  id: string;
  project_id: string;
  owner_id?: string;
  name: string;
  status: PetStatus;
  registration_id?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePetInput {
  owner_id?: string;
  name: string;
  status?: PetStatus;
  registration_id?: string;
  notes?: string;
}

export interface UpdatePetInput {
  owner_id?: string;
  name?: string;
  status?: PetStatus;
  registration_id?: string;
  notes?: string;
}

export interface ListPetsParams {
  page?: number;
  per_page?: number;
}

export interface ListPetsOutput {
  pets: Pet[];
  total: number;
  page: number;
  per_page: number;
}
```

```ts
// types/tag.ts
export interface Tag {
  id: string;
  project_id: string;
  name: string;
  created_at: string;
}

export interface CreateTagInput {
  name: string;
}
```

**Step 2: Verify build**

Run: `cd packages/observer-web && bun run build`

---

### Task 3: React Query hooks for project-scoped resources

**Files:**

- Create: `packages/observer-web/src/hooks/use-support-records.ts`
- Create: `packages/observer-web/src/hooks/use-migration-records.ts`
- Create: `packages/observer-web/src/hooks/use-households.ts`
- Create: `packages/observer-web/src/hooks/use-notes.ts`
- Create: `packages/observer-web/src/hooks/use-documents.ts`
- Create: `packages/observer-web/src/hooks/use-pets.ts`
- Create: `packages/observer-web/src/hooks/use-tags.ts`

Follow the `use-people.ts` pattern. Each hook file exports:

- `useList*` — query with params, `keepPreviousData`
- `useGet*` — single resource query (if applicable)
- `useCreate*` — mutation, invalidates list key
- `useUpdate*` — mutation, invalidates list key (if applicable)
- `useDelete*` — mutation, invalidates list key (if applicable)

**Step 1: Create hooks**

Example — support records hook:

```ts
// hooks/use-support-records.ts
import { useQuery, useMutation, useQueryClient, keepPreviousData } from "@tanstack/react-query";
import api from "@/lib/api";
import type {
  ListSupportRecordsOutput,
  ListSupportRecordsParams,
  SupportRecord,
  CreateSupportRecordInput,
  UpdateSupportRecordInput,
} from "@/types/support-record";

export function useSupportRecords(projectId: string, params?: ListSupportRecordsParams) {
  return useQuery({
    queryKey: ["support-records", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/support-records`, { searchParams: params as Record<string, string> })
        .json<ListSupportRecordsOutput>(),
    placeholderData: keepPreviousData,
  });
}

export function useSupportRecord(projectId: string, id: string) {
  return useQuery({
    queryKey: ["support-records", projectId, id],
    queryFn: () => api.get(`projects/${projectId}/support-records/${id}`).json<SupportRecord>(),
    enabled: !!id,
  });
}

export function useCreateSupportRecord(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateSupportRecordInput) =>
      api.post(`projects/${projectId}/support-records`, { json: data }).json<SupportRecord>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["support-records", projectId] }),
  });
}

export function useUpdateSupportRecord(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }: UpdateSupportRecordInput & { id: string }) =>
      api.patch(`projects/${projectId}/support-records/${id}`, { json: data }).json<SupportRecord>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["support-records", projectId] }),
  });
}

export function useDeleteSupportRecord(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`projects/${projectId}/support-records/${id}`).json(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["support-records", projectId] }),
  });
}
```

Other hooks follow identical pattern adjusting:

- Migration records: person-scoped paths (`projects/${projectId}/people/${personId}/migration-records`), no update/delete
- Notes: person-scoped, no update
- Documents: person-scoped for list, project-scoped for create/delete
- Households: project-scoped, includes `useAddMember`, `useRemoveMember`
- Tags: project-scoped, no update, list returns `{ tags: Tag[] }`
- Pets: project-scoped, full CRUD

**Step 2: Verify build**

Run: `cd packages/observer-web && bun run build`

---

### Task 4: Tags management page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/tags/index.tsx`

Simple list + create dialog + delete confirmation. No edit (tags are name-only, just delete and recreate).

**Step 1: Create tags page**

Follow the countries reference page pattern — DataTable + FormDialog for create + ConfirmDialog for delete.

Columns: name, created_at, actions (delete only).

**Step 2: Add i18n keys**

Add `project.tags.*` keys to all 4 locale files:

```json
{
  "project": {
    "tags": {
      "title": "Tegder",
      "add": "Teg koshuu",
      "name": "Atalyşy",
      "deleteConfirm": "Tegdi jokko çygarasyzby?",
      "deleted": "Teg öçürüldü"
    }
  }
}
```

**Step 3: Verify and commit**

Run: `cd packages/observer-web && bun run build`

---

### Task 5: Support records page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/support-records/index.tsx`
- Create: `packages/observer-web/src/components/support-record-drawer.tsx`

Full CRUD with drawer for create/edit. Most complex page — has type, sphere, referral status, office, consultant, dates.

**Step 1: Create support record drawer**

Follow person-drawer pattern. Form sections:

- Record info: type (select), sphere (select), provided_at (date input)
- People: person_id (search/select from project people)
- Referral: referral_status (select), referred_to_office (select from offices)
- Assignment: consultant_id (user search), office_id (select)
- Notes: textarea

**Step 2: Create support records list page**

DataTable columns: person name (or ID), type, sphere, provided_at, referral_status, actions.
Tab filtering by support type (all, legal, social, humanitarian, etc.).

**Step 3: Add i18n keys for support records**

Add `project.supportRecords.*` keys with all field labels, type values, sphere values, referral statuses.

**Step 4: Verify and commit**

Run: `cd packages/observer-web && bun run build`

---

### Task 6: Households page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/households/index.tsx`
- Create: `packages/observer-web/src/components/household-drawer.tsx`

CRUD with members management. Drawer shows household details + member list with add/remove.

**Step 1: Create household drawer**

Form sections:

- Household info: reference_number (text), head_person_id (person search/select)
- Members list: table of current members with relationship + remove button
- Add member form: person_id (search/select), relationship (select)

**Step 2: Create households list page**

DataTable columns: reference_number, head person name, member count, created_at, actions.
Pagination.

**Step 3: Add i18n keys for households**

Add `project.households.*` keys including relationship type labels.

**Step 4: Verify and commit**

Run: `cd packages/observer-web && bun run build`

---

### Task 7: Person detail page with tabs (migration records, notes, documents)

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId.tsx` (person detail layout)
- Create: `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId/index.tsx` (overview)
- Create: `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId/migration-records.tsx`
- Create: `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId/notes.tsx`
- Create: `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId/documents.tsx`

**Step 1: Create person detail layout**

Tabbed layout: Overview, Migration Records, Notes, Documents.

**Step 2: Create migration records tab**

List of migration records for this person + create form (dialog):

- from_place_id (cascading country > state > place select)
- destination_place_id (same cascade)
- migration_date (date)
- movement_reason (select)
- housing_at_destination (select)
- notes (textarea)

Append-only: create + view, no edit/delete.

**Step 3: Create notes tab**

Timeline-style list of notes (body, author, created_at). Textarea + submit button at top for new note. Delete button per note (confirm dialog). Append-only: create + delete, no edit.

**Step 4: Create documents tab**

Metadata-only for now (no file upload — deferred). DataTable: name, mime_type, size, uploaded_by, created_at, actions (delete). Create dialog: name, path, mime_type, size fields.

**Step 5: Add i18n keys**

Add `project.migrationRecords.*`, `project.notes.*`, `project.documents.*` keys including movement reason and housing labels.

**Step 6: Verify and commit**

---

### Task 8: Pets page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/pets/index.tsx`
- Create: `packages/observer-web/src/components/pet-drawer.tsx`

Full CRUD with drawer. Simple entity.

**Step 1: Create pet drawer**

Form fields: name, status (select), owner_id (person search/select), registration_id, notes.

**Step 2: Create pets list page**

DataTable columns: name, status (badge), owner, registration_id, actions.
Tab filtering by status (all, registered, adopted, etc.).

**Step 3: Add i18n keys**

Add `project.pets.*` keys including status labels.

**Step 4: Verify and commit**

Run: `cd packages/observer-web && bun run build`

---

### Task 9: Update dashboard with project-scoped navigation

**Files:**

- Modify: `packages/observer-web/src/routes/_app/index.tsx`
- Modify: `packages/observer-web/src/routes/_app/admin/projects/index.tsx`

**Step 1: Update project links**

"My Projects" cards on dashboard should link to `/projects/$projectId/people` (already correct).
Admin projects list "Browse" action should link to `/projects/$projectId/people`.

**Step 2: Verify and commit**

---

## Phase 2: Backend + Frontend — Auth Flows

### Task 10: Registration approval flow (users start inactive)

Registration is open, but users start as inactive (`is_active = false`). Admins approve users by activating them via `PATCH /admin/users/:id`. No email verification needed — admins vouch for users directly.

**Files:**

- Modify: `internal/usecase/auth/register_usecase.go` — set `IsActive: false`
- Modify: `packages/observer-web/src/stores/auth.tsx` — handle inactive user error on login
- Modify: `packages/observer-web/src/routes/login.tsx` — show "pending approval" message

**Step 1: Update register usecase**

Change `IsActive: true` to `IsActive: false` and update the success message:

```go
// In register_usecase.go, change:
newUser := &user.User{
    // ...
    IsActive:   false,  // was: true — admin must approve
    // ...
}

// Update response message:
return &RegisterOutput{
    Message: "Registration successful. Your account is pending admin approval.",
}, nil
```

**Step 2: Update frontend login error handling**

When login returns `errors.auth.userNotActive`, show a clear message like "Your account is pending admin approval" instead of a generic error.

**Step 3: Add admin user activation UI**

The admin user edit page already has an `is_active` toggle. Add a visual indicator on the admin users list for pending (inactive) users — e.g., a "Pending" badge using the `gold` variant.

**Step 4: Add i18n keys**

```json
{
  "auth": {
    "pendingApproval": "Sizdin akkauntunuz admin tarabyndan tastyktaluunu küttö"
  }
}
```

**Step 5: Verify and commit**

Run: `just test && cd packages/observer-web && bun run build`

---

### Task 11: User self-profile update endpoint

**Files:**

- Create: `internal/usecase/auth/update_profile_usecase.go`
- Create: `internal/usecase/auth/update_profile_types.go`
- Modify: `internal/handler/auth_handler.go` — add UpdateProfile method
- Modify: `internal/server/server.go` — add `PATCH /auth/me` route
- Modify: `internal/app/container.go` — wire UpdateProfileUseCase

**Step 1: Create update profile types**

```go
// internal/usecase/auth/update_profile_types.go
package auth

// UpdateProfileInput holds fields a user can update on their own profile.
type UpdateProfileInput struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
}
```

**Step 2: Create update profile usecase**

```go
// internal/usecase/auth/update_profile_usecase.go
package auth

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/repository"
)

type UpdateProfileUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateProfileUseCase(userRepo repository.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userRepo: userRepo}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, input UpdateProfileInput) (*UserDTO, error) {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if input.FirstName != nil {
		u.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		u.LastName = input.LastName
	}
	if input.Phone != nil {
		u.Phone = input.Phone
	}

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	dto := userToDTO(u)
	return &dto, nil
}
```

**Step 3: Add handler method**

Add `UpdateProfile` method to `AuthHandler` — extracts userID from context, binds JSON, calls usecase.

**Step 4: Wire in container and add route**

Add `PATCH /auth/me` in the authenticated auth group.

**Step 5: Create frontend profile page**

- Create: `packages/observer-web/src/routes/_app/profile.tsx`
- Simple form with first_name, last_name, phone fields, save button.

**Step 6: Verify and commit**

Run: `just test`

---

### Task 12: Password change endpoint (authenticated)

**Files:**

- Create: `internal/usecase/auth/change_password_usecase.go`
- Modify: `internal/handler/auth_handler.go` — add ChangePassword method
- Modify: `internal/server/server.go` — add `POST /auth/change-password` route
- Modify: `internal/app/container.go` — wire ChangePasswordUseCase

**Step 1: Create change password usecase**

```go
// internal/usecase/auth/change_password_usecase.go
package auth

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/crypto"
	domainuser "github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type ChangePasswordUseCase struct {
	userRepo repository.UserRepository
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

func NewChangePasswordUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
		credRepo: credRepo,
		hasher:   hasher,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, userID string, input ChangePasswordInput) error {
	cred, err := uc.credRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get credentials: %w", err)
	}

	if !uc.hasher.Verify(input.CurrentPassword, cred.PasswordHash) {
		return domainuser.ErrInvalidCredentials
	}

	hash, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	cred.PasswordHash = hash
	if err := uc.credRepo.Update(ctx, cred); err != nil {
		return fmt.Errorf("update credentials: %w", err)
	}

	return nil
}
```

**Step 2: Add handler method and route**

`POST /auth/change-password` — authenticated, extracts userID, returns 200 on success.

**Step 3: Add frontend change password form**

Add a "Change Password" section to the profile page with current_password, new_password, confirm_password fields.

**Step 4: Add i18n keys and verify**

Run: `just test`

---

### Task 13: Admin password reset for users

In humanitarian/NGO deployments, self-service password reset (forgot password + email) doesn't apply — field staff may lack reliable email, and the admin vouches for users directly. Instead, admins can reset a user's password.

**Files:**

- Create: `internal/usecase/admin/reset_password_usecase.go`
- Modify: `internal/handler/admin_handler.go` — add ResetPassword method
- Modify: `internal/server/server.go` — add `POST /admin/users/:id/reset-password` route
- Modify: `internal/app/container.go` — wire ResetPasswordUseCase

**Step 1: Create admin reset password usecase**

```go
// internal/usecase/admin/reset_password_usecase.go
package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/repository"
)

type ResetPasswordInput struct {
	NewPassword string `json:"new_password" binding:"required"`
}

type ResetPasswordUseCase struct {
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

func NewResetPasswordUseCase(
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{credRepo: credRepo, hasher: hasher}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, userID string, input ResetPasswordInput) error {
	cred, err := uc.credRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get credentials: %w", err)
	}

	hash, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	cred.PasswordHash = hash
	if err := uc.credRepo.Update(ctx, cred); err != nil {
		return fmt.Errorf("update credentials: %w", err)
	}

	return nil
}
```

**Step 2: Add handler method and route**

`POST /admin/users/:id/reset-password` — admin only, binds `ResetPasswordInput`, returns 200.

**Step 3: Add frontend "Reset Password" button to user edit page**

Add a dialog on the admin user edit page (`/admin/users/$userId`) with a new password field.

**Step 4: Add i18n keys and verify**

Run: `just test`

---

## Phase 3: Reports (ADR-005) with D3 Visualizations

### Task 14: Report domain types

**Files:**

- Create: `internal/domain/report/types.go`

**Step 1: Define report request/response types**

```go
package report

import "time"

// ReportFilter contains common filter parameters for all reports.
type ReportFilter struct {
	ProjectID string
	DateFrom  *time.Time
	DateTo    *time.Time
}

// CountResult represents a single count in a report breakdown.
type CountResult struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// GroupedReport contains a title and breakdown rows.
type GroupedReport struct {
	Title   string        `json:"title"`
	Total   int           `json:"total"`
	Rows    []CountResult `json:"rows"`
}
```

**Step 2: Verify and commit**

Run: `just test`

---

### Task 15: Report repository interface

**Files:**

- Modify: `internal/repository/interfaces.go` — add ReportRepository interface

**Step 1: Add report repository interface**

```go
// ReportRepository provides aggregation queries for ADR-005 reports.
type ReportRepository interface {
	// Group 1: General consultation counts
	CountConsultations(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 2: Sex breakdown
	CountBySex(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 3: IDP/geographic status
	CountByIDPStatus(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 4: Vulnerability categories
	CountByCategory(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 5: Current region
	CountByCurrentRegion(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 6: Support sphere
	CountBySphere(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 7: Office breakdown
	CountByOffice(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 8: Age group
	CountByAgeGroup(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 9: Tag search
	CountByTag(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	// Group 10: Family units
	CountFamilyUnits(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
}
```

**Step 2: Regenerate mocks**

Run: `just generate-mocks`

**Step 3: Verify and commit**

---

### Task 16: Report repository SQL implementation

**Files:**

- Create: `internal/repository/report_repository.go`

**Step 1: Implement report SQL queries**

Each method builds a query following ADR-005 patterns. Example for Group 1:

```go
func (r *reportRepository) CountConsultations(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	query := `
		SELECT
			sr.type AS label,
			COUNT(*) AS count
		FROM support_records sr
		WHERE sr.project_id = $1
	`
	args := []any{f.ProjectID}
	ix := 2

	if f.DateFrom != nil {
		query += fmt.Sprintf(" AND sr.provided_at >= $%d", ix)
		args = append(args, *f.DateFrom)
		ix++
	}
	if f.DateTo != nil {
		query += fmt.Sprintf(" AND sr.provided_at <= $%d", ix)
		args = append(args, *f.DateTo)
		ix++
	}

	query += " GROUP BY sr.type ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, fmt.Errorf("count consultations: %w", err)
	}
	return results, nil
}
```

Group 3 (IDP status) joins through places → states:

```sql
SELECT s.conflict_zone AS label, COUNT(DISTINCT p.id) AS count
FROM people p
JOIN places pl ON p.origin_place_id = pl.id
JOIN states s ON pl.state_id = s.id
WHERE p.project_id = $1
GROUP BY s.conflict_zone
```

Group 8 (age group) uses CASE for birth_date bucketing:

```sql
SELECT
  CASE
    WHEN age < 1 THEN 'infant'
    WHEN age < 3 THEN 'toddler'
    WHEN age < 6 THEN 'pre_school'
    WHEN age < 12 THEN 'middle_childhood'
    WHEN age < 14 THEN 'young_teen'
    WHEN age < 18 THEN 'teenager'
    WHEN age < 25 THEN 'young_adult'
    WHEN age < 35 THEN 'early_adult'
    WHEN age < 55 THEN 'middle_aged_adult'
    ELSE 'old_adult'
  END AS label,
  COUNT(*) AS count
FROM (
  SELECT EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) AS age
  FROM people p WHERE p.project_id = $1 AND p.birth_date IS NOT NULL
) sub
GROUP BY label ORDER BY label
```

Group 10 (family units) counts distinct households:

```sql
SELECT 'households' AS label, COUNT(DISTINCT h.id) AS count
FROM households h
JOIN household_members hm ON h.id = hm.household_id
JOIN people p ON hm.person_id = p.id
JOIN support_records sr ON sr.person_id = p.id
WHERE h.project_id = $1
```

**Step 2: Verify and commit**

Run: `just test`

---

### Task 17: Report usecase

**Files:**

- Create: `internal/usecase/report/report_usecase.go`
- Create: `internal/usecase/report/types.go`

**Step 1: Create report types**

```go
// internal/usecase/report/types.go
package report

// ReportInput is the query input for all report endpoints.
type ReportInput struct {
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}

// ReportOutput wraps a single report group result.
type ReportOutput struct {
	Group string             `json:"group"`
	Rows  []CountResultDTO   `json:"rows"`
	Total int                `json:"total"`
}

// CountResultDTO is the response row.
type CountResultDTO struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// FullReportOutput returns all 10 groups at once.
type FullReportOutput struct {
	Consultations  ReportOutput `json:"consultations"`
	BySex          ReportOutput `json:"by_sex"`
	ByIDPStatus    ReportOutput `json:"by_idp_status"`
	ByCategory     ReportOutput `json:"by_category"`
	ByRegion       ReportOutput `json:"by_region"`
	BySphere       ReportOutput `json:"by_sphere"`
	ByOffice       ReportOutput `json:"by_office"`
	ByAgeGroup     ReportOutput `json:"by_age_group"`
	ByTag          ReportOutput `json:"by_tag"`
	FamilyUnits    ReportOutput `json:"family_units"`
}
```

**Step 2: Create report usecase**

```go
// internal/usecase/report/report_usecase.go
package report

import (
	"context"
	"fmt"
	"time"

	domainreport "github.com/lbrty/observer/internal/domain/report"
	"github.com/lbrty/observer/internal/repository"
)

type ReportUseCase struct {
	repo repository.ReportRepository
}

func NewReportUseCase(repo repository.ReportRepository) *ReportUseCase {
	return &ReportUseCase{repo: repo}
}

func (uc *ReportUseCase) Generate(ctx context.Context, projectID string, input ReportInput) (*FullReportOutput, error) {
	f := domainreport.ReportFilter{ProjectID: projectID}

	if input.DateFrom != "" {
		t, err := time.Parse("2006-01-02", input.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("parse date_from: %w", err)
		}
		f.DateFrom = &t
	}
	if input.DateTo != "" {
		t, err := time.Parse("2006-01-02", input.DateTo)
		if err != nil {
			return nil, fmt.Errorf("parse date_to: %w", err)
		}
		f.DateTo = &t
	}

	out := &FullReportOutput{}

	consultations, err := uc.repo.CountConsultations(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("consultations report: %w", err)
	}
	out.Consultations = toReportOutput("consultations", consultations)

	bySex, err := uc.repo.CountBySex(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("sex report: %w", err)
	}
	out.BySex = toReportOutput("by_sex", bySex)

	byIDP, err := uc.repo.CountByIDPStatus(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("idp report: %w", err)
	}
	out.ByIDPStatus = toReportOutput("by_idp_status", byIDP)

	byCat, err := uc.repo.CountByCategory(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("category report: %w", err)
	}
	out.ByCategory = toReportOutput("by_category", byCat)

	byRegion, err := uc.repo.CountByCurrentRegion(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("region report: %w", err)
	}
	out.ByRegion = toReportOutput("by_region", byRegion)

	bySphere, err := uc.repo.CountBySphere(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("sphere report: %w", err)
	}
	out.BySphere = toReportOutput("by_sphere", bySphere)

	byOffice, err := uc.repo.CountByOffice(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("office report: %w", err)
	}
	out.ByOffice = toReportOutput("by_office", byOffice)

	byAge, err := uc.repo.CountByAgeGroup(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("age report: %w", err)
	}
	out.ByAgeGroup = toReportOutput("by_age_group", byAge)

	byTag, err := uc.repo.CountByTag(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("tag report: %w", err)
	}
	out.ByTag = toReportOutput("by_tag", byTag)

	families, err := uc.repo.CountFamilyUnits(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("family report: %w", err)
	}
	out.FamilyUnits = toReportOutput("family_units", families)

	return out, nil
}

func toReportOutput(group string, rows []domainreport.CountResult) ReportOutput {
	total := 0
	dtos := make([]CountResultDTO, len(rows))
	for i, r := range rows {
		dtos[i] = CountResultDTO{Label: r.Label, Count: r.Count}
		total += r.Count
	}
	return ReportOutput{Group: group, Rows: dtos, Total: total}
}
```

**Step 3: Verify and commit**

Run: `just test`

---

### Task 18: Report handler and routes

**Files:**

- Create: `internal/handler/report_handler.go`
- Modify: `internal/server/server.go` — add report routes
- Modify: `internal/app/container.go` — wire ReportUseCase + ReportRepository

**Step 1: Create report handler**

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucreport "github.com/lbrty/observer/internal/usecase/report"
)

type ReportHandler struct {
	uc *ucreport.ReportUseCase
}

func NewReportHandler(uc *ucreport.ReportUseCase) *ReportHandler {
	return &ReportHandler{uc: uc}
}

// Generate handles GET /projects/:project_id/reports.
func (h *ReportHandler) Generate(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucreport.ReportInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Generate(c.Request.Context(), projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, out)
}
```

**Step 2: Add route**

In `server.go`, in the project read group:

```go
readGroup.GET("/reports", reportHandler.Generate)
```

**Step 3: Wire in container**

Add `ReportRepository` and `ReportUseCase` to container.

**Step 4: Verify and commit**

Run: `just test`

---

### Task 19: Reports frontend — D3 charts and page

**Files:**

- Create: `packages/observer-web/src/routes/_app/projects/$projectId/reports/index.tsx`
- Create: `packages/observer-web/src/hooks/use-reports.ts`
- Create: `packages/observer-web/src/types/report.ts`
- Create: `packages/observer-web/src/components/charts/bar-chart.tsx`
- Create: `packages/observer-web/src/components/charts/pie-chart.tsx`
- Create: `packages/observer-web/src/components/charts/grouped-bar-chart.tsx`

**Step 1: Install D3**

Run: `cd packages/observer-web && bun add d3 @types/d3`

**Step 2: Create report types**

```ts
// types/report.ts
export interface CountResult {
  label: string;
  count: number;
}

export interface ReportGroup {
  group: string;
  rows: CountResult[];
  total: number;
}

export interface FullReport {
  consultations: ReportGroup;
  by_sex: ReportGroup;
  by_idp_status: ReportGroup;
  by_category: ReportGroup;
  by_region: ReportGroup;
  by_sphere: ReportGroup;
  by_office: ReportGroup;
  by_age_group: ReportGroup;
  by_tag: ReportGroup;
  family_units: ReportGroup;
}

export interface ReportParams {
  date_from?: string;
  date_to?: string;
}
```

**Step 3: Create report hook**

```ts
// hooks/use-reports.ts
import { useQuery, keepPreviousData } from "@tanstack/react-query";
import api from "@/lib/api";
import type { FullReport, ReportParams } from "@/types/report";

export function useReports(projectId: string, params?: ReportParams) {
  return useQuery({
    queryKey: ["reports", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/reports`, {
          searchParams: params as Record<string, string>,
        })
        .json<FullReport>(),
    placeholderData: keepPreviousData,
  });
}
```

**Step 4: Create D3 chart components**

Use D3 for rendering SVG charts that respect the Observer color system (CSS variables).

**BarChart** — horizontal bars for Groups 1, 3, 4, 5, 7, 9:

- Uses `d3.scaleBand()` for y-axis (labels), `d3.scaleLinear()` for x-axis (counts)
- Bars colored with `var(--accent)`, labels with `var(--fg)`
- Responsive: uses `useRef` + `ResizeObserver` to redraw on container resize
- Tooltip on hover showing exact count

**PieChart** — donut chart for Groups 2, 10 (sex breakdown, family units):

- Uses `d3.pie()` + `d3.arc()` with inner radius for donut style
- Color scale using `d3.scaleOrdinal()` with Observer semantic colors
- Legend below chart
- Center text showing total

**GroupedBarChart** — grouped vertical bars for Groups 6, 8 (sphere × type, age × type):

- Uses `d3.scaleBand()` with nested groups
- Each sub-group colored differently (legal, social, humanitarian)
- X-axis labels rotated for readability

All charts:

- Use `var(--fg)` for text, `var(--border)` for gridlines
- Support dark/light themes via CSS variables
- Animate on mount with D3 transitions
- Show "No data" placeholder when rows are empty

**Step 5: Create reports page**

Date range filter (date_from, date_to inputs) at the top. Grid of 10 chart cards, each with a title and the appropriate chart type:

| Report Group  | Chart Type      |
| ------------- | --------------- |
| Consultations | BarChart        |
| By Sex        | PieChart        |
| By IDP Status | BarChart        |
| By Category   | BarChart        |
| By Region     | BarChart        |
| By Sphere     | GroupedBarChart |
| By Office     | BarChart        |
| By Age Group  | GroupedBarChart |
| By Tag        | BarChart        |
| Family Units  | PieChart        |

**Step 6: Add project sidebar link for reports**

Update `project-sidebar.tsx` to include reports link with `ChartBar` icon from Phosphor.

**Step 7: Add i18n keys for reports**

Add `project.reports.*` keys with group titles and labels.

**Step 8: Verify and commit**

Run: `cd packages/observer-web && bun run build`

---

## Phase 4: Deferred Features (Documented for Future)

### Task 20: Document future work

**Files:**

- Create: `docs/adr/009-deferred-features.md`

Document the following as future work:

1. **Email service integration** — SMTP client for sending notifications (optional, not required for core flow)
2. **MFA/TOTP verification** — TOTP library integration, enrollment flow, verification endpoint
3. **Document file upload/download** — Storage abstraction (local FS / S3), multipart upload, binary download, MEK/DEK encryption
4. **Data export/import** — CSV/JSON batch export of people + support records, batch import with deduplication by external_id
5. **Advanced search/filtering** — Date range filters on list endpoints, multi-field sorting, saved filters
6. **Go → TypeScript type generation** — tygo integration per ADR-008

---

## Execution Order Summary

| Phase       | Tasks       | Description                                                        |
| ----------- | ----------- | ------------------------------------------------------------------ |
| **Phase 1** | Tasks 1–9   | Frontend project-scoped pages (backend APIs exist)                 |
| **Phase 2** | Tasks 10–13 | Auth: registration approval, profile, password change, admin reset |
| **Phase 3** | Tasks 14–19 | Reports: domain → repo → usecase → handler → D3 frontend           |
| **Phase 4** | Task 20     | Documentation of deferred work                                     |

Total: 20 tasks across 4 phases.
