import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { CreateTagInput, ListTagsOutput, Tag } from "@/types/tag";

export function useTags(projectId: string) {
  return useQuery({
    queryKey: ["tags", projectId],
    queryFn: () =>
      api.get(`projects/${projectId}/tags`).json<ListTagsOutput>(),
    enabled: !!projectId,
  });
}

export function useCreateTag(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTagInput) =>
      api.post(`projects/${projectId}/tags`, { json: data }).json<Tag>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tags", projectId] }),
  });
}

export function useDeleteTag(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      api.delete(`projects/${projectId}/tags/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tags", projectId] }),
  });
}
