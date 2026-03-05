import { useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { BarChart } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import { PageHeader } from "@/components/page-header";
import { useReport } from "@/hooks/use-reports";
import type { CountResult, ReportGroup, ReportParams } from "@/types/report";

export const Route = createFileRoute(
  "/_app/projects/$projectId/reports/",
)({
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
  return rows.map((r) => {
    const key = labelKeyMap[r.label];
    return key ? { ...r, label: t(key) } : r;
  });
}

function ReportCard({
  group,
  title,
  chart,
}: {
  group: ReportGroup;
  title: string;
  chart: "bar" | "pie";
}) {
  const rows = useTranslatedRows(group.rows);
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-1 flex items-baseline justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <span className="text-xs text-fg-tertiary">
          {group.total}
        </span>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart data={rows} />
        ) : (
          <PieChart data={rows} />
        )
      ) : (
        <p className="py-8 text-center text-sm text-fg-tertiary">—</p>
      )}
    </div>
  );
}

function ReportsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const [params, setParams] = useState<ReportParams>({});
  const { data, isLoading } = useReport(projectId, params);

  return (
    <div>
      <PageHeader title={t("project.reports.title")} />

      <div className="mb-6 flex items-end gap-4">
        <div className="space-y-1">
          <span className="text-xs font-medium text-fg-secondary">
            {t("project.reports.dateFrom")}
          </span>
          <DatePicker
            value={params.date_from ?? ""}
            onChange={(v) =>
              setParams((p) => ({ ...p, date_from: v || undefined }))
            }
          />
        </div>
        <div className="space-y-1">
          <span className="text-xs font-medium text-fg-secondary">
            {t("project.reports.dateTo")}
          </span>
          <DatePicker
            value={params.date_to ?? ""}
            onChange={(v) =>
              setParams((p) => ({ ...p, date_to: v || undefined }))
            }
          />
        </div>
        {(params.date_from || params.date_to) && (
          <button
            type="button"
            onClick={() => setParams({})}
            className="text-xs text-accent hover:underline"
          >
            {t("project.reports.clearFilters")}
          </button>
        )}
      </div>

      {isLoading && (
        <p className="py-12 text-center text-sm text-fg-tertiary">
          {t("project.reports.loading")}
        </p>
      )}

      {data && (
        <div className="grid gap-6 lg:grid-cols-2">
          <ReportCard
            group={data.consultations}
            title={t("project.reports.consultations")}
            chart="bar"
          />
          <ReportCard
            group={data.by_sex}
            title={t("project.reports.bySex")}
            chart="pie"
          />
          <ReportCard
            group={data.by_idp_status}
            title={t("project.reports.byIdpStatus")}
            chart="bar"
          />
          <ReportCard
            group={data.by_category}
            title={t("project.reports.byCategory")}
            chart="bar"
          />
          <ReportCard
            group={data.by_region}
            title={t("project.reports.byRegion")}
            chart="bar"
          />
          <ReportCard
            group={data.by_sphere}
            title={t("project.reports.bySphere")}
            chart="bar"
          />
          <ReportCard
            group={data.by_office}
            title={t("project.reports.byOffice")}
            chart="bar"
          />
          <ReportCard
            group={data.by_age_group}
            title={t("project.reports.byAgeGroup")}
            chart="bar"
          />
          <ReportCard
            group={data.by_tag}
            title={t("project.reports.byTag")}
            chart="bar"
          />
          <ReportCard
            group={data.family_units}
            title={t("project.reports.familyUnits")}
            chart="pie"
          />
        </div>
      )}
    </div>
  );
}
