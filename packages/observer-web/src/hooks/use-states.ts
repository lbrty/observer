import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreateStateInput,
  State,
  UpdateStateInput,
} from "@/types/reference";

export function useStates(countryId: string) {
  return useQuery({
    queryKey: ["states", countryId],
    queryFn: () =>
      api
        .get("admin/states", { searchParams: { country_id: countryId } })
        .json<{ states: State[] }>(),
    enabled: !!countryId,
  });
}

export function useCreateState() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      countryId,
      data,
    }: {
      countryId: string;
      data: CreateStateInput;
    }) =>
      api
        .post("admin/states", {
          json: data,
          searchParams: { country_id: countryId },
        })
        .json<State>(),
    onSuccess: (_d, vars) =>
      qc.invalidateQueries({ queryKey: ["states", vars.countryId] }),
  });
}

export function useUpdateState() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateStateInput }) =>
      api.patch(`admin/states/${id}`, { json: data }).json<State>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["states"] }),
  });
}

export function useDeleteState() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`admin/states/${id}`).json(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["states"] }),
  });
}
