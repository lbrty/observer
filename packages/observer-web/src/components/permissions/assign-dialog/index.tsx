import { type SyntheticEvent, useState } from "react";

import { Dialog } from "@base-ui/react/dialog";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { UISelect } from "@/components/ui/ui-select";
import { UserCombobox } from "@/components/users/user-combobox";
import type { AdminUser } from "@/types/admin";
import type { ProjectRole } from "@/types/permission";

import { RoleDescription, useRoleOptions } from "../role-select";
import { PermissionToggleRow } from "./permission-toggle-row";
import { SelectedUserCard } from "./selected-user-card";

export function AssignDialog({
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
    can_export: boolean;
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
  const [canExport, setCanExport] = useState(false);

  function handleSelectUser(user: AdminUser | null) {
    setSelectedUser(user);
    setCanExport(user?.role === "admin" || user?.role === "staff");
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!selectedUser) return;
    await onSubmit({
      user_id: selectedUser.id,
      role,
      can_view_contact: contact,
      can_view_personal: personal,
      can_view_documents: documents,
      can_export: canExport,
    });
    handleSelectUser(null);
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
                <SelectedUserCard user={selectedUser} onClear={() => handleSelectUser(null)} />
              ) : (
                <UserCombobox excludeIds={excludeIds} onSelect={handleSelectUser} />
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
              <PermissionToggleRow
                checked={contact}
                onCheckedChange={setContact}
                label={t("admin.permissions.contactAccess")}
                description={t("admin.permissions.contactAccessDesc")}
              />
              <PermissionToggleRow
                checked={personal}
                onCheckedChange={setPersonal}
                label={t("admin.permissions.personalAccess")}
                description={t("admin.permissions.personalAccessDesc")}
              />
              <PermissionToggleRow
                checked={documents}
                onCheckedChange={setDocuments}
                label={t("admin.permissions.documentAccess")}
                description={t("admin.permissions.documentAccessDesc")}
              />
              <PermissionToggleRow
                checked={canExport}
                onCheckedChange={setCanExport}
                label={t("admin.permissions.exportAccess")}
                description={t("admin.permissions.exportAccessDesc")}
              />
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
