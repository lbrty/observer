import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  CreateOfficeInput,
  Office,
  UpdateOfficeInput,
} from "@/types/reference";

export function useOffices() {
  return useQuery({
    queryKey: ["offices"],
    queryFn: () => api.get("admin/offices").json<{ offices: Office[] }>(),
  });
}

export function useCreateOffice() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateOfficeInput) =>
      api.post("admin/offices", { json: data }).json<Office>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["offices"] }),
  });
}

export function useUpdateOffice() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateOfficeInput }) =>
      api.patch(`admin/offices/${id}`, { json: data }).json<Office>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["offices"] }),
  });
}

export function useDeleteOffice() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`admin/offices/${id}`).json(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["offices"] }),
  });
}
