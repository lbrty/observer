import { useTranslation } from "react-i18next";

import { ReportDatePresets } from "@/components/report-date-presets";
import { FilterField, labelKeyMap } from "@/components/report";
import type { DatePreset } from "@/components/report";
import { DateRangePicker } from "@/components/date-picker";
import { UISelect } from "@/components/ui-select";
import { SEX_VALUES, AGE_GROUP_VALUES, CASE_STATUS_VALUES } from "@/constants/person";
import { useCategories } from "@/hooks/use-categories";
import { useOffices } from "@/hooks/use-offices";
import { toSelectOptions } from "@/lib/options";
import type { ReportParams } from "@/types/report";

import { ReportFilterChips } from "./report-filter-chips";

interface ReportFilterBarProps {
  params: ReportParams;
  activePreset: DatePreset | null;
  onParamsChange: (updater: (prev: ReportParams) => ReportParams) => void;
  onPresetChange: (preset: DatePreset | null) => void;
}

const SUPPORT_TYPE_OPTIONS = [
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
] as const;

export function ReportFilterBar({
  params,
  activePreset,
  onParamsChange,
  onPresetChange,
}: ReportFilterBarProps) {
  const { t } = useTranslation();
  const { data: offices } = useOffices();
  const { data: categories } = useCategories();

  const officeOptions = toSelectOptions(offices);
  const categoryOptions = toSelectOptions(categories);
  const supportTypeOptions = SUPPORT_TYPE_OPTIONS.map((s) => ({
    label: t(labelKeyMap[s] ?? s),
    value: s,
  }));
  const caseStatusOptions = CASE_STATUS_VALUES.map((s) => ({
    label: t(`project.people.${s}`),
    value: s,
  }));
  const sexOptions = SEX_VALUES.map((s) => ({
    label: t(`project.people.sex${s[0].toUpperCase()}${s.slice(1)}`),
    value: s,
  }));
  const ageGroupOptions = AGE_GROUP_VALUES.map((g) => ({
    label: t(labelKeyMap[g] ?? g),
    value: g,
  }));

  const clearDatePreset = () => onPresetChange(null);

  return (
    <>
      <ReportDatePresets
        activePreset={activePreset}
        onSelect={(preset, dates) => {
          onParamsChange((p) => ({ ...p, ...dates }));
          onPresetChange(preset);
        }}
      />

      <div className="flex flex-wrap items-start gap-4">
        <FilterField label={t("project.reports.dateRange")}>
          <DateRangePicker
            from={params.date_from ?? ""}
            to={params.date_to ?? ""}
            onChange={(range) => {
              onParamsChange((p) => ({
                ...p,
                date_from: range.from || undefined,
                date_to: range.to || undefined,
              }));
              clearDatePreset();
            }}
          />
        </FilterField>
        <FilterField label={t("project.reports.filterOffice")}>
          <UISelect
            fullWidth
            value={params.office_id ?? ""}
            onValueChange={(v) => onParamsChange((p) => ({ ...p, office_id: v || undefined }))}
            options={[{ label: t("project.reports.allValues"), value: "" }, ...officeOptions]}
            placeholder={t("project.reports.allValues")}
          />
        </FilterField>
        <FilterField label={t("project.reports.filterCategory")}>
          <UISelect
            fullWidth
            value={params.category_id ?? ""}
            onValueChange={(v) => onParamsChange((p) => ({ ...p, category_id: v || undefined }))}
            options={[{ label: t("project.reports.allValues"), value: "" }, ...categoryOptions]}
            placeholder={t("project.reports.allValues")}
          />
        </FilterField>
        <FilterField label={t("project.reports.filterCaseStatus")}>
          <UISelect
            fullWidth
            value={params.case_status ?? ""}
            onValueChange={(v) => onParamsChange((p) => ({ ...p, case_status: v || undefined }))}
            options={[{ label: t("project.reports.allValues"), value: "" }, ...caseStatusOptions]}
            placeholder={t("project.reports.allValues")}
          />
        </FilterField>
        <FilterField label={t("project.reports.filterSex")}>
          <UISelect
            fullWidth
            value={params.sex ?? ""}
            onValueChange={(v) => onParamsChange((p) => ({ ...p, sex: v || undefined }))}
            options={[{ label: t("project.reports.allValues"), value: "" }, ...sexOptions]}
            placeholder={t("project.reports.allValues")}
          />
        </FilterField>
        <FilterField label={t("project.reports.filterAgeGroup")}>
          <UISelect
            fullWidth
            value={params.age_group ?? ""}
            onValueChange={(v) => onParamsChange((p) => ({ ...p, age_group: v || undefined }))}
            options={[{ label: t("project.reports.allValues"), value: "" }, ...ageGroupOptions]}
            placeholder={t("project.reports.allValues")}
          />
        </FilterField>
        <FilterField label={t("project.reports.filterSupportType")}>
          <UISelect
            fullWidth
            value={params.support_type ?? ""}
            onValueChange={(v) => onParamsChange((p) => ({ ...p, support_type: v || undefined }))}
            options={[{ label: t("project.reports.allValues"), value: "" }, ...supportTypeOptions]}
            placeholder={t("project.reports.allValues")}
          />
        </FilterField>
      </div>

      <ReportFilterChips
        params={params}
        officeLabel={officeOptions.find((o) => o.value === params.office_id)?.label}
        categoryLabel={categoryOptions.find((c) => c.value === params.category_id)?.label}
        caseStatusLabel={caseStatusOptions.find((s) => s.value === params.case_status)?.label}
        sexLabel={sexOptions.find((s) => s.value === params.sex)?.label}
        ageGroupLabel={ageGroupOptions.find((g) => g.value === params.age_group)?.label}
        supportTypeLabel={supportTypeOptions.find((s) => s.value === params.support_type)?.label}
        onParamsChange={onParamsChange}
        onClearDatePreset={clearDatePreset}
      />
    </>
  );
}
