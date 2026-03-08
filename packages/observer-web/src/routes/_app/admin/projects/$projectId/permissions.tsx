import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { PencilSimpleIcon, TrashIcon, UsersIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { AssignDialog, EditDialog } from "@/components/permissions";
import { StatusBadge } from "@/components/status-badge";
import { UserInitials } from "@/components/user-initials";
import {
  useAssignPermission,
  usePermissions,
  useRevokePermission,
  useUpdatePermission,
} from "@/hooks/use-permissions";
import { useAuth } from "@/stores/auth";
import type { ProjectPermissionMember } from "@/types/permission";

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
        emptyState={
          <EmptyState
            icon={UsersIcon}
            title={t("admin.permissions.emptyTitle")}
            description={t("admin.permissions.emptyDescription")}
            action={
              <Button onClick={() => setAssignOpen(true)}>
                {t("admin.permissions.addMember")}
              </Button>
            }
          />
        }
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
