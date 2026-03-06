import { useDeferredValue, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { MagnifyingGlassIcon } from "@/components/icons";
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
  const [activeIndex, setActiveIndex] = useState(-1);
  const listRef = useRef<HTMLDivElement>(null);
  const deferred = useDeferredValue(search);
  const { data, isLoading } = useSearchPeople(projectId, deferred);

  const filtered = data?.people.filter((p) => !excludeIds.includes(p.id)) ?? [];
  const showDropdown = search.length >= 2;

  function select(person: Person) {
    onSelect(person);
    setSearch("");
    setActiveIndex(-1);
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (!showDropdown || filtered.length === 0) return;

    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIndex((i) => (i < filtered.length - 1 ? i + 1 : 0));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIndex((i) => (i > 0 ? i - 1 : filtered.length - 1));
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (activeIndex >= 0 && activeIndex < filtered.length) {
        select(filtered[activeIndex]);
      }
    } else if (e.key === "Escape") {
      setSearch("");
      setActiveIndex(-1);
    }
  }

  return (
    <div className="relative">
      <div className="relative">
        <MagnifyingGlassIcon
          size={16}
          className="absolute top-1/2 left-3 -translate-y-1/2 text-fg-tertiary"
        />
        <input
          type="text"
          role="combobox"
          aria-expanded={showDropdown}
          aria-activedescendant={activeIndex >= 0 ? `person-option-${activeIndex}` : undefined}
          aria-autocomplete="list"
          aria-controls="person-listbox"
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setActiveIndex(-1);
          }}
          onKeyDown={handleKeyDown}
          placeholder={placeholder ?? t("project.people.searchPeople")}
          className="block w-full rounded-lg border border-border-secondary bg-bg-secondary py-2 pr-3 pl-9 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </div>
      {showDropdown && (
        <div
          ref={listRef}
          id="person-listbox"
          role="listbox"
          className="absolute z-10 mt-1 max-h-48 w-full overflow-auto rounded-lg border border-border-secondary bg-bg-secondary shadow-elevated"
        >
          {isLoading && <p className="px-3 py-2 text-sm text-fg-tertiary">...</p>}
          {!isLoading && filtered.length === 0 && (
            <p className="px-3 py-2 text-sm text-fg-tertiary">{t("admin.common.noData")}</p>
          )}
          {filtered.map((p, i) => (
            <button
              key={p.id}
              id={`person-option-${i}`}
              role="option"
              aria-selected={i === activeIndex}
              type="button"
              onClick={() => select(p)}
              className={`flex w-full cursor-pointer flex-col px-3 py-2 text-left ${
                i === activeIndex ? "bg-bg-tertiary" : "hover:bg-bg-tertiary"
              }`}
            >
              <span className="text-sm text-fg">
                {p.first_name} {p.last_name ?? ""}
              </span>
              <span className="font-mono text-xs text-fg-tertiary">{p.id}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
