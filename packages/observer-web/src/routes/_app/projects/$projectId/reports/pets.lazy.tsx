import { useState } from "react";

import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { BarChart } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import { PET_STATUS_COLORS, PET_OWNERSHIP_COLORS } from "@/components/charts/colors";
import { DateRangePicker } from "@/components/date-picker";
import {
  CaretDownIcon,
  CaretUpIcon,
  FunnelIcon,
  PrinterIcon,
  DownloadSimpleIcon,
} from "@/components/icons";
import {
  FilterChip,
  FilterField,
  ReportSkeleton,
  getPresetDates,
  PRESET_KEYS,
  useTranslatedRows,
} from "@/components/report";
import type { DatePreset } from "@/components/report";
import { PetsKpiCards } from "@/components/reports/pets-kpi-cards";
import { UISelect } from "@/components/ui-select";
import { petStatusKeys } from "@/constants/pet";
import { usePetReport } from "@/hooks/use-pet-reports";
import { exportGroupCSV } from "@/lib/export-csv";
import type { CountResult, MonthlyStatusCount, PetReportParams } from "@/types/report";

export const Route = createLazyFileRoute("/_app/projects/$projectId/reports/pets")({
  component: PetReportsPage,
});

const PET_STATUS_OPTIONS = [
  "registered",
  "adopted",
  "owner_found",
  "needs_shelter",
  "unknown",
] as const;

function extractMonthlySeriesForStatus(data: MonthlyStatusCount[], status: string): CountResult[] {
  return data.filter((r) => r.status === status).map((r) => ({ label: r.month, count: r.count }));
}

function extractMonthlyTotals(data: MonthlyStatusCount[]): CountResult[] {
  const totals = new Map<string, number>();
  for (const r of data) {
    totals.set(r.month, (totals.get(r.month) ?? 0) + r.count);
  }
  return Array.from(totals, ([label, count]) => ({ label, count })).sort((a, b) =>
    a.label.localeCompare(b.label),
  );
}

function ReportCard({
  title,
  rows,
  chart,
  colorMap,
  direction,
  yAxisLabel,
  onExport,
  total,
}: {
  title: string;
  rows: CountResult[];
  chart: "bar" | "pie";
  colorMap?: Record<string, string>;
  direction?: "vertical" | "horizontal" | "auto";
  yAxisLabel?: string;
  onExport?: () => void;
  total?: number;
}) {
  const { t } = useTranslation();
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-3 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <div className="flex items-center gap-2">
          {onExport && (
            <button
              type="button"
              onClick={onExport}
              className="text-fg-tertiary transition-colors hover:text-fg"
              title={t("common.downloadCsv")}
            >
              <DownloadSimpleIcon size={14} />
            </button>
          )}
          {total != null && (
            <span className="tabular-nums text-xs font-medium text-fg-tertiary">
              {total.toLocaleString()}
            </span>
          )}
        </div>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart data={rows} colorMap={colorMap} direction={direction} yAxisLabel={yAxisLabel} />
        ) : (
          <PieChart data={rows} colorMap={colorMap} />
        )
      ) : (
        <p className="py-8 text-center text-sm text-fg-tertiary">&mdash;</p>
      )}
    </div>
  );
}

function PetReportsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const [params, setParams] = useState<PetReportParams>({});
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [activePreset, setActivePreset] = useState<DatePreset | null>(null);
  const { data, isLoading } = usePetReport(projectId, params);

  const statusOptions = PET_STATUS_OPTIONS.map((s) => ({
    label: t(petStatusKeys[s] ?? s),
    value: s,
  }));

  const hasFilters = Object.values(params).some((v) => v != null && v !== "");
  const axisLabel = t("project.reports.axisCount");
  const clearDatePreset = () => setActivePreset(null);

  const needsShelterCount =
    data?.by_status.rows.find((r) => r.label === "needs_shelter")?.count ?? 0;
  const adoptedCount = data?.by_status.rows.find((r) => r.label === "adopted")?.count ?? 0;

  const translatedStatus = useTranslatedRows(data?.by_status.rows ?? []);
  const translatedOwnership = useTranslatedRows(data?.by_ownership.rows ?? []);

  const needsShelterMonthly = data
    ? extractMonthlySeriesForStatus(data.by_status_by_month, "needs_shelter")
    : [];
  const adoptedMonthly = data
    ? extractMonthlySeriesForStatus(data.by_status_by_month, "adopted")
    : [];
  const totalMonthly = data ? extractMonthlyTotals(data.by_status_by_month) : [];

  return (
    <div>
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
            {data && (
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

      {isLoading && <ReportSkeleton kpiCount={3} />}

      {data && (
        <div className="grid gap-6 lg:grid-cols-2">
          <PetsKpiCards
            total={data.by_status.total}
            needsShelter={needsShelterCount}
            adopted={adoptedCount}
          />

          {/* Status & Ownership pies */}
          <ReportCard
            title={t("project.petReports.byStatus")}
            rows={translatedStatus}
            chart="pie"
            colorMap={PET_STATUS_COLORS}
            total={data.by_status.total}
            onExport={() => exportGroupCSV(t("project.petReports.byStatus"), translatedStatus)}
          />
          <ReportCard
            title={t("project.petReports.byOwnership")}
            rows={translatedOwnership}
            chart="pie"
            colorMap={PET_OWNERSHIP_COLORS}
            total={data.by_ownership.total}
            onExport={() =>
              exportGroupCSV(t("project.petReports.byOwnership"), translatedOwnership)
            }
          />

          {/* Monthly trends */}
          <div className="col-span-full">
            <ReportCard
              title={t("project.petReports.byMonth")}
              rows={data.by_month.rows}
              chart="bar"
              yAxisLabel={axisLabel}
              total={data.by_month.total}
              onExport={() => exportGroupCSV(t("project.petReports.byMonth"), data.by_month.rows)}
            />
          </div>

          <div className="col-span-full">
            <ReportCard
              title={t("project.petReports.totalByMonth")}
              rows={totalMonthly}
              chart="bar"
              yAxisLabel={axisLabel}
              onExport={() => exportGroupCSV(t("project.petReports.totalByMonth"), totalMonthly)}
            />
          </div>

          <ReportCard
            title={t("project.petReports.needsShelterByMonth")}
            rows={needsShelterMonthly}
            chart="bar"
            yAxisLabel={axisLabel}
            colorMap={{
              ...Object.fromEntries(needsShelterMonthly.map((r) => [r.label, "#ef4444"])),
            }}
            onExport={() =>
              exportGroupCSV(t("project.petReports.needsShelterByMonth"), needsShelterMonthly)
            }
          />
          <ReportCard
            title={t("project.petReports.adoptedByMonth")}
            rows={adoptedMonthly}
            chart="bar"
            yAxisLabel={axisLabel}
            colorMap={{ ...Object.fromEntries(adoptedMonthly.map((r) => [r.label, "#10b981"])) }}
            onExport={() => exportGroupCSV(t("project.petReports.adoptedByMonth"), adoptedMonthly)}
          />
        </div>
      )}
    </div>
  );
}
