import { useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { BarChart, type BarLegendItem } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import { SankeyChart } from "@/components/charts/sankey-chart";
import {
  SEX_COLORS,
  SUPPORT_TYPE_COLORS,
  SPHERE_COLORS,
  IDP_STATUS_COLORS,
  AGE_GROUP_COLORS,
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
import { useCategories } from "@/hooks/use-categories";
import { useOffices } from "@/hooks/use-offices";
import { useReport } from "@/hooks/use-reports";
import { exportGroupCSV, exportReportCSV } from "@/lib/export-csv";
import type { CountResult, ReportGroup, ReportParams } from "@/types/report";

export const Route = createFileRoute("/_app/projects/$projectId/reports/")({
  component: ReportsPage,
});

const labelKeyMap: Record<string, string> = {
  // support types
  humanitarian: "project.supportRecords.typeHumanitarian",
  legal: "project.supportRecords.typeLegal",
  social: "project.supportRecords.typeSocial",
  psychological: "project.supportRecords.typePsychological",
  medical: "project.supportRecords.typeMedical",
  general: "project.supportRecords.typeGeneral",
  // spheres
  housing_assistance: "project.supportRecords.sphereHousing",
  document_recovery: "project.supportRecords.sphereDocumentRecovery",
  social_benefits: "project.supportRecords.sphereSocialBenefits",
  property_rights: "project.supportRecords.spherePropertyRights",
  employment_rights: "project.supportRecords.sphereEmploymentRights",
  family_law: "project.supportRecords.sphereFamilyLaw",
  healthcare_access: "project.supportRecords.sphereHealthcareAccess",
  education_access: "project.supportRecords.sphereEducationAccess",
  financial_aid: "project.supportRecords.sphereFinancialAid",
  psychological_support: "project.supportRecords.spherePsychologicalSupport",
  other: "project.supportRecords.sphereOther",
  unspecified: "project.supportRecords.sphereOther",
  // sex
  male: "project.people.sexMale",
  female: "project.people.sexFemale",
  unknown: "project.people.sexUnknown",
  // age groups
  infant: "project.people.ageInfant",
  toddler: "project.people.ageToddler",
  pre_school: "project.people.agePreSchool",
  middle_childhood: "project.people.ageMiddleChildhood",
  young_teen: "project.people.ageYoungTeen",
  teenager: "project.people.ageTeenager",
  young_adult: "project.people.ageYoungAdult",
  early_adult: "project.people.ageEarlyAdult",
  middle_aged_adult: "project.people.ageMiddleAgedAdult",
  old_adult: "project.people.ageOldAdult",
  // referral
  pending: "project.supportRecords.referralPending",
  accepted: "project.supportRecords.referralAccepted",
  completed: "project.supportRecords.referralCompleted",
  declined: "project.supportRecords.referralDeclined",
  no_response: "project.supportRecords.referralNoResponse",
};

function useTranslatedRows(rows: CountResult[]): CountResult[] {
  const { t } = useTranslation();
  const translated = rows.map((r) => {
    const key = labelKeyMap[r.label];
    return key ? { ...r, label: t(key) } : r;
  });
  const merged = new Map<string, number>();
  for (const r of translated) {
    merged.set(r.label, (merged.get(r.label) ?? 0) + r.count);
  }
  return Array.from(merged, ([label, count]) => ({ label, count }));
}

function ReportCard({
  group,
  title,
  chart,
  yAxisLabel,
  legend,
  mapLabel,
  skipTranslation,
  colorMap,
  direction,
}: {
  group: ReportGroup;
  title: string;
  chart: "bar" | "pie";
  yAxisLabel?: string;
  legend?: BarLegendItem[];
  mapLabel?: (label: string) => string;
  skipTranslation?: boolean;
  colorMap?: Record<string, string>;
  direction?: "vertical" | "horizontal" | "auto";
}) {
  const translated = useTranslatedRows(group.rows);
  const source = skipTranslation ? group.rows : translated;
  const rows = mapLabel ? source.map((r) => ({ ...r, label: mapLabel(r.label) })) : source;
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-3 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => exportGroupCSV(title, rows)}
            className="text-fg-tertiary transition-colors hover:text-fg"
            title="Download CSV"
          >
            <DownloadSimpleIcon size={14} />
          </button>
          <span className="tabular-nums text-xs font-medium text-fg-tertiary">
            {group.total.toLocaleString()}
          </span>
        </div>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart
            data={rows}
            yAxisLabel={yAxisLabel}
            legend={legend}
            colorMap={colorMap}
            direction={direction}
          />
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

function SectionHeader({ title }: { title: string }) {
  return (
    <div className="col-span-full border-b border-border-secondary pb-1 pt-4">
      <h2 className="text-xs font-semibold uppercase tracking-wider text-fg-tertiary">{title}</h2>
    </div>
  );
}

function KpiCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-4">
      <p className="text-2xl font-bold tabular-nums text-fg">{value.toLocaleString()}</p>
      <p className="mt-0.5 text-xs font-medium text-fg-tertiary">{label}</p>
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
      <div className="grid grid-cols-3 gap-4 lg:grid-cols-6">
        {Array.from({ length: 6 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
      <div className="h-64 animate-pulse rounded-xl bg-bg-tertiary" />
      <div className="grid gap-6 lg:grid-cols-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-72 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
    </div>
  );
}

const AGE_RANGE_MAP: Record<string, string> = {
  infant: "0-1",
  toddler: "1-3",
  pre_school: "3-6",
  middle_childhood: "6-12",
  young_teen: "12-14",
  teenager: "14-18",
  young_adult: "18-25",
  early_adult: "25-35",
  middle_aged_adult: "35-55",
  old_adult: "55+",
};

const SUPPORT_TYPE_OPTIONS = [
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
] as const;
const CASE_STATUS_OPTIONS = ["new", "active", "closed", "archived"] as const;
const SEX_OPTIONS = ["male", "female", "other", "unknown"] as const;
const AGE_GROUP_OPTIONS = [
  "infant",
  "toddler",
  "pre_school",
  "middle_childhood",
  "young_teen",
  "teenager",
  "young_adult",
  "early_adult",
  "middle_aged_adult",
  "old_adult",
] as const;

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

function ReportsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const [params, setParams] = useState<ReportParams>({});
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [activePreset, setActivePreset] = useState<DatePreset | null>(null);
  const { data, isLoading } = useReport(projectId, params);
  const { data: offices } = useOffices();
  const { data: categories } = useCategories();
  const officeOptions = (offices ?? []).map((o) => ({
    label: o.name,
    value: o.id,
  }));

  const categoryOptions = (categories ?? []).map((c) => ({
    label: c.name,
    value: c.id,
  }));

  const supportTypeOptions = SUPPORT_TYPE_OPTIONS.map((s) => ({
    label: t(labelKeyMap[s] ?? s),
    value: s,
  }));

  const caseStatusOptions = CASE_STATUS_OPTIONS.map((s) => ({
    label: t(`project.people.${s}`),
    value: s,
  }));

  const sexOptions = SEX_OPTIONS.map((s) => ({
    label: t(`project.people.sex${s[0].toUpperCase()}${s.slice(1)}`),
    value: s,
  }));

  const ageGroupOptions = AGE_GROUP_OPTIONS.map((g) => ({
    label: t(labelKeyMap[g] ?? g),
    value: g,
  }));

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
                  onValueChange={(v) =>
                    setParams((p) => ({ ...p, category_id: v || undefined }))
                  }
                  options={[
                    { label: t("project.reports.allValues"), value: "" },
                    ...categoryOptions,
                  ]}
                  placeholder={t("project.reports.allValues")}
                  fullWidth
                />
              </FilterField>
              <FilterField label={t("project.reports.filterCaseStatus")}>
                <UISelect
                  value={params.case_status ?? ""}
                  onValueChange={(v) =>
                    setParams((p) => ({ ...p, case_status: v || undefined }))
                  }
                  options={[
                    { label: t("project.reports.allValues"), value: "" },
                    ...caseStatusOptions,
                  ]}
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
                  options={[
                    { label: t("project.reports.allValues"), value: "" },
                    ...ageGroupOptions,
                  ]}
                  placeholder={t("project.reports.allValues")}
                  fullWidth
                />
              </FilterField>
              <FilterField label={t("project.reports.filterSupportType")}>
                <UISelect
                  value={params.support_type ?? ""}
                  onValueChange={(v) =>
                    setParams((p) => ({ ...p, support_type: v || undefined }))
                  }
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

        {/* Active filter chips (always visible) */}
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
            {params.office_id && (
              <FilterChip
                label={t("project.reports.filterOffice")}
                value={
                  officeOptions.find((o) => o.value === params.office_id)?.label ??
                  params.office_id
                }
                onRemove={() => setParams((p) => ({ ...p, office_id: undefined }))}
              />
            )}
            {params.category_id && (
              <FilterChip
                label={t("project.reports.filterCategory")}
                value={
                  categoryOptions.find((c) => c.value === params.category_id)?.label ??
                  params.category_id
                }
                onRemove={() => setParams((p) => ({ ...p, category_id: undefined }))}
              />
            )}
            {params.case_status && (
              <FilterChip
                label={t("project.reports.filterCaseStatus")}
                value={
                  caseStatusOptions.find((s) => s.value === params.case_status)?.label ??
                  params.case_status
                }
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
                value={
                  ageGroupOptions.find((g) => g.value === params.age_group)?.label ??
                  params.age_group
                }
                onRemove={() => setParams((p) => ({ ...p, age_group: undefined }))}
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
              <KpiCard
                label={t("project.reports.kpiConsultations")}
                value={data.consultations.total}
              />
              <KpiCard
                label={t("project.reports.kpiActiveCases")}
                value={data.by_case_status?.rows.find((r) => r.label === "active")?.count ?? 0}
              />
              <KpiCard
                label={t("project.reports.kpiIdp")}
                value={
                  data.by_idp_status.rows.find((r) => r.label === "idp")?.count ??
                  data.by_idp_status.total
                }
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
