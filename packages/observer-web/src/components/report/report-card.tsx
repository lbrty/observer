import { BarChart, type BarLegendItem } from "@/components/charts/bar-chart";
import { PieChart } from "@/components/charts/pie-chart";
import { DownloadSimpleIcon } from "@/components/icons";
import { exportGroupCSV } from "@/lib/export-csv";
import { useTranslatedRows } from "@/components/report/label-maps";
import type { ReportGroup } from "@/types/report";

interface ReportCardProps {
  group: ReportGroup;
  title: string;
  chart: "bar" | "pie";
  yAxisLabel?: string;
  legend?: BarLegendItem[];
  mapLabel?: (label: string) => string;
  skipTranslation?: boolean;
  colorMap?: Record<string, string>;
  direction?: "vertical" | "horizontal" | "auto";
}

export function ReportCard({
  group, title, chart, yAxisLabel, legend, mapLabel, skipTranslation, colorMap, direction,
}: ReportCardProps) {
  const translated = useTranslatedRows(group.rows);
  const source = skipTranslation ? group.rows : translated;
  const rows = mapLabel ? source.map((r) => ({ ...r, label: mapLabel(r.label) })) : source;
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-3 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => exportGroupCSV(title, rows)}
            className="text-fg-tertiary transition-colors hover:text-fg"
            title="Download CSV"
          >
            <DownloadSimpleIcon size={14} />
          </button>
          <span className="tabular-nums text-xs font-medium text-fg-tertiary">
            {group.total.toLocaleString()}
          </span>
        </div>
      </div>
      {rows.length > 0 ? (
        chart === "bar" ? (
          <BarChart data={rows} yAxisLabel={yAxisLabel} legend={legend} colorMap={colorMap} direction={direction} />
        ) : (
          <PieChart data={rows} colorMap={colorMap} />
        )
      ) : (
        <p className="py-8 text-center text-sm text-fg-tertiary">&mdash;</p>
      )}
    </div>
  );
}
