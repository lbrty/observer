import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreateProjectInput,
  ListProjectsOutput,
  ListProjectsParams,
  Project,
  UpdateProjectInput,
} from "@/types/project";

export function useProjects(params: ListProjectsParams = {}) {
  return useQuery({
    queryKey: ["projects", params],
    queryFn: () =>
      api
        .get("admin/projects", {
          searchParams: params as Record<string, string>,
        })
        .json<ListProjectsOutput>(),
    placeholderData: keepPreviousData,
  });
}

export function useProject(id: string) {
  return useQuery({
    queryKey: ["projects", id],
    queryFn: () => api.get(`admin/projects/${id}`).json<Project>(),
    enabled: !!id,
  });
}

export function useCreateProject() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateProjectInput) =>
      api.post("admin/projects", { json: data }).json<Project>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["projects"] }),
  });
}

export function useUpdateProject() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateProjectInput }) =>
      api.patch(`admin/projects/${id}`, { json: data }).json<Project>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["projects"] }),
  });
}
