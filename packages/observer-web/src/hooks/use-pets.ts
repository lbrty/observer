import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreatePetInput,
  ListPetsOutput,
  ListPetsParams,
  Pet,
  UpdatePetInput,
} from "@/types/pet";

export function usePets(projectId: string, params: ListPetsParams = {}) {
  return useQuery({
    queryKey: ["pets", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/pets`, {
          searchParams: Object.fromEntries(
            Object.entries(params).filter(([, v]) => v != null),
          ) as Record<string, string>,
        })
        .json<ListPetsOutput>(),
    enabled: !!projectId,
    placeholderData: keepPreviousData,
  });
}

export function usePet(projectId: string, id: string) {
  return useQuery({
    queryKey: ["pets", projectId, id],
    queryFn: () => api.get(`projects/${projectId}/pets/${id}`).json<Pet>(),
    enabled: !!projectId && !!id,
  });
}

export function useCreatePet(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreatePetInput) =>
      api.post(`projects/${projectId}/pets`, { json: data }).json<Pet>(),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["pets", projectId] }),
  });
}

export function useUpdatePet(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePetInput }) =>
      api
        .patch(`projects/${projectId}/pets/${id}`, { json: data })
        .json<Pet>(),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["pets", projectId] }),
  });
}

export function useDeletePet(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      api.delete(`projects/${projectId}/pets/${id}`),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["pets", projectId] }),
  });
}
