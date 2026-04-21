import { useTranslation } from "react-i18next";

import { FilterChip } from "@/components/reports/shared";
import type { ReportParams } from "@/types/report";

interface ReportFilterChipsProps {
  params: ReportParams;
  officeLabel: string | undefined;
  categoryLabel: string | undefined;
  caseStatusLabel: string | undefined;
  sexLabel: string | undefined;
  ageGroupLabel: string | undefined;
  supportTypeLabel: string | undefined;
  onParamsChange: (updater: (prev: ReportParams) => ReportParams) => void;
  onClearDatePreset: () => void;
}

export function ReportFilterChips({
  params,
  officeLabel,
  categoryLabel,
  caseStatusLabel,
  sexLabel,
  ageGroupLabel,
  supportTypeLabel,
  onParamsChange,
  onClearDatePreset,
}: ReportFilterChipsProps) {
  const { t } = useTranslation();

  const hasFilters = Object.values(params).some((v) => v != null && v !== "");
  if (!hasFilters) return null;

  return (
    <div className="mt-3 flex flex-wrap items-center gap-1.5 border-t border-border-secondary pt-2.5">
      {params.date_from && (
        <FilterChip
          label={t("project.reports.dateFrom")}
          value={params.date_from}
          onRemove={() => {
            onParamsChange((p) => ({ ...p, date_from: undefined }));
            onClearDatePreset();
          }}
        />
      )}
      {params.date_to && (
        <FilterChip
          label={t("project.reports.dateTo")}
          value={params.date_to}
          onRemove={() => {
            onParamsChange((p) => ({ ...p, date_to: undefined }));
            onClearDatePreset();
          }}
        />
      )}
      {params.office_id && (
        <FilterChip
          label={t("project.reports.filterOffice")}
          value={officeLabel ?? params.office_id}
          onRemove={() => onParamsChange((p) => ({ ...p, office_id: undefined }))}
        />
      )}
      {params.category_id && (
        <FilterChip
          label={t("project.reports.filterCategory")}
          value={categoryLabel ?? params.category_id}
          onRemove={() => onParamsChange((p) => ({ ...p, category_id: undefined }))}
        />
      )}
      {params.case_status && (
        <FilterChip
          label={t("project.reports.filterCaseStatus")}
          value={caseStatusLabel ?? params.case_status}
          onRemove={() => onParamsChange((p) => ({ ...p, case_status: undefined }))}
        />
      )}
      {params.sex && (
        <FilterChip
          label={t("project.reports.filterSex")}
          value={sexLabel ?? params.sex}
          onRemove={() => onParamsChange((p) => ({ ...p, sex: undefined }))}
        />
      )}
      {params.age_group && (
        <FilterChip
          label={t("project.reports.filterAgeGroup")}
          value={ageGroupLabel ?? params.age_group}
          onRemove={() => onParamsChange((p) => ({ ...p, age_group: undefined }))}
        />
      )}
      {params.support_type && (
        <FilterChip
          label={t("project.reports.filterSupportType")}
          value={supportTypeLabel ?? params.support_type}
          onRemove={() => onParamsChange((p) => ({ ...p, support_type: undefined }))}
        />
      )}
      <button
        type="button"
        onClick={() => {
          onParamsChange(() => ({}));
          onClearDatePreset();
        }}
        className="ml-1 text-xs font-medium text-fg-tertiary underline transition-colors hover:text-fg"
      >
        {t("project.reports.clearAll")}
      </button>
    </div>
  );
}
