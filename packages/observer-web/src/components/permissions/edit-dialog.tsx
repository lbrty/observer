import { Dialog } from "@base-ui/react/dialog";
import { type SyntheticEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { UserInitials } from "@/components/user-initials";
import type { ProjectPermissionMember, ProjectRole } from "@/types/permission";

import { RoleDescription, useRoleOptions } from "./role-select";

export function EditDialog({
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

  async function handleSubmit(e: SyntheticEvent) {
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
                <p className="ml-11.5 text-xs text-fg-tertiary">
                  {t("admin.permissions.contactAccessDesc")}
                </p>
              </div>
              <div className="space-y-2">
                <UISwitch
                  checked={personal}
                  onCheckedChange={setPersonal}
                  label={t("admin.permissions.personalAccess")}
                />
                <p className="ml-11.5 text-xs text-fg-tertiary">
                  {t("admin.permissions.personalAccessDesc")}
                </p>
              </div>
              <div className="space-y-2">
                <UISwitch
                  checked={documents}
                  onCheckedChange={setDocuments}
                  label={t("admin.permissions.documentAccess")}
                />
                <p className="ml-11.5 text-xs text-fg-tertiary">
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
