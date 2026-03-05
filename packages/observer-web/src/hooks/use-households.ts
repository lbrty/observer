import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  AddMemberInput,
  CreateHouseholdInput,
  Household,
  ListHouseholdsOutput,
  ListHouseholdsParams,
  UpdateHouseholdInput,
} from "@/types/household";

export function useHouseholds(
  projectId: string,
  params: ListHouseholdsParams = {},
) {
  return useQuery({
    queryKey: ["households", projectId, params],
    queryFn: () =>
      api
        .get(`projects/${projectId}/households`, {
          searchParams: Object.fromEntries(
            Object.entries(params).filter(([, v]) => v != null),
          ) as Record<string, string>,
        })
        .json<ListHouseholdsOutput>(),
    enabled: !!projectId,
    placeholderData: keepPreviousData,
  });
}

export function useHousehold(projectId: string, id: string) {
  return useQuery({
    queryKey: ["households", projectId, id],
    queryFn: () =>
      api.get(`projects/${projectId}/households/${id}`).json<Household>(),
    enabled: !!projectId && !!id,
  });
}

export function useCreateHousehold(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateHouseholdInput) =>
      api
        .post(`projects/${projectId}/households`, { json: data })
        .json<Household>(),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["households", projectId] }),
  });
}

export function useUpdateHousehold(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateHouseholdInput }) =>
      api
        .patch(`projects/${projectId}/households/${id}`, { json: data })
        .json<Household>(),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["households", projectId] }),
  });
}

export function useDeleteHousehold(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      api.delete(`projects/${projectId}/households/${id}`),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["households", projectId] }),
  });
}

export function useAddHouseholdMember(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      householdId,
      data,
    }: {
      householdId: string;
      data: AddMemberInput;
    }) =>
      api
        .post(`projects/${projectId}/households/${householdId}/members`, {
          json: data,
        })
        .json(),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["households", projectId] }),
  });
}

export function useRemoveHouseholdMember(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      householdId,
      personId,
    }: {
      householdId: string;
      personId: string;
    }) =>
      api.delete(
        `projects/${projectId}/households/${householdId}/members/${personId}`,
      ),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["households", projectId] }),
  });
}
