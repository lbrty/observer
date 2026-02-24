import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreatePlaceInput,
  Place,
  UpdatePlaceInput,
} from "@/types/reference";

export function usePlaces(stateId: string) {
  return useQuery({
    queryKey: ["places", stateId],
    queryFn: () =>
      api
        .get("admin/places", { searchParams: { state_id: stateId } })
        .json<{ places: Place[] }>(),
    enabled: !!stateId,
  });
}

export function useCreatePlace() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      stateId,
      data,
    }: {
      stateId: string;
      data: CreatePlaceInput;
    }) =>
      api
        .post("admin/places", {
          json: data,
          searchParams: { state_id: stateId },
        })
        .json<Place>(),
    onSuccess: (_d, vars) =>
      qc.invalidateQueries({ queryKey: ["places", vars.stateId] }),
  });
}

export function useUpdatePlace() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePlaceInput }) =>
      api.patch(`admin/places/${id}`, { json: data }).json<Place>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["places"] }),
  });
}

export function useDeletePlace() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`admin/places/${id}`).json(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["places"] }),
  });
}
