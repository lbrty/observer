import { keepPreviousData, useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import { filterParams } from "@/lib/params";
import type { AuditListOutput, AuditListParams } from "@/types/audit";

export function useAuditLogs(params: AuditListParams) {
  return useQuery({
    queryKey: ["audit-logs", params],
    queryFn: () =>
      api
        .get("admin/audit-logs", { searchParams: filterParams(params as Record<string, unknown>) })
        .json<AuditListOutput>(),
    placeholderData: keepPreviousData,
  });
}

export function useProjectAuditLogs(projectId: string, params: AuditListParams) {
  return useQuery({
    queryKey: ["project-audit-logs", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/audit-logs`, {
          searchParams: filterParams(params as Record<string, unknown>),
        })
        .json<AuditListOutput>(),
    placeholderData: keepPreviousData,
  });
}
