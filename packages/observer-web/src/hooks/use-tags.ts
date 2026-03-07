import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type { CreateTagInput, ListTagsOutput, Tag, UpdateTagInput } from "@/types/tag";

export function useTags(projectId: string) {
  return useQuery({
    queryKey: ["tags", projectId],
    queryFn: () => api.get(`projects/${projectId}/tags`).json<ListTagsOutput>(),
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

export function useUpdateTag(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTagInput }) =>
      api.patch(`projects/${projectId}/tags/${id}`, { json: data }).json<Tag>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tags", projectId] }),
  });
}

export function useDeleteTag(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`projects/${projectId}/tags/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tags", projectId] }),
  });
}

export function usePersonTags(projectId: string, personId: string) {
  return useQuery({
    queryKey: ["personTags", projectId, personId],
    queryFn: () =>
      api.get(`projects/${projectId}/people/${personId}/tags`).json<{ tag_ids: string[] }>(),
    enabled: !!projectId && !!personId,
  });
}

export function useReplacePersonTags(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ personId, ids }: { personId: string; ids: string[] }) =>
      api
        .put(`projects/${projectId}/people/${personId}/tags`, { json: { ids } })
        .json<{ tag_ids: string[] }>(),
    onSuccess: (_data, { personId }) =>
      qc.invalidateQueries({ queryKey: ["personTags", projectId, personId] }),
  });
}

export function usePetTags(projectId: string, petId: string) {
  return useQuery({
    queryKey: ["petTags", projectId, petId],
    queryFn: () =>
      api.get(`projects/${projectId}/pets/${petId}/tags`).json<{ tag_ids: string[] }>(),
    enabled: !!projectId && !!petId,
  });
}

export function useReplacePetTags(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ petId, ids }: { petId: string; ids: string[] }) =>
      api
        .put(`projects/${projectId}/pets/${petId}/tags`, { json: { ids } })
        .json<{ tag_ids: string[] }>(),
    onSuccess: (_data, { petId }) =>
      qc.invalidateQueries({ queryKey: ["petTags", projectId, petId] }),
  });
}
