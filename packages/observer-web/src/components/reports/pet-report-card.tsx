import { useTranslation } from "react-i18next";

import { BarChart } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import { DownloadSimpleIcon } from "@/components/ui/icons";
import type { CountResult, MonthlyStatusCount } from "@/types/report";

export function extractMonthlySeriesForStatus(
  data: MonthlyStatusCount[],
  status: string,
): CountResult[] {
  return data.filter((r) => r.status === status).map((r) => ({ label: r.month, count: r.count }));
}

export function extractMonthlyTotals(data: MonthlyStatusCount[]): CountResult[] {
  const totals = new Map<string, number>();
  for (const r of data) {
    totals.set(r.month, (totals.get(r.month) ?? 0) + r.count);
  }
  return Array.from(totals, ([label, count]) => ({ label, count })).sort((a, b) =>
    a.label.localeCompare(b.label),
  );
}

interface ReportCardProps {
  title: string;
  rows: CountResult[];
  chart: "bar" | "pie";
  colorMap?: Record<string, string>;
  direction?: "vertical" | "horizontal" | "auto";
  yAxisLabel?: string;
  onExport?: () => void;
  total?: number;
}

export function ReportCard({
  title,
  rows,
  chart,
  colorMap,
  direction,
  yAxisLabel,
  onExport,
  total,
}: ReportCardProps) {
  const { t } = useTranslation();
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-3 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <div className="flex items-center gap-2">
          {onExport && (
            <button
              type="button"
              onClick={onExport}
              className="text-fg-tertiary transition-colors hover:text-fg"
              title={t("common.downloadCsv")}
            >
              <DownloadSimpleIcon size={14} />
            </button>
          )}
          {total != null && (
            <span className="tabular-nums text-xs font-medium text-fg-tertiary">
              {total.toLocaleString()}
            </span>
          )}
        </div>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart data={rows} colorMap={colorMap} direction={direction} yAxisLabel={yAxisLabel} />
        ) : (
          <PieChart data={rows} colorMap={colorMap} />
        )
      ) : (
        <p className="py-8 text-center text-sm text-fg-tertiary">&mdash;</p>
      )}
    </div>
  );
}
