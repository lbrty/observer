import { useTranslation } from "react-i18next";

import { type BarLegendItem } from "@/components/charts/bar-chart";
import {
  SEX_COLORS,
  SUPPORT_TYPE_COLORS,
  SPHERE_COLORS,
  AGE_GROUP_COLORS,
} from "@/components/charts/colors";
import { MyStatsKpiCards } from "@/components/my-stats/my-stats-kpi-cards";
import { ReportCard, AGE_RANGE_MAP } from "@/components/reports/shared";
import { useReport } from "@/hooks/reports/use-reports";

interface MyStatsChartSectionProps {
  data: NonNullable<ReturnType<typeof useReport>["data"]>;
  axisLabel: string;
  ageGroupLegend: BarLegendItem[];
}

export function MyStatsChartSection({ data, axisLabel, ageGroupLegend }: MyStatsChartSectionProps) {
  const { t } = useTranslation();

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <MyStatsKpiCards
        totalPeople={data.by_sex.total}
        totalConsultations={data.consultations.total}
        totalActiveCases={data.by_case_status?.rows.find((r) => r.label === "active")?.count ?? 0}
        totalHouseholds={data.family_units.total}
      />

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
        group={data.by_office}
        title={t("project.reports.byOffice")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />
      <ReportCard
        group={data.by_region}
        title={t("project.reports.byRegion")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />

      <ReportCard
        group={data.by_sex}
        title={t("project.reports.bySex")}
        chart="pie"
        colorMap={SEX_COLORS}
      />
      <ReportCard
        group={data.by_case_status}
        title={t("project.reports.byCaseStatus")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />

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

      <ReportCard
        group={data.by_category}
        title={t("project.reports.byCategory")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />
      <ReportCard
        group={data.by_tag}
        title={t("project.reports.byTag")}
        chart="bar"
        yAxisLabel={axisLabel}
        direction="auto"
      />
    </div>
  );
}
