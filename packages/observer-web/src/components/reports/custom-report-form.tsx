import { useTranslation } from "react-i18next";

import { DateRangePicker } from "@/components/date-picker";
import { UISelect } from "@/components/ui-select";
import { typeKeys } from "@/constants/support";
import type { CustomReportParams } from "@/types/report";

const METRICS = ["events", "people", "units", "pets"] as const;

const DIMENSIONS = [
  "sex",
  "age_group",
  "region",
  "conflict_zone",
  "office",
  "sphere",
  "category",
  "person_tag",
  "pet_tag",
  "pet_status",
] as const;

const SUPPORT_TYPE_OPTIONS = Object.keys(typeKeys);

export const DIMENSION_LABEL_KEYS: Record<string, string> = {
  sex: "project.customReport.dimSex",
  age_group: "project.customReport.dimAgeGroup",
  region: "project.customReport.dimRegion",
  conflict_zone: "project.customReport.dimConflictZone",
  office: "project.customReport.dimOffice",
  sphere: "project.customReport.dimSphere",
  category: "project.customReport.dimCategory",
  person_tag: "project.customReport.dimPersonTag",
  pet_tag: "project.customReport.dimPetTag",
  pet_status: "project.customReport.dimPetStatus",
};

interface CustomReportFormProps {
  metric: CustomReportParams["metric"];
  groupBy: string[];
  dateFrom: string;
  dateTo: string;
  supportType: string;
  isFetching: boolean;
  onMetricChange: (v: CustomReportParams["metric"]) => void;
  onToggleDimension: (dim: string) => void;
  onDateRangeChange: (from: string, to: string) => void;
  onSupportTypeChange: (v: string) => void;
  onGenerate: () => void;
}

export function CustomReportForm({
  metric,
  groupBy,
  dateFrom,
  dateTo,
  supportType,
  isFetching,
  onMetricChange,
  onToggleDimension,
  onDateRangeChange,
  onSupportTypeChange,
  onGenerate,
}: CustomReportFormProps) {
  const { t } = useTranslation();

  return (
    <div className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary">
      <div className="px-5 py-4">
        <h1 className="font-serif text-xl font-bold tracking-tight text-fg">
          {t("project.customReport.title")}
        </h1>
      </div>

      <div className="border-t border-border-secondary px-5 pb-5 pt-4 space-y-5">
        {/* Metric */}
        <div>
          <span className="mb-2 block text-xs font-medium text-fg-secondary">
            {t("project.customReport.metric")}
          </span>
          <div className="flex flex-wrap gap-2">
            {METRICS.map((m) => (
              <button
                key={m}
                type="button"
                onClick={() => onMetricChange(m)}
                className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                  metric === m
                    ? "bg-accent text-accent-fg"
                    : "bg-bg-tertiary text-fg-secondary hover:text-fg"
                }`}
              >
                {t(`project.customReport.metric_${m}`)}
              </button>
            ))}
          </div>
        </div>

        {/* Dimensions */}
        <div>
          <span className="mb-2 block text-xs font-medium text-fg-secondary">
            {t("project.customReport.dimensions")}
            <span className="ml-1 text-fg-tertiary">({groupBy.length}/2)</span>
          </span>
          <div className="flex flex-wrap gap-2">
            {DIMENSIONS.map((dim) => {
              const selected = groupBy.includes(dim);
              const disabled = !selected && groupBy.length >= 2;
              return (
                <button
                  key={dim}
                  type="button"
                  disabled={disabled}
                  onClick={() => onToggleDimension(dim)}
                  className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                    selected
                      ? "bg-accent text-accent-fg"
                      : disabled
                        ? "cursor-not-allowed bg-bg-tertiary text-fg-tertiary/50"
                        : "bg-bg-tertiary text-fg-secondary hover:text-fg"
                  }`}
                >
                  {t(DIMENSION_LABEL_KEYS[dim])}
                </button>
              );
            })}
          </div>
        </div>

        {/* Date range + support type */}
        <div className="flex flex-wrap items-end gap-4">
          <DateRangePicker
            from={dateFrom}
            to={dateTo}
            onChange={(range) => {
              onDateRangeChange(range.from ?? "", range.to ?? "");
            }}
          />
          <div className="min-w-56 space-y-1.5">
            <span className="block text-xs font-medium text-fg-secondary">
              {t("project.reports.filterSupportType")}
            </span>
            <UISelect
              value={supportType}
              onValueChange={onSupportTypeChange}
              options={[
                { label: t("project.reports.allValues"), value: "" },
                ...SUPPORT_TYPE_OPTIONS.map((s) => ({
                  label: t(typeKeys[s]),
                  value: s,
                })),
              ]}
              placeholder={t("project.reports.allValues")}
              fullWidth
            />
          </div>
        </div>

        {/* Generate button */}
        <button
          type="button"
          disabled={groupBy.length === 0 || isFetching}
          onClick={onGenerate}
          className="rounded-lg bg-accent px-5 py-2 text-sm font-semibold text-accent-fg transition-colors hover:bg-accent/90 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {isFetching ? t("project.customReport.generating") : t("project.customReport.generate")}
        </button>
      </div>
    </div>
  );
}
