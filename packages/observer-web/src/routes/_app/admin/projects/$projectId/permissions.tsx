import { Dialog } from "@base-ui/react/dialog";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { PencilSimpleIcon, TrashIcon, XIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { StatusBadge } from "@/components/status-badge";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { UserCombobox } from "@/components/user-combobox";
import { UserInitials } from "@/components/user-initials";
import {
  useAssignPermission,
  usePermissions,
  useRevokePermission,
  useUpdatePermission,
} from "@/hooks/use-permissions";
import { useAuth } from "@/stores/auth";
import type { AdminUser } from "@/types/admin";
import type { ProjectPermissionMember, ProjectRole } from "@/types/permission";

export const Route = createFileRoute("/_app/admin/projects/$projectId/permissions")({
  component: PermissionsPage,
});

function PermissionsPage() {
  const { t } = useTranslation();
  const { user } = useAuth();
  const { projectId } = Route.useParams();
  const { data, isLoading } = usePermissions(projectId);
  const assignPermission = useAssignPermission();
  const updatePermission = useUpdatePermission();
  const revokePermission = useRevokePermission();

  const [assignOpen, setAssignOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<ProjectPermissionMember | null>(null);
  const [revokeTarget, setRevokeTarget] = useState<ProjectPermissionMember | null>(null);

  const isStaff = user?.role === "staff";

  const columns: Column<ProjectPermissionMember>[] = [
    {
      key: "member",
      header: t("admin.permissions.member"),
      render: (p) => (
        <div className="flex items-center gap-3">
          <UserInitials firstName={p.user_first_name} lastName={p.user_last_name} />
          <div className="min-w-0">
            <p className="truncate text-sm font-medium text-fg">
              {p.user_first_name} {p.user_last_name}
            </p>
            <p className="truncate text-xs text-fg-tertiary">{p.user_email}</p>
          </div>
        </div>
      ),
    },
    {
      key: "platformRole",
      header: t("admin.permissions.platformRole"),
      render: (p) => <StatusBadge label={p.user_role} />,
    },
    {
      key: "projectRole",
      header: t("admin.permissions.projectRole"),
      render: (p) => <StatusBadge label={p.role} />,
    },
    {
      key: "access",
      header: t("admin.permissions.access"),
      render: (p) => (
        <div className="flex gap-1">
          {p.can_view_contact && (
            <StatusBadge label={t("admin.permissions.contact")} variant="foam" />
          )}
          {p.can_view_personal && (
            <StatusBadge label={t("admin.permissions.personal")} variant="foam" />
          )}
          {p.can_view_documents && (
            <StatusBadge label={t("admin.permissions.documents")} variant="foam" />
          )}
          {!p.can_view_contact && !p.can_view_personal && !p.can_view_documents && (
            <span className="text-xs text-fg-tertiary">-</span>
          )}
        </div>
      ),
    },
    {
      key: "actions",
      header: t("admin.common.actions"),
      render: (p) => {
        const restricted = isStaff && p.user_role === "admin";
        if (restricted) return null;
        return (
          <div className="flex gap-2">
            <Button
              variant="ghost"
              className="p-1"
              onClick={(e) => {
                e.stopPropagation();
                setEditTarget(p);
              }}
            >
              <PencilSimpleIcon size={16} />
            </Button>
            <Button
              variant="ghost"
              className="p-1 hover:text-rose"
              onClick={(e) => {
                e.stopPropagation();
                setRevokeTarget(p);
              }}
            >
              <TrashIcon size={16} />
            </Button>
          </div>
        );
      },
    },
  ];

  const existingUserIds = data?.permissions.map((p) => p.user_id) ?? [];

  return (
    <div>
      <PageHeader
        title={t("admin.permissions.title")}
        action={
          <Button onClick={() => setAssignOpen(true)}>
            {t("admin.permissions.addMember")}
          </Button>
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
        excludeIds={existingUserIds}
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
        description={t("admin.permissions.revokeConfirm", {
          name: revokeTarget
            ? `${revokeTarget.user_first_name} ${revokeTarget.user_last_name}`
            : "",
        })}
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

function useRoleOptions() {
  const { t } = useTranslation();
  return [
    { label: t("admin.permissions.roleOwner"), value: "owner" },
    { label: t("admin.permissions.roleManager"), value: "manager" },
    { label: t("admin.permissions.roleConsultant"), value: "consultant" },
    { label: t("admin.permissions.roleViewer"), value: "viewer" },
  ];
}

function RoleDescription({ role }: { role: string }) {
  const { t } = useTranslation();
  const key = `admin.permissions.role${role.charAt(0).toUpperCase() + role.slice(1)}Desc`;
  return <p className="text-xs text-fg-tertiary">{t(key)}</p>;
}

function AssignDialog({
  open,
  onOpenChange,
  excludeIds,
  onSubmit,
  loading,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  excludeIds: string[];
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
  const roleOptions = useRoleOptions();
  const [selectedUser, setSelectedUser] = useState<AdminUser | null>(null);
  const [role, setRole] = useState<ProjectRole>("viewer");
  const [contact, setContact] = useState(false);
  const [personal, setPersonal] = useState(false);
  const [documents, setDocuments] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!selectedUser) return;
    await onSubmit({
      user_id: selectedUser.id,
      role,
      can_view_contact: contact,
      can_view_personal: personal,
      can_view_documents: documents,
    });
    setSelectedUser(null);
    setRole("viewer");
    setContact(false);
    setPersonal(false);
    setDocuments(false);
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/40 backdrop-blur-xs" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-2xl border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="text-lg font-semibold text-fg">
            {t("admin.permissions.addMemberTitle")}
          </Dialog.Title>
          <form onSubmit={handleSubmit} className="mt-4 space-y-4">
            <div>
              <label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.permissions.searchUsers")}
              </label>
              {selectedUser ? (
                <div className="flex items-center justify-between rounded-lg border border-border-secondary bg-bg px-3 py-2">
                  <div>
                    <p className="text-sm text-fg">
                      {selectedUser.first_name} {selectedUser.last_name}
                    </p>
                    <p className="text-xs text-fg-tertiary">{selectedUser.email}</p>
                  </div>
                  <button
                    type="button"
                    onClick={() => setSelectedUser(null)}
                    className="cursor-pointer rounded p-1 text-fg-tertiary hover:text-fg"
                  >
                    <XIcon size={14} />
                  </button>
                </div>
              ) : (
                <UserCombobox excludeIds={excludeIds} onSelect={setSelectedUser} />
              )}
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.permissions.projectRole")}
              </label>
              <UISelect
                value={role}
                onValueChange={(v) => setRole(v as ProjectRole)}
                options={roleOptions}
                fullWidth
              />
              <RoleDescription role={role} />
            </div>

            <div className="space-y-3">
              <p className="text-sm font-medium text-fg-secondary">
                {t("admin.permissions.access")}
              </p>
              <div className="space-y-2">
                <UISwitch
                  checked={contact}
                  onCheckedChange={setContact}
                  label={t("admin.permissions.contactAccess")}
                />
                <p className="ml-[46px] text-xs text-fg-tertiary">
                  {t("admin.permissions.contactAccessDesc")}
                </p>
              </div>
              <div className="space-y-2">
                <UISwitch
                  checked={personal}
                  onCheckedChange={setPersonal}
                  label={t("admin.permissions.personalAccess")}
                />
                <p className="ml-[46px] text-xs text-fg-tertiary">
                  {t("admin.permissions.personalAccessDesc")}
                </p>
              </div>
              <div className="space-y-2">
                <UISwitch
                  checked={documents}
                  onCheckedChange={setDocuments}
                  label={t("admin.permissions.documentAccess")}
                />
                <p className="ml-[46px] text-xs text-fg-tertiary">
                  {t("admin.permissions.documentAccessDesc")}
                </p>
              </div>
            </div>

            <div className="flex justify-end gap-3 pt-2">
              <Button variant="secondary" asChild>
                <Dialog.Close>{t("admin.common.cancel")}</Dialog.Close>
              </Button>
              <Button type="submit" disabled={loading || !selectedUser}>
                {loading ? t("admin.permissions.saving") : t("admin.permissions.save")}
              </Button>
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
  permission: ProjectPermissionMember;
  onSubmit: (data: {
    role?: ProjectRole;
    can_view_contact?: boolean;
    can_view_personal?: boolean;
    can_view_documents?: boolean;
  }) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const roleOptions = useRoleOptions();
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
        <Dialog.Backdrop className="fixed inset-0 bg-black/40 backdrop-blur-xs" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-2xl border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="text-lg font-semibold text-fg">
            {t("admin.permissions.editMemberTitle")}
          </Dialog.Title>

          <div className="mt-3 flex items-center gap-3 rounded-lg border border-border-secondary bg-bg px-3 py-2">
            <UserInitials firstName={permission.user_first_name} lastName={permission.user_last_name} />
            <div>
              <p className="text-sm font-medium text-fg">
                {permission.user_first_name} {permission.user_last_name}
              </p>
              <p className="text-xs text-fg-tertiary">{permission.user_email}</p>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="mt-4 space-y-4">
            <div>
              <label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.permissions.projectRole")}
              </label>
              <UISelect
                value={role}
                onValueChange={(v) => setRole(v as ProjectRole)}
                options={roleOptions}
                fullWidth
              />
              <RoleDescription role={role} />
            </div>

            <div className="space-y-3">
              <p className="text-sm font-medium text-fg-secondary">
                {t("admin.permissions.access")}
              </p>
              <div className="space-y-2">
                <UISwitch
                  checked={contact}
                  onCheckedChange={setContact}
                  label={t("admin.permissions.contactAccess")}
                />
                <p className="ml-[46px] text-xs text-fg-tertiary">
                  {t("admin.permissions.contactAccessDesc")}
                </p>
              </div>
              <div className="space-y-2">
                <UISwitch
                  checked={personal}
                  onCheckedChange={setPersonal}
                  label={t("admin.permissions.personalAccess")}
                />
                <p className="ml-[46px] text-xs text-fg-tertiary">
                  {t("admin.permissions.personalAccessDesc")}
                </p>
              </div>
              <div className="space-y-2">
                <UISwitch
                  checked={documents}
                  onCheckedChange={setDocuments}
                  label={t("admin.permissions.documentAccess")}
                />
                <p className="ml-[46px] text-xs text-fg-tertiary">
                  {t("admin.permissions.documentAccessDesc")}
                </p>
              </div>
            </div>

            <div className="flex justify-end gap-3 pt-2">
              <Button variant="secondary" asChild>
                <Dialog.Close>{t("admin.common.cancel")}</Dialog.Close>
              </Button>
              <Button type="submit" disabled={loading}>
                {loading ? t("admin.permissions.saving") : t("admin.permissions.save")}
              </Button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
