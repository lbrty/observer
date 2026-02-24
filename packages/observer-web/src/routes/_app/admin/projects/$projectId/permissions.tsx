import { Dialog } from "@base-ui/react/dialog";
import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { StatusBadge } from "@/components/status-badge";
import {
  useAssignPermission,
  usePermissions,
  useRevokePermission,
  useUpdatePermission,
} from "@/hooks/use-permissions";
import type { ProjectPermission, ProjectRole } from "@/types/permission";

export const Route = createFileRoute(
  "/_app/admin/projects/$projectId/permissions",
)({
  component: PermissionsPage,
});

function PermissionsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const { data, isLoading } = usePermissions(projectId);
  const assignPermission = useAssignPermission();
  const updatePermission = useUpdatePermission();
  const revokePermission = useRevokePermission();

  const [assignOpen, setAssignOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<ProjectPermission | null>(null);
  const [revokeTarget, setRevokeTarget] = useState<ProjectPermission | null>(
    null,
  );

  const columns: Column<ProjectPermission>[] = [
    {
      key: "user_id",
      header: t("admin.permissions.userId"),
      render: (p) => (
        <span className="font-mono text-xs text-fg-secondary">{p.user_id}</span>
      ),
    },
    {
      key: "role",
      header: t("admin.permissions.role"),
      render: (p) => <StatusBadge label={p.role} />,
    },
    {
      key: "contact",
      header: t("admin.permissions.contact"),
      render: (p) => (
        <StatusBadge
          label={
            p.can_view_contact ? t("admin.users.yes") : t("admin.users.no")
          }
          variant={p.can_view_contact ? "foam" : "neutral"}
        />
      ),
    },
    {
      key: "personal",
      header: t("admin.permissions.personal"),
      render: (p) => (
        <StatusBadge
          label={
            p.can_view_personal ? t("admin.users.yes") : t("admin.users.no")
          }
          variant={p.can_view_personal ? "foam" : "neutral"}
        />
      ),
    },
    {
      key: "documents",
      header: t("admin.permissions.documents"),
      render: (p) => (
        <StatusBadge
          label={
            p.can_view_documents ? t("admin.users.yes") : t("admin.users.no")
          }
          variant={p.can_view_documents ? "foam" : "neutral"}
        />
      ),
    },
    {
      key: "actions",
      header: t("admin.common.actions"),
      render: (p) => (
        <div className="flex gap-2">
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setEditTarget(p);
            }}
            className="cursor-pointer text-xs text-accent hover:underline"
          >
            {t("admin.common.edit")}
          </button>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setRevokeTarget(p);
            }}
            className="cursor-pointer text-xs text-rose hover:underline"
          >
            {t("admin.permissions.revoke")}
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.permissions.title")}
        action={
          <button
            type="button"
            onClick={() => setAssignOpen(true)}
            className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90"
          >
            {t("admin.permissions.assign")}
          </button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.permissions ?? []}
        keyExtractor={(p) => p.id}
        isLoading={isLoading}
      />

      <AssignDialog
        open={assignOpen}
        onOpenChange={setAssignOpen}
        projectId={projectId}
        onSubmit={async (form) => {
          await assignPermission.mutateAsync({ projectId, data: form });
          setAssignOpen(false);
        }}
        loading={assignPermission.isPending}
      />

      {editTarget && (
        <EditDialog
          open={!!editTarget}
          onOpenChange={(open) => !open && setEditTarget(null)}
          permission={editTarget}
          onSubmit={async (form) => {
            await updatePermission.mutateAsync({
              projectId,
              id: editTarget.id,
              data: form,
            });
            setEditTarget(null);
          }}
          loading={updatePermission.isPending}
        />
      )}

      <ConfirmDialog
        open={!!revokeTarget}
        onOpenChange={(open) => !open && setRevokeTarget(null)}
        title={t("admin.permissions.revoke")}
        description={t("admin.permissions.revokeConfirm")}
        onConfirm={async () => {
          if (revokeTarget) {
            await revokePermission.mutateAsync({
              projectId,
              id: revokeTarget.id,
            });
            setRevokeTarget(null);
          }
        }}
        loading={revokePermission.isPending}
      />
    </div>
  );
}

function AssignDialog({
  open,
  onOpenChange,
  onSubmit,
  loading,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  onSubmit: (data: {
    user_id: string;
    role: ProjectRole;
    can_view_contact: boolean;
    can_view_personal: boolean;
    can_view_documents: boolean;
  }) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const [userId, setUserId] = useState("");
  const [role, setRole] = useState<ProjectRole>("viewer");
  const [contact, setContact] = useState(false);
  const [personal, setPersonal] = useState(false);
  const [documents, setDocuments] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await onSubmit({
      user_id: userId,
      role,
      can_view_contact: contact,
      can_view_personal: personal,
      can_view_documents: documents,
    });
    setUserId("");
    setRole("viewer");
    setContact(false);
    setPersonal(false);
    setDocuments(false);
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/40" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-sm -translate-x-1/2 -translate-y-1/2 rounded-lg border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="text-lg font-semibold text-fg">
            {t("admin.permissions.assignTitle")}
          </Dialog.Title>
          <form onSubmit={handleSubmit} className="mt-4 space-y-3">
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.permissions.userId")}
              </Field.Label>
              <Field.Control
                required
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
              />
            </Field.Root>
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.permissions.role")}
              </Field.Label>
              <select
                value={role}
                onChange={(e) => setRole(e.target.value as ProjectRole)}
                className="block w-full rounded-md border border-border-secondary bg-bg pl-3 pr-1 py-2 text-sm text-fg outline-none"
              >
                <option value="owner">owner</option>
                <option value="manager">manager</option>
                <option value="consultant">consultant</option>
                <option value="viewer">viewer</option>
              </select>
            </Field.Root>
            <div className="flex flex-col gap-2">
              <label className="flex items-center gap-2 text-sm text-fg-secondary">
                <input
                  type="checkbox"
                  checked={contact}
                  onChange={(e) => setContact(e.target.checked)}
                  className="accent-accent"
                />
                {t("admin.permissions.contact")}
              </label>
              <label className="flex items-center gap-2 text-sm text-fg-secondary">
                <input
                  type="checkbox"
                  checked={personal}
                  onChange={(e) => setPersonal(e.target.checked)}
                  className="accent-accent"
                />
                {t("admin.permissions.personal")}
              </label>
              <label className="flex items-center gap-2 text-sm text-fg-secondary">
                <input
                  type="checkbox"
                  checked={documents}
                  onChange={(e) => setDocuments(e.target.checked)}
                  className="accent-accent"
                />
                {t("admin.permissions.documents")}
              </label>
            </div>
            <div className="flex justify-end gap-3 pt-2">
              <Dialog.Close className="cursor-pointer rounded-md border border-border-secondary px-3 py-1.5 text-sm text-fg-secondary hover:bg-bg-tertiary">
                {t("admin.common.cancel")}
              </Dialog.Close>
              <button
                type="submit"
                disabled={loading}
                className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90 disabled:opacity-50"
              >
                {loading
                  ? t("admin.permissions.saving")
                  : t("admin.permissions.save")}
              </button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}

function EditDialog({
  open,
  onOpenChange,
  permission,
  onSubmit,
  loading,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  permission: ProjectPermission;
  onSubmit: (data: {
    role?: ProjectRole;
    can_view_contact?: boolean;
    can_view_personal?: boolean;
    can_view_documents?: boolean;
  }) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const [role, setRole] = useState<ProjectRole>(permission.role);
  const [contact, setContact] = useState(permission.can_view_contact);
  const [personal, setPersonal] = useState(permission.can_view_personal);
  const [documents, setDocuments] = useState(permission.can_view_documents);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await onSubmit({
      role,
      can_view_contact: contact,
      can_view_personal: personal,
      can_view_documents: documents,
    });
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/40" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-sm -translate-x-1/2 -translate-y-1/2 rounded-lg border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="text-lg font-semibold text-fg">
            {t("admin.permissions.editTitle")}
          </Dialog.Title>
          <form onSubmit={handleSubmit} className="mt-4 space-y-3">
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.permissions.role")}
              </Field.Label>
              <select
                value={role}
                onChange={(e) => setRole(e.target.value as ProjectRole)}
                className="block w-full rounded-md border border-border-secondary bg-bg pl-3 pr-1 py-2 text-sm text-fg outline-none"
              >
                <option value="owner">owner</option>
                <option value="manager">manager</option>
                <option value="consultant">consultant</option>
                <option value="viewer">viewer</option>
              </select>
            </Field.Root>
            <div className="flex flex-col gap-2">
              <label className="flex items-center gap-2 text-sm text-fg-secondary">
                <input
                  type="checkbox"
                  checked={contact}
                  onChange={(e) => setContact(e.target.checked)}
                  className="accent-accent"
                />
                {t("admin.permissions.contact")}
              </label>
              <label className="flex items-center gap-2 text-sm text-fg-secondary">
                <input
                  type="checkbox"
                  checked={personal}
                  onChange={(e) => setPersonal(e.target.checked)}
                  className="accent-accent"
                />
                {t("admin.permissions.personal")}
              </label>
              <label className="flex items-center gap-2 text-sm text-fg-secondary">
                <input
                  type="checkbox"
                  checked={documents}
                  onChange={(e) => setDocuments(e.target.checked)}
                  className="accent-accent"
                />
                {t("admin.permissions.documents")}
              </label>
            </div>
            <div className="flex justify-end gap-3 pt-2">
              <Dialog.Close className="cursor-pointer rounded-md border border-border-secondary px-3 py-1.5 text-sm text-fg-secondary hover:bg-bg-tertiary">
                {t("admin.common.cancel")}
              </Dialog.Close>
              <button
                type="submit"
                disabled={loading}
                className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90 disabled:opacity-50"
              >
                {loading
                  ? t("admin.permissions.saving")
                  : t("admin.permissions.save")}
              </button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
