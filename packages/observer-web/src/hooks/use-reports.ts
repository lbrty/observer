import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { FullReport, ReportParams } from "@/types/report";

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
