import { useState } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { type BarLegendItem } from "@/components/charts/bar-chart";
import { SankeyChart } from "@/components/charts/sankey-chart";
import {
  SEX_COLORS,
  SUPPORT_TYPE_COLORS,
  SPHERE_COLORS,
  IDP_STATUS_COLORS,
  AGE_GROUP_COLORS,
} from "@/components/charts/colors";
import {
  ReportCard,
  ReportSkeleton,
  labelKeyMap,
  AGE_RANGE_MAP,
} from "@/components/report";
import type { DatePreset } from "@/components/report";
import { ReportFilterBar } from "@/components/reports/report-filter-bar";
import { PeopleKpiCards } from "@/components/reports/people-kpi-cards";
import { PeopleChartSection } from "@/components/reports/people-chart-section";
import {
  CaretDownIcon,
  CaretUpIcon,
  DownloadSimpleIcon,
  FunnelIcon,
  PrinterIcon,
} from "@/components/icons";
import { useReport } from "@/hooks/use-reports";
import { exportReportCSV } from "@/lib/export-csv";
import type { ReportParams } from "@/types/report";

export const Route = createLazyFileRoute("/_app/projects/$projectId/reports/people")({
  component: ReportsPage,
});

function SectionHeader({ title }: { title: string }) {
  return (
    <div className="col-span-full border-b border-border-secondary pb-1 pt-4">
      <h2 className="text-xs font-semibold uppercase tracking-wider text-fg-tertiary">{title}</h2>
    </div>
  );
}

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
        <div className="report-grid grid gap-6 lg:grid-cols-2">
          {/* Overview: KPI cards + Sankey side by side */}
          <div className="col-span-full grid grid-cols-1 gap-6 lg:grid-cols-2">
            <PeopleKpiCards
              totalPeople={data.by_sex.total}
              totalConsultations={data.consultations.total}
              activeCases={data.by_case_status?.rows.find((r) => r.label === "active")?.count ?? 0}
              idpCount={
                data.by_idp_status.rows.find((r) => r.label === "idp")?.count ??
                data.by_idp_status.total
              }
              households={data.family_units.total}
              offices={data.by_office.rows.length}
            />

            {data.status_flow && data.status_flow.length > 0 && (
              <PeopleChartSection title={t("project.reports.statusFlow")}>
                <SankeyChart
                  data={data.status_flow}
                  translateLabel={(l) => {
                    const key = labelKeyMap[l];
                    return key ? t(key) : t(`project.people.${l}`, l);
                  }}
                />
              </PeopleChartSection>
            )}
          </div>

          {/* Services */}
          <SectionHeader title={t("project.reports.sectionServices")} />
          <div className="col-span-full">
            <ReportCard
              group={data.consultations}
              title={t("project.reports.consultations")}
              chart="bar"
              yAxisLabel={axisLabel}
              colorMap={SUPPORT_TYPE_COLORS}
            />
          </div>
          <ReportCard
            group={data.by_sphere}
            title={t("project.reports.bySphere")}
            chart="bar"
            yAxisLabel={axisLabel}
            colorMap={SPHERE_COLORS}
            direction="auto"
          />
          <ReportCard
            group={data.people_by_sphere}
            title={t("project.reports.peopleBySphere")}
            chart="bar"
            yAxisLabel={axisLabel}
            colorMap={SPHERE_COLORS}
            direction="auto"
          />
          <ReportCard
            group={data.by_office}
            title={t("project.reports.byOffice")}
            chart="bar"
            yAxisLabel={axisLabel}
            direction="auto"
          />

          {/* Demographics */}
          <SectionHeader title={t("project.reports.sectionDemographics")} />
          <div className="col-span-full grid grid-cols-1 gap-6 md:grid-cols-3">
            <ReportCard
              group={data.by_sex}
              title={t("project.reports.bySex")}
              chart="pie"
              colorMap={SEX_COLORS}
            />
            <ReportCard
              group={data.family_units}
              title={t("project.reports.familyUnits")}
              chart="pie"
            />
            <ReportCard
              group={data.by_idp_status}
              title={t("project.reports.byIdpStatus")}
              chart="pie"
              colorMap={IDP_STATUS_COLORS}
            />
          </div>
          <div className="col-span-full">
            <ReportCard
              group={data.by_age_group}
              title={t("project.reports.byAgeGroup")}
              chart="bar"
              yAxisLabel={axisLabel}
              skipTranslation
              mapLabel={(l) => AGE_RANGE_MAP[l] ?? l}
              legend={ageGroupLegend}
              colorMap={AGE_GROUP_COLORS}
            />
          </div>
          <div className="col-span-full">
            <ReportCard
              group={data.consultations_by_age_group}
              title={t("project.reports.consultationsByAgeGroup")}
              chart="bar"
              yAxisLabel={axisLabel}
              skipTranslation
              mapLabel={(l) => AGE_RANGE_MAP[l] ?? l}
              legend={ageGroupLegend}
              colorMap={AGE_GROUP_COLORS}
            />
          </div>

          {/* Geography & Taxonomy */}
          <SectionHeader title={t("project.reports.sectionGeography")} />
          <ReportCard
            group={data.by_region}
            title={t("project.reports.byRegion")}
            chart="bar"
            yAxisLabel={axisLabel}
            direction="auto"
          />
          <ReportCard
            group={data.by_category}
            title={t("project.reports.byCategory")}
            chart="bar"
            yAxisLabel={axisLabel}
            direction="auto"
          />
          <div className="col-span-full">
            <ReportCard
              group={data.by_tag}
              title={t("project.reports.byTag")}
              chart="bar"
              yAxisLabel={axisLabel}
              direction="auto"
            />
          </div>
        </div>
      )}
    </div>
  );
}
