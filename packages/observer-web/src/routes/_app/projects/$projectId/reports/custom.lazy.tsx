import { useState } from "react";

import { createLazyFileRoute } from "@tanstack/react-router";

import { CustomReportForm } from "@/components/reports/custom-report-form";
import { ReportResultTable } from "@/components/reports/report-result-table";
import { useCustomReport } from "@/hooks/use-reports";
import type { CustomReportParams } from "@/types/report";

export const Route = createLazyFileRoute("/_app/projects/$projectId/reports/custom")({
  component: CustomReportPage,
});

function CustomReportPage() {
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

  return (
    <div>
      <CustomReportForm
        metric={metric}
        groupBy={groupBy}
        dateFrom={dateFrom}
        dateTo={dateTo}
        supportType={supportType}
        isFetching={isFetching}
        onMetricChange={(v) => {
          setMetric(v);
          setSubmitted(false);
        }}
        onToggleDimension={toggleDimension}
        onDateRangeChange={(from, to) => {
          setDateFrom(from);
          setDateTo(to);
          setSubmitted(false);
        }}
        onSupportTypeChange={(v) => {
          setSupportType(v);
          setSubmitted(false);
        }}
        onGenerate={handleGenerate}
      />

      {/* Results */}
      {isLoading && submitted && (
        <div className="space-y-4">
          <div className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
          <div className="h-64 animate-pulse rounded-xl bg-bg-tertiary" />
        </div>
      )}

      {data && submitted && <ReportResultTable data={data} />}
    </div>
  );
}
