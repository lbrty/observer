import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreateMigrationRecordInput,
  ListMigrationRecordsOutput,
  MigrationRecord,
  UpdateMigrationRecordInput,
} from "@/types/migration-record";

export function useMigrationRecords(projectId: string, personId: string) {
  return useQuery({
    queryKey: ["migration-records", projectId, personId],
    queryFn: () =>
      api
        .get(`projects/${projectId}/people/${personId}/migration-records`)
        .json<ListMigrationRecordsOutput>(),
    enabled: !!projectId && !!personId,
  });
}

export function useMigrationRecord(projectId: string, personId: string, id: string) {
  return useQuery({
    queryKey: ["migration-records", projectId, personId, id],
    queryFn: () =>
      api
        .get(`projects/${projectId}/people/${personId}/migration-records/${id}`)
        .json<MigrationRecord>(),
    enabled: !!projectId && !!personId && !!id,
  });
}

export function useCreateMigrationRecord(projectId: string, personId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateMigrationRecordInput) =>
      api
        .post(`projects/${projectId}/people/${personId}/migration-records`, { json: data })
        .json<MigrationRecord>(),
    onSuccess: () =>
      qc.invalidateQueries({
        queryKey: ["migration-records", projectId, personId],
      }),
  });
}

export function useUpdateMigrationRecord(projectId: string, personId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMigrationRecordInput }) =>
      api
        .patch(`projects/${projectId}/people/${personId}/migration-records/${id}`, { json: data })
        .json<MigrationRecord>(),
    onSuccess: () =>
      qc.invalidateQueries({
        queryKey: ["migration-records", projectId, personId],
      }),
  });
}
