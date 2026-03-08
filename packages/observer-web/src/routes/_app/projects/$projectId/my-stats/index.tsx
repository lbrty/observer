import { useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { type BarLegendItem } from "@/components/charts/bar-chart";
import {
  SEX_COLORS,
  SUPPORT_TYPE_COLORS,
  SPHERE_COLORS,
  AGE_GROUP_COLORS,
} from "@/components/charts/colors";
import { DatePicker } from "@/components/date-picker";
import {
  CaretDownIcon,
  CaretUpIcon,
  DownloadSimpleIcon,
  FunnelIcon,
} from "@/components/icons";
import {
  ReportCard,
  KpiCard,
  FilterChip,
  FilterField,
  ReportSkeleton,
  labelKeyMap,
  AGE_RANGE_MAP,
  getPresetDates,
  PRESET_KEYS,
} from "@/components/report";
import type { DatePreset } from "@/components/report";
import { UISelect } from "@/components/ui-select";
import { useReport } from "@/hooks/use-reports";
import { exportReportCSV } from "@/lib/export-csv";
import { useAuth } from "@/stores/auth";
import type { ReportParams } from "@/types/report";

export const Route = createFileRoute("/_app/projects/$projectId/my-stats/")({
  component: MyStatsPage,
});

const SUPPORT_TYPE_OPTIONS = ["humanitarian", "legal", "social", "psychological", "medical", "general"] as const;

function MyStatsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const { user } = useAuth();
  const [params, setParams] = useState<ReportParams>({});
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [activePreset, setActivePreset] = useState<DatePreset | null>(null);

  const reportParams: ReportParams = { ...params, consultant_id: user?.id };
  const { data, isLoading } = useReport(projectId, reportParams);

  const supportTypeOptions = SUPPORT_TYPE_OPTIONS.map((s) => ({
    label: t(labelKeyMap[s] ?? s),
    value: s,
  }));

  const hasFilters = Object.entries(params).some(([, v]) => v != null && v !== "");
  const axisLabel = t("project.reports.axisCount");
  const clearDatePreset = () => setActivePreset(null);

  const ageGroupLegend: BarLegendItem[] = Object.entries(AGE_RANGE_MAP).map(([key, range]) => ({
    short: range,
    full: t(labelKeyMap[key] ?? key),
  }));

  return (
    <div>
      {/* Header + filters */}
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
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...supportTypeOptions]}
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
                onRemove={() => { setParams((p) => ({ ...p, date_from: undefined })); clearDatePreset(); }}
              />
            )}
            {params.date_to && (
              <FilterChip
                label={t("project.reports.dateTo")}
                value={params.date_to}
                onRemove={() => { setParams((p) => ({ ...p, date_to: undefined })); clearDatePreset(); }}
              />
            )}
            {params.support_type && (
              <FilterChip
                label={t("project.reports.filterSupportType")}
                value={supportTypeOptions.find((s) => s.value === params.support_type)?.label ?? params.support_type}
                onRemove={() => setParams((p) => ({ ...p, support_type: undefined }))}
              />
            )}
            <button
              type="button"
              onClick={() => { setParams({}); clearDatePreset(); }}
              className="ml-1 text-xs font-medium text-fg-tertiary underline transition-colors hover:text-fg"
            >
              {t("project.reports.clearAll")}
            </button>
          </div>
        )}
      </div>

      {isLoading && <ReportSkeleton kpiCount={4} />}

      {data && (
        <div className="grid gap-6 lg:grid-cols-2">
          {/* KPI overview */}
          <div className="col-span-full grid grid-cols-2 gap-3 sm:grid-cols-4">
            <KpiCard label={t("project.myStats.kpiPeople")} value={data.by_sex.total} />
            <KpiCard label={t("project.myStats.kpiConsultations")} value={data.consultations.total} />
            <KpiCard
              label={t("project.myStats.kpiActiveCases")}
              value={data.by_case_status?.rows.find((r) => r.label === "active")?.count ?? 0}
            />
            <KpiCard label={t("project.myStats.kpiHouseholds")} value={data.family_units.total} />
          </div>

          {/* Consultations */}
          <div className="col-span-full">
            <ReportCard group={data.consultations} title={t("project.reports.consultations")} chart="bar" yAxisLabel={axisLabel} colorMap={SUPPORT_TYPE_COLORS} />
          </div>

          {/* Service breakdown */}
          <ReportCard group={data.by_sphere} title={t("project.reports.bySphere")} chart="bar" yAxisLabel={axisLabel} colorMap={SPHERE_COLORS} direction="auto" />
          <ReportCard group={data.by_office} title={t("project.reports.byOffice")} chart="bar" yAxisLabel={axisLabel} direction="auto" />
          <ReportCard group={data.by_region} title={t("project.reports.byRegion")} chart="bar" yAxisLabel={axisLabel} direction="auto" />

          {/* Demographics */}
          <ReportCard group={data.by_sex} title={t("project.reports.bySex")} chart="pie" colorMap={SEX_COLORS} />
          <ReportCard group={data.by_case_status} title={t("project.reports.byCaseStatus")} chart="bar" yAxisLabel={axisLabel} direction="auto" />

          {/* Age distribution */}
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

          {/* Categories & tags */}
          <ReportCard group={data.by_category} title={t("project.reports.byCategory")} chart="bar" yAxisLabel={axisLabel} direction="auto" />
          <ReportCard group={data.by_tag} title={t("project.reports.byTag")} chart="bar" yAxisLabel={axisLabel} direction="auto" />
        </div>
      )}
    </div>
  );
}
