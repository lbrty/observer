import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreatePlaceInput,
  Place,
  UpdatePlaceInput,
} from "@/types/reference";

export function usePlaces(stateId?: string) {
  return useQuery({
    queryKey: ["places", stateId ?? "all"],
    queryFn: () =>
      api
        .get("admin/places", {
          searchParams: stateId ? { state_id: stateId } : {},
        })
        .json<{ places: Place[] }>(),
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
    onSuccess: () => qc.invalidateQueries({ queryKey: ["places"] }),
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
