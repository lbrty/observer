import { useTranslation } from "react-i18next";

import { KpiCard } from "@/components/reports/shared";

interface PeopleKpiCardsProps {
  totalPeople: number;
  totalConsultations: number;
  activeCases: number;
  idpCount: number;
  households: number;
  offices: number;
}

export function PeopleKpiCards({
  totalPeople,
  totalConsultations,
  activeCases,
  idpCount,
  households,
  offices,
}: PeopleKpiCardsProps) {
  const { t } = useTranslation();

  return (
    <div className="grid grid-cols-3 gap-3">
      <KpiCard label={t("project.reports.kpiPeople")} value={totalPeople} index={0} />
      <KpiCard label={t("project.reports.kpiConsultations")} value={totalConsultations} index={1} />
      <KpiCard label={t("project.reports.kpiActiveCases")} value={activeCases} index={2} />
      <KpiCard label={t("project.reports.kpiIdp")} value={idpCount} index={3} />
      <KpiCard label={t("project.reports.kpiHouseholds")} value={households} index={4} />
      <KpiCard label={t("project.reports.kpiOffices")} value={offices} index={5} />
    </div>
  );
}
