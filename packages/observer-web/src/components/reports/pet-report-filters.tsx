import { useTranslation } from "react-i18next";

import { DateRangePicker } from "@/components/date-picker";
import { CaretDownIcon, CaretUpIcon, FunnelIcon, PrinterIcon } from "@/components/ui/icons";
import {
  FilterChip,
  FilterField,
  getPresetDates,
  PRESET_KEYS,
  type DatePreset,
} from "@/components/reports/shared";
import { UISelect } from "@/components/ui/ui-select";
import type { PetReportParams } from "@/types/report";

interface PetReportFiltersProps {
  params: PetReportParams;
  setParams: React.Dispatch<React.SetStateAction<PetReportParams>>;
  filtersOpen: boolean;
  setFiltersOpen: React.Dispatch<React.SetStateAction<boolean>>;
  activePreset: DatePreset | null;
  setActivePreset: React.Dispatch<React.SetStateAction<DatePreset | null>>;
  statusOptions: { label: string; value: string }[];
  data: unknown;
}

export function PetReportFilters({
  params,
  setParams,
  filtersOpen,
  setFiltersOpen,
  activePreset,
  setActivePreset,
  statusOptions,
  data,
}: PetReportFiltersProps) {
  const { t } = useTranslation();
  const hasFilters = Object.values(params).some((v) => v != null && v !== "");
  const clearDatePreset = () => setActivePreset(null);

  return (
    <>
      {/* Print-only header */}
      <div className="print-header hidden">
        <h1 className="text-lg font-bold">{t("project.petReports.title")}</h1>
        {params.date_from && (
          <p>
            {params.date_from} &mdash; {params.date_to ?? "..."}
          </p>
        )}
      </div>

      {/* Header + filter panel */}
      <div
        data-print-hide
        className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary"
      >
        <div className="flex items-center justify-between px-5 py-3">
          <h1 className="font-serif text-xl font-bold tracking-tight text-fg">
            {t("project.petReports.title")}
          </h1>
          <div className="flex items-center gap-2">
            {!!data && (
              <button
                type="button"
                onClick={() => window.print()}
                className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
              >
                <PrinterIcon size={14} />
                {t("project.reports.print")}
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
            <div className="mb-3 flex flex-wrap gap-1.5">
              {PRESET_KEYS.map(({ key, i18n }) => (
                <button
                  key={key}
                  type="button"
                  onClick={() => {
                    const dates = getPresetDates(key);
                    setParams((p) => ({ ...p, ...dates }));
                    setActivePreset(key);
                  }}
                  className={`rounded-md px-2.5 py-1 text-xs font-medium transition-colors ${activePreset === key ? "bg-accent text-accent-fg" : "bg-bg-tertiary text-fg-secondary hover:text-fg"}`}
                >
                  {t(i18n)}
                </button>
              ))}
            </div>

            <div className="flex flex-wrap items-start gap-4">
              <FilterField label={t("project.reports.dateRange")}>
                <DateRangePicker
                  from={params.date_from ?? ""}
                  to={params.date_to ?? ""}
                  onChange={(range) => {
                    setParams((p) => ({
                      ...p,
                      date_from: range.from || undefined,
                      date_to: range.to || undefined,
                    }));
                    clearDatePreset();
                  }}
                />
              </FilterField>
              <FilterField label={t("project.petReports.filterStatus")}>
                <UISelect
                  fullWidth
                  value={params.status ?? ""}
                  onValueChange={(v) => setParams((p) => ({ ...p, status: v || undefined }))}
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...statusOptions]}
                  placeholder={t("project.reports.allValues")}
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
            {params.status && (
              <FilterChip
                label={t("project.petReports.filterStatus")}
                value={statusOptions.find((s) => s.value === params.status)?.label ?? params.status}
                onRemove={() => setParams((p) => ({ ...p, status: undefined }))}
              />
            )}
            <button
              type="button"
              onClick={() => {
                setParams({});
                clearDatePreset();
              }}
              className="ml-1 text-xs font-medium text-fg-tertiary underline transition-colors hover:text-fg"
            >
              {t("project.reports.clearAll")}
            </button>
          </div>
        )}
      </div>
    </>
  );
}
