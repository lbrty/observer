import type { ReactNode } from "react";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DateRangePicker } from "@/components/date-picker";
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
  onSearch?: () => void;
}

export function FilterBar({ filters, trailing, onSearch }: FilterBarProps) {
  const { t } = useTranslation();
  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
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
                className="h-9 rounded-lg border border-border-secondary bg-bg-secondary pr-3 pl-8 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
              />
            </div>
          );
        }
        if (f.type === "date-range") {
          return (
            <DateRangePicker
              key={i}
              from={f.fromValue}
              to={f.toValue}
              onChange={(range) => {
                f.onFromChange(range.from ?? "");
                f.onToChange(range.to ?? "");
              }}
              placeholderFrom={f.fromPlaceholder}
              placeholderTo={f.toPlaceholder}
            />
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
      {onSearch && (
        <Button variant="secondary" onClick={onSearch}>
          {t("common.search")}
        </Button>
      )}
    </div>
  );
}
