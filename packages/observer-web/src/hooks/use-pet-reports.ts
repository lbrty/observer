import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { PetReport, PetReportParams } from "@/types/report";

export function usePetReport(projectId: string, params: PetReportParams = {}) {
  return useQuery({
    queryKey: ["pet-reports", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/reports/pets`, {
          searchParams: Object.fromEntries(
            Object.entries(params).filter(([, v]) => v != null && v !== ""),
          ) as Record<string, string>,
        })
        .json<PetReport>(),
    enabled: !!projectId,
  });
}
