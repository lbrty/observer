import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import { filterParams } from "@/lib/params";
import type { PetReport, PetReportParams } from "@/types/report";

export function usePetReport(projectId: string, params: PetReportParams = {}) {
  return useQuery({
    queryKey: ["pet-reports", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/reports/pets`, {
          searchParams: filterParams(params as Record<string, unknown>),
        })
        .json<PetReport>(),
    enabled: !!projectId,
  });
}
