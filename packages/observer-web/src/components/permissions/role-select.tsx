import { useTranslation } from "react-i18next";

import { projectRoleKeys } from "@/constants/user";

export function useRoleOptions() {
  const { t } = useTranslation();
  return Object.entries(projectRoleKeys).map(([value, label]) => ({
    label: t(label),
    value,
  }));
}

export function RoleDescription({ role }: { role: string }) {
  const { t } = useTranslation();
  const key = `admin.permissions.role${role.charAt(0).toUpperCase() + role.slice(1)}Desc`;
  return <p className="text-xs text-fg-tertiary">{t(key)}</p>;
}
