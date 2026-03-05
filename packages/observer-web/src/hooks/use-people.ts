import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreatePersonInput,
  ListPeopleOutput,
  ListPeopleParams,
  Person,
  UpdatePersonInput,
} from "@/types/person";

export function usePeople(projectId: string, params: ListPeopleParams = {}) {
  return useQuery({
    queryKey: ["people", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/people`, {
          searchParams: Object.fromEntries(
            Object.entries(params).filter(([, v]) => v != null),
          ) as Record<string, string>,
        })
        .json<ListPeopleOutput>(),
    enabled: !!projectId,
    placeholderData: keepPreviousData,
  });
}

export function usePerson(projectId: string, personId: string) {
  return useQuery({
    queryKey: ["people", projectId, personId],
    queryFn: () => api.get(`projects/${projectId}/people/${personId}`).json<Person>(),
    enabled: !!projectId && !!personId,
  });
}

export function useCreatePerson(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreatePersonInput) =>
      api.post(`projects/${projectId}/people`, { json: data }).json<Person>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["people", projectId] }),
  });
}

export function useUpdatePerson(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ personId, data }: { personId: string; data: UpdatePersonInput }) =>
      api.patch(`projects/${projectId}/people/${personId}`, { json: data }).json<Person>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["people", projectId] }),
  });
}

export function useDeletePerson(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (personId: string) => api.delete(`projects/${projectId}/people/${personId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["people", projectId] }),
  });
}
