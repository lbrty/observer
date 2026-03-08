import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { CustomReportOutput, CustomReportParams, FullReport, ReportParams } from "@/types/report";

export function useReport(projectId: string, params: ReportParams = {}) {
  return useQuery({
    queryKey: ["reports", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/reports`, {
          searchParams: Object.fromEntries(
            Object.entries(params).filter(([, v]) => v != null && v !== ""),
          ) as Record<string, string>,
        })
        .json<FullReport>(),
    enabled: !!projectId,
  });
}

function buildCustomSearchParams(params: CustomReportParams): URLSearchParams {
  const sp = new URLSearchParams();
  sp.set("metric", params.metric);
  for (const g of params.group_by) {
    sp.append("group_by", g);
  }
  if (params.date_from) sp.set("date_from", params.date_from);
  if (params.date_to) sp.set("date_to", params.date_to);
  if (params.support_type) sp.set("support_type", params.support_type);
  return sp;
}

export function useCustomReport(projectId: string, params: CustomReportParams, enabled: boolean) {
  return useQuery({
    queryKey: ["custom-report", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/reports/custom`, {
          searchParams: buildCustomSearchParams(params),
        })
        .json<CustomReportOutput>(),
    enabled: enabled && params.group_by.length > 0,
  });
}
