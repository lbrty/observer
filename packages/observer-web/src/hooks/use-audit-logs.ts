import { keepPreviousData, useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { AuditListOutput, AuditListParams } from "@/types/audit";

function cleanParams(params: AuditListParams): Record<string, string> {
  const out: Record<string, string> = {};
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== "") {
      out[k] = String(v);
    }
  }
  return out;
}

export function useAuditLogs(params: AuditListParams) {
  return useQuery({
    queryKey: ["audit-logs", params],
    queryFn: () =>
      api
        .get("admin/audit-logs", { searchParams: cleanParams(params) })
        .json<AuditListOutput>(),
    placeholderData: keepPreviousData,
  });
}

export function useProjectAuditLogs(projectId: string, params: AuditListParams) {
  return useQuery({
    queryKey: ["project-audit-logs", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/audit-logs`, { searchParams: cleanParams(params) })
        .json<AuditListOutput>(),
    placeholderData: keepPreviousData,
  });
}
