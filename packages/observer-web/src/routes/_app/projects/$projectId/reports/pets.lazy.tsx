import { useState } from "react";

import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { ReportSkeleton, useTranslatedRows, type DatePreset } from "@/components/reports/shared";
import {
  extractMonthlySeriesForStatus,
  extractMonthlyTotals,
} from "@/components/reports/pet-report-card";
import { PetReportFilters } from "@/components/reports/pet-report-filters";
import { PetsChartSection } from "@/components/reports/pets-chart-section";
import { petStatusKeys } from "@/constants/i18n";
import { usePetReport } from "@/hooks/reports/use-pet-reports";

import type { PetReportParams } from "@/types/report";

export const Route = createLazyFileRoute("/_app/projects/$projectId/reports/pets")({
  component: PetReportsPage,
});

const PET_STATUS_OPTIONS = [
  "registered",
  "adopted",
  "owner_found",
  "needs_shelter",
  "unknown",
] as const;

function PetReportsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const [params, setParams] = useState<PetReportParams>({});
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [activePreset, setActivePreset] = useState<DatePreset | null>(null);
  const { data, isLoading } = usePetReport(projectId, params);

  const statusOptions = PET_STATUS_OPTIONS.map((s) => ({
    label: t(petStatusKeys[s] ?? s),
    value: s,
  }));

  const axisLabel = t("project.reports.axisCount");

  const translatedStatus = useTranslatedRows(data?.by_status.rows ?? []);
  const translatedOwnership = useTranslatedRows(data?.by_ownership.rows ?? []);

  const needsShelterMonthly = data
    ? extractMonthlySeriesForStatus(data.by_status_by_month, "needs_shelter")
    : [];
  const adoptedMonthly = data
    ? extractMonthlySeriesForStatus(data.by_status_by_month, "adopted")
    : [];
  const totalMonthly = data ? extractMonthlyTotals(data.by_status_by_month) : [];

  return (
    <div>
      <PetReportFilters
        params={params}
        setParams={setParams}
        filtersOpen={filtersOpen}
        setFiltersOpen={setFiltersOpen}
        activePreset={activePreset}
        setActivePreset={setActivePreset}
        statusOptions={statusOptions}
        data={data}
      />

      {isLoading && <ReportSkeleton kpiCount={3} />}

      {data && (
        <PetsChartSection
          data={data}
          axisLabel={axisLabel}
          translatedStatus={translatedStatus}
          translatedOwnership={translatedOwnership}
          needsShelterMonthly={needsShelterMonthly}
          adoptedMonthly={adoptedMonthly}
          totalMonthly={totalMonthly}
        />
      )}
    </div>
  );
}
