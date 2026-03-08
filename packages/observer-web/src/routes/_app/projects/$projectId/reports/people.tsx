import { useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
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
import { DatePicker } from "@/components/date-picker";
import {
  CaretDownIcon,
  CaretUpIcon,
  DownloadSimpleIcon,
  FunnelIcon,
  PrinterIcon,
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
import { SEX_VALUES, AGE_GROUP_VALUES, CASE_STATUS_VALUES } from "@/constants/person";
import { useCategories } from "@/hooks/use-categories";
import { useOffices } from "@/hooks/use-offices";
import { useReport } from "@/hooks/use-reports";
import { exportReportCSV } from "@/lib/export-csv";
import type { ReportParams } from "@/types/report";

export const Route = createFileRoute("/_app/projects/$projectId/reports/people")({
  component: ReportsPage,
});

const SUPPORT_TYPE_OPTIONS = ["humanitarian", "legal", "social", "psychological", "medical", "general"] as const;

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
  const { data: offices } = useOffices();
  const { data: categories } = useCategories();

  const officeOptions = (offices ?? []).map((o) => ({ label: o.name, value: o.id }));
  const categoryOptions = (categories ?? []).map((c) => ({ label: c.name, value: c.id }));
  const supportTypeOptions = SUPPORT_TYPE_OPTIONS.map((s) => ({ label: t(labelKeyMap[s] ?? s), value: s }));
  const caseStatusOptions = CASE_STATUS_VALUES.map((s) => ({ label: t(`project.people.${s}`), value: s }));
  const sexOptions = SEX_VALUES.map((s) => ({ label: t(`project.people.sex${s[0].toUpperCase()}${s.slice(1)}`), value: s }));
  const ageGroupOptions = AGE_GROUP_VALUES.map((g) => ({ label: t(labelKeyMap[g] ?? g), value: g }));

  const ageGroupLegend: BarLegendItem[] = Object.entries(AGE_RANGE_MAP).map(([key, range]) => ({
    short: range,
    full: t(labelKeyMap[key] ?? key),
  }));

  const hasFilters = Object.values(params).some((v) => v != null && v !== "");
  const axisLabel = t("project.reports.axisCount");
  const clearDatePreset = () => setActivePreset(null);

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
      <div data-print-hide className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary">
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
            {/* Date presets */}
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

            {/* Filter grid */}
            <div className="grid grid-cols-2 gap-x-4 gap-y-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-7">
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
              <FilterField label={t("project.reports.filterOffice")}>
                <UISelect
                  value={params.office_id ?? ""}
                  onValueChange={(v) => setParams((p) => ({ ...p, office_id: v || undefined }))}
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...officeOptions]}
                  placeholder={t("project.reports.allValues")}
                  fullWidth
                />
              </FilterField>
              <FilterField label={t("project.reports.filterCategory")}>
                <UISelect
                  value={params.category_id ?? ""}
                  onValueChange={(v) => setParams((p) => ({ ...p, category_id: v || undefined }))}
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...categoryOptions]}
                  placeholder={t("project.reports.allValues")}
                  fullWidth
                />
              </FilterField>
              <FilterField label={t("project.reports.filterCaseStatus")}>
                <UISelect
                  value={params.case_status ?? ""}
                  onValueChange={(v) => setParams((p) => ({ ...p, case_status: v || undefined }))}
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...caseStatusOptions]}
                  placeholder={t("project.reports.allValues")}
                  fullWidth
                />
              </FilterField>
              <FilterField label={t("project.reports.filterSex")}>
                <UISelect
                  value={params.sex ?? ""}
                  onValueChange={(v) => setParams((p) => ({ ...p, sex: v || undefined }))}
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...sexOptions]}
                  placeholder={t("project.reports.allValues")}
                  fullWidth
                />
              </FilterField>
              <FilterField label={t("project.reports.filterAgeGroup")}>
                <UISelect
                  value={params.age_group ?? ""}
                  onValueChange={(v) => setParams((p) => ({ ...p, age_group: v || undefined }))}
                  options={[{ label: t("project.reports.allValues"), value: "" }, ...ageGroupOptions]}
                  placeholder={t("project.reports.allValues")}
                  fullWidth
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

        {/* Active filter chips (always visible) */}
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
            {params.office_id && (
              <FilterChip
                label={t("project.reports.filterOffice")}
                value={officeOptions.find((o) => o.value === params.office_id)?.label ?? params.office_id}
                onRemove={() => setParams((p) => ({ ...p, office_id: undefined }))}
              />
            )}
            {params.category_id && (
              <FilterChip
                label={t("project.reports.filterCategory")}
                value={categoryOptions.find((c) => c.value === params.category_id)?.label ?? params.category_id}
                onRemove={() => setParams((p) => ({ ...p, category_id: undefined }))}
              />
            )}
            {params.case_status && (
              <FilterChip
                label={t("project.reports.filterCaseStatus")}
                value={caseStatusOptions.find((s) => s.value === params.case_status)?.label ?? params.case_status}
                onRemove={() => setParams((p) => ({ ...p, case_status: undefined }))}
              />
            )}
            {params.sex && (
              <FilterChip
                label={t("project.reports.filterSex")}
                value={sexOptions.find((s) => s.value === params.sex)?.label ?? params.sex}
                onRemove={() => setParams((p) => ({ ...p, sex: undefined }))}
              />
            )}
            {params.age_group && (
              <FilterChip
                label={t("project.reports.filterAgeGroup")}
                value={ageGroupOptions.find((g) => g.value === params.age_group)?.label ?? params.age_group}
                onRemove={() => setParams((p) => ({ ...p, age_group: undefined }))}
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

      {/* Loading skeleton */}
      {isLoading && <ReportSkeleton />}

      {/* Dashboard content */}
      {data && (
        <div className="report-grid grid gap-6 lg:grid-cols-2">
          {/* Overview: KPI cards + Sankey side by side */}
          <div className="col-span-full grid grid-cols-1 gap-6 lg:grid-cols-2">
            {/* KPI cards — 3x2 grid on the left */}
            <div className="grid grid-cols-3 gap-3">
              <KpiCard label={t("project.reports.kpiPeople")} value={data.by_sex.total} />
              <KpiCard label={t("project.reports.kpiConsultations")} value={data.consultations.total} />
              <KpiCard
                label={t("project.reports.kpiActiveCases")}
                value={data.by_case_status?.rows.find((r) => r.label === "active")?.count ?? 0}
              />
              <KpiCard
                label={t("project.reports.kpiIdp")}
                value={data.by_idp_status.rows.find((r) => r.label === "idp")?.count ?? data.by_idp_status.total}
              />
              <KpiCard label={t("project.reports.kpiHouseholds")} value={data.family_units.total} />
              <KpiCard label={t("project.reports.kpiOffices")} value={data.by_office.rows.length} />
            </div>

            {/* Sankey on the right */}
            {data.status_flow && data.status_flow.length > 0 && (
              <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
                <h3 className="mb-3 text-sm font-semibold text-fg">
                  {t("project.reports.statusFlow")}
                </h3>
                <SankeyChart
                  data={data.status_flow}
                  translateLabel={(l) => {
                    const key = labelKeyMap[l];
                    return key ? t(key) : t(`project.people.${l}`, l);
                  }}
                />
              </div>
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
          <ReportCard group={data.by_sphere} title={t("project.reports.bySphere")} chart="bar" yAxisLabel={axisLabel} colorMap={SPHERE_COLORS} direction="auto" />
          <ReportCard group={data.people_by_sphere} title={t("project.reports.peopleBySphere")} chart="bar" yAxisLabel={axisLabel} colorMap={SPHERE_COLORS} direction="auto" />
          <ReportCard group={data.by_office} title={t("project.reports.byOffice")} chart="bar" yAxisLabel={axisLabel} direction="auto" />

          {/* Demographics */}
          <SectionHeader title={t("project.reports.sectionDemographics")} />
          <div className="col-span-full grid grid-cols-1 gap-6 md:grid-cols-3">
            <ReportCard group={data.by_sex} title={t("project.reports.bySex")} chart="pie" colorMap={SEX_COLORS} />
            <ReportCard group={data.family_units} title={t("project.reports.familyUnits")} chart="pie" />
            <ReportCard group={data.by_idp_status} title={t("project.reports.byIdpStatus")} chart="pie" colorMap={IDP_STATUS_COLORS} />
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
          <ReportCard group={data.by_region} title={t("project.reports.byRegion")} chart="bar" yAxisLabel={axisLabel} direction="auto" />
          <ReportCard group={data.by_category} title={t("project.reports.byCategory")} chart="bar" yAxisLabel={axisLabel} direction="auto" />
          <div className="col-span-full">
            <ReportCard group={data.by_tag} title={t("project.reports.byTag")} chart="bar" yAxisLabel={axisLabel} direction="auto" />
          </div>
        </div>
      )}
    </div>
  );
}
