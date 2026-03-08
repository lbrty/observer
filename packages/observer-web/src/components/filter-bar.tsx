import type { ReactNode } from "react";

import { MagnifyingGlassIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";

export interface SearchFilter {
  type: "search";
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
}

export interface SelectFilter {
  type: "select";
  value: string;
  onValueChange: (value: string) => void;
  options: { label: string; value: string }[];
  placeholder?: string;
}

export interface DateRangeFilter {
  type: "date-range";
  fromValue: string;
  toValue: string;
  onFromChange: (value: string) => void;
  onToChange: (value: string) => void;
  fromPlaceholder?: string;
  toPlaceholder?: string;
}

export type FilterDef = SearchFilter | SelectFilter | DateRangeFilter;

interface FilterBarProps {
  filters: FilterDef[];
  trailing?: ReactNode;
}

export function FilterBar({ filters, trailing }: FilterBarProps) {
  return (
    <div className="mb-4 flex gap-3">
      {filters.map((f, i) => {
        if (f.type === "search") {
          return (
            <div key={i} className="relative">
              <MagnifyingGlassIcon
                size={14}
                className="absolute top-1/2 left-3 -translate-y-1/2 text-fg-tertiary"
              />
              <input
                placeholder={f.placeholder}
                value={f.value}
                onChange={(e) => f.onChange(e.target.value)}
                className="rounded-lg border border-border-secondary bg-bg-secondary py-1.5 pr-3 pl-8 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
              />
            </div>
          );
        }
        if (f.type === "date-range") {
          return (
            <div key={i} className="flex items-center gap-1.5">
              <input
                type="date"
                value={f.fromValue}
                onChange={(e) => f.onFromChange(e.target.value)}
                placeholder={f.fromPlaceholder}
                className="rounded-lg border border-border-secondary bg-bg-secondary px-3 py-1.5 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
              />
              <span className="text-xs text-fg-tertiary">&ndash;</span>
              <input
                type="date"
                value={f.toValue}
                onChange={(e) => f.onToChange(e.target.value)}
                placeholder={f.toPlaceholder}
                className="rounded-lg border border-border-secondary bg-bg-secondary px-3 py-1.5 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
              />
            </div>
          );
        }
        return (
          <UISelect
            key={i}
            value={f.value}
            onValueChange={f.onValueChange}
            options={f.options}
            placeholder={f.placeholder}
          />
        );
      })}
      {trailing}
    </div>
  );
}
