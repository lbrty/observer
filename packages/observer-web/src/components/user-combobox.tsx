import { MagnifyingGlassIcon } from "@/components/icons";
import { useDeferredValue, useState } from "react";
import { useTranslation } from "react-i18next";

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
  const showDropdown = search.length >= 2;

  return (
    <div className="relative">
      <div className="relative">
        <MagnifyingGlassIcon
          size={16}
          className="absolute top-1/2 left-3 -translate-y-1/2 text-fg-tertiary"
        />
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder={t("admin.permissions.searchUsers")}
          className="block w-full rounded-md border border-border-secondary bg-bg py-2 pr-3 pl-9 text-sm text-fg outline-none focus:border-accent"
        />
      </div>
      {showDropdown && (
        <div className="absolute z-10 mt-1 max-h-48 w-full overflow-auto rounded-md border border-border-secondary bg-bg-secondary shadow-elevated">
          {isLoading && (
            <p className="px-3 py-2 text-sm text-fg-tertiary">...</p>
          )}
          {!isLoading && filtered.length === 0 && (
            <p className="px-3 py-2 text-sm text-fg-tertiary">
              {t("admin.common.noData")}
            </p>
          )}
          {filtered.map((u) => (
            <button
              key={u.id}
              type="button"
              onClick={() => {
                onSelect(u);
                setSearch("");
              }}
              className="flex w-full cursor-pointer flex-col px-3 py-2 text-left hover:bg-bg-tertiary"
            >
              <span className="text-sm text-fg">
                {u.first_name} {u.last_name}
              </span>
              <span className="text-xs text-fg-tertiary">{u.email}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
