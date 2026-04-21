import { useTranslation } from "react-i18next";

import { KpiCard } from "@/components/report";

interface PetsKpiCardsProps {
  total: number;
  needsShelter: number;
  adopted: number;
}

export function PetsKpiCards({ total, needsShelter, adopted }: PetsKpiCardsProps) {
  const { t } = useTranslation();
  return (
    <div className="col-span-full grid grid-cols-3 gap-3">
      <KpiCard label={t("project.petReports.kpiTotal")} value={total} />
      <KpiCard label={t("project.petReports.kpiNeedsShelter")} value={needsShelter} />
      <KpiCard label={t("project.petReports.kpiAdopted")} value={adopted} />
    </div>
  );
}
