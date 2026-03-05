import { useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { BarChart, type BarLegendItem } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import { SankeyChart } from "@/components/charts/sankey-chart";
import { UISelect } from "@/components/ui-select";
import { XIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { useCategories } from "@/hooks/use-categories";
import { useOffices } from "@/hooks/use-offices";
import { usePermissions } from "@/hooks/use-permissions";
import { useReport } from "@/hooks/use-reports";
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
  middle_aged: "project.people.ageMiddleAged",
  older_adult: "project.people.ageOlderAdult",
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
}: {
  group: ReportGroup;
  title: string;
  chart: "bar" | "pie";
  yAxisLabel?: string;
  legend?: BarLegendItem[];

  mapLabel?: (label: string) => string;
  skipTranslation?: boolean;
}) {
  const translated = useTranslatedRows(group.rows);
  const source = skipTranslation ? group.rows : translated;
  const rows = mapLabel ? source.map((r) => ({ ...r, label: mapLabel(r.label) })) : source;
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-3 flex items-baseline justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <span className="tabular-nums text-xs font-medium text-fg-tertiary">
          {group.total.toLocaleString()}
        </span>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart data={rows} yAxisLabel={yAxisLabel} legend={legend} />
        ) : (
          <PieChart data={rows} />
        )
      ) : (
        <p className="py-8 text-center text-sm text-fg-tertiary">—</p>
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

const AGE_RANGE_MAP: Record<string, string> = {
  infant: "<1",
  toddler: "1-2",
  pre_school: "3-5",
  middle_childhood: "6-11",
  young_teen: "12-14",
  teenager: "15-17",
  young_adult: "18-24",
  early_adult: "25-44",
  middle_aged: "45-64",
  older_adult: "65+",
};

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

function ReportsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const [params, setParams] = useState<ReportParams>({});
  const { data, isLoading } = useReport(projectId, params);
  const { data: offices } = useOffices();
  const { data: categories } = useCategories();
  const { data: permissionsData } = usePermissions(projectId);

  const officeOptions = (offices ?? []).map((o) => ({
    label: o.name,
    value: o.id,
  }));

  const categoryOptions = (categories ?? []).map((c) => ({
    label: c.name,
    value: c.id,
  }));

  const consultantOptions = (permissionsData?.permissions ?? []).map((p) => ({
    label: `${p.user_first_name} ${p.user_last_name}`,
    value: p.user_id,
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

  return (
    <div>
      <PageHeader title={t("project.reports.title")} />

      {/* Filter panel */}
      <div className="mb-8 rounded-xl border border-border-secondary bg-bg-secondary p-4">
        <div className="grid grid-cols-2 gap-x-4 gap-y-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-8">
          <FilterField label={t("project.reports.dateFrom")}>
            <DatePicker
              value={params.date_from ?? ""}
              onChange={(v) => setParams((p) => ({ ...p, date_from: v || undefined }))}
            />
          </FilterField>
          <FilterField label={t("project.reports.dateTo")}>
            <DatePicker
              value={params.date_to ?? ""}
              onChange={(v) => setParams((p) => ({ ...p, date_to: v || undefined }))}
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
          <FilterField label={t("project.reports.filterConsultant")}>
            <UISelect
              value={params.consultant_id ?? ""}
              onValueChange={(v) => setParams((p) => ({ ...p, consultant_id: v || undefined }))}
              options={[{ label: t("project.reports.allValues"), value: "" }, ...consultantOptions]}
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
        </div>
        {hasFilters && (
          <div className="mt-3 border-t border-border-secondary pt-3">
            <button
              type="button"
              onClick={() => setParams({})}
              className="inline-flex items-center gap-1.5 rounded-md bg-bg-tertiary px-2.5 py-1 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
            >
              <XIcon size={12} />
              {t("project.reports.clearFilters")}
            </button>
          </div>
        )}
      </div>

      {isLoading && (
        <p className="py-12 text-center text-sm text-fg-tertiary">{t("project.reports.loading")}</p>
      )}

      {data && (
        <div className="grid gap-6 lg:grid-cols-2">
          <ReportCard
            group={data.consultations}
            title={t("project.reports.consultations")}
            chart="bar"
            yAxisLabel={axisLabel}
          />
          <ReportCard group={data.by_sex} title={t("project.reports.bySex")} chart="pie" />
          <ReportCard
            group={data.by_idp_status}
            title={t("project.reports.byIdpStatus")}
            chart="bar"
            yAxisLabel={axisLabel}
          />
          <ReportCard
            group={data.by_category}
            title={t("project.reports.byCategory")}
            chart="bar"
            yAxisLabel={axisLabel}
          />
          <ReportCard
            group={data.by_region}
            title={t("project.reports.byRegion")}
            chart="bar"
            yAxisLabel={axisLabel}
          />
          <ReportCard
            group={data.by_sphere}
            title={t("project.reports.bySphere")}
            chart="bar"
            yAxisLabel={axisLabel}
          />
          <ReportCard
            group={data.by_office}
            title={t("project.reports.byOffice")}
            chart="bar"
            yAxisLabel={axisLabel}
          />
          <ReportCard
            group={data.by_age_group}
            title={t("project.reports.byAgeGroup")}
            chart="bar"
            yAxisLabel={axisLabel}
            skipTranslation
            mapLabel={(l) => AGE_RANGE_MAP[l] ?? l}
            legend={ageGroupLegend}
          />
          <ReportCard
            group={data.by_tag}
            title={t("project.reports.byTag")}
            chart="bar"
            yAxisLabel={axisLabel}
          />
          <ReportCard
            group={data.family_units}
            title={t("project.reports.familyUnits")}
            chart="pie"
          />
          {data.status_flow && data.status_flow.length > 0 && (
            <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5 lg:col-span-2">
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
      )}
    </div>
  );
}
