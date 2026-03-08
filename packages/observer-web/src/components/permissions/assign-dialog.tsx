import { Dialog } from "@base-ui/react/dialog";
import { type SyntheticEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { XIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { UserCombobox } from "@/components/user-combobox";
import type { AdminUser } from "@/types/admin";
import type { ProjectRole } from "@/types/permission";

import { RoleDescription, useRoleOptions } from "./role-select";

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

  async function handleSubmit(e: SyntheticEvent) {
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
