import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreateSupportRecordInput,
  ListSupportRecordsOutput,
  ListSupportRecordsParams,
  SupportRecord,
  UpdateSupportRecordInput,
} from "@/types/support-record";

export function useSupportRecords(projectId: string, params: ListSupportRecordsParams = {}) {
  return useQuery({
    queryKey: ["support-records", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/support-records`, {
          searchParams: Object.fromEntries(
            Object.entries(params).filter(([, v]) => v != null),
          ) as Record<string, string>,
        })
        .json<ListSupportRecordsOutput>(),
    enabled: !!projectId,
    placeholderData: keepPreviousData,
  });
}

export function useSupportRecord(projectId: string, id: string) {
  return useQuery({
    queryKey: ["support-records", projectId, id],
    queryFn: () => api.get(`projects/${projectId}/support-records/${id}`).json<SupportRecord>(),
    enabled: !!projectId && !!id,
  });
}

export function useCreateSupportRecord(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateSupportRecordInput) =>
      api.post(`projects/${projectId}/support-records`, { json: data }).json<SupportRecord>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["support-records", projectId] }),
  });
}

export function useUpdateSupportRecord(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateSupportRecordInput }) =>
      api
        .patch(`projects/${projectId}/support-records/${id}`, { json: data })
        .json<SupportRecord>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["support-records", projectId] }),
  });
}

export function useDeleteSupportRecord(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`projects/${projectId}/support-records/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["support-records", projectId] }),
  });
}
