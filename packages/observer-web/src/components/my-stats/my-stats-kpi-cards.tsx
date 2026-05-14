import { useTranslation } from "react-i18next";

import { KpiCard } from "@/components/reports/shared";

interface MyStatsKpiCardsProps {
  totalPeople: number;
  totalConsultations: number;
  totalActiveCases: number;
  totalHouseholds: number;
}

export function MyStatsKpiCards({
  totalPeople,
  totalConsultations,
  totalActiveCases,
  totalHouseholds,
}: MyStatsKpiCardsProps) {
  const { t } = useTranslation();

  return (
    <div className="col-span-full grid grid-cols-2 gap-3 sm:grid-cols-4">
      <KpiCard label={t("project.myStats.kpiPeople")} value={totalPeople} index={0} />
      <KpiCard label={t("project.myStats.kpiConsultations")} value={totalConsultations} index={1} />
      <KpiCard label={t("project.myStats.kpiActiveCases")} value={totalActiveCases} index={2} />
      <KpiCard label={t("project.myStats.kpiHouseholds")} value={totalHouseholds} index={3} />
    </div>
  );
}
