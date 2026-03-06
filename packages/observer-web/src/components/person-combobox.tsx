import { useDeferredValue, useState } from "react";
import { useTranslation } from "react-i18next";

import { BaseCombobox } from "@/components/base-combobox";
import { useSearchPeople } from "@/hooks/use-people";
import type { Person } from "@/types/person";

interface PersonComboboxProps {
  projectId: string;
  excludeIds?: string[];
  onSelect: (person: Person) => void;
  placeholder?: string;
}

export function PersonCombobox({
  projectId,
  excludeIds = [],
  onSelect,
  placeholder,
}: PersonComboboxProps) {
  const { t } = useTranslation();
  const [search, setSearch] = useState("");
  const deferred = useDeferredValue(search);
  const { data, isLoading } = useSearchPeople(projectId, deferred);

  const filtered = data?.people.filter((p) => !excludeIds.includes(p.id)) ?? [];

  return (
    <BaseCombobox
      items={filtered}
      onSelect={onSelect}
      getItemKey={(p) => p.id}
      search={search}
      onSearchChange={setSearch}
      isLoading={isLoading}
      placeholder={placeholder ?? t("project.people.searchPeople")}
      noDataLabel={t("admin.common.noData")}
      listboxId="person-listbox"
      optionIdPrefix="person-option"
      renderItem={(p) => (
        <span className="flex flex-col">
          <span className="text-sm text-fg">
            {p.first_name} {p.last_name ?? ""}
          </span>
          <span className="font-mono text-xs text-fg-tertiary">{p.id}</span>
        </span>
      )}
    />
  );
}
