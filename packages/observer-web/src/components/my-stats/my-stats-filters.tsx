import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { CaretDownIcon, CaretUpIcon, DownloadSimpleIcon, FunnelIcon } from "@/components/ui/icons";
import { FilterChip, FilterField, labelKeyMap } from "@/components/reports/shared";
import type { DatePreset } from "@/components/reports/shared";
import { ReportDatePresets } from "@/components/reports/report-date-presets";
import { UISelect } from "@/components/ui/ui-select";
import { useReport } from "@/hooks/reports/use-reports";
import { exportReportCSV } from "@/lib/export-csv";
import type { ReportParams } from "@/types/report";

interface MyStatsFiltersProps {
  projectId: string;
  params: ReportParams;
  setParams: (updater: (p: ReportParams) => ReportParams) => void;
  filtersOpen: boolean;
  setFiltersOpen: (updater: (o: boolean) => boolean) => void;
  activePreset: DatePreset | null;
  setActivePreset: (preset: DatePreset | null) => void;
  data: ReturnType<typeof useReport>["data"];
  supportTypeOptions: Array<{ label: string; value: string }>;
}

export function MyStatsFilters({
  projectId,
  params,
  setParams,
  filtersOpen,
  setFiltersOpen,
  activePreset,
  setActivePreset,
  data,
  supportTypeOptions,
}: MyStatsFiltersProps) {
  const { t } = useTranslation();

  const hasFilters = Object.entries(params).some(([, v]) => v != null && v !== "");
  const clearDatePreset = () => setActivePreset(null);

  return (
    <div className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary">
      <div className="flex items-center justify-between px-5 py-3">
        <h1 className="font-serif text-xl font-bold tracking-tight text-fg">
          {t("project.myStats.title")}
        </h1>
        <div className="flex items-center gap-2">
          {data && (
            <button
              type="button"
              onClick={() => exportReportCSV(data, projectId)}
              className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
            >
              <DownloadSimpleIcon size={14} />
              {t("project.reports.exportCsv")}
            </button>
          )}
          <button
            type="button"
            onClick={() => setFiltersOpen((o) => !o)}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
          >
            <FunnelIcon size={14} />
            {t("project.reports.toggleFilters")}
            {filtersOpen ? <CaretUpIcon size={12} /> : <CaretDownIcon size={12} />}
          </button>
        </div>
      </div>

      {filtersOpen && (
        <div className="border-t border-border-secondary px-5 pb-4 pt-3">
          <ReportDatePresets
            activePreset={activePreset}
            onSelect={(key, dates) => {
              setParams((p) => ({ ...p, ...dates }));
              setActivePreset(key);
            }}
          />

          <div className="grid grid-cols-3 gap-x-4 gap-y-3">
            <FilterField label={t("project.reports.dateFrom")}>
              <DatePicker
                value={params.date_from ?? ""}
                onChange={(v) => {
                  setParams((p) => ({ ...p, date_from: v || undefined }));
                  clearDatePreset();
                }}
              />
            </FilterField>
            <FilterField label={t("project.reports.dateTo")}>
              <DatePicker
                value={params.date_to ?? ""}
                onChange={(v) => {
                  setParams((p) => ({ ...p, date_to: v || undefined }));
                  clearDatePreset();
                }}
              />
            </FilterField>
            <FilterField label={t("project.reports.filterSupportType")}>
              <UISelect
                value={params.support_type ?? ""}
                onValueChange={(v) => setParams((p) => ({ ...p, support_type: v || undefined }))}
                options={[
                  { label: t("project.reports.allValues"), value: "" },
                  ...supportTypeOptions,
                ]}
                placeholder={t("project.reports.allValues")}
                fullWidth
              />
            </FilterField>
          </div>
        </div>
      )}

      {hasFilters && (
        <div className="flex flex-wrap items-center gap-1.5 border-t border-border-secondary px-5 py-2.5">
          {params.date_from && (
            <FilterChip
              label={t("project.reports.dateFrom")}
              value={params.date_from}
              onRemove={() => {
                setParams((p) => ({ ...p, date_from: undefined }));
                clearDatePreset();
              }}
            />
          )}
          {params.date_to && (
            <FilterChip
              label={t("project.reports.dateTo")}
              value={params.date_to}
              onRemove={() => {
                setParams((p) => ({ ...p, date_to: undefined }));
                clearDatePreset();
              }}
            />
          )}
          {params.support_type && (
            <FilterChip
              label={t("project.reports.filterSupportType")}
              value={
                supportTypeOptions.find((s) => s.value === params.support_type)?.label ??
                params.support_type
              }
              onRemove={() => setParams((p) => ({ ...p, support_type: undefined }))}
            />
          )}
          <button
            type="button"
            onClick={() => {
              setParams(() => ({}));
              clearDatePreset();
            }}
            className="ml-1 text-xs font-medium text-fg-tertiary underline transition-colors hover:text-fg"
          >
            {t("project.reports.clearAll")}
          </button>
        </div>
      )}
    </div>
  );
}
