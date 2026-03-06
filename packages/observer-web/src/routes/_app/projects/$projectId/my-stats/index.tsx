import { useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { BarChart, type BarLegendItem } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
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
  XIcon,
} from "@/components/icons";
import { useReport } from "@/hooks/use-reports";
import { exportGroupCSV, exportReportCSV } from "@/lib/export-csv";
import { useAuth } from "@/stores/auth";
import type { CountResult, ReportGroup, ReportParams } from "@/types/report";

export const Route = createFileRoute("/_app/projects/$projectId/my-stats/")({
  component: MyStatsPage,
});

const labelKeyMap: Record<string, string> = {
  humanitarian: "project.supportRecords.typeHumanitarian",
  legal: "project.supportRecords.typeLegal",
  social: "project.supportRecords.typeSocial",
  psychological: "project.supportRecords.typePsychological",
  medical: "project.supportRecords.typeMedical",
  general: "project.supportRecords.typeGeneral",
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
  male: "project.people.sexMale",
  female: "project.people.sexFemale",
  unknown: "project.people.sexUnknown",
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

function StatsCard({
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

function StatsSkeleton() {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
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

function MyStatsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const { user } = useAuth();
  const [params, setParams] = useState<ReportParams>({});
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [activePreset, setActivePreset] = useState<DatePreset | null>(null);

  const reportParams: ReportParams = { ...params, consultant_id: user?.id };
  const { data, isLoading } = useReport(projectId, reportParams);

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

            <div className="grid grid-cols-2 gap-x-4 gap-y-3">
              <div className="space-y-1.5">
                <span className="block text-xs font-medium text-fg-secondary">
                  {t("project.reports.dateFrom")}
                </span>
                <DatePicker
                  value={params.date_from ?? ""}
                  onChange={(v) => {
                    setParams((p) => ({ ...p, date_from: v || undefined }));
                    clearDatePreset();
                  }}
                />
              </div>
              <div className="space-y-1.5">
                <span className="block text-xs font-medium text-fg-secondary">
                  {t("project.reports.dateTo")}
                </span>
                <DatePicker
                  value={params.date_to ?? ""}
                  onChange={(v) => {
                    setParams((p) => ({ ...p, date_to: v || undefined }));
                    clearDatePreset();
                  }}
                />
              </div>
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

      {isLoading && <StatsSkeleton />}

      {data && (
        <div className="grid gap-6 lg:grid-cols-2">
          {/* KPI overview */}
          <div className="col-span-full grid grid-cols-2 gap-3 sm:grid-cols-4">
            <KpiCard label={t("project.myStats.kpiPeople")} value={data.by_sex.total} />
            <KpiCard
              label={t("project.myStats.kpiConsultations")}
              value={data.consultations.total}
            />
            <KpiCard
              label={t("project.myStats.kpiActiveCases")}
              value={data.by_case_status?.rows.find((r) => r.label === "active")?.count ?? 0}
            />
            <KpiCard label={t("project.myStats.kpiHouseholds")} value={data.family_units.total} />
          </div>

          {/* Consultations */}
          <div className="col-span-full">
            <StatsCard
              group={data.consultations}
              title={t("project.reports.consultations")}
              chart="bar"
              yAxisLabel={axisLabel}
              colorMap={SUPPORT_TYPE_COLORS}
            />
          </div>

          {/* Service breakdown */}
          <StatsCard
            group={data.by_sphere}
            title={t("project.reports.bySphere")}
            chart="bar"
            yAxisLabel={axisLabel}
            colorMap={SPHERE_COLORS}
            direction="auto"
          />

          {/* Demographics */}
          <StatsCard
            group={data.by_sex}
            title={t("project.reports.bySex")}
            chart="pie"
            colorMap={SEX_COLORS}
          />

          {/* Age distribution */}
          <div className="col-span-full">
            <StatsCard
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
          <StatsCard
            group={data.by_category}
            title={t("project.reports.byCategory")}
            chart="bar"
            yAxisLabel={axisLabel}
            direction="auto"
          />
          <StatsCard
            group={data.by_tag}
            title={t("project.reports.byTag")}
            chart="bar"
            yAxisLabel={axisLabel}
            direction="auto"
          />
        </div>
      )}
    </div>
  );
}
