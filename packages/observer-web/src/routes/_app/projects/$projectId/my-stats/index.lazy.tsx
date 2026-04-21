import { useState } from "react";

import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { type BarLegendItem } from "@/components/charts/bar-chart";
import { MyStatsChartSection } from "@/components/my-stats/my-stats-chart-section";
import { MyStatsFilters } from "@/components/my-stats/my-stats-filters";
import { AGE_RANGE_MAP, ReportSkeleton, labelKeyMap } from "@/components/reports/shared";
import type { DatePreset } from "@/components/reports/shared";
import { useReport } from "@/hooks/reports/use-reports";
import { useAuth } from "@/stores/auth";
import type { ReportParams } from "@/types/report";

export const Route = createLazyFileRoute("/_app/projects/$projectId/my-stats/")({
  component: MyStatsPage,
});

const SUPPORT_TYPE_OPTIONS = [
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
] as const;

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

  const axisLabel = t("project.reports.axisCount");

  const ageGroupLegend: BarLegendItem[] = Object.entries(AGE_RANGE_MAP).map(([key, range]) => ({
    short: range,
    full: t(labelKeyMap[key] ?? key),
  }));

  return (
    <div>
      <MyStatsFilters
        projectId={projectId}
        params={params}
        setParams={setParams}
        filtersOpen={filtersOpen}
        setFiltersOpen={setFiltersOpen}
        activePreset={activePreset}
        setActivePreset={setActivePreset}
        data={data}
        supportTypeOptions={supportTypeOptions}
      />

      {isLoading && <ReportSkeleton kpiCount={4} />}

      {data && (
        <MyStatsChartSection data={data} axisLabel={axisLabel} ageGroupLegend={ageGroupLegend} />
      )}
    </div>
  );
}
