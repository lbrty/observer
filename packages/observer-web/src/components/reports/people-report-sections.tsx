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
import { ReportCard, labelKeyMap, AGE_RANGE_MAP } from "@/components/reports/shared";
import { PeopleKpiCards } from "@/components/reports/people-kpi-cards";
import { PeopleChartSection } from "@/components/reports/people-chart-section";
import { useReport } from "@/hooks/reports/use-reports";

function SectionHeader({ title }: { title: string }) {
  return (
    <div className="col-span-full border-b border-border-secondary pb-1 pt-4">
      <h2 className="text-xs font-semibold uppercase tracking-wider text-fg-tertiary">{title}</h2>
    </div>
  );
}

interface PeopleReportSectionsProps {
  data: NonNullable<ReturnType<typeof useReport>["data"]>;
  axisLabel: string;
  ageGroupLegend: BarLegendItem[];
}

export function PeopleReportSections({
  data,
  axisLabel,
  ageGroupLegend,
}: PeopleReportSectionsProps) {
  const { t } = useTranslation();

  return (
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
  );
}
