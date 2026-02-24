import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";
import type {
  Country,
  CreateCountryInput,
  UpdateCountryInput,
} from "@/types/reference";

export function useCountries() {
  return useQuery({
    queryKey: ["countries"],
    queryFn: () => api.get("admin/countries").json<{ countries: Country[] }>(),
  });
}

export function useCreateCountry() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateCountryInput) =>
      api.post("admin/countries", { json: data }).json<Country>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["countries"] }),
  });
}

export function useUpdateCountry() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCountryInput }) =>
      api.patch(`admin/countries/${id}`, { json: data }).json<Country>(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["countries"] }),
  });
}

export function useDeleteCountry() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`admin/countries/${id}`).json(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["countries"] }),
  });
}
