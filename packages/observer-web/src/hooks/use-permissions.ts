import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  AssignPermissionInput,
  PermissionListOutput,
  ProjectPermission,
  UpdatePermissionInput,
} from "@/types/permission";

export function usePermissions(projectId: string) {
  return useQuery({
    queryKey: ["permissions", projectId],
    queryFn: () => api.get(`admin/projects/${projectId}/permissions`).json<PermissionListOutput>(),
    enabled: !!projectId,
  });
}

export function useAssignPermission() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, data }: { projectId: string; data: AssignPermissionInput }) =>
      api.post(`admin/projects/${projectId}/permissions`, { json: data }).json<ProjectPermission>(),
    onSuccess: (_d, vars) => qc.invalidateQueries({ queryKey: ["permissions", vars.projectId] }),
  });
}

export function useUpdatePermission() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      projectId,
      id,
      data,
    }: {
      projectId: string;
      id: string;
      data: UpdatePermissionInput;
    }) =>
      api
        .patch(`admin/projects/${projectId}/permissions/${id}`, { json: data })
        .json<ProjectPermission>(),
    onSuccess: (_d, vars) => qc.invalidateQueries({ queryKey: ["permissions", vars.projectId] }),
  });
}

export function useRevokePermission() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, id }: { projectId: string; id: string }) =>
      api.delete(`admin/projects/${projectId}/permissions/${id}`).json(),
    onSuccess: (_d, vars) => qc.invalidateQueries({ queryKey: ["permissions", vars.projectId] }),
  });
}
