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
      <KpiCard label={t("project.reports.kpiPeople")} value={totalPeople} />
      <KpiCard label={t("project.reports.kpiConsultations")} value={totalConsultations} />
      <KpiCard label={t("project.reports.kpiActiveCases")} value={activeCases} />
      <KpiCard label={t("project.reports.kpiIdp")} value={idpCount} />
      <KpiCard label={t("project.reports.kpiHouseholds")} value={households} />
      <KpiCard label={t("project.reports.kpiOffices")} value={offices} />
    </div>
  );
}
