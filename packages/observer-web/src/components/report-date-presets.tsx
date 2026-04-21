import { useTranslation } from "react-i18next";

import { PRESET_KEYS } from "@/components/report";
import type { DatePreset } from "@/components/report";

interface ReportDatePresetsProps {
  activePreset: DatePreset | null;
  onSelect: (preset: DatePreset) => void;
}

export function ReportDatePresets({ activePreset, onSelect }: ReportDatePresetsProps) {
  const { t } = useTranslation();

  return (
    <div className="mb-3 flex flex-wrap gap-1.5">
      {PRESET_KEYS.map(({ key, i18n }) => (
        <button
          key={key}
          type="button"
          onClick={() => onSelect(key)}
          className={`rounded-md px-2.5 py-1 text-xs font-medium transition-colors ${
            activePreset === key
              ? "bg-accent text-accent-fg"
              : "bg-bg-tertiary text-fg-secondary hover:text-fg"
          }`}
        >
          {t(i18n)}
        </button>
      ))}
    </div>
  );
}
