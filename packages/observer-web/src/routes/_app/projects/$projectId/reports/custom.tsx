import { useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { DatePicker } from "@/components/date-picker";
import { DownloadSimpleIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";
import { useCustomReport } from "@/hooks/use-reports";
import type { CustomReportParams } from "@/types/report";

export const Route = createFileRoute("/_app/projects/$projectId/reports/custom")({
  component: CustomReportPage,
});

const METRICS = ["events", "people", "units", "pets"] as const;

const DIMENSIONS = [
  "sex",
  "age_group",
  "region",
  "conflict_zone",
  "office",
  "sphere",
  "category",
  "person_tag",
  "pet_tag",
  "pet_status",
] as const;

const SUPPORT_TYPE_OPTIONS = ["legal", "social"] as const;

const DIMENSION_LABEL_KEYS: Record<string, string> = {
  sex: "project.customReport.dimSex",
  age_group: "project.customReport.dimAgeGroup",
  region: "project.customReport.dimRegion",
  conflict_zone: "project.customReport.dimConflictZone",
  office: "project.customReport.dimOffice",
  sphere: "project.customReport.dimSphere",
  category: "project.customReport.dimCategory",
  person_tag: "project.customReport.dimPersonTag",
  pet_tag: "project.customReport.dimPetTag",
  pet_status: "project.customReport.dimPetStatus",
};

function CustomReportPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();

  const [metric, setMetric] = useState<CustomReportParams["metric"]>("events");
  const [groupBy, setGroupBy] = useState<string[]>([]);
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const [supportType, setSupportType] = useState("");
  const [submitted, setSubmitted] = useState(false);

  const params: CustomReportParams = {
    metric,
    group_by: groupBy,
    date_from: dateFrom || undefined,
    date_to: dateTo || undefined,
    support_type: supportType || undefined,
  };

  const { data, isLoading, isFetching } = useCustomReport(projectId, params, submitted);

  function toggleDimension(dim: string) {
    setGroupBy((prev) => {
      if (prev.includes(dim)) return prev.filter((d) => d !== dim);
      if (prev.length >= 2) return prev;
      return [...prev, dim];
    });
    setSubmitted(false);
  }

  function handleGenerate() {
    if (groupBy.length === 0) return;
    setSubmitted(true);
  }

  function exportCSV() {
    if (!data || data.rows.length === 0) return;
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
    <div>
      <div className="mb-6 rounded-xl border border-border-secondary bg-bg-secondary">
        <div className="px-5 py-4">
          <h1 className="font-serif text-xl font-bold tracking-tight text-fg">
            {t("project.customReport.title")}
          </h1>
        </div>

        <div className="border-t border-border-secondary px-5 pb-5 pt-4 space-y-5">
          {/* Metric */}
          <div>
            <span className="mb-2 block text-xs font-medium text-fg-secondary">
              {t("project.customReport.metric")}
            </span>
            <div className="flex flex-wrap gap-2">
              {METRICS.map((m) => (
                <button
                  key={m}
                  type="button"
                  onClick={() => { setMetric(m); setSubmitted(false); }}
                  className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                    metric === m
                      ? "bg-accent text-accent-fg"
                      : "bg-bg-tertiary text-fg-secondary hover:text-fg"
                  }`}
                >
                  {t(`project.customReport.metric_${m}`)}
                </button>
              ))}
            </div>
          </div>

          {/* Dimensions */}
          <div>
            <span className="mb-2 block text-xs font-medium text-fg-secondary">
              {t("project.customReport.dimensions")}
              <span className="ml-1 text-fg-tertiary">({groupBy.length}/2)</span>
            </span>
            <div className="flex flex-wrap gap-2">
              {DIMENSIONS.map((dim) => {
                const selected = groupBy.includes(dim);
                const disabled = !selected && groupBy.length >= 2;
                return (
                  <button
                    key={dim}
                    type="button"
                    disabled={disabled}
                    onClick={() => toggleDimension(dim)}
                    className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                      selected
                        ? "bg-accent text-accent-fg"
                        : disabled
                          ? "cursor-not-allowed bg-bg-tertiary text-fg-quaternary"
                          : "bg-bg-tertiary text-fg-secondary hover:text-fg"
                    }`}
                  >
                    {t(DIMENSION_LABEL_KEYS[dim])}
                  </button>
                );
              })}
            </div>
          </div>

          {/* Date range + support type */}
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
            <div className="space-y-1.5">
              <span className="block text-xs font-medium text-fg-secondary">
                {t("project.reports.dateFrom")}
              </span>
              <DatePicker
                value={dateFrom}
                onChange={(v) => { setDateFrom(v); setSubmitted(false); }}
              />
            </div>
            <div className="space-y-1.5">
              <span className="block text-xs font-medium text-fg-secondary">
                {t("project.reports.dateTo")}
              </span>
              <DatePicker
                value={dateTo}
                onChange={(v) => { setDateTo(v); setSubmitted(false); }}
              />
            </div>
            <div className="space-y-1.5">
              <span className="block text-xs font-medium text-fg-secondary">
                {t("project.reports.filterSupportType")}
              </span>
              <UISelect
                value={supportType}
                onValueChange={(v) => { setSupportType(v); setSubmitted(false); }}
                options={[
                  { label: t("project.reports.allValues"), value: "" },
                  ...SUPPORT_TYPE_OPTIONS.map((s) => ({
                    label: t(`project.customReport.supportType_${s}`),
                    value: s,
                  })),
                ]}
                placeholder={t("project.reports.allValues")}
                fullWidth
              />
            </div>
          </div>

          {/* Generate button */}
          <button
            type="button"
            disabled={groupBy.length === 0 || isFetching}
            onClick={handleGenerate}
            className="rounded-lg bg-accent px-5 py-2 text-sm font-semibold text-accent-fg transition-colors hover:bg-accent/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {isFetching ? t("project.customReport.generating") : t("project.customReport.generate")}
          </button>
        </div>
      </div>

      {/* Results */}
      {isLoading && submitted && (
        <div className="space-y-4">
          <div className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
          <div className="h-64 animate-pulse rounded-xl bg-bg-tertiary" />
        </div>
      )}

      {data && submitted && (
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
                          {row.dimensions[dim] ?? "\u2014"}
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
      )}
    </div>
  );
}

function escapeCSV(value: string): string {
  if (value.includes(",") || value.includes('"') || value.includes("\n")) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}
