import { keepPreviousData, useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import { filterParams } from "@/lib/params";
import type { AuditListOutput, AuditListParams } from "@/types/audit";

export function useAuditLogs(params: AuditListParams) {
  return useQuery({
    queryKey: ["audit-logs", params],
    queryFn: ({ signal }) =>
      api
        .get("admin/audit-logs", {
          signal,
          searchParams: filterParams(params as Record<string, unknown>),
        })
        .json<AuditListOutput>(),
    placeholderData: keepPreviousData,
  });
}

export function useProjectAuditLogs(projectId: string, params: AuditListParams) {
  return useQuery({
    queryKey: ["project-audit-logs", projectId, params],
    queryFn: ({ signal }) =>
      api
        .get(`projects/${projectId}/audit-logs`, {
          signal,
          searchParams: filterParams(params as Record<string, unknown>),
        })
        .json<AuditListOutput>(),
    placeholderData: keepPreviousData,
  });
}
