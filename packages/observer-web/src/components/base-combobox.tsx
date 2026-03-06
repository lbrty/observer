import { type ReactNode, useEffect, useRef, useState } from "react";

import { MagnifyingGlassIcon } from "@/components/icons";

interface BaseComboboxProps<T> {
  items: T[];
  onSelect: (item: T) => void;
  getItemKey: (item: T) => string;
  renderItem: (item: T, isActive: boolean) => ReactNode;
  search: string;
  onSearchChange: (value: string) => void;
  placeholder?: string;
  minChars?: number;
  isLoading?: boolean;
  noDataLabel?: string;
  listboxId?: string;
  optionIdPrefix?: string;
  inputClassName?: string;
  listClassName?: string;
  renderGroupHeader?: (item: T, index: number) => ReactNode;
}

export function BaseCombobox<T>({
  items,
  onSelect,
  getItemKey,
  renderItem,
  search,
  onSearchChange,
  placeholder,
  minChars = 2,
  isLoading = false,
  noDataLabel = "...",
  listboxId = "combobox-listbox",
  optionIdPrefix = "combobox-option",
  inputClassName = "block w-full rounded-lg border border-border-secondary bg-bg-secondary py-2 pr-3 pl-9 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg",
  listClassName = "absolute z-10 mt-1 max-h-48 w-full overflow-auto rounded-lg border border-border-secondary bg-bg-secondary shadow-elevated",
  renderGroupHeader,
}: BaseComboboxProps<T>) {
  const [activeIndex, setActiveIndex] = useState(-1);
  const listRef = useRef<HTMLDivElement>(null);
  const showDropdown = search.length >= minChars;

  useEffect(() => {
    setActiveIndex(-1);
  }, [items]);

  useEffect(() => {
    if (activeIndex < 0 || !listRef.current) return;
    const el = listRef.current.querySelector(`#${optionIdPrefix}-${activeIndex}`);
    el?.scrollIntoView({ block: "nearest" });
  }, [activeIndex, optionIdPrefix]);

  function select(item: T) {
    onSelect(item);
    onSearchChange("");
    setActiveIndex(-1);
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (!showDropdown || items.length === 0) return;

    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIndex((i) => (i < items.length - 1 ? i + 1 : 0));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIndex((i) => (i > 0 ? i - 1 : items.length - 1));
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (activeIndex >= 0 && activeIndex < items.length) {
        select(items[activeIndex]);
      }
    } else if (e.key === "Escape") {
      onSearchChange("");
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
          aria-activedescendant={activeIndex >= 0 ? `${optionIdPrefix}-${activeIndex}` : undefined}
          aria-autocomplete="list"
          aria-controls={listboxId}
          value={search}
          onChange={(e) => {
            onSearchChange(e.target.value);
            setActiveIndex(-1);
          }}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          className={inputClassName}
        />
      </div>
      {showDropdown && (
        <div ref={listRef} id={listboxId} role="listbox" className={listClassName}>
          {isLoading && <p className="px-3 py-2 text-sm text-fg-tertiary">...</p>}
          {!isLoading && items.length === 0 && (
            <p className="px-3 py-2 text-sm text-fg-tertiary">{noDataLabel}</p>
          )}
          {items.map((item, i) => (
            <div key={getItemKey(item)}>
              {renderGroupHeader?.(item, i)}
              <button
                id={`${optionIdPrefix}-${i}`}
                role="option"
                aria-selected={i === activeIndex}
                type="button"
                onClick={() => select(item)}
                className={`flex w-full cursor-pointer px-3 py-2 text-left ${
                  i === activeIndex ? "bg-bg-tertiary" : "hover:bg-bg-tertiary"
                }`}
              >
                {renderItem(item, i === activeIndex)}
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
