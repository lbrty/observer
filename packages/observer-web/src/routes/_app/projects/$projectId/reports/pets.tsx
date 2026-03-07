import { useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { BarChart } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import {
  PET_STATUS_COLORS,
  PET_OWNERSHIP_COLORS,
} from "@/components/charts/colors";
import { UISelect } from "@/components/ui-select";
import {
  CaretDownIcon,
  CaretUpIcon,
  DownloadSimpleIcon,
  FunnelIcon,
  PrinterIcon,
  XIcon,
} from "@/components/icons";
import { usePetReport } from "@/hooks/use-pet-reports";
import { exportGroupCSV } from "@/lib/export-csv";
import type { CountResult, MonthlyStatusCount, PetReportParams, ReportGroup } from "@/types/report";

export const Route = createFileRoute("/_app/projects/$projectId/reports/pets")({
  component: PetReportsPage,
});

const PET_STATUS_OPTIONS = [
  "registered",
  "adopted",
  "owner_found",
  "needs_shelter",
  "unknown",
] as const;

const statusLabelKeyMap: Record<string, string> = {
  registered: "project.pets.statusRegistered",
  adopted: "project.pets.statusAdopted",
  owner_found: "project.pets.statusOwnerFound",
  needs_shelter: "project.pets.statusNeedsShelter",
  unknown: "project.pets.statusUnknown",
};

const ownershipLabelKeyMap: Record<string, string> = {
  with_owner: "project.petReports.withOwner",
  without_owner: "project.petReports.withoutOwner",
};

function useTranslatedRows(rows: CountResult[], keyMap: Record<string, string>): CountResult[] {
  const { t } = useTranslation();
  return rows.map((r) => {
    const key = keyMap[r.label];
    return key ? { ...r, label: t(key) } : r;
  });
}

function extractMonthlySeriesForStatus(
  data: MonthlyStatusCount[],
  status: string,
): CountResult[] {
  return data
    .filter((r) => r.status === status)
    .map((r) => ({ label: r.month, count: r.count }));
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

type DatePreset = "month" | "quarter" | "year" | "all";

function getPresetDates(preset: DatePreset): { date_from?: string; date_to?: string } {
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

const PRESET_KEYS: { key: DatePreset; i18n: string }[] = [
  { key: "month", i18n: "project.reports.presetMonth" },
  { key: "quarter", i18n: "project.reports.presetQuarter" },
  { key: "year", i18n: "project.reports.presetYear" },
  { key: "all", i18n: "project.reports.presetAll" },
];

function KpiCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-4">
      <p className="text-2xl font-bold tabular-nums text-fg">{value.toLocaleString()}</p>
      <p className="mt-0.5 text-xs font-medium text-fg-tertiary">{label}</p>
    </div>
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
              title="Download CSV"
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

function FilterField({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="space-y-1.5">
      <span className="block text-xs font-medium text-fg-secondary">{label}</span>
      {children}
    </div>
  );
}

function FilterChip({
  label,
  value,
  onRemove,
}: {
  label: string;
  value: string;
  onRemove: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onRemove}
      className="inline-flex items-center gap-1 rounded-md bg-bg-tertiary px-2 py-0.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
    >
      <span className="text-fg-tertiary">{label}:</span> {value}
      <XIcon size={10} />
    </button>
  );
}

function ReportSkeleton() {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-3 gap-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
      <div className="grid gap-6 lg:grid-cols-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-72 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
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
    label: t(statusLabelKeyMap[s] ?? s),
    value: s,
  }));

  const hasFilters = Object.values(params).some((v) => v != null && v !== "");
  const axisLabel = t("project.reports.axisCount");
  const clearDatePreset = () => setActivePreset(null);

  const needsShelterCount =
    data?.by_status.rows.find((r) => r.label === "needs_shelter")?.count ?? 0;
  const adoptedCount =
    data?.by_status.rows.find((r) => r.label === "adopted")?.count ?? 0;

  const translatedStatus = useTranslatedRows(data?.by_status.rows ?? [], statusLabelKeyMap);
  const translatedOwnership = useTranslatedRows(data?.by_ownership.rows ?? [], ownershipLabelKeyMap);

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
      <div data-print-hide className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary">
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

            <div className="grid grid-cols-2 gap-x-4 gap-y-3 sm:grid-cols-3">
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
              <FilterField label={t("project.petReports.filterStatus")}>
                <UISelect
                  value={params.status ?? ""}
                  onValueChange={(v) => setParams((p) => ({ ...p, status: v || undefined }))}
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...statusOptions]}
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

      {isLoading && <ReportSkeleton />}

      {data && (
        <div className="grid gap-6 lg:grid-cols-2">
          {/* KPI Cards */}
          <div className="col-span-full grid grid-cols-3 gap-3">
            <KpiCard label={t("project.petReports.kpiTotal")} value={data.by_status.total} />
            <KpiCard label={t("project.petReports.kpiNeedsShelter")} value={needsShelterCount} />
            <KpiCard label={t("project.petReports.kpiAdopted")} value={adoptedCount} />
          </div>

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
            onExport={() => exportGroupCSV(t("project.petReports.byOwnership"), translatedOwnership)}
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
            colorMap={{ ...Object.fromEntries(needsShelterMonthly.map((r) => [r.label, "#ef4444"])) }}
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
            onExport={() =>
              exportGroupCSV(t("project.petReports.adoptedByMonth"), adoptedMonthly)
            }
          />
        </div>
      )}
    </div>
  );
}
