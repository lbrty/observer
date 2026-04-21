import { useState } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { type BarLegendItem } from "@/components/charts/bar-chart";
import { ReportSkeleton, labelKeyMap, AGE_RANGE_MAP } from "@/components/reports/shared";
import type { DatePreset } from "@/components/reports/shared";
import { ReportFilterBar } from "@/components/reports/report-filter-bar";
import { PeopleReportSections } from "@/components/reports/people-report-sections";
import {
  CaretDownIcon,
  CaretUpIcon,
  DownloadSimpleIcon,
  FunnelIcon,
  PrinterIcon,
} from "@/components/ui/icons";
import { useReport } from "@/hooks/reports/use-reports";
import { exportReportCSV } from "@/lib/export-csv";
import type { ReportParams } from "@/types/report";

export const Route = createLazyFileRoute("/_app/projects/$projectId/reports/people")({
  component: ReportsPage,
});

function ReportsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const [params, setParams] = useState<ReportParams>({});
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [activePreset, setActivePreset] = useState<DatePreset | null>(null);
  const { data, isLoading } = useReport(projectId, params);

  const ageGroupLegend: BarLegendItem[] = Object.entries(AGE_RANGE_MAP).map(([key, range]) => ({
    short: range,
    full: t(labelKeyMap[key] ?? key),
  }));

  const axisLabel = t("project.reports.axisCount");

  return (
    <div>
      {/* Print-only header */}
      <div className="print-header hidden">
        <h1 className="text-lg font-bold">{t("project.reports.title")}</h1>
        {params.date_from && (
          <p>
            {params.date_from} &mdash; {params.date_to ?? "..."}
          </p>
        )}
      </div>

      {/* Unified header + filter panel */}
      <div
        data-print-hide
        className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary"
      >
        {/* Top bar */}
        <div className="flex items-center justify-between px-5 py-3">
          <h1 className="font-serif text-xl font-bold tracking-tight text-fg">
            {t("project.reports.title")}
          </h1>
          <div className="flex items-center gap-2">
            {data && (
              <>
                <button
                  type="button"
                  onClick={() => exportReportCSV(data, projectId)}
                  className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
                >
                  <DownloadSimpleIcon size={14} />
                  {t("project.reports.exportCsv")}
                </button>
                <button
                  type="button"
                  onClick={() => window.print()}
                  className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
                >
                  <PrinterIcon size={14} />
                  {t("project.reports.print")}
                </button>
              </>
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

        {/* Collapsible filter panel */}
        {filtersOpen && (
          <div className="border-t border-border-secondary px-5 pb-4 pt-3">
            <ReportFilterBar
              params={params}
              activePreset={activePreset}
              onParamsChange={setParams}
              onPresetChange={setActivePreset}
            />
          </div>
        )}
      </div>

      {/* Loading skeleton */}
      {isLoading && <ReportSkeleton />}

      {/* Dashboard content */}
      {data && (
        <PeopleReportSections data={data} axisLabel={axisLabel} ageGroupLegend={ageGroupLegend} />
      )}
    </div>
  );
}
