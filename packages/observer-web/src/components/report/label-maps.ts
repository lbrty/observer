import { useTranslation } from "react-i18next";

import { typeKeys, sphereKeys, referralKeys } from "@/constants/support";
import { sexKeys, ageGroupKeys, AGE_RANGE_MAP as AGE_RANGES } from "@/constants/person";
import type { CountResult } from "@/types/report";

export { AGE_RANGES as AGE_RANGE_MAP };

export const labelKeyMap: Record<string, string> = {
  ...typeKeys,
  ...sphereKeys,
  unspecified: "project.supportRecords.sphereOther",
  ...referralKeys,
  ...sexKeys,
  ...ageGroupKeys,
};

export function useTranslatedRows(rows: CountResult[]): CountResult[] {
  const { t } = useTranslation();
  const translated = rows.map((r) => {
    const key = labelKeyMap[r.label];
    return key ? { ...r, label: t(key) } : r;
  });
  const merged = new Map<string, number>();
  for (const r of translated) {
    merged.set(r.label, (merged.get(r.label) ?? 0) + r.count);
  }
  return Array.from(merged, ([label, count]) => ({ label, count }));
}

export type DatePreset = "month" | "quarter" | "year" | "all";

export function getPresetDates(preset: DatePreset): { date_from?: string; date_to?: string } {
  const now = new Date();
  const fmt = (d: Date) => d.toISOString().slice(0, 10);
  const today = fmt(now);

  switch (preset) {
    case "month": {
      const from = new Date(now.getFullYear(), now.getMonth(), 1);
      return { date_from: fmt(from), date_to: today };
    }
    case "quarter": {
      const qMonth = Math.floor(now.getMonth() / 3) * 3 - 3;
      const from = new Date(now.getFullYear(), qMonth, 1);
      const to = new Date(now.getFullYear(), qMonth + 3, 0);
      return { date_from: fmt(from), date_to: fmt(to) };
    }
    case "year": {
      const from = new Date(now.getFullYear(), 0, 1);
      return { date_from: fmt(from), date_to: today };
    }
    case "all":
      return { date_from: undefined, date_to: undefined };
  }
}

export const PRESET_KEYS: { key: DatePreset; i18n: string }[] = [
  { key: "month", i18n: "project.reports.presetMonth" },
  { key: "quarter", i18n: "project.reports.presetQuarter" },
  { key: "year", i18n: "project.reports.presetYear" },
  { key: "all", i18n: "project.reports.presetAll" },
];
