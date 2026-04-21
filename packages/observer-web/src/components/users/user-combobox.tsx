import { useDeferredValue, useState } from "react";
import { useTranslation } from "react-i18next";

import { BaseCombobox } from "@/components/base-combobox";
import { useSearchUsers } from "@/hooks/use-users";
import type { AdminUser } from "@/types/admin";

interface UserComboboxProps {
  excludeIds?: string[];
  onSelect: (user: AdminUser) => void;
}

export function UserCombobox({ excludeIds = [], onSelect }: UserComboboxProps) {
  const { t } = useTranslation();
  const [search, setSearch] = useState("");
  const deferred = useDeferredValue(search);
  const { data, isLoading } = useSearchUsers(deferred);

  const filtered = data?.users.filter((u) => !excludeIds.includes(u.id)) ?? [];

  return (
    <BaseCombobox
      items={filtered}
      onSelect={onSelect}
      getItemKey={(u) => u.id}
      search={search}
      onSearchChange={setSearch}
      isLoading={isLoading}
      placeholder={t("admin.permissions.searchUsers")}
      noDataLabel={t("admin.common.noData")}
      listboxId="user-listbox"
      optionIdPrefix="user-option"
      inputClassName="block w-full rounded-lg border border-border-secondary bg-bg py-2 pr-3 pl-9 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
      renderItem={(u) => (
        <span className="flex flex-col">
          <span className="text-sm text-fg">
            {u.first_name} {u.last_name}
          </span>
          <span className="text-xs text-fg-tertiary">{u.email}</span>
        </span>
      )}
    />
  );
}
