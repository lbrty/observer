import { useTranslation } from "react-i18next";

import { DownloadSimpleIcon } from "@/components/ui/icons";
import { labelKeyMap } from "@/components/reports/shared";
import type { CustomReportOutput } from "@/types/report";

import { DIMENSION_LABEL_KEYS } from "./custom-report-form";

function escapeCSV(value: string): string {
  if (value.includes(",") || value.includes('"') || value.includes("\n")) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}

interface ReportResultTableProps {
  data: CustomReportOutput;
}

export function ReportResultTable({ data }: ReportResultTableProps) {
  const { t } = useTranslation();

  function exportCSV() {
    if (data.rows.length === 0) return;
    const dims = data.group_by;
    const header = [...dims, "count"].join(",");
    const lines = [header];
    for (const row of data.rows) {
      const vals = dims.map((d) => escapeCSV(row.dimensions[d] ?? ""));
      lines.push([...vals, String(row.count)].join(","));
    }
    lines.push([...dims.map(() => ""), String(data.total)].join(","));
    const blob = new Blob([lines.join("\n")], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `custom-report-${new Date().toISOString().slice(0, 10)}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  }

  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary">
      {/* Total + export */}
      <div className="flex items-center justify-between px-5 py-4">
        <div>
          <p className="text-3xl font-bold tabular-nums text-fg">{data.total.toLocaleString()}</p>
          <p className="mt-0.5 text-xs font-medium text-fg-tertiary">
            {t("project.customReport.total")}
          </p>
        </div>
        {data.rows.length > 0 && (
          <button
            type="button"
            onClick={exportCSV}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border-secondary px-3 py-1.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
          >
            <DownloadSimpleIcon size={14} />
            {t("project.reports.exportCsv")}
          </button>
        )}
      </div>

      {/* Table */}
      {data.rows.length > 0 ? (
        <div className="overflow-x-auto border-t border-border-secondary">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-secondary bg-bg-tertiary text-left">
                {data.group_by.map((dim) => (
                  <th key={dim} className="px-4 py-2.5 font-medium text-fg-secondary">
                    {t(DIMENSION_LABEL_KEYS[dim] ?? dim)}
                  </th>
                ))}
                <th className="px-4 py-2.5 text-right font-medium text-fg-secondary">
                  {t("project.reports.axisCount")}
                </th>
              </tr>
            </thead>
            <tbody>
              {data.rows.map((row, ix) => (
                <tr key={ix} className="border-b border-border-secondary last:border-b-0">
                  {data.group_by.map((dim) => (
                    <td key={dim} className="px-4 py-2.5 text-fg">
                      {(() => {
                        const val = row.dimensions[dim];
                        if (!val) return "—";
                        const key = labelKeyMap[val];
                        return key ? t(key) : val;
                      })()}
                    </td>
                  ))}
                  <td className="px-4 py-2.5 text-right tabular-nums font-medium text-fg">
                    {row.count.toLocaleString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="border-t border-border-secondary px-5 py-12 text-center text-sm text-fg-tertiary">
          {t("project.customReport.noData")}
        </div>
      )}
    </div>
  );
}
