import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api";

export function makeReferenceHooks<T extends { id: string }, C, U>(
  resource: string,
) {
  function useList() {
    return useQuery({
      queryKey: [resource],
      queryFn: async () => {
        const res = await api
          .get(`admin/${resource}`)
          .json<Record<string, T[]>>();
        return res[resource];
      },
    });
  }

  function useCreate() {
    const qc = useQueryClient();
    return useMutation({
      mutationFn: (data: C) =>
        api.post(`admin/${resource}`, { json: data }).json<T>(),
      onSuccess: () => qc.invalidateQueries({ queryKey: [resource] }),
    });
  }

  function useUpdate() {
    const qc = useQueryClient();
    return useMutation({
      mutationFn: ({ id, data }: { id: string; data: U }) =>
        api.patch(`admin/${resource}/${id}`, { json: data }).json<T>(),
      onSuccess: () => qc.invalidateQueries({ queryKey: [resource] }),
    });
  }

  function useDelete() {
    const qc = useQueryClient();
    return useMutation({
      mutationFn: (id: string) => api.delete(`admin/${resource}/${id}`).json(),
      onSuccess: () => qc.invalidateQueries({ queryKey: [resource] }),
    });
  }

  return { useList, useCreate, useUpdate, useDelete };
}
