import { useTranslation } from "react-i18next";

import { PET_STATUS_COLORS, PET_OWNERSHIP_COLORS } from "@/components/charts/colors";
import { ReportCard } from "@/components/reports/pet-report-card";
import { PetsKpiCards } from "@/components/reports/pets-kpi-cards";
import { exportGroupCSV } from "@/lib/export-csv";
import type { CountResult, PetReport } from "@/types/report";

interface PetsChartSectionProps {
  data: PetReport;
  axisLabel: string;
  translatedStatus: CountResult[];
  translatedOwnership: CountResult[];
  needsShelterMonthly: CountResult[];
  adoptedMonthly: CountResult[];
  totalMonthly: CountResult[];
}

export function PetsChartSection({
  data,
  axisLabel,
  translatedStatus,
  translatedOwnership,
  needsShelterMonthly,
  adoptedMonthly,
  totalMonthly,
}: PetsChartSectionProps) {
  const { t } = useTranslation();
  const needsShelterCount =
    data.by_status.rows.find((r) => r.label === "needs_shelter")?.count ?? 0;
  const adoptedCount = data.by_status.rows.find((r) => r.label === "adopted")?.count ?? 0;

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <PetsKpiCards
        total={data.by_status.total}
        needsShelter={needsShelterCount}
        adopted={adoptedCount}
      />

      <ReportCard
        title={t("project.petReports.byStatus")}
        rows={translatedStatus}
        chart="pie"
        colorMap={PET_STATUS_COLORS}
        total={data.by_status.total}
        onExport={() => exportGroupCSV(t("project.petReports.byStatus"), translatedStatus)}
      />
      <ReportCard
        title={t("project.petReports.byOwnership")}
        rows={translatedOwnership}
        chart="pie"
        colorMap={PET_OWNERSHIP_COLORS}
        total={data.by_ownership.total}
        onExport={() => exportGroupCSV(t("project.petReports.byOwnership"), translatedOwnership)}
      />

      <div className="col-span-full">
        <ReportCard
          title={t("project.petReports.byMonth")}
          rows={data.by_month.rows}
          chart="bar"
          yAxisLabel={axisLabel}
          total={data.by_month.total}
          onExport={() => exportGroupCSV(t("project.petReports.byMonth"), data.by_month.rows)}
        />
      </div>

      <div className="col-span-full">
        <ReportCard
          title={t("project.petReports.totalByMonth")}
          rows={totalMonthly}
          chart="bar"
          yAxisLabel={axisLabel}
          onExport={() => exportGroupCSV(t("project.petReports.totalByMonth"), totalMonthly)}
        />
      </div>

      <ReportCard
        title={t("project.petReports.needsShelterByMonth")}
        rows={needsShelterMonthly}
        chart="bar"
        yAxisLabel={axisLabel}
        colorMap={{
          ...Object.fromEntries(needsShelterMonthly.map((r) => [r.label, "#e5534b"])),
        }}
        onExport={() =>
          exportGroupCSV(t("project.petReports.needsShelterByMonth"), needsShelterMonthly)
        }
      />
      <ReportCard
        title={t("project.petReports.adoptedByMonth")}
        rows={adoptedMonthly}
        chart="bar"
        yAxisLabel={axisLabel}
        colorMap={{ ...Object.fromEntries(adoptedMonthly.map((r) => [r.label, "#30a46c"])) }}
        onExport={() => exportGroupCSV(t("project.petReports.adoptedByMonth"), adoptedMonthly)}
      />
    </div>
  );
}
